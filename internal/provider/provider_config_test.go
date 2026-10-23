package provider

import (
	"context"
	"reflect"
	"strings"
	"testing"

	providerconfig "github.com/CiscoDevNet/terraform-provider-mso/internal/config"
	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/ciscoecosystem/mso-go-client/client"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
)

// Capture the resolved settings while still exercising the production Configure
// method through the Framework protocol server, including Terraform null values.
type configurationTestProvider struct {
	*MSOProvider
	config         Config
	called         bool
	resourceData   interface{}
	dataSourceData interface{}
}

func (p *configurationTestProvider) Configure(ctx context.Context, request frameworkprovider.ConfigureRequest, response *frameworkprovider.ConfigureResponse) {
	p.called = true
	var model MSOProviderModel
	if diagnostics := request.Config.Get(ctx, &model); !diagnostics.HasError() {
		p.config, _ = model.Config()
	}
	p.MSOProvider.Configure(ctx, request, response)
	p.resourceData = response.ResourceData
	p.dataSourceData = response.DataSourceData
}

func TestFrameworkProviderConfiguration(t *testing.T) {
	defaults := Config{
		Username: "environment-user", Password: "environment-password", URL: "https://environment.test",
		IsInsecure: true, MaxRetries: 2,
	}
	tests := []struct {
		name       string
		env        map[string]string
		attributes map[string]interface{}
		want       Config
		wantError  string
	}{
		{name: "null attributes use defaults", want: defaults},
		{
			name: "null attributes use environment",
			env:  map[string]string{"MSO_INSECURE": "false", "MSO_RETRIES": "4", "MSO_DOMAIN": "environment-domain", "MSO_PROXY_URL": "http://environment-proxy.test:8080", "MSO_PLATFORM": "nd"},
			want: Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: false, MaxRetries: 4, Domain: "environment-domain", ProxyUrl: "http://environment-proxy.test:8080", Platform: "nd"},
		},
		{
			name: "environment still accepts legacy mso platform",
			env:  map[string]string{"MSO_PLATFORM": "mso"},
			want: Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: true, MaxRetries: 2, Platform: "mso"},
		},
		{
			name:       "explicit values override environment including false and zero",
			env:        map[string]string{"MSO_INSECURE": "true", "MSO_RETRIES": "4", "MSO_DOMAIN": "environment-domain", "MSO_PROXY_URL": "http://environment-proxy.test:8080", "MSO_PLATFORM": "mso"},
			attributes: map[string]interface{}{"username": "configured-user", "password": "configured-password", "url": "https://configured.test", "insecure": false, "retries": "0", "domain": "configured-domain", "proxy_url": "http://configured-proxy.test:8080", "platform": "nd"},
			want:       Config{Username: "configured-user", Password: "configured-password", URL: "https://configured.test", IsInsecure: false, MaxRetries: 0, Domain: "configured-domain", ProxyUrl: "http://configured-proxy.test:8080", Platform: "nd"},
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
			want:       Config{Username: defaults.Username, Password: defaults.Password, URL: defaults.URL, IsInsecure: false, MaxRetries: 0, Platform: "mso"},
		},
		{name: "explicit empty username does not fall back", attributes: map[string]interface{}{"username": ""}, wantError: "Username must be provided"},
		{name: "unknown insecure is rejected", attributes: map[string]interface{}{"insecure": tftypes.UnknownValue}, wantError: "insecure must be known"},
		{name: "unknown username is rejected", attributes: map[string]interface{}{"username": tftypes.UnknownValue}, wantError: "MSO_USERNAME must be known"},
		{name: "unknown platform does not use environment", env: map[string]string{"MSO_PLATFORM": "nd"}, attributes: map[string]interface{}{"platform": tftypes.UnknownValue}, wantError: "MSO_PLATFORM must be known"},
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

			p := &configurationTestProvider{MSOProvider: &MSOProvider{version: "test"}}
			server := providerserver.NewProtocol6(p)()
			ctx := context.Background()
			providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
			if err != nil || len(providerSchema.Diagnostics) != 0 {
				t.Fatalf("getting Framework provider schema: response=%v, error=%v", providerSchema, err)
			}
			configType := providerSchema.Provider.ValueType().(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(configType.AttributeTypes))
			for name, attrType := range configType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attrType, test.attributes[name])
			}
			config, err := tfprotov6.NewDynamicValue(configType, tftypes.NewValue(configType, attributes))
			if err != nil {
				t.Fatalf("encoding provider configuration: %v", err)
			}
			response, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{Config: &config, TerraformVersion: "1.1.7"})
			if err != nil {
				t.Fatalf("configuring Framework provider: %v", err)
			}
			if !p.called {
				t.Fatal("Framework protocol server did not call provider configuration")
			}
			wantEnvironmentWarnings := 0
			if attributes["platform"].IsNull() && test.env["MSO_PLATFORM"] != "" {
				wantEnvironmentWarnings = 1
			}
			errors, warnings := 0, 0
			for _, diagnostic := range response.Diagnostics {
				switch diagnostic.Severity {
				case tfprotov6.DiagnosticSeverityWarning:
					warnings++
					if diagnostic.Summary != "Environment variable is deprecated" || diagnostic.Detail != providerconfig.PlatformEnvironmentDeprecationMessage {
						t.Fatalf("unexpected configuration warning: %v", diagnostic)
					}
				case tfprotov6.DiagnosticSeverityError:
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
				if p.resourceData != nil || p.dataSourceData != nil {
					t.Fatal("invalid configuration returned provider data")
				}
				return
			}
			if !reflect.DeepEqual(p.config, test.want) {
				t.Fatalf("resolved configuration mismatch: got %#v, want %#v", p.config, test.want)
			}
			apiClient, ok := p.resourceData.(*client.Client)
			if !ok || apiClient == nil || p.dataSourceData != apiClient {
				t.Fatal("successful configuration must share its API client with resources and data sources")
			}
		})
	}
}

