package deploy

import "strings"

// EnvConfig holds the deployment configuration for a target environment.
type EnvConfig struct {
	Profile        string // e.g. "draftea-dev", "draftea-prod", "draftea-feature"
	ExtraSLSParams string // non-empty only for feature: --param="stage=feature"
}

// Stage derives the stage name from the profile by stripping the "draftea-" prefix.
func (e EnvConfig) Stage() string {
	return strings.TrimPrefix(e.Profile, "draftea-")
}

var (
	DevEnv = EnvConfig{
		Profile: "draftea-dev",
	}

	ProdEnv = EnvConfig{
		Profile: "draftea-prod",
	}

	FeatureEnv = EnvConfig{
		Profile:        "draftea-feature",
		ExtraSLSParams: `--param="stage=feature"`,
	}
)
