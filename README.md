# veed-sdk-go

Official Go SDK for the [VEED API](https://api.veed.io) — programmatic access to VEED's AI video models.

> Generated from VEED's OpenAPI spec. Do not edit by hand — changes are overwritten on the next release.

## Install

```bash
go get github.com/veedstudio/veed-sdk-go
```

## Usage

```go
client, err := veed.NewClient() // reads VEED_API_KEY
if err != nil {
	log.Fatal(err)
}

// One-shot: submit and wait for the finished video.
job, err := client.Fabric.Generate(ctx, veed.FabricInput{
	ImageURL:   "https://example.com/face.png",
	AudioURL:   "https://example.com/speech.mp3",
	Resolution: veed.FabricResolution720p,
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(job.Result.Video.URL)
```

Or keep control of the lifecycle:

```go
job, _ := client.Lipsync20.Submit(ctx, veed.Lipsync20Input{
	VideoURL: "https://example.com/source.mp4",
	AudioURL: "https://example.com/dub.mp3",
})
// ... persist job.JobID, come back later ...
job, err = client.Lipsync20.Wait(ctx, job.JobID,
	veed.WithPollInterval(5*time.Second),
	veed.WithWaitTimeout(30*time.Minute),
)
```

### Errors

HTTP-layer rejections (invalid key, bad inputs, rate limits) are `*veed.Error` with
`StatusCode`, `Code` and `RequestID`. Accepted jobs that fail during rendering are reported by
`Wait`/`Generate` as `*veed.JobFailedError` with a model-specific failure code:

```go
var failed *veed.JobFailedError
if errors.As(err, &failed) && failed.Code == string(veed.FabricJobErrorCodeContentModeration) {
	// handle moderation rejection
}
```

Retries (429/5xx, honoring `Retry-After`) are built in; submits are safe to retry because a
rejected submit never creates a job.

## License

MIT — see [LICENSE](LICENSE).
