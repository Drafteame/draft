package aws //nolint:typecheck

import (
	"github.com/Drafteame/draft/internal/pkg/exec"
)

func GetAccountID() (string, error) {
	cmd := "aws sts get-caller-identity --query Account --output text --profile draftea-dev"

	output, err := exec.Command(cmd)
	if err != nil {
		return "", err
	}

	return output, nil
}
