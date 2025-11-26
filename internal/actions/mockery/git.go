package mockery

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

func (m *Mockery) findConfigsFromStaged() ([]string, error) {
	return m.findConfigsFromGitSource("staged", m.getStagedFiles)
}

func (m *Mockery) findConfigsFromCommitted() ([]string, error) {
	return m.findConfigsFromGitSource("committed", m.getCommittedFiles)
}

func (m *Mockery) findConfigsFromModified() ([]string, error) {
	stagedFiles, err := m.getStagedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	committedFiles, err := m.getCommittedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get committed files: %w", err)
	}

	allFiles := append(stagedFiles, committedFiles...)

	if len(allFiles) == 0 {
		log.Info("No staged or committed files found")
		return nil, nil
	}

	configs := m.findConfigsFromFiles(allFiles)

	if len(configs) == 0 {
		log.Warn("No .mockery.pkg.yml files found in modified directories")
	}

	return configs, nil
}

func (m *Mockery) findConfigsFromGitSource(name string, getFiles func() ([]string, error)) ([]string, error) {
	files, err := getFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s files: %w", name, err)
	}

	if len(files) == 0 {
		log.Infof("No %s files found", name)
		return nil, nil
	}

	configs := m.findConfigsFromFiles(files)

	if len(configs) == 0 {
		log.Warnf("No .mockery.pkg.yml files found in %s directories", name)
	}

	return configs, nil
}

func (m *Mockery) getStagedFiles() ([]string, error) {
	cmd := "git diff --cached --name-only --diff-filter=AM"
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run git diff --cached: %w\nOutput: %s\nTip: Ensure you're in a git repository", err, output)
	}

	return parseGitOutput(output), nil
}

func (m *Mockery) getCommittedFiles() ([]string, error) {
	mainBranch := m.detectMainBranch()

	cmd := fmt.Sprintf("git diff --name-only --diff-filter=AM origin/%s...HEAD", mainBranch)
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w\nOutput: %s\nTip: Ensure you're in a git repository and origin/%s exists", err, output, mainBranch)
	}

	return parseGitOutput(output), nil
}

func (m *Mockery) detectMainBranch() string {
	if _, err := exec.Command("git rev-parse --verify origin/main"); err == nil {
		return "main"
	}

	if _, err := exec.Command("git rev-parse --verify origin/master"); err == nil {
		return "master"
	}

	return "main"
}

func parseGitOutput(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

func (m *Mockery) findConfigsFromFiles(filePaths []string) []string {
	directories := m.extractDirectories(filePaths)
	return m.searchConfigsInDirs(directories)
}

func (m *Mockery) extractDirectories(filePaths []string) []string {
	dirMap := make(map[string]struct{})

	for _, file := range filePaths {
		dir := filepath.Dir(file)
		if dir != "" && dir != "." {
			dirMap[dir] = struct{}{}
		}
	}

	var directories []string
	for dir := range dirMap {
		directories = append(directories, dir)
	}

	return directories
}

func (m *Mockery) searchConfigsInDirs(directories []string) []string {
	configMap := make(map[string]struct{})

	for _, dir := range directories {
		currentDir := dir
		for {
			configPath := filepath.Join(currentDir, pkgConfigSuffix)
			if files.Exists(configPath) {
				normalized, err := filepath.Abs(configPath)
				if err != nil {
					normalized = configPath
				}
				configMap[normalized] = struct{}{}
			}

			parent := filepath.Dir(currentDir)
			if parent == currentDir || parent == "." || parent == "/" {
				break
			}
			currentDir = parent
		}
	}

	var configs []string
	for config := range configMap {
		configs = append(configs, config)
	}

	return configs
}
