package aws //nolint:typecheck

import (
	"fmt"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

func GetAccountID(profile string) (string, error) {
	cmd := fmt.Sprintf("aws sts get-caller-identity --query Account --output text --profile %s", profile)

	output, err := exec.Command(cmd)
	if err != nil {
		return "", err
	}

	return output, nil
}
