// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import "time"

// RequestOption customizes a single API call.
type RequestOption func(*requestConfig)

type requestConfig struct {
	storeIO                *bool
	mediaExpirationSeconds *int
	pollInterval           time.Duration
	waitTimeout            time.Duration
}

// pollInterval zero means "use the model's spec-defined default"; the
// generated Wait methods resolve it.
func newRequestConfig(opts []RequestOption) *requestConfig {
	cfg := &requestConfig{
		waitTimeout: 15 * time.Minute,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithStoreIO(false) asks the API not to store this request's and response's
// bodies in your workspace request logs (X-Veed-Store-IO: 0).
func WithStoreIO(store bool) RequestOption {
	return func(cfg *requestConfig) { cfg.storeIO = &store }
}

// WithMediaExpirationSeconds sets how long signed media URLs returned for this
// request stay valid (X-Veed-Media-Expiration-Seconds; API default: one day).
func WithMediaExpirationSeconds(seconds int) RequestOption {
	return func(cfg *requestConfig) { cfg.mediaExpirationSeconds = &seconds }
}

// WithPollInterval sets the delay between status polls in Wait and Generate.
// Default: the model's spec-defined interval (x-veed-poll-interval-seconds).
// A Retry-After from a rate-limited poll overrides it.
func WithPollInterval(d time.Duration) RequestOption {
	return func(cfg *requestConfig) {
		if d > 0 {
			cfg.pollInterval = d
		}
	}
}

// WithWaitTimeout caps how long Wait and Generate poll before returning a
// WaitTimeoutError. Default: 15m. A context deadline also applies if set.
func WithWaitTimeout(d time.Duration) RequestOption {
	return func(cfg *requestConfig) {
		if d > 0 {
			cfg.waitTimeout = d
		}
	}
}
