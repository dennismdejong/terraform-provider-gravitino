package auth

import (
	"context"
	"net/http"
)

type BasicProvider struct {
	username string
	password string
}

func NewBasicProvider(username, password string) *BasicProvider {
	return &BasicProvider{username: username, password: password}
}

func (p *BasicProvider) Header(ctx context.Context) (string, string, error) {
	req := &http.Request{Header: make(http.Header)}
	req.SetBasicAuth(p.username, p.password)
	return "Authorization", req.Header.Get("Authorization"), nil
}
