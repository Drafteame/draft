package templates

import (
	appdata "github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

type Domain struct {
	DomainGo []byte
	ErrorsGo []byte
	Options  DomainOptions
}

func loadDomainsDomain(v *Domain, data any) error {
	var loaders []func(*Domain, any) error

	// Only load domain for postgres
	input, ok := data.(dtos.DomainInput)
	if ok && input.DBType == appdata.DBTypePostgres {
		loaders = []func(*Domain, any) error{
			loadDomainDomainGo,
			loadDomainErrorsGo,
			loadDomainOptions,
		}
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadDomainDomainGo(v *Domain, data any) error {
	name := "domains/domain.go"
	path := "tmpl/domain/domain/domain.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DomainGo = content

	return nil
}

func loadDomainErrorsGo(v *Domain, data any) error {
	name := "domains/errors.go"
	path := "tmpl/domain/domain/errors.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ErrorsGo = content

	return nil
}

func loadDomainOptions(v *Domain, data any) error {
	return loadDomainsDomainOptions(&v.Options, data)
}

type DomainOptions struct {
	SearchGo           []byte
	SearchFiltersGo    []byte
	SearchOrdersGo     []byte
	SearchPaginationGo []byte
	UpdateFieldsGo     []byte
}

func loadDomainsDomainOptions(v *DomainOptions, data any) error {
	loaders := []func(*DomainOptions, any) error{
		loadDomainOptionsSearchGo,
		loadDomainOptionsSearchFiltersGo,
		loadDomainOptionsSearchOrdersGo,
		loadDomainOptionsSearchPaginationGo,
		loadDomainOptionsUpdateFieldsGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadDomainOptionsSearchGo(v *DomainOptions, data any) error {
	name := "domains/options/search.go"
	path := "tmpl/domain/domain/options/search.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchGo = content

	return nil
}

func loadDomainOptionsSearchFiltersGo(v *DomainOptions, data any) error {
	name := "domains/options/search_filters.go"
	path := "tmpl/domain/domain/options/search_filters.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchFiltersGo = content

	return nil
}

func loadDomainOptionsSearchOrdersGo(v *DomainOptions, data any) error {
	name := "domains/options/search_orders.go"
	path := "tmpl/domain/domain/options/search_orders.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOrdersGo = content

	return nil
}

func loadDomainOptionsSearchPaginationGo(v *DomainOptions, data any) error {
	name := "domains/options/search_pagination.go"
	path := "tmpl/domain/domain/options/search_pagination.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchPaginationGo = content

	return nil
}

func loadDomainOptionsUpdateFieldsGo(v *DomainOptions, data any) error {
	name := "domains/options/update_fields.go"
	path := "tmpl/domain/domain/options/update_fields.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateFieldsGo = content

	return nil
}
