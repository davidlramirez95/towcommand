package config

import (
	"os"
	"testing"
)

// setEnv is a test helper that sets environment variables and returns a cleanup function.
func setEnv(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

// requiredEnvs returns the minimum env vars needed for Load() to succeed.
func requiredEnvs() map[string]string {
	return map[string]string{
		"AWS_REGION":           "ap-southeast-1",
		"STAGE":                "dev",
		"DYNAMODB_TABLE":       "TowCommand-dev",
		"EVENT_BUS_NAME":       "towcommand-events",
		"S3_BUCKET":            "towcommand-uploads-dev",
		"COGNITO_USER_POOL_ID": "ap-southeast-1_abc123",
	}
}

func TestLoad_AllRequired(t *testing.T) {
	setEnv(t, requiredEnvs())

	cfg := Load()

	if cfg.Region != "ap-southeast-1" {
		t.Errorf("Region = %q, want %q", cfg.Region, "ap-southeast-1")
	}
	if cfg.Stage != "dev" {
		t.Errorf("Stage = %q, want %q", cfg.Stage, "dev")
	}
	if cfg.DynamoDBTable != "TowCommand-dev" {
		t.Errorf("DynamoDBTable = %q, want %q", cfg.DynamoDBTable, "TowCommand-dev")
	}
	if cfg.EventBusName != "towcommand-events" {
		t.Errorf("EventBusName = %q, want %q", cfg.EventBusName, "towcommand-events")
	}
	if cfg.S3Bucket != "towcommand-uploads-dev" {
		t.Errorf("S3Bucket = %q, want %q", cfg.S3Bucket, "towcommand-uploads-dev")
	}
	if cfg.CognitoPoolID != "ap-southeast-1_abc123" {
		t.Errorf("CognitoPoolID = %q, want %q", cfg.CognitoPoolID, "ap-southeast-1_abc123")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, requiredEnvs())

	cfg := Load()

	if cfg.RedisHost != "localhost" {
		t.Errorf("RedisHost = %q, want %q", cfg.RedisHost, "localhost")
	}
	if cfg.RedisPort != 6379 {
		t.Errorf("RedisPort = %d, want %d", cfg.RedisPort, 6379)
	}
	if cfg.LogLevel != "INFO" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "INFO")
	}
	if cfg.FunctionName != "unknown" {
		t.Errorf("FunctionName = %q, want %q", cfg.FunctionName, "unknown")
	}
	if cfg.FunctionVersion != "$LATEST" {
		t.Errorf("FunctionVersion = %q, want %q", cfg.FunctionVersion, "$LATEST")
	}
}

func TestLoad_OptionalOverrides(t *testing.T) {
	envs := requiredEnvs()
	envs["REDIS_HOST"] = "redis.cluster.local"
	envs["REDIS_PORT"] = "6380"
	envs["LOG_LEVEL"] = "debug"
	envs["DYNAMODB_ENDPOINT"] = "http://localhost:4566"
	envs["AWS_LAMBDA_FUNCTION_NAME"] = "create-booking"
	envs["AWS_LAMBDA_FUNCTION_VERSION"] = "42"
	setEnv(t, envs)

	cfg := Load()

	if cfg.RedisHost != "redis.cluster.local" {
		t.Errorf("RedisHost = %q, want %q", cfg.RedisHost, "redis.cluster.local")
	}
	if cfg.RedisPort != 6380 {
		t.Errorf("RedisPort = %d, want %d", cfg.RedisPort, 6380)
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "DEBUG")
	}
	if cfg.DynamoDBEndpoint != "http://localhost:4566" {
		t.Errorf("DynamoDBEndpoint = %q, want %q", cfg.DynamoDBEndpoint, "http://localhost:4566")
	}
	if cfg.FunctionName != "create-booking" {
		t.Errorf("FunctionName = %q, want %q", cfg.FunctionName, "create-booking")
	}
	if cfg.FunctionVersion != "42" {
		t.Errorf("FunctionVersion = %q, want %q", cfg.FunctionVersion, "42")
	}
}

func TestLoad_PanicsOnMissingRequired(t *testing.T) {
	required := []string{
		"AWS_REGION",
		"STAGE",
		"DYNAMODB_TABLE",
		"EVENT_BUS_NAME",
		"S3_BUCKET",
		"COGNITO_USER_POOL_ID",
	}

	for _, key := range required {
		t.Run(key, func(t *testing.T) {
			// Set all required, then unset the one we're testing
			envs := requiredEnvs()
			delete(envs, key)
			// Clear all env vars first to avoid leaking from parallel tests
			for k := range requiredEnvs() {
				_ = os.Unsetenv(k)
			}
			setEnv(t, envs)

			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("Load() did not panic when %s is missing", key)
				}
			}()
			Load()
		})
	}
}

func TestLoad_PanicsOnInvalidRedisPort(t *testing.T) {
	envs := requiredEnvs()
	envs["REDIS_PORT"] = "not-a-number"
	setEnv(t, envs)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Load() did not panic on invalid REDIS_PORT")
		}
	}()
	Load()
}

func TestConfig_IsLocal(t *testing.T) {
	tests := []struct {
		name     string
		stage    string
		endpoint string
		want     bool
	}{
		{"dev with endpoint", "dev", "http://localhost:4566", true},
		{"local with endpoint", "local", "http://localhost:4566", true},
		{"dev without endpoint", "dev", "", false},
		{"prod with endpoint", "prod", "http://localhost:4566", false},
		{"staging without endpoint", "staging", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Stage: tt.stage, DynamoDBEndpoint: tt.endpoint}
			if got := cfg.IsLocal(); got != tt.want {
				t.Errorf("IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name  string
		stage string
		want  bool
	}{
		{"prod", "prod", true},
		{"dev", "dev", false},
		{"staging", "staging", false},
		{"local", "local", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Stage: tt.stage}
			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad_EndpointOverrides(t *testing.T) {
	envs := requiredEnvs()
	envs["DYNAMODB_ENDPOINT"] = "http://localhost:4566"
	envs["EVENTBRIDGE_ENDPOINT"] = "http://localhost:4566"
	envs["S3_ENDPOINT"] = "http://localhost:4566"
	envs["SNS_ENDPOINT"] = "http://localhost:4566"
	envs["SES_ENDPOINT"] = "http://localhost:4566"
	envs["BEDROCK_ENDPOINT"] = "http://localhost:4566"
	envs["REKOGNITION_ENDPOINT"] = "http://localhost:4566"
	envs["COGNITO_ENDPOINT"] = "http://localhost:4566"
	envs["APIGATEWAY_MANAGEMENT_ENDPOINT"] = "http://localhost:4566"
	setEnv(t, envs)

	cfg := Load()

	endpoints := map[string]string{
		"DynamoDBEndpoint":             cfg.DynamoDBEndpoint,
		"EventBridgeEndpoint":          cfg.EventBridgeEndpoint,
		"S3Endpoint":                   cfg.S3Endpoint,
		"SNSEndpoint":                  cfg.SNSEndpoint,
		"SESEndpoint":                  cfg.SESEndpoint,
		"BedrockEndpoint":              cfg.BedrockEndpoint,
		"RekognitionEndpoint":          cfg.RekognitionEndpoint,
		"CognitoEndpoint":              cfg.CognitoEndpoint,
		"APIGatewayManagementEndpoint": cfg.APIGatewayManagementEndpoint,
	}
	for name, val := range endpoints {
		if val != "http://localhost:4566" {
			t.Errorf("%s = %q, want %q", name, val, "http://localhost:4566")
		}
	}
}
