package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	ratinguc "github.com/davidlramirez95/towcommand/internal/usecase/rating"
)

// SubmitRatingRequest is the expected JSON body for POST /bookings/{id}/rating.
type SubmitRatingRequest struct {
	Score   int      `json:"score" validate:"required,min=1,max=5"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
}

// SubmitRatingHandler handles POST /bookings/{id}/rating requests.
type SubmitRatingHandler struct {
	uc *ratinguc.SubmitRatingUseCase
}

// NewSubmitRatingHandler constructs a SubmitRatingHandler.
func NewSubmitRatingHandler(uc *ratinguc.SubmitRatingUseCase) *SubmitRatingHandler {
	return &SubmitRatingHandler{uc: uc}
}

// Handle processes a submit-rating API Gateway event.
func (h *SubmitRatingHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	body, err := ParseBody[SubmitRatingRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &ratinguc.SubmitRatingInput{
		BookingID:  bookingID,
		CustomerID: userID,
		Score:      body.Score,
		Comment:    body.Comment,
		Tags:       body.Tags,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
