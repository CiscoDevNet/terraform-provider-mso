package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteServiceGraphNode() *schema.Resource {
	return &schema.Resource{
		Read: datasourceMSOSchemaSiteServiceGraphNodeRead,

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
			"service_node_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"site_nodes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"tenant_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"node_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		}),
	}
}

func datasourceMSOSchemaSiteServiceGraphNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Datasource MSO Schema Site Service GraphNode: Beginning Read")

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	stateTemplate := d.Get("template_name").(string)
	graphName := d.Get("service_graph_name").(string)

	sgCont, _, err := getTemplateServiceGraphCont(cont, stateTemplate, graphName)

	log.Printf("GraphCont %v", sgCont)
	if err != nil {
		log.Printf("graphcont err %v", err)
		return err
	}

	nodeType := d.Get("service_node_type").(string)

	nodeId, err := getNodeIdFromName(msoClient, nodeType)
	if err != nil {
		return err
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)

	nodeIdSt := d.Get("service_node_name").(string)
	_, _, err = getTemplateServiceNodeCont(sgCont, nodeIdSt, nodeId)

	if err != nil {
		d.Set("service_node_type", "")
		log.Printf("nodecont err %v", err)
		return err
	} else {
		d.Set("service_node_type", nodeType)
	}

	var siteParams []interface{}

	if tempVar, ok := d.GetOk("site_nodes"); ok {
		siteParams = tempVar.([]interface{})

		for ind, site := range siteParams {
			siteMap := site.(map[string]interface{})
			if siteMap["site_id"] == "" {
				return fmt.Errorf("site_id is required in site_nodes list")
			}

			graphCont, _, err := getSiteServiceGraphCont(
				cont,
				schemaId,
				stateTemplate,
				siteMap["site_id"].(string),
				graphName,
			)

			if err != nil {
				log.Printf("sitegraphcont err %v", err)
				return err
			}

			nodeCont, _, err := getSiteServiceNodeCont(
				graphCont,
				schemaId,
				stateTemplate,
				graphName,
				nodeIdSt,
			)

			if err != nil {
				log.Printf("sitenodecont err %v", err)
				return err
			}

			deviceDn := models.StripQuotes(nodeCont.S("device", "dn").String())

			dnSplit := strings.Split(deviceDn, "/")

			tnName := strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")
			siteMap["tenant_name"] = tnName
			siteMap["node_name"] = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")

			siteParams[ind] = siteMap

		}

		d.Set("site_nodes", siteParams)
	} else {
		d.Set("site_nodes", nil)
	}

	d.SetId(nodeIdSt)
	return nil
}
