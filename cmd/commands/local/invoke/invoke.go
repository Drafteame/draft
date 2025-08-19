package invoke

import (
	"fmt"
	"maps"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/build"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/pkl"
)

var invokeCmd = cobra.Command{
	Use:   "local:invoke [flags] <path-to-lambda>",
	Short: "Invoke a lambda",
	Long:  "Invoke a lambda by compiling code and its configuration and running it locally",
	Run:   run,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
}

func GetCmd() *cobra.Command {
	return &invokeCmd
}

func init() {
	invokeCmd.Flags().StringP("body", "b", "", "Body to send to the lambda")
	invokeCmd.Flags().StringP("body-file", "f", "", "File to read body from")
}

func run(cmd *cobra.Command, args []string) {
	if workDir := cmd.Parent().Flag("working-dir").Value.String(); workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			panic(err)
		}
	}

	lambdaPath := getPath(args[0])
	body := getBody(cmd, lambdaPath)

	println("Getting AWS credentials...")

	awsEnvs := getAWSEnvs()

	println("Identifying service...")

	serviceName, parent, err := getService(lambdaPath)
	if err != nil {
		println("Failed to get service name:", err.Error())
		os.Exit(1)
	}

	pklOutFile := buildPkl(parent, awsEnvs)

	println("Building lambda...")

	if errBuild := build.Exec(cmd.Context(), lambdaPath); errBuild != nil {
		println("Failed to build:", errBuild.Error())
		os.Exit(1)
	}

	binCmd := buildBinCmd(lambdaPath, pklOutFile, body)

	serviceEnvs := map[string]string{
		"STAGE":    "dev",
		"APP_NAME": serviceName,
		"DEBUG":    "true",
	}

	maps.Copy(serviceEnvs, awsEnvs)

	execOpts := []exec.CommandOption{
		exec.WithEnvs(serviceEnvs),
		exec.WithStdout(os.Stdout),
		exec.WithStderr(os.Stderr),
	}

	println("Running lambda...")

	if _, errRun := exec.Command(binCmd, execOpts...); errRun != nil {
		_, _ = fmt.Printf("Failed to run command '%s': %v", binCmd, errRun)
		os.Exit(1)
	}
}

// buildBinCmd builds a binary execution command by appending optional config and body arguments to the given path.
func buildBinCmd(lambdaPath, pklOutFile, body string) string {
	execName := path.Join(build.BinPath, lambdaPath)

	if pklOutFile != "" {
		pklOutFile = fmt.Sprintf("--config %s", pklOutFile)
	}

	binCmd := fmt.Sprintf("./%s --local --logger.colored %s", execName, pklOutFile)

	if body != "" {
		binCmd = fmt.Sprintf("%s --body '%s'", binCmd, body)
	}

	return binCmd
}

// buildPkl generates a YAML configuration file from a specified PKL file and environment variables.
// It takes the parent directory path and a map of AWS environment variables as input.
// If the required PKL file does not exist, it returns an empty string.
// On success, it returns the path to the output YAML file.
func buildPkl(parent string, awsEnvs map[string]string) string {
	pklConfigFile := fmt.Sprintf("%s/config/app/app.pkl", parent)
	pklOutFile := fmt.Sprintf("%s/.app-config.yaml", parent)

	if !files.Exists(pklConfigFile) {
		return ""
	}

	println("Building pkl...")

	_, errPkl := pkl.Build(pklConfigFile, pklOutFile, "yaml", pkl.WithEnvs(awsEnvs))
	if errPkl != nil {
		println("Failed to build pkl:", errPkl.Error())
		os.Exit(1)
	}

	copyToEmbedDir(pklOutFile, parent)

	return pklOutFile
}

// copyToEmbedDir copies the generated YAML configuration file to the embed directory.
// It takes the path to the output YAML file and the parent directory path as input.
func copyToEmbedDir(pklOutFile, parent string) {
	if !files.Exists(path.Join(parent, "embed")) {
		if errEmbedCreate := dirs.Create(path.Join(parent, "embed")); errEmbedCreate != nil {
			println("Failed to create embed directory:", errEmbedCreate.Error())
			return
		}
	}

	if errCopy := files.Copy(pklOutFile, path.Join(parent, "embed", ".app-config.yaml")); errCopy != nil {
		println("Failed to copy pkl to embed dir:", errCopy.Error())
		return
	}
}

// getPath returns the full absolute path by appending the passed argument to the current working directory.
// It exits the program if the current working directory cannot be retrieved.
func getPath(argPath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		println("Failed to get current working directory:", err.Error())
		os.Exit(1)
	}

	println("Current working directory:", cwd)

	argPath = cwd + "/" + argPath

	return argPath
}

// getAWSEnvs retrieves AWS environment variables using local AWS credentials.
// The function exits the program if credentials cannot be retrieved.
func getAWSEnvs() map[string]string {
	awsEnvs, err := aws.GetLocalCredentials()
	if err != nil {
		println("Failed to get AWS credentials:", err.Error())
		os.Exit(1)
	}

	return awsEnvs
}

// getBody retrieves the body content for a command from either a file, a flag, or defaults to a local file in the
// service path.
func getBody(cmd *cobra.Command, servicePath string) string {
	bodyFile := cmd.Flag("body-file").Value.String()
	if bodyFile != "" {
		return readBodyFromFile(bodyFile)
	}

	body := cmd.Flag("body").Value.String()

	if body == "" {
		return searchLocalBodyFile(servicePath)
	}

	return body
}

// searchLocalBodyFile checks for the existence of "local-body.json" in the provided service path and reads its content.
// If the file does not exist, it returns an empty string.
func searchLocalBodyFile(servicePath string) string {
	filePath := path.Join(servicePath, "local-body.json")

	if !files.Exists(filePath) {
		return ""
	}

	return readBodyFromFile(filePath)
}

// readBodyFromFile reads the content of the specified file, removes newlines and tabs, and returns it as a string.
// It exits the program with an error if the file does not exist or cannot be read.
func readBodyFromFile(file string) string {
	if !files.Exists(file) {
		println("Failed to get body:", fmt.Errorf("file not found: %s", file).Error())
		os.Exit(1)
	}

	body, err := files.Read(file)
	if err != nil {
		println("Failed to read body:", err.Error())
		os.Exit(1)
	}

	strBody := string(body)
	strBody = strings.ReplaceAll(strBody, "\n", "")
	strBody = strings.ReplaceAll(strBody, "\t", "")

	return strBody
}
