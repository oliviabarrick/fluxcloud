package apis

import (
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/formatters"
	"net/http"
	"time"
)

// All of the configuration necessary to run a fluxcloud API
type APIConfig struct {
	Server    *http.ServeMux
	Client    *http.Client
	Exporter  exporters.Exporter
	Formatter formatters.Formatter
	Config    config.Config
}

// Initialize API configuration
func NewAPIConfig(f formatters.Formatter, e exporters.Exporter, c config.Config) APIConfig {
	return APIConfig{
		Server: http.NewServeMux(),
		Client: &http.Client{
			Timeout: 120 * time.Second,
		},
		Formatter: f,
		Exporter:  e,
		Config:    c,
	}
}

// Listen on addr
func (a *APIConfig) Listen(addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: a.Server,
	}

	return server.ListenAndServe()
}
