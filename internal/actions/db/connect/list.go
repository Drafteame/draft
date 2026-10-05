package connect

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/log"
)

// List prints all connections that can be derived from the config
// by combining every service with every environment.
func List() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	type row struct {
		env        string
		dbType     string
		name       string
		host       string
		serverless bool
		remotePort int
		localPort  int
	}

	var rows []row

	// Deterministic env ordering: dev before prod, then alphabetical for others.
	envOrder := sortedEnvs(cfg.Environments)

	for dbType, typeConfig := range cfg.Connections {
		remotePort := cfg.Defaults[dbType].RemotePort

		for _, svc := range typeConfig.Instances {
			for _, env := range envOrder {
				envCfg := cfg.Environments[env]

				localPort, ok := svc.LocalPorts[env]
				if !ok {
					continue // env not defined for this service, skip
				}

				host, err := buildHost(dbType, &svc, env, envCfg.Clusters)
				if err != nil {
					log.Warnf("skipping %s/%s-%s: %s", dbType, svc.Name, env, err)
					continue
				}

				rows = append(rows, row{
					env:        env,
					dbType:     dbType,
					name:       svc.Name + "-" + env,
					host:       host,
					serverless: dbType == "redis" && svc.Serverless.For(env),
					remotePort: remotePort,
					localPort:  localPort,
				})
			}
		}
	}

	if len(rows) == 0 {
		fmt.Println("No connections configured.")
		return nil
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].dbType != rows[j].dbType {
			return rows[i].dbType < rows[j].dbType
		}

		if rows[i].name != rows[j].name {
			return rows[i].name < rows[j].name
		}

		return rows[i].env < rows[j].env
	})

	wEnv := len("ENV")
	wType := len("TYPE")
	wName := len("NAME")
	wHost := len("HOST")

	for _, r := range rows {
		if len(r.env) > wEnv {
			wEnv = len(r.env)
		}

		if len(r.dbType) > wType {
			wType = len(r.dbType)
		}

		if len(r.name) > wName {
			wName = len(r.name)
		}

		if len(r.host) > wHost {
			wHost = len(r.host)
		}
	}

	wEnv += 2
	wType += 2
	wName += 2
	wHost += 2

	const wServerless = len("SERVERLESS") + 2

	fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-13s %s\n",
		wEnv, "ENV", wType, "TYPE", wName, "NAME", wServerless, "SERVERLESS", wHost, "HOST", "REMOTE PORT", "LOCAL PORT")
	fmt.Println(strings.Repeat("-", wEnv+wType+wName+wServerless+wHost+27))

	for _, r := range rows {
		serverless := "-"
		if r.serverless {
			serverless = "yes"
		}

		fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-13d %d\n",
			wEnv, r.env, wType, r.dbType, wName, r.name, wServerless, serverless, wHost, r.host, r.remotePort, r.localPort)
	}

	return nil
}

// sortedEnvs returns environment keys with "dev" first, "prod" second, then alphabetical.
func sortedEnvs(envs map[string]EnvConfig) []string {
	keys := make([]string, 0, len(envs))
	for k := range envs {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		order := map[string]int{"dev": 0, "prod": 1}
		oi, iOk := order[keys[i]]
		oj, jOk := order[keys[j]]

		if iOk && jOk {
			return oi < oj
		}

		if iOk {
			return true
		}

		if jOk {
			return false
		}

		return keys[i] < keys[j]
	})

	return keys
}
