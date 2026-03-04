package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

// ListBookingsHandler handles GET /bookings requests.
type ListBookingsHandler struct {
	uc *bookinguc.ListBookingsUseCase
}

// NewListBookingsHandler constructs a ListBookingsHandler.
func NewListBookingsHandler(uc *bookinguc.ListBookingsUseCase) *ListBookingsHandler {
	return &ListBookingsHandler{uc: uc}
}

// Handle processes a list-bookings API Gateway event.
func (h *ListBookingsHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	userType := ExtractUserType(event)

	var limit int32 = 25
	if ls := ParseQueryParam(event, "limit"); ls != "" {
		if n, err := strconv.ParseInt(ls, 10, 32); err == nil {
			limit = int32(n)
		}
	}

	statusFilter := ParseQueryParam(event, "status")

	result, err := h.uc.Execute(ctx, bookinguc.ListBookingsInput{
		CallerID:     userID,
		CallerType:   userType,
		Limit:        limit,
		StatusFilter: statusFilter,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
