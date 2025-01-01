package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"mine.local/ocr-gallery/image-collector/api/ocrserver"
	"mine.local/ocr-gallery/image-collector/conf"
	"mine.local/ocr-gallery/image-collector/entity"
)

type OcrSerivce interface {
	DoOcr(ctx context.Context, imageMetadata *entity.StorageMetaData, imageData *string) (*entity.OcrResultBulk, error)
}

type OcrServiceImpl struct {
	ocrclient ocrserver.ClientWithResponsesInterface
}

func (ocr *OcrServiceImpl) DoOcr(
	ctx context.Context,
	imageMetadata *entity.StorageMetaData,
	imageData *string) (*entity.OcrResultBulk, error) {

	log.Printf("DoOcr request: imageName=%s", imageMetadata.Name)

	request := ocrserver.OcrRequestDto{
		ImageId:   &imageMetadata.Id,
		ImageName: &imageMetadata.Name,
		Image:     imageData,
	}

	response, err := ocr.ocrclient.PostApiV1OcrProcessWithResponse(ctx, request)
	if err != nil {
		return nil, wrapError(
			err,
			"request to ocr service failed: imageid=%s",
			imageMetadata.Id)

	}

	if response.StatusCode() != 200 {
		return nil, wrapError(
			err,
			"request to ocr service: imageid=%s status=%d",
			imageMetadata.Id,
			response.StatusCode())
	}

	responseJson := response.JSON200

	ocrResultItems := make([]entity.OcrResult, len(*responseJson.Texts))
	for i, item := range *responseJson.Texts {
		ocrResultItems[i] = entity.OcrResult{
			ProcessorKey: *item.ProcessorKey,
			Text:         *item.Text,
		}
	}

	retVal := entity.OcrResultBulk{
		Texts: &ocrResultItems,
	}

	return &retVal, nil
}

func wrapError(e error, text string, arg ...any) error {
	return errors.Join(fmt.Errorf(text, arg...), e)
}

func NewOcrService(conf *conf.OcrConfig) (OcrSerivce, error) {
	ocrServiceUrl := conf.Url
	log.Printf("Creating ocr service url=%s\n", ocrServiceUrl)

	client, err := ocrserver.NewClientWithResponses(ocrServiceUrl)
	if err != nil {
		return nil, err
	}
	return &OcrServiceImpl{
			ocrclient: client,
		},
		nil
}
