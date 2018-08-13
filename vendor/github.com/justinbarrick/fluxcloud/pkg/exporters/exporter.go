package exporters

import (
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"net/http"
)

// An exporter sends a formatted event to an upstream.
type Exporter interface {
	// Send a message through the exporter.
	Send(client *http.Client, message msg.Message) error

	// Return a new line as a string for the exporter.
	NewLine() string

	// Return a link formatted for the exporter.
	FormatLink(link string, name string) string
}
