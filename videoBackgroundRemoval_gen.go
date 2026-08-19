// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// VideoBackgroundRemovalOutputCodec: Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
type VideoBackgroundRemovalOutputCodec string

const (
	VideoBackgroundRemovalOutputCodecVp9  VideoBackgroundRemovalOutputCodec = "vp9"
	VideoBackgroundRemovalOutputCodecH264 VideoBackgroundRemovalOutputCodec = "h264"
)

// VideoBackgroundRemovalJobErrorCode: stable machine-readable failure codes for video-background-removal jobs.
type VideoBackgroundRemovalJobErrorCode string

const (
	VideoBackgroundRemovalJobErrorCodeInputValidation   VideoBackgroundRemovalJobErrorCode = "input_validation"
	VideoBackgroundRemovalJobErrorCodeContentModeration VideoBackgroundRemovalJobErrorCode = "content_moderation"
	VideoBackgroundRemovalJobErrorCodeInvalidFile       VideoBackgroundRemovalJobErrorCode = "invalid_file"
	VideoBackgroundRemovalJobErrorCodeAudioTooLong      VideoBackgroundRemovalJobErrorCode = "audio_too_long"
	VideoBackgroundRemovalJobErrorCodeTransloadFailed   VideoBackgroundRemovalJobErrorCode = "transload_failed"
	VideoBackgroundRemovalJobErrorCodeGenerationFailed  VideoBackgroundRemovalJobErrorCode = "generation_failed"
	VideoBackgroundRemovalJobErrorCodeTimeout           VideoBackgroundRemovalJobErrorCode = "timeout"
)

// VideoBackgroundRemovalInput holds the inputs for a video-background-removal job.
type VideoBackgroundRemovalInput struct {
	// Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
	OutputCodec VideoBackgroundRemovalOutputCodec `json:"output_codec,omitempty"`
	// Improves the quality of the extracted subject's edges.
	RefineForegroundEdges *bool `json:"refine_foreground_edges,omitempty"`
	// Set to false when the subject is not a person.
	SubjectIsPerson *bool `json:"subject_is_person,omitempty"`
	// URL of the source video to remove the background from.
	VideoURL string `json:"video_url"`
}

// VideoBackgroundRemovalFiles is the resource produced by a completed video-background-removal job.
type VideoBackgroundRemovalFiles struct {
	// Rendered background-removed file(s): one webm with alpha for vp9; the RGB video and the alpha matte (two files) for h264.
	Files []File `json:"files"`
}

// VideoBackgroundRemovalJobError describes why a video-background-removal job FAILED.
type VideoBackgroundRemovalJobError struct {
	Code    VideoBackgroundRemovalJobErrorCode `json:"code"`
	Message string                             `json:"message"`
	Details []JobErrorDetail                   `json:"details,omitempty"`
}

// VideoBackgroundRemovalJob is the job envelope for the video-background-removal model.
type VideoBackgroundRemovalJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *VideoBackgroundRemovalFiles `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *VideoBackgroundRemovalJobError `json:"error,omitempty"`
}

const videoBackgroundRemovalPath = "/v1/video-background-removal"

// Default Wait polling interval for video-background-removal (x-veed-poll-interval-seconds).
const videoBackgroundRemovalPollInterval = 10 * time.Second

// VideoBackgroundRemovalService accesses the video-background-removal model.
type VideoBackgroundRemovalService struct {
	client *Client
}

// Submit starts an asynchronous video-background-removal job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *VideoBackgroundRemovalService) Submit(ctx context.Context, input VideoBackgroundRemovalInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodPost, videoBackgroundRemovalPath, cfg, input)
}

// Get returns one snapshot of a video-background-removal job.
func (s *VideoBackgroundRemovalService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *VideoBackgroundRemovalService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = videoBackgroundRemovalPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*VideoBackgroundRemovalJob, error) {
			return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *VideoBackgroundRemovalJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *VideoBackgroundRemovalService) Generate(ctx context.Context, input VideoBackgroundRemovalInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
