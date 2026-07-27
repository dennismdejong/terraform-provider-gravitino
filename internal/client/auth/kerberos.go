package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

type KerberosProvider struct {
	principal string
	krbClient *client.Client
	transport http.RoundTripper
}

func NewKerberosProvider(principal, keytabPath string, useTicketCache bool) (*KerberosProvider, error) {
	if principal == "" {
		return nil, fmt.Errorf("kerberos principal is required")
	}

	var krbClient *client.Client

	cfg, err := loadKerberosConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load krb5 config: %w", err)
	}

	if useTicketCache {
		ccachePath := os.Getenv("KRB5CCNAME")
		if ccachePath == "" {
			u, err := user.Current()
			if err != nil {
				return nil, fmt.Errorf("failed to determine current user for default ccache path: %w", err)
			}
			ccachePath = "/tmp/krb5cc_" + u.Uid
		}

		cc, err := credentials.LoadCCache(ccachePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load Kerberos CCache from %s: %w", ccachePath, err)
		}

		krbClient, err = client.NewFromCCache(cc, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kerberos client from CCache: %w", err)
		}
	} else {
		if keytabPath == "" {
			return nil, fmt.Errorf("keytab path is required when use_ticket_cache is false")
		}

		kt, err := keytab.Load(keytabPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load keytab from %s: %w", keytabPath, err)
		}

		realm := extractRealm(principal)
		krbClient = client.NewWithKeytab(principal, realm, kt, cfg)
	}

	return &KerberosProvider{
		principal: principal,
		krbClient: krbClient,
		transport: &spnegoRoundTripper{
			client: krbClient,
		},
	}, nil
}

func (p *KerberosProvider) Header(ctx context.Context) (string, string, error) {
	return "", "", nil
}

func (p *KerberosProvider) WrapTransport(base http.RoundTripper) http.RoundTripper {
	if rt, ok := p.transport.(*spnegoRoundTripper); ok {
		if rt.base == nil {
			rt.base = base
		}
	}
	return p.transport
}

func (p *KerberosProvider) Close() error {
	p.krbClient.Destroy()
	return nil
}

type spnegoRoundTripper struct {
	base   http.RoundTripper
	client *client.Client
}

func (rt *spnegoRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := spnego.SetSPNEGOHeader(rt.client, req, ""); err != nil {
		return nil, fmt.Errorf("kerberos authentication failed: %w", err)
	}
	if rt.base != nil {
		return rt.base.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}

func extractRealm(principal string) string {
	parts := strings.SplitN(principal, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func loadKerberosConfig() (*config.Config, error) {
	if cfgPath := os.Getenv("KRB5_CONFIG"); cfgPath != "" {
		cfg, err := config.Load(cfgPath)
		if err == nil {
			return cfg, nil
		}
	}

	cfg, err := config.Load("/etc/krb5.conf")
	if err == nil {
		return cfg, nil
	}

	return config.New(), nil
}
