package auth_test

import (
	"context"
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
)

func TestOAuthStaticProvider_Header(t *testing.T) {
	p := auth.NewOAuthStaticProvider("my-token")
	key, value, err := p.Header(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if key != "Authorization" {
		t.Fatalf("expected Authorization, got %s", key)
	}
	if value != "Bearer my-token" {
		t.Fatalf("expected 'Bearer my-token', got %s", value)
	}
}
