package mockery

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

type (
	Mockery struct {
		ctx         context.Context
		configFiles []string
		jobsNum     int
		dry         bool
		gitMod      bool
		tmpFiles    []string // Track for cleanup
		mu          sync.Mutex
	}

	mockeryJob struct {
		configFile string
		tmpFile    string
		err        error
		duration   time.Duration
	}

	executionStats struct {
		total     int
		succeeded int
		failed    int
		duration  time.Duration
	}

	progressUpdate struct {
		current    int
		total      int
		configFile string
		success    bool
		err        error
		duration   time.Duration
	}
)

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
	// Validate inputs
	if err := m.validate(); err != nil {
		return err
	}

	if m.dry {
		log.Info("Running in dry-run mode - mockery commands will not be executed")
	}

	startTime := time.Now()

	// Find and validate config files
	configFiles, err := m.resolveConfigFiles()
	if err != nil {
		return err
	}

	if len(configFiles) == 0 {
		log.Warn("No .mockery.pkg.yml files found")
		log.Info("Tip: Create package-specific configs with the naming pattern: .mockery.pkg.yml")
		return nil
	}

	// Load base config
	baseConfig, err := m.loadBaseConfig()
	if err != nil {
		return err
	}

	// Create temporary config files
	if errTmp := m.createTempConfigs(configFiles, baseConfig); errTmp != nil {
		m.cleanup() // Clean up any temp files created before the error
		return errTmp
	}

	// Ensure cleanup on exit
	defer m.cleanup()

	// Execute mockery concurrently with the progress spinner
	results := m.executeConcurrentWithProgress(configFiles)

	// Check if context was canceled
	if m.ctx.Err() != nil {
		stats := m.calculateStats(results, startTime)
		m.displayCancellationSummary(stats, results)
		return fmt.Errorf("operation cancelled: %w", m.ctx.Err())
	}

	// Calculate and display stats
	stats := m.calculateStats(results, startTime)
	m.displaySummary(stats, results)

	// Return error if any executions failed
	if stats.failed > 0 {
		return fmt.Errorf("mockery execution failed for %d package(s)", stats.failed)
	}

	return nil
}

// validate validates input parameters
func (m *Mockery) validate() error {
	if m.jobsNum <= 0 {
		return fmt.Errorf("invalid --jobs-num value: %d (must be greater than 0)", m.jobsNum)
	}

	if m.jobsNum > 100 {
		log.Warnf("Very high concurrency (%d) may cause performance issues", m.jobsNum)
	}

	if m.gitMod && len(m.configFiles) > 0 {
		return fmt.Errorf("cannot use --git-mod with explicit config file paths")
	}

	return nil
}

// resolveConfigFiles resolves, validates, and deduplicates config files
func (m *Mockery) resolveConfigFiles() ([]string, error) {
	log.Info("Resolving config files...")
	var configFiles []string

	// If git-mod is enabled, find configs from modified files
	if m.gitMod {
		found, err := m.findConfigsFromGitDiff()
		if err != nil {
			return nil, err
		}
		configFiles = found
	} else if len(m.configFiles) > 0 {
		// If files provided, validate them
		validated, err := m.validateProvidedConfigs(m.configFiles)
		if err != nil {
			return nil, err
		}
		configFiles = validated
	} else {
		// Search for config files
		found, err := m.findPackageConfigs()
		if err != nil {
			return nil, err
		}
		configFiles = found
	}

	// Deduplicate
	configFiles = m.deduplicateFiles(configFiles)

	return configFiles, nil
}

// validateProvidedConfigs validates user-provided config file paths
func (m *Mockery) validateProvidedConfigs(paths []string) ([]string, error) {
	var validated []string

	for _, path := range paths {
		// Check if a file exists
		if !files.Exists(path) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}

		// Check if it's a file
		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", path, err)
		}

		if stat.IsDir() {
			return nil, fmt.Errorf("path is a directory, not a file: %s", path)
		}

		// Warn if it doesn't follow naming convention
		if !strings.HasSuffix(path, pkgConfigSuffix) {
			log.Warnf("Config file %s doesn't follow naming convention (.mockery.pkg.yml)", path)
		}

		validated = append(validated, path)
	}

	return validated, nil
}

// deduplicateFiles removes duplicate file paths
func (m *Mockery) deduplicateFiles(files []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, file := range files {
		// Normalize a path
		normalized, err := filepath.Abs(file)
		if err != nil {
			normalized = file
		}

		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, file)
		}
	}

	if len(result) < len(files) {
		log.Warnf("Removed %d duplicate config file(s)", len(files)-len(result))
	}

	return result
}

