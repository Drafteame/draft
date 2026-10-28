package token

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestClaimsAreExplicitlyUnverified(t *testing.T) {
	var output bytes.Buffer
	jwt := "header." + base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"test-user"}`)) + ".signature"
	if err := Claims(jwt, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "signature NOT verified") || !strings.Contains(output.String(), "test-user") {
		t.Fatal("missing claims disclaimer or payload")
	}
	for _, token := range []string{"private-opaque-token", "a.bad-base64!.c", "a." + base64.RawURLEncoding.EncodeToString([]byte("private-invalid-json")) + ".c"} {
		err := Claims(token, &output)
		if err == nil || strings.Contains(err.Error(), token) || strings.Contains(err.Error(), "private-invalid-json") {
			t.Fatal("invalid claims accepted or leaked")
		}
	}
}
