// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

// JobStatus is the lifecycle state of an asynchronous generation job.
type JobStatus string

const (
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
	JobStatusCancelled  JobStatus = "CANCELLED"
)

// Terminal reports whether the job has finished, successfully or not.
func (s JobStatus) Terminal() bool {
	return s == JobStatusCompleted || s == JobStatusFailed || s == JobStatusCancelled
}

// File is a downloadable media file produced by a job.
type File struct {
	// The URL where the file can be downloaded from. Signed; expires after the
	// duration set via WithMediaExpirationSeconds (API default: one day).
	URL         string `json:"url"`
	ContentType string `json:"content_type,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
}

// JobErrorDetail is one structured entry of a job-level failure.
type JobErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	// Dotted path to the offending input, when applicable.
	Field string `json:"field,omitempty"`
}
