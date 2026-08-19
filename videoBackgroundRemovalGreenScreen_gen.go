// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// VideoBackgroundRemovalGreenScreenOutputCodec: Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
type VideoBackgroundRemovalGreenScreenOutputCodec string

const (
	VideoBackgroundRemovalGreenScreenOutputCodecVp9  VideoBackgroundRemovalGreenScreenOutputCodec = "vp9"
	VideoBackgroundRemovalGreenScreenOutputCodecH264 VideoBackgroundRemovalGreenScreenOutputCodec = "h264"
)

// VideoBackgroundRemovalGreenScreenInput holds the inputs for a video-background-removal-green-screen job.
type VideoBackgroundRemovalGreenScreenInput struct {
	// Output encoding. vp9 (the default) yields a single webm video with an alpha channel; h264 yields two files (the RGB video and an alpha matte) and is recommended for better RGB quality.
	OutputCodec VideoBackgroundRemovalGreenScreenOutputCodec `json:"output_codec,omitempty"`
	// How strongly green cast is removed from the kept subject. Raise it when green spots remain, lower it when colours shift on the subject.
	SpillSuppressionStrength string `json:"spill_suppression_strength,omitempty"`
	// URL of the source video, shot against a green screen.
	VideoURL string `json:"video_url"`
}

const videoBackgroundRemovalGreenScreenPath = "/v1/video-background-removal-green-screen"

// Default Wait polling interval for video-background-removal-green-screen (x-veed-poll-interval-seconds).
const videoBackgroundRemovalGreenScreenPollInterval = 10 * time.Second

// VideoBackgroundRemovalGreenScreenService accesses the video-background-removal-green-screen model.
type VideoBackgroundRemovalGreenScreenService struct {
	client *Client
}

// Submit starts an asynchronous video-background-removal-green-screen job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *VideoBackgroundRemovalGreenScreenService) Submit(ctx context.Context, input VideoBackgroundRemovalGreenScreenInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodPost, videoBackgroundRemovalGreenScreenPath, cfg, input)
}

// Get returns one snapshot of a video-background-removal-green-screen job.
func (s *VideoBackgroundRemovalGreenScreenService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalGreenScreenPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *VideoBackgroundRemovalGreenScreenService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = videoBackgroundRemovalGreenScreenPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*VideoBackgroundRemovalJob, error) {
			return doResource[VideoBackgroundRemovalJob](ctx, s.client, http.MethodGet, videoBackgroundRemovalGreenScreenPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *VideoBackgroundRemovalJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *VideoBackgroundRemovalGreenScreenService) Generate(ctx context.Context, input VideoBackgroundRemovalGreenScreenInput, opts ...RequestOption) (*VideoBackgroundRemovalJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
