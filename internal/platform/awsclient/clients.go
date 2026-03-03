// Package awsclient provides singleton AWS SDK v2 client factories.
// Clients are initialized once per cold start and reused across warm Lambda invocations.
// It follows 12-Factor IV: backing services as attached resources (endpoint swappable).
package awsclient

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	appconfig "github.com/davidlramirez95/towcommand/internal/platform/config"
)

var (
	awsCfg     aws.Config
	awsCfgOnce sync.Once

	dynamoClient  *dynamodb.Client
	dynamoOnce    sync.Once
	ebClient      *eventbridge.Client
	ebOnce        sync.Once
	s3Client      *s3.Client
	s3Once        sync.Once
	snsClient     *sns.Client
	snsOnce       sync.Once
	sesClient     *ses.Client
	sesOnce       sync.Once
	brClient      *bedrockruntime.Client
	brOnce        sync.Once
	rekClient     *rekognition.Client
	rekOnce       sync.Once
	cognitoClient *cognitoidentityprovider.Client
	cognitoOnce   sync.Once
)

// Reset clears all cached clients, forcing re-initialization on next access.
// This is intended for testing only.
func Reset() {
	awsCfgOnce = sync.Once{}
	dynamoClient = nil
	dynamoOnce = sync.Once{}
	ebClient = nil
	ebOnce = sync.Once{}
	s3Client = nil
	s3Once = sync.Once{}
	snsClient = nil
	snsOnce = sync.Once{}
	sesClient = nil
	sesOnce = sync.Once{}
	brClient = nil
	brOnce = sync.Once{}
	rekClient = nil
	rekOnce = sync.Once{}
	cognitoClient = nil
	cognitoOnce = sync.Once{}
}

func loadAWSConfig(region string) aws.Config {
	awsCfgOnce.Do(func() {
		var err error
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
		)
		if err != nil {
			panic("awsclient: failed to load AWS config: " + err.Error())
		}
	})
	return awsCfg
}

// DynamoDBClient returns a singleton DynamoDB client.
func DynamoDBClient(cfg *appconfig.Config) *dynamodb.Client {
	dynamoOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*dynamodb.Options){}
		if cfg.DynamoDBEndpoint != "" {
			opts = append(opts, func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
			})
		}
		dynamoClient = dynamodb.NewFromConfig(awsCfg, opts...)
	})
	return dynamoClient
}

// EventBridgeClient returns a singleton EventBridge client.
func EventBridgeClient(cfg *appconfig.Config) *eventbridge.Client {
	ebOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*eventbridge.Options){}
		if cfg.EventBridgeEndpoint != "" {
			opts = append(opts, func(o *eventbridge.Options) {
				o.BaseEndpoint = aws.String(cfg.EventBridgeEndpoint)
			})
		}
		ebClient = eventbridge.NewFromConfig(awsCfg, opts...)
	})
	return ebClient
}

// S3Client returns a singleton S3 client.
func S3Client(cfg *appconfig.Config) *s3.Client {
	s3Once.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*s3.Options){}
		if cfg.S3Endpoint != "" {
			opts = append(opts, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(cfg.S3Endpoint)
				o.UsePathStyle = true
			})
		}
		s3Client = s3.NewFromConfig(awsCfg, opts...)
	})
	return s3Client
}

// SNSClient returns a singleton SNS client.
func SNSClient(cfg *appconfig.Config) *sns.Client {
	snsOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*sns.Options){}
		if cfg.SNSEndpoint != "" {
			opts = append(opts, func(o *sns.Options) {
				o.BaseEndpoint = aws.String(cfg.SNSEndpoint)
			})
		}
		snsClient = sns.NewFromConfig(awsCfg, opts...)
	})
	return snsClient
}

// SESClient returns a singleton SES client.
func SESClient(cfg *appconfig.Config) *ses.Client {
	sesOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*ses.Options){}
		if cfg.SESEndpoint != "" {
			opts = append(opts, func(o *ses.Options) {
				o.BaseEndpoint = aws.String(cfg.SESEndpoint)
			})
		}
		sesClient = ses.NewFromConfig(awsCfg, opts...)
	})
	return sesClient
}

// BedrockRuntimeClient returns a singleton Bedrock Runtime client.
func BedrockRuntimeClient(cfg *appconfig.Config) *bedrockruntime.Client {
	brOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*bedrockruntime.Options){}
		if cfg.BedrockEndpoint != "" {
			opts = append(opts, func(o *bedrockruntime.Options) {
				o.BaseEndpoint = aws.String(cfg.BedrockEndpoint)
			})
		}
		brClient = bedrockruntime.NewFromConfig(awsCfg, opts...)
	})
	return brClient
}

// RekognitionClient returns a singleton Rekognition client.
func RekognitionClient(cfg *appconfig.Config) *rekognition.Client {
	rekOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*rekognition.Options){}
		if cfg.RekognitionEndpoint != "" {
			opts = append(opts, func(o *rekognition.Options) {
				o.BaseEndpoint = aws.String(cfg.RekognitionEndpoint)
			})
		}
		rekClient = rekognition.NewFromConfig(awsCfg, opts...)
	})
	return rekClient
}

// CognitoClient returns a singleton Cognito Identity Provider client.
func CognitoClient(cfg *appconfig.Config) *cognitoidentityprovider.Client {
	cognitoOnce.Do(func() {
		awsCfg := loadAWSConfig(cfg.Region)
		opts := []func(*cognitoidentityprovider.Options){}
		if cfg.CognitoEndpoint != "" {
			opts = append(opts, func(o *cognitoidentityprovider.Options) {
				o.BaseEndpoint = aws.String(cfg.CognitoEndpoint)
			})
		}
		cognitoClient = cognitoidentityprovider.NewFromConfig(awsCfg, opts...)
	})
	return cognitoClient
}

// APIGatewayManagementClient creates a new API Gateway Management API client.
// Unlike other clients, this is NOT a singleton because the endpoint varies per
// API Gateway stage and must be provided at call time.
func APIGatewayManagementClient(cfg *appconfig.Config, endpoint string) *apigatewaymanagementapi.Client {
	ep := endpoint
	if ep == "" {
		ep = cfg.APIGatewayManagementEndpoint
	}
	awsCfg := loadAWSConfig(cfg.Region)
	opts := []func(*apigatewaymanagementapi.Options){}
	if ep != "" {
		opts = append(opts, func(o *apigatewaymanagementapi.Options) {
			o.BaseEndpoint = aws.String(ep)
		})
	}
	return apigatewaymanagementapi.NewFromConfig(awsCfg, opts...)
}
