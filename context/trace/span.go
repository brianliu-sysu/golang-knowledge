package trace

import (
	"context"
	"time"
)

type Span struct {
	TraceID  string
	SpanID   string
	ParentID string

	Operation string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	Tags map[string]string
	Logs []LogEntry
	Err  bool
}

type LogEntry struct {
	Time    time.Time
	Message string
	Fields  map[string]any
}

func StartSpan(ctx context.Context, operation string) (context.Context, *Span) {
	tc, ok := FromContext(ctx)
	if !ok {
		tc = NewTraceContext()
		ctx = WithContext(ctx, tc)
	}

	span := &Span{
		TraceID:  tc.TraceID,
		ParentID: tc.SpanID,
		SpanID:   generateID(),

		Operation: operation,
		StartTime: time.Now(),

		Tags: make(map[string]string),
		Logs: make([]LogEntry, 0),
	}

	childTC := &TraceContext{
		TraceID:  span.TraceID,
		SpanID:   span.SpanID,
		ParentID: span.ParentID,

		Baggage: tc.CopyBaggage(),
	}

	return WithContext(ctx, childTC), span
}

func (s *Span) Finish() {
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)

	DefaultExporter.Export(s)
}

func (s *Span) SetTag(k, v string) {
	s.Tags[k] = v
}

func (s *Span) LogEvent(message string, fields map[string]any) {
	s.Logs = append(s.Logs, LogEntry{
		Time:    time.Now(),
		Message: message,
		Fields:  fields,
	})
}

func (s *Span) SetError(err error) {
	s.Err = true
	s.Tags["error"] = err.Error()

	s.Logs = append(s.Logs, LogEntry{
		Time:    time.Now(),
		Message: "error",
		Fields: map[string]any{
			"error": err.Error(),
		},
	})
}
