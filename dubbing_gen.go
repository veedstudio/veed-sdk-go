// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// DubbingJobErrorCode: stable machine-readable failure codes for dubbing jobs.
type DubbingJobErrorCode string

const (
	DubbingJobErrorCodeInputValidation     DubbingJobErrorCode = "input_validation"
	DubbingJobErrorCodeContentModeration   DubbingJobErrorCode = "content_moderation"
	DubbingJobErrorCodeInvalidFile         DubbingJobErrorCode = "invalid_file"
	DubbingJobErrorCodeAudioTooLong        DubbingJobErrorCode = "audio_too_long"
	DubbingJobErrorCodeTransloadFailed     DubbingJobErrorCode = "transload_failed"
	DubbingJobErrorCodeGenerationFailed    DubbingJobErrorCode = "generation_failed"
	DubbingJobErrorCodeTimeout             DubbingJobErrorCode = "timeout"
	DubbingJobErrorCodeInsufficientCredits DubbingJobErrorCode = "insufficient_credits"
)

// DubbingInput holds the inputs for a dubbing job.
type DubbingInput struct {
	// Language spoken in the video, as a BCP 47 tag. Omit it, or send auto, to detect it.
	SourceLanguage string `json:"source_language,omitempty"`
	// Language to dub into, as a BCP 47 tag.
	TargetLanguage string `json:"target_language"`
	// URL of the video to dub, up to 10 minutes.
	VideoURL string `json:"video_url"`
}

// Dubbing is the resource produced by a completed dubbing job.
type Dubbing struct {
	// The input video with its speech dubbed into the target language and the speakers' lips synced to it.
	Video File `json:"video"`
}

// DubbingJobError describes why a dubbing job FAILED.
type DubbingJobError struct {
	Code    DubbingJobErrorCode `json:"code"`
	Message string              `json:"message"`
	Details []JobErrorDetail    `json:"details,omitempty"`
}

// DubbingJob is the job envelope for the dubbing model.
type DubbingJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *Dubbing `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *DubbingJobError `json:"error,omitempty"`
}

const dubbingPath = "/v1/dubbing"

// Default Wait polling interval for dubbing (x-veed-poll-interval-seconds).
const dubbingPollInterval = 15 * time.Second

// DubbingService accesses the dubbing model.
type DubbingService struct {
	client *Client
}

// Submit starts an asynchronous dubbing job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *DubbingService) Submit(ctx context.Context, input DubbingInput, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[DubbingJob](ctx, s.client, http.MethodPost, dubbingPath, cfg, input)
}

// Get returns one snapshot of a dubbing job.
func (s *DubbingService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[DubbingJob](ctx, s.client, http.MethodGet, dubbingPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *DubbingService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = dubbingPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*DubbingJob, error) {
			return doResource[DubbingJob](ctx, s.client, http.MethodGet, dubbingPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *DubbingJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *DubbingService) Generate(ctx context.Context, input DubbingInput, opts ...RequestOption) (*DubbingJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
