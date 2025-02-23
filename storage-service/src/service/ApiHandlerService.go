package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"mine.local/ocr-gallery/apispec/meme-storage/server"
	"mine.local/ocr-gallery/storage-service/entity"
	"mine.local/ocr-gallery/storage-service/helper"
)

type ApiHandler struct {
	metaStorage  MetadataStorageService
	imageStorage ImageStorageService
	ocr          OcrSerivce
	validate     *validator.Validate
}

func (a *ApiHandler) iterateDocuments(ctx context.Context, accountId uuid.UUID, callback func(context.Context, *entity.ElasticMatchedContent) error) error {
	items, err := a.metaStorage.Search(ctx, accountId, "", nil, addr(100))
	for err == nil && len(items) > 0 {
		for _, item := range items {
			err = callback(ctx, item)
		}
		if len(items) > 0 {
			items, err = a.metaStorage.Search(ctx, accountId, "", &items[len(items)-1].Metadata.Created, addr(100))
		}
	}

	if err != nil {
		return err
	}
	return nil
}

// CheckDuplicates implements server.StrictServerInterface.
func (a *ApiHandler) CheckDuplicates(ctx context.Context, request server.CheckDuplicatesRequestObject) (server.CheckDuplicatesResponseObject, error) {
	deleted := map[uuid.UUID]any{}

	return server.CheckDuplicates200Response{},
		a.iterateDocuments(
			ctx,
			request.AccountId,
			func(ctx2 context.Context, emc *entity.ElasticMatchedContent) error {
				id := emc.Metadata.ImageId
				embedding := emc.Metadata.EmbeddingV1

				if _, ok := deleted[id]; ok {
					return nil
				}

				embeddingFoundImage, err := a.metaStorage.GetByEmbeddingV1(ctx2, embedding, 100)
				if err != nil {
					log.Printf("Failed to search image embedding duplicates : id=%s, err=%v", id, err)
					return nil
				}

				for i, item := range embeddingFoundImage {
					if i == 0 {
						continue
					}
					if _, ok := deleted[item.ImageId]; ok {
						continue
					}
					a.metaStorage.Delete(ctx2, item.ImageId)
					deleted[item.ImageId] = ""
				}
				return nil
			})

}

// UpdateOcr implements server.StrictServerInterface.
func (a *ApiHandler) UpdateOcr(ctx context.Context, request server.UpdateOcrRequestObject) (server.UpdateOcrResponseObject, error) {
	return server.UpdateOcr200Response{},
		a.iterateDocuments(ctx, request.AccountId, func(ctx2 context.Context, emc *entity.ElasticMatchedContent) error {

			id := emc.Metadata.ImageId
			accountId := emc.Metadata.AccountId
			hash := emc.Metadata.Hash
			s3id := emc.Metadata.S3Id
			created := emc.Metadata.Created

			log.Printf("UpdateOcr: checking image id=%s", id)

			img, err := a.imageStorage.GetImage(ctx2, s3id)
			if err != nil {
				log.Printf("Failed to read image from storage : id=%s, err=%v", id, err)
				return nil
			}

			ocrResult, err := a.ocr.DoOcr(ctx2, id, img)
			if err != nil {
				log.Printf("Failed to doOcr for image id=%s, err=%v", id, err)
				return nil
			}

			elasticObject := OcrResultToElastic(id, accountId, hash, created, ocrResult)
			err = a.metaStorage.Save(ctx2, elasticObject)
			if err != nil {
				log.Printf("Failed to save new metadata for image id=%s, err=%v", id, err)
				return nil
			}
			return nil
		})
}

// GetMemeImageThumbUrl implements server.StrictServerInterface.
func (a *ApiHandler) GetMemeImageThumbUrl(ctx context.Context, request server.GetMemeImageThumbUrlRequestObject) (server.GetMemeImageThumbUrlResponseObject, error) {
	memeMetadata, err := a.metaStorage.GetById(ctx, request.MemeId)
	if err != nil {
		return nil, err
	}

	if memeMetadata.AccountId != request.AccountId {
		return nil, echo.ErrNotFound
	}

	url, err := a.imageStorage.GetUrlThumb(ctx, memeMetadata.S3Id)
	if err != nil {
		return nil, err
	}

	resp := server.GetMemeImageThumbUrl200JSONResponse{}
	resp.Url = &url
	return resp, nil
}

// GetMemeImageUrl implements server.StrictServerInterface.
func (a *ApiHandler) GetMemeImageUrl(ctx context.Context, request server.GetMemeImageUrlRequestObject) (server.GetMemeImageUrlResponseObject, error) {
	memeMetadata, err := a.metaStorage.GetById(ctx, request.MemeId)
	if err != nil {
		return nil, err
	}

	if memeMetadata.AccountId != request.AccountId {
		return nil, echo.ErrNotFound
	}

	url, err := a.imageStorage.GetUrl(ctx, memeMetadata.S3Id)
	if err != nil {
		return nil, err
	}

	resp := server.GetMemeImageUrl200JSONResponse{}
	resp.Url = &url
	return resp, nil
}

func (a *ApiHandler) findHashDuplicates(
	ctx context.Context,
	hash string,
) (*entity.ElasticImageMetaData, error) {
	return a.metaStorage.GetByHash(ctx, hash)
}

func (a *ApiHandler) findContentDuplicates(
	ctx context.Context,
	ocrResult *OcrProcessedResult,
) (*entity.ElasticImageMetaData, error) {
	embedding := &entity.ElasticEmbeddingV1{
		Data:  &ocrResult.Embedding.Data,
		Model: ocrResult.Embedding.Model,
	}
	results, err := a.metaStorage.GetByEmbeddingV1(ctx, embedding, 1)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}
	return results[0], err
}

