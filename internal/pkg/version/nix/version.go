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
	"github.com/Drafteame/draft/internal/pkg/log"
	nixmetadata "github.com/Drafteame/draft/internal/pkg/nix-metadata"
)

const nixModulesVersion = "https://draftea-cdn-prod.s3.us-east-2.amazonaws.com/xb/nix/version.json"

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
		log.Debugf("Error getting current nix modules version: %v", err)
		return
	}

	if !nixMetadata.ShouldRunUpdateScript() || nixMetadata.CurrentVersion == nil {
		return
	}

	latestVersion, err := getLatestVersion(nixMetadata)
	if err != nil {
		log.Debugf("Error getting latest nix modules version: %v", err)
		return
	}

	shouldUpdateLastNegativeResponse := handleUpdate(nixMetadata.CurrentVersion, latestVersion)
	if !shouldUpdateLastNegativeResponse {
		if err = nixMetadata.UpdateLastNegativeUpdatePrompt(); err != nil {
			log.Debugf("Error updating last negative update prompt: %v", err)
		}
	}
}

func handleUpdate(currentVersion, latestVersion *semver.Version) (confirmedUpdate bool) {
	switch updateType := getUpdateType(currentVersion, latestVersion); updateType {
	case MajorUpdate, MinorUpdate:
		userResponseToUpdatePrompt := confirmUpdate(currentVersion, latestVersion)
		if userResponseToUpdatePrompt {
			performNixUpdate()
		}
		return userResponseToUpdatePrompt
	case PatchUpdate:
		go performNixUpdate()
		return true
	default:
		return true
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
		log.Debugf("Error reading user for nix modules input: %v", err)
		return false
	}
	return input.ShouldUpdateNixModules
}

func performNixUpdate() {
	cmd := exec.Command("hms-update", "--silent=true")

	err := spinner.New().Type(spinner.Dots).
		Title("Updating nix-modules...").
		Action(func() {
			if err := cmd.Run(); err != nil {
				log.Debugf("Error updating nix-modules: %v", err)
				return
			}
		}).
		Run()

	if err != nil {
		log.Debugf("Error updating nix-modules: %v", err)
		return
	}

	log.Info("Nix modules update completed successfully")
}

func getLatestVersion(nixMetadata nixmetadata.NixMetadata) (*semver.Version, error) {
	if nixMetadata.ShouldFetchNixVersion() {
		version, err := fetchLatestVersion()
		if err != nil {
			return nil, err
		}

		nixMetadata.LastNixVersion = version.String()

		if err := nixMetadata.UpdateLastNixVersionCheck(); err != nil {
			// Best effort, if we can't update the nix-metadata, we still want to update the version
			return version, err
		}
	}
	return semver.MustParse(nixMetadata.LastNixVersion), nil
}

func fetchLatestVersion() (*semver.Version, error) {
	client := resty.New().
		SetTimeout(2 * time.Second)

	response, err := client.R().
		Get(nixModulesVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting nix version: %s", response.Status())
	}

	var latestVersion struct {
		NixVersion string `json:"nixVersion"`
	}

	if err := json.Unmarshal(response.Body(), &latestVersion); err != nil {
		return nil, fmt.Errorf("failed to parse github response: %w", err)
	}

	return semver.NewVersion(strings.TrimPrefix(latestVersion.NixVersion, "v"))
}
