package newlambda

import (
	"fmt"
	"regexp"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func cronForm(input *dtos.ServiceInput) error {
	return inputs.Text("Cron Expression:",
		inputs.WithDescription[string]("Enter the cron expression:"),
		inputs.WithValue(&input.CronExpression),
		inputs.WithValidation(func(val string) error {
			pattern := `^(rate\(\d+\s+(minute|minutes|hour|hours|day|days)\)|cron\((\*|\?|[\d\/,\-LW#]+)\s+(\*|\?|[\d\/,\-LW#]+)\s+(\*|\?|[\d\/,\-LW#]+)\s+(\*|\?|[\d\/,\-LW#]+)\s+(\*|\?|[\d\/,\-LW#]+)\s+(\*|\?|[\d\/,\-LW#]+)\))$`
			regex := regexp.MustCompile(pattern)

			if !regex.MatchString(val) {
				return fmt.Errorf("invalid cron expression: %s\nshould match: %s", val, pattern)
			}

			return nil
		}),
	)
}
