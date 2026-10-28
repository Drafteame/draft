package token

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func testFlow(t *testing.T, handler http.HandlerFunc) (Config, loginDependencies) {
	t.Helper()
	endpoint := httptest.NewTLSServer(handler)
	t.Cleanup(endpoint.Close)
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	cfg := Config{ClientID: "public-client", Provider: "GoogleOAuth", Port: listener.Addr().(*net.TCPAddr).Port, AuthURL: "https://login.example/authorize", TokenURL: endpoint.URL}
	return cfg, loginDependencies{client: endpoint.Client(), listen: func(string, string) (net.Listener, error) { return listener, nil }}
}

func callback(t *testing.T, loginURL string, alter func(url.Values)) int {
	t.Helper()
	auth, err := url.Parse(loginURL)
	if err != nil {
		t.Fatal(err)
	}
	target, err := url.Parse(auth.Query().Get("redirect_uri"))
	if err != nil {
		t.Fatal(err)
	}
	q := url.Values{"state": {auth.Query().Get("state")}, "code": {"one-time-code"}}
	if alter != nil {
		alter(q)
	}
	target.RawQuery = q.Encode()
	response, err := http.Get(target.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "secret-access-token") {
		t.Fatal("token leaked into callback response")
	}
	return response.StatusCode
}

func TestPKCELoginRejectsInvalidCallbacksAndReplay(t *testing.T) {
	observed := make(chan map[string]string, 1)
	cfg, deps := testFlow(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
			t.Error("unexpected exchange request")
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		observed <- payload
		_, _ = io.WriteString(w, `{"access_token":"secret-access-token","refresh_token":"ignored-refresh-token"}`)
	})
	var authorize url.Values
	deps.openBrowser = func(_ context.Context, loginURL string) error {
		u, _ := url.Parse(loginURL)
		authorize = u.Query()
		if got := callback(t, loginURL, func(q url.Values) { q.Set("state", "attacker-state") }); got != 400 {
			t.Errorf("invalid state: %d", got)
		}
		if got := callback(t, loginURL, func(q url.Values) { q.Del("code") }); got != 400 {
			t.Errorf("missing code: %d", got)
		}
		if got := callback(t, loginURL, nil); got != 200 {
			t.Errorf("successful callback: %d", got)
		}
		if got := callback(t, loginURL, nil); got != 409 {
			t.Errorf("replayed callback: %d", got)
		}
		return nil
	}
	var output strings.Builder
	accessToken, err := login(context.Background(), cfg, LoginOptions{Timeout: time.Second * 5, Progress: &output}, deps)
	if err != nil || accessToken != "secret-access-token" {
		t.Fatalf("login failed: %v", err)
	}
	payload := <-observed
	sum := sha256.Sum256([]byte(payload["code_verifier"]))
	if authorize.Get("code_challenge") != base64.RawURLEncoding.EncodeToString(sum[:]) || authorize.Get("code_challenge_method") != "S256" {
		t.Fatal("PKCE challenge does not match verifier")
	}
	if payload["grant_type"] != "authorization_code" || payload["client_id"] != cfg.ClientID || payload["code"] != "one-time-code" || payload["redirect_uri"] != authorize.Get("redirect_uri") {
		t.Fatal("incorrect token exchange payload")
	}
	if _, ok := payload["client_secret"]; ok {
		t.Fatal("public client must not send client secret")
	}
	if strings.Contains(output.String(), accessToken) || strings.Contains(output.String(), payload["code_verifier"]) {
		t.Fatal("token or verifier leaked into progress")
	}
	if authorize.Get("state") == "" {
		t.Fatal("missing OAuth state")
	}
}

func TestExchangeFailuresDoNotExposeResponseSecrets(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"provider error", 401, `{"error_description":"secret-access-token"}`},
		{"invalid JSON", 200, `secret-access-token`},
		{"missing token", 200, `{"refresh_token":"secret-access-token"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, deps := testFlow(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			})
			deps.openBrowser = func(_ context.Context, u string) error {
				if callback(t, u, nil) != 502 {
					t.Error("failed exchange was not rejected")
				}
				return nil
			}
			_, err := login(context.Background(), cfg, LoginOptions{Timeout: 5 * time.Second}, deps)
			if err == nil || strings.Contains(err.Error(), "secret-access-token") {
				t.Fatalf("missing or unsafe error: %v", err)
			}
		})
	}
}

func TestOAuthRejection(t *testing.T) {
	cfg, deps := testFlow(t, func(http.ResponseWriter, *http.Request) { t.Error("rejected login should not exchange a token") })
	deps.openBrowser = func(_ context.Context, u string) error {
		if callback(t, u, func(q url.Values) { q.Set("error", "access_denied"); q.Del("code") }) != 400 {
			t.Error("provider rejection was not handled")
		}
		return nil
	}
	if _, err := login(context.Background(), cfg, LoginOptions{Timeout: time.Second}, deps); err == nil {
		t.Fatal("rejected authorization accepted")
	}
}

func TestTimeoutAndCancellationReleaseListener(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		t.Run(fmt.Sprint(cancelled), func(t *testing.T) {
			cfg, deps := testFlow(t, func(http.ResponseWriter, *http.Request) { t.Error("unexpected exchange") })
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if cancelled {
				cancel()
			}
			_, err := login(ctx, cfg, LoginOptions{NoBrowser: true, Timeout: 20 * time.Millisecond}, deps)
			expected := context.DeadlineExceeded
			if cancelled {
				expected = context.Canceled
			}
			if !errors.Is(err, expected) {
				t.Fatalf("unexpected cancellation error: %v", err)
			}
			listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", cfg.Port))
			if err != nil {
				t.Fatalf("callback listener leaked: %v", err)
			}
			_ = listener.Close()
		})
	}
}

func TestBrowserFailureAllowsManualLogin(t *testing.T) {
	cfg, deps := testFlow(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"access_token":"secret-access-token"}`)
	})
	deps.openBrowser = func(_ context.Context, u string) error { callback(t, u, nil); return errors.New("no browser") }
	var progress strings.Builder
	if _, err := login(context.Background(), cfg, LoginOptions{Timeout: time.Second, Progress: &progress}, deps); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(progress.String(), "open the URL above") {
		t.Fatal("manual fallback was not explained")
	}
}

func TestTokenRedirectDoesNotReceiveAuthorizationCode(t *testing.T) {
	redirected := false
	sink := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { redirected = true }))
	defer sink.Close()
	cfg, deps := testFlow(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, sink.URL, http.StatusTemporaryRedirect)
	})
	deps.client.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirect forbidden") }
	deps.openBrowser = func(_ context.Context, u string) error { callback(t, u, nil); return nil }
	if _, err := login(context.Background(), cfg, LoginOptions{Timeout: time.Second}, deps); err == nil {
		t.Fatal("redirect accepted")
	}
	if redirected {
		t.Fatal("token exchange credentials followed a redirect")
	}
}

func TestOccupiedPortDoesNotOpenBrowser(t *testing.T) {
	cfg, deps := testFlow(t, func(http.ResponseWriter, *http.Request) { t.Error("unexpected exchange") })
	deps.listen = net.Listen
	deps.openBrowser = func(context.Context, string) error { t.Error("opened browser with unavailable port"); return nil }
	if _, err := login(context.Background(), cfg, LoginOptions{Timeout: time.Second}, deps); err == nil {
		t.Fatal("occupied callback port accepted")
	}
}
