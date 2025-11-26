package mockery

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/log"
)

type mockeryJob struct {
	configFile string
	tmpFile    string
	err        error
	duration   time.Duration
}

type executionStats struct {
	total     int
	succeeded int
	failed    int
	duration  time.Duration
}

type progressUpdate struct {
	current    int
	total      int
	configFile string
	success    bool
	err        error
	duration   time.Duration
}

func (m *Mockery) executeConcurrent(configFiles []string) []mockeryJob {
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

		go func() {
			defer progressWg.Done()
			for update := range progressChan {
				completed++
				status := "✓"
				if !update.success {
					status = "✗"
				}

				shortName := m.shortenConfigPath(update.configFile)
				spin.Title(fmt.Sprintf("[%s] [%2d / %d] %s (%.2fs)",
					status, completed, total, shortName, update.duration.Seconds()))
			}
		}()

		for idx := range m.tmpFiles {
			if m.ctx.Err() != nil {
				if !cancelled {
					log.Warn("Operation cancelled by user, waiting for ongoing tasks to complete...")
					cancelled = true
					execErr = m.ctx.Err()
				}
				goto waitForCompletion
			}

			select {
			case semaphore <- struct{}{}:
			case <-m.ctx.Done():
				if !cancelled {
					log.Warn("Operation cancelled by user, waiting for ongoing tasks to complete...")
					cancelled = true
					execErr = m.ctx.Err()
				}
				goto waitForCompletion
			}

			wg.Add(1)
			go m.executeJob(
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
		wg.Wait()
		close(progressChan)
		close(resultsChan)
		progressWg.Wait()
	}

	if err := spin.Action(action).Run(); err != nil {
		execErr = fmt.Errorf("execution error: %w", err)
	}

	<-doneChan

	if execErr != nil && !errors.Is(execErr, context.Canceled) {
		log.Errorf("Execution encountered errors: %v", execErr)
	}

	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (m *Mockery) executeJob(
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
	defer func() { <-sem }()

	select {
	case <-m.ctx.Done():
		result := mockeryJob{
			configFile: configFile,
			tmpFile:    tmpFile,
			err:        m.ctx.Err(),
			duration:   0,
		}
		resultChan <- result
		return
	default:
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

	progressChan <- progressUpdate{
		current:    idx + 1,
		total:      total,
		configFile: configFile,
		success:    err == nil,
		err:        err,
		duration:   duration,
	}
}

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

func (m *Mockery) shortenConfigPath(configPath string) string {
	path := strings.TrimSuffix(configPath, "/.mockery.pkg.yml")

	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return ".../" + strings.Join(parts[len(parts)-3:], "/")
	}

	return path
}

func (m *Mockery) calculateStats(results []mockeryJob, startTime time.Time) executionStats {
	stats := executionStats{
		total:    len(m.tmpFiles),
		duration: time.Since(startTime),
	}

	for _, result := range results {
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
