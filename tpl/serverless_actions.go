package tpl

var (
	ServerlessSQSEvent = `# engine:serverless:functions

  {{SnakeCaseName}}:
    handler: bin/{{PackageName}}
    events:
      - sqs:
          batchSize: 1
          functionResponseType: ReportBatchItemFailures
          arn: {{TriggerArn}}`

	ServerlessPlainEvent = `# engine:serverless:functions

  {{SnakeCaseName}}:
    handler: bin/{{PackageName}}`
)