func TestProviderPlatformValidation(t *testing.T) {
	ctx := context.Background()
	servers := []struct {
		name                string
		deprecationWarnings int
		new                 func() (tfprotov6.ProviderServer, error)
	}{
		{
			name:                "SDK upgraded to protocol 6",
			deprecationWarnings: 1,
			new: func() (tfprotov6.ProviderServer, error) {
				return tf5to6server.UpgradeServer(ctx, mso.Provider().GRPCProvider)
			},
		},
		{
			name:                "Framework protocol 6",
			deprecationWarnings: 1,
			new: func() (tfprotov6.ProviderServer, error) {
				return providerserver.NewProtocol6(New("test")())(), nil
			},
		},
		{
			name:                "production mux",
			deprecationWarnings: 2,
			new: func() (tfprotov6.ProviderServer, error) {
				factory, err := MuxProviderServer(ctx, "test", mso.Provider().GRPCProvider)
				if err != nil {
					return nil, err
				}
				return factory(), nil
			},
		},
	}
	tests := []struct {
		name       string
		attributes map[string]interface{}
		wantError  bool
	}{
		{name: "omitted"},
		{name: "null", attributes: map[string]interface{}{"platform": nil}},
		{name: "unknown", attributes: map[string]interface{}{"platform": tftypes.UnknownValue}},
		{name: "mso", attributes: map[string]interface{}{"platform": "mso"}},
		{name: "nd", attributes: map[string]interface{}{"platform": "nd"}},
		{name: "invalid", attributes: map[string]interface{}{"platform": "invalid"}, wantError: true},
		{name: "empty", attributes: map[string]interface{}{"platform": ""}, wantError: true},
		{name: "case mismatch", attributes: map[string]interface{}{"platform": "MSO"}, wantError: true},
	}

	// Schema validation applies to Terraform values. Environment fallback is
	// resolved and validated later by Configure, including for null attributes.
	t.Setenv("MSO_PLATFORM", "invalid")
	for _, serverTest := range servers {
		t.Run(serverTest.name, func(t *testing.T) {
			server, err := serverTest.new()
			if err != nil {
				t.Fatalf("constructing provider server: %v", err)
			}
			providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
			if err != nil || len(providerSchema.Diagnostics) != 0 {
				t.Fatalf("getting provider schema: response=%v, error=%v", providerSchema, err)
			}
			var platformSchema *tfprotov6.SchemaAttribute
			for _, attribute := range providerSchema.Provider.Block.Attributes {
				if attribute.Name == "platform" {
					platformSchema = attribute
					break
				}
			}
			if platformSchema == nil || !platformSchema.Deprecated || platformSchema.DeprecationMessage != providerconfig.PlatformDeprecationMessage {
				t.Fatalf("platform schema must advertise its shared deprecation message: %v", platformSchema)
			}
			configType := providerSchema.Provider.ValueType().(tftypes.Object)
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					attributes := make(map[string]tftypes.Value, len(configType.AttributeTypes))
					for name, attrType := range configType.AttributeTypes {
						attributes[name] = tftypes.NewValue(attrType, test.attributes[name])
					}
					config, err := tfprotov6.NewDynamicValue(configType, tftypes.NewValue(configType, attributes))
					if err != nil {
						t.Fatalf("encoding provider configuration: %v", err)
					}
					response, err := server.ValidateProviderConfig(ctx, &tfprotov6.ValidateProviderConfigRequest{Config: &config})
					if err != nil {
						t.Fatalf("validating provider configuration: %v", err)
					}
					errors, warnings := 0, 0
					for _, diagnostic := range response.Diagnostics {
						if !diagnostic.Attribute.Equal(tftypes.NewAttributePath().WithAttributeName("platform")) {
							t.Errorf("validation diagnostic must identify platform: %v", diagnostic)
						}
						switch diagnostic.Severity {
						case tfprotov6.DiagnosticSeverityError:
							errors++
						case tfprotov6.DiagnosticSeverityWarning:
							warnings++
							if diagnostic.Detail != providerconfig.PlatformDeprecationMessage ||
								(diagnostic.Summary != "Argument is deprecated" && diagnostic.Summary != "Attribute Deprecated") {
								t.Errorf("unexpected validation warning: %v", diagnostic)
							}
						default:
							t.Errorf("unexpected validation diagnostic severity: %v", diagnostic)
						}
					}
					if (errors > 0) != test.wantError {
						t.Fatalf("validation errors = %d, want error = %t; diagnostics: %v", errors, test.wantError, response.Diagnostics)
					}
					if attributes["platform"].IsKnown() && !test.wantError {
						wantWarnings := 0
						if !attributes["platform"].IsNull() {
							wantWarnings = serverTest.deprecationWarnings
						}
						if warnings != wantWarnings {
							t.Fatalf("platform deprecation warnings = %d, want %d", warnings, wantWarnings)
						}
					}
				})
			}
		})
	}
}

