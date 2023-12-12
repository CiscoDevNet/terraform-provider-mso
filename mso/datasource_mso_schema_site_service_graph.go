package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteServiceGraph() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteServiceGraphRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configure service nodes for the site service graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_dn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_interface": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_interface": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func dataSourceMSOSchemaSiteServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning datasource Read")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	graphName := d.Get("service_graph_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	graphCont, _, err := getSiteServiceGraphCont(cont, schemaId, templateName, siteId, graphName)
	if err != nil {
		d.SetId("")
		return nil
	}

	serviceNodeList, err := setServiceNodeList(graphCont)
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(fmt.Sprintf("%s/sites/%s/template/%s/serviceGraphs/%s", schemaId, siteId, templateName, graphName))
	log.Printf("[DEBUG] %s: Datasource read finished successfully", d.Id())
	return nil
}
