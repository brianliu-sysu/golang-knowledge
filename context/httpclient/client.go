package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	MaxResponseSize = 10 * 1024 * 1024 // 10MB
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	if client == nil {
		return NewDefaultClient()
	}

	return &Client{
		client: client,
	}
}

func NewDefaultClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

func (c *Client) Do(ctx context.Context, method, url string, body io.Reader, header map[string]string) (data []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("http new request failed [%s, %s]: err:%w", method, url, err)
	}

	if len(header) > 0 {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request failed [%s, %s]: err:%w", method, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http request failed [%s, %s]: err:status code is %v", method, url, resp.StatusCode)
	}

	limiterReader := io.LimitReader(resp.Body, MaxResponseSize)
	data, err = io.ReadAll(limiterReader)
	if err != nil {
		return nil, fmt.Errorf("http request failed [%s, %s]: err:%w", method, url, err)
	}

	if int64(len(data)) == MaxResponseSize {
		return nil, fmt.Errorf("http request failed [%s, %s]: err:exceed the max size", method, url)
	}
	return
}
