// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// jobFailure is the model-agnostic view of a job-level error, extracted from
// the model-specific error types by generated code.
type jobFailure struct {
	Code    string
	Message string
	Details []JobErrorDetail
}

func waitForJob[J any](
	ctx context.Context,
	cfg *requestConfig,
	jobID string,
	get func(context.Context) (*J, error),
	inspect func(*J) (JobStatus, *jobFailure),
) (*J, error) {
	deadline := time.Now().Add(cfg.waitTimeout)
	var last *J
	for {
		sleep := cfg.pollInterval
		job, err := get(ctx)
		if err != nil {
			// The job keeps rendering server-side, so a rate-limited poll is
			// not fatal to the wait; anything else is.
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests {
				return last, err
			}
			if apiErr.RetryAfter > 0 {
				sleep = apiErr.RetryAfter
			}
		} else {
			last = job
			status, failure := inspect(job)
			switch status {
			case JobStatusCompleted:
				return job, nil
			case JobStatusFailed:
				failed := &JobFailedError{JobID: jobID}
				if failure != nil {
					failed.Code = failure.Code
					failed.Message = failure.Message
					failed.Details = failure.Details
				}
				return job, failed
			case JobStatusCancelled:
				return job, &JobCancelledError{JobID: jobID}
			}
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return last, &WaitTimeoutError{JobID: jobID, Timeout: cfg.waitTimeout}
		}
		// Or a poll answered with Retry-After decides the wait's real length.
		if err := sleepCtx(ctx, min(sleep, remaining)); err != nil {
			return last, err
		}
	}
}
