package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

// CreateBookingRequest is the expected JSON body for POST /bookings.
type CreateBookingRequest struct {
	VehicleID       string              `json:"vehicleId" validate:"required"`
	ServiceType     booking.ServiceType `json:"serviceType" validate:"required,oneof=FLATBED_TOW WHEEL_LIFT JUMPSTART TIRE_CHANGE FUEL_DELIVERY LOCKOUT ACCIDENT_RECOVERY"`
	PickupLocation  booking.GeoLocation `json:"pickupLocation" validate:"required"`
	DropoffLocation booking.GeoLocation `json:"dropoffLocation" validate:"required"`
	EstimateID      string              `json:"estimateId" validate:"required"`
	Notes           string              `json:"notes"`
}

// CreateBookingHandler handles POST /bookings requests.
type CreateBookingHandler struct {
	uc *bookinguc.CreateBookingUseCase
}

// NewCreateBookingHandler constructs a CreateBookingHandler.
func NewCreateBookingHandler(uc *bookinguc.CreateBookingUseCase) *CreateBookingHandler {
	return &CreateBookingHandler{uc: uc}
}

// Handle processes a create-booking API Gateway event.
func (h *CreateBookingHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	body, err := ParseBody[CreateBookingRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &bookinguc.CreateBookingInput{
		CustomerID:      userID,
		VehicleID:       body.VehicleID,
		ServiceType:     body.ServiceType,
		PickupLocation:  body.PickupLocation,
		DropoffLocation: body.DropoffLocation,
		EstimateID:      body.EstimateID,
		Notes:           body.Notes,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
