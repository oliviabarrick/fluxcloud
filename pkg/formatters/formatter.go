package formatters

import (
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	fluxevent "github.com/weaveworks/flux/event"
)

// Formats a flux event for an exporter
type Formatter interface {
	FormatEvent(event fluxevent.Event, exporter exporters.Exporter) msg.Message
}
