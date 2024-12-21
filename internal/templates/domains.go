package templates

type Domains struct {
	Service    Service
	Repository Repository
	Providers  Providers
	Domain     Domain
}

func NewDomains(data any) (Domains, error) {
	d := Domains{}

	loaders := []func(*Domains, any) error{
		loadService,
		loadRepository,
		loadProviders,
		loadDomain,
	}

	for _, loader := range loaders {
		if err := loader(&d, data); err != nil {
			return d, err
		}
	}

	return d, nil
}

func loadService(d *Domains, data any) error {
	return loadDomainsService(&d.Service, data)
}

func loadRepository(d *Domains, data any) error {
	return loadDomainsRepository(&d.Repository, data)
}

func loadProviders(d *Domains, data any) error {
	return loadDomainsProviders(&d.Providers, data)
}

func loadDomain(d *Domains, data any) error {
	return loadDomainsDomain(&d.Domain, data)
}
