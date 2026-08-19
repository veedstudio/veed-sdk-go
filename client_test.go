// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := NewClient(WithAPIKey("vp_test"), WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSubmitSendsAuthAndBody(t *testing.T) {
	var gotAuth, gotPath, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.Method + " " + r.URL.Path
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		gotBody = string(buf)
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
	})
	job, err := c.Fabric.Submit(context.Background(), FabricInput{
		ImageURL:   "https://in/i.png",
		AudioURL:   "https://in/a.mp3",
		Resolution: FabricResolution720p,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer vp_test" {
		t.Errorf("auth = %q", gotAuth)
	}
	if gotPath != "POST /v1/fabric-1.0" {
		t.Errorf("path = %q", gotPath)
	}
	want := `{"audio_url":"https://in/a.mp3","image_url":"https://in/i.png","resolution":"720p"}`
	if gotBody != want {
		t.Errorf("body = %q, want %q", gotBody, want)
	}
	if job.JobID != "j1" || job.Status != JobStatusProcessing {
		t.Errorf("job = %+v", job)
	}
}

func TestErrorMapping(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-1")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":"not_found","message":"no such job","request_id":"req-1"}}`))
	})
	_, err := c.Fabric.Get(context.Background(), "missing")
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 404 || apiErr.Code != "not_found" || apiErr.RequestID != "req-1" {
		t.Errorf("err = %+v", apiErr)
	}
}

func TestRetriesOn429ThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"slow down","request_id":"r"}}`))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
	})
	job, err := c.Lipsync20.Submit(context.Background(), Lipsync20Input{
		VideoURL: "https://in/v.mp4",
		AudioURL: "https://in/a.mp3",
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.JobID != "j1" {
		t.Errorf("job = %+v", job)
	}
	if got := calls.Load(); got != 2 {
		t.Errorf("calls = %d, want 2", got)
	}
}

func TestWaitPollsUntilCompleted(t *testing.T) {
	var gets atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if gets.Add(1) < 3 {
			_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"COMPLETED","result":{"video":{"url":"https://cdn/out.mp4"}}}}`))
	})
	job, err := c.Fabric.Wait(context.Background(), "j1", WithPollInterval(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if job.Result == nil || job.Result.Video.URL != "https://cdn/out.mp4" {
		t.Errorf("job = %+v", job)
	}
	if got := gets.Load(); got != 3 {
		t.Errorf("gets = %d, want 3", got)
	}
}

func TestWaitReturnsJobFailedError(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"FAILED","error":{"code":"content_moderation","message":"rejected"}}}`))
	})
	job, err := c.Fabric.Wait(context.Background(), "j1", WithPollInterval(time.Millisecond))
	var failed *JobFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("want *JobFailedError, got %T: %v", err, err)
	}
	if failed.Code != string(FabricJobErrorCodeContentModeration) || failed.JobID != "j1" {
		t.Errorf("failed = %+v", failed)
	}
	if job == nil || job.Status != JobStatusFailed {
		t.Errorf("job = %+v", job)
	}
}

func TestWaitTimesOut(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
	})
	_, err := c.Fabric.Wait(context.Background(), "j1",
		WithPollInterval(time.Millisecond), WithWaitTimeout(20*time.Millisecond))
	var timeout *WaitTimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("want *WaitTimeoutError, got %T: %v", err, err)
	}
	if timeout.JobID != "j1" {
		t.Errorf("timeout = %+v", timeout)
	}
}

