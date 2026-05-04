package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

// DeployResult holds the outcome of deploying a single service.
type DeployResult struct {
	Name string // original arg (name or path as passed by user)
	Err  error
}

// DeployService deploys one or more services (by name or path) using the given env config.
// If force is true, .deployignore is ignored.
func DeployService(env EnvConfig, args []string, force bool) []DeployResult {
	results := make([]DeployResult, 0, len(args))

	log.Info("Fetching AWS Account ID...")
	accountID, err := aws.GetAccountID(env.Profile)
	if err != nil {
		err = fmt.Errorf("failed to get AWS account ID: %w", err)
		for _, arg := range args {
			results = append(results, DeployResult{Name: arg, Err: err})
		}
		return results
	}

	for _, arg := range args {
		log.Infof("\n─── Deploying: %s ───", arg)

		absPath, err := resolveService(arg)
		if err != nil {
			results = append(results, DeployResult{Name: arg, Err: err})
			continue
		}

		err = deployServiceToDir(env, absPath, accountID, force)
		results = append(results, DeployResult{Name: arg, Err: err})
	}

	return results
}

func deployServiceToDir(env EnvConfig, absPath, accountID string, force bool) error {
	if err := validateServiceDir(absPath); err != nil {
		return err
	}

	if !force {
		skip, reason, err := shouldSkipForStage(absPath, env.Stage())
		if err != nil {
			return err
		}
		if skip {
			log.Warnf("Skipping %s: %s (use --force to override)", absPath, reason)
			return nil
		}
	}

	stage := env.Stage()
	slsParams := fmt.Sprintf("--aws-profile=%s", env.Profile)
	if env.ExtraSLSParams != "" {
		slsParams = fmt.Sprintf("%s %s", slsParams, env.ExtraSLSParams)
	}

	syncSecretsDry := "false"
	if env.SyncSecretsDry {
		syncSecretsDry = "true"
	}

	script := fmt.Sprintf(
		`cd %q && npm install && env STAGE=%s AWS_ACCOUNT=%s SLS_PARAMS=%q SLS_SYNC_SECRETS_DRY=%s npm run deploy`,
		absPath, stage, accountID, slsParams, syncSecretsDry,
	)

	if _, err := exec.Command(script, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
		return fmt.Errorf("deploy failed: %w", err)
	}

	log.Successf("✓ Deployed: %s", absPath)
	return nil
}

// validateServiceDir asserts that absPath looks like a deployable serverless
// service. It returns nil when serverless.yml is present and a descriptive
// error otherwise. .deployignore handling is done separately by
// shouldSkipForStage so that --force can bypass the skip without bypassing
// this structural check.
func validateServiceDir(absPath string) error {
	if !files.Exists(filepath.Join(absPath, "serverless.yml")) {
		return fmt.Errorf("serverless.yml not found in %s", absPath)
	}
	return nil
}
