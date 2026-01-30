package dtos

type ServiceInput struct {
	PackageName           string
	ServiceFramework      string
	ServiceName           string
	NormalizedServiceName string
	ServicePath           string
	ServicePackage        string
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
	HasSentry             bool
	SentryDSN             string
	NextImportTag         string
	NextLambdaImportTag   string
	IsLegacy              bool
	UseDig                bool
	ReservedConcurrency   string
	RoleName              string
}
