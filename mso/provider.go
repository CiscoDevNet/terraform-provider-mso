package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_USERNAME", nil),
				Description: "Username for the MSO Account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_PASSWORD", nil),
				Description: "Password for the MSO Account",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_URL", nil),
				Description: "URL of the Cisco MSO web interface",
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure HTTPS client",
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_PROXY_URL", nil),
				Description: "Proxy Server URL with port number",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mso_schema":      resourceMSOSchema(),
			"mso_schema_site": resourceMSOSchemaSite(),
			"mso_site":        resourceMSOSite(),
			"mso_role":        resourceMSORole(),
			"mso_user":        resourceMSOUser(),
			"mso_label":       resourceMSOLabel(),
			"mso_tenant":      resourceMSOTenant(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mso_schema":      datasourceMSOSchema(),
			"mso_schema_site": datasourceMSOSchemaSite(),
			"mso_site":        datasourceMSOSite(),
			"mso_role":        datasourceMSORole(),
			"mso_user":        datasourceMSOUser(),
			"mso_label":       datasourceMSOLabel(),
			"mso_tenant":      datasourceMSOTenant(),
		},

		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		URL:        d.Get("url").(string),
		IsInsecure: d.Get("insecure").(bool),
		ProxyUrl:   d.Get("proxy_url").(string),
	}

	if err := config.Valid(); err != nil {
		return nil, err
	}

	return config.getClient(), nil
}

func (c Config) Valid() error {

	if c.Username == "" {
		return fmt.Errorf("Username must be provided for the MSO provider")
	}

	if c.Password == "" {

		return fmt.Errorf("Password must be provided for the MSO provider")
	}
	if c.URL == "" {
		return fmt.Errorf("URL must be provided for MSO provider")
	}

	return nil
}

func (c Config) getClient() interface{} {
	if c.Password != "" {

		return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl))

	}
	return nil
}

// Config
type Config struct {
	Username   string
	Password   string
	IsInsecure bool
	ProxyUrl   string
	URL        string
}
