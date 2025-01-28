package nixversion

import (
	"github.com/Masterminds/semver"

	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.UpdateNixModules, currentVersion, latestVersion *semver.Version) error {
	return baseForm(input, currentVersion, latestVersion)
}
