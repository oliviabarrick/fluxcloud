package exporters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"log"
	"net/http"
)

// The Slack exporter sends Flux events to a Slack channel via a webhook.
type Slack struct {
	Url       string
	Username  string
	Channel   string
	IconEmoji string
}

// Represents a slack message sent to the API
type SlackMessage struct {
	Channel     string            `json:"channel"`
	IconEmoji   string            `json:"icon_emoji"`
	Username    string            `json:"username"`
	Attachments []SlackAttachment `json:"attachments"`
}

// Represents a section of a slack message that is sent to the API
type SlackAttachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}

// Initialize a new Slack instance
func NewSlack(config config.Config) (*Slack, error) {
	var err error
	s := Slack{}

	s.Url, err = config.Required("slack_url")
	if err != nil {
		return nil, err
	}

	s.Channel, err = config.Required("slack_channel")
	if err != nil {
		return nil, err
	}

	s.Username = config.Optional("slack_username", "Flux Deployer")
	s.IconEmoji = config.Optional("slack_icon_emoji", ":star-struck:")

	return &s, nil
}

// Send a SlackMessage to Slack
func (s *Slack) Send(client *http.Client, message msg.Message) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(s.NewSlackMessage(message))
	if err != nil {
		log.Print("Could encode message to slack:", err)
		return err
	}

	log.Print(string(b.Bytes()))
	res, err := client.Post(s.Url, "application/json", b)
	if err != nil {
		log.Print("Could not post to slack:", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to slack, status: ", res.Status)
		return errors.New(fmt.Sprintf("Could not post to slack, status: %d", res.StatusCode))
	}

	return nil
}

// Return the new line character for Slack messages
func (s *Slack) NewLine() string {
	return "\n"
}

// Return a formatted link for Slack.
func (s *Slack) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

// Convert a flux event into a Slack message
func (s *Slack) NewSlackMessage(message msg.Message) SlackMessage {
	return SlackMessage{
		Channel:   s.Channel,
		IconEmoji: s.IconEmoji,
		Username:  s.Username,
		Attachments: []SlackAttachment{
			SlackAttachment{
				Color:     "#4286f4",
				TitleLink: message.TitleLink,
				Title:     message.Title,
				Text:      message.Body,
			},
		},
	}
}

// Return the name of the exporter.
func (s *Slack) Name() string {
	return "Slack"
}
