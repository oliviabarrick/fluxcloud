package main

import (
	"log"

	"github.com/justinbarrick/fluxcloud/pkg/apis"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
)

func main() {
	log.SetFlags(0)

	config := &config.DefaultConfig{}

	formatter, err := formatters.NewDefaultFormatter(config)
	if err != nil {
		log.Fatal(err)
	}

	exporterType := config.Optional("Exporter_type", "slack")

	var slackExporter *exporters.Slack
	var webhookExporter *exporters.Webhook

	var apiConfig apis.APIConfig

	switch exporterType {
	case "webhook":
		log.Println("Using Webhook exporter")

		webhookExporter, err = exporters.NewWebhook(config)
		if err != nil {
			log.Fatal(err)
		}

		apiConfig = apis.NewAPIConfig(formatter, webhookExporter, config)
	default:
		log.Println("Using Slack exporter")

		slackExporter, err = exporters.NewSlack(config)
		if err != nil {
			log.Fatal(err)
		}

		apiConfig = apis.NewAPIConfig(formatter, slackExporter, config)
	}

	apis.HandleWebsocket(apiConfig)
	apis.HandleV6(apiConfig)
	log.Fatal(apiConfig.Listen(":3031"))
}
