package project

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/files"
)

func GetServices() (map[string]string, error) {
	type sls struct {
		Service string `yaml:"service"`
	}

	if !files.Exists("services") {
		return nil, errors.New("services directory not found")
	}

	omit := []string{"node_modules", "cmd", "config"}

	services, err := files.Search("services", "serverless.yml", files.WithOmit(omit...))
	if err != nil {
		return nil, err
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(runtime.NumCPU())

	mu := &sync.Mutex{}

	serviceList := make(map[string]string)

	for _, service := range services {
		g.Go(func() error {
			ymlContent, errRead := files.Read(service)
			if errRead != nil {
				return errRead
			}

			var sl sls

			if errUnm := yaml.Unmarshal(ymlContent, &sl); errUnm != nil {
				return errUnm
			}

			basePath := strings.TrimSuffix(service, "/serverless.yml")

			mu.Lock()
			serviceList[sl.Service] = basePath
			mu.Unlock()

			return nil
		})
	}

	if errWait := g.Wait(); errWait != nil {
		return nil, errWait
	}

	return serviceList, nil
}
