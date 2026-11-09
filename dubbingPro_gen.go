// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// DubbingProInput holds the inputs for a dubbing-pro job.
type DubbingProInput struct {
	// Language spoken in the video, as a BCP 47 tag. Omit it, or send auto, to detect it.
	SourceLanguage string `json:"source_language,omitempty"`
	// Language to dub into, as a BCP 47 tag.
	TargetLanguage string `json:"target_language"`
	// URL of the video to dub, up to 10 minutes.
	VideoURL string `json:"video_url"`
	// How closely the dubbed voices match the speakers. Higher keeps more of each speaker's voice and accent; lower sounds more natural in the target language.
	VoiceCloningStrength int64 `json:"voice_cloning_strength,omitempty"`
}

const dubbingProPath = "/v1/dubbing-pro"

// Default Wait polling interval for dubbing-pro (x-veed-poll-interval-seconds).
const dubbingProPollInterval = 15 * time.Second

// DubbingProService accesses the dubbing-pro model.
type DubbingProService struct {
	client *Client
}

// Submit starts an asynchronous dubbing-pro job. It returns immediately with
// the job in status PROCESSING; rendering happens in the background.
func (s *DubbingProService) Submit(ctx context.Context, input DubbingProInput, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[DubbingJob](ctx, s.client, http.MethodPost, dubbingProPath, cfg, input)
}

// Get returns one snapshot of a dubbing-pro job.
func (s *DubbingProService) Get(ctx context.Context, jobID string, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	return doResource[DubbingJob](ctx, s.client, http.MethodGet, dubbingProPath+"/"+url.PathEscape(jobID), cfg, nil)
}

// Wait polls until the job reaches a terminal state. On COMPLETED it returns
// the job with its Result set; on FAILED or CANCELLED it returns the job
// alongside a JobFailedError or JobCancelledError; if the wait budget runs
// out it returns a WaitTimeoutError (the job may still finish — call Wait
// again to resume). See WithPollInterval and WithWaitTimeout.
func (s *DubbingProService) Wait(ctx context.Context, jobID string, opts ...RequestOption) (*DubbingJob, error) {
	cfg := newRequestConfig(opts)
	if cfg.pollInterval == 0 {
		cfg.pollInterval = dubbingProPollInterval
	}
	return waitForJob(ctx, cfg, jobID,
		func(ctx context.Context) (*DubbingJob, error) {
			return doResource[DubbingJob](ctx, s.client, http.MethodGet, dubbingProPath+"/"+url.PathEscape(jobID), cfg, nil)
		},
		func(j *DubbingJob) (JobStatus, *jobFailure) {
			if j.Error == nil {
				return j.Status, nil
			}
			return j.Status, &jobFailure{Code: string(j.Error.Code), Message: j.Error.Message, Details: j.Error.Details}
		})
}

// Generate is Submit followed by Wait: inputs in, finished job out.
func (s *DubbingProService) Generate(ctx context.Context, input DubbingProInput, opts ...RequestOption) (*DubbingJob, error) {
	job, err := s.Submit(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return s.Wait(ctx, job.JobID, opts...)
}
