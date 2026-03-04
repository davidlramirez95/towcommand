package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

// CancelBookingRequest is the expected JSON body for POST /bookings/{id}/cancel.
type CancelBookingRequest struct {
	Reason string `json:"reason"`
}

// CancelBookingHandler handles POST /bookings/{id}/cancel requests.
type CancelBookingHandler struct {
	uc *bookinguc.CancelBookingUseCase
}

// NewCancelBookingHandler constructs a CancelBookingHandler.
func NewCancelBookingHandler(uc *bookinguc.CancelBookingUseCase) *CancelBookingHandler {
	return &CancelBookingHandler{uc: uc}
}

// Handle processes a cancel-booking API Gateway event.
func (h *CancelBookingHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("missing path parameter: id")), nil
	}

	var body CancelBookingRequest
	if event.Body != "" {
		parsed, err := ParseBody[CancelBookingRequest](event)
		if err != nil {
			return ErrorResponse(err), nil
		}
		body = parsed
	}

	result, err := h.uc.Execute(ctx, bookinguc.CancelBookingInput{
		BookingID: bookingID,
		CallerID:  userID,
		Reason:    body.Reason,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
