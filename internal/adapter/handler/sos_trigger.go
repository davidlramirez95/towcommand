package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
	safetyuc "github.com/davidlramirez95/towcommand/internal/usecase/safety"
)

// TriggerSOSRequest is the expected JSON body for POST /bookings/{id}/sos.
type TriggerSOSRequest struct {
	TriggerType safety.TriggerType `json:"triggerType" validate:"required,oneof=TRIPLE_TAP SHAKE CODE_WORD BUTTON"`
	Lat         float64            `json:"lat" validate:"required"`
	Lng         float64            `json:"lng" validate:"required"`
}

// TriggerSOSHandler handles POST /bookings/{id}/sos requests.
type TriggerSOSHandler struct {
	uc *safetyuc.TriggerSOSUseCase
}

// NewTriggerSOSHandler constructs a TriggerSOSHandler.
func NewTriggerSOSHandler(uc *safetyuc.TriggerSOSUseCase) *TriggerSOSHandler {
	return &TriggerSOSHandler{uc: uc}
}

// Handle processes a trigger-SOS API Gateway event.
func (h *TriggerSOSHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	body, err := ParseBody[TriggerSOSRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &safetyuc.TriggerSOSInput{
		BookingID:   bookingID,
		TriggeredBy: userID,
		TriggerType: body.TriggerType,
		Lat:         body.Lat,
		Lng:         body.Lng,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
