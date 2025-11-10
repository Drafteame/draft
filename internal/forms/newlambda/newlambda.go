package newlambda

import (
	"fmt"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.LambdaInput) error {
	input.NextImportTag = data.NextImportTag
	input.NextLambdaImportTag = data.NextLambdaImportTag

	if err := baseForm(input); err != nil {
		return err
	}

	if input.LambdaType == "plain" {
		return nil
	}

	switch input.LambdaType {
	case "sqs", "snssqs":
		return queueForm(input)
	case "http":
		return httpForm(input)
	case "cron":
		return cronForm(input)
	case "custom":
		return customForm(input)
	default:
		return fmt.Errorf("unknown lambda type: %s", input.LambdaType)
	}
}
