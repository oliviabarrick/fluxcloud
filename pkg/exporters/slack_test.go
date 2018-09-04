package exporters

import (
	"encoding/json"
	"fmt"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlackDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")
	config.Set("slack_channel", "#mychannel")

	slack, err := NewSlack(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://myslack/", slack.Url)
	assert.Equal(t, "#mychannel", slack.Channel)
	assert.Equal(t, "Flux Deployer", slack.Username)
	assert.Equal(t, ":star-struck:", slack.IconEmoji)
}

func TestSlackOverrides(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")
	config.Set("slack_channel", "#mychannel")
	config.Set("slack_username", "my user")
	config.Set("slack_icon_emoji", ":weave:")

	slack, err := NewSlack(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://myslack/", slack.Url)
	assert.Equal(t, "#mychannel", slack.Channel)
	assert.Equal(t, "my user", slack.Username)
	assert.Equal(t, ":weave:", slack.IconEmoji)
}

func TestSlackMissingChannel(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")

	_, err := NewSlack(config)
	assert.NotNil(t, err)
}

func TestSlackMissingSlackUrl(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_channel", "#channel")

	_, err := NewSlack(config)
	assert.NotNil(t, err)
}

func TestSlackNewLine(t *testing.T) {
	slack := Slack{}
	assert.Equal(t, "\n", slack.NewLine())
}

func TestSlackFormatLink(t *testing.T) {
	slack := Slack{}
	assert.Equal(t, "<https://google.com|title>", slack.FormatLink("https://google.com", "title"))
}

func TestNewSlackMessage(t *testing.T) {
	slack := Slack{
		Channel:   "#mychannel",
		IconEmoji: ":weave:",
		Username:  "My Username",
	}

	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	slackMessage := slack.NewSlackMessage(message)
	assert.Equal(t, slack.Channel, slackMessage.Channel)
	assert.Equal(t, slack.IconEmoji, slackMessage.IconEmoji)
	assert.Equal(t, slack.Username, slackMessage.Username)

	attach := slackMessage.Attachments[0]
	assert.Equal(t, "#4286f4", attach.Color)
	assert.Equal(t, message.TitleLink, attach.TitleLink)
	assert.Equal(t, message.Title, attach.Title)
}

func TestSlackSend(t *testing.T) {
	slack := Slack{
		Channel:   "#mychannel",
		IconEmoji: ":weave:",
		Username:  "My Username",
	}

	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	slackMessage := SlackMessage{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&slackMessage)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	slack.Url = ts.URL

	err := slack.Send(&http.Client{}, message)
	assert.Nil(t, err)
	assert.Equal(t, slack.NewSlackMessage(message), slackMessage)
}

func TestSlackSendNon200(t *testing.T) {
	slack := Slack{}
	message := msg.Message{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	slack.Url = ts.URL

	err := slack.Send(&http.Client{}, message)
	assert.NotNil(t, err)
}

func TestSlackSendHTTPError(t *testing.T) {
	slack := Slack{}
	message := msg.Message{}

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	slack.Url = ts.URL

	err := slack.Send(&http.Client{}, message)
	assert.NotNil(t, err)
}

func TestSlackName(t *testing.T) {

	slack := Slack{}
	assert.Equal(t, "Slack", slack.Name())
}

func TestSlackImplementsExporter(t *testing.T) {
	_ = Exporter(&Slack{})
}
