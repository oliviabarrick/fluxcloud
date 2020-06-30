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
	"github.com/zorkian/go-datadog-api"
)

var testDatadog = Datadog{
	datadogClient: datadog.NewClient("app", "api"),
}

func TestDatadogDefault(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("datadog_app_key", "app_xxx")
	config.Set("datadog_api_key", "api_xxx")

	_, err := NewDatadog(config)
	assert.Nil(t, err)
}

func TestDatadogMissingAppKey(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("api_key", "api_xxx")

	_, err := NewDatadog(config)
	assert.NotNil(t, err)
}

func TestDatadogMissingApiKey(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("app_key", "app_xxx")

	_, err := NewDatadog(config)
	assert.NotNil(t, err)
}

func TestDatadogNewLine(t *testing.T) {
	Datadog := Datadog{}
	assert.Equal(t, "\n", Datadog.NewLine())
}

func TestDatadogFormatLink(t *testing.T) {
	Datadog := Datadog{}
	assert.Equal(t, "<https://google.com|title>", Datadog.FormatLink("https://google.com", "title"))
}

func TestNewDatadogMessage(t *testing.T) {
	defaultResourceID, _ := flux.ParseResourceID("default:resource/name")
	nsResourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			Type: "sync",
			ServiceIDs: []flux.ResourceID{
				defaultResourceID,
				nsResourceID,
			},
		},
	}

	DatadogEvents := testDatadog.NewDatadogEvent(message)
	assert.Len(t, DatadogEvents, 2)
}

func TestDatadogTags(t *testing.T) {
	deployResourceID, _ := flux.ParseResourceID("ns1:deploy/name-1")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			Type: "sync",
			ServiceIDs: []flux.ResourceID{
				deployResourceID,
			},
		},
	}
	DatadogEvents := testDatadog.NewDatadogEvent(message)
	assert.Equal(t, DatadogEvents[0].Tags, []string{"application:flux", "fluxEventType:sync", "fluxnamespace:ns1", "fluxkind:deploy", "fluxresourcename:name-1"})
}

func TestDatadogSend(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
		Event: fluxevent.Event{
			Type: "sync",
			ServiceIDs: []flux.ResourceID{
				resourceID,
			},
		},
	}

	DatadogMessage := DatadogEvent{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&DatadogMessage)
		fmt.Println(r.Body)
		fmt.Println(DatadogMessage)
	}))
	defer ts.Close()
	testDatadog.datadogClient.SetBaseUrl(ts.URL)

	err := testDatadog.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)

	assert.Contains(t, testDatadog.NewDatadogEvent(message), DatadogMessage)
}

func TestDatadogSendNon200(t *testing.T) {
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

	testDatadog.datadogClient.SetBaseUrl(ts.URL)

	err := testDatadog.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestDatadogSendHTTPError(t *testing.T) {
	resourceID, _ := flux.ParseResourceID("namespace:resource/name")
	message := msg.Message{
		Event: fluxevent.Event{
			Type: "sync",
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

	testDatadog.datadogClient.SetBaseUrl(ts.URL)

	err := testDatadog.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestDatadogName(t *testing.T) {

	Datadog := Datadog{}
	assert.Equal(t, "Datadog Events", Datadog.Name())
}

func TestDatadogImplementsExporter(t *testing.T) {
	_ = Exporter(&Datadog{})
}
