package main

import (
	"log"

	"github.com/justinbarrick/fluxcloud/pkg/apis"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
)

func initExporter(config config.Config) (exporter exporters.Exporter) {
	exporterType := config.Optional("Exporter_type", "slack")

	var err error

	switch exporterType {
	case "webhook":
		exporter, err = exporters.NewWebhook(config)
	default:
		exporter, err = exporters.NewSlack(config)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using %s exporter", exporter.Name())
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
