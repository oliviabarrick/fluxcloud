package apis

import (
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"time"
)

// All of the configuration necessary to run a fluxcloud API
type APIConfig struct {
	Server    *http.ServeMux
	Client    *http.Client
	Exporter  []exporters.Exporter
	Formatter formatters.Formatter
	Config    config.Config
}

// Initialize API configuration
func NewAPIConfig(f formatters.Formatter, e []exporters.Exporter, c config.Config) APIConfig {
	return APIConfig{
		Server: http.NewServeMux(),
		Client: &http.Client{
			Timeout:   120 * time.Second,
			Transport: &ochttp.Transport{},
		},
		Formatter: f,
		Exporter:  e,
		Config:    c,
	}
}

// Listen on addr
func (a *APIConfig) Listen(addr string) error {
	if os.Getenv("JAEGER_ENDPOINT") != "" {
		exporter, err := jaeger.NewExporter(jaeger.Options{
			CollectorEndpoint: os.Getenv("JAEGER_ENDPOINT"),
			Process: jaeger.Process{
				ServiceName: "fluxcloud",
			},
		})
		if err != nil {
			return err
		}

		trace.RegisterExporter(exporter)
	}

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	server := &http.Server{
		Addr:    addr,
		Handler: &ochttp.Handler{Handler: a.Server, IsPublicEndpoint: false},
	}

	return server.ListenAndServe()
}
