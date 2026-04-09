package deploy

import "strings"

const profilePrefix = "draftea-"

// EnvConfig holds the deployment configuration for a target environment.
type EnvConfig struct {
	Profile        string // e.g. "draftea-dev", "draftea-prod", "draftea-feature"
	ExtraSLSParams string // non-empty only for feature: --param="stage=feature"
}

// Stage derives the stage name from the profile by stripping the profilePrefix.
func (e EnvConfig) Stage() string {
	return strings.TrimPrefix(e.Profile, profilePrefix)
}

var (
	DevEnv = EnvConfig{
		Profile: profilePrefix + "dev",
	}

	ProdEnv = EnvConfig{
		Profile: profilePrefix + "prod",
	}

	FeatureEnv = EnvConfig{
		Profile:        profilePrefix + "feature",
		ExtraSLSParams: `--param="stage=feature"`,
	}
)
