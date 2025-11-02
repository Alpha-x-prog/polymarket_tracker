package polymarket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const baseURL = "https://gamma-api.polymarket.com"

type Client struct {
	http    *http.Client
	limiter <-chan time.Time // простейший rate-limit
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{Timeout: 10 * time.Second},
		// 10 req/sec как пример (настроишь под себя)
		limiter: time.Tick(100 * time.Millisecond),
	}
}

func (c *Client) doJSON(ctx context.Context, path string, out any) error {
	<-c.limiter

	u := baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
