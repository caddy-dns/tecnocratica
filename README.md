Tecnocratica module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with Tecnocratica-powered services (Neodigit and Virtualname).

## Caddy module name

```
dns.providers.tecnocratica
```

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

This module supports Tecnocratica-powered DNS services. You'll need an API token from your control panel.

The `api_url` parameter specifies which service endpoint to use:

- **Neodigit**: `https://api.neodigit.net/v1`
- **Virtualname**: `https://api.virtualname.net/v1`

See the [associated README in the libdns package](https://github.com/libdns/tecnocratica) for more information about authentication.