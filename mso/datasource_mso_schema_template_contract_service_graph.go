package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaTemplateContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaTemplateContractServiceGraphRead,
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
			"service_graph_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceMSOSchemaTemplateContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning datasource Read")

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s", schemaId, templateName, contractName))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSchemaTemplateContractServiceGraphAttrs(cont, d, true)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Datasource read finished successfully", d.Id())
	return nil
}
