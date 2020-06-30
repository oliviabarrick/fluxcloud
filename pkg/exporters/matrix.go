package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"
)

type MatrixMessage struct {
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
	Body          string `json:"body"`
}

// The Matrix exporter sends Flux events to a Matrix channel.
type Matrix struct {
	url         string
	fullUrl     string
	roomId      string
	accessToken string
}

// Initialize a new Matrix instance
func NewMatrix(config config.Config) (*Matrix, error) {
	var err error
	s := Matrix{}

	s.url, err = config.Required("Matrix_url")
	if err != nil {
		return nil, err
	}

	s.accessToken, err = config.Required("Matrix_token")
	if err != nil {
		return nil, err
	}

	s.roomId, err = config.Required("Matrix_room_id")
	if err != nil {
		return nil, err
	}

	s.fullUrl, err = s.GetUrl()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Matrix) GetUrl() (string, error) {
	parsed, err := url.Parse(s.url)
	if err != nil {
		return "", err
	}

	pathPart := fmt.Sprintf("_matrix/client/r0/rooms/%s/send/m.room.message", s.roomId)
	parsed.Path = filepath.Join(parsed.Path, pathPart)

	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return "", err
	}

	query.Add("access_token", s.accessToken)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

// Send a Message to Matrix
func (s *Matrix) Send(c context.Context, client *http.Client, message msg.Message) error {
	b := new(bytes.Buffer)

	body := fmt.Sprintf("<a href='%s'>%s</a><br>%s", message.TitleLink, message.Title, message.Body)

	err := json.NewEncoder(b).Encode(MatrixMessage{
		MsgType:       "m.text",
		Format:        "org.matrix.custom.html",
		FormattedBody: body,
		Body:          message.Title,
	})
	if err != nil {
		log.Print("Could encode message to Matrix:", err)
		return err
	}

	log.Print(string(b.Bytes()))

	req, _ := http.NewRequest("POST", s.fullUrl, b)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(c)

	res, err := client.Do(req)
	if err != nil {
		log.Print("Could not post to matrix:", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to Matrix, status code:", res.StatusCode)
		return fmt.Errorf("Could not post to Matrix, status code: %d", res.StatusCode)
	}

	return nil
}

// Return the new line character for Matrix messages
func (s *Matrix) NewLine() string {
	return "</br>"
}

// Return a formatted link for Matrix.
func (s *Matrix) FormatLink(link string, name string) string {
	return fmt.Sprintf("<a href='%s'>%s</a>", link, name)
}

// Return the name of the exporter.
func (s *Matrix) Name() string {
	return "Matrix"
}
