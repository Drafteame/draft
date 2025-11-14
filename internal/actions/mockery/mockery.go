package mockery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var (
	ErrConfigFileNotFound = errors.New(".mockery.yml configuration file not found")
)

type (
	Mockery struct {
		servicePkg  string
		servicePath string
		runAll      bool
	}
)

func New(servicePkg, servicePath string) *Mockery {
	runAll := servicePkg == "" && servicePath == ""
	return &Mockery{
		servicePkg:  servicePkg,
		servicePath: servicePath,
		runAll:      runAll,
	}
}

func (m *Mockery) Exec() error {
	if m.runAll {
		return m.execAll()
	}
	return m.exec()
}

func (m *Mockery) exec() error {
	spin := spinner.New().Title("Resolving configuration file path")

	var configPath string
	var err error
	action := func() {
		configPath, err = m.resolveConfigPath()
	}
	spinErr := spin.Action(action).Run()
	if spinErr != nil || err != nil {
		return errors.Join(spinErr, err)
	}

	if configPath == "" {
		log.Infof("Service %s skipped", m.servicePkg)
		return nil
	}

	return m.runMockery(configPath)
}

func (m *Mockery) resolveConfigPath() (string, error) {
	if m.servicePkg != "" && m.servicePath != "" {
		return "", fmt.Errorf("only one of --service (-s) or --path (-p) flags can be provided")
	}

	path := filepath.Join("services", m.servicePkg)
	if m.servicePath != "" {
		path = m.servicePath
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path not found: %s", path)
		}
		return "", fmt.Errorf("failed searching path directory: %w", err)
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", path)
	}

	configPath := filepath.Join(path, ".mockery.yml")
	if !files.Exists(configPath) {
		return "", ErrConfigFileNotFound
	}

	return configPath, nil
}

func (m *Mockery) runMockery(configPath string) error {
	log.Info("Running mockery...")

	command := fmt.Sprintf("mockery --config %s", configPath)
	if _, err := exec.Command(command, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr)); err != nil {
		return fmt.Errorf("mockery command failed: %w", err)
	}

	return nil
}

func (m *Mockery) execAll() error {
	services, err := searchServices()
	if err != nil {
		return err
	}

	var errs []error
	for i, service := range services {
		log.Infof("[%d/%d] Processing service: %s", i+1, len(services), service)

		svc := New(service, "")
		if err = svc.exec(); err != nil {
			switch {
			case errors.Is(err, ErrConfigFileNotFound):
				log.Warnf("Service %s skipped: %v", service, err)
			default:
				log.Errorf("Service %s failed: %v", service, err)
				errs = append(errs, err)
			}
			continue
		}

		log.Successf("Service %s completed", service)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func searchServices() ([]string, error) {
	servicesDir := "services"
	stat, err := os.Stat(servicesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("services directory not found: %s", servicesDir)
		}
		return nil, fmt.Errorf("failed searching services directory: %w", err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("services path is not a directory: %s", servicesDir)
	}

	servicesEntries, err := os.ReadDir(servicesDir)
	if err != nil {
		return nil, fmt.Errorf("failed reading services directory: %w", err)
	}

	services := make([]string, 0, len(servicesEntries))
	for _, entry := range servicesEntries {
		if entry.IsDir() {
			services = append(services, entry.Name())
		}
	}

	if len(services) == 0 {
		log.Warn("No services found")
		return nil, nil
	}

	log.Infof("Found %d services", len(services))
	return services, nil
}
