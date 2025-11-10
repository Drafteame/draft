package newlambda

import "errors"

func (nl *NewLambda) exec() error {
	switch nl.input.LambdaType {
	case "plain":
		return nl.createPlain()
	case "sqs":
		return nl.createSqs()
	case "http":
		return nl.createHttp()
	case "snssqs":
		return nl.createSnsSqs()
	case "cron":
		return nl.createCron()
	case "custom":
		return nl.createCustom()
	default:
		return errors.New("unsupported lambda type")
	}
}
