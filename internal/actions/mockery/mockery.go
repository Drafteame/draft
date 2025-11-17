package mockery

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"gopkg.in/yaml.v3"

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
		configFiles []string
		jobsNum     int
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

func New(configFiles []string, jobsNum int) *Mockery {
	return &Mockery{
		configFiles: configFiles,
		jobsNum:     jobsNum,
		tmpFiles:    make([]string, 0),
	}
}

func (m *Mockery) Exec() error {
	// Validate inputs
	if err := m.validate(); err != nil {
		return err
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
	if err := m.createTempConfigs(configFiles, baseConfig); err != nil {
		m.cleanup() // Clean up any temp files created before the error
		return err
	}

	// Ensure cleanup on exit
	defer m.cleanup()

	// Execute mockery concurrently with progress spinner
	results := m.executeConcurrentWithProgress(configFiles)

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

	return nil
}

// resolveConfigFiles resolves, validates, and deduplicates config files
func (m *Mockery) resolveConfigFiles() ([]string, error) {
	var configFiles []string

	// If files provided, validate them
	if len(m.configFiles) > 0 {
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
		// Check if file exists
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
		// Normalize path
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
	var searchErr error

	spin := spinner.New().Title("Searching for package config files...")
	action := func() {
		searchErr = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip hidden directories and vendor (but not the root "." directory)
			if info.IsDir() {
				// Don't skip the root directory
				if path == "." {
					return nil
				}

				if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" || info.Name() == "node_modules" {
					return filepath.SkipDir
				}
				return nil
			}

			if strings.HasSuffix(path, pkgConfigSuffix) {
				configs = append(configs, path)
			}

			return nil
		})
	}

	if err := spin.Action(action).Run(); err != nil {
		return nil, fmt.Errorf("failed to search for configs: %w", err)
	}

	if searchErr != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", searchErr)
	}

	return configs, nil
}

// loadBaseConfig loads the base configuration file
func (m *Mockery) loadBaseConfig() (map[string]any, error) {
	if !files.Exists(baseConfigFile) {
		return make(map[string]any), nil
	}

	data, err := os.ReadFile(baseConfigFile)
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
	for i, configFile := range configFiles {
		// Load package config
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", configFile, err)
		}

		var pkgConfig map[string]any
		if err := yaml.Unmarshal(data, &pkgConfig); err != nil {
			return fmt.Errorf("failed to parse %s: %w (check YAML syntax)", configFile, err)
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

		// Write temp file
		if err := os.WriteFile(tmpFile, mergedData, 0644); err != nil {
			return fmt.Errorf("failed to write temp config %s: %w", tmpFile, err)
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

	// Create descriptive temp file name: .mockery.tmp.[pkg-name].[idx].[rand].yml
	tmpFile := fmt.Sprintf("%s%s.%d.%s%s", tmpConfigPrefix, pkgName, index, randID, tmpConfigSuffix)

	return tmpFile, nil
}

// executeConcurrentWithProgress executes mockery commands concurrently with progress spinner
func (m *Mockery) executeConcurrentWithProgress(configFiles []string) []mockeryJob {
	var (
		wg           sync.WaitGroup
		mu           sync.Mutex
		results      []mockeryJob
		semaphore    = make(chan struct{}, m.jobsNum)
		progressChan = make(chan progressUpdate, len(configFiles))
		completed    = 0
		execErr      error
	)

	total := len(configFiles)

	spin := spinner.New().Title(fmt.Sprintf("[ 0 / %d ] Preparing...", total))

	action := func() {
		// Start goroutines to process configs
		for i := range m.tmpFiles {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				configFile := configFiles[idx]
				tmpFile := m.tmpFiles[idx]

				startTime := time.Now()
				err := m.runMockery(tmpFile, configFile)
				duration := time.Since(startTime)

				result := mockeryJob{
					configFile: configFile,
					tmpFile:    tmpFile,
					err:        err,
					duration:   duration,
				}

				mu.Lock()
				results = append(results, result)
				mu.Unlock()

				// Send progress update
				progressChan <- progressUpdate{
					current:    idx + 1,
					total:      total,
					configFile: configFile,
					success:    err == nil,
					err:        err,
					duration:   duration,
				}
			}(i)
		}

		// Wait for all goroutines and close progress channel
		go func() {
			wg.Wait()
			close(progressChan)
		}()

		// Process progress updates
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
	}

	if err := spin.Action(action).Run(); err != nil {
		execErr = fmt.Errorf("execution error: %w", err)
	}

	if execErr != nil {
		log.Errorf("Execution encountered errors: %v", execErr)
	}

	return results
}

// shortenConfigPath extracts a meaningful short name from config file path
func (m *Mockery) shortenConfigPath(configPath string) string {
	// Remove .mockery.pkg.yml suffix
	path := strings.TrimSuffix(configPath, "/.mockery.pkg.yml")

	// If path is too long, show only last 2-3 segments
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return ".../" + strings.Join(parts[len(parts)-3:], "/")
	}

	return path
}

// runMockery executes mockery with the given config file
func (m *Mockery) runMockery(configPath, originalPath string) error {
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
		total:    len(results),
		duration: time.Since(startTime),
	}

	for _, result := range results {
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
		log.Successf("✓ All %d package(s) completed successfully (%.2fs)", stats.total, stats.duration.Seconds())
	}
}

// generateRandomID generates a random hex string for temp file naming
func generateRandomID() (string, error) {
	bytes := make([]byte, 4) // Reduced from 8 for shorter names
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
