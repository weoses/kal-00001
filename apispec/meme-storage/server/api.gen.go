// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ImageDto defines model for ImageDto.
type ImageDto struct {
	ImageBase64 *string `json:"ImageBase64,omitempty"`
	MimeType    *string `json:"MimeType,omitempty"`
}

// ImageUrlDto defines model for ImageUrlDto.
type ImageUrlDto struct {
	Url *string `json:"Url,omitempty"`
}

// MemeDto defines model for MemeDto.
type MemeDto struct {
	Hash      *string             `json:"Hash,omitempty"`
	Id        *openapi_types.UUID `json:"Id,omitempty"`
	OcrResult *string             `json:"OcrResult,omitempty"`
}

// SearchMemeDto defines model for SearchMemeDto.
type SearchMemeDto struct {
	Hash               *string             `json:"Hash,omitempty"`
	Id                 *openapi_types.UUID `json:"Id,omitempty"`
	ImageThumbUrl      *string             `json:"ImageThumbUrl,omitempty"`
	ImageUrl           *string             `json:"ImageUrl,omitempty"`
	OcrResult          *string             `json:"OcrResult,omitempty"`
	OcrResultHighlight *[]string           `json:"OcrResultHighlight,omitempty"`
}

// AccountId defines model for AccountId.
type AccountId = openapi_types.UUID

// MemeId defines model for MemeId.
type MemeId = openapi_types.UUID

// MemeQuery defines model for MemeQuery.
type MemeQuery = string

// PageSize defines model for PageSize.
type PageSize = int

// SearchAfterId defines model for SearchAfterId.
type SearchAfterId = openapi_types.UUID

// SearchMemeParams defines parameters for SearchMeme.
type SearchMemeParams struct {
	MemeQuery     MemeQuery      `form:"MemeQuery" json:"MemeQuery"`
	SearchAfterId *SearchAfterId `form:"SearchAfterId,omitempty" json:"SearchAfterId,omitempty"`
	PageSize      *PageSize      `form:"PageSize,omitempty" json:"PageSize,omitempty"`
}

// CreateMemeJSONRequestBody defines body for CreateMeme for application/json ContentType.
type CreateMemeJSONRequestBody = ImageDto

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/accounts/{AccountId}/meme)
	SearchMeme(ctx echo.Context, accountId AccountId, params SearchMemeParams) error

	// (POST /api/v1/accounts/{AccountId}/meme)
	CreateMeme(ctx echo.Context, accountId AccountId) error

	// (GET /api/v1/accounts/{AccountId}/meme/{MemeId}/image/thumb/url)
	GetMemeImageThumbUrl(ctx echo.Context, accountId AccountId, memeId MemeId) error

	// (GET /api/v1/accounts/{AccountId}/meme/{MemeId}/image/url)
	GetMemeImageUrl(ctx echo.Context, accountId AccountId, memeId MemeId) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SearchMeme converts echo context to params.
func (w *ServerInterfaceWrapper) SearchMeme(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "AccountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, ctx.Param("AccountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter AccountId: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchMemeParams
	// ------------- Required query parameter "MemeQuery" -------------

	err = runtime.BindQueryParameter("form", true, true, "MemeQuery", ctx.QueryParams(), &params.MemeQuery)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter MemeQuery: %s", err))
	}

	// ------------- Optional query parameter "SearchAfterId" -------------

	err = runtime.BindQueryParameter("form", true, false, "SearchAfterId", ctx.QueryParams(), &params.SearchAfterId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter SearchAfterId: %s", err))
	}

	// ------------- Optional query parameter "PageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "PageSize", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter PageSize: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SearchMeme(ctx, accountId, params)
	return err
}

// CreateMeme converts echo context to params.
func (w *ServerInterfaceWrapper) CreateMeme(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "AccountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, ctx.Param("AccountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter AccountId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateMeme(ctx, accountId)
	return err
}

// GetMemeImageThumbUrl converts echo context to params.
func (w *ServerInterfaceWrapper) GetMemeImageThumbUrl(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "AccountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, ctx.Param("AccountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter AccountId: %s", err))
	}

	// ------------- Path parameter "MemeId" -------------
	var memeId MemeId

	err = runtime.BindStyledParameterWithLocation("simple", false, "MemeId", runtime.ParamLocationPath, ctx.Param("MemeId"), &memeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter MemeId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMemeImageThumbUrl(ctx, accountId, memeId)
	return err
}

