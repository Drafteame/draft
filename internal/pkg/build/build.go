package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
)

type buildEntry struct {
	mainPath string
	outPath  string
}

// Exec builds binaries for all main.go files in the cmd folder git the given env variables.
func Exec(ctx context.Context, path string) error {
	fileList, err := getMainFiles(path)
	if err != nil {
		return err
	}

	entries, err := getBuildEntries(fileList)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)
	group.SetLimit(runtime.NumCPU())

	for _, entry := range entries {
		current := entry

		group.Go(func() error {
			return compileEntry(current)
		})
	}

	return group.Wait()
}

func compileEntry(entry buildEntry) error {
	cmd := "go build -tags local -ldflags=\"-s -w\" -o %s %s"
	cmd = fmt.Sprintf(cmd, entry.outPath, entry.mainPath)

	_, err := exec.Command(cmd, exec.WithEnvs(envs))
	if err != nil {
		return err
	}

	cmd = "chmod +x %s"
	cmd = fmt.Sprintf(cmd, entry.outPath)

	_, err = exec.Command(cmd)
	if err != nil {
		return err
	}

	return nil
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
