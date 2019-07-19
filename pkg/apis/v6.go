package apis

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/justinbarrick/fluxcloud/pkg/utils"
)

// Handle Flux events
func HandleV6(config APIConfig) error {
	config.Server.HandleFunc("/v6/events", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request for:", r.URL)

		eventStr, err := ioutil.ReadAll(r.Body)
		log.Print(string(eventStr))

		event, err := utils.ParseFluxEvent(bytes.NewBuffer(eventStr))
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), 400)
			return
		}

		for _, exporter := range config.Exporter {
			message := config.Formatter.FormatEvent(event, exporter)
			if message.Title == "" {
				w.WriteHeader(200)
				return
			}

			err = exporter.Send(r.Context(), config.Client, message)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
		}
		w.WriteHeader(200)
	})

	return nil
}
