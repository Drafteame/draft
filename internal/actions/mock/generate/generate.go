package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const (
	baseConfigFile    = ".mockery-base.yml"
	packageConfigFile = ".mockery-packages.yml"
	mockeryConfigFile = ".mockery.yml"
)

type MockGenerate struct {
	input          dtos.MockGenerateInput
	baseConfigPath string
	tempConfigPath string
	packagesFound  int
}

func New(input dtos.MockGenerateInput) *MockGenerate {
	return &MockGenerate{
		input: input,
	}
}

func (mg *MockGenerate) Exec() error {
	if err := mg.preCreate(); err != nil {
		return err
	}

	if err := mg.exec(); err != nil {
		return err
	}

	return mg.postCreate()
}

// preCreate validates that the base config file exists
func (mg *MockGenerate) preCreate() error {
	log.Info("🔍 Validating project root...")

	// Validate we're in the project root by checking for app.toml
	appTomlPath := filepath.Join(mg.input.WorkingDir, "app.toml")
	if _, err := os.Stat(appTomlPath); os.IsNotExist(err) {
		return fmt.Errorf("this command must be run from the project root directory")
	}

	log.Info("🔍 Validating base mockery configuration...")

	mg.baseConfigPath = filepath.Join(mg.input.WorkingDir, baseConfigFile)
	mg.tempConfigPath = filepath.Join(mg.input.WorkingDir, mockeryConfigFile)

	if _, err := os.Stat(mg.baseConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("base config file not found: %s", mg.baseConfigPath)
	}

	// Remove .mockery.yml if it already exists (defensive cleanup)
	if _, err := os.Stat(mg.tempConfigPath); err == nil {
		log.Info("🧹 Removing existing .mockery.yml file...")
		if err := os.Remove(mg.tempConfigPath); err != nil {
			return fmt.Errorf("failed to remove existing config: %w", err)
		}
		log.Success("✓ Existing configuration removed")
	}

	log.Success("✓ Base configuration found")
	return nil
}

// exec finds package configs, concatenates them, and runs mockery
func (mg *MockGenerate) exec() error {
	// Find all package config files using fd and create temp config
	if err := mg.findPackageConfigs(); err != nil {
		return err
	}

	if mg.packagesFound == 0 {
		log.Warn("⚠️  No package configuration files found")
		return nil
	}

	// Run mockery
	if err := mg.runMockery(); err != nil {
		return err
	}

	return nil
}

// postCreate removes the temporary config file
func (mg *MockGenerate) postCreate() error {
	if mg.packagesFound == 0 {
		return nil
	}

	log.Info("🧹 Cleaning up temporary configuration...")

	if err := os.Remove(mg.tempConfigPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove temporary config: %w", err)
	}

	log.Success("✓ Temporary configuration removed")
	return nil
}

// findPackageConfigs uses fd to find all package config files
func (mg *MockGenerate) findPackageConfigs() error {
	log.Info(fmt.Sprintf("🔎 Searching for %s files...", packageConfigFile))

	// Use fd with absolute hidden file flag (-H) to find files
	cmd := exec.Command("fd", "-t", "f", "-H", "-I", packageConfigFile, ".")
	cmd.Dir = mg.input.WorkingDir

	log.Debug(fmt.Sprintf("Running command: %s in directory: %s", cmd.String(), mg.input.WorkingDir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if fd is installed
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("fd command not found. Please install fd: https://github.com/sharkdp/fd")
		}
		// Don't fail if fd returns non-zero but has output (sometimes happens with no matches)
		if len(output) == 0 {
			return fmt.Errorf("failed to search for package configs: %w\nOutput: %s", err, string(output))
		}
	}

	log.Debug(fmt.Sprintf("fd output: %s", string(output)))

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		mg.packagesFound = 0
		return nil
	}

	mg.packagesFound = len(files)
	log.Info(fmt.Sprintf("📦 Found %d package configuration file(s)", mg.packagesFound))

	// Read and concatenate all package configs
	var allPackages strings.Builder
	for i, file := range files {
		if file == "" {
			continue
		}

		fullPath := filepath.Join(mg.input.WorkingDir, file)
		log.Info(fmt.Sprintf("  %d. %s", i+1, file))

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		allPackages.Write(content)
		allPackages.WriteString("\n")
	}

	// Create the temporary .mockery.yml file
	if err := mg.createTempConfig([]byte(allPackages.String())); err != nil {
		return err
	}

	return nil
}

// createTempConfig creates a temporary .mockery.yml with base + packages content
func (mg *MockGenerate) createTempConfig(packagesContent []byte) error {
	log.Info("📝 Creating temporary mockery configuration...")

	// Read base config
	baseContent, err := os.ReadFile(mg.baseConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read base config: %w", err)
	}

	// Concatenate base + packages
	var combined strings.Builder
	combined.Write(baseContent)
	combined.WriteString("\n")
	combined.Write(packagesContent)

	// Write to .mockery.yml
	if err := os.WriteFile(mg.tempConfigPath, []byte(combined.String()), 0644); err != nil {
		return fmt.Errorf("failed to write temporary config: %w", err)
	}

	log.Success("✓ Temporary configuration created")
	return nil
}

// runMockery executes the mockery command
func (mg *MockGenerate) runMockery() error {
	log.Info("🚀 Running mockery...")

	cmd := exec.Command("mockery")
	cmd.Dir = mg.input.WorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Check if mockery is installed
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("mockery command not found. Please install mockery: https://github.com/vektra/mockery")
		}
		return fmt.Errorf("mockery command failed: %w", err)
	}

	log.Success("✓ Mockery execution completed")
	return nil
}
