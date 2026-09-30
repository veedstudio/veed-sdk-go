// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// VideoBackgroundRemovalFastOutputCodec: Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
type VideoBackgroundRemovalFastOutputCodec string

const (
	VideoBackgroundRemovalFastOutputCodecVp9  VideoBackgroundRemovalFastOutputCodec = "vp9"
	VideoBackgroundRemovalFastOutputCodecH264 VideoBackgroundRemovalFastOutputCodec = "h264"
)

// VideoBackgroundRemovalFastInput holds the inputs for a video-background-removal-fast job.
type VideoBackgroundRemovalFastInput struct {
	// Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
	OutputCodec VideoBackgroundRemovalFastOutputCodec `json:"output_codec,omitempty"`
	// Improves the quality of the extracted subject's edges.
	RefineForegroundEdges *bool `json:"refine_foreground_edges,omitempty"`
	// Set to false when the subject is not a person.
	SubjectIsPerson *bool `json:"subject_is_person,omitempty"`
	// URL of the source video to remove the background from.
	VideoURL string `json:"video_url"`
}

const videoBackgroundRemovalFastPath = "/v1/video-background-removal-fast"

// Default Wait polling interval for video-background-removal-fast (x-veed-poll-interval-seconds).
const videoBackgroundRemovalFastPollInterval = 10 * time.Second

// VideoBackgroundRemovalFastService accesses the video-background-removal-fast model.
type VideoBackgroundRemovalFastService struct {
	client *Client
}

// Submit starts an asynchronous video-background-removal-fast job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *VideoBackgroundRemovalFastService) Submit(ctx context.Context, input VideoBackgroundRemovalFastInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodPost, videoBackgroundRemovalFastPath, cfg, input)
}

// Get returns one snapshot of a video-background-removal-fast job.
func (s *VideoBackgroundRemovalFastService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalFastPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *VideoBackgroundRemovalFastService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = videoBackgroundRemovalFastPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*VideoBackgroundRemovalJob, error) {
			return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalFastPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *VideoBackgroundRemovalJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *VideoBackgroundRemovalFastService) Generate(ctx context.Context, input VideoBackgroundRemovalFastInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
