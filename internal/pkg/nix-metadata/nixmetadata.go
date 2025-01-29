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
	SysUserName           string `json:"sysUserName"`
	SysArch               string `json:"sysArch"`
	SysOs                 string `json:"sysOs"`
	GitUserName           string `json:"gitUserName"`
	GitUser               string `json:"gitUser"`
	GitEmail              string `json:"gitEmail"`
	Shell                 string `json:"shell"`
	ShellRcFile           string `json:"shellRcFile"`
	Version               string `json:"version"`
	CurrentVersion        *semver.Version
	LastUpdated           time.Time `json:"last_updated"`
	LastNixVersionCheckAt time.Time `json:"last_nix_version_check"`
	LastNixVersion        string    `json:"last_nix_version"`
}

func (n *NixMetadata) ShouldRunUpdateScript() bool {
	return time.Since(n.LastUpdated) > 24*time.Hour
}

func (n *NixMetadata) ParseCurrentVersion() {
	if n.Version == "" {
		fmt.Println("Current version not found in nix metadata")
		return
	}

	parsedVersion, err := semver.NewVersion(n.Version)
	if err != nil {
		fmt.Printf("Error parsing current version: %v\n", err)
		return
	}

	n.CurrentVersion = parsedVersion
}

func (n *NixMetadata) UpdateLastNixVersionCheck(latestVersion string) error {
	n.LastNixVersionCheckAt = time.Now()
	n.LastNixVersion = latestVersion

	metadata, err := json.Marshal(n)
	if err != nil {
		fmt.Printf("Error updating nix-metadata: %v\n", err)
		return err
	}

	if err := files.Create(nixMetadataPath, metadata); err != nil {
		fmt.Printf("Error updating nix-metadata: %v\n", err)
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
