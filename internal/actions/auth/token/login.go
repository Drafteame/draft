package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// LoginOptions controls interaction without persisting session information.
type LoginOptions struct {
	Timeout   time.Duration
	NoBrowser bool
	Progress  io.Writer
}

type loginDependencies struct {
	client      *http.Client
	listen      func(string, string) (net.Listener, error)
	openBrowser func(context.Context, string) error
}

type loginResult struct {
	token string
	err   error
}

// Login obtains an access token through a one-shot loopback authorization-code flow.
func Login(ctx context.Context, cfg Config, opts LoginOptions) (string, error) {
	return login(ctx, cfg, opts, loginDependencies{
		client: &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return errors.New("token endpoint redirects are not allowed")
		}},
		listen: net.Listen, openBrowser: openBrowser,
	})
}

func login(ctx context.Context, cfg Config, opts LoginOptions, deps loginDependencies) (string, error) {
	if err := cfg.Validate(); err != nil {
		return "", err
	}
	if opts.Timeout <= 0 {
		return "", fmt.Errorf("login timeout must be positive")
	}
	if opts.Progress == nil {
		opts.Progress = io.Discard
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()
	verifier, err := randomValue()
	if err != nil {
		return "", err
	}
	state, err := randomValue()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	listener, err := deps.listen("tcp4", fmt.Sprintf("127.0.0.1:%d", cfg.Port))
	if err != nil {
		return "", fmt.Errorf("listen on callback port %d (choose --port and register its callback in WorkOS if occupied): %w", cfg.Port, err)
	}
	defer func() { _ = listener.Close() }()
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", listener.Addr().(*net.TCPAddr).Port)
	results := make(chan loginResult, 1)
	var mu sync.Mutex
	exchanging := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if r.URL.Path != "/callback" || r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		query := r.URL.Query()
		if subtle.ConstantTimeCompare([]byte(query.Get("state")), []byte(state)) != 1 {
			http.Error(w, "Invalid OAuth state; use the original login link.", http.StatusBadRequest)
			return
		}
		mu.Lock()
		if exchanging {
			mu.Unlock()
			http.Error(w, "This login callback has already been processed.", http.StatusConflict)
			return
		}
		if query.Get("error") == "" && query.Get("code") == "" {
			mu.Unlock()
			http.Error(w, "Authorization code missing.", http.StatusBadRequest)
			return
		}
		exchanging = true
		mu.Unlock()
		if query.Get("error") != "" {
			http.Error(w, "Authentication was rejected. Return to the terminal.", http.StatusBadRequest)
			results <- loginResult{err: fmt.Errorf("WorkOS authentication was rejected")}
			return
		}
		token, err := exchange(ctx, deps.client, cfg, query.Get("code"), verifier, redirectURI)
		if err != nil {
			http.Error(w, "Token exchange failed. Return to the terminal.", http.StatusBadGateway)
		} else {
			_, _ = io.WriteString(w, "Authentication complete. Return to the terminal; you can close this window.")
		}
		results <- loginResult{token: token, err: err}
	})
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 10 * time.Second}
	defer func() {
		cancel()
		shutdownCtx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
		}
	}()
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- server.Serve(listener) }()
	authURL, _ := url.Parse(cfg.AuthURL)
	q := authURL.Query()
	for k, v := range map[string]string{"client_id": cfg.ClientID, "provider": cfg.Provider, "redirect_uri": redirectURI, "response_type": "code", "state": state, "code_challenge": challenge, "code_challenge_method": "S256"} {
		q.Set(k, v)
	}
	authURL.RawQuery = q.Encode()
	_, _ = fmt.Fprintf(opts.Progress, "Callback: %s\nOpen this URL to sign in:\n%s\n", redirectURI, authURL.String())
	if !opts.NoBrowser {
		if err := deps.openBrowser(ctx, authURL.String()); err != nil {
			_, _ = fmt.Fprintln(opts.Progress, "Could not open the browser automatically; open the URL above.")
		}
	}
	select {
	case result := <-results:
		return result.token, result.err
	case <-ctx.Done():
		return "", fmt.Errorf("login did not complete: %w", ctx.Err())
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return "", fmt.Errorf("login callback server closed")
		}
		return "", fmt.Errorf("serve login callback: %w", err)
	}
}

func randomValue() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate OAuth random value: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func exchange(ctx context.Context, client *http.Client, cfg Config, code, verifier, redirectURI string) (string, error) {
	body, err := json.Marshal(map[string]string{"grant_type": "authorization_code", "client_id": cfg.ClientID, "code": code, "code_verifier": verifier, "redirect_uri": redirectURI})
	if err != nil {
		return "", fmt.Errorf("encode token exchange request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.TokenURL, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("create token exchange request")
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("WorkOS token exchange request failed")
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("WorkOS token exchange failed (HTTP %d)", response.StatusCode)
	}
	var data struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&data); err != nil {
		return "", fmt.Errorf("WorkOS returned an invalid token response")
	}
	if data.AccessToken == "" {
		return "", fmt.Errorf("WorkOS response did not contain an access token")
	}
	return data.AccessToken, nil
}
