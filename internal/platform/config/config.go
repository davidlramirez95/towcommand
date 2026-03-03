// Package config loads application configuration from environment variables.
// It follows 12-Factor III: config strictly from environment, never hardcoded.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Required
	Region          string
	Stage           string
	DynamoDBTable   string
	EventBusName    string
	S3Bucket        string
	CognitoPoolID   string
	FunctionName    string
	FunctionVersion string

	// Optional (have defaults or can be empty)
	DynamoDBEndpoint             string
	EventBridgeEndpoint          string
	S3Endpoint                   string
	SNSEndpoint                  string
	SESEndpoint                  string
	BedrockEndpoint              string
	RekognitionEndpoint          string
	CognitoEndpoint              string
	APIGatewayManagementEndpoint string
	RedisHost                    string
	RedisPort                    int
	LogLevel                     string
}

// IsLocal returns true if running against LocalStack (dev/local stage with endpoints).
func (c *Config) IsLocal() bool {
	return (c.Stage == "local" || c.Stage == "dev") && c.DynamoDBEndpoint != ""
}

// IsProduction returns true if running in production.
func (c *Config) IsProduction() bool {
	return c.Stage == "prod"
}

// Load reads configuration from environment variables.
// It panics if any required variable is missing or empty.
func Load() *Config {
	cfg := &Config{
		// Required
		Region:          requireEnv("AWS_REGION"),
		Stage:           requireEnv("STAGE"),
		DynamoDBTable:   requireEnv("DYNAMODB_TABLE"),
		EventBusName:    requireEnv("EVENT_BUS_NAME"),
		S3Bucket:        requireEnv("S3_BUCKET"),
		CognitoPoolID:   requireEnv("COGNITO_USER_POOL_ID"),
		FunctionName:    getEnv("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		FunctionVersion: getEnv("AWS_LAMBDA_FUNCTION_VERSION", "$LATEST"),

		// Endpoint overrides for LocalStack / testing
		DynamoDBEndpoint:             os.Getenv("DYNAMODB_ENDPOINT"),
		EventBridgeEndpoint:          os.Getenv("EVENTBRIDGE_ENDPOINT"),
		S3Endpoint:                   os.Getenv("S3_ENDPOINT"),
		SNSEndpoint:                  os.Getenv("SNS_ENDPOINT"),
		SESEndpoint:                  os.Getenv("SES_ENDPOINT"),
		BedrockEndpoint:              os.Getenv("BEDROCK_ENDPOINT"),
		RekognitionEndpoint:          os.Getenv("REKOGNITION_ENDPOINT"),
		CognitoEndpoint:              os.Getenv("COGNITO_ENDPOINT"),
		APIGatewayManagementEndpoint: os.Getenv("APIGATEWAY_MANAGEMENT_ENDPOINT"),

		// Optional with defaults
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnvInt("REDIS_PORT", 6379),
		LogLevel:  strings.ToUpper(getEnv("LOG_LEVEL", "INFO")),
	}
	return cfg
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: required environment variable %s is not set", key))
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("config: environment variable %s must be an integer, got %q", key, v))
	}
	return n
}
