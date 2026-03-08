// Package main is the composition root for the api-push-register Lambda handler.
// It handles POST /users/{id}/push-token to register push notification tokens
// for mobile devices via FCM (Android) or APNs (iOS).
package main

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	snsClient := awsclient.SNSClient(cfg)

	tokenRepo := repository.NewPushRepository(ddb, cfg.DynamoDBTable)
	pushSender := gateway.NewSNSPushSender(snsClient)

	h := handler.NewRegisterPushTokenHandler(tokenRepo, pushSender)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
