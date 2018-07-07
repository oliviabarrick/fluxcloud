package formatters

import (
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	fluxevent "github.com/weaveworks/flux/event"
	"testing"
)

func TestNewDefaultFormatter(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("github_url", "https://github.com/")

	formatter, err := NewDefaultFormatter(config)
	assert.Nil(t, err)
	assert.Equal(t, "https://github.com/", formatter.vcsLink)
	assert.Equal(t, config, formatter.config)
}

func TestNewDefaultFormatterNoGithubLink(t *testing.T) {
	config := config.NewFakeConfig()

	_, err := NewDefaultFormatter(config)
	assert.NotNil(t, err)
}

func TestDefaultFormatterImplementsFormatter(t *testing.T) {
	_ = Formatter(&DefaultFormatter{})
}

func TestDefaultFormatterFormatSyncEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}

	event := test_utils.NewFluxSyncEvent()

	msg := d.FormatEvent(event, &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/810c2e6f22ac5ab7c831fe0dd697fe32997b098f", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventSync, msg.Type)
	assert.Equal(t, `Event: Sync: 810c2e6, default:deployment/test
Commits: 

* <https://github.com/commit/810c2e6f22ac5ab7c831fe0dd697fe32997b098f|810c2e6>: change test image

Resources updated:

* default:deployment/test`, msg.Body)
	assert.Equal(t, event, msg.Event)
}

func TestDefaultFormatterFormatCommitEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxCommitEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventCommit, msg.Type)
	assert.Equal(t, `Event: Commit: d644e1a, default:deployment/test

Resources updated:

* default:deployment/test`, msg.Body)
}

func TestDefaultFormatterFormatAutoReleaseEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxAutoReleaseEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventAutoRelease, msg.Type)
	assert.Equal(t, `Event: Automated release of justinbarrick/nginx:test3

Resources updated:

* default:deployment/test`, msg.Body)
}

func TestDefaultFormatterFormatUpdatePolicyEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxUpdatePolicyEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventSync, msg.Type)
	assert.Equal(t, `Event: Sync: d644e1a, default:deployment/test
Commits: 

* <https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036|d644e1a>: Automated: default:deployment/test

Resources updated:

* default:deployment/test`, msg.Body)
}
