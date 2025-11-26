package mockery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Drafteame/draft/internal/pkg/dirs"
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
	staged      bool
	committed   bool
	modified    bool
	tmpFiles    []string
	mu          sync.Mutex
}

func New(ctx context.Context, configFiles []string, jobsNum int, dry bool, staged bool, committed bool, modified bool) *Mockery {
	return &Mockery{
		ctx:         ctx,
		configFiles: configFiles,
		jobsNum:     jobsNum,
		dry:         dry,
		staged:      staged,
		committed:   committed,
		modified:    modified,
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

	if err := m.createTempConfigs(configFiles, baseConfig); err != nil {
		m.cleanup()
		return err
	}

	defer m.cleanup()

	results := m.executeConcurrent(configFiles)

	if m.ctx.Err() != nil {
		stats := m.calculateStats(results, startTime)
		m.displayCancellationSummary(stats, results)
		return fmt.Errorf("operation cancelled: %w", m.ctx.Err())
	}

	stats := m.calculateStats(results, startTime)
	m.displaySummary(stats, results)

	if stats.failed > 0 {
		return fmt.Errorf("mockery execution failed for %d package(s)", stats.failed)
	}

	return nil
}

func (m *Mockery) validate() error {
	if m.jobsNum <= 0 {
		return fmt.Errorf("invalid --jobs-num value: %d (must be greater than 0)", m.jobsNum)
	}

	if m.jobsNum > 100 {
		log.Warnf("Very high concurrency (%d) may cause performance issues", m.jobsNum)
	}

	gitFlagsCount := 0
	if m.staged {
		gitFlagsCount++
	}
	if m.committed {
		gitFlagsCount++
	}
	if m.modified {
		gitFlagsCount++
	}

	if gitFlagsCount > 1 {
		return fmt.Errorf("cannot combine --staged, --committed, and --modified flags (use --modified for staged + committed)")
	}

	if gitFlagsCount > 0 && len(m.configFiles) > 0 {
		return fmt.Errorf("cannot use --staged, --committed, or --modified with explicit config file paths")
	}

	return nil
}

func (m *Mockery) resolveConfigFiles() ([]string, error) {
	log.Info("Resolving config files...")

	var configFiles []string
	var err error

	if m.staged {
		configFiles, err = m.findConfigsFromStaged()
	} else if m.committed {
		configFiles, err = m.findConfigsFromCommitted()
	} else if m.modified {
		configFiles, err = m.findConfigsFromModified()
	} else if len(m.configFiles) > 0 {
		configFiles, err = m.validateProvidedConfigs(m.configFiles)
	} else {
		configFiles, err = m.findPackageConfigs()
	}

	if err != nil {
		return nil, err
	}

	return m.deduplicateFiles(configFiles), nil
}

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

func (m *Mockery) findPackageConfigs() ([]string, error) {
	var configs []string

	err := dirs.Walk(".", func(path string, info os.FileInfo, err error) error {
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

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return configs, nil
}