func TestGenerateSubmitsAndWaits(t *testing.T) {
	var gets atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"data":{"job_id":"j9","status":"PROCESSING"}}`))
			return
		}
		if gets.Add(1) < 2 {
			_, _ = w.Write([]byte(`{"data":{"job_id":"j9","status":"PROCESSING"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"job_id":"j9","status":"COMPLETED","result":{"video":{"url":"https://cdn/out.mp4"}}}}`))
	})
	job, err := c.Lipsync20.Generate(context.Background(), Lipsync20Input{
		VideoURL: "https://in/v.mp4",
		AudioURL: "https://in/a.mp3",
	}, WithPollInterval(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if job.JobID != "j9" || job.Result == nil {
		t.Errorf("job = %+v", job)
	}
}

func TestMissingAPIKey(t *testing.T) {
	t.Setenv("VEED_API_KEY", "")
	if _, err := NewClient(); err == nil {
		t.Fatal("want error for missing API key")
	}
}

func TestClientDefaultHeadersAndPerCallOverride(t *testing.T) {
	var storeIO, mediaExp string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storeIO = r.Header.Get("X-Veed-Store-IO")
		mediaExp = r.Header.Get("X-Veed-Media-Expiration-Seconds")
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(WithAPIKey("vp_test"), WithBaseURL(srv.URL),
		WithDefaultStoreIO(false), WithDefaultMediaExpirationSeconds(600))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Fabric.Get(context.Background(), "j1"); err != nil {
		t.Fatal(err)
	}
	if storeIO != "0" || mediaExp != "600" {
		t.Errorf("client defaults: storeIO = %q, mediaExp = %q", storeIO, mediaExp)
	}

	if _, err := c.Fabric.Get(context.Background(), "j1", WithStoreIO(true), WithMediaExpirationSeconds(60)); err != nil {
		t.Fatal(err)
	}
	if storeIO != "1" || mediaExp != "60" {
		t.Errorf("per-call override: storeIO = %q, mediaExp = %q", storeIO, mediaExp)
	}
}

func TestRequestHeaders(t *testing.T) {
	var storeIO, mediaExp string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		storeIO = r.Header.Get("X-Veed-Store-IO")
		mediaExp = r.Header.Get("X-Veed-Media-Expiration-Seconds")
		_, _ = w.Write([]byte(`{"data":{"job_id":"j1","status":"PROCESSING"}}`))
	})
	_, err := c.Fabric.Get(context.Background(), "j1",
		WithStoreIO(false), WithMediaExpirationSeconds(3600))
	if err != nil {
		t.Fatal(err)
	}
	if storeIO != "0" || mediaExp != "3600" {
		t.Errorf("storeIO = %q, mediaExp = %q", storeIO, mediaExp)
	}
}

// Retry-After: 0 fails the `> 0` guard every caller applies, so the suite left
// the parser unexercised and it shipped without a ceiling.
func TestRetryAfterIsClamped(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "999999999")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"slow down","request_id":"r"}}`))
	})
	c.maxRetries = 0

	_, err := c.Fabric.Get(context.Background(), "j1")
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error, got %T: %v", err, err)
	}
	if want := maxRetryAfterSecs * time.Second; apiErr.RetryAfter != want {
		t.Errorf("RetryAfter = %s, want %s", apiErr.RetryAfter, want)
	}
}

// A wait must not outlive its timeout because a poll asked it to.
func TestWaitTimeoutSurvivesRetryAfter(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "600")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"slow down","request_id":"r"}}`))
	})
	c.maxRetries = 0

	// Backstop: without the clamp this fails in seconds instead of hanging.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Fabric.Wait(ctx, "j1", WithWaitTimeout(50*time.Millisecond))
	var timeout *WaitTimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("want *WaitTimeoutError, got %T: %v", err, err)
	}
}

// Following a redirect is how the key ends up on an http:// wire.
func TestRedirectIsNotFollowed(t *testing.T) {
	var leaked atomic.Bool
	elsewhere := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			leaked.Store(true)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer elsewhere.Close()

	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, elsewhere.URL+"/v1/fabric-1.0/j1", http.StatusFound)
	})

	_, err := c.Fabric.Get(context.Background(), "j1")
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *Error for a redirect, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusFound)
	}
	if leaked.Load() {
		t.Error("the redirect target received the Authorization header")
	}
}
