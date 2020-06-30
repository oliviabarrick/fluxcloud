package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"
)

// The MSTeams exporter sends Flux events to a Microsoft Teams channel via a webhook.
type MSTeams struct {
	Url string
}

// Represents a MS Teams message sent to the API
type MSTeamsMessage struct {
	Context    string          `json:"@context"`
	Type       string          `json:"@type"`
	ThemeColor string          `json:"themeColor"`
	Title      string          `json:"title"`
	Text       string          `json:"text"`
	Actions    []MSTeamsAction `json:"potentialAction"`
}

// Represents an action embedded in a MS Teams message
type MSTeamsAction struct {
	Type    string                `json:"@type"`
	Name    string                `json:"name"`
	Targets []MSTeamsActionTarget `json:"targets"`
}

// Represents a slack channel and the Kubernetes namespace linked to it
type MSTeamsActionTarget struct {
	OS  string `json:"os"`
	URI string `json:"uri"`
}

// Initialize a new MSTeams instance
func NewMSTeams(config config.Config) (*MSTeams, error) {
	var err error
	t := MSTeams{}

	t.Url, err = config.Required("msteams_url")
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Send a MSTeamsMessage to MS Teams
func (s *MSTeams) Send(ctx context.Context, client *http.Client, message msg.Message) error {
	msTeamsMessage := s.NewMSTeamsMessage(message)
	fmt.Println(msTeamsMessage)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(msTeamsMessage)
	if err != nil {
		log.Print("Error when encoding MS Teams message in JSON:", err)
		return err
	}

	req, _ := http.NewRequest("POST", s.Url, b)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	log.Print(string(b.Bytes()))
	res, err := client.Do(req)
	if err != nil {
		log.Print("Could not post to MS Teams:", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to MS Teams, status: ", res.Status)
		return fmt.Errorf("Could not post to MS Teams, status: %d", res.StatusCode)
	}

	return nil
}

// Return the new line character for MS Teams messages
func (s *MSTeams) NewLine() string {
	return "\n"
}

// Return a formatted link for MS Teams.
func (s *MSTeams) FormatLink(link string, name string) string {
	return fmt.Sprintf("[%s](%s)", name, link)
}

// Convert a flux event into a MS Teams message
func (s *MSTeams) NewMSTeamsMessage(message msg.Message) MSTeamsMessage {
	result := MSTeamsMessage{
		Context:    "https://schema.org/extensions",
		Type:       "MessageCard",
		ThemeColor: "4286f4",
		Title:      message.Title,
		Text:       message.Body,
		Actions: []MSTeamsAction{
			MSTeamsAction{
				Type: "OpenUri",
				Name: "Open related GIT repository",
				Targets: []MSTeamsActionTarget{
					MSTeamsActionTarget{
						OS:  "default",
						URI: message.TitleLink,
					},
				},
			},
		},
	}
	return result
}

// Return the name of the exporter.
func (s *MSTeams) Name() string {
	return "MS Teams"
}
