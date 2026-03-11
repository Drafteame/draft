package deploy

// EnvConfig holds the deployment configuration for a target environment.
type EnvConfig struct {
	Stage          string
	AWSAccount     string
	AWSProfile     string
	ExtraSLSParams string // non-empty only for feature: --params="stage=feature"
}

var (
	DevEnv = EnvConfig{
		Stage:      "dev",
		AWSAccount: "776658659836",
		AWSProfile: "draftea-dev",
	}

	ProdEnv = EnvConfig{
		Stage:      "prod",
		AWSAccount: "632258128187",
		AWSProfile: "draftea-prod",
	}

	FeatureEnv = EnvConfig{
		Stage:          "feature",
		AWSAccount:     "636385746594",
		AWSProfile:     "draftea-feature",
		ExtraSLSParams: `--params="stage=feature"`,
	}
)