// GetMemeImageUrl converts echo context to params.
func (w *ServerInterfaceWrapper) GetMemeImageUrl(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "AccountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, ctx.Param("AccountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter AccountId: %s", err))
	}

	// ------------- Path parameter "MemeId" -------------
	var memeId MemeId

	err = runtime.BindStyledParameterWithLocation("simple", false, "MemeId", runtime.ParamLocationPath, ctx.Param("MemeId"), &memeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter MemeId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMemeImageUrl(ctx, accountId, memeId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v1/accounts/:AccountId/meme", wrapper.SearchMeme)
	router.POST(baseURL+"/api/v1/accounts/:AccountId/meme", wrapper.CreateMeme)
	router.GET(baseURL+"/api/v1/accounts/:AccountId/meme/:MemeId/image/thumb/url", wrapper.GetMemeImageThumbUrl)
	router.GET(baseURL+"/api/v1/accounts/:AccountId/meme/:MemeId/image/url", wrapper.GetMemeImageUrl)

}

type SearchMemeRequestObject struct {
	AccountId AccountId `json:"AccountId"`
	Params    SearchMemeParams
}

type SearchMemeResponseObject interface {
	VisitSearchMemeResponse(w http.ResponseWriter) error
}

type SearchMeme200JSONResponse []SearchMemeDto

func (response SearchMeme200JSONResponse) VisitSearchMemeResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type CreateMemeRequestObject struct {
	AccountId AccountId `json:"AccountId"`
	Body      *CreateMemeJSONRequestBody
}

type CreateMemeResponseObject interface {
	VisitCreateMemeResponse(w http.ResponseWriter) error
}

type CreateMeme200JSONResponse MemeDto

func (response CreateMeme200JSONResponse) VisitCreateMemeResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetMemeImageThumbUrlRequestObject struct {
	AccountId AccountId `json:"AccountId"`
	MemeId    MemeId    `json:"MemeId"`
}

type GetMemeImageThumbUrlResponseObject interface {
	VisitGetMemeImageThumbUrlResponse(w http.ResponseWriter) error
}

type GetMemeImageThumbUrl200JSONResponse ImageUrlDto

func (response GetMemeImageThumbUrl200JSONResponse) VisitGetMemeImageThumbUrlResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetMemeImageUrlRequestObject struct {
	AccountId AccountId `json:"AccountId"`
	MemeId    MemeId    `json:"MemeId"`
}

type GetMemeImageUrlResponseObject interface {
	VisitGetMemeImageUrlResponse(w http.ResponseWriter) error
}

type GetMemeImageUrl200JSONResponse ImageUrlDto

func (response GetMemeImageUrl200JSONResponse) VisitGetMemeImageUrlResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /api/v1/accounts/{AccountId}/meme)
	SearchMeme(ctx context.Context, request SearchMemeRequestObject) (SearchMemeResponseObject, error)

	// (POST /api/v1/accounts/{AccountId}/meme)
	CreateMeme(ctx context.Context, request CreateMemeRequestObject) (CreateMemeResponseObject, error)

	// (GET /api/v1/accounts/{AccountId}/meme/{MemeId}/image/thumb/url)
	GetMemeImageThumbUrl(ctx context.Context, request GetMemeImageThumbUrlRequestObject) (GetMemeImageThumbUrlResponseObject, error)

	// (GET /api/v1/accounts/{AccountId}/meme/{MemeId}/image/url)
	GetMemeImageUrl(ctx context.Context, request GetMemeImageUrlRequestObject) (GetMemeImageUrlResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// SearchMeme operation middleware
func (sh *strictHandler) SearchMeme(ctx echo.Context, accountId AccountId, params SearchMemeParams) error {
	var request SearchMemeRequestObject

	request.AccountId = accountId
	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.SearchMeme(ctx.Request().Context(), request.(SearchMemeRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SearchMeme")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(SearchMemeResponseObject); ok {
		return validResponse.VisitSearchMemeResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// CreateMeme operation middleware
func (sh *strictHandler) CreateMeme(ctx echo.Context, accountId AccountId) error {
	var request CreateMemeRequestObject

	request.AccountId = accountId

	var body CreateMemeJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.CreateMeme(ctx.Request().Context(), request.(CreateMemeRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateMeme")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(CreateMemeResponseObject); ok {
		return validResponse.VisitCreateMemeResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetMemeImageThumbUrl operation middleware
func (sh *strictHandler) GetMemeImageThumbUrl(ctx echo.Context, accountId AccountId, memeId MemeId) error {
	var request GetMemeImageThumbUrlRequestObject

	request.AccountId = accountId
	request.MemeId = memeId

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetMemeImageThumbUrl(ctx.Request().Context(), request.(GetMemeImageThumbUrlRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetMemeImageThumbUrl")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetMemeImageThumbUrlResponseObject); ok {
		return validResponse.VisitGetMemeImageThumbUrlResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetMemeImageUrl operation middleware
func (sh *strictHandler) GetMemeImageUrl(ctx echo.Context, accountId AccountId, memeId MemeId) error {
	var request GetMemeImageUrlRequestObject

	request.AccountId = accountId
	request.MemeId = memeId

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetMemeImageUrl(ctx.Request().Context(), request.(GetMemeImageUrlRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetMemeImageUrl")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetMemeImageUrlResponseObject); ok {
		return validResponse.VisitGetMemeImageUrlResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}
