package apis

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

// Handle Flux WebSocket connections
func HandleWebsocket(config APIConfig) error {
	config.Server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	return nil
}
