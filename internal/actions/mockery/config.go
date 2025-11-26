package mockery

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

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

func (m *Mockery) deepMerge(base, override map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		if vMap, ok := v.(map[string]any); ok {
			if baseMap, ok := result[k].(map[string]any); ok {
				result[k] = m.deepMerge(baseMap, vMap)
				continue
			}
		}
		result[k] = v
	}

	return result
}

func (m *Mockery) createTempConfigs(configFiles []string, baseConfig map[string]any) error {
	log.Info("Creating temporary config files...")

	for i, configFile := range configFiles {
		var pkgConfig map[string]any

		if err := files.LoadYAML(configFile, &pkgConfig); err != nil {
			return fmt.Errorf("failed to load %s: %w", configFile, err)
		}

		merged := m.deepMerge(baseConfig, pkgConfig)

		tmpFile, err := m.generateTempFileName(configFile, i)
		if err != nil {
			return fmt.Errorf("failed to generate temp file name: %w", err)
		}

		mergedData, err := yaml.Marshal(merged)
		if err != nil {
			return fmt.Errorf("failed to marshal merged config for %s: %w", configFile, err)
		}

		if err := files.Create(tmpFile, mergedData); err != nil {
			return fmt.Errorf("failed to write temp config %s: %w", tmpFile, err)
		}

		m.mu.Lock()
		m.tmpFiles = append(m.tmpFiles, tmpFile)
		m.mu.Unlock()
	}

	return nil
}

func (m *Mockery) generateTempFileName(configFile string, index int) (string, error) {
	randID, err := generateRandomID()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(configFile)
	pkgName := filepath.Base(dir)
	if pkgName == "." {
		pkgName = "root"
	}

	tmpFile := fmt.Sprintf("%s%s.%d.%s%s", tmpConfigPrefix, pkgName, index, randID, tmpConfigSuffix)

	return tmpFile, nil
}

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

func generateRandomID() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
