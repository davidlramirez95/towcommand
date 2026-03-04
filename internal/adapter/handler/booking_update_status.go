package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

// UpdateBookingStatusRequest is the expected JSON body for PATCH /bookings/{id}/status.
type UpdateBookingStatusRequest struct {
	Status   booking.BookingStatus `json:"status" validate:"required"`
	Metadata map[string]any        `json:"metadata"`
}

// UpdateBookingStatusHandler handles PATCH /bookings/{id}/status requests.
type UpdateBookingStatusHandler struct {
	uc *bookinguc.UpdateBookingStatusUseCase
}

// NewUpdateBookingStatusHandler constructs an UpdateBookingStatusHandler.
func NewUpdateBookingStatusHandler(uc *bookinguc.UpdateBookingStatusUseCase) *UpdateBookingStatusHandler {
	return &UpdateBookingStatusHandler{uc: uc}
}

// Handle processes an update-booking-status API Gateway event.
func (h *UpdateBookingStatusHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("missing path parameter: id")), nil
	}

	body, err := ParseBody[UpdateBookingStatusRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	userType := ExtractUserType(event)

	result, err := h.uc.Execute(ctx, bookinguc.UpdateBookingStatusInput{
		BookingID:  bookingID,
		CallerID:   userID,
		CallerType: userType,
		NewStatus:  body.Status,
		Metadata:   body.Metadata,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
