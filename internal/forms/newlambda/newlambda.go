package newlambda

import (
	"fmt"

	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.ServiceInput) error {
	if err := baseForm(input); err != nil {
		return err
	}

	switch input.LambdaType {
	case "sqs", "snssqs":
		return queueForm(input)
	case "http":
		return httpForm(input)
	case "cron":
		return cronForm(input)
	default:
		return fmt.Errorf("unknown lambda type: %s", input.LambdaType)
	}
}
