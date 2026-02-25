package mockery

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const (
	baseConfigFile  = ".mockery.base.yml"
	pkgConfigSuffix = ".mockery.pkg.yml"
	tmpConfigPrefix = ".mockery.tmp."
	tmpConfigSuffix = ".yml"
)

type Mockery struct {
	ctx         context.Context
	configFiles []string
	jobsNum     int
	dry         bool
	gitMod      bool
	tmpFiles    []string
}

func New(ctx context.Context, configFiles []string, jobsNum int, dry bool, gitMod bool) *Mockery {
	return &Mockery{
		ctx:         ctx,
		configFiles: configFiles,
		jobsNum:     jobsNum,
		dry:         dry,
		gitMod:      gitMod,
		tmpFiles:    make([]string, 0),
	}
}

func (m *Mockery) Exec() error {
	log.Info("Running mockery...")

	if err := m.validate(); err != nil {
		return err
	}

	if m.dry {
		log.Info("Running in dry-run mode - mockery commands will not be executed")
	}

	startTime := time.Now()

	configFiles, err := m.resolveConfigFiles()
	if err != nil {
		return err
	}

	if len(configFiles) == 0 {
		log.Warn("No .mockery.pkg.yml files found")
		log.Info("Tip: Create package-specific configs with the naming pattern: .mockery.pkg.yml")
		return nil
	}

	baseConfig, err := m.loadBaseConfig()
	if err != nil {
		return err
	}

	// Merge all package configs into a single unified config and run mockery once.
	// This avoids N separate mockery process startups and duplicate `go list` + type-checking
	// costs across packages.
	mergedConfig, err := m.mergeSingleConfig(configFiles, baseConfig)
	if err != nil {
		return err
	}

	tmpFile, err := m.createSingleTempConfig(mergedConfig, len(configFiles))
	if err != nil {
		return err
	}

	defer m.cleanup()

	execErr := m.runSingleInvocation(tmpFile, len(configFiles))

	if m.ctx.Err() != nil {
		log.Warnf("⚠ Operation cancelled by user after %.2fs", time.Since(startTime).Seconds())
		return fmt.Errorf("operation cancelled: %w", m.ctx.Err())
	}

	duration := time.Since(startTime)

	if execErr != nil {
		log.Errorf("✗ mockery failed for %d package(s) (%.2fs)", len(configFiles), duration.Seconds())
		log.Errorf("  %v", execErr)
		log.Info("Tip: Check the error messages above for details on how to fix the configurations")
		return execErr
	}

	if m.dry {
		log.Successf("✓ All %d package(s) validated successfully (%.2fs)", len(configFiles), duration.Seconds())
		log.Info("Dry run completed - no mockery commands were executed")
	} else {
		log.Successf("✓ All %d package(s) completed successfully (%.2fs)", len(configFiles), duration.Seconds())
	}

	return nil
}

// validate validates input parameters.
func (m *Mockery) validate() error {
	if m.jobsNum <= 0 {
		return fmt.Errorf("invalid --jobs-num value: %d (must be greater than 0)", m.jobsNum)
	}

	if m.gitMod && len(m.configFiles) > 0 {
		return fmt.Errorf("cannot use --git-mod with explicit config file paths")
	}

	return nil
}

// resolveConfigFiles resolves, validates, and deduplicates config files.
func (m *Mockery) resolveConfigFiles() ([]string, error) {
	log.Info("Resolving config files...")
	var configFiles []string

	if m.gitMod {
		found, err := m.findConfigsFromGitDiff()
		if err != nil {
			return nil, err
		}
		configFiles = found
	} else if len(m.configFiles) > 0 {
		validated, err := m.validateProvidedConfigs(m.configFiles)
		if err != nil {
			return nil, err
		}
		configFiles = validated
	} else {
		found, err := m.findPackageConfigs()
		if err != nil {
			return nil, err
		}
		configFiles = found
	}

	configFiles = m.deduplicateFiles(configFiles)

	return configFiles, nil
}

