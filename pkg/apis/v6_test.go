package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	test_utils "github.com/justinbarrick/fluxcloud/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/formatters"
)

func TestHandleV6(t *testing.T) {
	fakeExporter := &exporters.FakeExporter{}
	config := config.NewFakeConfig()
	config.Set("github_url", "https://github.com")

	formatter, _ := formatters.NewDefaultFormatter(config)

	apiConfig := APIConfig{
		Server:    http.NewServeMux(),
		Exporter:  []exporters.Exporter{fakeExporter},
		Formatter: formatter,
	}

	HandleV6(apiConfig)

	event := test_utils.NewFluxSyncEvent()
	data, _ := json.Marshal(event)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:3030/v6/events", bytes.NewBuffer(data))

	recorder := httptest.NewRecorder()
	apiConfig.Server.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)

	formatted := formatter.FormatEvent(event, fakeExporter)
	assert.Equal(t, formatted.Title, fakeExporter.Sent[0].Title, formatted.Title)
	assert.Equal(t, formatted.Body, fakeExporter.Sent[0].Body, formatted.Body)
}
