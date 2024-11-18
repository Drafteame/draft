package reviveextrarules

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func Test_IncorrectConfigImport(t *testing.T) {
	incorrectConfigImportRule := &IncorrectConfigImport{}
	for i, testCase := range []struct {
		importStrings []string
		filePath      string
		expectedError bool
	}{
		{
			importStrings: []string{"cmd/foo-service/config"},
			filePath:      "cmd/foo-service/plain/main.go",
			expectedError: false,
		},
		{
			importStrings: []string{"cmd/foo-service/config"},
			filePath:      "cmd/foo-service/http/api/main.go",
			expectedError: false,
		},
		{
			importStrings: []string{"cmd/foo-service/config"},
			filePath:      "cmd/bar-service/http/api/main.go",
			expectedError: true,
		},
		{
			importStrings: []string{"cmd/foo-service/config"},
			filePath:      "bar-service/http/api/main.go",
			expectedError: false,
		},
		{
			importStrings: []string{"cmd/foo-service/config"},
			filePath:      "cmd/bar-service/http/api/baz.go",
			expectedError: false,
		},
	} {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			result := incorrectConfigImportRule.Apply(&lint.File{
				Name: testCase.filePath,
				AST: &ast.File{
					Imports: lo.Map(testCase.importStrings, func(s string, _ int) *ast.ImportSpec {
						return &ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("\"%s\"", s)}}
					}),
				},
			}, nil)
			assert.Equal(t, testCase.expectedError, len(result) > 0, "Expected error")
		})
	}
}
