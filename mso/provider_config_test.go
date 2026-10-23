package mso

import (
	"context"
	"reflect"
	"strings"
	"testing"

	providerconfig "github.com/CiscoDevNet/terraform-provider-mso/internal/config"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSDKProviderConfiguration(t *testing.T) {
	defaults := providerconfig.Config{
		Username: "environment-user", Password: "environment-password", URL: "https://environment.test",
		IsInsecure: true, MaxRetries: 2,
	}
	tests := []struct {
		name       string
		env        map[string]string
		attributes map[string]interface{}
		want       providerconfig.Config
		wantError  string
	}{
		{name: "null attributes use defaults", want: defaults},
		{
			name: "null attributes use environment",
			env:  map[string]string{"MSO_INSECURE": "false", "MSO_RETRIES": "4", "MSO_DOMAIN": "environment-domain", "MSO_PROXY_URL": "http://environment-proxy.test:8080", "MSO_PLATFORM": "nd"},
			want: providerconfig.Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: false, MaxRetries: 4, Domain: "environment-domain", ProxyUrl: "http://environment-proxy.test:8080", Platform: "nd"},
		},
		{
			name: "environment still accepts legacy mso platform",
			env:  map[string]string{"MSO_PLATFORM": "mso"},
			want: providerconfig.Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: true, MaxRetries: 2, Platform: "mso"},
		},
		{
			name:       "explicit values override environment including false and zero",
			env:        map[string]string{"MSO_INSECURE": "true", "MSO_RETRIES": "4", "MSO_DOMAIN": "environment-domain", "MSO_PROXY_URL": "http://environment-proxy.test:8080", "MSO_PLATFORM": "mso"},
			attributes: map[string]interface{}{"username": "configured-user", "password": "configured-password", "url": "https://configured.test", "insecure": false, "retries": "0", "domain": "configured-domain", "proxy_url": "http://configured-proxy.test:8080", "platform": "nd"},
			want:       providerconfig.Config{Username: "configured-user", Password: "configured-password", URL: "https://configured.test", IsInsecure: false, MaxRetries: 0, Domain: "configured-domain", ProxyUrl: "http://configured-proxy.test:8080", Platform: "nd"},
		},
		{
			name:       "explicit empty optional strings override environment",
			env:        map[string]string{"MSO_RETRIES": "4", "MSO_DOMAIN": "environment-domain", "MSO_PROXY_URL": "http://environment-proxy.test:8080"},
			attributes: map[string]interface{}{"retries": "", "domain": "", "proxy_url": ""},
			want:       defaults,
		},
		{name: "invalid insecure environment", env: map[string]string{"MSO_INSECURE": "invalid"}, wantError: "invalid value for MSO_INSECURE"},
		{name: "invalid retries environment", env: map[string]string{"MSO_RETRIES": "invalid"}, wantError: "invalid value for retries"},
		{name: "invalid platform environment", env: map[string]string{"MSO_PLATFORM": "invalid"}, wantError: "invalid value for platform"},
		{
			name:       "explicit values override invalid environment",
			env:        map[string]string{"MSO_INSECURE": "invalid", "MSO_RETRIES": "invalid", "MSO_PLATFORM": "invalid"},
			attributes: map[string]interface{}{"insecure": false, "retries": "0", "platform": "mso"},
			want:       providerconfig.Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: false, MaxRetries: 0, Platform: "mso"},
		},
		{name: "explicit empty username does not fall back", attributes: map[string]interface{}{"username": ""}, wantError: "Username must be provided"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for name, value := range map[string]string{
				"MSO_USERNAME": defaults.Username, "MSO_PASSWORD": defaults.Password, "MSO_URL": defaults.URL,
				"MSO_INSECURE": "", "MSO_RETRIES": "", "MSO_DOMAIN": "", "MSO_PROXY_URL": "", "MSO_PLATFORM": "",
			} {
				t.Setenv(name, value)
			}
			for name, value := range test.env {
				t.Setenv(name, value)
			}

			var got providerconfig.Config
			var called bool
			p := Provider()
			configure := p.ConfigureContextFunc
			p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
				called = true
				got, _ = configFromSDKResourceData(data)
				return configure(ctx, data)
			}
			server := p.GRPCProvider()
			ctx := context.Background()
			providerSchema, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
			if err != nil || len(providerSchema.Diagnostics) != 0 {
				t.Fatalf("getting SDK provider schema: response=%v, error=%v", providerSchema, err)
			}
			configType := providerSchema.Provider.ValueType().(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(configType.AttributeTypes))
			for name, attrType := range configType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attrType, test.attributes[name])
			}
			config, err := tfprotov5.NewDynamicValue(configType, tftypes.NewValue(configType, attributes))
			if err != nil {
				t.Fatalf("encoding provider configuration: %v", err)
			}
			prepared, err := server.PrepareProviderConfig(ctx, &tfprotov5.PrepareProviderConfigRequest{Config: &config})
			if err != nil {
				t.Fatalf("preparing provider configuration: %v", err)
			}
			wantWarnings := 0
			if !attributes["platform"].IsNull() {
				wantWarnings = 1
			}
			if len(prepared.Diagnostics) != wantWarnings {
				t.Fatalf("preparation diagnostics = %v, want %d platform deprecation warnings", prepared.Diagnostics, wantWarnings)
			}
			for _, diagnostic := range prepared.Diagnostics {
				if diagnostic.Severity != tfprotov5.DiagnosticSeverityWarning ||
					diagnostic.Summary != "Argument is deprecated" ||
					diagnostic.Detail != providerconfig.PlatformDeprecationMessage ||
					!diagnostic.Attribute.Equal(tftypes.NewAttributePath().WithAttributeName("platform")) {
					t.Fatalf("unexpected preparation diagnostic: %v", diagnostic)
				}
			}
			response, err := server.ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{Config: prepared.PreparedConfig, TerraformVersion: "1.1.7"})
			if err != nil {
				t.Fatalf("configuring SDK provider: %v", err)
			}
			if !called {
				t.Fatal("SDK protocol server did not call provider configuration")
			}
			wantEnvironmentWarnings := 0
			if attributes["platform"].IsNull() && test.env["MSO_PLATFORM"] != "" {
				wantEnvironmentWarnings = 1
			}
			errors, warnings := 0, 0
			for _, diagnostic := range response.Diagnostics {
				switch diagnostic.Severity {
				case tfprotov5.DiagnosticSeverityWarning:
					warnings++
					if diagnostic.Summary != "Environment variable is deprecated" || diagnostic.Detail != providerconfig.PlatformEnvironmentDeprecationMessage {
						t.Fatalf("unexpected configuration warning: %v", diagnostic)
					}
				case tfprotov5.DiagnosticSeverityError:
					errors++
					if test.wantError == "" || !strings.Contains(diagnostic.Summary+diagnostic.Detail, test.wantError) {
						t.Fatalf("unexpected configuration error: %v", diagnostic)
					}
				default:
					t.Fatalf("unexpected configuration diagnostic: %v", diagnostic)
				}
			}
			if warnings != wantEnvironmentWarnings {
				t.Fatalf("environment deprecation warnings = %d, want %d", warnings, wantEnvironmentWarnings)
			}
			if test.wantError != "" {
				if errors != 1 {
					t.Fatalf("configuration errors = %d, want one containing %q", errors, test.wantError)
				}
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("resolved configuration mismatch: got %#v, want %#v", got, test.want)
			}
		})
	}
}