func (a *ApiHandler) HandleDuplicate(
	ctx context.Context,
	status server.DuplicateStatus,
	duplicate *entity.ElasticImageMetaData,
	request server.CreateMemeRequestObject,
) (server.CreateMemeResponseObject, error) {
	response := server.CreateMeme200JSONResponse{}
	if duplicate.AccountId != request.AccountId {
		copyMetadata := *duplicate
		log.Printf("Found meme duplicate in another account: id=%s", duplicate.ImageId)
		copyMetadata.ImageId, _ = uuid.NewRandom()
		copyMetadata.AccountId = request.AccountId
		err := a.metaStorage.Save(ctx, &copyMetadata)
		if err != nil {
			return nil, err
		}
		helper.ElasticToCreateResponse(&copyMetadata, status, &response)
	} else {
		log.Printf("Found meme duplicate in this account: id=%s", duplicate.ImageId)
		helper.ElasticToCreateResponse(duplicate, status, &response)
	}

	return response, nil
}

// CreateMeme implements server.StrictServerInterface.
func (a *ApiHandler) CreateMeme(
	ctx context.Context,
	request server.CreateMemeRequestObject,
) (server.CreateMemeResponseObject, error) {
	image := request.Body

	idUuid, _ := uuid.NewRandom()
	if len(*image.ImageBase64) == 0 {
		return nil, errors.New("image is empty")
	}

	hash := helper.CalcHash(request.Body.ImageBase64)
	hashDuplicate, err := a.findHashDuplicates(ctx, hash)
	if err != nil {
		return nil, err
	}

	if hashDuplicate != nil {
		return a.HandleDuplicate(ctx, server.DuplicateHash, hashDuplicate, request)
	}

	reqImage := helper.ImageToEntity2(request.Body)
	ocrResult, err := a.ocr.DoOcr(ctx, idUuid, reqImage)
	if err != nil {
		return nil, err
	}

	contentDuplicate, err := a.findContentDuplicates(ctx, ocrResult)
	if err != nil {
		return nil, err
	}

	if contentDuplicate != nil {
		return a.HandleDuplicate(ctx, server.DuplicateImage, contentDuplicate, request)
	}

	if strings.TrimSpace(ocrResult.OcrText) == "" {
		return nil, errors.New("no text on image")
	}

	err = a.imageStorage.Save(ctx, idUuid, ocrResult.Image, ocrResult.Thumbnail.Image)
	if err != nil {
		return nil, err
	}

	elasticMetaData := OcrResultToElastic(
		idUuid,
		request.AccountId,
		hash,
		time.Now().UnixMicro(),
		ocrResult,
	)

	err = a.validate.Struct(elasticMetaData)
	if err != nil {
		//TODO handle fail
		return nil, err
	}

	err = a.metaStorage.Save(ctx, elasticMetaData)
	if err != nil {
		//TODO handle fail
		return nil, err
	}

	response := server.CreateMeme200JSONResponse{}
	helper.ElasticToCreateResponse(elasticMetaData, server.New, &response)
	return response, nil
}

// SearchMeme implements server.StrictServerInterface.
func (a *ApiHandler) SearchMeme(ctx context.Context, request server.SearchMemeRequestObject) (server.SearchMemeResponseObject, error) {
	query := request.Params.MemeQuery

	log.Printf("SearchMeme: query=%s", query)

	matchedMetadata, err := a.metaStorage.Search(
		ctx,
		request.AccountId,
		query,
		request.Params.SearchAfterSortId,
		request.Params.PageSize,
	)

	if err != nil {
		return nil, err
	}

	response := make(server.SearchMeme200JSONResponse, len(matchedMetadata))
	for index, metadataItem := range matchedMetadata {

		imageThumbUrl, err := a.imageStorage.GetUrlThumb(ctx, metadataItem.Metadata.S3Id)
		if err != nil {
			return nil, err
		}
		imageUrl, err := a.imageStorage.GetUrl(ctx, metadataItem.Metadata.S3Id)
		if err != nil {
			return nil, err
		}

		dto := server.SearchMemeDto{}
		helper.ElasticToSearchMemeDto(metadataItem, &dto)
		dto.ImageUrl = &imageUrl
		dto.Thumbnail = new(server.SearchMemeThumb)
		dto.Thumbnail.ThumbUrl = &imageThumbUrl
		dto.Thumbnail.ThumbHeight = &metadataItem.Metadata.ThumbSize.Height
		dto.Thumbnail.ThumbWidth = &metadataItem.Metadata.ThumbSize.Width
		response[index] = dto
	}

	return response, nil
}

func OcrResultToElastic(idUuid uuid.UUID, accountId uuid.UUID, hash string, created int64, ocrResult *OcrProcessedResult) *entity.ElasticImageMetaData {
	return &entity.ElasticImageMetaData{
		ImageId:   idUuid,
		S3Id:      idUuid,
		AccountId: accountId,
		ThumbSize: &entity.ElasticSizes{
			Height: ocrResult.Thumbnail.Height,
			Width:  ocrResult.Thumbnail.Width,
		},
		Created: created,
		Hash:    hash,
		Result:  ocrResult.OcrText,
		EmbeddingV1: &entity.ElasticEmbeddingV1{
			Data:  &ocrResult.Embedding.Data,
			Model: ocrResult.Embedding.Model,
		},
	}
}

func NewApiHandler(
	metaStorage MetadataStorageService,
	imageStorage ImageStorageService,
	ocr OcrSerivce,
	validator *validator.Validate,
) server.StrictServerInterface {
	return &ApiHandler{
		metaStorage:  metaStorage,
		imageStorage: imageStorage,
		ocr:          ocr,
		validate:     validator,
	}
}
