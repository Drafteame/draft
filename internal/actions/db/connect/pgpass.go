package connect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/dirs"
)

const (
	// pgpassFileEnv is honoured by libpq to override the location of the password file.
	pgpassFileEnv = "PGPASSFILE"

	// pgpassHost is the host written to the password file. Clients connect to the
	// local end of the tunnel, not to the real RDS hostname.
	pgpassHost = "localhost"

	// libpq refuses the file unless it is readable only by the owner.
	pgpassFileMode = os.FileMode(0600)
)

// pgpassPath returns the password file psql will read: $PGPASSFILE or ~/.pgpass.
func pgpassPath() (string, error) {
	if custom := strings.TrimSpace(os.Getenv(pgpassFileEnv)); custom != "" {
		return custom, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(home, ".pgpass"), nil
}

// escapePgpass escapes the characters that are field separators in .pgpass.
// The IAM token contains ':' and '&', only ':' and '\' need escaping.
func escapePgpass(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)

	return strings.ReplaceAll(value, ":", `\:`)
}

// pgpassPrefix builds the "host:port:database:user:" prefix that identifies the
// entry managed by draft for one tunnel + user. The database is always '*'.
func pgpassPrefix(localPort int, user string) string {
	return fmt.Sprintf("%s:%d:*:%s:", pgpassHost, localPort, escapePgpass(user))
}

// pgpassLine builds the full .pgpass entry for a tunnel + user + token.
func pgpassLine(localPort int, user, token string) string {
	return pgpassPrefix(localPort, user) + escapePgpass(token)
}

// upsertPgpassEntry replaces (or appends) the entry for localhost:<localPort>
// and user in the given password file, leaving every other line untouched.
func upsertPgpassEntry(path string, localPort int, user, token string) error {
	lines, err := readPgpassLines(path)
	if err != nil {
		return err
	}

	prefix := pgpassPrefix(localPort, user)
	kept := make([]string, 0, len(lines)+1)

	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			continue
		}

		kept = append(kept, line)
	}

	kept = append(kept, pgpassLine(localPort, user, token))

	return writePgpassLines(path, kept)
}

// removePgpassEntry deletes the entry for localhost:<localPort> and user from the
// given password file. A missing file or entry is not an error.
func removePgpassEntry(path string, localPort int, user string) error {
	lines, err := readPgpassLines(path)
	if err != nil {
		return err
	}

	prefix := pgpassPrefix(localPort, user)
	kept := make([]string, 0, len(lines))
	changed := false

	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			changed = true
			continue
		}

		kept = append(kept, line)
	}

	if !changed {
		return nil
	}

	return writePgpassLines(path, kept)
}

func readPgpassLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return nil, nil
	}

	return strings.Split(content, "\n"), nil
}

func writePgpassLines(path string, lines []string) error {
	if err := dirs.Create(filepath.Dir(path)); err != nil {
		return err
	}

	content := ""
	if len(lines) > 0 {
		content = strings.Join(lines, "\n") + "\n"
	}

	if err := os.WriteFile(path, []byte(content), pgpassFileMode); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	// WriteFile only applies the mode on creation; enforce it on existing files too.
	if err := os.Chmod(path, pgpassFileMode); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", path, err)
	}

	return nil
}
