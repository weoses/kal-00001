// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
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
	ImageUrl           *string             `json:"ImageUrl,omitempty"`
	OcrResult          *string             `json:"OcrResult,omitempty"`
	OcrResultHighlight *[]string           `json:"OcrResultHighlight,omitempty"`
	Thumbnail          *SearchMemeThumb    `json:"Thumbnail,omitempty"`
}

// SearchMemeThumb defines model for SearchMemeThumb.
type SearchMemeThumb struct {
	ThumbHeight *int    `json:"ThumbHeight,omitempty"`
	ThumbUrl    *string `json:"ThumbUrl,omitempty"`
	ThumbWidth  *int    `json:"ThumbWidth,omitempty"`
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

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SearchMeme request
	SearchMeme(ctx context.Context, accountId AccountId, params *SearchMemeParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateMemeWithBody request with any body
	CreateMemeWithBody(ctx context.Context, accountId AccountId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateMeme(ctx context.Context, accountId AccountId, body CreateMemeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetMemeImageThumbUrl request
	GetMemeImageThumbUrl(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetMemeImageUrl request
	GetMemeImageUrl(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SearchMeme(ctx context.Context, accountId AccountId, params *SearchMemeParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchMemeRequest(c.Server, accountId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMemeWithBody(ctx context.Context, accountId AccountId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMemeRequestWithBody(c.Server, accountId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMeme(ctx context.Context, accountId AccountId, body CreateMemeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMemeRequest(c.Server, accountId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMemeImageThumbUrl(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMemeImageThumbUrlRequest(c.Server, accountId, memeId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMemeImageUrl(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMemeImageUrlRequest(c.Server, accountId, memeId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSearchMemeRequest generates requests for SearchMeme
func NewSearchMemeRequest(server string, accountId AccountId, params *SearchMemeParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/accounts/%s/meme", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "MemeQuery", runtime.ParamLocationQuery, params.MemeQuery); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if params.SearchAfterId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "SearchAfterId", runtime.ParamLocationQuery, *params.SearchAfterId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageSize != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "PageSize", runtime.ParamLocationQuery, *params.PageSize); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateMemeRequest calls the generic CreateMeme builder with application/json body
func NewCreateMemeRequest(server string, accountId AccountId, body CreateMemeJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateMemeRequestWithBody(server, accountId, "application/json", bodyReader)
}

// NewCreateMemeRequestWithBody generates requests for CreateMeme with any type of body
func NewCreateMemeRequestWithBody(server string, accountId AccountId, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/accounts/%s/meme", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetMemeImageThumbUrlRequest generates requests for GetMemeImageThumbUrl
func NewGetMemeImageThumbUrlRequest(server string, accountId AccountId, memeId MemeId) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "MemeId", runtime.ParamLocationPath, memeId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/accounts/%s/meme/%s/image/thumb/url", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetMemeImageUrlRequest generates requests for GetMemeImageUrl
func NewGetMemeImageUrlRequest(server string, accountId AccountId, memeId MemeId) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "AccountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "MemeId", runtime.ParamLocationPath, memeId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/accounts/%s/meme/%s/image/url", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SearchMemeWithResponse request
	SearchMemeWithResponse(ctx context.Context, accountId AccountId, params *SearchMemeParams, reqEditors ...RequestEditorFn) (*SearchMemeResponse, error)

	// CreateMemeWithBodyWithResponse request with any body
	CreateMemeWithBodyWithResponse(ctx context.Context, accountId AccountId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMemeResponse, error)

	CreateMemeWithResponse(ctx context.Context, accountId AccountId, body CreateMemeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMemeResponse, error)

	// GetMemeImageThumbUrlWithResponse request
	GetMemeImageThumbUrlWithResponse(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*GetMemeImageThumbUrlResponse, error)

	// GetMemeImageUrlWithResponse request
	GetMemeImageUrlWithResponse(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*GetMemeImageUrlResponse, error)
}

type SearchMemeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]SearchMemeDto
}

// Status returns HTTPResponse.Status
func (r SearchMemeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchMemeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateMemeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *MemeDto
}

// Status returns HTTPResponse.Status
func (r CreateMemeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateMemeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetMemeImageThumbUrlResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ImageUrlDto
}

// Status returns HTTPResponse.Status
func (r GetMemeImageThumbUrlResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMemeImageThumbUrlResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetMemeImageUrlResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ImageUrlDto
}

// Status returns HTTPResponse.Status
func (r GetMemeImageUrlResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMemeImageUrlResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SearchMemeWithResponse request returning *SearchMemeResponse
func (c *ClientWithResponses) SearchMemeWithResponse(ctx context.Context, accountId AccountId, params *SearchMemeParams, reqEditors ...RequestEditorFn) (*SearchMemeResponse, error) {
	rsp, err := c.SearchMeme(ctx, accountId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchMemeResponse(rsp)
}

// CreateMemeWithBodyWithResponse request with arbitrary body returning *CreateMemeResponse
func (c *ClientWithResponses) CreateMemeWithBodyWithResponse(ctx context.Context, accountId AccountId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMemeResponse, error) {
	rsp, err := c.CreateMemeWithBody(ctx, accountId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMemeResponse(rsp)
}

func (c *ClientWithResponses) CreateMemeWithResponse(ctx context.Context, accountId AccountId, body CreateMemeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMemeResponse, error) {
	rsp, err := c.CreateMeme(ctx, accountId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMemeResponse(rsp)
}

// GetMemeImageThumbUrlWithResponse request returning *GetMemeImageThumbUrlResponse
func (c *ClientWithResponses) GetMemeImageThumbUrlWithResponse(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*GetMemeImageThumbUrlResponse, error) {
	rsp, err := c.GetMemeImageThumbUrl(ctx, accountId, memeId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMemeImageThumbUrlResponse(rsp)
}

// GetMemeImageUrlWithResponse request returning *GetMemeImageUrlResponse
func (c *ClientWithResponses) GetMemeImageUrlWithResponse(ctx context.Context, accountId AccountId, memeId MemeId, reqEditors ...RequestEditorFn) (*GetMemeImageUrlResponse, error) {
	rsp, err := c.GetMemeImageUrl(ctx, accountId, memeId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMemeImageUrlResponse(rsp)
}

// ParseSearchMemeResponse parses an HTTP response from a SearchMemeWithResponse call
func ParseSearchMemeResponse(rsp *http.Response) (*SearchMemeResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SearchMemeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []SearchMemeDto
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateMemeResponse parses an HTTP response from a CreateMemeWithResponse call
func ParseCreateMemeResponse(rsp *http.Response) (*CreateMemeResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateMemeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest MemeDto
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetMemeImageThumbUrlResponse parses an HTTP response from a GetMemeImageThumbUrlWithResponse call
func ParseGetMemeImageThumbUrlResponse(rsp *http.Response) (*GetMemeImageThumbUrlResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMemeImageThumbUrlResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ImageUrlDto
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetMemeImageUrlResponse parses an HTTP response from a GetMemeImageUrlWithResponse call
func ParseGetMemeImageUrlResponse(rsp *http.Response) (*GetMemeImageUrlResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMemeImageUrlResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ImageUrlDto
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
