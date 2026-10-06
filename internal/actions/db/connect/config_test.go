package connect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const testConfigYAML = `
defaults:
  postgres: { remote_port: 5432 }
  redis:    { remote_port: 6379 }
  mongo:    { remote_port: 27017 }

environments:
  dev:
    bastion: { target: i-dev, profile: dev-profile, region: us-east-2 }
    clusters:
      rds:              aaa.us-east-2.rds.amazonaws.com
      cache:            6erhpv.0001.use2.cache.amazonaws.com
      cache_serverless: 6erhpv.serverless.use2.cache.amazonaws.com
      docdb:            ccc.us-east-2.docdb.amazonaws.com
  prod:
    bastion: { target: i-prod, profile: prod-profile, region: us-east-2 }
    clusters:
      rds:   ddd.us-east-2.rds.amazonaws.com
      cache: kct7ey.0001.use2.cache.amazonaws.com
      docdb: fff.us-east-2.docdb.amazonaws.com

connections:
  postgres:
    instances:
      - name: turbo
        local_ports: { dev: 56000, prod: 56011 }
      - name: bad-pg
        serverless: true
        local_ports: { dev: 56002 }
  redis:
    instances:
      - name: turbo-cache
        local_ports: { dev: 56100, prod: 56101 }
      - name: api-cache
        serverless: true
        local_ports: { dev: 56150, prod: 56151 }
      - name: dev-only-cache
        serverless: true
        local_ports: { dev: 56160 }
  mongo:
    instances:
      - name: main
        local_ports: { dev: 56200, prod: 56201 }
`

func loadTestConfig(t *testing.T) ConnConfig {
	t.Helper()

	var cfg ConnConfig
	require.NoError(t, yaml.Unmarshal([]byte(testConfigYAML), &cfg))

	return cfg
}

func TestResolveConnection(t *testing.T) {
	cfg := loadTestConfig(t)

	tests := []struct {
		name           string
		dbType         string
		connName       string
		wantHost       string
		wantServerless bool
		wantLocalPort  int
		wantErr        string
	}{
		{
			name:          "postgres unchanged",
			dbType:        "postgres",
			connName:      "turbo-dev",
			wantHost:      "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com",
			wantLocalPort: 56000,
		},
		{
			name:          "mongo unchanged",
			dbType:        "mongo",
			connName:      "main-prod",
			wantHost:      "draftea-prod-maincluster.cluster-fff.us-east-2.docdb.amazonaws.com",
			wantLocalPort: 56201,
		},
		{
			name:          "redis regular",
			dbType:        "redis",
			connName:      "turbo-cache-dev",
			wantHost:      "turbo-cache-dev.6erhpv.0001.use2.cache.amazonaws.com",
			wantLocalPort: 56100,
		},
		{
			name:           "redis serverless",
			dbType:         "redis",
			connName:       "api-cache-dev",
			wantHost:       "api-cache-dev-6erhpv.serverless.use2.cache.amazonaws.com",
			wantServerless: true,
			wantLocalPort:  56150,
		},
		{
			name:     "redis serverless without cache_serverless",
			dbType:   "redis",
			connName: "api-cache-prod",
			wantErr:  `redis service "api-cache" is marked serverless but environments.prod.clusters.cache_serverless is not set`,
		},
		{
			name:     "redis serverless without local port reports missing port first",
			dbType:   "redis",
			connName: "dev-only-cache-prod",
			wantErr:  `no local_port defined for env "prod" in service "redis"/dev-only-cache`,
		},
		{
			name:     "serverless on postgres is rejected",
			dbType:   "postgres",
			connName: "bad-pg-dev",
			wantErr:  `postgres service "bad-pg" has serverless set, but serverless is only supported for redis`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cfg.ResolveConnection(tt.dbType, tt.connName)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantHost, got.Host)
			assert.Equal(t, tt.wantServerless, got.Serverless)
			assert.Equal(t, tt.wantLocalPort, got.LocalPort)
		})
	}
}

func TestBuildHost(t *testing.T) {
	clusters := ClustersConfig{
		RDS:             "aaa.rds.amazonaws.com",
		Cache:           "6erhpv.0001.use2.cache.amazonaws.com",
		CacheServerless: "6erhpv.serverless.use2.cache.amazonaws.com",
		DocDB:           "ccc.docdb.amazonaws.com",
	}

	tests := []struct {
		name     string
		dbType   string
		svc      ServiceConfig
		clusters ClustersConfig
		want     string
		wantErr  bool
	}{
		{
			name:     "postgres",
			dbType:   "postgres",
			svc:      ServiceConfig{Name: "turbo"},
			clusters: clusters,
			want:     "turbo-dev.cluster-aaa.rds.amazonaws.com",
		},
		{
			name:     "mongo",
			dbType:   "mongo",
			svc:      ServiceConfig{Name: "main"},
			clusters: clusters,
			want:     "draftea-dev-maincluster.cluster-ccc.docdb.amazonaws.com",
		},
		{
			name:     "redis regular",
			dbType:   "redis",
			svc:      ServiceConfig{Name: "turbo"},
			clusters: clusters,
			want:     "turbo-dev.6erhpv.0001.use2.cache.amazonaws.com",
		},
		{
			name:     "redis serverless",
			dbType:   "redis",
			svc:      ServiceConfig{Name: "api-cache", Serverless: true},
			clusters: clusters,
			want:     "api-cache-dev-6erhpv.serverless.use2.cache.amazonaws.com",
		},
		{
			name:     "redis serverless missing suffix",
			dbType:   "redis",
			svc:      ServiceConfig{Name: "api-cache", Serverless: true},
			clusters: ClustersConfig{Cache: clusters.Cache},
			wantErr:  true,
		},
		{
			name:     "mongo serverless rejected",
			dbType:   "mongo",
			svc:      ServiceConfig{Name: "main", Serverless: true},
			clusters: clusters,
			wantErr:  true,
		},
		{
			name:     "unknown type",
			dbType:   "mysql",
			svc:      ServiceConfig{Name: "x"},
			clusters: clusters,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHost(tt.dbType, &tt.svc, "dev", tt.clusters)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
