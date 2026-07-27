package auth

import "context"

type OAuthStaticProvider struct {
	token string
}

func NewOAuthStaticProvider(token string) *OAuthStaticProvider {
	return &OAuthStaticProvider{token: token}
}

func (p *OAuthStaticProvider) Header(ctx context.Context) (string, string, error) {
	return "Authorization", "Bearer " + p.token, nil
}
