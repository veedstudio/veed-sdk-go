// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

// Package veed is the official Go SDK for the VEED API — programmatic access
// to VEED's AI video models (https://api.veed.io).
//
// Construct a client once and reuse it:
//
//	client, err := veed.NewClient() // reads VEED_API_KEY
//	job, err := client.Fabric.Generate(ctx, veed.FabricInput{
//		ImageURL:   "https://example.com/face.png",
//		AudioURL:   "https://example.com/speech.mp3",
//		Resolution: veed.FabricResolution720p,
//	})
//	fmt.Println(job.Result.Video.URL)
package veed

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Version is the SDK version, sent in the User-Agent header.
const Version = "0.1.0"

const defaultBaseURL = "https://api.veed.io"

// The server chooses Retry-After and no request timeout covers a sleep, so an
// unclamped header is a remote hang.
const maxRetryAfterSecs = 60

// Client is the VEED API client. It is safe for concurrent use.
type Client struct {
	services

	apiKey                 string
	baseURL                *url.URL
	httpClient             *http.Client
	maxRetries             int
	userAgent              string
	storeIO                *bool
	mediaExpirationSeconds *int
}

// Option configures a Client.
type Option func(*Client) error

// WithAPIKey sets the workspace API key (vp_...). Defaults to $VEED_API_KEY.
func WithAPIKey(key string) Option {
	return func(c *Client) error {
		c.apiKey = key
		return nil
	}
}

// WithBaseURL overrides the API base URL. Defaults to $VEED_BASE_URL or
// https://api.veed.io.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) error {
		u, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("veed: invalid base URL %q: %w", rawURL, err)
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient replaces the underlying HTTP client (default: 60s timeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("veed: http client must not be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithMaxRetries sets how many times a failed request is retried (default 2).
// Retried statuses: 408, 429 and 5xx — all safe for submits too, since a
// rejected submit never creates a job.
func WithMaxRetries(n int) Option {
	return func(c *Client) error {
		if n < 0 {
			return errors.New("veed: max retries must be >= 0")
		}
		c.maxRetries = n
		return nil
	}
}

// WithDefaultStoreIO sets a client-wide default for whether the API stores
// request/response bodies in your request logs (X-Veed-Store-IO). Override
// per call with WithStoreIO.
func WithDefaultStoreIO(store bool) Option {
	return func(c *Client) error {
		c.storeIO = &store
		return nil
	}
}

// WithDefaultMediaExpirationSeconds sets a client-wide default for how long
// signed media URLs stay valid (X-Veed-Media-Expiration-Seconds). Override
// per call with WithMediaExpirationSeconds.
func WithDefaultMediaExpirationSeconds(seconds int) Option {
	return func(c *Client) error {
		c.mediaExpirationSeconds = &seconds
		return nil
	}
}

// NewClient returns a configured client, or an error if no API key is
// available from options or $VEED_API_KEY.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		apiKey: os.Getenv("VEED_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			// Go's default policy keeps Authorization for the same host and any
			// subdomain, and ignores the scheme — a 302 to http:// would put the
			// key on the wire. The API does not redirect.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		maxRetries: 2,
		userAgent:  "veed-sdk-go/" + Version,
	}
	base := defaultBaseURL
	if v := os.Getenv("VEED_BASE_URL"); v != "" {
		base = v
	}
	if err := WithBaseURL(base)(c); err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.apiKey == "" {
		return nil, errors.New("veed: missing API key: pass veed.WithAPIKey(...) or set VEED_API_KEY (create keys in your VEED workspace)")
	}
	c.initServices()
	return c, nil
}

// resource is the {"data": ...} envelope every 2xx response is wrapped in.
type resource[T any] struct {
	Data *T `json:"data"`
}

func doResource[T any](ctx context.Context, c *Client, method, path string, cfg *requestConfig, body any) (*T, error) {
	var env resource[T]
	if err := c.do(ctx, method, path, cfg, body, &env); err != nil {
		return nil, err
	}
	if env.Data == nil {
		return nil, fmt.Errorf("veed: response for %s %s is missing the data envelope", method, path)
	}
	return env.Data, nil
}

func (c *Client) do(ctx context.Context, method, path string, cfg *requestConfig, body, out any) error {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("veed: encode request: %w", err)
		}
	}
	u := c.baseURL.JoinPath(path)

	for attempt := 0; ; attempt++ {
		var bodyReader io.Reader
		if payload != nil {
			bodyReader = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
		if err != nil {
			return fmt.Errorf("veed: build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		storeIO := c.storeIO
		mediaExpiration := c.mediaExpirationSeconds
		if cfg != nil {
			if cfg.storeIO != nil {
				storeIO = cfg.storeIO
			}
			if cfg.mediaExpirationSeconds != nil {
				mediaExpiration = cfg.mediaExpirationSeconds
			}
		}
		if storeIO != nil {
			v := "1"
			if !*storeIO {
				v = "0"
			}
			req.Header.Set("X-Veed-Store-IO", v)
		}
		if mediaExpiration != nil {
			req.Header.Set("X-Veed-Media-Expiration-Seconds", strconv.Itoa(*mediaExpiration))
		}

		resp, err := c.httpClient.Do(req) //nolint:gosec // G704: the target is the configured API base URL — calling user-supplied hosts is the SDK's purpose
		if err != nil {
			if attempt < c.maxRetries && ctx.Err() == nil {
				if err := sleepCtx(ctx, backoff(attempt, 0)); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("veed: request failed: %w", err)
		}

		if resp.StatusCode < 300 {
			data, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				return fmt.Errorf("veed: read response: %w", err)
			}
			if out == nil {
				return nil
			}
			if err := json.Unmarshal(data, out); err != nil {
				return fmt.Errorf("veed: decode response: %w", err)
			}
			return nil
		}

		apiErr := parseError(resp)
		if retryableStatus(resp.StatusCode) && attempt < c.maxRetries {
			if err := sleepCtx(ctx, backoff(attempt, apiErr.RetryAfter)); err != nil {
				return err
			}
			continue
		}
		return apiErr
	}
}

type errorBody struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	RequestID string        `json:"request_id"`
	Details   []ErrorDetail `json:"details"`
}

type errorResponse struct {
	Error *errorBody `json:"error"`
}

func parseError(resp *http.Response) *Error {
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	apiErr := &Error{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Request-ID"),
	}
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		// Clamped before scaling: near MaxInt64 the multiplication overflows.
		if secs, err := strconv.Atoi(ra); err == nil && secs >= 0 {
			apiErr.RetryAfter = time.Duration(min(secs, maxRetryAfterSecs)) * time.Second
		}
	}
	var body errorResponse
	if err := json.Unmarshal(data, &body); err == nil && body.Error != nil {
		apiErr.Code = body.Error.Code
		apiErr.Message = body.Error.Message
		apiErr.Details = body.Error.Details
		if body.Error.RequestID != "" {
			apiErr.RequestID = body.Error.RequestID
		}
	}
	return apiErr
}

func retryableStatus(code int) bool {
	return code == http.StatusRequestTimeout || code == http.StatusTooManyRequests || code >= 500
}

// backoff returns the delay before the given retry attempt: Retry-After when
// the server provided one, otherwise full-jitter exponential backoff
// (base 500ms, doubling, capped at 8s).
func backoff(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}
	max := 500 * time.Millisecond << attempt
	if max > 8*time.Second {
		max = 8 * time.Second
	}
	return time.Duration(rand.Int64N(int64(max) + 1)) //nolint:gosec // G404: jitter, not cryptography
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
