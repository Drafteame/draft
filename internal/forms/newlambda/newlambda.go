package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.ServiceInput) error {
	if err := baseForm(input); err != nil {
		return err
	}

	if input.LambdaType == "sqs" || input.LambdaType == "snssqs" {
		if err := queueForm(input); err != nil {
			return err
		}
	}

	if input.LambdaType == "http" {
		if err := httpForm(input); err != nil {
			return err
		}
	}

	if input.LambdaType == "cron" {
		if err := cronForm(input); err != nil {
			return err
		}
	}

	return nil
}