// validateProvidedConfigs validates user-provided config file paths.
func (m *Mockery) validateProvidedConfigs(paths []string) ([]string, error) {
	var validated []string

	for _, path := range paths {
		if !files.Exists(path) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}

		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", path, err)
		}

		if stat.IsDir() {
			return nil, fmt.Errorf("path is a directory, not a file: %s", path)
		}

		if !strings.HasSuffix(path, pkgConfigSuffix) {
			log.Warnf("Config file %s doesn't follow naming convention (.mockery.pkg.yml)", path)
		}

		validated = append(validated, path)
	}

	return validated, nil
}

// deduplicateFiles removes duplicate file paths.
func (m *Mockery) deduplicateFiles(fileList []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, file := range fileList {
		normalized, err := filepath.Abs(file)
		if err != nil {
			normalized = file
		}

		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, file)
		}
	}

	if len(result) < len(fileList) {
		log.Warnf("Removed %d duplicate config file(s)", len(fileList)-len(result))
	}

	return result
}

// findPackageConfigs searches for all .mockery.pkg.yml files in the project.
func (m *Mockery) findPackageConfigs() ([]string, error) {
	var configs []string

	searchErr := dirs.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path == "." {
				return nil
			}
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(path, pkgConfigSuffix) {
			configs = append(configs, path)
		}

		return nil
	})

	if searchErr != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", searchErr)
	}

	return configs, nil
}

// findConfigsFromGitDiff finds .mockery.pkg.yml files in directories with modified files.
func (m *Mockery) findConfigsFromGitDiff() ([]string, error) {
	modifiedFiles, err := m.getModifiedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	if len(modifiedFiles) == 0 {
		log.Info("No modified files found in git diff")
		return nil, nil
	}

	directories := m.extractDirectories(modifiedFiles)
	configs := m.findConfigsInDirectories(directories)

	if len(configs) == 0 {
		log.Warn("No .mockery.pkg.yml files found in modified directories")
	}

	return configs, nil
}

// getModifiedFiles returns a list of modified files using git diff.
func (m *Mockery) getModifiedFiles() ([]string, error) {
	mainBranch := "main"

	checkCmd := "git rev-parse --verify origin/main"
	if _, err := exec.Command(checkCmd); err != nil {
		checkCmd = "git rev-parse --verify origin/master"
		if _, err := exec.Command(checkCmd); err == nil {
			mainBranch = "master"
		}
	}

	cmd := fmt.Sprintf("git diff --name-only --diff-filter=AM origin/%s...HEAD", mainBranch)
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w\nOutput: %s\nTip: Ensure you're in a git repository and origin/%s exists", err, output, mainBranch)
	}

	fileList := strings.Split(strings.TrimSpace(output), "\n")
	var result []string
	for _, file := range fileList {
		file = strings.TrimSpace(file)
		if file != "" {
			result = append(result, file)
		}
	}

	return result, nil
}

