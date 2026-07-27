package auth_test

import (
	"context"
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
)

func TestBasicProvider_Header(t *testing.T) {
	p := auth.NewBasicProvider("admin", "secret")
	key, value, err := p.Header(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if key != "Authorization" {
		t.Fatalf("expected Authorization, got %s", key)
	}
	if value == "" {
		t.Fatal("expected non-empty header value")
	}
}
