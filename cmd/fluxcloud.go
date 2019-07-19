package main

import (
	"log"
	"strings"

	"github.com/justinbarrick/fluxcloud/pkg/apis"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
)

func initExporter(config config.Config) (exporter []exporters.Exporter) {
	exporterType := config.Optional("Exporter_type", "slack")

	exporterTypes := strings.Split(exporterType, ",")

	for _, v := range exporterTypes {
		if v == "webhook" {
			webhook, err := exporters.NewWebhook(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, webhook)
		}

		if v == "matrix" {
			matrix, err := exporters.NewMatrix(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, matrix)
		}

		if v == "msteams" {
			msteams, err := exporters.NewMSTeams(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, msteams)
		}

		if v == "slack" {
			slack, err := exporters.NewSlack(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, slack)
		}
	}

	for _, e := range exporter {
		log.Printf("Using %s exporter", e.Name())
	}

	return exporter
}

func main() {
	log.SetFlags(0)

	config := &config.DefaultConfig{}

	formatter, err := formatters.NewDefaultFormatter(config)
	if err != nil {
		log.Fatal(err)
	}

	apiConfig := apis.NewAPIConfig(formatter, initExporter(config), config)

	apis.HandleWebsocket(apiConfig)
	apis.HandleV6(apiConfig)
	log.Fatal(apiConfig.Listen(config.Optional("listen_address", ":3031")))
}
