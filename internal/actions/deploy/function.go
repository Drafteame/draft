package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const deployRegion = "us-east-2"

// DeployFunction packages and deploys one or more Lambda functions using the given env config.
// serviceArg can be a service name or a path. All functions must belong to the same service.
func DeployFunction(env EnvConfig, serviceArg string, functionNames []string) error {
	absPath, err := resolveService(serviceArg)
	if err != nil {
		return err
	}

	skip, err := validateServiceDir(absPath)
	if err != nil {
		return err
	}
	if skip {
		log.Warnf("Skipping %s: .deployignore found", absPath)
		return nil
	}

	stage := env.Stage()

	log.Info("Fetching AWS Account ID...")
	accountID, err := aws.GetAccountID(env.Profile)
	if err != nil {
		return fmt.Errorf("failed to get AWS account ID: %w", err)
	}

	if files.Exists(filepath.Join(absPath, "package.json")) {
		log.Info("Installing dependencies...")
		installScript := fmt.Sprintf("cd %q && npm install", absPath)
		if _, err := exec.Command(installScript, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	syncSecretsDry := "false"
	if env.SyncSecretsDry {
		syncSecretsDry = "true"
	}

	log.Info("Packaging service...")
	packageScript := fmt.Sprintf(
		`cd %q && env STAGE=%s AWS_ACCOUNT=%s SLS_SYNC_SECRETS_DRY=%s sls package --stage %s --verbose --aws-profile %s`,
		absPath, stage, accountID, syncSecretsDry, stage, env.Profile,
	)
	if _, err := exec.Command(packageScript, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
		return fmt.Errorf("sls package failed: %w", err)
	}

	slsFile := filepath.Join(absPath, "serverless.yml")
	serviceName, err := parseServiceName(slsFile)
	if err != nil {
		return err
	}

	type result struct {
		name string
		err  error
	}

	results := make([]result, len(functionNames))
	var wg sync.WaitGroup

	for i, functionName := range functionNames {
		wg.Add(1)
		go func(i int, functionName string) {
			defer wg.Done()
			results[i] = result{
				name: functionName,
				err:  deployOneFunction(env, absPath, serviceName, stage, functionName),
			}
		}(i, functionName)
	}

	wg.Wait()

	hasError := false
	if len(functionNames) > 1 {
		log.Info("\n─── Deploy Summary ───")
	}
	for _, r := range results {
		if r.err != nil {
			log.Errorf("✗ %s: %v", r.name, r.err)
			hasError = true
		} else {
			log.Successf("✓ %s", r.name)
		}
	}

	if hasError {
		return fmt.Errorf("one or more function deploys failed")
	}
	return nil
}

func deployOneFunction(env EnvConfig, absPath, serviceName, stage, functionName string) error {
	fullLambdaName := fmt.Sprintf("%s-%s-%s", serviceName, stage, functionName)
	log.Infof("Lambda function: %s", fullLambdaName)

	zipFile := filepath.Join(absPath, ".bin", functionName+".zip")
	if !files.Exists(zipFile) {
		return fmt.Errorf(".bin/%s.zip not found — check that the function name matches serverless.yml", functionName)
	}

	log.Info("Updating Lambda code...")
	updateScript := fmt.Sprintf(
		`aws lambda update-function-code --function-name %q --zip-file "fileb://%s" --profile %s --region %s --output json`,
		fullLambdaName, zipFile, env.Profile, deployRegion,
	)
	out, err := exec.Command(updateScript)
	if err != nil {
		return fmt.Errorf("lambda update failed: %w\n%s", err, out)
	}

	log.Info("Waiting for Lambda update to complete...")
	waitScript := fmt.Sprintf(
		`aws lambda wait function-updated --function-name %q --profile %s --region %s`,
		fullLambdaName, env.Profile, deployRegion,
	)
	if out, err := exec.Command(waitScript); err != nil {
		return fmt.Errorf("lambda wait failed: %w\n%s", err, out)
	}

	log.Info("Publishing new version...")
	publishScript := fmt.Sprintf(
		`aws lambda publish-version --function-name %q --profile %s --region %s --output json`,
		fullLambdaName, env.Profile, deployRegion,
	)
	publishOut, err := exec.Command(publishScript)
	if err != nil {
		return fmt.Errorf("lambda publish-version failed: %w\n%s", err, publishOut)
	}

	var publishResult struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal([]byte(publishOut), &publishResult); err != nil {
		return fmt.Errorf("failed to parse publish-version response: %w", err)
	}

	log.Info("Fetching aliases...")
	listAliasesScript := fmt.Sprintf(
		`aws lambda list-aliases --function-name %q --profile %s --region %s --output json`,
		fullLambdaName, env.Profile, deployRegion,
	)
	listOut, err := exec.Command(listAliasesScript)
	if err != nil {
		return fmt.Errorf("lambda list-aliases failed: %w\n%s", err, listOut)
	}

	var aliasesResult struct {
		Aliases []struct {
			Name string `json:"Name"`
		} `json:"Aliases"`
	}
	if err := json.Unmarshal([]byte(listOut), &aliasesResult); err != nil {
		return fmt.Errorf("failed to parse list-aliases response: %w", err)
	}

	if len(aliasesResult.Aliases) == 0 {
		log.Infof("Function %s has no aliases, skipping alias update", functionName)
		return nil
	}

	if len(aliasesResult.Aliases) > 1 {
		return fmt.Errorf("function %s has multiple aliases, cannot determine which one to update", functionName)
	}

	aliasName := aliasesResult.Aliases[0].Name

	log.Infof("Updating alias %q → version %s...", aliasName, publishResult.Version)
	updateAliasScript := fmt.Sprintf(
		`aws lambda update-alias --function-name %q --name %s --function-version %s --profile %s --region %s`,
		fullLambdaName, aliasName, publishResult.Version, env.Profile, deployRegion,
	)
	if out, err := exec.Command(updateAliasScript); err != nil {
		return fmt.Errorf("lambda update-alias failed: %w\n%s", err, out)
	}

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
