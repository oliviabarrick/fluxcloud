package utils

import (
	"encoding/json"
	"io"

	fluxevent "github.com/weaveworks/flux/event"
)

// Parse a flux event from Json into a flux Event struct.
func ParseFluxEvent(reader io.Reader) (event fluxevent.Event, err error) {
	err = json.NewDecoder(reader).Decode(&event)
	return
}
