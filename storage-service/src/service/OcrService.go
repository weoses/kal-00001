package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"mine.local/ocr-gallery/apispec/ocr-server/client"
	"mine.local/ocr-gallery/storage-service/conf"
	"mine.local/ocr-gallery/storage-service/entity"
	"mine.local/ocr-gallery/storage-service/helper"
)

type OcrSerivce interface {
	DoOcr(ctx context.Context,
		id uuid.UUID,
		incomingImage *entity.Image,
	) (*entity.OcrProcessedResult, error)
}

type OcrServiceImpl struct {
	ocrclient client.ClientWithResponsesInterface
}

func (ocr *OcrServiceImpl) DoOcr(
	ctx context.Context,
	id uuid.UUID,
	incomingImage *entity.Image,
) (*entity.OcrProcessedResult, error) {

	idStr := id.String()

	request := client.OcrRequestDto{
		ImageId: &idStr,
		Image: &client.ImageDto{
			ImageBase64: incomingImage.ImageBase64,
			MimeType:    &incomingImage.MimeType,
		},
	}

	response, err := ocr.ocrclient.PostApiV1OcrProcessWithResponse(ctx, request)
	if err != nil {
		return nil, err

	}

	if response.StatusCode() != 200 {
		return nil, errors.New("status code fault")
	}

	responseJson := response.JSON200

	textVariants := responseJson.ImageText
	thumbnail := responseJson.ImageThumb
	image := responseJson.Image

	retval := new(entity.OcrProcessedResult)
	retval.OcrText = textVariantsToString(textVariants)
	retval.Thumbnail = helper.ImageToEntity(thumbnail)
	retval.Image = helper.ImageToEntity(image)

	return retval, nil
}

func textVariantsToString(textVariants *[]client.OcrResponseItem) string {
	builder := strings.Builder{}
	for _, textVariant := range *textVariants {
		builder.WriteString(*textVariant.Text)
		builder.WriteString(" ")
	}
	return builder.String()
}

func NewOcrService(conf *conf.OcrConfig) (OcrSerivce, error) {
	ocrServiceUrl := conf.Url
	log.Printf("Creating ocr service url=%s\n", ocrServiceUrl)

	client, err := client.NewClientWithResponses(ocrServiceUrl)
	if err != nil {
		return nil, err
	}
	return &OcrServiceImpl{
			ocrclient: client,
		},
		nil
}
