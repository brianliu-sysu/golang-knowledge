package main

import (
	"context"
	"fmt"
	"golang-knowledge/context/trace"
	"net/http"
	"time"
)

func httpHandler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	go func() {
		requestDB(ctx)
	}()

	select {
	case <-time.After(5 * time.Second):
		resp.Write([]byte("complete"))
		fmt.Println("complete")
	case <-ctx.Done():
		fmt.Println("client is close, err:", ctx.Err())
	}

}

func requestDB(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "requestDB")
	defer span.Finish()
	span.SetTag("db.type", "mysql")
	select {
	case <-time.After(time.Second * 10):
		span.LogEvent("done", map[string]any{"time": time.Now()})
		fmt.Println("done")
	case <-ctx.Done():
		span.SetError(ctx.Err())
		fmt.Println("cancel")
	}
}

func main() {
	http.ListenAndServe(":8080", trace.HTTPMiddleware(http.HandlerFunc(httpHandler)))
}
