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
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"node_index": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"service_node_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
							ForceNew:     true,
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
		log.Printf("graphcont err %v", err)
		return err
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)
	nodeInd := d.Get("node_index").(int)

	tempNodeCont, _, err := getTemplateServiceNodeContFromIndex(sgCont, nodeInd)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("node_index", nodeInd)

	nodeId := models.StripQuotes(tempNodeCont.S("serviceNodeTypeId").String())

	nodeName := models.StripQuotes(tempNodeCont.S("name").String())

	nodeIdHuman, err := getNodeNameFromId(msoClient, nodeId)

	if err != nil {
		return err
	}

	d.Set("service_node_type", nodeIdHuman)

	siteParams := make([]interface{}, 0, 1)

	sitesCount, err := cont.ArrayCount("sites")

	if err != nil {
		d.SetId(graphName)
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

		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			schemaId,
			stateTemplate,
			apiSiteId,
			graphName,
		)

		if err == nil {
			nodeCont, _, nodeerr := getSiteServiceNodeCont(
				graphCont,
				schemaId,
				stateTemplate,
				graphName,
				nodeName,
			)

			if nodeerr == nil {
				siteMap := make(map[string]interface{})

				deviceDn := models.StripQuotes(nodeCont.S("device", "dn").String())

				dnSplit := strings.Split(deviceDn, "/")

				tnName := strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")
				siteMap["tenant_name"] = tnName
				siteMap["node_name"] = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")
				siteMap["site_id"] = apiSiteId

				siteParams = append(siteParams, siteMap)
			}

		}

	}

	d.Set("site_nodes", siteParams)
	d.SetId(graphName)
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
