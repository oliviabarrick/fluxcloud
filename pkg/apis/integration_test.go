package apis

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	test_utils "github.com/justinbarrick/fluxcloud/pkg/utils/test"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/fluxcloud/pkg/msg"

	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/formatters"
)

func TestSlackIntegrationTest(t *testing.T) {
	var exporter *exporters.Slack
	var formatter *formatters.DefaultFormatter

	event := test_utils.NewFluxSyncEvent()
	data, err := json.Marshal(event)
	assert.Nil(t, err)

	reqCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/slack", r.URL.Path)

		sent := exporters.SlackMessage{}
		json.NewDecoder(r.Body).Decode(&sent)
		formatted := exporter.NewSlackMessage(formatter.FormatEvent(event, exporter))
		assert.Equal(t, sent, formatted[0])
		reqCount += 1
	}))
	defer ts.Close()

	config := config.NewFakeConfig()
	config.Set("slack_url", ts.URL+"/slack")
	config.Set("slack_channel", "#kubernetes")
	config.Set("github_url", "https://github.com")

	exporter, err = exporters.NewSlack(config)
	assert.Nil(t, err)

	formatter, err = formatters.NewDefaultFormatter(config)
	assert.Nil(t, err)

	apiConfig := NewAPIConfig(formatter, []exporters.Exporter{exporter}, config)
	HandleV6(apiConfig)
	HandleWebsocket(apiConfig)

	apiServer := httptest.NewServer(apiConfig.Server)
	defer apiServer.Close()

	resp, err := http.Post(apiServer.URL+"/v6/events", "application/json", bytes.NewBuffer(data))
	assert.Nil(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, reqCount)
}

func TestWebhookAndSlackIntegrationTest(t *testing.T) {
	require := require.New(t)
	var slackExporter *exporters.Slack
	var webhookExporter *exporters.Webhook
	var formatter *formatters.DefaultFormatter

	event := test_utils.NewFluxSyncEvent()
	data, err := json.Marshal(event)
	require.NoError(err)

	reqCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal("/slack", r.URL.Path)

		sent := exporters.SlackMessage{}
		json.NewDecoder(r.Body).Decode(&sent)
		formatted := slackExporter.NewSlackMessage(formatter.FormatEvent(event, slackExporter))
		assert.Equal(t, sent, formatted[0])
		reqCount += 1
	}))
	defer ts.Close()

	webhookReceiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal("/", r.URL.Path)
		require.NotEmpty(r.Body)
		body, _ := ioutil.ReadAll(r.Body)
		res := msg.Message{}
		err = json.Unmarshal([]byte(body), &res)
		require.NoError(err)
		require.Equal(res.Event.ID, event.ID)
	}))
	defer webhookReceiver.Close()

	config := config.NewFakeConfig()
	config.Set("slack_url", ts.URL+"/slack")
	config.Set("slack_channel", "#kubernetes")
	config.Set("github_url", "https://github.com")
	config.Set("webhook_url", webhookReceiver.URL)
	config.Set("exporter_type", "slack,webhook")

	slackExporter, err = exporters.NewSlack(config)
	require.NoError(err)
	webhookExporter, err = exporters.NewWebhook(config)
	require.NoError(err)

	formatter, err = formatters.NewDefaultFormatter(config)
	require.NoError(err)

	apiConfig := NewAPIConfig(formatter, []exporters.Exporter{slackExporter, webhookExporter}, config)
	HandleV6(apiConfig)
	HandleWebsocket(apiConfig)

	apiServer := httptest.NewServer(apiConfig.Server)
	defer apiServer.Close()

	resp, err := http.Post(apiServer.URL+"/v6/events", "application/json", bytes.NewBuffer(data))
	require.NoError(err)

	require.Equal(200, resp.StatusCode)
	require.Equal(1, reqCount)
}
