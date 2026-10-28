package token

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	auth "github.com/Drafteame/draft/internal/actions/auth/token"
)

func TestTokenOutputModes(t *testing.T) {
	jwt := "header." + base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"test-user"}`)) + ".signature"
	for _, tc := range []struct {
		name   string
		flags  []string
		copy   bool
		claims bool
		raw    bool
	}{
		{name: "clipboard", copy: true},
		{name: "raw", flags: []string{"--raw"}, raw: true},
		{name: "raw and copy", flags: []string{"--raw", "--copy"}, raw: true, copy: true},
		{name: "claims", flags: []string{"--raw", "--claims"}, raw: true, claims: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			var stdout, stderr bytes.Buffer
			copied := false
			cmd := newCommand(func(_ context.Context, _ auth.Config, opts auth.LoginOptions) (string, error) {
				_, _ = opts.Progress.Write([]byte("login progress\n"))
				return jwt, nil
			}, func(_ context.Context, value string) error {
				copied = true
				if value != jwt {
					t.Error("incorrect clipboard value")
				}
				return nil
			})
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(append([]string{"dev"}, tc.flags...))
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if tc.raw && stdout.String() != jwt+"\n" {
				t.Fatal("raw stdout contains extra output")
			}
			if !tc.raw && stdout.Len() != 0 {
				t.Fatal("default output exposed token")
			}
			if strings.Contains(stderr.String(), jwt) {
				t.Fatal("token leaked to stderr")
			}
			if copied != tc.copy {
				t.Fatal("unexpected clipboard behavior")
			}
			if strings.Contains(stderr.String(), "signature NOT verified") != tc.claims {
				t.Fatal("unexpected claims output")
			}
		})
	}
}

func TestFlagsOverrideFileAndEnvironment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".draftea", "draft")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte("[auth.workos.dev]\nclient_id = 'file-client'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DRAFT_WORKOS_CLIENT_ID", "env-client")
	cmd := newCommand(func(_ context.Context, cfg auth.Config, _ auth.LoginOptions) (string, error) {
		if cfg.ClientID != "flag-client" {
			t.Error("flag did not override environment and file")
		}
		return "test-token", nil
	}, func(context.Context, string) error { return nil })
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"dev", "--client-id", "flag-client"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestAliasAndErrors(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := &cobra.Command{Use: "draft"}
	root.AddCommand(newCommand(func(context.Context, auth.Config, auth.LoginOptions) (string, error) { return "token", nil }, func(context.Context, string) error { return nil }))
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"token", "prod", "--raw"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{}, {"dev", "prod"}, {"unknown"}} {
		cmd := newCommand(func(context.Context, auth.Config, auth.LoginOptions) (string, error) {
			t.Error("invalid arguments started login")
			return "", nil
		}, func(context.Context, string) error { return nil })
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs(args)
		if cmd.Execute() == nil {
			t.Errorf("invalid arguments accepted: %v", args)
		}
	}
	cmd := newCommand(func(context.Context, auth.Config, auth.LoginOptions) (string, error) {
		return "", errors.New("cancelled")
	}, func(context.Context, string) error { t.Error("copied after login failure"); return nil })
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"dev"})
	if cmd.Execute() == nil {
		t.Fatal("login failure ignored")
	}
}

func TestClipboardFailureDoesNotPrintToken(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	var stdout, stderr bytes.Buffer
	cmd := newCommand(func(context.Context, auth.Config, auth.LoginOptions) (string, error) { return "private-token", nil }, func(context.Context, string) error { return errors.New("clipboard unavailable; use --raw") })
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"dev"})
	if cmd.Execute() == nil {
		t.Fatal("clipboard failure ignored")
	}
	if stdout.Len() != 0 || strings.Contains(stderr.String(), "private-token") {
		t.Fatal("clipboard failure leaked token")
	}
}
