package nixversion

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func baseForm(input *dtos.UpdateNixModules, currentVersion, latestVersion *semver.Version) error {
	err := inputs.Confirm(fmt.Sprintf("There is a new version of nix-modules (%s -> %s)", currentVersion.String(), latestVersion.String()),
		inputs.WithDescription[bool]("Do you want to update?"),
		inputs.WithValue(&input.ShouldUpdateNixModules),
	)

	if err != nil {
		return err
	}

	return nil
}
