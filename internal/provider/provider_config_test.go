package provider

import (
	"context"
	"reflect"
	"strings"
	"testing"

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
			if test.wantError != "" {
				for _, diagnostic := range response.Diagnostics {
					if diagnostic.Severity == tfprotov6.DiagnosticSeverityError && strings.Contains(diagnostic.Summary+diagnostic.Detail, test.wantError) {
						if p.resourceData != nil || p.dataSourceData != nil {
							t.Fatal("invalid configuration returned provider data")
						}
						return
					}
				}
				t.Fatalf("expected error containing %q, got %v", test.wantError, response.Diagnostics)
			}
			if len(response.Diagnostics) != 0 {
				t.Fatalf("unexpected configuration diagnostics: %v", response.Diagnostics)
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
		name string
		new  func() (tfprotov6.ProviderServer, error)
	}{
		{
			name: "SDK upgraded to protocol 6",
			new: func() (tfprotov6.ProviderServer, error) {
				return tf5to6server.UpgradeServer(ctx, mso.Provider().GRPCProvider)
			},
		},
		{
			name: "Framework protocol 6",
			new: func() (tfprotov6.ProviderServer, error) {
				return providerserver.NewProtocol6(New("test")())(), nil
			},
		},
		{
			name: "production mux",
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
					errors := 0
					for _, diagnostic := range response.Diagnostics {
						if diagnostic.Severity != tfprotov6.DiagnosticSeverityError {
							continue
						}
						errors++
						if !diagnostic.Attribute.Equal(tftypes.NewAttributePath().WithAttributeName("platform")) {
							t.Errorf("validation error must identify platform: %v", diagnostic)
						}
					}
					if (errors > 0) != test.wantError {
						t.Fatalf("validation errors = %d, want error = %t; diagnostics: %v", errors, test.wantError, response.Diagnostics)
					}
				})
			}
		})
	}
}
