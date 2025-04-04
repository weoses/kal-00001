package service

import (
	"context"

	"github.com/h2non/bimg"
	"mine.local/ocr-gallery/ocr-server/conf"
	"mine.local/ocr-gallery/ocr-server/entity"
)

type ImageConveter interface {
	ConvertImage(ctx context.Context, image *entity.Image) (*entity.Image, error)
	MakeThumb(ctx context.Context, image *entity.Image) (*entity.Image, *entity.ImageSizes, error)
}

type ImageConveterImpl struct {
	config *conf.ImageConverterConfig
}

// ConvertImage implements ImageConveter.
func (i *ImageConveterImpl) ConvertImage(ctx context.Context, image *entity.Image) (*entity.Image, error) {
	img := bimg.NewImage(*image.Data)
	bytesData, err := img.Convert(bimg.JPEG)
	if err != nil {
		return nil, err
	}

	retval := new(entity.Image)
	retval.MimeType = "image/jpeg"
	retval.Data = &bytesData
	return retval, nil
}

// MakeThumb implements ImageConveter.
func (i *ImageConveterImpl) MakeThumb(ctx context.Context, image *entity.Image) (*entity.Image, *entity.ImageSizes, error) {
	img := bimg.NewImage(*image.Data)

	size, err := img.Size()
	if err != nil {
		return nil, nil, err
	}

	sizes := new(entity.ImageSizes)

	sizes.Width = i.config.ThumbSize
	sizes.Height = int(float64(i.config.ThumbSize) / float64(size.Width) * float64(size.Height))

	bytesData, err := img.Resize(size.Width, size.Height)
	if err != nil {
		return nil, nil, err
	}

	retval := new(entity.Image)
	retval.MimeType = image.MimeType
	retval.Data = &bytesData
	return retval, sizes, nil
}

func NewImageConverter(config *conf.ImageConverterConfig) (ImageConveter, error) {
	return &ImageConveterImpl{
		config: config,
	}, nil
}
