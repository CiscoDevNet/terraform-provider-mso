package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"

	providerconfig "github.com/CiscoDevNet/terraform-provider-mso/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*MSOProvider)(nil)

// MSOProvider hosts new resources while the existing provider remains
// implemented by SDKv2 behind the mux server.
type MSOProvider struct {
	version string
}

type MSOProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	URL      types.String `tfsdk:"url"`
	Insecure types.Bool   `tfsdk:"insecure"`
	Domain   types.String `tfsdk:"domain"`
	ProxyURL types.String `tfsdk:"proxy_url"`
	Platform types.String `tfsdk:"platform"`
	Retries  types.String `tfsdk:"retries"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MSOProvider{version: version}
	}
}

func (p *MSOProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "mso"
	response.Version = p.version
}

func (p *MSOProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = providerschema.Schema{
		Attributes: map[string]providerschema.Attribute{
			"username": providerschema.StringAttribute{
				Optional:    true,
				Description: "Username for the MSO Account",
			},
			"password": providerschema.StringAttribute{
				Optional:    true,
				Description: "Password for the MSO Account",
			},
			"url": providerschema.StringAttribute{
				Optional:    true,
				Description: "URL of the Cisco MSO web interface",
			},
			"insecure": providerschema.BoolAttribute{
				Optional:    true,
				Description: "Allow insecure HTTPS client",
			},
			"domain": providerschema.StringAttribute{
				Optional:    true,
				Description: "Domain name for remote user authentication",
			},
			"proxy_url": providerschema.StringAttribute{
				Optional:    true,
				Description: "Proxy Server URL with port number",
			},
			"platform": providerschema.StringAttribute{
				Optional:    true,
				Description: "Parameter that specifies where MSO is installed",
			},
			"retries": providerschema.StringAttribute{
				Optional:    true,
				Description: "Number of retries for REST API calls. Defaults to 2.",
			},
		},
	}
}

func (p *MSOProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data MSOProviderModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	config, err := data.Config()
	if err != nil {
		response.Diagnostics.AddError("Unable to configure MSO provider", err.Error())
		return
	}

	client := config.GetClient()
	response.ResourceData = client
	response.DataSourceData = client
}

func (p *MSOProvider) Resources(context.Context) []func() resource.Resource {
	return providerResources()
}

func (p *MSOProvider) DataSources(context.Context) []func() datasource.DataSource {
	return providerDataSources()
}

func (m MSOProviderModel) Config() (Config, error) {
	username, err := stringValue(m.Username, "MSO_USERNAME")
	if err != nil {
		return Config{}, err
	}
	password, err := stringValue(m.Password, "MSO_PASSWORD")
	if err != nil {
		return Config{}, err
	}
	url, err := stringValue(m.URL, "MSO_URL")
	if err != nil {
		return Config{}, err
	}
	domain, err := stringValue(m.Domain, "MSO_DOMAIN")
	if err != nil {
		return Config{}, err
	}
	proxyURL, err := stringValue(m.ProxyURL, "MSO_PROXY_URL")
	if err != nil {
		return Config{}, err
	}
	platform, err := stringValue(m.Platform, "MSO_PLATFORM")
	if err != nil {
		return Config{}, err
	}
	retries, err := stringValue(m.Retries, "MSO_RETRIES")
	if err != nil {
		return Config{}, err
	}

	isInsecure := true
	if m.Insecure.IsUnknown() {
		return Config{}, fmt.Errorf("insecure must be known during provider configuration")
	}
	if !m.Insecure.IsNull() {
		isInsecure = m.Insecure.ValueBool()
	} else if value := os.Getenv("MSO_INSECURE"); value != "" {
		isInsecure, err = strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid value for MSO_INSECURE")
		}
	}

	return BuildConfig(username, password, url, isInsecure, proxyURL, domain, platform, retries)
}

func stringValue(value types.String, environmentVariable string) (string, error) {
	if value.IsUnknown() {
		return "", fmt.Errorf("provider attribute for %s must be known during provider configuration", environmentVariable)
	}
	if !value.IsNull() {
		return value.ValueString(), nil
	}

	return os.Getenv(environmentVariable), nil
}

// Config and BuildConfig remain available from the provider package while
// their implementation is shared with the legacy SDKv2 provider.
type Config = providerconfig.Config

func BuildConfig(username, password, url string, isInsecure bool, proxyURL, domain, platform, retries string) (Config, error) {
	return providerconfig.BuildConfig(username, password, url, isInsecure, proxyURL, domain, platform, retries)
}
