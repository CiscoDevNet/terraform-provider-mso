package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteContractServiceGraphRead,

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_relationship": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_connector_cluster_interface": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_redirect_policy_tenant": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_redirect_policy": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_cluster_interface": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_redirect_policy_tenant": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_redirect_policy": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_subnet_ips": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOSchemaSiteContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning datasource Read")

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	siteID := d.Get("site_id").(string)

	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s", schemaID, siteID, templateName, contractName))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSiteContractServiceGraphAttrs(cont, d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Datasource read finished successfully", d.Id())
	return nil
}
