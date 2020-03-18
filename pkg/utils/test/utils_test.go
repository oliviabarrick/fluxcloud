package test_utils

import (
	"github.com/stretchr/testify/assert"
	fluxevent "github.com/weaveworks/flux/event"
	"testing"
)

func TestParseFluxEventSync(t *testing.T) {
	event := NewFluxSyncEvent()

	assert.Equal(t, fluxevent.EventSync, event.Type)

	assert.Equal(t, fluxevent.EventID(0), event.ID)
	assert.Equal(t, "default:deployment/test", event.ServiceIDs[0].String())
	assert.Equal(t, "info", event.LogLevel)

	metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
	commit := metadata.Commits[0]

	assert.Equal(t, "810c2e6f22ac5ab7c831fe0dd697fe32997b098f", commit.Revision)
	assert.Equal(t, "change test image", commit.Message)
	assert.Equal(t, true, metadata.Includes["other"])
}

func TestParseFluxEventDeleteSync(t *testing.T) {
	event := NewFluxDeleteEvent()

	assert.Equal(t, fluxevent.EventSync, event.Type)

	assert.Equal(t, fluxevent.EventID(0), event.ID)
	assert.Len(t, event.ServiceIDs, 0)
	assert.Equal(t, "info", event.LogLevel)

	metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
	commit := metadata.Commits[0]

	assert.Equal(t, "c6b7c44b4300f92b788bbc9bb6cb7282852300b4", commit.Revision)
	assert.Equal(t, "deleted k8s-global-objects", commit.Message)
	assert.Equal(t, true, metadata.Includes["other"])
}

func TestParseFluxEventCommit(t *testing.T) {
	event := NewFluxCommitEvent()

	assert.Equal(t, fluxevent.EventCommit, event.Type)

	_ = event.Metadata.(*fluxevent.CommitEventMetadata)
}

func TestParseFluxEventAutoRelease(t *testing.T) {
	event := NewFluxAutoReleaseEvent()

	assert.Equal(t, fluxevent.EventAutoRelease, event.Type)
	_ = event.Metadata.(*fluxevent.AutoReleaseEventMetadata)
}

func TestParseFluxEventSyncError(t *testing.T) {
	event := NewFluxSyncErrorEvent()

	assert.Equal(t, fluxevent.EventSync, event.Type)

	assert.Equal(t, fluxevent.EventID(0), event.ID)
	assert.Equal(t, "default:persistentvolumeclaim/test", event.ServiceIDs[0].String())
	assert.Equal(t, "info", event.LogLevel)

	metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
	commit := metadata.Commits[0]

	assert.Equal(t, "4997efcd4ac6255604d0d44eeb7085c5b0eb9d48", commit.Revision)
	assert.Equal(t, "create invalid resource", commit.Message)
	assert.Equal(t, true, metadata.Includes["other"])
	assert.Equal(t, "running kubectl: The PersistentVolumeClaim \"test\" is invalid: spec: Forbidden: field is immutable after creation", metadata.Errors[0].Error)
}
