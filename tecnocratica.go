package tecnocratica

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/tecnocratica"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *tecnocratica.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.tecnocratica",
		New: func() caddy.Module { return &Provider{new(tecnocratica.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.APIURL = repl.ReplaceAll(p.Provider.APIURL, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	tecnocratica [<api_token>] {
//	    api_token <api_token>
//	    api_url <api_url>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	if d.NextArg() {
		p.Provider.APIToken = d.Val()
	}

	for d.NextBlock(0) {
		switch d.Val() {
		case "api_token":
			if p.Provider.APIToken != "" {
				return d.Err("API token already set")
			}
			if d.NextArg() {
				p.Provider.APIToken = d.Val()
			} else {
				return d.ArgErr()
			}
		case "api_url":
			if d.NextArg() {
				p.Provider.APIURL = d.Val()
			} else {
				return d.ArgErr()
			}
		default:
			return d.Errf("unrecognized subdirective '%s'", d.Val())
		}
	}

	if d.NextArg() {
		return d.Errf("unexpected argument '%s'", d.Val())
	}

	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
