package formatters

import (
	"fmt"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	fluxevent "github.com/weaveworks/flux/event"
)

// The default formatter formats a message for a chat webhook
type DefaultFormatter struct {
	config  config.Config
	vcsLink string
}

// Create a DefaultFormatter
func NewDefaultFormatter(config config.Config) (*DefaultFormatter, error) {
	vcsLink, err := config.Required("github_url")
	if err != nil {
		return nil, err
	}

	return &DefaultFormatter{
		config:  config,
		vcsLink: vcsLink,
	}, nil
}

// Format plaintext message for an exporter for Flux event
func (d DefaultFormatter) FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message {
	newLine := exporter.NewLine()
	body := fmt.Sprintf("Event: %s%s", event.String(), newLine)
	errorBody := ""
	commit_id := ""

	if len(event.ServiceIDs) == 0 {
		return msg.Message{}
	}

	titleLink := d.vcsLink

	switch event.Type {
	case fluxevent.EventSync:
		metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
		commit_id = metadata.Commits[0].Revision

		body += fmt.Sprintf("Commits: %s", newLine)
		for _, commit := range metadata.Commits {
			vcsLinkWithCommit := d.vcsLink + "/commit/" + commit.Revision
			link := exporter.FormatLink(vcsLinkWithCommit, commit.Revision[:7])
			body += fmt.Sprintf("%s* %s: %s", newLine, link, commit.Message)
		}

		body += newLine

		if len(metadata.Errors) > 0 {
			errorBody += "Errors:\n"

			for _, err := range metadata.Errors {
				errorBody = fmt.Sprintf("%s\nResource %s, file: %s:\n\n> %s\n", errorBody, err.ID, err.Path, err.Error)
			}
		}
	case fluxevent.EventCommit:
		metadata := event.Metadata.(*fluxevent.CommitEventMetadata)
		commit_id = metadata.Revision
	default:
	}

	body += fmt.Sprintf("%sResources updated:%s", newLine, newLine)
	for _, serviceId := range event.ServiceIDStrings() {
		body += fmt.Sprintf("%s* %s", newLine, serviceId)
	}

	if len(errorBody) != 0 {
		body += fmt.Sprintf("\n\n%s", errorBody)
	}

	if commit_id != "" {
		titleLink = d.vcsLink + "/commit/" + commit_id
	}

	return msg.Message{
		TitleLink: titleLink,
		Title:     "Applied flux changes to cluster",
		Body:      body,
		Type:      event.Type,
		Event:     event,
	}
}