func TestMuxProviderPlatformEnvironmentDeprecation(t *testing.T) {
	for _, name := range []string{"MSO_INSECURE", "MSO_RETRIES", "MSO_DOMAIN", "MSO_PROXY_URL"} {
		t.Setenv(name, "")
	}
	tests := []struct {
		name         string
		environment  string
		platform     interface{}
		wantWarnings int
		wantError    bool
	}{
		{name: "environment unset"},
		{name: "legacy mso environment", environment: "mso", wantWarnings: 2},
		{name: "nd environment", environment: "nd", wantWarnings: 2},
		{name: "explicit mso overrides nd environment", environment: "nd", platform: "mso"},
		{name: "explicit nd overrides invalid environment", environment: "invalid", platform: "nd"},
		// The mux stops when its first provider returns an error, so invalid
		// environment configuration returns that provider's warning and error.
		{name: "invalid environment still fails", environment: "invalid", wantWarnings: 1, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("MSO_PLATFORM", test.environment)
			ctx := context.Background()
			factory, err := MuxProviderServer(ctx, "test", mso.Provider().GRPCProvider)
			if err != nil {
				t.Fatalf("constructing provider mux: %v", err)
			}
			server := factory()
			providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
			if err != nil || len(providerSchema.Diagnostics) != 0 {
				t.Fatalf("getting provider schema: response=%v, error=%v", providerSchema, err)
			}
			configType := providerSchema.Provider.ValueType().(tftypes.Object)
			values := map[string]interface{}{
				"username": "configured-user", "password": "configured-password", "url": "https://configured.test", "platform": test.platform,
			}
			attributes := make(map[string]tftypes.Value, len(configType.AttributeTypes))
			for name, attrType := range configType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attrType, values[name])
			}
			config, err := tfprotov6.NewDynamicValue(configType, tftypes.NewValue(configType, attributes))
			if err != nil {
				t.Fatalf("encoding provider configuration: %v", err)
			}
			response, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{Config: &config, TerraformVersion: "1.1.7"})
			if err != nil {
				t.Fatalf("configuring provider mux: %v", err)
			}
			errors, warnings := 0, 0
			for _, diagnostic := range response.Diagnostics {
				switch diagnostic.Severity {
				case tfprotov6.DiagnosticSeverityWarning:
					warnings++
					if diagnostic.Summary != "Environment variable is deprecated" || diagnostic.Detail != providerconfig.PlatformEnvironmentDeprecationMessage {
						t.Fatalf("unexpected configuration warning: %v", diagnostic)
					}
				case tfprotov6.DiagnosticSeverityError:
					errors++
					if !test.wantError || !strings.Contains(diagnostic.Summary+diagnostic.Detail, "invalid value for platform") {
						t.Fatalf("unexpected configuration error: %v", diagnostic)
					}
				default:
					t.Fatalf("unexpected configuration diagnostic: %v", diagnostic)
				}
			}
			if warnings != test.wantWarnings || (errors > 0) != test.wantError {
				t.Fatalf("got %d warnings and %d errors, want %d warnings and error=%t", warnings, errors, test.wantWarnings, test.wantError)
			}
		})
	}
}
