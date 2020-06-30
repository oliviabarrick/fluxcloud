package exporters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/topfreegames/fluxcloud/pkg/msg"
)

type FakeExporter struct {
	Sent []msg.Message
}

func (f *FakeExporter) Send(_ context.Context, _ *http.Client, message msg.Message) error {
	f.Sent = append(f.Sent, message)
	return nil
}

func (f *FakeExporter) NewLine() string {
	return "\n"
}

func (f *FakeExporter) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

func (f *FakeExporter) Name() string {
	return "Fake"
}
