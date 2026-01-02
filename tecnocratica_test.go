package tecnocratica

import (
	"context"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
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
			p := &Provider{Provider: new(libdnstecnocratica.Provider)}
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

// TestUnmarshalCaddyfileEdgeCases tests error handling and edge cases in Caddyfile parsing
func TestUnmarshalCaddyfileEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{
			name:        "duplicate api_token in block and inline",
			input:       "tecnocratica inline_token {\n  api_token block_token\n}",
			wantErr:     true,
			errContains: "API token already set",
		},
		{
			name:        "duplicate api_token in block",
			input:       "tecnocratica {\n  api_token token1\n  api_token token2\n}",
			wantErr:     true,
			errContains: "API token already set",
		},
		{
			name:        "unrecognized subdirective",
			input:       "tecnocratica {\n  api_token test_token\n  invalid_directive value\n}",
			wantErr:     true,
			errContains: "unrecognized subdirective",
		},
		{
			name:        "api_token without value",
			input:       "tecnocratica {\n  api_token\n}",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "api_url without value",
			input:       "tecnocratica {\n  api_token test_token\n  api_url\n}",
			wantErr:     true,
			errContains: "",
		},
		{
			name:    "environment variable placeholder (not expanded during unmarshal)",
			input:   "tecnocratica {\n  api_token {env.TECNOCRATICA_TOKEN}\n  api_url {env.TECNOCRATICA_URL}\n}",
			wantErr: false,
		},
		{
			name:        "too many inline arguments",
			input:       "tecnocratica token1 token2",
			wantErr:     true,
			errContains: "unexpected argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.input)
			p := &Provider{Provider: new(libdnstecnocratica.Provider)}
			err := p.UnmarshalCaddyfile(d)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("UnmarshalCaddyfile() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

// TestProvision tests the Provision method with placeholder replacement
func TestProvision(t *testing.T) {
	tests := []struct {
		name          string
		apiToken      string
		apiURL        string
		envVars       map[string]string
		expectedToken string
		expectedURL   string
	}{
		{
			name:          "plain strings without placeholders",
			apiToken:      "plain_token",
			apiURL:        "https://api.neodigit.net/v1",
			envVars:       nil,
			expectedToken: "plain_token",
			expectedURL:   "https://api.neodigit.net/v1",
		},
		{
			name:          "with environment variable placeholders",
			apiToken:      "{env.TEST_TOKEN}",
			apiURL:        "{env.TEST_URL}",
			envVars:       map[string]string{"TEST_TOKEN": "env_token_value", "TEST_URL": "https://env.url"},
			expectedToken: "env_token_value",
			expectedURL:   "https://env.url",
		},
		{
			name:          "mixed placeholders and plain text",
			apiToken:      "prefix_{env.TOKEN_SUFFIX}",
			apiURL:        "https://api.neodigit.net/v1",
			envVars:       map[string]string{"TOKEN_SUFFIX": "suffix_value"},
			expectedToken: "prefix_suffix_value",
			expectedURL:   "https://api.neodigit.net/v1",
		},
		{
			name:          "empty api_url is allowed",
			apiToken:      "test_token",
			apiURL:        "",
			envVars:       nil,
			expectedToken: "test_token",
			expectedURL:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables if provided
			if tt.envVars != nil {
				for k, v := range tt.envVars {
					t.Setenv(k, v)
				}
			}

			p := &Provider{Provider: &libdnstecnocratica.Provider{
				APIToken: tt.apiToken,
				APIURL:   tt.apiURL,
			}}

			ctx, cancel := caddy.NewContext(caddy.Context{Context: context.Background()})
			defer cancel()

			err := p.Provision(ctx)
			if err != nil {
				t.Errorf("Provision() unexpected error = %v", err)
				return
			}

			if p.Provider.APIToken != tt.expectedToken {
				t.Errorf("APIToken after Provision() = %v, want %v", p.Provider.APIToken, tt.expectedToken)
			}

			if p.Provider.APIURL != tt.expectedURL {
				t.Errorf("APIURL after Provision() = %v, want %v", p.Provider.APIURL, tt.expectedURL)
			}
		})
	}
}

// TestInterfaceGuards verifies that Provider implements required interfaces
func TestInterfaceGuards(t *testing.T) {
	var p any = (*Provider)(nil)

	// Verify Caddy interfaces
	if _, ok := p.(caddy.Provisioner); !ok {
		t.Error("Provider does not implement caddy.Provisioner")
	}

	if _, ok := p.(caddyfile.Unmarshaler); !ok {
		t.Error("Provider does not implement caddyfile.Unmarshaler")
	}

	// Verify libdns interfaces are accessible through embedded Provider
	provider := &Provider{Provider: new(libdnstecnocratica.Provider)}

	if _, ok := any(provider.Provider).(libdns.RecordGetter); !ok {
		t.Error("Embedded Provider does not implement libdns.RecordGetter")
	}

	if _, ok := any(provider.Provider).(libdns.RecordAppender); !ok {
		t.Error("Embedded Provider does not implement libdns.RecordAppender")
	}

	if _, ok := any(provider.Provider).(libdns.RecordSetter); !ok {
		t.Error("Embedded Provider does not implement libdns.RecordSetter")
	}

	if _, ok := any(provider.Provider).(libdns.RecordDeleter); !ok {
		t.Error("Embedded Provider does not implement libdns.RecordDeleter")
	}
}

// TestCaddyModule verifies module registration and metadata
func TestCaddyModule(t *testing.T) {
	p := Provider{}
	info := p.CaddyModule()

	expectedID := "dns.providers.tecnocratica"
	if info.ID != caddy.ModuleID(expectedID) {
		t.Errorf("CaddyModule().ID = %v, want %v", info.ID, expectedID)
	}

	// Verify New() returns correct type
	instance := info.New()
	if _, ok := instance.(*Provider); !ok {
		t.Errorf("CaddyModule().New() returned %T, want *Provider", instance)
	}

	// Verify the new instance has initialized embedded provider
	if newProvider, ok := instance.(*Provider); ok {
		if newProvider.Provider == nil {
			t.Error("CaddyModule().New() returned Provider with nil embedded Provider")
		}
	}
}

// TestProviderIntegration tests end-to-end workflow
func TestProviderIntegration(t *testing.T) {
	// Set up environment variables
	t.Setenv("TEST_API_TOKEN", "integration_test_token")
	t.Setenv("TEST_API_URL", "https://api.test.example.com/v1")

	// Parse Caddyfile configuration
	input := `tecnocratica {
		api_token {env.TEST_API_TOKEN}
		api_url {env.TEST_API_URL}
	}`

	d := caddyfile.NewTestDispenser(input)
	p := &Provider{Provider: new(libdnstecnocratica.Provider)}

	// Unmarshal configuration
	err := p.UnmarshalCaddyfile(d)
	if err != nil {
		t.Fatalf("UnmarshalCaddyfile() error = %v", err)
	}

	// Verify placeholders are not yet expanded
	if p.Provider.APIToken != "{env.TEST_API_TOKEN}" {
		t.Errorf("APIToken before Provision() = %v, want {env.TEST_API_TOKEN}", p.Provider.APIToken)
	}

	// Provision the provider
	ctx, cancel := caddy.NewContext(caddy.Context{Context: context.Background()})
	defer cancel()

	err = p.Provision(ctx)
	if err != nil {
		t.Fatalf("Provision() error = %v", err)
	}

	// Verify placeholders are expanded
	if p.Provider.APIToken != "integration_test_token" {
		t.Errorf("APIToken after Provision() = %v, want integration_test_token", p.Provider.APIToken)
	}

	if p.Provider.APIURL != "https://api.test.example.com/v1" {
		t.Errorf("APIURL after Provision() = %v, want https://api.test.example.com/v1", p.Provider.APIURL)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
