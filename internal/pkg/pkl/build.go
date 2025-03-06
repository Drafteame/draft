package pkl

import (
	"fmt"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
)

type buildOptions struct {
	envs map[string]string
}

type BuildOption func(*buildOptions)

func WithEnvs(envs map[string]string) BuildOption {
	return func(o *buildOptions) {
		o.envs = envs
	}
}

func Build(path, out, format string, opts ...BuildOption) (string, error) {
	options := buildOptions{}

	for _, o := range opts {
		o(&options)
	}

	if !files.Exists(path) {
		return "", nil
	}

	cmdf := "pkl eval -f %s %s > %s"
	cmd := fmt.Sprintf(cmdf, format, path, out)

	out, err := exec.Command(cmd)
	if err != nil {
		return "", err
	}

	return out, nil
}
