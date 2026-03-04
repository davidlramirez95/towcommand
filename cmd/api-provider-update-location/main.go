// Package main is the composition root for the provider location update Lambda.
package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	provider "github.com/davidlramirez95/towcommand/internal/usecase/provider"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)

	dynamoClient := awsclient.DynamoDBClient(cfg)
	ebClient := awsclient.EventBridgeClient(cfg)
	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	repo := repository.NewProviderRepository(dynamoClient, cfg.DynamoDBTable)
	geo := cache.NewRedisGeoCache(redisClient)
	pub := gateway.NewEventBridgePublisher(ebClient, cfg.EventBusName, log)

	uc := provider.NewUpdateLocationUseCase(repo, geo, pub, log)

	h := handler.WithRecover(handler.WithCorrelationID(handler.WithLogging(
		func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			input, err := handler.ParseBody[provider.UpdateLocationInput](event)
			if err != nil {
				return handler.ErrorResponse(err), nil
			}

			input.ProviderID = handler.ParsePathParam(event, "id")

			result, err := uc.Execute(ctx, input)
			if err != nil {
				return handler.ErrorResponse(err), nil
			}

			return handler.SuccessResponse(http.StatusOK, result), nil
		},
	)))

	lambda.Start(h)
}
