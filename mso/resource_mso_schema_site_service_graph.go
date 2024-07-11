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

func resourceMSOSchemaSiteServiceGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteServiceGraphCreate,
		Read:   resourceMSOSchemaSiteServiceGraphRead,
		Update: resourceMSOSchemaSiteServiceGraphUpdate,
		Delete: resourceMSOSchemaSiteServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteServiceGraphImport,
		},

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
			"site_id": &schema.Schema{
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
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configure service nodes for the service graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_dn": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"redir",
							}, false),
							Default: "none",
						},
						"provider_connector_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "none",
						},
						"consumer_interface": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"provider_interface": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			/* This function validates the user input for service_node.provider_connector_type when
			the template_service_graph.service_node.type is 'other' or 'firewall'.

			- The user input for site_service_graph.service_node.other_provider_connector_type should be one of 'none' or 'redir',
			when the corresponding template_service_graph.service_node.type is 'other'.

			- The user input for site_service_graph.servicenode.firewall_provider_connector_type_list should be one of 'none', 'redir', 'snat', 'dnat' or 'snat_dnat',
			when the corresponding template_service_graph.service_node.type is 'firewall'.
			*/

			// Create a list of service node types using the user input(template service graph).
			msoClient := v.(*client.Client)
			_, schemaId := diff.GetChange("schema_id")
			_, templateName := diff.GetChange("template_name")
			_, graphName := diff.GetChange("service_graph_name")
			cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
			if err != nil {
				return err
			}

			sgCont, _, err := getTemplateServiceGraphCont(cont, templateName.(string), graphName.(string))
			if strings.Contains(fmt.Sprint(err), "No Template found") {
				// The function getTemplateServiceGraphCont() is not required when the template is attched to physical site.
				return nil
			} else if err != nil {
				log.Printf("graphcont err %v", err)
				return err
			} else {
				/* The function getTemplateServiceGraphCont() is required when the template is attached to cloud sites.
				provider_connector_type is applicable only for cloud sites. */
				var templateServiceNodeList []string
				serviceNodes := sgCont.S("serviceNodes").Data().([]interface{})
				for _, val := range serviceNodes {
					serviceNodeValues := val.(map[string]interface{})
					nodeId := models.StripQuotes(serviceNodeValues["serviceNodeTypeId"].(string))

					nodeType, err := getNodeNameFromId(msoClient, nodeId)
					if err != nil {
						return err
					}

					templateServiceNodeList = append(templateServiceNodeList, nodeType)
				}

				/* Loop trough the templateServiceNodeList and validate the site level user input(provider_connector_type)
				to verify it's value for nodetype 'other' and 'firewall'. */
				_, siteServiceNodes := diff.GetChange("service_node")

				for i, val := range siteServiceNodes.([]interface{}) {
					serviceNode := val.(map[string]interface{})
					if templateServiceNodeList[i] == "other" && !valueInSliceofStrings(serviceNode["provider_connector_type"].(string), []string{"none", "redir"}) {
						return fmt.Errorf("The expected value for service_node.%d.provider_connector_type have to be one of [none, redir] when template's service node type is other, got %s.", i, serviceNode["provider_connector_type"])
					} else if templateServiceNodeList[i] == "firewall" && !valueInSliceofStrings(serviceNode["provider_connector_type"].(string), []string{"none", "redir", "snat", "dnat", "snat_dnat"}) {
						return fmt.Errorf("The expected value for service_node.%d.provider_connector_type have to be one of [none, redir, snat, dnat, snat_dnat] when template's service node type is firewall, got %s.", i, serviceNode["provider_connector_type"])
					}
				}
				return nil
			}
		},
	}
}

func resourceMSOSchemaSiteServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	siteId := get_attribute[2]
	templateName := get_attribute[4]
	graphName := get_attribute[6]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	graphCont, _, err := getSiteServiceGraphCont(cont, schemaId, templateName, siteId, graphName)
	if err != nil {
		d.SetId("")
		return nil, err
	}

	serviceNodeList, err := setServiceNodeList(graphCont)
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(fmt.Sprintf("%s/templates/%s/serviceGraphs/%s", schemaId, templateName, graphName))
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Site Service Graph")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	graphName := d.Get("service_graph_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	graphCont, _, err := getTemplateServiceGraphCont(cont, templateName, graphName)
	if err != nil {
		return err
	}

	var siteServiceNodeList []interface{}

	if siteServiceNodes, ok := d.GetOk("service_node"); ok {
		siteServiceNodeList, err = createSiteServiceNodeList(msoClient, siteServiceNodes, graphCont)
		if err != nil {
			return err
		}
	}
	serviceNodePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%s/serviceNodes", siteId, templateName, graphName)
	siteServiceGraphPayload := models.GetPatchPayloadList("add", serviceNodePath, siteServiceNodeList)
	_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s/template/%s/serviceGraphs/%s", schemaId, siteId, templateName, graphName))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaSiteServiceGraphRead(d, m)
}

func resourceMSOSchemaSiteServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	nodeIdSt := d.Id()
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	graphName := d.Get("service_graph_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
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

	d.SetId(nodeIdSt)
	return nil
}

func resourceMSOSchemaSiteServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Site Service Graph")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	graphName := d.Get("service_graph_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	if d.HasChange("service_node") {
		graphCont, _, err := getTemplateServiceGraphCont(cont, templateName, graphName)
		if err != nil {
			return err
		}

		if siteServiceNodes, ok := d.GetOk("service_node"); ok {
			siteServiceNodeList, err := createSiteServiceNodeList(msoClient, siteServiceNodes, graphCont)
			if err != nil {
				return err
			}

			serviceNodePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%s/serviceNodes", siteId, templateName, graphName)
			siteServiceGraphPayload := models.GetPatchPayloadList("replace", serviceNodePath, siteServiceNodeList)
			_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload)
			if err != nil {
				return err
			}
		}
	}

	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteServiceGraphRead(d, m)
}

func resourceMSOSchemaSiteServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO]: Deletion of site Service Graph is not supported by the API.  Site Service Graph will be removed when site is disassociated from the template or when Service Graph is removed at the template level.")
	return nil
}

func createSiteServiceNodeList(msoClient *client.Client, siteServiceNodes interface{}, graphCont *container.Container) ([]interface{}, error) {
	siteServiceNodeList := make([]interface{}, 0, 1)
	for index, serviceNode := range graphCont.S("serviceNodes").Data().([]interface{}) {
		siteServiceNodeMap := siteServiceNodes.([]interface{})[index].(map[string]interface{})

		serviceNodeMap := map[string]interface{}{
			"serviceNodeRef": serviceNode.(map[string]interface{})["serviceNodeRef"],
			"device": map[string]interface{}{
				"dn": siteServiceNodeMap["device_dn"],
			},
			"consumerConnectorType": siteServiceNodeMap["consumer_connector_type"],
			"providerConnectorType": siteServiceNodeMap["provider_connector_type"],
			"consumerInterface":     siteServiceNodeMap["consumer_interface"],
			"providerInterface":     siteServiceNodeMap["provider_interface"],
		}
		siteServiceNodeList = append(siteServiceNodeList, serviceNodeMap)
	}
	return siteServiceNodeList, nil
}

func setServiceNodeList(graphCont *container.Container) ([]interface{}, error) {
	serviceNodeList := make([]interface{}, 0, 1)
	for _, val := range graphCont.S("serviceNodes").Data().([]interface{}) {
		serviceNodeMap := map[string]interface{}{
			"device_dn":               val.(map[string]interface{})["device"].(map[string]interface{})["dn"],
			"consumer_connector_type": val.(map[string]interface{})["consumerConnectorType"],
			"provider_connector_type": val.(map[string]interface{})["providerConnectorType"],
			"consumer_interface":      val.(map[string]interface{})["consumerInterface"],
			"provider_interface":      val.(map[string]interface{})["providerInterface"],
		}

		serviceNodeList = append(serviceNodeList, serviceNodeMap)
	}
	return serviceNodeList, nil
}
