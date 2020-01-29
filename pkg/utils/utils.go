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

func IsErrorEvent(event fluxevent.Event) bool {
	switch event.Type {
	case fluxevent.EventRelease:
		metadata := event.Metadata.(*fluxevent.ReleaseEventMetadata)
		return len(metadata.Result.Error()) > 0
	case fluxevent.EventAutoRelease:
		metadata := event.Metadata.(*fluxevent.AutoReleaseEventMetadata)
		return len(metadata.Result.Error()) > 0
	case fluxevent.EventCommit:
		metadata := event.Metadata.(*fluxevent.CommitEventMetadata)
		return len(metadata.Result.Error()) > 0
	case fluxevent.EventSync:
		metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
		return len(metadata.Errors) > 0
	default:
		return false
	}
}
