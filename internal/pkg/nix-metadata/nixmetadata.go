package nixmetadata

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/semver"

	"github.com/Drafteame/draft/internal/pkg/files"
)

const nixMetadataPath = "$HOME/.config/.nix-metadata.json"

type NixMetadata struct {
	SysUserName              string          `json:"sysUserName"`
	SysArch                  string          `json:"sysArch"`
	SysOs                    string          `json:"sysOs"`
	GitUserName              string          `json:"gitUserName"`
	GitUser                  string          `json:"gitUser"`
	GitEmail                 string          `json:"gitEmail"`
	Shell                    string          `json:"shell"`
	ShellRcFile              string          `json:"shellRcFile"`
	Version                  string          `json:"version"`
	CurrentVersion           *semver.Version `json:"current_semver_version"`
	LastUpdated              time.Time       `json:"last_updated"`
	LastNixVersionCheckAt    time.Time       `json:"last_nix_version_check"`
	LastNixVersion           string          `json:"last_nix_version"`
	LastNegativeUpdatePrompt time.Time       `json:"last_negative_update_prompt"`
}

func (n *NixMetadata) ShouldRunUpdateScript() bool {
	return time.Since(n.LastUpdated) > 24*time.Hour && time.Since(n.LastNegativeUpdatePrompt) > 2*time.Hour
}

func (n *NixMetadata) ParseCurrentVersion() {
	if n.Version == "" {
		println("Current version not found in nix metadata")
		return
	}

	parsedVersion, err := semver.NewVersion(n.Version)
	if err != nil {
		_, _ = fmt.Printf("Error parsing current version: %v\n", err)
		return
	}

	n.CurrentVersion = parsedVersion
}

func (n *NixMetadata) ShouldFetchNixVersion() bool {
	return time.Since(n.LastNixVersionCheckAt) > 3*time.Hour || n.LastNixVersionCheckAt.Equal(time.Time{}) || n.LastNixVersion == ""
}

func (n *NixMetadata) UpdateLastNixVersionCheck() error {
	n.LastNixVersionCheckAt = time.Now()

	metadata, err := json.MarshalIndent(n, "", "	")
	if err != nil {
		_, _ = fmt.Printf("Error updating nix-metadata: %v\n", err)
		return err
	}

	if err := files.Create(nixMetadataPath, metadata); err != nil {
		_, _ = fmt.Printf("Error updating nix-metadata: %v\n", err)
		return err
	}

	return nil
}

func (n *NixMetadata) UpdateLastNegativeUpdatePrompt() error {
	n.LastNegativeUpdatePrompt = time.Now()

	metadata, err := json.MarshalIndent(n, "", "	")
	if err != nil {
		_, _ = fmt.Printf("Error updating nix-metadata: %v\n", err)
		return err
	}

	if err := files.Create(nixMetadataPath, metadata); err != nil {
		_, _ = fmt.Printf("Error updating nix-metadata: %v\n", err)
		return err
	}

	return nil
}

func Get() (NixMetadata, error) {
	nixMetadataContent, err := readNixMetadataFile()
	if err != nil {
		return NixMetadata{}, err
	}

	nixMetadataContent.ParseCurrentVersion()

	return nixMetadataContent, nil
}

func readNixMetadataFile() (NixMetadata, error) {
	if !files.Exists(nixMetadataPath) {
		return NixMetadata{}, fmt.Errorf("nix metadata file not found, expected at %s", nixMetadataPath)
	}

	fileContent, err := files.Read(nixMetadataPath)
	if err != nil {
		return NixMetadata{}, err
	}

	var nixMetadata NixMetadata
	if err := json.Unmarshal(fileContent, &nixMetadata); err != nil {
		return NixMetadata{}, err
	}

	return nixMetadata, nil
}
