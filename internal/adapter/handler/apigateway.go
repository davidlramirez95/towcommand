// Package handler provides shared helpers for parsing Lambda event inputs
// and building responses across API Gateway REST, WebSocket, EventBridge,
// and Cognito trigger handler types.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// validate is the package-level validator instance reused across all parse helpers.
var validate = validator.New()

// corsHeaders are the standard CORS headers applied to every API Gateway response.
var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "Content-Type,Authorization,X-Correlation-ID",
	"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	"Content-Type":                 "application/json",
}

// errorResponseBody is the standard JSON envelope for error responses.
type errorResponseBody struct {
	Error errorDetail `json:"error"`
}

// errorDetail carries the error code and human-readable message.
type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ParseBody unmarshals the JSON request body from an API Gateway event
// into T and runs struct validation using go-playground/validator tags.
func ParseBody[T any](event *events.APIGatewayProxyRequest) (T, error) {
	var body T
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return body, domainerrors.NewValidationError("invalid JSON body").WithCause(err)
	}
	if err := validate.Struct(body); err != nil {
		return body, domainerrors.NewValidationError(err.Error()).WithCause(err)
	}
	return body, nil
}

// ParsePathParam returns a path parameter value by key.
func ParsePathParam(event *events.APIGatewayProxyRequest, key string) string {
	return event.PathParameters[key]
}

// ParseQueryParam returns a query string parameter value by key.
func ParseQueryParam(event *events.APIGatewayProxyRequest, key string) string {
	return event.QueryStringParameters[key]
}

// ExtractUserID extracts the Cognito sub claim from the API Gateway authorizer.
func ExtractUserID(event *events.APIGatewayProxyRequest) string {
	claims, ok := event.RequestContext.Authorizer["claims"]
	if !ok {
		return ""
	}
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return ""
	}
	sub, _ := claimsMap["sub"].(string)
	return sub
}

// SuccessResponse builds an API Gateway response with the given status code
// and JSON-marshalled body. CORS headers are applied automatically.
func SuccessResponse(statusCode int, body any) events.APIGatewayProxyResponse {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ErrorResponse(domainerrors.NewInternalError("failed to marshal response body").WithCause(err))
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    corsHeaders,
		Body:       string(jsonBody),
	}
}

// ErrorResponse maps an error to an API Gateway error response.
// AppErrors are mapped to their HTTP status code; unknown errors become 500.
func ErrorResponse(err error) events.APIGatewayProxyResponse {
	var appErr *domainerrors.AppError
	if errors.As(err, &appErr) {
		body := errorResponseBody{
			Error: errorDetail{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			},
		}
		jsonBody, _ := json.Marshal(body)
		return events.APIGatewayProxyResponse{
			StatusCode: appErr.HTTPStatusCode(),
			Headers:    corsHeaders,
			Body:       string(jsonBody),
		}
	}

	body := errorResponseBody{
		Error: errorDetail{
			Code:    string(domainerrors.CodeInternalError),
			Message: "An unexpected error occurred",
		},
	}
	jsonBody, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    corsHeaders,
		Body:       string(jsonBody),
	}
}
