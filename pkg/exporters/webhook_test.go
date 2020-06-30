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

func TestWebhookDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("webhook_url", "https://mywebhook/")

	webhook, err := NewWebhook(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://mywebhook/", webhook.Url)
}

func TestWehookMissingURL(t *testing.T) {
	config := config.NewFakeConfig()

	_, err := NewSlack(config)
	assert.NotNil(t, err)
}

func TestWebhookNewLine(t *testing.T) {
	webhook := Webhook{}
	assert.Equal(t, "\n", webhook.NewLine())
}

func TestWebhookFormatLink(t *testing.T) {
	webhook := Webhook{}
	assert.Equal(t, "<https://google.com|title>", webhook.FormatLink("https://google.com", "title"))
}

func TestWebhookSend(t *testing.T) {
	webhook := Webhook{}

	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	receivedMessage := msg.Message{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedMessage)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	webhook.Url = ts.URL

	err := webhook.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Equal(t, receivedMessage, message)
}

func TestWebhookSendNon200(t *testing.T) {
	webhook := Webhook{}
	message := msg.Message{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	webhook.Url = ts.URL

	err := webhook.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestWebhookSendHTTPError(t *testing.T) {
	webhook := Webhook{}
	message := msg.Message{}

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	webhook.Url = ts.URL

	err := webhook.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestWebhookName(t *testing.T) {
	webhook := Webhook{}
	assert.Equal(t, "Webhook", webhook.Name())
}

func TestWebhookImplementsExporter(t *testing.T) {
	_ = Exporter(&Webhook{})
}
