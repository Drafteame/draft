package tpl

var (
	PlainHandlerGo = `package {{PackageName}}

import (
	"context"
)

func Handle(ctx context.Context, event map[string]interface{}) error {
	return nil
}

`

	SqsHandlerGo = `package {{PackageName}}

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func Handle(ctx context.Context, records []events.SQSMessage) (map[string]interface{}, error) {
	return nil, nil
}

`
)
