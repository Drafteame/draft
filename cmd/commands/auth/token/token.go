package token

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	auth "github.com/Drafteame/draft/internal/actions/auth/token"
)

// GetCmd returns a fresh command with no AWS or Sentry initialization.
func GetCmd() *cobra.Command { return newCommand(auth.Login, auth.Copy) }

func newCommand(login func(context.Context, auth.Config, auth.LoginOptions) (string, error), copyToClipboard func(context.Context, string) error) *cobra.Command {
	var raw, copyToken, claims, noBrowser bool
	var timeout time.Duration
	var clientID, provider, authURL, tokenURL string
	var port int
	cmd := &cobra.Command{
		Use: "auth:token <dev|prod>", Aliases: []string{"token"}, Short: "Sign in with WorkOS and copy an API access token",
		Long: "Sign in in your browser using PKCE. Copies the access token by default; --raw prints only the token to stdout. Permissions remain those of the signed-in account. No tokens are saved to disk.",
		Args: cobra.ExactArgs(1), ValidArgs: []string{"dev", "prod"}, SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := auth.LoadConfig(args[0])
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("client-id") {
				cfg.ClientID = clientID
			}
			if cmd.Flags().Changed("provider") {
				cfg.Provider = provider
			}
			if cmd.Flags().Changed("port") {
				cfg.Port = port
			}
			if cmd.Flags().Changed("auth-url") {
				cfg.AuthURL = authURL
			}
			if cmd.Flags().Changed("token-url") {
				cfg.TokenURL = tokenURL
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Environment: %s\n", args[0])
			accessToken, err := login(cmd.Context(), cfg, auth.LoginOptions{Timeout: timeout, NoBrowser: noBrowser, Progress: cmd.ErrOrStderr()})
			if err != nil {
				return err
			}
			if raw {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), accessToken); err != nil {
					return err
				}
			}
			if !raw || copyToken {
				if err := copyToClipboard(cmd.Context(), accessToken); err != nil {
					return err
				}
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Access token copied to clipboard.")
			}
			if claims {
				return auth.Claims(accessToken, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "Print only the access token to stdout instead of copying it")
	cmd.Flags().BoolVar(&copyToken, "copy", false, "Also copy to clipboard when using --raw")
	cmd.Flags().BoolVar(&claims, "claims", false, "Print decoded JWT claims to stderr (signature is not verified)")
	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Print the login URL without opening the browser")
	cmd.Flags().DurationVar(&timeout, "timeout", 10*time.Minute, "Maximum time to complete login")
	cmd.Flags().StringVar(&clientID, "client-id", "", "Override the public WorkOS client ID")
	cmd.Flags().StringVar(&provider, "provider", "", "Override the WorkOS provider (default: GoogleOAuth)")
	cmd.Flags().IntVar(&port, "port", 0, "Override the registered loopback callback port (default: 3001)")
	cmd.Flags().StringVar(&authURL, "auth-url", "", "Override the HTTPS authorization endpoint")
	cmd.Flags().StringVar(&tokenURL, "token-url", "", "Override the HTTPS token endpoint")
	return cmd
}
