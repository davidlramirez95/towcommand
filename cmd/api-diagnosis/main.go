package main

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	br := awsclient.BedrockRuntimeClient(cfg)

	engine := gateway.NewBedrockDiagnosisEngine(br)
	uc := diagnosisuc.NewDiagnoseUseCase(engine)
	h := handler.NewDiagnoseHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
