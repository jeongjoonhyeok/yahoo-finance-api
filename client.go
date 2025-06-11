package yahoofinanceapi

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// RotateThreshold defines how many requests are allowed
// before we force-refresh cookies and crumb.
const RotateThreshold int64 = 1000

type Client struct {
	client    *http.Client
	cookies   []*http.Cookie
	crumb     string
	callCount int64
	mu        sync.Mutex
}

var instance *Client
var once sync.Once

func GetClient() *Client {
	once.Do(func() {
		instance = &Client{client: &http.Client{}}
	})
	return instance
}

// Get is the public entry. It automatically rotates session
// every RotateThreshold calls to avoid Yahoo 429 limits.
func (c *Client) Get(url string, params url.Values) (*http.Response, error) {
	c.maybeRotateSession()
	c.getCrumb()
	return c.get(url, params)
}

// maybeRotateSession increments the counter and clears
// cookie / crumb once the threshold is reached.
func (c *Client) maybeRotateSession() {
	c.mu.Lock()
	c.callCount++
	if c.callCount%RotateThreshold == 0 {
		slog.Info("rotating Yahoo Finance session", "calls", c.callCount)
		c.cookies = nil
		c.crumb = ""
	}
	c.mu.Unlock()
}

func (c *Client) get(url string, params url.Values) (*http.Response, error) {
	if c.crumb != "" {
		params.Add("crumb", c.crumb)
	}
	url = fmt.Sprintf("%s?%s", url, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "err", err)
		return nil, err
	}

	// attach cookies
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// realistic browser headers
	req.Header.Set("User-Agent", USER_AGENTS[rand.Intn(len(USER_AGENTS))])
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ko;q=0.8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://finance.yahoo.com/")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to get data from Yahoo Finance API", "err", err)
		return nil, err
	}

	return resp, nil
}

// getCookie fetches fresh cookies if none are cached.
func (c *Client) getCookie() {
	if len(c.cookies) > 0 {
		return
	}

	endpoint := "https://fc.yahoo.com"
	resp, err := c.get(endpoint, url.Values{})
	if err != nil {
		slog.Error("Failed to get cookie", "err", err)
		return
	}

	c.cookies = resp.Cookies()
}

// getCrumb fetches crumb lazily.
func (c *Client) getCrumb() {
	if c.crumb != "" {
		return
	}

	c.getCookie()
	endpoint := fmt.Sprintf("%s/v1/test/getcrumb", BASE_URL)
	resp, err := c.get(endpoint, url.Values{})
	if err != nil {
		slog.Error("Failed to get crumb", "err", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "err", err)
		return
	}

	c.crumb = string(body)
}
