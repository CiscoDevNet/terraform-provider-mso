package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateContractServiceChaining() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceMSOSchemaTemplateContractServiceChainingRead,
		SchemaVersion: version,

		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of service nodes in the service chaining graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_redirect": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"provider_connector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_redirect": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOSchemaTemplateContractServiceChainingRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Read Service Chaining (data source)")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	if err := setServiceChainingFromSchema(d, schemaCont, schemaId, templateName, contractName); err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Read Service Chaining (data source)")
	return nil
}
