package ffmpeg

import (
	"context"

	commonservice "github.com/weoses/memelo/common/service"
	"github.com/weoses/memelo/common/temp"
	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/storage-service/conf"
	"github.com/weoses/memelo/storage-service/ocr"
)

type video2FrameExtractorAdapter struct {
	adapter ffmpegJobAdapter
}

var _ ocr.Video2FrameExtractor = (*video2FrameExtractorAdapter)(nil)

func (a *video2FrameExtractorAdapter) ExtractOneFrame(ctx context.Context, video temp.Data) (temp.Data, error) {
	return a.adapter.run(ctx, video, func(inputS3Path string) *v1.SubmitFfmpegJobRequest {
		return &v1.SubmitFfmpegJobRequest{
			InputS3Path: inputS3Path,
			Action:      &v1.SubmitFfmpegJobRequest_ExtractThumbnail{ExtractThumbnail: &v1.ExtractThumbnailAction{}},
		}
	})
}

func NewVideo2FrameExtractorAdapter(cfg *conf.FfmpegServiceConfig, tmpDataService commonservice.TmpDataService) ocr.Video2FrameExtractor {
	return &video2FrameExtractorAdapter{
		adapter: newFfmpegJobAdapter(cfg, tmpDataService, "Video2FrameExtractorAdapter"),
	}
}
