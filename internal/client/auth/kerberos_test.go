package auth_test

import (
	"testing"

	"github.com/gravitino/terraform-provider-gravitino/internal/client/auth"
)

func TestKerberosProvider_TransportProvider(t *testing.T) {
	var _ auth.TransportProvider = (*auth.KerberosProvider)(nil)
}

func TestKerberosProvider_KeytabFileNotFound(t *testing.T) {
	_, err := auth.NewKerberosProvider("HTTP/test@REALM", "/nonexistent/keytab", false)
	if err == nil {
		t.Fatal("expected error for nonexistent keytab")
	}
}

func TestKerberosProvider_WrapTransport(t *testing.T) {
	t.Skip("requires Kerberos KDC; run integration tests manually with TEST_KERBEROS=1")
}
