package templates

type Service struct {
	CreateGo        []byte
	CreateTestGo    []byte
	DeleteGo        []byte
	DeleteTestGo    []byte
	GetGo           []byte
	GetTestGo       []byte
	InterfacesGo    []byte
	SearchGo        []byte
	SearchTestGo    []byte
	SearchOneGo     []byte
	SearchOneTestGo []byte
	ServiceGo       []byte
	ServiceTestGo   []byte
	UpdateGo        []byte
	UpdateTestGo    []byte
}

func loadDomainsService(v *Service, data any) error {
	loaders := []func(*Service, any) error{
		loadServiceCreateGo,
		loadServiceCreateTestGo,
		loadServiceDeleteGo,
		loadServiceDeleteTestGo,
		loadServiceGetGo,
		loadServiceGetTestGo,
		loadServiceInterfaces,
		loadServiceSearchGo,
		loadServiceSearchTestGo,
		loadServiceSearchOneGo,
		loadServiceSearchOneTestGo,
		loadServiceServiceGo,
		loadServiceServiceTestGo,
		loadServiceUpdateGo,
		loadServiceUpdateTestGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadServiceCreateGo(v *Service, data any) error {
	name := "domains/service/create.go"
	path := "tmpl/domain/service/create.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateGo = content

	return nil
}

func loadServiceCreateTestGo(v *Service, data any) error {
	name := "domains/service/create_test.go"
	path := "tmpl/domain/service/create_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateTestGo = content

	return nil
}

func loadServiceDeleteGo(v *Service, data any) error {
	name := "domains/service/delete.go"
	path := "tmpl/domain/service/delete.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteGo = content

	return nil
}

func loadServiceDeleteTestGo(v *Service, data any) error {
	name := "domains/service/delete_test.go"
	path := "tmpl/domain/service/delete_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteTestGo = content

	return nil
}

func loadServiceGetGo(v *Service, data any) error {
	name := "domains/service/get.go"
	path := "tmpl/domain/service/get.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetGo = content

	return nil
}

func loadServiceGetTestGo(v *Service, data any) error {
	name := "domains/service/get_test.go"
	path := "tmpl/domain/service/get_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetTestGo = content

	return nil
}

func loadServiceInterfaces(v *Service, data any) error {
	name := "domains/service/interfaces.go"
	path := "tmpl/domain/service/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.InterfacesGo = content

	return nil
}

func loadServiceSearchGo(v *Service, data any) error {
	name := "domains/service/search.go"
	path := "tmpl/domain/service/search.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchGo = content

	return nil
}

func loadServiceSearchTestGo(v *Service, data any) error {
	name := "domains/service/search_test.go"
	path := "tmpl/domain/service/search_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchTestGo = content

	return nil
}

func loadServiceSearchOneGo(v *Service, data any) error {
	name := "domains/service/search_one.go"
	path := "tmpl/domain/service/search_one.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneGo = content

	return nil
}

func loadServiceSearchOneTestGo(v *Service, data any) error {
	name := "domains/service/search_one_test.go"
	path := "tmpl/domain/service/search_one_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneTestGo = content

	return nil
}

func loadServiceServiceGo(v *Service, data any) error {
	name := "domains/service/service.go"
	path := "tmpl/domain/service/service.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceGo = content

	return nil
}

func loadServiceServiceTestGo(v *Service, data any) error {
	name := "domains/service/service_test.go"
	path := "tmpl/domain/service/service_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceTestGo = content

	return nil
}

func loadServiceUpdateGo(v *Service, data any) error {
	name := "domains/service/update.go"
	path := "tmpl/domain/service/update.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateGo = content

	return nil
}

func loadServiceUpdateTestGo(v *Service, data any) error {
	name := "domains/service/update_test.go"
	path := "tmpl/domain/service/update_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateTestGo = content

	return nil
}
