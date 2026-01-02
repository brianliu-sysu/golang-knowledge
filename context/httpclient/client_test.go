package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDo_timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.Write([]byte("done"))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()
	_, err := NewDefaultClient().Do(ctx, "GET", server.URL, nil, nil)
	if err == nil {
		t.Fatal("should be error")
	}

	t.Logf("error:%v", err)
}

func TestDo_cancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.Write([]byte("done"))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	go func() {
		time.Sleep(time.Microsecond)
		cancel()
	}()
	_, err := NewDefaultClient().Do(ctx, "GET", server.URL, nil, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatal("should be cancel error")
	}
}
