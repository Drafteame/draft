package nixversion

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/charmbracelet/huh/spinner"
	"github.com/go-resty/resty/v2"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/files"
	nixmetadata "github.com/Drafteame/draft/internal/pkg/nix-metadata"
)

const nixModulesGithubURL = "https://api.github.com/repos/Drafteame/nix-modules/releases/latest"

const ghToken = ""

const cacheFile = "$HOME/.cache/draft/latest_nix_version.json"

const cacheDuration = 6 * time.Hour

type VersionCache struct {
	LatestVersion string    `json:"latest_nix_version"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type UpdateType int

const (
	NoUpdate UpdateType = iota
	PatchUpdate
	MinorUpdate
	MajorUpdate
)

func CheckNixModulesVersion() {
	nixMetadata, err := nixmetadata.Get()
	if err != nil {
		fmt.Printf("Error getting current version: %v\n", err)
		return
	}

	if !nixMetadata.ShouldRunUpdateScript() || nixMetadata.CurrentVersion == nil {
		return
	}

	latestVersion, err := getLatestVersion(nixMetadata)
	if err != nil {
		fmt.Printf("Error getting latest version: %v\n", err)
		return
	}

	handleUpdate(nixMetadata.CurrentVersion, latestVersion)
}

func handleUpdate(currentVersion, latestVersion *semver.Version) {
	switch updateType := getUpdateType(currentVersion, latestVersion); updateType {
	case MajorUpdate:
		return
	case MinorUpdate:
		if confirmUpdate(currentVersion, latestVersion) {
			performNixUpdate()
		}
	case PatchUpdate:
		go performNixUpdate()
	default:
		return
	}
}

func getUpdateType(current, latest *semver.Version) UpdateType {
	// If the major version is greater, we skip the update to avoid breaking changes
	if latest.Major() > current.Major() {
		return MajorUpdate
	}
	// If the minor version is greater, we prompt the user to update
	if latest.Minor() > current.Minor() {
		return MinorUpdate
	}
	// If the patch version is greater, we force the update
	if latest.Patch() > current.Patch() {
		return PatchUpdate
	}
	return NoUpdate
}

func confirmUpdate(current, latest *semver.Version) bool {
	input := dtos.UpdateNixModules{}
	err := forms.UpdateNixModules(&input, current, latest)
	if err != nil {
		fmt.Printf("Error reading user input: %v\n", err)
		return false
	}
	return input.ShouldUpdateNixModules
}

func performNixUpdate() {
	cmd := exec.Command("hms-update", "--silent=true")

	if err := spinner.New().Type(spinner.Dots).Title("Updating nix-modules...").Action(func() {
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error updating nix-modules: %v\n", err)
			return
		}
	}).Run(); err != nil {
		fmt.Printf("Error updating nix-modules: %v\n", err)
		return
	}

	fmt.Println("Nix modules update completed successfully")
}

func getLatestVersion(nixMetadata nixmetadata.NixMetadata) (*semver.Version, error) {
	if version, err := getVersionFromCache(); err == nil {
		return version, nil
	}

	version, err := fetchLatestVersion()
	if err != nil {
		return nil, err
	}

	if err := nixMetadata.UpdateLastNixVersionCheck(version.String()); err != nil {
		// Best effort, if we can't update the nix-metadata, we still want to update the version
		return version, err
	}

	return version, nil
}

func getVersionFromCache() (*semver.Version, error) {
	if !files.Exists(cacheFile) {
		return nil, fmt.Errorf("cache file not found")
	}

	content, err := files.Read(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache VersionCache
	if err := json.Unmarshal(content, &cache); err != nil {
		return nil, err
	}

	if time.Since(cache.FetchedAt) > cacheDuration {
		return nil, fmt.Errorf("cache expired")
	}

	return semver.NewVersion(cache.LatestVersion)
}

func fetchLatestVersion() (*semver.Version, error) {
	client := resty.New().
		SetTimeout(2 * time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(100 * time.Millisecond)

	response, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("token %s", ghToken)).
		Get(nixModulesGithubURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("github API error: %s", response.Status())
	}

	var latestVersion struct {
		Tag string `json:"tag_name"`
	}
	if err := json.Unmarshal(response.Body(), &latestVersion); err != nil {
		return nil, fmt.Errorf("failed to parse github response: %w", err)
	}

	return semver.NewVersion(strings.TrimPrefix(latestVersion.Tag, "v"))
}
