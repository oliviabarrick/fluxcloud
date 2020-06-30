package exporters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"
)

func TestMSTeamsDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("msteams_url", "https://msteams.com/uuid/aaa")

	msteams, err := NewMSTeams(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://msteams.com/uuid/aaa", msteams.Url)
}

func TestMSTeamsMissingURL(t *testing.T) {
	config := config.NewFakeConfig()

	_, err := NewMSTeams(config)
	assert.NotNil(t, err)
}

func TestMSTeamsNewLine(t *testing.T) {
	msteams := MSTeams{}
	assert.Equal(t, "\n", msteams.NewLine())
}

func TestMSTeamsFormatLink(t *testing.T) {
	msteams := MSTeams{}
	assert.Equal(t, "[title](https://test.com/aaa)", msteams.FormatLink("https://test.com/aaa", "title"))
}

func TestNewMSTeamsMessage(t *testing.T) {
	msteams := MSTeams{}
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	msteamsMsg := msteams.NewMSTeamsMessage(message)

	assert.Equal(t, "https://schema.org/extensions", msteamsMsg.Context)
	assert.Equal(t, "MessageCard", msteamsMsg.Type)
	assert.Equal(t, "4286f4", msteamsMsg.ThemeColor)
	assert.Equal(t, message.Title, msteamsMsg.Title)
	assert.Equal(t, message.Body, msteamsMsg.Text)
}

func TestMSTeamsSend(t *testing.T) {
	msteams := MSTeams{}

	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	msteamsMsg := MSTeamsMessage{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&msteamsMsg)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	msteams.Url = ts.URL

	err := msteams.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Equal(t, msteams.NewMSTeamsMessage(message), msteamsMsg)
}

func TestMSTeamsSendNon200(t *testing.T) {
	msteams := MSTeams{}
	message := msg.Message{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	msteams.Url = ts.URL

	err := msteams.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMSTeamsSendHTTPError(t *testing.T) {
	msteams := MSTeams{}
	message := msg.Message{}

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	msteams.Url = ts.URL

	err := msteams.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMSTeamsName(t *testing.T) {
	msteams := MSTeams{}
	assert.Equal(t, "MS Teams", msteams.Name())
}

func TestMSTeamsImplementsExporter(t *testing.T) {
	_ = Exporter(&MSTeams{})
}
