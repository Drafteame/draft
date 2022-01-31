package tpl

var (
	SqsMainGo = `package main

import (
	"context"
	"errors"

	"github.com/Drafteame/framework/engine"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	handler "{{Namespace}}/internal/handlers/{{PackageName}}"
)

var app engine.App

func init() {
	app = engine.NewApp(func(cfg engine.AppConfig, s *engine.Services) {
		s.Router = false
		s.Mongo = true
		s.QLDB = true
		s.Sentry = true
	})
}

func main() {
	lambda.Start(worker)
}

func worker(ctx context.Context, event events.SQSEvent) (map[string]interface{}, error) {
	if len(event.Records) == 0 {
		return nil, errors.New("No SQS message passed to function")
	}

	return handler.Handle(ctx, event.Records)
}
`

	PlainMainGo = `package main

import (
	"context"

	"github.com/Drafteame/framework/engine"

	"github.com/aws/aws-lambda-go/lambda"

	handler "{{Namespace}}/internal/handlers/{{PackageName}}"
)

var app engine.App

func init() {
	app = engine.NewApp(func(cfg engine.AppConfig, s *engine.Services) {
		s.Router = false
		s.Mongo = true
		s.QLDB = true
		s.Sentry = true
	})
}

func main() {
	lambda.Start(worker)
}

func worker(ctx context.Context, event map[string]interface{}) error {
	return handler.Handle(ctx, event)
}
`

	HTTPMainGo = `package main

import (
	"github.com/Drafteame/framework/engine"
	"{{Namespace}}/internal/routes"
)

func main() {
	app := engine.NewApp()

	app.Start(func(a engine.App) {
		routes.Register(a)
	})
}
`
)
