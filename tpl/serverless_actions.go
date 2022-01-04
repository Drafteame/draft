package tpl

var (
	ServerlessSQSEvent = `# engine:serverless:functions

  {{SnakeCaseName}}:
    handler: bin/{{PackageName}}
    events:
      - sqs:
          batchSize: 1
          functionResponseType: ReportBatchItemFailures
          arn:
            Fn::Join:
              - ':'
              - - arn
                - aws
                - sqs
                - Ref: AWS::Region
                - Ref: AWS::AccountId
                - {{CammelCaseName}}
                - "-${self:provider.stage}"`

	ServerlessPlainEvent = `# engine:serverless:functions

  {{SnakeCaseName}}:
    handler: bin/{{PackageName}}`
)
