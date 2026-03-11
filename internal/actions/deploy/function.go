package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const deployRegion = "us-east-2"

// DeployFunction packages and deploys a single Lambda function using the given env config.
// serviceArg can be a service name or a path.
func DeployFunction(env EnvConfig, serviceArg, functionName string) error {
	absPath, err := resolveService(serviceArg)
	if err != nil {
		return err
	}

	if fileExists(filepath.Join(absPath, ".deployignore")) {
		log.Warnf("Skipping %s: .deployignore found", absPath)
		return nil
	}

	slsFile := filepath.Join(absPath, "serverless.yml")
	if !fileExists(slsFile) {
		return fmt.Errorf("serverless.yml not found in %s", absPath)
	}

	log.Info("Fetching AWS Account ID...")
	accountID, err := aws.GetAccountID(env.AWSProfile)
	if err != nil {
		return fmt.Errorf("failed to get AWS account ID: %w", err)
	}

	if fileExists(filepath.Join(absPath, "package.json")) {
		log.Info("Installing dependencies...")
		installScript := fmt.Sprintf("cd %q && npm install", absPath)
		if _, err := exec.Command(installScript, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	log.Info("Packaging service...")
	packageScript := fmt.Sprintf(
		`cd %q && env STAGE=%s AWS_ACCOUNT=%s sls package --stage %s --verbose --aws-profile %s`,
		absPath, env.Stage, accountID, env.Stage, env.AWSProfile,
	)
	if _, err := exec.Command(packageScript, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
		return fmt.Errorf("sls package failed: %w", err)
	}

	serviceName, err := parseServiceName(slsFile)
	if err != nil {
		return err
	}

	fullLambdaName := fmt.Sprintf("%s-%s-%s", serviceName, env.Stage, functionName)
	log.Infof("Lambda function: %s", fullLambdaName)

	zipFile := filepath.Join(absPath, ".bin", functionName+".zip")
	if !fileExists(zipFile) {
		return fmt.Errorf(".bin/%s.zip not found — check that the function name matches serverless.yml", functionName)
	}

	log.Info("Updating Lambda code...")
	updateScript := fmt.Sprintf(
		`aws lambda update-function-code --function-name %q --zip-file "fileb://%s" --profile %s --region %s --output json`,
		fullLambdaName, zipFile, env.AWSProfile, deployRegion,
	)
	out, err := exec.Command(updateScript)
	if err != nil {
		return fmt.Errorf("lambda update failed: %w\n%s", err, out)
	}

	log.Success("✓ Lambda deployed:", fullLambdaName)
	return nil
}

func parseServiceName(slsFile string) (string, error) {
	content, err := files.Read(slsFile)
	if err != nil {
		return "", fmt.Errorf("failed to read serverless.yml: %w", err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "service:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("could not find 'service:' field in serverless.yml")
}
