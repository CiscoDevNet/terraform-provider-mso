package mso

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				DefaultFunc: schema.EnvDefaultFunc("MSO_INSECURE", true),
				Description: "Allow insecure HTTPS client",
			},
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_DOMAIN", nil),
				Description: "Domain name for remote user authentication",
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_PROXY_URL", nil),
				Description: "Proxy Server URL with port number",
			},
			"platform": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_PLATFORM", nil),
				Description: "Parameter that specifies where MSO is installed", // defaults to "mso"
				ValidateFunc: validation.StringInSlice([]string{
					"mso",
					"nd",
				}, false),
			},
			"retries": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSO_RETRIES", nil),
				Description: "Number of retries for REST API calls.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mso_schema":                                      resourceMSOSchema(),
			"mso_schema_site":                                 resourceMSOSchemaSite(),
			"mso_site":                                        resourceMSOSite(),
			"mso_remote_location":                             resourceMSORemoteLocation(),
			"mso_user":                                        resourceMSOUser(),
			"mso_label":                                       resourceMSOLabel(),
			"mso_schema_template":                             resourceMSOSchemaTemplate(),
			"mso_tenant":                                      resourceMSOTenant(),
			"mso_schema_template_bd":                          resourceMSOTemplateBD(),
			"mso_schema_template_vrf":                         resourceMSOSchemaTemplateVrf(),
			"mso_schema_template_bd_subnet":                   resourceMSOTemplateBDSubnet(),
			"mso_schema_template_anp":                         resourceMSOSchemaTemplateAnp(),
			"mso_schema_template_anp_epg":                     resourceMSOSchemaTemplateAnpEpg(),
			"mso_schema_template_anp_epg_contract":            resourceMSOTemplateAnpEpgContract(),
			"mso_schema_template_contract":                    resourceMSOTemplateContract(),
			"mso_schema_template_anp_epg_subnet":              resourceMSOSchemaTemplateAnpEpgSubnet(),
			"mso_schema_template_l3out":                       resourceMSOTemplateL3out(),
			"mso_schema_template_external_epg":                resourceMSOTemplateExtenalepg(),
			"mso_schema_template_contract_filter":             resourceMSOTemplateContractFilter(),
			"mso_schema_template_external_epg_contract":       resourceMSOTemplateExternalEpgContract(),
			"mso_schema_template_filter_entry":                resourceMSOSchemaTemplateFilterEntry(),
			"mso_schema_template_external_epg_subnet":         resourceMSOTemplateExtenalepgSubnet(),
			"mso_schema_site_anp_epg_static_leaf":             resourceMSOSchemaSiteAnpEpgStaticleaf(),
			"mso_schema_site_anp_epg_static_port":             resourceMSOSchemaSiteAnpEpgStaticPort(),
			"mso_schema_site_anp_epg_bulk_staticport":         resourceMSOSchemaSiteAnpEpgBulkStaticPort(),
			"mso_schema_site_bd":                              resourceMSOSchemaSiteBd(),
			"mso_schema_site_anp_epg_subnet":                  resourceMSOSchemaSiteAnpEpgSubnet(),
			"mso_schema_site_anp_epg_domain":                  resourceMSOSchemaSiteAnpEpgDomain(),
			"mso_schema_site_bd_l3out":                        resourceMSOSchemaSiteBdL3out(),
			"mso_schema_site_vrf":                             resourceMSOSchemaSiteVrf(),
			"mso_schema_site_vrf_route_leak":                  resourceMSOSchemaSiteVrfRouteLeak(),
			"mso_schema_site_vrf_region":                      resourceMSOSchemaSiteVrfRegion(),
			"mso_schema_site_bd_subnet":                       resourceMSOSchemaSiteBdSubnet(),
			"mso_rest":                                        resourceMSORest(),
			"mso_schema_template_deploy":                      resourceMSOSchemaTemplateDeploy(),
			"mso_schema_template_deploy_ndo":                  resourceNDOSchemaTemplateDeploy(),
			"mso_schema_site_vrf_region_cidr_subnet":          resourceMSOSchemaSiteVrfRegionCidrSubnet(),
			"mso_schema_site_vrf_region_cidr":                 resourceMSOSchemaSiteVrfRegionCidr(),
			"mso_schema_site_anp":                             resourceMSOSchemaSiteAnp(),
			"mso_schema_site_anp_epg":                         resourceMSOSchemaSiteAnpEpg(),
			"mso_schema_template_anp_epg_selector":            resourceMSOSchemaTemplateAnpEpgSelector(),
			"mso_schema_site_external_epg":                    resourceMSOSchemaSiteExternalEpg(),
			"mso_schema_template_external_epg_selector":       resourceSchemaTemplateExternalEPGSelector(),
			"mso_schema_template_anp_epg_useg_attr":           resourceMSOSchemaTemplateAnpEpgUsegAttr(),
			"mso_schema_site_anp_epg_selector":                resourceMSOSchemaSiteAnpEpgSelector(),
			"mso_schema_template_vrf_contract":                resourceMSOTemplateVRFContract(),
			"mso_schema_site_external_epg_selector":           resourceMSOSchemaSiteExternalEpgSelector(),
			"mso_schema_template_service_graph":               resourceMSOSchemaTemplateServiceGraphs(),
			"mso_schema_site_service_graph_node":              resourceMSOSchemaSiteServiceGraphNode(),
			"mso_schema_site_service_graph":                   resourceMSOSchemaSiteServiceGraph(),
			"mso_service_node_type":                           resourceMSOServiceNodeType(),
			"mso_schema_template_contract_service_graph":      resourceMSOSchemaTemplateContractServiceGraph(),
			"mso_system_config":                               resourceMSOSystemConfig(),
			"mso_schema_site_contract_service_graph":          resourceMSOSchemaSiteContractServiceGraph(),
			"mso_schema_site_contract_service_graph_listener": resourceMSOSchemaSiteContractServiceGraphListener(),
			"mso_template":                                    resourceMSOTemplate(),
			"mso_tenant_policies_ipsla_monitoring_policy":     resourceMSOIPSLAMonitoringPolicy(),
			"mso_tenant_policies_route_map_policy_multicast":  resourceMSOMcastRouteMapPolicy(),
			"mso_tenant_policies_dhcp_relay_policy":           resourceMSOTenantPoliciesDHCPRelayPolicy(),
			"mso_fabric_policies_vlan_pool":                   resourceMSOVlanPool(),
			"mso_fabric_policies_physical_domain":             resourceMSOPhysicalDomain(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mso_schema":                                      datasourceMSOSchema(),
			"mso_schema_site":                                 datasourceMSOSchemaSite(),
			"mso_site":                                        datasourceMSOSite(),
			"mso_remote_location":                             datasourceMSORemoteLocation(),
			"mso_role":                                        datasourceMSORole(),
			"mso_user":                                        datasourceMSOUser(),
			"mso_label":                                       datasourceMSOLabel(),
			"mso_schema_template":                             datasourceMSOSchemaTemplate(),
			"mso_tenant":                                      datasourceMSOTenant(),
			"mso_schema_template_bd":                          dataSourceMSOTemplateBD(),
			"mso_schema_template_vrf":                         datasourceMSOSchemaTemplateVrf(),
			"mso_schema_template_bd_subnet":                   dataSourceMSOTemplateSubnetBD(),
			"mso_schema_template_anp":                         datasourceMSOSchemaTemplateAnp(),
			"mso_schema_template_anp_epg":                     datasourceMSOSchemaTemplateAnpEpg(),
			"mso_schema_template_anp_epg_contract":            dataSourceMSOTemplateAnpEpgContract(),
			"mso_schema_template_contract":                    dataSourceMSOTemplateContract(),
			"mso_schema_template_anp_epg_subnet":              dataSourceMSOSchemaTemplateAnpEpgSubnet(),
			"mso_schema_template_l3out":                       dataSourceMSOTemplateL3out(),
			"mso_schema_template_external_epg":                dataSourceMSOTemplateExternalepg(),
			"mso_schema_template_contract_filter":             dataSourceMSOTemplateContractFilter(),
			"mso_schema_template_external_epg_contract":       dataSourceMSOTemplateExternalEpgContract(),
			"mso_schema_template_filter_entry":                dataSourceMSOSchemaTemplateFilterEntry(),
			"mso_schema_template_external_epg_subnet":         dataSourceMSOTemplateExternalEpgSubnet(),
			"mso_schema_site_anp":                             dataSourceMSOSchemaSiteAnp(),
			"mso_schema_site_anp_epg":                         dataSourceMSOSchemaSiteAnpEpg(),
			"mso_schema_site_anp_epg_static_leaf":             dataSourceMSOSchemaSiteAnpEpgStaticleaf(),
			"mso_schema_site_anp_epg_static_port":             datasourceMSOSchemaSiteAnpEpgStaticPort(),
			"mso_schema_site_anp_epg_bulk_staticport":         datasourceMSOSchemaSiteAnpEpgBulkStaticPort(),
			"mso_schema_site_bd":                              dataSourceMSOSchemaSiteBd(),
			"mso_schema_site_anp_epg_subnet":                  datasourceMSOSchemaSiteAnpEpgSubnet(),
			"mso_schema_site_anp_epg_domain":                  dataSourceMSOSchemaSiteAnpEpgDomain(),
			"mso_schema_site_bd_l3out":                        dataSourceMSOSchemaSiteBdL3out(),
			"mso_schema_site_vrf":                             dataSourceMSOSchemaSiteVrf(),
			"mso_schema_site_vrf_region":                      dataSourceMSOSchemaSiteVrfRegion(),
			"mso_schema_site_vrf_route_leak":                  dataSourceMSOSchemaSiteVrfRouteLeak(),
			"mso_schema_site_bd_subnet":                       dataSourceMSOSchemaSiteBdSubnet(),
			"mso_schema_site_vrf_region_cidr_subnet":          dataSourceMSOSchemaSiteVrfRegionCidrSubnet(),
			"mso_schema_site_vrf_region_cidr":                 dataSourceMSOSchemaSiteVrfRegionCidr(),
			"mso_schema_template_anp_epg_selector":            datasourceMSOSchemaTemplateAnpEpgSelector(),
			"mso_schema_site_external_epg":                    dataSourceMSOSchemaSiteExternalEpg(),
			"mso_schema_template_external_epg_selector":       datasourceSchemaTemplateExternalEPGSelector(),
			"mso_schema_template_anp_epg_useg_attr":           dataSourceMSOSchemaTemplateAnpEpgUsegAttr(),
			"mso_schema_site_anp_epg_selector":                datasourceMSOSchemaSiteAnpEpgSelector(),
			"mso_schema_template_vrf_contract":                dataSourceMSOTemplateVRFContract(),
			"mso_schema_site_external_epg_selector":           datasourceMSOSchemaSiteExternalEpgSelector(),
			"mso_schema_template_service_graph":               dataSourceMSOSchemaTemplateServiceGraph(),
			"mso_service_node_type":                           dataSourceMSOServiceNodeType(),
			"mso_schema_site_service_graph":                   datasourceMSOSchemaSiteServiceGraph(),
			"mso_schema_template_contract_service_graph":      dataSourceMSOSchemaTemplateContractServiceGraph(),
			"mso_system_config":                               dataSourceMSOSystemConfig(),
			"mso_rest":                                        datasourceMSORest(),
			"mso_schema_site_contract_service_graph":          dataSourceMSOSchemaSiteContractServiceGraph(),
			"mso_schema_site_contract_service_graph_listener": dataSourceMSOSchemaSiteContractServiceGraphListener(),
			"mso_template":                                    datasourceMSOTemplate(),
			"mso_tenant_policies_ipsla_monitoring_policy":     datasourceMSOIPSLAMonitoringPolicy(),
			"mso_tenant_policies_route_map_policy_multicast":  datasourceMSOMcastRouteMapPolicy(),
			"mso_tenant_policies_dhcp_relay_policy":           datasourceMSOTenantPoliciesDHCPRelayPolicy(),
			"mso_fabric_policies_vlan_pool":                   datasourceMSOVlanPool(),
			"mso_fabric_policies_physical_domain":             datasourceMSOPhysicalDomain(),
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
		Domain:     d.Get("domain").(string),
		Platform:   d.Get("platform").(string),
	}

	config.MaxRetries = 2
	if d.Get("retries").(string) != "" {
		maxRetries, err := strconv.Atoi(d.Get("retries").(string))
		if err != nil {
			return nil, fmt.Errorf("Invalid value for retries")
		}
		config.MaxRetries = maxRetries
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

		return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl), client.Domain(c.Domain), client.Platform(c.Platform), client.MaxRetries(c.MaxRetries))

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
	Domain     string
	Platform   string
	MaxRetries int
}
