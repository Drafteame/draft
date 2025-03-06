package invoke

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/build"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/pkl"
)

var invokeCmd = cobra.Command{
	Use:   "invoke [flags] <path-to-lambda>",
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

	path := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		println("Failed to get current working directory:", err.Error())
		os.Exit(1)
	}

	println("Current working directory:", cwd)
	println("Getting AWS credentials...")

	awsEnvs, err := aws.GetLocalCredentials()
	if err != nil {
		println("Failed to get AWS credentials:", err.Error())
		os.Exit(1)
	}

	path = cwd + "/" + path

	println("Identifying service...")

	serviceName, parent, err := getService(path)
	if err != nil {
		println("Failed to get service name:", err.Error())
		os.Exit(1)
	}

	pklConfigFile := fmt.Sprintf("%s/config/app/app.pkl", parent)
	pklOutFile := fmt.Sprintf("%s/.app-config.yaml", parent)

	if files.Exists(pklConfigFile) {
		println("Building pkl...")
		_, errPkl := pkl.Build(pklConfigFile, pklOutFile, "yaml", pkl.WithEnvs(awsEnvs))
		if errPkl != nil {
			println("Failed to build pkl:", errPkl.Error())
			os.Exit(1)
		}
	}

	println("Building lambda...")

	if errBuild := build.Exec(cmd.Context(), path); errBuild != nil {
		println("Failed to build:", errBuild.Error())
		os.Exit(1)
	}

	execName := build.BinPath + path

	binCmd := fmt.Sprintf("./%s --local --logger.colored --config %s", execName, pklOutFile)

	body, err := getBody(cmd)
	if err != nil {
		println("Failed to get body:", err.Error())
		os.Exit(1)
	}

	if body != "" {
		binCmd = fmt.Sprintf("%s --body '%s'", binCmd, body)
	}

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

func getBody(cmd *cobra.Command) (string, error) {
	bodyFile := cmd.Flag("body-file").Value.String()
	if bodyFile != "" {
		return readBodyFromFile(bodyFile)
	}

	body := cmd.Flag("body").Value.String()

	return body, nil
}

func readBodyFromFile(file string) (string, error) {
	if !files.Exists(file) {
		return "", fmt.Errorf("file not found: %s", file)
	}

	body, err := files.Read(file)
	if err != nil {
		return "", err
	}

	strBody := string(body)
	strBody = strings.ReplaceAll(strBody, "\n", "")
	strBody = strings.ReplaceAll(strBody, "\t", "")

	return strBody, nil
}
