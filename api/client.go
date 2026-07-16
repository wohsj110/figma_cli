package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.figma.com/v1"

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	Verbose    bool
	VerboseOut io.Writer
}

func New(token string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		Token:      token,
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		VerboseOut: os.Stderr,
	}
}

func (c *Client) Me(ctx context.Context) (*Me, []byte, error) {
	var out Me
	raw, err := c.getJSON(ctx, "/me", &out)
	return &out, raw, err
}

func (c *Client) File(ctx context.Context, key string, depth int) (*File, []byte, error) {
	path := "/files/" + url.PathEscape(key)
	if depth > 0 {
		path += "?depth=" + url.QueryEscape(fmt.Sprintf("%d", depth))
	}
	var out File
	raw, err := c.getJSON(ctx, path, &out)
	return &out, raw, err
}

func (c *Client) Nodes(ctx context.Context, key string, ids []string) (*NodesResponse, []byte, error) {
	q := url.Values{}
	q.Set("ids", strings.Join(ids, ","))
	path := "/files/" + url.PathEscape(key) + "/nodes?" + q.Encode()
	var out NodesResponse
	raw, err := c.getJSON(ctx, path, &out)
	return &out, raw, err
}

func (c *Client) Images(ctx context.Context, key string, ids []string, format string, scale float64) (*ImagesResponse, []byte, error) {
	if format == "" {
		format = "png"
	}
	if scale <= 0 {
		scale = 1
	}
	q := url.Values{}
	q.Set("ids", strings.Join(ids, ","))
	q.Set("format", format)
	q.Set("scale", fmt.Sprintf("%g", scale))
	path := "/images/" + url.PathEscape(key) + "?" + q.Encode()
	var out ImagesResponse
	raw, err := c.getJSON(ctx, path, &out)
	return &out, raw, err
}

func (c *Client) Comments(ctx context.Context, key string) (*CommentsResponse, []byte, error) {
	path := "/files/" + url.PathEscape(key) + "/comments"
	var out CommentsResponse
	raw, err := c.getJSON(ctx, path, &out)
	return &out, raw, err
}

func (c *Client) Download(ctx context.Context, rawURL, outPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = io.Copy(f, resp.Body)
	return err
}

func (c *Client) getJSON(ctx context.Context, path string, out any) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Figma-Token", c.Token)
	req.Header.Set("Accept", "application/json")
	if c.Verbose {
		w := c.VerboseOut
		if w == nil {
			w = os.Stderr
		}
		fmt.Fprintf(w, "-> GET %s\n", req.URL.String())
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		w := c.VerboseOut
		if w == nil {
			w = os.Stderr
		}
		fmt.Fprintf(w, "<- %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("figma api error: %s: %s", resp.Status, truncate(body, 512))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}
	return body, nil
}

func truncate(b []byte, max int) string {
	s := strings.TrimSpace(string(b))
	if len(s) <= max {
		return s
	}
	return s[:max] + "...[truncated]"
}
