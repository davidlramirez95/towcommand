// Package main is the composition root for the nearby providers query Lambda.
package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	provider "github.com/davidlramirez95/towcommand/internal/usecase/provider"
)

func main() {
	cfg := config.Load()
	_ = logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)

	dynamoClient := awsclient.DynamoDBClient(cfg)
	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	repo := repository.NewProviderRepository(dynamoClient, cfg.DynamoDBTable)
	geo := cache.NewRedisGeoCache(redisClient)

	uc := provider.NewGetNearbyUseCase(geo, repo)

	h := handler.WithRecover(handler.WithCorrelationID(handler.WithLogging(
		func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			lat, _ := strconv.ParseFloat(handler.ParseQueryParam(event, "lat"), 64)
			lng, _ := strconv.ParseFloat(handler.ParseQueryParam(event, "lng"), 64)
			radiusKm, _ := strconv.ParseFloat(handler.ParseQueryParam(event, "radius"), 64)
			limit, _ := strconv.Atoi(handler.ParseQueryParam(event, "limit"))

			input := provider.GetNearbyInput{
				Lat:      lat,
				Lng:      lng,
				RadiusKm: radiusKm,
				Limit:    limit,
			}

			result, err := uc.Execute(ctx, input)
			if err != nil {
				return handler.ErrorResponse(err), nil
			}

			return handler.SuccessResponse(http.StatusOK, result), nil
		},
	)))

	lambda.Start(h)
}
