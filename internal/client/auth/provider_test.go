package auth_test

import (
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
)

func TestAuthProviderInterface(t *testing.T) {
	var _ auth.AuthProvider = (*auth.SimpleProvider)(nil)
}
