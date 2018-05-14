package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/websocket"
	fluxevent "github.com/weaveworks/flux/event"
)

var (
	slack_url = os.Getenv("SLACK_URL")
	slack_username = os.Getenv("SLACK_USERNAME")
	slack_channel = os.Getenv("SLACK_CHANNEL")
	slack_icon_emoji = os.Getenv("SLACK_ICON_EMOJI")
	github_link = os.Getenv("GITHUB_URL")
)

var upgrader = websocket.Upgrader{}

type SlackMessage struct {
	Channel     string `json:"channel"`
	IconEmoji   string `json:"icon_emoji"`
	Username    string `json:"username"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}

// Send a SlackMessage to Slack
func slack(message SlackMessage) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(message)
	if err != nil {
		log.Print("Could encode message to slack:", err)
	}

	res, err := http.Post(slack_url, "application/json", b)
	if err != nil {
		log.Print("Could not post to slack:", err)
		return
	}

	if res.StatusCode != 200 {
		log.Print("Could not post to slack, status: ", res.Status)
	}
}

// Send Event to Slack
func slackEvent(event fluxevent.Event) {
	body := "Event: " + event.String() + "\n"
	commit_id := ""

	if len(event.ServiceIDs) == 0 {
		return
	}

	switch event.Type{
		case fluxevent.EventSync:
			metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
			commit_id = metadata.Commits[0].Revision

			body += "Commits:\n"
			for _, commit := range metadata.Commits {
				link := github_link + "/commit/" + commit.Revision
				body += "\n* <" + link + "|" + commit.Revision + ">: " + commit.Message
			}
		default:
	}

	body += "\n\nResources updated:\n"
	for _, serviceId := range event.ServiceIDStrings() {
		body += "\n* " + serviceId
	}

	message := SlackMessage{
		Channel: slack_channel,
		IconEmoji: slack_icon_emoji,
		Username: slack_username,
		Attachments: []SlackAttachment{
			SlackAttachment{
				Color: "#4286f4",
				TitleLink: github_link + "/commit/" + commit_id,
				Title: "Applied flux changes to cluster",
				Text: body,
			},
		},
	}

	slack(message)
}

// Handle Flux WebSocket connections
func fluxcloud(w http.ResponseWriter, r *http.Request) {
	log.Print("Request for:", r.URL)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		log.Println("client disconnected")
		c.Close()
	}()

	log.Println("client connected!")

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// Handle Flux webhooks
func fluxcloudEvent(w http.ResponseWriter, r *http.Request) {
	log.Print("Request for:", r.URL)

	event := fluxevent.Event{}

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	go slackEvent(event)
  w.WriteHeader(200)
}

func settings() {
	if slack_username == "" {
		slack_username = "Flux Deployer"
	}

	if slack_icon_emoji == "" {
		slack_icon_emoji = ":star-struck:"
	}

	if slack_channel == "" {
		log.Fatal("Set the SLACK_CHANNEL environment variable to the channel to notify.")
	}

	if slack_url == "" {
		log.Fatal("Set the SLACK_URL environment variable to the Slack webhook URL.")
	}
}

func main() {
	log.SetFlags(0)
	settings()
	http.HandleFunc("/v6/events", fluxcloudEvent)
	http.HandleFunc("/", fluxcloud)
	log.Fatal(http.ListenAndServe(":3031", nil))
}
