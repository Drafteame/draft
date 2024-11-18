package dtos

type Input struct {
	ServiceFramework string
	ServiceName      string
	ServicePath      string
	LambdaName       string
	LambdaType       string
	CustomDomain     bool
	DomainPath       string
	WarmupEnabled    bool
	FrameVersion     string
	QueueARN         string
	HTTPPath         string
	HTTPMethod       string
	CronExpression   string
}
