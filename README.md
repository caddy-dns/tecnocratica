Tecnocratica module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with Tecnocratica (Neodigit and Virtualname).

## Caddy module name

```
dns.providers.tecnocratica
```

## Building

**Note:** The libdns package for Tecnocratica is currently in development. Until it is published at `github.com/libdns/tecnocratica`, this module uses a replace directive in `go.mod` to point to the development repository at `github.com/aalmenar/libdns-tecnocratica`. Once the official package is published, the replace directive can be removed.

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "tecnocratica",
				"api_token": "YOUR_PROVIDER_API_TOKEN",
				"api_url": "https://api.neodigit.net/v1"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns tecnocratica {
		api_token <token>
		api_url https://api.neodigit.net/v1
	}
}
```

```
# one site
tls {
	dns tecnocratica {
		api_token <token>
		api_url https://api.neodigit.net/v1
	}
}
```

You can also use environment variable placeholders:

```
tls {
	dns tecnocratica {
		api_token {env.TECNOCRATICA_API_TOKEN}
		api_url {env.TECNOCRATICA_API_URL}
	}
}
```

## Authenticating

This module uses the Tecnocratica API. You can obtain an API token from your control panel:

- **Neodigit**: `https://api.neodigit.net/v1`
- **Virtualname**: `https://api.virtualname.net/v1`

See the [associated README in the libdns package](https://github.com/libdns/tecnocratica) for more information.