package exporters

import (
	"fmt"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
	"net/http"
)

type FakeExporter struct {
	Sent []msg.Message
}

func (f *FakeExporter) Send(_ *http.Client, message msg.Message) error {
	f.Sent = append(f.Sent, message)
	return nil
}

func (f *FakeExporter) NewLine() string {
	return "\n"
}

func (f *FakeExporter) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}
