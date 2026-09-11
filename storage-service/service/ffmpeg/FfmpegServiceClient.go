package ffmpeg

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/weoses/memelo/common/helper"
	commonservice "github.com/weoses/memelo/common/service"
	"github.com/weoses/memelo/common/temp"
	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"
	"github.com/weoses/memelo/storage-service/conf"
)

func newClient(cfg *conf.FfmpegServiceConfig) v1connect.FfmpegServiceClient {
	return v1connect.NewFfmpegServiceClient(http.DefaultClient, cfg.Uri)
}

// ffmpegJobAdapter holds the client/polling config shared by adapters that
// submit a single ffmpeg job and wrap its single S3 output as temp.Data.
type ffmpegJobAdapter struct {
	cl             v1connect.FfmpegServiceClient
	tmpDataService commonservice.TmpDataService
	pollInterval   time.Duration
	pollMaxWait    time.Duration
	log            *slog.Logger
	errPrefix      string
}

func newFfmpegJobAdapter(cfg *conf.FfmpegServiceConfig, tmpDataService commonservice.TmpDataService, errPrefix string) ffmpegJobAdapter {
	return ffmpegJobAdapter{
		cl:             newClient(cfg),
		tmpDataService: tmpDataService,
		pollInterval:   time.Duration(cfg.PollIntervalMs) * time.Millisecond,
		pollMaxWait:    time.Duration(cfg.PollMaxWaitSec) * time.Second,
		log:            slog.With("service", errPrefix),
		errPrefix:      errPrefix,
	}
}

// run submits a single ffmpeg job built by buildRequest, waits for it to
// finish, and wraps its single S3 output as temp.Data.
func (a *ffmpegJobAdapter) run(ctx context.Context, video temp.Data, buildRequest func(inputS3Path string) *v1.SubmitFfmpegJobRequest) (temp.Data, error) {
	input, err := resolveInput(ctx, a.tmpDataService, video)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", a.errPrefix, err)
	}
	if input.Owned {
		defer helper.QuietClose(input.Created, a.log)
	}

	submitResp, err := a.cl.SubmitJob(ctx, buildRequest(input.S3Path))
	if err != nil {
		return nil, fmt.Errorf("%s: submit job: %w", a.errPrefix, err)
	}

	status, err := pollJob(ctx, a.cl, submitResp.JobId, a.pollInterval, a.pollMaxWait)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", a.errPrefix, err)
	}

	result, err := a.tmpDataService.WrapInternalS3Path(ctx, status.GetOutputS3Path())
	if err != nil {
		return nil, fmt.Errorf("%s: wrap result: %w", a.errPrefix, err)
	}
	return result, nil
}

// pollJob polls GetJobStatus until the job reaches a terminal state, ctx is
// cancelled, or maxWait elapses.
func pollJob(
	ctx context.Context,
	cl v1connect.FfmpegServiceClient,
	jobId string,
	interval, maxWait time.Duration,
) (*v1.GetFfmpegJobStatusResponse, error) {
	deadline := time.Now().Add(maxWait)
	for {
		resp, err := cl.GetJobStatus(ctx, &v1.GetFfmpegJobStatusRequest{JobId: jobId})
		if err != nil {
			return nil, fmt.Errorf("pollJob: get status: %w", err)
		}

		switch resp.State {
		case v1.FfmpegJobState_FFMPEG_JOB_STATE_DONE:
			return resp, nil
		case v1.FfmpegJobState_FFMPEG_JOB_STATE_FAILED:
			return resp, fmt.Errorf("pollJob: job %s failed: %s", jobId, resp.GetError())
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("pollJob: job %s did not finish within %s", jobId, maxWait)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// resolvedInput is the S3 path an adapter should submit as a job's input,
// plus whether the adapter created that object and must delete it once the
// job is done.
type resolvedInput struct {
	S3Path  string
	Owned   bool
	Created temp.S3BackedData
}

// resolveInput returns an S3 path for video. If video is already S3-backed,
// its existing path is reused as-is — the adapter does not own it and must
// never delete it (the caller may read it again later, e.g. storage-service's
// pipeline reads the same VideoMp4 object from more than one step). Otherwise
// a fresh single-use object is uploaded, which the adapter does own and must
// close once the job finishes.
func resolveInput(ctx context.Context, tmpDataService commonservice.TmpDataService, video temp.Data) (resolvedInput, error) {
	if s3Data, ok := video.(temp.S3BackedData); ok {
		path, err := s3Data.GetS3Path(ctx)
		if err != nil {
			return resolvedInput{}, fmt.Errorf("resolveInput: get s3 path: %w", err)
		}
		return resolvedInput{S3Path: path}, nil
	}

	reader, err := video.Reader()
	if err != nil {
		return resolvedInput{}, fmt.Errorf("resolveInput: get reader: %w", err)
	}
	defer helper.QuietClose(reader, slog.With("func", "resolveInput"))

	uploaded, err := tmpDataService.ByReaderUpload(ctx, "video/mp4", reader)
	if err != nil {
		return resolvedInput{}, fmt.Errorf("resolveInput: upload: %w", err)
	}
	path, err := uploaded.GetS3Path(ctx)
	if err != nil {
		return resolvedInput{}, fmt.Errorf("resolveInput: get uploaded s3 path: %w", err)
	}
	return resolvedInput{S3Path: path, Owned: true, Created: uploaded}, nil
}