// findPackageConfigs searches for all .mockery.pkg.yml files in the project
func (m *Mockery) findPackageConfigs() ([]string, error) {
	var configs []string

	searchErr := dirs.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories (but not the root "." directory)
		if info.IsDir() {
			// Don't skip the root directory
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

// findConfigsFromGitDiff finds .mockery.pkg.yml files in directories with modified files
func (m *Mockery) findConfigsFromGitDiff() ([]string, error) {
	// Get modified files comparing HEAD with the main branch
	modifiedFiles, err := m.getModifiedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	if len(modifiedFiles) == 0 {
		log.Info("No modified files found in git diff")
		return nil, nil
	}

	// Extract and deduplicate base directories
	directories := m.extractDirectories(modifiedFiles)

	// Search for .mockery.pkg.yml in each directory and its parents
	configs := m.findConfigsInDirectories(directories)

	if len(configs) == 0 {
		log.Warn("No .mockery.pkg.yml files found in modified directories")
	}

	return configs, nil
}

// getModifiedFiles returns a list of modified files using git diff
func (m *Mockery) getModifiedFiles() ([]string, error) {
	// Try to get the main branch name
	mainBranch := "main"

	// Check if origin/main exists
	checkCmd := "git rev-parse --verify origin/main"
	if _, err := exec.Command(checkCmd); err != nil {
		// Try master as fallback
		checkCmd = "git rev-parse --verify origin/master"
		if _, err := exec.Command(checkCmd); err == nil {
			mainBranch = "master"
		}
	}

	// Get modified and new files
	cmd := fmt.Sprintf("git diff --name-only --diff-filter=AM origin/%s...HEAD", mainBranch)
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w\nOutput: %s\nTip: Ensure you're in a git repository and origin/%s exists", err, output, mainBranch)
	}

	// Parse output into a file list
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

// extractDirectories extracts unique directories from file paths
func (m *Mockery) extractDirectories(files []string) []string {
	dirMap := make(map[string]struct{})

	for _, file := range files {
		// Get the directory of the file
		dir := filepath.Dir(file)
		if dir != "" && dir != "." {
			dirMap[dir] = struct{}{}
		}
	}

	// Convert map to slice
	var directories []string
	for dir := range dirMap {
		directories = append(directories, dir)
	}

	return directories
}

// findConfigsInDirectories searches for .mockery.pkg.yml files in directories and their parents
func (m *Mockery) findConfigsInDirectories(directories []string) []string {
	configMap := make(map[string]struct{})

	for _, dir := range directories {
		// Check the current directory and walk up to find .mockery.pkg.yml
		currentDir := dir
		for {
			configPath := filepath.Join(currentDir, pkgConfigSuffix)
			if files.Exists(configPath) {
				// Normalize a path
				normalized, err := filepath.Abs(configPath)
				if err != nil {
					normalized = configPath
				}
				configMap[normalized] = struct{}{}
			}

			// Move to the parent directory
			parent := filepath.Dir(currentDir)
			if parent == currentDir || parent == "." || parent == "/" {
				break
			}
			currentDir = parent
		}
	}

	// Convert map to slice
	var configs []string
	for config := range configMap {
		configs = append(configs, config)
	}

	return configs
}

// loadBaseConfig loads the base configuration file
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

// deepMerge performs a deep merge of two maps, with b taking precedence
func (m *Mockery) deepMerge(a, b map[string]any) map[string]any {
	result := make(map[string]any)

	// Copy all from a
	for k, v := range a {
		result[k] = v
	}

	// Merge from b
	for k, v := range b {
		if vMap, ok := v.(map[string]any); ok {
			// If b[k] is a map, try to deep merge with a[k]
			if aMap, ok := result[k].(map[string]any); ok {
				result[k] = m.deepMerge(aMap, vMap)
				continue
			}
		}
		// Otherwise, b[k] overwrites result[k]
		result[k] = v
	}

	return result
}

// createTempConfigs creates temporary config files for each package
func (m *Mockery) createTempConfigs(configFiles []string, baseConfig map[string]any) error {
	log.Info("Creating temporary config files...")
	for i, configFile := range configFiles {
		var pkgConfig map[string]any

		if errLoad := files.LoadYAML(configFile, &pkgConfig); errLoad != nil {
			return fmt.Errorf("failed to load %s: %w", configFile, errLoad)
		}

		// Deep merge configs (package takes precedence)
		merged := m.deepMerge(baseConfig, pkgConfig)

		// Generate temp file name with context
		tmpFile, err := m.generateTempFileName(configFile, i)
		if err != nil {
			return fmt.Errorf("failed to generate temp file name: %w", err)
		}

		// Marshal merged config
		mergedData, err := yaml.Marshal(merged)
		if err != nil {
			return fmt.Errorf("failed to marshal merged config for %s: %w", configFile, err)
		}

		// Write the temp file
		if errCreate := files.Create(tmpFile, mergedData); errCreate != nil {
			return fmt.Errorf("failed to write temp config %s: %w", tmpFile, errCreate)
		}

		// Track for cleanup
		m.mu.Lock()
		m.tmpFiles = append(m.tmpFiles, tmpFile)
		m.mu.Unlock()
	}

	return nil
}

// generateTempFileName generates a descriptive temp file name
func (m *Mockery) generateTempFileName(configFile string, index int) (string, error) {
	randID, err := generateRandomID()
	if err != nil {
		return "", err
	}

	// Extract a meaningful name from the config file path
	dir := filepath.Dir(configFile)
	pkgName := filepath.Base(dir)
	if pkgName == "." {
		pkgName = "root"
	}

	// Create a descriptive temp file name: .mockery.tmp.[pkg-name].[idx].[rand].yml
	tmpFile := fmt.Sprintf("%s%s.%d.%s%s", tmpConfigPrefix, pkgName, index, randID, tmpConfigSuffix)

	return tmpFile, nil
}

// executeConcurrentWithProgress executes mockery commands concurrently with progress spinner
func (m *Mockery) executeConcurrentWithProgress(configFiles []string) []mockeryJob {
	log.Info("Executing mockery commands...")
	var (
		wg           sync.WaitGroup
		results      = make([]mockeryJob, 0, len(configFiles))
		resultsChan  = make(chan mockeryJob, len(configFiles))
		semaphore    = make(chan struct{}, m.jobsNum)
		progressChan = make(chan progressUpdate, len(configFiles))
		completed    = 0
		execErr      error
		doneChan     = make(chan struct{})
	)

	total := len(configFiles)

	spin := spinner.New().Title(fmt.Sprintf("[0 / %d] Preparing...", total))

	action := func() {
		defer close(doneChan)

		var progressWg sync.WaitGroup
		var cancelled bool
		progressWg.Add(1)

		// Start a progress reader goroutine before spawning work goroutines
		go func() {
			defer progressWg.Done()
			// Process progress updates as they arrive
			for update := range progressChan {
				completed++
				status := "✓"
				if !update.success {
					status = "✗"
				}

				// Extract a shorter name from the config file path
				shortName := m.shortenConfigPath(update.configFile)

				spin.Title(fmt.Sprintf("[%s] [%2d / %d] %s (%.2fs)",
					status, completed, total, shortName, update.duration.Seconds()))
			}
		}()

		// Start goroutines to process configs
		for idx := range m.tmpFiles {
			// Check if context is cancelled before starting new goroutine
			if m.ctx.Err() != nil {
				if !cancelled {
					log.Warn("Operation cancelled by user, waiting for ongoing tasks to complete...")
					cancelled = true
					execErr = m.ctx.Err()
				}
				goto waitForCompletion
			}

			// Acquire semaphore before spawning a goroutine
			select {
			case semaphore <- struct{}{}:
				// Successfully acquired semaphore
			case <-m.ctx.Done():
				// Context cancelled while waiting for semaphore
				if !cancelled {
					log.Warn("Operation cancelled by user, waiting for ongoing tasks to complete...")
					cancelled = true
					execErr = m.ctx.Err()
				}
				goto waitForCompletion
			}

			wg.Add(1)
			go m.execute(
				idx,
				resultsChan,
				progressChan,
				&wg,
				semaphore,
				configFiles[idx],
				m.tmpFiles[idx],
				total,
			)
		}

	waitForCompletion:
		// Wait for all work goroutines to complete
		wg.Wait()
		close(progressChan)
		close(resultsChan)

		// Wait for the progress reader to finish processing all updates
		progressWg.Wait()
	}

	if err := spin.Action(action).Run(); err != nil {
		execErr = fmt.Errorf("execution error: %w", err)
	}

	// Wait for action to complete
	<-doneChan

	if execErr != nil && !errors.Is(execErr, context.Canceled) {
		log.Errorf("Execution encountered errors: %v", execErr)
	}

	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// execute runs a mockery command for a specific configuration file and communicates results and progress updates.
func (m *Mockery) execute(
	idx int,
	resultChan chan mockeryJob,
	progressChan chan progressUpdate,
	wg *sync.WaitGroup,
	sem chan struct{},
	configFile string,
	tmpFile string,
	total int,
) {
	defer wg.Done()
	defer func() { <-sem }() // Release semaphore after work

	// Check if context is cancelled before starting work
	select {
	case <-m.ctx.Done():
		// Context cancelled, skip execution
		result := mockeryJob{
			configFile: configFile,
			tmpFile:    tmpFile,
			err:        m.ctx.Err(),
			duration:   0,
		}
		resultChan <- result
		return
	default:
		// Continue with execution
	}

	startTime := time.Now()
	err := m.runMockery(tmpFile, configFile)
	duration := time.Since(startTime)

	result := mockeryJob{
		configFile: configFile,
		tmpFile:    tmpFile,
		err:        err,
		duration:   duration,
	}

	resultChan <- result

	// Send progress update
	progressChan <- progressUpdate{
		current:    idx + 1,
		total:      total,
		configFile: configFile,
		success:    err == nil,
		err:        err,
		duration:   duration,
	}
}

// shortenConfigPath extracts a meaningful short name from a config file path
func (m *Mockery) shortenConfigPath(configPath string) string {
	// Remove .mockery.pkg.yml suffix
	path := strings.TrimSuffix(configPath, "/.mockery.pkg.yml")

	// If a path is too long, show only the last 2-3 segments
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return ".../" + strings.Join(parts[len(parts)-3:], "/")
	}

	return path
}

// runMockery executes mockery with the given config file
func (m *Mockery) runMockery(configPath, originalPath string) error {
	if m.dry {
		// Dry run - skip execution and just validate the config file exists
		log.Debugf("Dry run: would execute mockery --config %s", configPath)
		return nil
	}

	command := fmt.Sprintf("mockery --config %s", configPath)
	output, err := exec.Command(command)
	if err != nil {
		// Provide more context in error
		return fmt.Errorf("mockery failed for %s: %w\nOutput: %s\nTip: Check the config syntax and package paths", originalPath, err, output)
	}
	return nil
}

// cleanup removes all temporary config files
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

// calculateStats calculates execution statistics
func (m *Mockery) calculateStats(results []mockeryJob, startTime time.Time) executionStats {
	stats := executionStats{
		total:    len(m.tmpFiles), // Use total tmpFiles, not just results length
		duration: time.Since(startTime),
	}

	for _, result := range results {
		// Don't count context cancellation as failure
		if errors.Is(result.err, context.Canceled) || errors.Is(result.err, context.DeadlineExceeded) {
			continue
		}

		if result.err != nil {
			stats.failed++
		} else {
			stats.succeeded++
		}
	}

	return stats
}

// displaySummary displays execution summary
func (m *Mockery) displaySummary(stats executionStats, results []mockeryJob) {
	if stats.failed > 0 {
		log.Errorf("✗ Failed: %d/%d packages (%.2fs)", stats.failed, stats.total, stats.duration.Seconds())
		log.Errorf("Failed packages:")

		for _, result := range results {
			if result.err != nil {
				log.Errorf("  • %s", result.configFile)
				log.Errorf("    %v", result.err)
			}
		}

		log.Info("Tip: Check the error messages above for details on how to fix the configurations")
	} else {
		if m.dry {
			log.Successf("✓ All %d package(s) validated successfully (%.2fs)", stats.total, stats.duration.Seconds())
			log.Info("Dry run completed - no mockery commands were executed")
		} else {
			log.Successf("✓ All %d package(s) completed successfully (%.2fs)", stats.total, stats.duration.Seconds())
		}
	}
}

// displayCancellationSummary displays summary when operation is cancelled
func (m *Mockery) displayCancellationSummary(stats executionStats, results []mockeryJob) {
	completed := stats.succeeded + stats.failed
	cancelled := stats.total - completed

	log.Warnf("⚠ Operation cancelled by user")
	log.Infof("Completed: %d/%d packages", completed, stats.total)
	log.Infof("Cancelled: %d packages", cancelled)
	log.Infof("Duration: %.2fs", stats.duration.Seconds())

	if stats.failed > 0 {
		log.Warnf("Failed packages before cancellation:")

		for _, result := range results {
			if result.err != nil && !errors.Is(result.err, context.Canceled) && !errors.Is(result.err, context.DeadlineExceeded) {
				log.Errorf("  • %s", result.configFile)
				log.Errorf("    %v", result.err)
			}
		}
	}

	log.Info("Temporary files have been cleaned up")
}

// generateRandomID generates a random hex string for temp file naming
func generateRandomID() (string, error) {
	bytes := make([]byte, 4) // Reduced from 8 for shorter names
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
