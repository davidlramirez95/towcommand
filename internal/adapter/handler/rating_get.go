package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	ratinguc "github.com/davidlramirez95/towcommand/internal/usecase/rating"
)

// GetRatingHandler handles GET /bookings/{id}/rating requests.
type GetRatingHandler struct {
	uc *ratinguc.GetBookingRatingUseCase
}

// NewGetRatingHandler constructs a GetRatingHandler.
func NewGetRatingHandler(uc *ratinguc.GetBookingRatingUseCase) *GetRatingHandler {
	return &GetRatingHandler{uc: uc}
}

// Handle processes a get-rating API Gateway event.
func (h *GetRatingHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	result, err := h.uc.Execute(ctx, bookingID)
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
