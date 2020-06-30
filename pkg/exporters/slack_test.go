package exporters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/weaveworks/flux"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"
	fluxevent "github.com/weaveworks/flux/event"
)

var testSlack = Slack{
	Channels: []SlackChannel{
		SlackChannel{"#channel", "*"},
		SlackChannel{"#namespace", "namespace"},
	},
	IconEmoji: ":weave:",
	Username:  "My Username",
}

func TestSlackDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")
	config.Set("slack_channel", "#mychannel")

	slack, err := NewSlack(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://myslack/", slack.Url)
	assert.Equal(t, []SlackChannel{SlackChannel{"#mychannel", "*"}}, slack.Channels)
	assert.Equal(t, "Flux Deployer", slack.Username)
	assert.Equal(t, ":star-struck:", slack.IconEmoji)
}

func TestSlackOverrides(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")
	config.Set("slack_channel", "#mychannel=namespace")
	config.Set("slack_username", "my user")
	config.Set("slack_icon_emoji", ":weave:")
	config.Set("slack_token", "mytoken")

	slack, err := NewSlack(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://myslack/", slack.Url)
	assert.Equal(t, []SlackChannel{SlackChannel{"#mychannel", "namespace"}}, slack.Channels)
	assert.Equal(t, "my user", slack.Username)
	assert.Equal(t, ":weave:", slack.IconEmoji)
	assert.Equal(t, "mytoken", slack.Token)
}

func TestSlackChannel(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("slack_url", "https://myslack/")
	config.Set("slack_channel", "#mychannel=*, #namespace=namespace")

	slack, err := NewSlack(config)
	assert.Nil(t, err)

	assert.Equal(t, []SlackChannel{
		SlackChannel{"#mychannel", "*"},
		SlackChannel{"#namespace", "namespace"},
	}, slack.Channels)
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
	defaultResourceID, _ := flux.ParseResourceID("default:resource/name")
	nsResourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			ServiceIDs: []flux.ResourceID{
				defaultResourceID,
				nsResourceID,
			},
		},
	}

	slackMessages := testSlack.NewSlackMessage(message)
	assert.Len(t, slackMessages, 2)

	assert.Equal(t, "#channel", slackMessages[0].Channel)
	assert.Equal(t, "#namespace", slackMessages[1].Channel)
	assert.Equal(t, testSlack.IconEmoji, slackMessages[0].IconEmoji)
	assert.Equal(t, testSlack.Username, slackMessages[0].Username)

	attach := slackMessages[0].Attachments[0]
	assert.Equal(t, "#4286f4", attach.Color)
	assert.Equal(t, message.TitleLink, attach.TitleLink)
	assert.Equal(t, message.Title, attach.Title)
}

func TestSlackSend(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			ServiceIDs: []flux.ResourceID{
				resourceID,
			},
		},
	}

	slackMessage := SlackMessage{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&slackMessage)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	testSlack.Url = ts.URL

	err := testSlack.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Contains(t, testSlack.NewSlackMessage(message), slackMessage)
}

func TestSlackSendNon200(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		Event: fluxevent.Event{
			ServiceIDs: []flux.ResourceID{
				resourceID,
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	testSlack.Url = ts.URL

	err := testSlack.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestSlackSendHTTPError(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		Event: fluxevent.Event{
			ServiceIDs: []flux.ResourceID{
				resourceID,
			},
		},
	}

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	testSlack.Url = ts.URL

	err := testSlack.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestSlackName(t *testing.T) {

	slack := Slack{}
	assert.Equal(t, "Slack", slack.Name())
}

func TestSlackImplementsExporter(t *testing.T) {
	_ = Exporter(&Slack{})
}

func TestSlackSendAuthToken(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			ServiceIDs: []flux.ResourceID{
				resourceID,
			},
		},
	}

	authHeader := ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	testSlack.Url = ts.URL
	testSlack.Token = "mytoken"

	err := testSlack.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Bearer %s", testSlack.Token), authHeader)
}
