package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

// DiagnoseRequest is the expected JSON body for POST /diagnosis.
type DiagnoseRequest struct {
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	PhotoURLs   []string `json:"photoUrls"`
	VehicleType string   `json:"vehicleType"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

// DiagnoseHandler handles POST /diagnosis requests for AI-powered vehicle
// issue diagnosis.
type DiagnoseHandler struct {
	uc *diagnosisuc.DiagnoseUseCase
}

// NewDiagnoseHandler constructs a DiagnoseHandler.
func NewDiagnoseHandler(uc *diagnosisuc.DiagnoseUseCase) *DiagnoseHandler {
	return &DiagnoseHandler{uc: uc}
}

// Handle processes a diagnosis API Gateway event.
func (h *DiagnoseHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	body, err := ParseBody[DiagnoseRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	input := &diagnosisuc.DiagnosisInput{
		Description: body.Description,
		PhotoURLs:   body.PhotoURLs,
		VehicleType: body.VehicleType,
	}

	if body.Lat != nil && body.Lng != nil {
		input.Location = &diagnosisuc.LatLng{
			Lat: *body.Lat,
			Lng: *body.Lng,
		}
	}

	result, err := h.uc.Execute(ctx, input)
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
