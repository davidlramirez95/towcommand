package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	evidenceuc "github.com/davidlramirez95/towcommand/internal/usecase/evidence"
)

// CreateConditionReportRequest is the expected JSON body for POST /bookings/{id}/condition-report.
type CreateConditionReportRequest struct {
	Phase string `json:"phase" validate:"required,oneof=pickup dropoff"`
	Notes string `json:"notes"`
}

// CreateConditionReportHandler handles POST /bookings/{id}/condition-report requests.
type CreateConditionReportHandler struct {
	uc *evidenceuc.CreateConditionReportUseCase
}

// NewCreateConditionReportHandler constructs a CreateConditionReportHandler.
func NewCreateConditionReportHandler(uc *evidenceuc.CreateConditionReportUseCase) *CreateConditionReportHandler {
	return &CreateConditionReportHandler{uc: uc}
}

// Handle processes a create-condition-report API Gateway event.
func (h *CreateConditionReportHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	body, err := ParseBody[CreateConditionReportRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &evidenceuc.CreateConditionReportInput{
		BookingID:  bookingID,
		ProviderID: userID,
		Phase:      body.Phase,
		Notes:      body.Notes,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
