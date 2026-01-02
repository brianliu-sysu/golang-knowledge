package trace

import (
	"encoding/json"
	"fmt"
)

var (
	DefaultExporter Exporter = &ConsoleExporter{}
)

type Exporter interface {
	Export(span *Span) error
}

func SetExporter(exporter Exporter) {
	DefaultExporter = exporter
}

type ConsoleExporter struct {
}

func (*ConsoleExporter) Export(span *Span) error {
	data, _ := json.Marshal(span)
	fmt.Println(string(data))
	return nil
}
