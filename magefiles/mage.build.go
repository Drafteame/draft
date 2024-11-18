package main

import (
	"errors"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"magefiles/build"
	"magefiles/files"
)

// Build is a namespace for build related targets.
type Build mg.Namespace

var awsEnvs = map[string]string{}

// Generated Execute automatic generation of code.
func (Build) Generated() error {
	return sh.Run("go", "generate", "./...")
}

// Local builds a local build starting from the given path inside the cmd folder.
func (Build) Local(path string) error {
	println("Building local: ", path)
	return build.Exec(path, build.LocalEnvType)
}

func (b Build) Run(path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path = cwd + "/" + path

	serviceName, parent, err := getService(path)
	if err != nil {
		return err
	}

	confFile, err := b.Pkl(parent)
	if err != nil {
		return err
	}

	println("Running:", path)

	if errLocal := b.Local(path); errLocal != nil {
		return errLocal
	}

	env, errEnv := getAWSEnvs()
	if errEnv != nil {
		return errEnv
	}

	env["STAGE"] = "dev"
	env["APP_NAME"] = serviceName
	env["DEBUG"] = "true"

	execName := ".bin/" + path

	cmd := "./" + execName

	args := []string{"--local", "--logger.colored"}

	if confFile != "" {
		println("Using config file:", confFile)
		args = append(args, "--config", confFile)
	}

	if err := sh.RunWith(env, cmd, args...); err != nil {
		out, errOut := sh.OutputWith(env, cmd, args...)
		if out != "" {
			println(out)
		}

		return errOut
	}

	return nil
}

func (b Build) Pkl(service string) (string, error) {
	mainFile := service + "/config/app/app.pkl"
	outFile := service + "/.app-config.yaml"

	println("Searching pkl file:", mainFile)

	if !files.Exists(mainFile) {
		return "", nil
	}

	println("Building pkl:", service)

	cmd := "pkl eval -f yaml " + mainFile + " > " + outFile

	envs, err := getAWSEnvs()
	if err != nil {
		return "", err
	}

	return outFile, sh.RunWith(envs, "bash", "-c", cmd)
}

func getAWSEnvs() (map[string]string, error) {
	if len(awsEnvs) > 0 {
		return awsEnvs, nil
	}

	awsProfile, err := files.GetAWSProfileCredentials("draftea-dev")
	if err != nil {
		return nil, err
	}

	accountID, err := getAWSAccountID()
	if err != nil {
		return nil, err
	}

	awsEnvs = map[string]string{
		"AWS_ACCESS_KEY_ID":     awsProfile["aws_access_key_id"],
		"AWS_SECRET_ACCESS_KEY": awsProfile["aws_secret_access_key"],
		"AWS_DEFAULT_REGION":    "us-east-2",
		"AWS_ACCOUNT":           accountID,
		"STAGE":                 "dev",
	}

	return awsEnvs, nil
}

func getAWSAccountID() (string, error) {
	cmd := "aws sts get-caller-identity --query Account --output text --profile draftea-dev"
	output, err := sh.Output("bash", "-c", cmd)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func getService(path string) (string, string, error) {
	parts := strings.Split(path, "/")
	// traverse the path backwards to find the service name
	for i := len(parts) - 1; i >= 0; i-- {
		parent := strings.Join(parts[:i+1], "/")

		if isServerlessService(parent) {
			serviceName := getServiceNameFormServerlessFile(parent)
			return serviceName, parent, nil
		}

		if isCDKService(parent) {
			serviceName := getServiceNameFormCDKFile(parent)
			return serviceName, parent, nil
		}
	}

	return "", "", errors.New("service is not serverless or cdk")
}

func getServiceNameFormServerlessFile(path string) string {
	file := path + "/serverless.yml"

	type serverless struct {
		Service string `yaml:"service"`
	}

	data := serverless{}

	if err := files.LoadYAML(file, &data); err != nil {
		panic("Failed to load serverless.yml: " + err.Error())
	}

	return data.Service
}

func isServerlessService(path string) bool {
	if !files.Exists(path) {
		return false
	}

	return files.Exists(path + "/serverless.yml")
}

func isCDKService(path string) bool {
	if !files.Exists(path) {
		return false
	}

	return files.Exists(path + "/cdk")
}

func getServiceNameFormCDKFile(path string) string {
	file := path + "/cdk/cdk.json"

	type cdk struct {
		Service string `json:"service"`
	}

	data := cdk{}

	if err := files.LoadJSON(file, &data); err != nil {
		panic("Failed to load cdk.json: " + err.Error())
	}

	return data.Service
}
