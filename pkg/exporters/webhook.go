package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"
)

// The Webhook exporter sends Flux events to a Webhook channel via a webhook.
type Webhook struct {
	Url string
}

// Initialize a new Webhook instance
func NewWebhook(config config.Config) (*Webhook, error) {
	var err error
	s := Webhook{}

	s.Url, err = config.Required("Webhook_url")
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Send a WebhookMessage to Webhook
func (s *Webhook) Send(c context.Context, client *http.Client, message msg.Message) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(message)
	if err != nil {
		log.Print("Could encode message to Webhook:", err)
		return err
	}

	log.Print(string(b.Bytes()))

	req, _ := http.NewRequest("POST", s.Url, b)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(c)

	res, err := client.Do(req)
	if err != nil {
		log.Print("Could not post to Webhook:", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to Webhook, status: ", res.Status)
		return errors.New(fmt.Sprintf("Could not post to Webhook, status: %d", res.StatusCode))
	}

	return nil
}

// Return the new line character for Webhook messages
func (s *Webhook) NewLine() string {
	return "\n"
}

// Return a formatted link for Webhook.
func (s *Webhook) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

// Return the name of the exporter.
func (s *Webhook) Name() string {
	return "Webhook"
}
