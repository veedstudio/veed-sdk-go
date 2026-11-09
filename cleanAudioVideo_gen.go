// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// CleanAudioVideoInput holds the inputs for a clean-audio-video job.
type CleanAudioVideoInput struct {
	// Set to false to skip loudness normalization and keep the input level.
	NormalizeLoudness *bool `json:"normalize_loudness,omitempty"`
	// How much of the original is allowed to remain under speech: the suppression floor is 1 - strength. Lower keeps more room tone behind the voice; silence between words is always fully cleaned.
	Strength string `json:"strength,omitempty"`
	// Integrated loudness of the output in LUFS (ITU-R BS.1770); true peak is capped at -1.1 dBTP. Ignored when normalize_loudness is false. A null reads as omitted: set normalize_loudness to false to skip normalization.
	TargetLufs string `json:"target_lufs,omitempty"`
	// URL of the video to clean: MP4, MOV, MKV or WebM with an H.264, HEVC, AV1, MPEG-4, VP8 or VP9 video stream and an audio track, up to 30 minutes and 2 GB. Only the first audio track is kept.
	VideoURL string `json:"video_url"`
}

// CleanAudioVideo is the resource produced by a completed clean-audio-video job.
type CleanAudioVideo struct {
	// The input video with its audio track replaced by the denoised one. The picture is copied unchanged; the audio is 48 kHz mono, AAC in an MP4 for H.264, HEVC, AV1 or MPEG-4 input, Opus in a WebM for VP8 or VP9.
	Video File `json:"video"`
}

// CleanAudioVideoJob is the job envelope for the clean-audio-video model.
type CleanAudioVideoJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *CleanAudioVideo `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *CleanAudioJobError `json:"error,omitempty"`
}

const cleanAudioVideoPath = "/v1/clean-audio-video"

// Default Wait polling interval for clean-audio-video (x-veed-poll-interval-seconds).
const cleanAudioVideoPollInterval = 15 * time.Second

// CleanAudioVideoService accesses the clean-audio-video model.
type CleanAudioVideoService struct {
	client *Client
}

// Submit starts an asynchronous clean-audio-video job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *CleanAudioVideoService) Submit(ctx context.Context, input CleanAudioVideoInput, opts ...RequestOption) (*CleanAudioVideoJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[CleanAudioVideoJob](ctx, s.client, http.MethodPost, cleanAudioVideoPath, cfg, input)
}

// Get returns one snapshot of a clean-audio-video job.
func (s *CleanAudioVideoService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*CleanAudioVideoJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[CleanAudioVideoJob](ctx, s.client, http.MethodGet, cleanAudioVideoPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *CleanAudioVideoService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*CleanAudioVideoJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = cleanAudioVideoPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*CleanAudioVideoJob, error) {
			return doResource[CleanAudioVideoJob](ctx, s.client, http.MethodGet, cleanAudioVideoPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *CleanAudioVideoJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *CleanAudioVideoService) Generate(ctx context.Context, input CleanAudioVideoInput, opts ...RequestOption) (*CleanAudioVideoJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
