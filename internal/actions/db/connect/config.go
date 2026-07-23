package connect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/files"
)

const configFileName = "dbconnect.yml"

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(home, ".draft", configFileName), nil
}

func loadConfig() (ConnConfig, error) {
	path, err := configFilePath()
	if err != nil {
		return ConnConfig{}, err
	}

	if !files.Exists(path) {
		return ConnConfig{}, fmt.Errorf("config file not found: %s", path)
	}

	var cfg ConnConfig
	if err := files.LoadYAML(path, &cfg); err != nil {
		return ConnConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// ResolveConnection builds a ResolvedConnection from (dbType, name) where
// name = "{service}-{env}", e.g. "turbo-dev" or "turbo-prod".
func (cfg ConnConfig) ResolveConnection(dbType, name string) (ResolvedConnection, error) {
	service, env, err := splitServiceEnv(name)
	if err != nil {
		return ResolvedConnection{}, err
	}

	envCfg, ok := cfg.Environments[env]
	if !ok {
		return ResolvedConnection{}, fmt.Errorf("environment %q not found in config (available: %s)",
			env, joinKeys(cfg.Environments))
	}

	typeConfig, ok := cfg.Connections[dbType]
	if !ok {
		return ResolvedConnection{}, fmt.Errorf("unknown DB type %q (supported: postgres, redis, mongo)", dbType)
	}

	svc := findService(typeConfig.Instances, service)
	if svc == nil {
		return ResolvedConnection{}, fmt.Errorf("service %q not found in %s connections", service, dbType)
	}

	host, err := buildHost(dbType, service, env, envCfg.Clusters)
	if err != nil {
		return ResolvedConnection{}, err
	}

	localPort, ok := svc.LocalPorts[env]
	if !ok {
		return ResolvedConnection{}, fmt.Errorf("no local_port defined for env %q in service %q/%s", env, dbType, service)
	}

	remotePort := cfg.Defaults[dbType].RemotePort

	return ResolvedConnection{
		DBType:     dbType,
		Name:       name,
		Service:    service,
		Env:        env,
		Bastion:    envCfg.Bastion,
		Host:       host,
		RemotePort: remotePort,
		LocalPort:  localPort,
	}, nil
}

// splitServiceEnv splits "turbo-dev" → ("turbo", "dev") and "turbo-prod" → ("turbo", "prod").
func splitServiceEnv(name string) (service, env string, err error) {
	if strings.HasSuffix(name, "-dev") {
		return strings.TrimSuffix(name, "-dev"), "dev", nil
	}

	if strings.HasSuffix(name, "-prod") {
		return strings.TrimSuffix(name, "-prod"), "prod", nil
	}

	return "", "", fmt.Errorf("connection name %q must end with -dev or -prod", name)
}

// buildHost constructs the remote host string for the given db type, service, and environment.
//
// PostgreSQL : {service}-{env}.cluster-{clusters.rds}
// Redis      : {service}-{env}.{clusters.cache}
// MongoDB    : draftea-{env}-maincluster.cluster-{clusters.docdb}
func buildHost(dbType, service, env string, clusters ClustersConfig) (string, error) {
	switch dbType {
	case "postgres":
		return fmt.Sprintf("%s-%s.cluster-%s", service, env, clusters.RDS), nil
	case "redis":
		return fmt.Sprintf("%s-%s.%s", service, env, clusters.Cache), nil
	case "mongo":
		return fmt.Sprintf("draftea-%s-maincluster.cluster-%s", env, clusters.DocDB), nil
	default:
		return "", fmt.Errorf("unknown DB type: %s", dbType)
	}
}

func findService(services []ServiceConfig, name string) *ServiceConfig {
	for i := range services {
		if services[i].Name == name {
			return &services[i]
		}
	}

	return nil
}

func joinKeys[V any](m map[string]V) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return strings.Join(keys, ", ")
}