// extractDirectories extracts unique directories from file paths.
func (m *Mockery) extractDirectories(fileList []string) []string {
	dirMap := make(map[string]struct{})

	for _, file := range fileList {
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

// findConfigsInDirectories searches for .mockery.pkg.yml files in directories and their parents.
func (m *Mockery) findConfigsInDirectories(directories []string) []string {
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

// loadBaseConfig loads the base configuration file.
func (m *Mockery) loadBaseConfig() (map[string]any, error) {
	log.Info("Loading base configuration file...")
	if !files.Exists(baseConfigFile) {
		return make(map[string]any), nil
	}

	data, err := files.Read(baseConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read base config %s: %w", baseConfigFile, err)
	}

	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse base config %s: %w (check YAML syntax)", baseConfigFile, err)
	}

	return config, nil
}

// mergeSingleConfig merges all package configs into a single unified config.
// The base config provides all top-level settings (dir, filename, etc.) and each
// pkg config contributes only its packages: entries. This works correctly when dir
// uses mockery template variables like {{.InterfaceDir}}, which resolve per-interface
// at runtime regardless of where mockery is invoked from.
func (m *Mockery) mergeSingleConfig(configFiles []string, baseConfig map[string]any) (map[string]any, error) {
	log.Infof("Merging %d package config(s) into single mockery invocation...", len(configFiles))

	result := m.deepMerge(baseConfig, map[string]any{})
	allPackages := make(map[string]any)

	for _, configFile := range configFiles {
		var pkgConfig map[string]any
		if err := files.LoadYAML(configFile, &pkgConfig); err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", configFile, err)
		}

		pkgs, ok := pkgConfig["packages"]
		if !ok {
			continue
		}

		pkgMap, ok := pkgs.(map[string]any)
		if !ok {
			continue
		}

		for pkgPath, pkgDef := range pkgMap {
			if existing, dup := allPackages[pkgPath]; dup {
				// Same package in multiple configs: deep merge their definitions so
				// interfaces from both are generated.
				if existingMap, ok := existing.(map[string]any); ok {
					if newMap, ok := pkgDef.(map[string]any); ok {
						allPackages[pkgPath] = m.deepMerge(existingMap, newMap)
						continue
					}
				}
				log.Warnf("Duplicate package %q in %s, overwriting previous definition", pkgPath, configFile)
			}
			allPackages[pkgPath] = pkgDef
		}
	}

	if len(allPackages) > 0 {
		result["packages"] = allPackages
	}

	return result, nil
}

// createSingleTempConfig writes the merged config to a single temporary file.
func (m *Mockery) createSingleTempConfig(config map[string]any, pkgCount int) (string, error) {
	randID, err := generateRandomID()
	if err != nil {
		return "", fmt.Errorf("failed to generate temp file name: %w", err)
	}

	tmpFile := fmt.Sprintf("%smerged.%s%s", tmpConfigPrefix, randID, tmpConfigSuffix)

	data, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged config: %w", err)
	}

	if err := files.Create(tmpFile, data); err != nil {
		return "", fmt.Errorf("failed to write temp config %s: %w", tmpFile, err)
	}

	m.tmpFiles = append(m.tmpFiles, tmpFile)

	log.Debugf("Merged config written to %s (%d packages)", tmpFile, pkgCount)

	return tmpFile, nil
}

// runSingleInvocation runs mockery once with the merged config.
func (m *Mockery) runSingleInvocation(tmpFile string, pkgCount int) error {
	var execErr error
	doneChan := make(chan struct{})

	spin := spinner.New().Title(fmt.Sprintf("Running mockery for %d package(s)...", pkgCount))

	action := func() {
		defer close(doneChan)

		select {
		case <-m.ctx.Done():
			execErr = m.ctx.Err()
			return
		default:
		}

		execErr = m.runMockery(tmpFile, tmpFile)
	}

	if err := spin.Action(action).Run(); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	<-doneChan

	return execErr
}

// deepMerge performs a deep merge of two maps, with b taking precedence.
func (m *Mockery) deepMerge(a, b map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range a {
		result[k] = v
	}

	for k, v := range b {
		if vMap, ok := v.(map[string]any); ok {
			if aMap, ok := result[k].(map[string]any); ok {
				result[k] = m.deepMerge(aMap, vMap)
				continue
			}
		}
		result[k] = v
	}

	return result
}

// runMockery executes mockery with the given config file.
func (m *Mockery) runMockery(configPath, originalPath string) error {
	if m.dry {
		log.Debugf("Dry run: would execute mockery --config %s", configPath)
		return nil
	}

	command := fmt.Sprintf("mockery --config %s", configPath)
	output, err := exec.Command(command)
	if err != nil {
		return fmt.Errorf("mockery failed for %s: %w\nOutput: %s\nTip: Check the config syntax and package paths", originalPath, err, output)
	}

	return nil
}

// cleanup removes all temporary config files.
func (m *Mockery) cleanup() {
	if len(m.tmpFiles) == 0 {
		return
	}

	var failed int
	for _, tmpFile := range m.tmpFiles {
		if err := os.Remove(tmpFile); err != nil && !os.IsNotExist(err) {
			failed++
		}
	}

	if failed > 0 {
		log.Warnf("Failed to clean up %d temporary file(s)", failed)
	}
}

// generateRandomID generates a random hex string for temp file naming.
func generateRandomID() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(b), nil
}
