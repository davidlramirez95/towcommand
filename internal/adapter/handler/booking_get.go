package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

// GetBookingHandler handles GET /bookings/{id} requests.
type GetBookingHandler struct {
	uc *bookinguc.GetBookingUseCase
}

// NewGetBookingHandler constructs a GetBookingHandler.
func NewGetBookingHandler(uc *bookinguc.GetBookingUseCase) *GetBookingHandler {
	return &GetBookingHandler{uc: uc}
}

// Handle processes a get-booking API Gateway event.
func (h *GetBookingHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("missing path parameter: id")), nil
	}

	userType := ExtractUserType(event)

	result, err := h.uc.Execute(ctx, bookinguc.GetBookingInput{
		BookingID:  bookingID,
		CallerID:   userID,
		CallerType: userType,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
