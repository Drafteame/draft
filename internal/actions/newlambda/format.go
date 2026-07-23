package newlambda

import (
	"github.com/Drafteame/draft/internal/pkg/format"
)

func (nl *NewLambda) format() error {
	depsPath := nl.input.ServicePath + "/deps.go"
	return format.Code(nl.lambdaPath, depsPath)
}
