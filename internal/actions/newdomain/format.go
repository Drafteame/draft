package newdomain

import (
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/format"
)

func (nd *NewDomain) format() error {
	paths := []string{nd.input.DomainPath}

	// For Postgres, also format the modified provider files
	if nd.input.DBType == data.DBTypePostgres {
		paths = append(paths, "pkg/providers")
	}

	return format.Code(paths...)
}
