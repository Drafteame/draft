package connect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscapePgpass(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "name.lastname@draftea.com", "name.lastname@draftea.com"},
		{"colon", "a:b", `a\:b`},
		{"backslash", `a\b`, `a\\b`},
		{"token-like", `host:5432/?Action=connect&X-Amz-Signature=ab\cd`, `host\:5432/?Action=connect&X-Amz-Signature=ab\\cd`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, escapePgpass(tt.in))
		})
	}
}

func TestPgpassLine(t *testing.T) {
	line := pgpassLine(56000, "dev@draftea.com", "tok:en")
	assert.Equal(t, `localhost:56000:*:dev@draftea.com:tok\:en`, line)
}

func TestUpsertPgpassEntry_CreatesFileWithOwnerOnlyPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", ".pgpass")

	require.NoError(t, upsertPgpassEntry(path, 56000, "dev@draftea.com", "token-1"))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "localhost:56000:*:dev@draftea.com:token-1\n", string(data))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, pgpassFileMode, info.Mode().Perm())
}

func TestUpsertPgpassEntry_ReplacesManagedLineAndKeepsOthers(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".pgpass")

	initial := "db.example.com:5432:app:alice:secret\n" +
		"localhost:56000:*:dev@draftea.com:old-token\n" +
		"localhost:56011:*:dev@draftea.com:prod-token\n" +
		"localhost:56000:*:other@draftea.com:other-token\n"
	require.NoError(t, os.WriteFile(path, []byte(initial), 0644))

	require.NoError(t, upsertPgpassEntry(path, 56000, "dev@draftea.com", "new:token"))

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	want := "db.example.com:5432:app:alice:secret\n" +
		"localhost:56011:*:dev@draftea.com:prod-token\n" +
		"localhost:56000:*:other@draftea.com:other-token\n" +
		`localhost:56000:*:dev@draftea.com:new\:token` + "\n"
	assert.Equal(t, want, string(data))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, pgpassFileMode, info.Mode().Perm(), "existing file must be tightened to 0600")
}

func TestRemovePgpassEntry(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".pgpass")

	initial := "db.example.com:5432:app:alice:secret\n" +
		"localhost:56000:*:dev@draftea.com:token\n"
	require.NoError(t, os.WriteFile(path, []byte(initial), 0600))

	require.NoError(t, removePgpassEntry(path, 56000, "dev@draftea.com"))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "db.example.com:5432:app:alice:secret\n", string(data))
}

func TestRemovePgpassEntry_MissingFileOrEntryIsNoop(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, removePgpassEntry(filepath.Join(dir, "missing"), 56000, "dev@draftea.com"))
	assert.NoFileExists(t, filepath.Join(dir, "missing"))

	path := filepath.Join(dir, ".pgpass")
	require.NoError(t, os.WriteFile(path, []byte("other:1:*:u:p\n"), 0600))
	require.NoError(t, removePgpassEntry(path, 56000, "dev@draftea.com"))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "other:1:*:u:p\n", string(data))
}

func TestPgpassPath(t *testing.T) {
	t.Setenv("PGPASSFILE", "")
	t.Setenv("HOME", "/tmp/fake-home")

	path, err := pgpassPath()
	require.NoError(t, err)
	assert.Equal(t, "/tmp/fake-home/.pgpass", path)

	t.Setenv("PGPASSFILE", "/custom/pgpass")

	path, err = pgpassPath()
	require.NoError(t, err)
	assert.Equal(t, "/custom/pgpass", path)
}
