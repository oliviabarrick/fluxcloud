package exporters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/weaveworks/flux"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"github.com/stretchr/testify/assert"
	fluxevent "github.com/weaveworks/flux/event"
)

var testMattermost = Mattermost{
	Channels: []MattermostChannel{
		MattermostChannel{"#channel", "*"},
		MattermostChannel{"#namespace", "namespace"},
	},
	IconURL:  "https://user-images.githubusercontent.com/27962005/35868977-0d5f85f6-0b2c-11e8-9fa8-8e4eaf35161a.png",
	Username: "My Username",
}

func TestMattermostDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("mattermost_url", "https://my_mattermost/")
	config.Set("mattermost_channel", "#mychannel")

	mattermost, err := NewMattermost(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://my_mattermost/", mattermost.Url)
	assert.Equal(t, []MattermostChannel{MattermostChannel{"#mychannel", "*"}}, mattermost.Channels)
	assert.Equal(t, "Flux Deployer", mattermost.Username)
	assert.Equal(t, "https://user-images.githubusercontent.com/27962005/35868977-0d5f85f6-0b2c-11e8-9fa8-8e4eaf35161a.png", mattermost.IconURL)
}

func TestMattermostOverrides(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("mattermost_url", "https://my_mattermost/")
	config.Set("mattermost_channel", "#mychannel=namespace")
	config.Set("mattermost_username", "my user")
	config.Set("mattermost_icon_url", "https://test-icon.png")

	mattermost, err := NewMattermost(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://my_mattermost/", mattermost.Url)
	assert.Equal(t, []MattermostChannel{MattermostChannel{"#mychannel", "namespace"}}, mattermost.Channels)
	assert.Equal(t, "my user", mattermost.Username)
	assert.Equal(t, "https://test-icon.png", mattermost.IconURL)
}

func TestMattermostChannel(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("mattermost_url", "https://my_mattermost/")
	config.Set("mattermost_channel", "#mychannel=*, #namespace=namespace")

	mattermost, err := NewMattermost(config)
	assert.Nil(t, err)

	assert.Equal(t, []MattermostChannel{
		MattermostChannel{"#mychannel", "*"},
		MattermostChannel{"#namespace", "namespace"},
	}, mattermost.Channels)
}

func TestMattermostMissingChannel(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("mattermost_url", "https://my_mattermost/")

	_, err := NewMattermost(config)
	assert.NotNil(t, err)
}

func TestMattermostMissingSlackUrl(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("mattermost_channel", "#channel")

	_, err := NewMattermost(config)
	assert.NotNil(t, err)
}

func TestMattermostNewLine(t *testing.T) {
	mattermost := Mattermost{}
	assert.Equal(t, "\n", mattermost.NewLine())
}

func TestMattermostFormatLink(t *testing.T) {
	mattermost := Mattermost{}
	assert.Equal(t, "<https://google.com|title>", mattermost.FormatLink("https://google.com", "title"))
}

func TestNewMattermostMessage(t *testing.T) {
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

	mattermostMessages := testMattermost.NewMattermostMessage(message)
	assert.Len(t, mattermostMessages, 2)

	assert.Equal(t, "#channel", mattermostMessages[0].ChannelName)
	assert.Equal(t, "#namespace", mattermostMessages[1].ChannelName)
	assert.Equal(t, testMattermost.IconURL, mattermostMessages[0].IconURL)
	assert.Equal(t, testMattermost.Username, mattermostMessages[0].Username)

	attach := mattermostMessages[0].Attachments[0]
	assert.Equal(t, "#4286f4", attach.Color)
	assert.Equal(t, message.TitleLink, attach.TitleLink)
	assert.Equal(t, message.Title, attach.Title)
}

func TestMattermostSend(t *testing.T) {
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

	mattermostMessage := MattermostMessage{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&mattermostMessage)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	testMattermost.Url = ts.URL

	err := testMattermost.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Contains(t, testMattermost.NewMattermostMessage(message), mattermostMessage)
}

func TestMattermostSendNon200(t *testing.T) {
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

	testMattermost.Url = ts.URL

	err := testMattermost.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMattermostSendHTTPError(t *testing.T) {
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

	testMattermost.Url = ts.URL

	err := testMattermost.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMattermostName(t *testing.T) {
	Mattermost := Mattermost{}

	assert.Equal(t, "Mattermost", Mattermost.Name())
}

func TestMattermostImplementsExporter(t *testing.T) {
	_ = Exporter(&Mattermost{})
}
