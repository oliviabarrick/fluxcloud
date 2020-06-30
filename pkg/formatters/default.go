package formatters

import (
	"bytes"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/msg"
	"github.com/weaveworks/flux"
	fluxevent "github.com/weaveworks/flux/event"
	"github.com/weaveworks/flux/update"
)

const (
	titleTemplate = `Applied flux changes to cluster`
	bodyTemplate  = `
Event: {{ .EventString }}
{{ if and (ne .EventType "commit") (gt (len .Commits) 0) }}Commits:
{{ range .Commits }}
* {{ call $.FormatLink (print $.VCSLink "/commit/" .Revision) (truncate .Revision 7) }}: {{ .Message }}
{{end}}{{end}}
{{ if (gt (len .EventServiceIDs) 0) }}Resources updated:
{{ range .EventServiceIDs }}
* {{ . }}
{{ end }}{{ end }}
{{ if gt (len .Errors) 0 }}Errors:
{{ range .Errors }}
Resource {{ .ID }}, file: {{ .Path }}:

> {{ .Error }}
{{ end }}{{ end }}
`
	commitTemplate = `{{ .VCSLink }}/commit/{{ .Commit }}`
)

// The default formatter formats a message for a chat webhook
type DefaultFormatter struct {
	config         config.Config
	vcsLink        string
	bodyTemplate   string
	titleTemplate  string
	commitTemplate string
}

type tplValues struct {
	VCSLink            string
	EventID            fluxevent.EventID
	EventServiceIDs    []flux.ResourceID
	EventChangedImages []string
	EventResult        update.Result
	EventType          string
	EventStartedAt     time.Time
	EventEndedAt       time.Time
	EventLogLevel      string
	EventMessage       string
	EventString        string
	Commits            []fluxevent.Commit
	Errors             []fluxevent.ResourceError
	FormatLink         func(string, string) string
}

type commitTemplateValues struct {
	VCSLink string
	Commit  string
}

var (
	tplFuncMap = template.FuncMap{
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
		"trim": func(input string) string {
			return strings.TrimSpace(input)
		},
		"contains": func(input, substr string) bool {
			return strings.Contains(input, substr)
		},
		"truncate": func(s string, max int) string {
			var numRunes = 0
			for i := range s {
				numRunes++
				if numRunes > max {
					return s[:i]
				}
			}
			return s
		},
	}
)

// Create a DefaultFormatter
func NewDefaultFormatter(config config.Config) (*DefaultFormatter, error) {
	vcsLink, err := config.Required("github_url")
	if err != nil {
		return nil, err
	}

	bodyTemplate := config.Optional("body_template", bodyTemplate)
	titleTemplate := config.Optional("title_template", titleTemplate)
	commitTemplate := config.Optional("commit_template", commitTemplate)

	if err := checkTemplate(bodyTemplate); err != nil {
		log.Println(bodyTemplate)
		return nil, err
	}

	if err := checkTemplate(titleTemplate); err != nil {
		log.Println(titleTemplate)
		return nil, err
	}

	if err := checkTemplate(commitTemplate); err != nil {
		log.Println(commitTemplate)
		return nil, err
	}

	return &DefaultFormatter{
		config:         config,
		vcsLink:        vcsLink,
		bodyTemplate:   bodyTemplate,
		titleTemplate:  titleTemplate,
		commitTemplate: commitTemplate,
	}, nil
}

// Format plaintext message for an exporter for Flux event
func (d DefaultFormatter) FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message {
	if len(event.ServiceIDs) == 0 {
		return msg.Message{}
	}

	values := &tplValues{
		VCSLink:            d.vcsLink,
		EventID:            event.ID,
		EventServiceIDs:    event.ServiceIDs,
		EventChangedImages: getChangedImages(event.Metadata),
		EventResult:        getResult(event.Metadata),
		EventType:          event.Type,
		EventStartedAt:     event.StartedAt,
		EventEndedAt:       event.EndedAt,
		EventLogLevel:      event.LogLevel,
		EventMessage:       event.Message,
		EventString:        event.String(),
		Commits:            getCommits(event.Metadata),
		Errors:             getErrors(event.Metadata),
		FormatLink: func(link, text string) string {
			return exporter.FormatLink(link, text)
		},
	}

	nl := exporter.NewLine()

	message := msg.Message{
		TitleLink: d.vcsLink,
		Title:     execTemplate(d.titleTemplate, values, nl),
		Body:      execTemplate(d.bodyTemplate, values, nl),
		Type:      event.Type,
		Event:     event,
	}

	if message.Title == "" || message.Body == "" {
		return msg.Message{}
	}

	commits := getCommits(event.Metadata)
	if len(commits) > 0 {
		message.TitleLink = execTemplate(d.commitTemplate, &commitTemplateValues{
			VCSLink: d.vcsLink,
			Commit:  commits[0].Revision,
		}, nl)
	}

	return message
}

func checkTemplate(tpl string) error {
	bodyTpl := template.New("tpl").Funcs(tplFuncMap)
	_, err := bodyTpl.Parse(tpl)
	return err
}

func execTemplate(tpl string, values interface{}, nl string) string {
	bodyBytes := &bytes.Buffer{}

	var err error
	bodyTpl := template.New("tpl").Funcs(tplFuncMap)
	bodyTpl, err = bodyTpl.Parse(tpl)
	if err != nil {
		log.Panicln("could not parse template")
	}

	if err := bodyTpl.Execute(bodyBytes, values); err != nil {
		log.Println("could not execute template:", err)
		return ""
	}

	body := bodyBytes.String()
	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	body = strings.Join(lines, nl)

	return body
}

func getCommits(meta fluxevent.EventMetadata) []fluxevent.Commit {
	switch v := meta.(type) {
	case *fluxevent.CommitEventMetadata:
		return []fluxevent.Commit{
			fluxevent.Commit{
				Revision: v.Revision,
			},
		}
	case *fluxevent.SyncEventMetadata:
		return v.Commits
	default:
		return []fluxevent.Commit{}
	}
}

func getResult(meta fluxevent.EventMetadata) update.Result {
	switch v := meta.(type) {
	case *fluxevent.AutoReleaseEventMetadata:
		return v.Result
	case *fluxevent.ReleaseEventMetadata:
		return v.Result
	default:
		return update.Result{}
	}
}

func getChangedImages(meta fluxevent.EventMetadata) []string {
	switch v := meta.(type) {
	case *fluxevent.AutoReleaseEventMetadata:
		return v.Result.ChangedImages()
	case *fluxevent.ReleaseEventMetadata:
		return v.Result.ChangedImages()
	default:
		return []string{}
	}
}

func getErrors(meta fluxevent.EventMetadata) []fluxevent.ResourceError {
	switch v := meta.(type) {
	case *fluxevent.SyncEventMetadata:
		return v.Errors
	default:
		return []fluxevent.ResourceError{}
	}
}
