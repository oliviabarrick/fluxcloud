package exporters

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"

	"github.com/zorkian/go-datadog-api"
)

// Datadog wraps the Datadog Client
type Datadog struct {
	datadogClient *datadog.Client
}

// DatadogEvent represents a Datadog Event
type DatadogEvent struct {
	Title      string   `json:"Title"`
	Text       string   `json:"Text"`
	Tags       []string `json:"Tags"`
	SourceType string   `json:"SourceType"`
}

// NewDatadog inits the client
func NewDatadog(config config.Config) (*Datadog, error) {
	var err error
	d := Datadog{}

	appKey, err := config.Required("datadog_app_key")
	if err != nil {
		return nil, err
	}
	apiKey, err := config.Required("datadog_api_key")
	if err != nil {
		return nil, err
	}
	d.datadogClient = datadog.NewClient(apiKey, appKey)
	return &d, nil
}

// NewDatadogEvent Convert a flux event into a Datadog Event
func (d *Datadog) NewDatadogEvent(message msg.Message) []DatadogEvent {
	var events []DatadogEvent
	// ServiceID should be something like ns/resource/resourcename
	for _, event := range message.Event.ServiceIDs {
		tags := []string{"application:flux", "fluxEventType:" + message.Event.Type}
		namespace, kind, name := event.Components()
		tags = append(tags, "fluxnamespace:"+namespace)
		tags = append(tags, "fluxkind:"+kind)
		tags = append(tags, "fluxresourcename:"+name)
		additionalTags, exists := os.LookupEnv("DATADOG_ADDITIONAL_TAGS")
		if exists {
			for _, tag := range strings.Split(additionalTags, ",") {
				tags = append(tags, tag)
			}
		}
		events = append(events, DatadogEvent{message.Title, message.Body, tags, "API"})
	}
	return events
}

// Send uses the datadog api client to post the event
func (d *Datadog) Send(c context.Context, client *http.Client, message msg.Message) error {
	for _, event := range d.NewDatadogEvent(message) {
		_, err := d.datadogClient.PostEvent(&datadog.Event{Title: &event.Title,
			Text:       &event.Text,
			Tags:       event.Tags,
			SourceType: &event.SourceType,
		})
		if err != nil {
			log.Print("Could not create datadog event")
			return err
		}
	}
	return nil
}

// NewLine returns the new line character
func (d *Datadog) NewLine() string {
	return "\n"
}

// FormatLink format the git link
func (d *Datadog) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

// Name returns the exporter name
func (d *Datadog) Name() string {
	return "Datadog Events"
}
