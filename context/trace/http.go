package trace

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	TraceIDHeader  = "X-Trace-ID"
	SpanIDHeader   = "X-Span-ID"
	ParentIDHeader = "X-Parent-ID" // Fixed typo: ParendIDHeader -> ParentIDHeader
)

// isValidTraceID 校验 trace ID 格式（非空且长度合理）
func isValidTraceID(id string) bool {
	return len(id) > 0 && len(id) <= 64
}

func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		spanID := r.Header.Get(SpanIDHeader)

		var tc *TraceContext
		// 校验上游 trace header
		if !isValidTraceID(traceID) {
			tc = NewTraceContext()
		} else {
			parentID := spanID
			if parentID == "" {
				parentID = generateID()
			}
			tc = &TraceContext{
				TraceID:  traceID,
				SpanID:   generateID(),
				ParentID: parentID,
			}
		}

		// 将trace context 传递给下一个handler
		ctx := WithContext(r.Context(), tc)

		// construct a new span
		ctx, span := StartSpan(ctx, r.Method+" "+r.URL.Path)

		// panic recovery
		defer func() {
			if err := recover(); err != nil {
				span.SetTag("error", "true")
				span.SetTag("error.message", fmt.Sprintf("panic: %v", err))
				span.Finish()
				panic(err) // re-panic after recording
			}
		}()
		defer span.Finish()

		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())
		span.SetTag("http.host", r.Host)
		span.SetTag("http.remote_addr", r.RemoteAddr)

		// set header to the response
		w.Header().Set(TraceIDHeader, tc.TraceID)
		w.Header().Set(SpanIDHeader, tc.SpanID)
		w.Header().Set(ParentIDHeader, tc.ParentID)

		wrappedWriter := &ResponseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		// 记录状态码
		span.SetTag("http.status_code", fmt.Sprintf("%d", wrappedWriter.statusCode))

		// 对 4xx/5xx 错误状态码设置 error tag
		if wrappedWriter.statusCode >= 400 {
			span.SetTag("error", "true")
		}

		// 使用更符合 OpenTelemetry 语义的命名
		span.SetTag("http.duration_ms", fmt.Sprintf("%d", time.Since(span.StartTime).Milliseconds()))
	})
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.statusCode = statusCode
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriterWrapper) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.statusCode = http.StatusOK
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(data)
}

// Flush 实现 http.Flusher 接口（用于 SSE 等场景）
func (w *ResponseWriterWrapper) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack 实现 http.Hijacker 接口（用于 WebSocket 等场景）
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
}

// Push 实现 http.Pusher 接口（用于 HTTP/2 Server Push）
func (w *ResponseWriterWrapper) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("ResponseWriter does not implement http.Pusher")
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{
		client: client,
	}
}

func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	ctx, span := StartSpan(ctx, req.Method+" "+req.URL.Path)
	defer span.Finish()

	req.Header.Set(TraceIDHeader, span.TraceID)
	req.Header.Set(SpanIDHeader, span.SpanID)
	req.Header.Set(ParentIDHeader, span.ParentID)

	span.SetTag("http.method", req.Method)
	span.SetTag("http.url", req.URL.String())
	span.SetTag("http.host", req.Host)

	resp, err = c.client.Do(req.WithContext(ctx))
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	span.SetTag("http.status_code", fmt.Sprintf("%d", resp.StatusCode))
	if resp.StatusCode >= 400 {
		span.SetTag("error", "true")
	}

	return resp, nil
}
