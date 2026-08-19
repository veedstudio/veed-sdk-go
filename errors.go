// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"fmt"
	"net/http"
	"time"
)

// ErrorDetail is one structured entry of an HTTP-layer error response.
type ErrorDetail struct {
	Type         string `json:"type"`
	Field        string `json:"field,omitempty"`
	Message      string `json:"message,omitempty"`
	Reason       string `json:"reason,omitempty"`
	RetryAfterMS int64  `json:"retry_after_ms,omitempty"`
	Value        any    `json:"value,omitempty"`
}

// Error is an HTTP-layer API error: the request itself was rejected. For
// submits this means no job was created (per the API contract). Match with
// errors.As and switch on StatusCode or Code.
type Error struct {
	StatusCode int
	// Code is the machine-readable API error code, e.g. "rate_limited".
	Code    string
	Message string
	// RequestID correlates this request in VEED's logs; include it in support requests.
	RequestID string
	Details   []ErrorDetail
	// RetryAfter is set from the Retry-After header on rate-limited responses.
	RetryAfter time.Duration
}

func (e *Error) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	if e.Code != "" {
		return fmt.Sprintf("veed: %d %s: %s (request_id: %s)", e.StatusCode, e.Code, msg, e.RequestID)
	}
	return fmt.Sprintf("veed: %d: %s", e.StatusCode, msg)
}

// JobFailedError is a job-layer failure: the job was accepted, but rendering
// failed. This is reported by Wait/Generate, not as an HTTP error. Code holds
// the model-specific failure code (compare against the generated
// <Model>JobErrorCode constants, e.g. FabricJobErrorCodeContentModeration).
type JobFailedError struct {
	JobID   string
	Code    string
	Message string
	Details []JobErrorDetail
}

func (e *JobFailedError) Error() string {
	return fmt.Sprintf("veed: job %s failed (%s): %s", e.JobID, e.Code, e.Message)
}

// JobCancelledError is returned by Wait/Generate when the job was cancelled.
type JobCancelledError struct {
	JobID string
}

func (e *JobCancelledError) Error() string {
	return fmt.Sprintf("veed: job %s was cancelled", e.JobID)
}

// WaitTimeoutError is returned when Wait/Generate gives up before the job
// reaches a terminal state. The job may still finish server-side; resume with
// Wait(ctx, JobID).
type WaitTimeoutError struct {
	JobID   string
	Timeout time.Duration
}

func (e *WaitTimeoutError) Error() string {
	return fmt.Sprintf("veed: timed out after %s waiting for job %s (the job may still finish; resume with Wait)", e.Timeout, e.JobID)
}
