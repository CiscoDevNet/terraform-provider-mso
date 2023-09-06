package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaTemplateServiceGraph() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaTemplateServiceGrapRead,

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
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node_type": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"service_node"},
				Deprecated:    "Use service_node to configure service nodes.",
			},
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configure service nodes for the service graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaTemplateServiceGrapRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	stateTemplate := d.Get("template_name").(string)
	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	sgCont, _, err := getTemplateServiceGraphCont(cont, stateTemplate, graphName)

	if err != nil {
		d.SetId("")
		return err
	}

	if tempVar, ok := d.GetOk("service_node_type"); ok {
		serviceNodeType := tempVar.(string)
		d.Set("service_node_type", serviceNodeType)
	} else {
		serviceNodeList := make([]interface{}, 0, 1)
		serviceNodes := sgCont.S("serviceNodes").Data().([]interface{})
		for _, val := range serviceNodes {
			serviceNodeValues := val.(map[string]interface{})
			serviceNodeMap := make(map[string]interface{})
			nodeId := models.StripQuotes(serviceNodeValues["serviceNodeTypeId"].(string))

			nodeType, err := getNodeNameFromId(msoClient, nodeId)
			if err != nil {
				return err
			}
			serviceNodeMap["type"] = nodeType

			serviceNodeList = append(serviceNodeList, serviceNodeMap)
		}
		d.Set("service_node", serviceNodeList)
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)
	d.SetId(fmt.Sprintf("%s/templates/%s/serviceGraphs/%s", schemaId, stateTemplate, graphName))
	return nil
}
