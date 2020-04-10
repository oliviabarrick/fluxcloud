package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
)

type Facebook struct {
	Url string
	ThreadKey string
	AccessToken string
}

type FacebookMessage struct {
	Recipient FacebookRecipientPayload `json:"recipient"`
	Message FacebookMessagePayload `json:"message"`
}

type FacebookRecipientPayload struct {
	ThreadKey string `json:"thread_key"`
}

type FacebookMessagePayload struct {
	Text string `json:"text"`
}

func NewFacebook(config config.Config) (*Facebook, error) {
	var err error
	s := Facebook{}

	s.ThreadKey, err = config.Required("facebook_thread_key")
	if err != nil {
		return nil, err
	}

	s.AccessToken, err = config.Required("facebook_access_token")
	if err != nil {
		return nil, err
	}

	baseUrl := config.Optional("facebook_base_url", "https://graph.facebook.com/v5.0/me/messages")

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("access_token", s.AccessToken)
	u.RawQuery = q.Encode()

	s.Url = u.String()

	return &s, nil
}

func (s *Facebook) Send(c context.Context, client *http.Client, message msg.Message) error {
	facebookMessage := FacebookMessage{
		Recipient: FacebookRecipientPayload{
			ThreadKey: s.ThreadKey,
		},
		Message: FacebookMessagePayload{
			Text: message.Body,
		},
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(facebookMessage)
	if err != nil {
		log.Print("Could encode message to Facebook:", err)
		return err
	}

	log.Print(string(b.Bytes()))

	req, _ := http.NewRequest("POST", s.Url, b)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(c)

	res, err := client.Do(req)
	if err != nil {
		log.Print("Could not post to Facebook:", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to Facebook, status: ", res.Status)
		return errors.New(fmt.Sprintf("Could not post to Facebook, status: %d", res.StatusCode))
	}

	return nil
}

// Return the new line character for Facebook messages
func (s *Facebook) NewLine() string {
	return "\n"
}

// Return a formatted link for Facebook.
func (s *Facebook) FormatLink(link string, name string) string {
	return link
}

// Return the name of the exporter.
func (s *Facebook) Name() string {
	return "Facebook"
}
