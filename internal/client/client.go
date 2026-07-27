package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
	"github.com/gravitino/terraform-provider-gravitino/internal/models"
)

const (
	defaultTimeout   = 30 * time.Second
	contentType      = "application/vnd.gravitino.v1+json"
	plainContentType = "application/json"
)

type Client struct {
	baseURL      string
	httpClient   *http.Client
	authProvider auth.AuthProvider
}

func New(uri string, authProvider auth.AuthProvider) (*Client, error) {
	u, err := url.Parse(strings.TrimRight(uri, "/"))
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	c := &Client{
		baseURL: u.String(),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		authProvider: authProvider,
	}

	if tp, ok := authProvider.(auth.TransportProvider); ok {
		c.httpClient.Transport = tp.WrapTransport(http.DefaultTransport)
	}

	return c, nil
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, c.baseURL+"/api"+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", contentType)
	if body != nil {
		req.Header.Set("Content-Type", plainContentType)
	}

	if c.authProvider != nil {
		key, value, err := c.authProvider.Header(context.Background())
		if err != nil {
			return nil, fmt.Errorf("auth header failed: %w", err)
		}
		if key != "" && value != "" {
			req.Header.Set(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) do(method, path string, body, result interface{}) error {
	resp, err := c.doRequest(method, path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		errResp.Code = resp.StatusCode
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return &errResp
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(path string, result interface{}) error {
	return c.do(http.MethodGet, path, nil, result)
}

func (c *Client) Post(path string, body, result interface{}) error {
	return c.do(http.MethodPost, path, body, result)
}

func (c *Client) Put(path string, body, result interface{}) error {
	return c.do(http.MethodPut, path, body, result)
}

func (c *Client) Delete(path string, result interface{}) error {
	return c.do(http.MethodDelete, path, nil, result)
}

func (c *Client) Patch(path string, body, result interface{}) error {
	return c.do(http.MethodPatch, path, body, result)
}

func (c *Client) Close() error {
	if closer, ok := c.authProvider.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (c *Client) BaseURL() string {
	return c.baseURL
}
