package tecnocratica

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"github.com/libdns/tecnocratica"
	"go.uber.org/zap"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct {
	*tecnocratica.Provider
	logger *zap.Logger
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.tecnocratica",
		New: func() caddy.Module { return &Provider{Provider: new(tecnocratica.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger()
	repl := caddy.NewReplacer()
	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.APIURL = repl.ReplaceAll(p.Provider.APIURL, "")

	p.logger.Info("tecnocratica DNS provider provisioned",
		zap.String("api_url", p.Provider.APIURL),
		zap.Bool("has_token", p.Provider.APIToken != ""))

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

// AppendRecords appends records to the zone. Logs the operation for debugging.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.logger.Info("attempting to append DNS records",
		zap.String("zone", zone),
		zap.Int("record_count", len(records)))

	for i, rec := range records {
		rr := rec.RR()
		p.logger.Debug("record to append",
			zap.Int("index", i),
			zap.String("name", rr.Name),
			zap.String("type", rr.Type),
			zap.String("value", rr.Data))
	}

	result, err := p.Provider.AppendRecords(ctx, zone, records)
	if err != nil {
		p.logger.Error("failed to append DNS records",
			zap.String("zone", zone),
			zap.Error(err))
		return nil, err
	}

	p.logger.Info("successfully appended DNS records",
		zap.String("zone", zone),
		zap.Int("appended_count", len(result)))

	return result, nil
}

// GetRecords gets records from the zone. Logs the operation for debugging.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.logger.Debug("getting DNS records", zap.String("zone", zone))

	result, err := p.Provider.GetRecords(ctx, zone)
	if err != nil {
		p.logger.Error("failed to get DNS records",
			zap.String("zone", zone),
			zap.Error(err))
		return nil, err
	}

	p.logger.Debug("successfully got DNS records",
		zap.String("zone", zone),
		zap.Int("record_count", len(result)))

	return result, nil
}

// SetRecords sets records in the zone. Logs the operation for debugging.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.logger.Info("attempting to set DNS records",
		zap.String("zone", zone),
		zap.Int("record_count", len(records)))

	result, err := p.Provider.SetRecords(ctx, zone, records)
	if err != nil {
		p.logger.Error("failed to set DNS records",
			zap.String("zone", zone),
			zap.Error(err))
		return nil, err
	}

	p.logger.Info("successfully set DNS records",
		zap.String("zone", zone),
		zap.Int("set_count", len(result)))

	return result, nil
}

// DeleteRecords deletes records from the zone. Logs the operation for debugging.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.logger.Info("attempting to delete DNS records",
		zap.String("zone", zone),
		zap.Int("record_count", len(records)))

	result, err := p.Provider.DeleteRecords(ctx, zone, records)
	if err != nil {
		p.logger.Error("failed to delete DNS records",
			zap.String("zone", zone),
			zap.Error(err))
		return nil, err
	}

	p.logger.Info("successfully deleted DNS records",
		zap.String("zone", zone),
		zap.Int("deleted_count", len(result)))

	return result, nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
