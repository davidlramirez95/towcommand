package handler_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

// errorBody represents the standard JSON error envelope returned by ErrorResponse.
type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// apiEventWithAuth builds a minimal API Gateway event with Cognito auth claims.
func apiEventWithAuth(userID string) *events.APIGatewayProxyRequest {
	return &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub": userID,
				},
			},
		},
	}
}

// parseErrorBody unmarshals a JSON error response body into an errorBody.
func parseErrorBody(t *testing.T, body string) errorBody {
	t.Helper()
	var eb errorBody
	require.NoError(t, json.Unmarshal([]byte(body), &eb))
	return eb
}
