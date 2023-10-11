package mso

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceTemplateContractServiceGraph_New() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTemplateContractServiceGraphRead_New,

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_site_id": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
			},
			"service_graph_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_relationship": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_connector_bd_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_bd_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_bd_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_cluster_interface": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_connector_redirect_policy_tenant": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_connector_redirect_policy": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_subnet_ips": &schema.Schema{
							Type:       schema.TypeList,
							Computed:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_cluster_interface": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_redirect_policy_tenant": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_redirect_policy": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_subnet_ips": &schema.Schema{
							Type:       schema.TypeList,
							Computed:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
					},
				},
			},
		},
	}
}

func datasourceTemplateContractServiceGraphRead_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Data source begining read Template Contract Service Graph")
	setTemplateContractServiceGraphAttrs(d, m, true)
	log.Printf("[DEBUG] Data source completed read Template Contract Service Graph")
	return nil
}
