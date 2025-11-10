package newlambda

import (
	"fmt"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

func (nl *NewLambda) postCreate() error {
	if nl.input.IsLegacy {
		log.Warn("Command executed in legacy mode. No deps and serverless.yml changes created. Please add manually")
		return nil
	}

	actions := []func() error{
		nl.addToDepsGo,
		nl.addToServerlessYAML,
		nl.format,
		nl.restoreDepsTag, // Restore tag after formatting
	}

	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}

	return nil
}

func (nl *NewLambda) addToServerlessYAML() error {
	path := nl.input.ServicePath + "/serverless.yml"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	typePath := nl.input.LambdaType
	if nl.input.LambdaType == "custom" {
		typePath = nl.input.CustomTypePath
	}

	line := "- ${file(cmd/%s/%s/lambda-config.yml):function}\n  " + data.NextLambdaImportTag
	line = fmt.Sprintf(line, typePath, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), data.NextLambdaImportTag, line)

	return files.Create(path, []byte(newContent))
}

func (nl *NewLambda) addToDepsGo() error {
	path := nl.input.ServicePath + "/deps.go"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	typePath := nl.input.LambdaType
	if nl.input.LambdaType == "custom" {
		typePath = nl.input.CustomTypePath
	}

	importLine := fmt.Sprintf("_ \"%s/%s/cmd/%s/%s/handler\"",
		nl.input.PackageName, nl.input.ServicePath, typePath, nl.input.LambdaName)

	// Check if tag exists inside import block
	if strings.Contains(string(content), data.NextImportTag) {
		// Replace tag with import + tag on same line to keep it in import block
		line := importLine + "\n\t" + data.NextImportTag
		newContent := strings.ReplaceAll(string(content), data.NextImportTag, line)
		return files.Create(path, []byte(newContent))
	}

	// Fallback: add import before closing parenthesis
	lines := strings.Split(string(content), "\n")
	var newLines []string
	for i, line := range lines {
		newLines = append(newLines, line)
		// Find closing parenthesis of import block
		if strings.TrimSpace(line) == ")" && i > 0 {
			// Insert before closing paren
			newLines[len(newLines)-1] = "\t" + importLine
			newLines = append(newLines, ")")
			// Add remaining lines
			newLines = append(newLines, lines[i+1:]...)
			break
		}
	}

	return files.Create(path, []byte(strings.Join(newLines, "\n")))
}

// restoreDepsTag ensures the draft:next-import tag is present after formatting
func (nl *NewLambda) restoreDepsTag() error {
	path := nl.input.ServicePath + "/deps.go"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// If tag already exists, nothing to do
	if strings.Contains(contentStr, data.NextImportTag) {
		return nil
	}

	// Tag was removed by formatter, re-add it
	// Find the import block and add tag before closing parenthesis
	lines := strings.Split(contentStr, "\n")
	var newLines []string
	var inImportBlock bool
	var importBlockStart int

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect import block start
		if strings.HasPrefix(trimmed, "import (") ||
			(strings.Contains(trimmed, "import") && strings.Contains(trimmed, "(")) {
			inImportBlock = true
			importBlockStart = i
		}

		// Find closing parenthesis of import block
		if inImportBlock && trimmed == ")" {
			// Add tag before closing paren
			newLines = append(newLines, "\t"+data.NextImportTag)
			inImportBlock = false
		}

		newLines = append(newLines, line)
	}

	// If no import block found, return as-is
	if importBlockStart == 0 {
		return nil
	}

	return files.Create(path, []byte(strings.Join(newLines, "\n")))
}
