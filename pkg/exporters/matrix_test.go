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

func TestMatrix(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("matrix_url", "https://mymatrix/lol")
	config.Set("matrix_token", "mytoken")
	config.Set("matrix_room_id", "!myroom:mymatrix")

	matrix, err := NewMatrix(config)
	assert.Nil(t, err)

	assert.Equal(t, "https://mymatrix/lol", matrix.url)
	assert.Equal(t, "mytoken", matrix.accessToken)
	assert.Equal(t, "!myroom:mymatrix", matrix.roomId)
	assert.Equal(t, "https://mymatrix/lol/_matrix/client/r0/rooms/%21myroom:mymatrix/send/m.room.message?access_token=mytoken", matrix.fullUrl)
}

func TestMatrixNoProtocol(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("matrix_url", "mymatrix")
	config.Set("matrix_token", "mytoken")
	config.Set("matrix_room_id", "!myroom:mymatrix")

	matrix, err := NewMatrix(config)
	assert.Nil(t, err)

	assert.Equal(t, "mymatrix/_matrix/client/r0/rooms/%21myroom:mymatrix/send/m.room.message?access_token=mytoken", matrix.fullUrl)
}

func TestMatrixMissingURL(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("matrix_token", "mytoken")
	config.Set("matrix_room_id", "!myroom:mymatrix")

	_, err := NewMatrix(config)
	assert.NotNil(t, err)
}

func TestMatrixMissingAccessToken(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("matrix_url", "https://mymatrix/lol")
	config.Set("matrix_room_id", "!myroom:mymatrix")

	_, err := NewMatrix(config)
	assert.NotNil(t, err)
}

func TestMatrixMissingRoomId(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("matrix_url", "https://mymatrix/lol")
	config.Set("matrix_token", "mytoken")

	_, err := NewMatrix(config)
	assert.NotNil(t, err)
}

func TestMatrixNewLine(t *testing.T) {
	matrix := Matrix{}
	assert.Equal(t, "</br>", matrix.NewLine())
}

func TestMatrixFormatLink(t *testing.T) {
	matrix := Matrix{}
	assert.Equal(t, "<a href='https://google.com'>title</a>", matrix.FormatLink("https://google.com", "title"))
}

func TestMatrixSend(t *testing.T) {
	matrix := Matrix{}

	message := msg.Message{
		TitleLink: "https://myvcslink/",
		Title:     "The title of the message",
		Body:      "this is the message body",
	}

	receivedMessage := MatrixMessage{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedMessage)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	matrix.url = ts.URL
	matrix.roomId = "!myroom:myserver"
	matrix.accessToken = "myaccesstoken"
	matrix.fullUrl, _ = matrix.GetUrl()

	err := matrix.Send(context.TODO(), &http.Client{}, message)
	assert.Nil(t, err)
	assert.Equal(t, receivedMessage, MatrixMessage{
		MsgType:       "m.text",
		Format:        "org.matrix.custom.html",
		FormattedBody: "<a href='https://myvcslink/'>The title of the message</a><br>this is the message body",
		Body:          "The title of the message",
	})
}

func TestMatrixSendNon200(t *testing.T) {
	matrix := Matrix{}
	message := msg.Message{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	matrix.url = ts.URL
	matrix.roomId = "!myroom:myserver"
	matrix.accessToken = "myaccesstoken"
	matrix.fullUrl, _ = matrix.GetUrl()

	err := matrix.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMatrixSendHTTPError(t *testing.T) {
	matrix := Matrix{}
	message := msg.Message{}

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	}))
	defer ts.Close()

	matrix.url = ts.URL
	matrix.roomId = "!myroom:myserver"
	matrix.accessToken = "myaccesstoken"
	matrix.fullUrl, _ = matrix.GetUrl()

	err := matrix.Send(context.TODO(), &http.Client{}, message)
	assert.NotNil(t, err)
}

func TestMatrixName(t *testing.T) {
	matrix := Matrix{}
	assert.Equal(t, "Matrix", matrix.Name())
}

func TestMatrixImplementsExporter(t *testing.T) {
	_ = Exporter(&Matrix{})
}
