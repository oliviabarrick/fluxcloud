package formatters

import (
	"bytes"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"github.com/weaveworks/flux"
	fluxevent "github.com/weaveworks/flux/event"
)

const (
	titleTemplate = `Applied flux changes to cluster`
	bodyTemplate  = `
Event: {{ .EventString }}
{{ if and (ne .EventType "commit") (gt (len .Commits) 0) }}Commits:
{{ range .Commits }}
* {{ call $.FormatLink (print $.VCSLink "/commit/" .Revision) (printf "%.7s" .Revision) }}: {{ .Message }}
{{end}}{{end}}
{{ if (gt (len .EventServiceIDs) 0) }}Resources updated:
{{ range .EventServiceIDs }}
* {{ . }}
{{ end }}{{ end }}
{{ if gt (len .Errors) 0 }}Errors:
{{ range .Errors }}
Resource {{ .ID }}, file: {{ .Path }}:

> {{ call $.FormatError .Error }}
{{ end }}{{ end }}
`
)

// The default formatter formats a message for a chat webhook
type DefaultFormatter struct {
	config        config.Config
	vcsLink       string
	bodyTemplate  string
	titleTemplate string
}

// Create a DefaultFormatter
func NewDefaultFormatter(config config.Config) (*DefaultFormatter, error) {
	vcsLink, err := config.Required("github_url")
	if err != nil {
		return nil, err
	}

	return &DefaultFormatter{
		config:        config,
		vcsLink:       vcsLink,
		bodyTemplate:  config.Optional("body_template", bodyTemplate),
		titleTemplate: config.Optional("title_template", titleTemplate),
	}, nil
}

// Format plaintext message for an exporter for Flux event
func (d DefaultFormatter) FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message {
	if len(event.ServiceIDs) == 0 {
		return msg.Message{}
	}

	values := struct {
		VCSLink         string
		EventID         fluxevent.EventID
		EventServiceIDs []flux.ResourceID
		EventType       template.HTML
		EventStartedAt  time.Time
		EventEndedAt    time.Time
		EventLogLevel   template.HTML
		EventMessage    template.HTML
		EventString     template.HTML
		Commits         []fluxevent.Commit
		Errors          []fluxevent.ResourceError
		FormatError     func(string) template.HTML
		FormatLink      func(string, string) template.HTML
	}{
		VCSLink:         d.vcsLink,
		EventID:         event.ID,
		EventServiceIDs: event.ServiceIDs,
		EventType:       template.HTML(event.Type),
		EventStartedAt:  event.StartedAt,
		EventEndedAt:    event.EndedAt,
		EventLogLevel:   template.HTML(event.LogLevel),
		EventMessage:    template.HTML(event.Message),
		EventString:     template.HTML(event.String()),
		Commits:         getCommits(event.Metadata),
		Errors:          getErrors(event.Metadata),
		FormatError: func(text string) template.HTML {
			return template.HTML(text)
		},
		FormatLink: func(link, text string) template.HTML {
			return template.HTML(exporter.FormatLink(link, text))
		},
	}

	nl := exporter.NewLine()

	message := msg.Message{
		TitleLink: d.vcsLink,
		Title:     parseTemplate(d.titleTemplate, values, nl),
		Body:      parseTemplate(d.bodyTemplate, values, nl),
		Type:      event.Type,
		Event:     event,
	}

	if message.Title == "" || message.Body == "" {
		return msg.Message{}
	}

	commits := getCommits(event.Metadata)
	if len(commits) > 0 {
		message.TitleLink = d.vcsLink + "/commit/" + commits[0].Revision
	}

	return message
}

func parseTemplate(tpl string, values interface{}, nl string) string {
	bodyBytes := &bytes.Buffer{}
	bodyTpl := template.New("")
	bodyTpl, _ = bodyTpl.Parse(tpl)
	if err := bodyTpl.Execute(bodyBytes, values); err != nil {
		log.Println("could not execute body template:", err)
		return ""
	}

	body := bodyBytes.String()
	body = strings.TrimSpace(body)
	body = strings.Replace(body, "\n", nl, -1)

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

func getErrors(meta fluxevent.EventMetadata) []fluxevent.ResourceError {
	switch v := meta.(type) {
	case *fluxevent.SyncEventMetadata:
		return v.Errors
	default:
		return []fluxevent.ResourceError{}
	}
}
