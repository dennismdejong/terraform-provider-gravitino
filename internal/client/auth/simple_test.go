package auth_test

import (
	"context"
	"encoding/base64"
	"os"
	"os/user"
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
)

func TestSimpleProvider_Header(t *testing.T) {
	t.Run("uses configured username", func(t *testing.T) {
		p := auth.NewSimpleProvider("testuser")
		key, value, err := p.Header(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if key != "Authorization" {
			t.Fatalf("expected Authorization, got %s", key)
		}
		expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:"))
		if value != expected {
			t.Fatalf("expected %q, got %q", expected, value)
		}
	})

	t.Run("falls back to GRAVITINO_USER env", func(t *testing.T) {
		os.Setenv("GRAVITINO_USER", "envuser")
		defer os.Unsetenv("GRAVITINO_USER")
		p := auth.NewSimpleProvider("")
		_, value, err := p.Header(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("envuser:"))
		if value != expected {
			t.Fatalf("expected %q, got %q", expected, value)
		}
	})

	t.Run("falls back to OS user when empty", func(t *testing.T) {
		p := auth.NewSimpleProvider("")
		key, value, err := p.Header(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if key != "Authorization" {
			t.Fatalf("expected Authorization, got %s", key)
		}
		u, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		expectedUser := u.Username

		prefix := "Basic "
		decoded, err := base64.StdEncoding.DecodeString(value[len(prefix):])
		if err != nil {
			t.Fatal(err)
		}
		gotUser := string(decoded[:len(decoded)-1])
		if gotUser != expectedUser {
			t.Fatalf("expected user %q, got %q", expectedUser, gotUser)
		}
	})
}
