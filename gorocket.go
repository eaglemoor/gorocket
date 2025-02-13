package gorocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	userID     string
	xToken     string
	apiVersion string
	HTTPClient *http.Client

	timeout time.Duration
}

// NewClient creates new Facest.io client with given API key
func NewClient(url string) *Client {
	return &Client{
		//userID: user,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		//xToken:     token,
		baseURL:    url,
		apiVersion: "api/v1",
	}
}

// NewClient creates new Facest.io client with given API key
func NewWithOptions(url string, opts ...Option) *Client {
	c := &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL:    url,
		apiVersion: "api/v1",
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

func WithUserID(userID string) Option {
	return func(c *Client) {
		c.userID = userID
	}
}

func WithXToken(xtoken string) Option {
	return func(c *Client) {
		c.xToken = xtoken
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("X-Auth-Token", c.xToken)
	req.Header.Add("X-User-Id", c.userID)

	if c.timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), c.timeout)
		defer cancel()

		req = req.WithContext(ctx)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer res.Body.Close()

	resp := v
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return err
	}

	return nil
}
