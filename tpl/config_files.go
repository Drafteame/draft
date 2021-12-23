package tpl

var (
	EngineYaml = `name: {{Name}}
namespace: {{Namespace}}

# Base environment variables
#
APP_NAME: {{Name}}
CORS: true
STAGE: local
PORT: 8080

# This should match with the custom domain configuration path
#
BASE_PATH: {{Name}}

# Debug mode
#
# DEBUG: true

# MongoDB configuration
#
# MONGO_URI: mongodb://localhost:27017
# MONGO_AUTH_SOURCE: '' # e.g. 'admin'
# MONGO_AUTH_MECHANISM: '' # e.g. 'SCRAM-SHA-256'
# MONGO_REPLICA_SET: '' # e.g. 'rs0'
# MONGO_SSL_VERIFY_CERTIFICATE: false
# MONGO_READ_PREFERENCES: '' # e.g. 'secondaryPreferred'
# MONGO_MIN_POOL_SIZE: 1
# MONGO_MAX_POOL_SIZE: 10
# MONGO_CONNECT_TIMEOUT: 5
# MONGO_QUERY_TIMEOUT: 30
# MONGO_USERNAME: ''
# MONGO_PASSWORD: ''
# MONGO_CLUSTER_ENDPOINT: ''
# MONGO_CERT_PATH: ''
# MONGO_DB_NAME: ''

# QLDB driver configuration
#
# QLDB_LEDGER_NAME: {{Name}}
# QLDB_REGION: us-west-2`

	ConcurrencyYaml = `.default: &defaultConfig
  enabled: true
  prewarm: true
  tracing: true

dev:
  warmup:
    default:
      <<: *defaultConfig
      concurrency: 1

prod:
  warmup:
    default:
      <<: *defaultConfig
      concurrency: 20`

	IamYaml = `.iam: &service_roles
  deploymentRole: ${ssm:/service/deployment-role}

  # Service level permissions
  #
  # role:
  #   statements:
  #     - Effect: Allow
  #       Action: sts:AssumeRole
  #       Resource: "*"

dev:
  roles:
    <<: *service_roles

prod:
  roles:
    <<: *service_roles`

	VpcYaml = `dev:
  vpc:
    securityGroupIds:
      - sg-0aae9f61bf9da1239
    subnetIds:
      - subnet-028b58cb52f3ccd40
      - subnet-090ab0a33f3c981b7
      - subnet-0160e052eab417d32

prod:
  vpc:
    securityGroupIds:
      - sg-0d7e919c8a276d635
    subnetIds:
      - subnet-098d1660c20c1395d
      - subnet-05162e0b22ed3dd52
      - subnet-0a25518bed8ed1701`

	CorsYaml = `.cors: &default
  origin: '*'
  allowCredentials: true
  headers:
    - Origin
    - Content-Type
    - Accept
    - Authorization
    - Startlower
    - Text
    - X-Amz-Date
    - X-Api-Key
    - X-Amz-Security-Token
    - X-Amz-User-Agent
    - Access-Control-Allow-Headers
    - Access-Control-Allow-Origin
    - Access-Control-Allow-Methods
    - Access-Control-Allow-Credentials
    - X-User-Id
    - X-Trace-Id

dev:
  cors:
    <<: *default

prod:
  cors:
    <<: *default`

	DomainsYaml = `.config: &default
  basePath: '{{Name}}'
  stage: ${self:provider.stage}
  certificateName: '*.draftea.com'
  securityPolicy: tls_1_2
  createRoute53Record: true
  endpointType: 'edge'

dev:
  customDomain:
    <<: *default
    domainName: 'dev-api.draftea.com'

prod:
  customDomain:
    <<: *default
    domainName: 'api.draftea.com'`

	EnvironmentYaml = `.envs: &default
  APP_NAME: {{Name}}
  CORS: true
  STAGE: ${self:provider.stage}
  PORT: 8080
  BASE_PATH: {{Name}}

  DEBUG: true

  # MONGO_URI: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_URI}
  # MONGO_AUTH_SOURCE: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_AUTH_SOURCE}
  # MONGO_AUTH_MECHANISM: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_AUTH_MECHANISM}
  # MONGO_REPLICA_SET: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_REPLICA_SET}
  # MONGO_READ_PREFERENCE: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_READ_PREFERENCE}
  # MONGO_SSL_VERIFY_CERTIFICATE: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_SSL_VERIFY_CERTIFICATE}
  # MONGO_MIN_POOL_SIZE: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_MIN_POOL_SIZE}
  # MONGO_MAX_POOL_SIZE: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_MAX_POOL_SIZE}
  # MONGO_CONNECT_TIMEOUT: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_CONNECT_TIMEOUT}
  # MONGO_QUERY_TIMEOUT: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_QUERY_TIMEOUT}
  # MONGO_CERT_PATH: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_CERT_PATH}
  # MONGO_USERNAME: ${ssm:/service/mongo/${self:provider.stage}/MONGO_USERNAME}
  # MONGO_PASSWORD: ${ssm:/service/mongo/${self:provider.stage}/MONGO_PASSWORD}
  # MONGO_CLUSTER_ENDPOINT: ${ssm:/service/mongo/${self:provider.stage}/MONGO_CLUSTER_ENDPOINT}
  # MONGO_DB_NAME: ${ssm:/service/{{Name}}/${self:provider.stage}/MONGO_DB_NAME}

  # QLDB_LEDGER_NAME: ${ssm:/service/{{Name}}/${self:provider.stage}/QLDB_LEDGER_NAME}
  # QLDB_REGION: ${ssm:/service/{{Name}}/${self:provider.stage}/QLDB_REGION}


dev:
  environment:
    <<: *default

prod:
  environment:
    <<: *default`

	ServerlessYaml = `service: {{Name}}


frameworkVersion: '2'
variablesResolutionMode: 20210326


custom:
  customDomain: ${file(config/domains.yml):${self:provider.stage}.customDomain}
  warmup: ${file(config/concurrency.yml):${self:provider.stage}.warmup}


provider:
  name: aws
  runtime: go1.x
  timeout: 20
  stage: ${opt:stage, "dev"}
  region: ${opt:region, "us-east-2"}
  lambdaHashingVersion: '20201221'
  environment: ${file(config/environment.yml):${self:provider.stage}.environment}
  iam: ${file(config/iam.yml):${self:provider.stage}.roles}
  vpc: ${file(config/vpc.yml):${self:provider.stage}.vpc}

  # Enable X-Ray tracing
  #
  tracing:
    apiGateway: true
    lambda: true


package:
  individually: true
  exclude:
    - ./**
  include:
    - ./bin/**


functions:
  # engine:serverless:functions

  api:
    handler: bin/api

    events:
      - http:
          path: /
          method: ANY
          cors: ${file(config/cors.yml):${self:provider.stage}.cors}

      - http:
          path: /{proxy+}
          method: ANY
          cors: ${file(config/cors.yml):${self:provider.stage}.cors}


plugins:
  - serverless-domain-manager
  - serverless-plugin-warmup`
)
