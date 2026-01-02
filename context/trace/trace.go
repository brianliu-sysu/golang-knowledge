package trace

import "context"

type TraceContext struct {
	TraceID  string
	SpanID   string
	ParentID string

	Baggage map[string]string
}

func NewTraceContext() *TraceContext {
	return &TraceContext{
		TraceID: generateID(),
		SpanID:  generateID(),
		Baggage: make(map[string]string),
	}
}

func (t *TraceContext) SetBaggage(k, v string) {
	t.Baggage[k] = v
}

func (t *TraceContext) GetBaggage(k string) (string, bool) {
	v, ok := t.Baggage[k]
	return v, ok
}

func (t *TraceContext) CopyBaggage() map[string]string {
	baggage := make(map[string]string)
	for k, v := range t.Baggage {
		baggage[k] = v
	}

	return baggage
}

type contextKey struct{}

var traceContextKey contextKey

func FromContext(ctx context.Context) (*TraceContext, bool) {
	tx, ok := ctx.Value(traceContextKey).(*TraceContext)
	return tx, ok
}

func WithContext(ctx context.Context, tc *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey, tc)
}
