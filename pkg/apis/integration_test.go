package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
	"github.com/justinbarrick/fluxcloud/pkg/utils/test"
	"github.com/stretchr/testify/assert"
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

	apiConfig := NewAPIConfig(formatter, exporter, config)
	HandleV6(apiConfig)
	HandleWebsocket(apiConfig)

	apiServer := httptest.NewServer(apiConfig.Server)
	defer apiServer.Close()

	resp, err := http.Post(apiServer.URL+"/v6/events", "application/json", bytes.NewBuffer(data))
	assert.Nil(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, reqCount)
}
