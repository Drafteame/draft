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
	"github.com/Drafteame/draft/internal/pkg/auth"
	nixmetadata "github.com/Drafteame/draft/internal/pkg/nix-metadata"
)

const nixModulesGithubURL = "https://api.github.com/repos/Drafteame/nix-modules/releases/latest"

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
		_, _ = fmt.Printf("Error getting current version: %v\n", err)
		return
	}

	if !nixMetadata.ShouldRunUpdateScript() || nixMetadata.CurrentVersion == nil {
		return
	}

	latestVersion, err := getLatestVersion(nixMetadata)
	if err != nil {
		_, _ = fmt.Printf("Error getting latest version: %v\n", err)
		return
	}

	shouldUpdateLastNegativeResponse := handleUpdate(nixMetadata.CurrentVersion, latestVersion)
	if !shouldUpdateLastNegativeResponse {
		if err = nixMetadata.UpdateLastNegativeUpdatePrompt(); err != nil {
			_, _ = fmt.Printf("Error updating last negative update prompt: %v\n", err)
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
		_, _ = fmt.Printf("Error reading user input: %v\n", err)
		return false
	}
	return input.ShouldUpdateNixModules
}

func performNixUpdate() {
	cmd := exec.Command("hms-update", "--silent=true")

	if err := spinner.New().Type(spinner.Dots).Title("Updating nix-modules...").Action(func() {
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Printf("Error updating nix-modules: %v\n", err)
			return
		}
	}).Run(); err != nil {
		_, _ = fmt.Printf("Error updating nix-modules: %v\n", err)
		return
	}

	println("Nix modules update completed successfully")
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

	ghToken, err := auth.GetGithubToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get github token: %w", err)
	}

	response, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("token %s", ghToken)).
		Get(nixModulesGithubURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if response.IsError() && response.StatusCode() != 401 {
		return nil, fmt.Errorf("github API error: %s", response.Status())
	} else if response.StatusCode() == 401 {
		ghToken, errToken := auth.RefreshGithubToken()
		if errToken != nil {
			return nil, fmt.Errorf("failed to refresh github token: %w", errToken)
		}

		var errRes error
		response, errRes = client.R().
			SetHeader("Authorization", fmt.Sprintf("token %s", ghToken)).
			Get(nixModulesGithubURL)

		if errRes != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", errRes)
		}
	}

	var latestVersion struct {
		Tag string `json:"tag_name"`
	}
	if err := json.Unmarshal(response.Body(), &latestVersion); err != nil {
		return nil, fmt.Errorf("failed to parse github response: %w", err)
	}

	return semver.NewVersion(strings.TrimPrefix(latestVersion.Tag, "v"))
}
