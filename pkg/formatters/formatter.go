package formatters

import (
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/msg"
	fluxevent "github.com/weaveworks/flux/event"
)

// Formats a flux event for an exporter
type Formatter interface {
	FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message
}
