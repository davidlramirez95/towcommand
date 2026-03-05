package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
)

// APIGatewayHandler is the function signature for API Gateway proxy handlers.
type APIGatewayHandler func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// WithLogging wraps a handler with request/response logging.
// It logs HTTP method, path, status code, and duration using slog.Default().
func WithLogging(next APIGatewayHandler) APIGatewayHandler {
	return func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		start := time.Now()
		log := logger.WithContext(ctx, slog.Default())

		log.InfoContext(ctx, "request started",
			slog.String("method", event.HTTPMethod),
			slog.String("path", event.Path),
		)

		resp, err := next(ctx, event)

		log.InfoContext(ctx, "request completed",
			slog.Int("status", resp.StatusCode),
			slog.Duration("duration", time.Since(start)),
		)

		return resp, err
	}
}

// WithCorrelationID injects a correlation ID into the context.
// It checks the X-Correlation-ID header first, then falls back to the
// Lambda request ID, and finally generates a random ID.
func WithCorrelationID(next APIGatewayHandler) APIGatewayHandler {
	return func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		correlationID := event.Headers["X-Correlation-ID"]
		if correlationID == "" {
			correlationID = event.Headers["x-correlation-id"]
		}
		if correlationID == "" {
			if lc, ok := lambdacontext.FromContext(ctx); ok {
				correlationID = lc.AwsRequestID
			}
		}
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}
		ctx = logger.SetCorrelationID(ctx, correlationID)
		return next(ctx, event)
	}
}

// WithRecover wraps a handler with panic recovery that returns a 500 response.
func WithRecover(next APIGatewayHandler) APIGatewayHandler {
	return func(ctx context.Context, event *events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Default().ErrorContext(ctx, "panic recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
				resp = ErrorResponse(domainerrors.NewInternalError("internal server error"))
				err = nil
			}
		}()
		return next(ctx, event)
	}
}

// RequireRole returns middleware that checks the user's role from Cognito JWT claims.
// It allows the request if the user's type matches any of the allowed roles.
// Returns 403 Forbidden if the role is not in the allowed list.
func RequireRole(allowedRoles ...string) func(APIGatewayHandler) APIGatewayHandler {
	return func(next APIGatewayHandler) APIGatewayHandler {
		return func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			userType := ExtractUserType(event)
			if userType == "" {
				return ErrorResponse(domainerrors.NewForbiddenError("missing user type")), nil
			}
			for _, role := range allowedRoles {
				if userType == role {
					return next(ctx, event)
				}
			}
			return ErrorResponse(domainerrors.NewForbiddenError("insufficient permissions")), nil
		}
	}
}

// generateCorrelationID produces a random hex string for use as a correlation ID.
func generateCorrelationID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
