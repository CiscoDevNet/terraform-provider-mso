package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
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
			"site_nodes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_index": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_node_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)

	siteParams := make([]interface{}, 0, 1)
	for _, node := range sgCont.S("serviceNodes").Data().([]interface{}) {
		nodeIndex := node.(map[string]interface{})["index"].(float64)
		nodeCont, _, _ := getTemplateServiceNodeContFromIndex(sgCont, int(nodeIndex))
		serviceNodeTypeId := models.StripQuotes(nodeCont.S("serviceNodeTypeId").String())
		serviceNodeTypeIdHuman, err := getNodeNameFromId(msoClient, serviceNodeTypeId)
		if err != nil {
			return err
		}
		sitesCount, err := cont.ArrayCount("sites")
		if err != nil {
			d.SetId(fmt.Sprintf("%s/templates/%s/serviceGraphs/%s", schemaId, stateTemplate, graphName))
			d.Set("site_nodes", nil)
			log.Printf("Unable to find sites")
			return nil
		}
		for i := 0; i < sitesCount; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return fmt.Errorf("Unable to load site element")
			}
			apiSiteId := models.StripQuotes(siteCont.S("siteId").String())
			graphCont, _, err := getSiteServiceGraphCont(cont, schemaId, stateTemplate, apiSiteId, graphName)
			if err == nil {
				nodeName := models.StripQuotes(nodeCont.S("name").String())
				siteServiceNodeCont, _, nodeerr := getSiteServiceNodeCont(graphCont, schemaId, stateTemplate, graphName, nodeName)
				if nodeerr == nil {
					deviceDn := models.StripQuotes(siteServiceNodeCont.S("device", "dn").String())
					dnSplit := strings.Split(deviceDn, "/")
					siteMap := make(map[string]interface{})
					siteMap["node_index"] = int(nodeIndex)
					siteMap["service_node_type"] = serviceNodeTypeIdHuman
					siteMap["tenant_name"] = strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")
					siteMap["node_name"] = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")
					siteMap["site_id"] = apiSiteId
					siteParams = append(siteParams, siteMap)
				}
			}
		}
	}
	d.Set("site_nodes", siteParams)
	d.SetId(fmt.Sprintf("%s/templates/%s/serviceGraphs/%s", schemaId, stateTemplate, graphName))
	return nil
}

func getTemplateServiceNodeContFromIndex(cont *container.Container, ind int) (*container.Container, int, error) {

	nodeCount, err := cont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load node count")
	}

	for i := 0; i < nodeCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load node element")
		}

		apiInd := int(nodeCont.S("index").Data().(float64))

		if apiInd == ind {
			return nodeCont, i, nil
		}

	}

	return nil, -1, fmt.Errorf("Unable to find the service node")
}
