package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCommand_PnpmWhenLockfileExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pnpm-lock.yaml"), []byte("lockfileVersion: '9.0'\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if got := installCommand(dir); got != "pnpm install" {
		t.Fatalf("installCommand() = %q, want %q", got, "pnpm install")
	}
}

func TestInstallCommand_NpmWithoutPnpmLockfile(t *testing.T) {
	t.Parallel()

	if got := installCommand(t.TempDir()); got != "npm install" {
		t.Fatalf("installCommand() = %q, want %q", got, "npm install")
	}
}
