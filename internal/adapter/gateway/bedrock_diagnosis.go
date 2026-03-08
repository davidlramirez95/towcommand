package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"

	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

const bedrockModelID = "anthropic.claude-sonnet-4-5-20250929"

const diagnosisSystemPrompt = `You are an expert automotive diagnostic assistant for TowCommand PH, a roadside assistance platform in the Philippines. Based on the customer's description and photos, provide a diagnosis and service recommendation.

Respond ONLY with valid JSON in this exact format:
{
  "recommendedService": "SERVICE_TYPE",
  "urgencyLevel": "LEVEL",
  "estimatedCostMin": NUMBER_IN_CENTAVOS,
  "estimatedCostMax": NUMBER_IN_CENTAVOS,
  "description": "Brief explanation",
  "safetyWarnings": ["warning1", "warning2"]
}

Valid service types: FLATBED_TOWING, WHEEL_LIFT_TOWING, MOTORCYCLE_TOWING, JUMPSTART, TIRE_CHANGE, LOCKOUT, FUEL_DELIVERY, WINCH_RECOVERY
Valid urgency levels: LOW, MEDIUM, HIGH, CRITICAL
Cost estimates should be in Philippine Peso centavos (e.g., 250000 = PHP 2,500.00)`

// BedrockRuntimeAPI is the subset of the Bedrock Runtime client needed by
// the diagnosis engine.
type BedrockRuntimeAPI interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// bedrockRequest models the Anthropic Messages API request format used by
// Bedrock InvokeModel.
type bedrockRequest struct {
	AnthropicVersion string           `json:"anthropic_version"`
	MaxTokens        int              `json:"max_tokens"`
	System           string           `json:"system"`
	Messages         []bedrockMessage `json:"messages"`
}

// bedrockMessage is a single message in the Anthropic Messages API conversation.
type bedrockMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// bedrockResponse models the Anthropic Messages API response from Bedrock.
type bedrockResponse struct {
	Content []bedrockContentBlock `json:"content"`
}

// bedrockContentBlock represents a single content block in the response.
type bedrockContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// BedrockDiagnosisEngine implements DiagnosisEngine using AWS Bedrock with
// Claude Sonnet.
type BedrockDiagnosisEngine struct {
	client BedrockRuntimeAPI
}

// NewBedrockDiagnosisEngine creates a new Bedrock-powered diagnosis engine.
func NewBedrockDiagnosisEngine(client BedrockRuntimeAPI) *BedrockDiagnosisEngine {
	return &BedrockDiagnosisEngine{client: client}
}

// Diagnose sends the vehicle issue description to Claude Sonnet via Bedrock
// and returns a structured DiagnosisResult.
func (e *BedrockDiagnosisEngine) Diagnose(ctx context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
	userPrompt := buildUserPrompt(input)

	reqBody := bedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        1024,
		System:           diagnosisSystemPrompt,
		Messages: []bedrockMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling bedrock request: %w", err)
	}

	output, err := e.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bedrockModelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        reqJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("invoking bedrock model: %w", err)
	}

	var resp bedrockResponse
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshalling bedrock response: %w", err)
	}

	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("bedrock returned empty content")
	}

	// Extract the text content from the first block.
	text := resp.Content[0].Text

	var result diagnosisuc.DiagnosisResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parsing diagnosis JSON from AI response: %w", err)
	}

	return &result, nil
}

// buildUserPrompt constructs the user message from the diagnosis input.
func buildUserPrompt(input *diagnosisuc.DiagnosisInput) string {
	var sb strings.Builder

	sb.WriteString("Vehicle issue description: ")
	sb.WriteString(input.Description)

	if input.VehicleType != "" {
		sb.WriteString("\nVehicle type: ")
		sb.WriteString(input.VehicleType)
	}

	if input.Location != nil {
		sb.WriteString(fmt.Sprintf("\nLocation: lat=%.6f, lng=%.6f", input.Location.Lat, input.Location.Lng))
	}

	if len(input.PhotoURLs) > 0 {
		sb.WriteString(fmt.Sprintf("\nCustomer uploaded %d photo(s) of the issue.", len(input.PhotoURLs)))
	}

	return sb.String()
}
