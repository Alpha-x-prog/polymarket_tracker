package polymarket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://gamma-api.polymarket.com"
const dataAPIBaseURL = "https://data-api.polymarket.com"

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

func (c *Client) doJSONData(ctx context.Context, path string, out any) error {
	<-c.limiter

	u := dataAPIBaseURL + path
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
func (c *Client) doDataAPIJSON(ctx context.Context, path string, out any) error {
	<-c.limiter

	u := dataAPIBaseURL + path
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

// GET https://data-api.polymarket.com/activity?user=0x...&limit=20
func (c *Client) GetUserActivity(ctx context.Context, addr string, limit int) ([]Activity, error) {
	if limit <= 0 {
		limit = 10
	}
	path := fmt.Sprintf("/activity?user=%s&limit=%d", addr, limit)

	var acts []Activity
	if err := c.doJSONData(ctx, path, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

// GetUserLastActivity — берём только ОДНО самое свежее действие
func (c *Client) GetUserLastActivity(ctx context.Context, addr string) (*Activity, error) {
	acts, err := c.GetUserActivity(ctx, addr, 1)
	if err != nil {
		return nil, err
	}
	if len(acts) == 0 {
		return nil, nil
	}
	return &acts[0], nil
}

// SearchMarkets searches markets by text (slug, question, part of name)
func (c *Client) SearchMarkets(ctx context.Context, term string, limit int) ([]Market, error) {
	if limit <= 0 || limit > 25 {
		limit = 5
	}
	q := url.Values{}
	q.Set("limit", fmt.Sprint(limit))
	q.Set("search", term)

	var resp MarketsResponse
	if err := c.doJSON(ctx, "/markets?"+q.Encode(), &resp); err != nil {
		return nil, err
	}
	if len(resp.Markets) > limit {
		resp.Markets = resp.Markets[:limit]
	}
	return resp.Markets, nil
}

// GetMarketByID gets a market by condition ID
func (c *Client) GetMarketByID(ctx context.Context, id string) (*Market, error) {
	q := url.Values{}
	q.Set("conditionId", id)

	var resp MarketsResponse
	if err := c.doJSON(ctx, "/markets?"+q.Encode(), &resp); err != nil {
		return nil, err
	}

	if len(resp.Markets) == 0 {
		return nil, nil
	}

	return &resp.Markets[0], nil
}

func (c *Client) GetUserPositions(ctx context.Context, addr string) ([]UserPosition, error) {
	var positions []UserPosition
	if err := c.doJSONData(ctx, "/positions?user="+addr, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

// GetUserTotalValue запрашивает "Get total value of a user's positions"
func (c *Client) GetUserTotalValue(ctx context.Context, addr string) (float64, error) {
	var resp []UserValue
	path := "/value?user=" + addr
	if err := c.doJSONData(ctx, path, &resp); err != nil {
		return 0, err
	}
	if len(resp) == 0 {
		return 0, nil
	}
	return resp[0].Value, nil
}
