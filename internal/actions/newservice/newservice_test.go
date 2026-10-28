package newservice

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/templates"
)

func TestCreateFiles_WritesPnpmFiles(t *testing.T) {
	t.Parallel()

	input := dtos.ServiceInput{
		PackageName:           "github.com/Drafteame/api",
		ServiceFramework:      "sls",
		ServiceName:           "testsvc",
		NormalizedServiceName: "testsvc",
		ServicePath:           t.TempDir(),
		ServicePackage:        "testsvc",
		LambdaName:            "helloworld",
		LambdaType:            "plain",
		FrameVersion:          "v2",
		ReservedConcurrency:   "medium.http",
		RoleName:              "Testsvc",
	}

	tmpl, err := templates.NewServiceTemplates(input)
	require.NoError(t, err)

	ns := &NewService{tmpl: tmpl, input: input}
	require.NoError(t, ns.createAllDirs())
	require.NoError(t, ns.createLambdaFolders())
	require.NoError(t, ns.createFiles())

	pkg, err := os.ReadFile(filepath.Join(input.ServicePath, "package.json"))
	require.NoError(t, err)
	assert.Contains(t, string(pkg), `"packageManager": "pnpm@`)
	assert.NotContains(t, string(pkg), `"overrides"`)

	ws, err := os.ReadFile(filepath.Join(input.ServicePath, "pnpm-workspace.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(ws), `"simple-git": "~3.33.0"`)
	assert.Contains(t, string(ws), "blockExoticSubdeps: true")
	assert.Contains(t, string(ws), "serverless: false")

	assert.NoFileExists(t, filepath.Join(input.ServicePath, "package-lock.json"))
}
