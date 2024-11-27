package dtos

type Input struct {
	PackageName           string
	ServiceFramework      string
	ServiceName           string
	NormalizedServiceName string
	ServicePath           string
	LambdaName            string
	LambdaType            string
	CustomDomain          bool
	DomainPath            string
	WarmupEnabled         bool
	FrameVersion          string
	QueueARN              string
	HTTPPath              string
	HTTPMethod            string
	CronExpression        string
}
