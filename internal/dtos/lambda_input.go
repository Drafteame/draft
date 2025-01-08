package dtos

type LambdaInput struct {
	PackageName         string
	ServiceFramework    string
	ServicePath         string
	LambdaName          string
	LambdaType          string
	FrameVersion        string
	QueueARN            string
	HTTPPath            string
	HTTPPathAPIGateway  string
	HTTPPathEcho        string
	HTTPMethod          string
	CronExpression      string
	NextImportTag       string
	NextLambdaImportTag string
	IsLegacy            bool
}
