package build

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/sh"
	"golang.org/x/sync/errgroup"
	"magefiles/files"
)

type buildEntry struct {
	mainPath string
	outPath  string
}

const maxBuildCores = 3

// Exec builds binaries for all main.go files in the cmd folder git the given env variables.
func Exec(path string, env EnvType) error {
	fileList, err := getMainFiles(path)
	if err != nil {
		return err
	}

	entries, err := getBuildEntries(fileList)
	if err != nil {
		return err
	}

	buildCores := runtime.NumCPU() - 1

	if buildCores > maxBuildCores {
		buildCores = maxBuildCores
	}

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(buildCores)

	for _, entry := range entries {
		current := entry

		group.Go(func() error {
			return compileEntry(current, env.Envs())
		})
	}

	return group.Wait()
}

func compileEntry(entry buildEntry, env map[string]string) error {
	println("Compiling: ", entry.mainPath)

	if env == nil {
		env = make(map[string]string)
	}

	cmd := "go"
	args := []string{"build", "-tags", "local", "-ldflags", "-s -w", "-o", entry.outPath, entry.mainPath}

	if err := sh.RunWith(env, cmd, args...); err != nil {
		return err
	}

	println("Compiled: ", entry.outPath)

	return sh.Run("chmod", "+x", entry.outPath)
}

func getMainFiles(path string) ([]string, error) {
	omit := []string{
		"config",
		"node_modules",
		".serverless",
		".bin",
	}

	return files.Search(path, "main.go", files.WithOmit(omit...))
}

func getBuildEntries(fileList []string) ([]buildEntry, error) {
	entries := make([]buildEntry, 0, len(fileList))

	for _, file := range fileList {
		workingDir := filepath.Dir(file)
		mainPath := filepath.Join(workingDir, "main.go")

		entries = append(entries, buildEntry{
			mainPath: mainPath,
			outPath:  getOutPath(workingDir),
		})
	}

	return entries, nil
}

func getOutPath(workingDir string) string {
	splitPath := strings.Split(workingDir, string(filepath.Separator))
	cleanedPath := splitPath[1:]

	fullPath := []string{BinPath}
	fullPath = append(fullPath, cleanedPath...)

	return filepath.Join(fullPath...)
}
