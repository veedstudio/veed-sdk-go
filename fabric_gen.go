// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// FabricResolution: Output video resolution.
type FabricResolution string

const (
	FabricResolution720p FabricResolution = "720p"
	FabricResolution480p FabricResolution = "480p"
)

// FabricJobErrorCode: stable machine-readable failure codes for fabric-1.0 jobs.
type FabricJobErrorCode string

const (
	FabricJobErrorCodeInputValidation   FabricJobErrorCode = "input_validation"
	FabricJobErrorCodeContentModeration FabricJobErrorCode = "content_moderation"
	FabricJobErrorCodeInvalidFile       FabricJobErrorCode = "invalid_file"
	FabricJobErrorCodeAudioTooLong      FabricJobErrorCode = "audio_too_long"
	FabricJobErrorCodeTransloadFailed   FabricJobErrorCode = "transload_failed"
	FabricJobErrorCodeGenerationFailed  FabricJobErrorCode = "generation_failed"
	FabricJobErrorCodeTimeout           FabricJobErrorCode = "timeout"
)

// FabricInput holds the inputs for a fabric-1.0 job.
type FabricInput struct {
	// URL of the audio track to lip-sync to.
	AudioURL string `json:"audio_url"`
	// URL of the source image to animate.
	ImageURL string `json:"image_url"`
	// Output video resolution.
	Resolution FabricResolution `json:"resolution"`
}

// FabricVideo is the resource produced by a completed fabric-1.0 job.
type FabricVideo struct {
	// Generated lip-synced video.
	Video File `json:"video"`
}

// FabricJobError describes why a fabric-1.0 job FAILED.
type FabricJobError struct {
	Code    FabricJobErrorCode `json:"code"`
	Message string             `json:"message"`
	Details []JobErrorDetail   `json:"details,omitempty"`
}

// FabricJob is the job envelope for the fabric-1.0 model.
type FabricJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *FabricVideo `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *FabricJobError `json:"error,omitempty"`
}

const fabricPath = "/v1/fabric-1.0"

// Default Wait polling interval for fabric-1.0 (x-veed-poll-interval-seconds).
const fabricPollInterval = 10 * time.Second

// FabricService accesses the fabric-1.0 model.
type FabricService struct {
	client *Client
}

// Submit starts an asynchronous fabric-1.0 job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *FabricService) Submit(ctx context.Context, input FabricInput, opts ...RequestOption) (*FabricJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[FabricJob](ctx, s.client, http.MethodPost, fabricPath, cfg, input)
}

// Get returns one snapshot of a fabric-1.0 job.
func (s *FabricService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*FabricJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[FabricJob](ctx, s.client, http.MethodGet, fabricPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *FabricService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*FabricJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = fabricPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*FabricJob, error) {
			return doResource[FabricJob](ctx, s.client, http.MethodGet, fabricPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *FabricJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *FabricService) Generate(ctx context.Context, input FabricInput, opts ...RequestOption) (*FabricJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
