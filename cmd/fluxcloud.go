package main

import (
	"github.com/justinbarrick/fluxcloud/pkg/apis"
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
	"log"
)

func main() {
	log.SetFlags(0)

	config := &config.DefaultConfig{}

	formatter, err := formatters.NewDefaultFormatter(config)
	if err != nil {
		log.Fatal(err)
	}

	exporter, err := exporters.NewSlack(config)
	if err != nil {
		log.Fatal(err)
	}

	apiConfig := apis.NewAPIConfig(formatter, exporter, config)
	apis.HandleWebsocket(apiConfig)
	apis.HandleV6(apiConfig)
	log.Fatal(apiConfig.Listen(":3031"))
}
