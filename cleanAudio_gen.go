// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// CleanAudioOutputFormat: Container for the 48 kHz mono 16-bit output. FLAC is lossless at about half the size of WAV.
type CleanAudioOutputFormat string

const (
	CleanAudioOutputFormatFlac CleanAudioOutputFormat = "flac"
	CleanAudioOutputFormatWav  CleanAudioOutputFormat = "wav"
)

// CleanAudioJobErrorCode: stable machine-readable failure codes for clean-audio jobs.
type CleanAudioJobErrorCode string

const (
	CleanAudioJobErrorCodeInputValidation   CleanAudioJobErrorCode = "input_validation"
	CleanAudioJobErrorCodeContentModeration CleanAudioJobErrorCode = "content_moderation"
	CleanAudioJobErrorCodeInvalidFile       CleanAudioJobErrorCode = "invalid_file"
	CleanAudioJobErrorCodeAudioTooLong      CleanAudioJobErrorCode = "audio_too_long"
	CleanAudioJobErrorCodeTransloadFailed   CleanAudioJobErrorCode = "transload_failed"
	CleanAudioJobErrorCodeGenerationFailed  CleanAudioJobErrorCode = "generation_failed"
	CleanAudioJobErrorCodeTimeout           CleanAudioJobErrorCode = "timeout"
)

// CleanAudioInput holds the inputs for a clean-audio job.
type CleanAudioInput struct {
	// URL of the recording to clean: any audio or video file, up to 30 minutes and 512 MB. A video's audio track is used; multi-channel audio is mixed down to mono.
	AudioURL string `json:"audio_url"`
	// Set to false to skip loudness normalization and keep the input level.
	NormalizeLoudness *bool `json:"normalize_loudness,omitempty"`
	// Container for the 48 kHz mono 16-bit output. FLAC is lossless at about half the size of WAV.
	OutputFormat CleanAudioOutputFormat `json:"output_format,omitempty"`
	// How much of the original is allowed to remain under speech: the suppression floor is 1 - strength. Lower keeps more room tone behind the voice; silence between words is always fully cleaned.
	Strength string `json:"strength,omitempty"`
	// Integrated loudness of the output in LUFS (ITU-R BS.1770); true peak is capped at -1.1 dBTP. Ignored when normalize_loudness is false. A null reads as omitted: set normalize_loudness to false to skip normalization.
	TargetLufs string `json:"target_lufs,omitempty"`
}

// CleanAudio is the resource produced by a completed clean-audio job.
type CleanAudio struct {
	// The denoised recording: 48 kHz mono, the same duration as the input.
	Audio File `json:"audio"`
}

// CleanAudioJobError describes why a clean-audio job FAILED.
type CleanAudioJobError struct {
	Code    CleanAudioJobErrorCode `json:"code"`
	Message string                 `json:"message"`
	Details []JobErrorDetail       `json:"details,omitempty"`
}

// CleanAudioJob is the job envelope for the clean-audio model.
type CleanAudioJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
	// Present once the job is COMPLETED.
	Result *CleanAudio `json:"result,omitempty"`
	// Present once the job has FAILED.
	Error *CleanAudioJobError `json:"error,omitempty"`
}

const cleanAudioPath = "/v1/clean-audio"

// Default Wait polling interval for clean-audio (x-veed-poll-interval-seconds).
const cleanAudioPollInterval = 15 * time.Second

// CleanAudioService accesses the clean-audio model.
type CleanAudioService struct {
	client *Client
}

// Submit starts an asynchronous clean-audio job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *CleanAudioService) Submit(ctx context.Context, input CleanAudioInput, opts ...RequestOption) (*CleanAudioJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[CleanAudioJob](ctx, s.client, http.MethodPost, cleanAudioPath, cfg, input)
}

// Get returns one snapshot of a clean-audio job.
func (s *CleanAudioService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*CleanAudioJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[CleanAudioJob](ctx, s.client, http.MethodGet, cleanAudioPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *CleanAudioService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*CleanAudioJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = cleanAudioPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*CleanAudioJob, error) {
			return doResource[CleanAudioJob](ctx, s.client, http.MethodGet, cleanAudioPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *CleanAudioJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *CleanAudioService) Generate(ctx context.Context, input CleanAudioInput, opts ...RequestOption) (*CleanAudioJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
