// Package main is the composition root for the Cognito PreSignUp trigger Lambda.
package main

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/auth"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	uc := auth.NewPreSignUpUseCase(cfg.Stage)

	lambda.Start(func(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Default().ErrorContext(ctx, "panic in pre-signup",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
			}
		}()
		return uc.Execute(ctx, &event)
	})
}
