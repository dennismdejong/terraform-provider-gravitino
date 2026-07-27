package auth

import (
	"context"
	"encoding/base64"
	"os"
	"os/user"
)

type SimpleProvider struct {
	username string
}

func NewSimpleProvider(username string) *SimpleProvider {
	return &SimpleProvider{username: username}
}

func resolveUsername(configured string) string {
	if configured != "" {
		return configured
	}
	if env := os.Getenv("GRAVITINO_USER"); env != "" {
		return env
	}
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return "anonymous"
}

func (p *SimpleProvider) Header(ctx context.Context) (string, string, error) {
	u := resolveUsername(p.username)
	auth := base64.StdEncoding.EncodeToString([]byte(u + ":"))
	return "Authorization", "Basic " + auth, nil
}
