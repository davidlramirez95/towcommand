package awsclient

import (
	"testing"

	appconfig "github.com/davidlramirez95/towcommand/internal/platform/config"
)

func newTestConfig() *appconfig.Config {
	return &appconfig.Config{
		Region:              "ap-southeast-1",
		Stage:               "dev",
		DynamoDBEndpoint:    "http://localhost:4566",
		EventBridgeEndpoint: "http://localhost:4566",
		S3Endpoint:          "http://localhost:4566",
		SNSEndpoint:         "http://localhost:4566",
		SESEndpoint:         "http://localhost:4566",
		BedrockEndpoint:     "http://localhost:4566",
		RekognitionEndpoint: "http://localhost:4566",
		CognitoEndpoint:     "http://localhost:4566",
	}
}

func TestDynamoDBClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := DynamoDBClient(cfg)
	c2 := DynamoDBClient(cfg)

	if c1 == nil {
		t.Fatal("DynamoDBClient returned nil")
	}
	if c1 != c2 {
		t.Error("DynamoDBClient did not return the same instance")
	}
}

func TestEventBridgeClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := EventBridgeClient(cfg)
	c2 := EventBridgeClient(cfg)

	if c1 == nil {
		t.Fatal("EventBridgeClient returned nil")
	}
	if c1 != c2 {
		t.Error("EventBridgeClient did not return the same instance")
	}
}

func TestS3Client_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := S3Client(cfg)
	c2 := S3Client(cfg)

	if c1 == nil {
		t.Fatal("S3Client returned nil")
	}
	if c1 != c2 {
		t.Error("S3Client did not return the same instance")
	}
}

func TestSNSClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := SNSClient(cfg)
	c2 := SNSClient(cfg)

	if c1 == nil {
		t.Fatal("SNSClient returned nil")
	}
	if c1 != c2 {
		t.Error("SNSClient did not return the same instance")
	}
}

func TestSESClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := SESClient(cfg)
	c2 := SESClient(cfg)

	if c1 == nil {
		t.Fatal("SESClient returned nil")
	}
	if c1 != c2 {
		t.Error("SESClient did not return the same instance")
	}
}

func TestBedrockRuntimeClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := BedrockRuntimeClient(cfg)
	c2 := BedrockRuntimeClient(cfg)

	if c1 == nil {
		t.Fatal("BedrockRuntimeClient returned nil")
	}
	if c1 != c2 {
		t.Error("BedrockRuntimeClient did not return the same instance")
	}
}

func TestRekognitionClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := RekognitionClient(cfg)
	c2 := RekognitionClient(cfg)

	if c1 == nil {
		t.Fatal("RekognitionClient returned nil")
	}
	if c1 != c2 {
		t.Error("RekognitionClient did not return the same instance")
	}
}

func TestCognitoClient_Singleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := CognitoClient(cfg)
	c2 := CognitoClient(cfg)

	if c1 == nil {
		t.Fatal("CognitoClient returned nil")
	}
	if c1 != c2 {
		t.Error("CognitoClient did not return the same instance")
	}
}

func TestAPIGatewayManagementClient_NotSingleton(t *testing.T) {
	Reset()
	cfg := newTestConfig()

	c1 := APIGatewayManagementClient(cfg, "https://abc123.execute-api.ap-southeast-1.amazonaws.com/prod")
	c2 := APIGatewayManagementClient(cfg, "https://xyz789.execute-api.ap-southeast-1.amazonaws.com/prod")

	if c1 == nil || c2 == nil {
		t.Fatal("APIGatewayManagementClient returned nil")
	}
	if c1 == c2 {
		t.Error("APIGatewayManagementClient should create new instances per call")
	}
}

func TestReset_ClearsClients(t *testing.T) {
	cfg := newTestConfig()

	c1 := DynamoDBClient(cfg)
	Reset()
	c2 := DynamoDBClient(cfg)

	if c1 == c2 {
		t.Error("Reset() did not clear DynamoDB client — got same instance")
	}
}

func TestClients_WithLocalStackEndpoint(t *testing.T) {
	Reset()
	cfg := &appconfig.Config{
		Region:           "ap-southeast-1",
		Stage:            "dev",
		DynamoDBEndpoint: "http://localhost:4566",
		S3Endpoint:       "http://localhost:4566",
	}

	dynamo := DynamoDBClient(cfg)
	s3c := S3Client(cfg)

	if dynamo == nil {
		t.Error("DynamoDBClient with LocalStack endpoint returned nil")
	}
	if s3c == nil {
		t.Error("S3Client with LocalStack endpoint returned nil")
	}
}

func TestClients_WithoutEndpointOverride(t *testing.T) {
	Reset()
	cfg := &appconfig.Config{
		Region: "ap-southeast-1",
		Stage:  "prod",
	}

	dynamo := DynamoDBClient(cfg)
	if dynamo == nil {
		t.Error("DynamoDBClient without endpoint override returned nil")
	}
}
