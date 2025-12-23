package tecnocratica

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	libdnstecnocratica "github.com/libdns/tecnocratica"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantToken string
		wantURL   string
		wantErr   bool
	}{
		{
			name:      "with api_token only",
			input:     "tecnocratica {\n  api_token test_token\n}",
			wantToken: "test_token",
			wantURL:   "",
			wantErr:   false,
		},
		{
			name:      "with api_token and api_url",
			input:     "tecnocratica {\n  api_token test_token\n  api_url https://api.neodigit.net/v1\n}",
			wantToken: "test_token",
			wantURL:   "https://api.neodigit.net/v1",
			wantErr:   false,
		},
		{
			name:      "with inline api_token",
			input:     "tecnocratica test_token",
			wantToken: "test_token",
			wantURL:   "",
			wantErr:   false,
		},
		{
			name:      "missing api_token",
			input:     "tecnocratica",
			wantToken: "",
			wantURL:   "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.input)
			p := &Provider{new(libdnstecnocratica.Provider)}
			err := p.UnmarshalCaddyfile(d)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if p.Provider.APIToken != tt.wantToken {
					t.Errorf("APIToken = %v, want %v", p.Provider.APIToken, tt.wantToken)
				}
				if p.Provider.APIURL != tt.wantURL {
					t.Errorf("APIURL = %v, want %v", p.Provider.APIURL, tt.wantURL)
				}
			}
		})
	}
}
