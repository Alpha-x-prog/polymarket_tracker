package polymarket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
// 1) точечный фетч по slug
func (c *Client) GetMarketBySlug(ctx context.Context, slug string) (*Market, error) {
	var m Market

	path := "/markets/slug/" + url.PathEscape(slug) // ВАЖНО: /markets/slug/...

	if err := c.doJSON(ctx, path, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// 2) текстовый поиск — /public-search?q=...
type publicSearchResp struct {
	Events []struct {
		Markets []Market `json:"markets"`
	} `json:"events"`
}

// хелпер: отличить slug от свободного текста
func LooksLikeSlug(s string) bool {
	s = strings.TrimSpace(s)
	return s != "" && !strings.Contains(s, " ") && strings.Contains(s, "-")
}

// 1) открытые позиции
func (c *Client) GetUserPositions(ctx context.Context, addr string) ([]UserPosition, error) {
	var positions []UserPosition
	if err := c.doJSONData(ctx, "/positions?user="+addr, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

// 2) закрытые позиции
func (c *Client) GetUserClosedPositions(ctx context.Context, addr string, limit int) ([]ClosedPosition, error) {
	if limit <= 0 {
		limit = 50
	}
	var closed []ClosedPosition
	path := fmt.Sprintf("/closed-positions?user=%s&limit=%d", addr, limit)
	if err := c.doJSONData(ctx, path, &closed); err != nil {
		return nil, err
	}
	return closed, nil
}

// 3) total value
func (c *Client) GetUserTotalValue(ctx context.Context, addr string) (float64, error) {
	var resp []UserValue
	if err := c.doJSONData(ctx, "/value?user="+addr, &resp); err != nil {
		return 0, err
	}
	if len(resp) == 0 {
		return 0, nil
	}
	return resp[0].Value, nil
}

// 4) активность
func (c *Client) GetUserActivity(ctx context.Context, addr string, limit int) ([]Activity, error) {
	if limit <= 0 {
		limit = 10
	}
	var acts []Activity
	path := fmt.Sprintf("/activity?user=%s&limit=%d", addr, limit)
	if err := c.doJSONData(ctx, path, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

// 5) сколько рынков трогал
func (c *Client) GetUserTraded(ctx context.Context, addr string) (int, error) {
	var res UserTraded
	if err := c.doJSONData(ctx, "/traded?user="+addr, &res); err != nil {
		return 0, err
	}
	return res.Traded, nil
}

// если у тебя уже есть Client и doJSON — просто добавь эту функцию

func (c *Client) GetMarketByID(ctx context.Context, id string) (*Market, error) {
	var m Market
	if err := c.doJSON(ctx, "/markets/"+url.PathEscape(id), &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *Client) SearchMarkets(ctx context.Context, text string, limit int) ([]Market, error) {
	if limit <= 0 || limit > 25 {
		limit = 5
	}
	q := url.Values{}
	q.Set("limit", fmt.Sprint(limit))
	q.Set("text", text)

	var resp MarketsResponse
	if err := c.doJSON(ctx, "/markets?"+q.Encode(), &resp); err != nil {
		return nil, err
	}
	return resp.Markets, nil
}
