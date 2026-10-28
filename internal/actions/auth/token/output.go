package token

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func openBrowser(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var command string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		command = "open"
		args = []string{url}
	case "windows":
		command = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		command = "xdg-open"
		args = []string{url}
	}
	return exec.CommandContext(ctx, command, args...).Run()
}

// Copy copies a token using platform tools, passing its value only through stdin.
func Copy(ctx context.Context, token string) error {
	var commands [][]string
	switch runtime.GOOS {
	case "darwin":
		commands = [][]string{{"pbcopy"}}
	case "windows":
		commands = [][]string{{"powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "$input | Set-Clipboard"}}
	default:
		commands = [][]string{{"wl-copy"}, {"xclip", "-selection", "clipboard"}, {"xsel", "--clipboard", "--input"}}
	}
	for _, args := range commands {
		copyCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		cmd := exec.CommandContext(copyCtx, args[0], args[1:]...)
		cmd.Stdin = strings.NewReader(token)
		err := cmd.Run()
		cancel()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("could not copy the token; install a clipboard tool or rerun with --raw to print it explicitly")
}

// Claims decodes JWT claims for inspection only; it does not verify the signature.
func Claims(token string, w io.Writer) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("access token is not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("JWT claims are not valid base64url")
	}
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, payload, "", "  "); err != nil {
		return fmt.Errorf("JWT claims are not valid JSON")
	}
	_, err = fmt.Fprintf(w, "JWT claims (decoded only; signature NOT verified):\n%s\n", formatted.String())
	return err
}
