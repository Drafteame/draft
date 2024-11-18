package reviveextrarules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mgechev/revive/lint"
	"github.com/samber/lo"
)

// IncorrectConfigImport lints given else constructs.
type IncorrectConfigImport struct{}

// Apply applies the rule to given file.
// In a file with path matching **/cmd/<path>/main.go, check that the any import statement
// matching **/cmd/<service-name>/config verifies <path> contains "cmd/<service-name>" as substring
func (*IncorrectConfigImport) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	filePath := file.Name
	filePathItems := strings.Split(filePath, "/")
	isLambdaMainFile := slices.Contains(filePathItems, "cmd") && filePathItems[len(filePathItems)-1] == "main.go"
	if !isLambdaMainFile {
		return nil
	}

	for _, imp := range file.AST.Imports {
		// check if imp.Name ends is of glob cmd/**/config
		impPath := strings.Replace(imp.Path.Value, "\"", "", -1)
		impPathItems := strings.Split(impPath, "/")
		isConfigImport := slices.Contains(impPathItems, "cmd") && impPathItems[len(impPathItems)-1] == "config"
		if !isConfigImport {
			continue
		}

		_, cmdPathIndex, _ := lo.FindIndexOf(impPathItems, func(s string) bool { return s == "cmd" })
		serviceNameFromConfigImport := strings.Join(impPathItems[(cmdPathIndex+1):len(impPathItems)-1], "/")
		if strings.Contains(filePath, "cmd/"+serviceNameFromConfigImport+"/") {
			continue
		}

		return []lint.Failure{
			{
				Confidence: 1,
				Failure:    fmt.Sprintf("Importing incorrect config: %s", serviceNameFromConfigImport),
				Node:       imp,
				Category:   "imports",
			},
		}
	}

	return failures
}

// Name returns the rule name.
func (*IncorrectConfigImport) Name() string {
	return "incorrect-config-import"
}
