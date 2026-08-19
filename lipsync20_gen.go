// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Lipsync20JobErrorCode: stable machine-readable failure codes for lipsync-2.0 jobs.
type Lipsync20JobErrorCode string

const (
	Lipsync20JobErrorCodeInputValidation   Lipsync20JobErrorCode = "input_validation"
	Lipsync20JobErrorCodeContentModeration Lipsync20JobErrorCode = "content_moderation"
	Lipsync20JobErrorCodeInvalidFile       Lipsync20JobErrorCode = "invalid_file"
	Lipsync20JobErrorCodeAudioTooLong      Lipsync20JobErrorCode = "audio_too_long"
	Lipsync20JobErrorCodeTransloadFailed   Lipsync20JobErrorCode = "transload_failed"
	Lipsync20JobErrorCodeGenerationFailed  Lipsync20JobErrorCode = "generation_failed"
	Lipsync20JobErrorCodeTimeout           Lipsync20JobErrorCode = "timeout"
)

// Lipsync20Input holds the inputs for a lipsync-2.0 job.
type Lipsync20Input struct {
	// URL of the new audio track to lip-sync the video to.
	AudioURL string `json:"audio_url"`
	// URL of the source video to dub.
	VideoURL string `json:"video_url"`
}

// Lipsync20Video is the resource produced by a completed lipsync-2.0 job.
type Lipsync20Video struct {
	// Generated re-lip-synced video.
	Video File `json:"video"`
}

// Lipsync20JobError describes why a lipsync-2.0 job FAILED.
type Lipsync20JobError struct {
	Code    Lipsync20JobErrorCode `json:"code"`
	Message string                `json:"message"`
	Details []JobErrorDetail      `json:"details,omitempty"`
}

// Lipsync20Job is the job envelope for the lipsync-2.0 model.
type Lipsync20Job struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *Lipsync20Video `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *Lipsync20JobError `json:"error,omitempty"`
}

const lipsync20Path = "/v1/lipsync-2.0"

// Default Wait polling interval for lipsync-2.0 (x-veed-poll-interval-seconds).
const lipsync20PollInterval = 10 * time.Second

// Lipsync20Service accesses the lipsync-2.0 model.
type Lipsync20Service struct {
	client *Client
}

// Submit starts an asynchronous lipsync-2.0 job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *Lipsync20Service) Submit(ctx context.Context, input Lipsync20Input, opts ...RequestOption) (*Lipsync20Job, error) {
	cfg := newRequestConfig(opts)
	return doResource[Lipsync20Job](ctx, s.client, http.MethodPost, lipsync20Path, cfg, input)
}

// Get returns one snapshot of a lipsync-2.0 job.
func (s *Lipsync20Service) Get(ctx context.Context, jobID string, opts ...RequestOption) (*Lipsync20Job, error) {
	cfg := newRequestConfig(opts)
	return doResource[Lipsync20Job](ctx, s.client, http.MethodGet, lipsync20Path+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *Lipsync20Service) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*Lipsync20Job, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = lipsync20PollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*Lipsync20Job, error) {
			return doResource[Lipsync20Job](ctx, s.client, http.MethodGet, lipsync20Path+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *Lipsync20Job) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *Lipsync20Service) Generate(ctx context.Context, input Lipsync20Input, opts ...RequestOption) (*Lipsync20Job, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
