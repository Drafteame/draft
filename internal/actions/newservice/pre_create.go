package newservice

func (ns *NewService) preCreate() error {
	return ns.sentry()
}
