package dtos

type DomainInput struct {
	PackageName        string
	DomainPath         string
	DomainName         string
	DomainNamePascal   string
	DomainNameLower    string
	DBPrefix           string
	TableName          string
	DBProviderFuncName string
	DBType             string
	DBName             string
}
