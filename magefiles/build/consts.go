package build

type EnvType string

const (
	BinPath = ".bin"

	// LocalEnvType is the environment type for local builds.
	LocalEnvType EnvType = "local"
)

var lambdaEnvs = map[string]string{
	"GOOS":        "linux",
	"GOARCH":      "amd64",
	"CGO_ENABLED": "0",
}

var localEnvs = map[string]string{
	"CGO_ENABLED": "0",
}

func (et EnvType) Envs() map[string]string {
	switch et {
	case "lambda":
		return lambdaEnvs
	case "local":
		return localEnvs
	default:
		return nil
	}
}
