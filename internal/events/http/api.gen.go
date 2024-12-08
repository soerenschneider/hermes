//go:build go1.22

// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// NotificationRequest defines model for NotificationRequest.
type NotificationRequest struct {
	// Message The message content of the notification.
	Message string `json:"message" validate:"required"`

	// RoutingKey The routing key used to categorize or route the notification.
	RoutingKey *string `json:"routing_key,omitempty"`

	// ServiceId The ID of the service requesting the notification.
	ServiceId string `json:"service_id" validate:"required"`

	// Subject The subject of the notification.
	Subject string `json:"subject" validate:"required"`
}

// SendNotificationJSONRequestBody defines body for SendNotification for application/json ContentType.
type SendNotificationJSONRequestBody = NotificationRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Send a notification
	// (POST /notify)
	SendNotification(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// SendNotification operation middleware
func (siw *ServerInterfaceWrapper) SendNotification(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendNotification(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/notify", wrapper.SendNotification)

	return m
}

type SendNotificationRequestObject struct {
	Body *SendNotificationJSONRequestBody
}

type SendNotificationResponseObject interface {
	VisitSendNotificationResponse(w http.ResponseWriter) error
}

type SendNotification200Response struct {
}

func (response SendNotification200Response) VisitSendNotificationResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type SendNotification400JSONResponse struct {
	Details *[]string `json:"details,omitempty"`
	Error   string    `json:"error,omitempty"`
}

func (response SendNotification400JSONResponse) VisitSendNotificationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type SendNotification500JSONResponse struct {
	Error *string `json:"error,omitempty"`
}

func (response SendNotification500JSONResponse) VisitSendNotificationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Send a notification
	// (POST /notify)
	SendNotification(ctx context.Context, request SendNotificationRequestObject) (SendNotificationResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// SendNotification operation middleware
func (sh *strictHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var request SendNotificationRequestObject

	var body SendNotificationJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.SendNotification(ctx, request.(SendNotificationRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SendNotification")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SendNotificationResponseObject); ok {
		if err := validResponse.VisitSendNotificationResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xVwY4bNwz9FYHtoQXGYydpL3PbtAXqS7FoewsWgSrRNrMzkkJy3HUX/vdCGtmx40ET",
	"oHuTOSL1Hvn4/AwuDikGDCrQPYO4HQ62HH+LShtyVimG3/HjiKI5bL2nHLL9PceErIQC3cb2gg2ki9Az",
	"DChit5iPHsUxpZwIHfy5Q1M/GheDYlATN0Z3aMLFqy00gE92SD1CB3fGMSk52xtkjmx2Vkx0bmRGbyiU",
	"dEHek8OcqYeU00SZwhYaeFpEm2jhoscthgU+KduF2m2Burc9eas5gfHjSIwejscGOI5KYfv+EQ/zNOoF",
	"84gHMwp6o9E4q7iNTP+giVxu4Be4FUKt7ZH1BvmxgcrqPfl5EOufT+2rNw1PA8vI/vvlmvDq9ZuXapmM",
	"f31Ap/NI68cvT/un06x/ya15GWx5nqdf3bvLrn5C3ZxV+3B+M06fSgEKm3hL7e5+bTaRjWDwueeXvMTs",
	"yRa2yxI+GAw+RQpaVEpa+P6KPKCYu/s1NLBHlqnwq3bVrrIEYsJgE0EHb9pVm4eVrO4K2Vo2H1OUmcb/",
	"gcGLsVeozN+ku4IqcdyTR288qqVesoKLlBI62hD6y6XK613y174WvrQJmBqMom+jL4DqdhfjSKmv95Yf",
	"JAM72U0+fcu4gQ6+WX7yo2U1o+WcE5Vh3OrriqK3ajObPJUWyvQlxSCTO71erW5bdflUzlMjo3Moshn7",
	"/pBr/DClfTWxa0usPc7Hs9YvlWhIzFmjZy1eRR8aIMWh1LhxihqwzLagLcZy9RqsQ9kOQyGN2piBRLJk",
	"T/XNhrD3Mu+g27jIwYU8UlrENP0PLIqYkaFTHvF4nN2b6za/tf5kUea7uqy54wWufJ+R//i/+jzLW5GD",
	"7Yuckae3bml+Ff7ZWtX9hsHyoS7HZ0s3lZpyBLp3n6vvnqMfXdVevgQNjNxDBzvVJN1yaRO1u2IVbWXW",
	"ujjA8eH4bwAAAP//ZpYjv8oHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
