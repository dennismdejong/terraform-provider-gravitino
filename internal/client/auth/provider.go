package auth

import (
	"context"
	"net/http"
)

type AuthProvider interface {
	Header(ctx context.Context) (string, string, error)
}

type TransportProvider interface {
	AuthProvider
	WrapTransport(base http.RoundTripper) http.RoundTripper
}
