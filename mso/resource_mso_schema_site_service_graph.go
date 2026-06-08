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

// resourceMSOSchemaSiteServiceGraph manages an mso_schema_site_service_graph
// entry (the site-level instantiation of a template service graph, binding
// physical or cloud L4-L7 service devices to each graph node).
//
// Delete implementation note: the NDO API does not support deleting a site
// service graph in isolation. The Delete handler is a no-op — the graph entry
// persists in the schema until the parent mso_schema_template_service_graph is
// deleted or the site association is removed. This is intentional API behaviour
// and not a provider limitation.
//
// ForceNew on identity fields: schema_id, template_name, site_id, and
// service_graph_name together identify the exact API path
// (/sites/{siteId}-{template}/serviceGraphs/{graphName}) used for all PATCH
// operations. Changing any of them would target a different object in the
// schema document, which is not achievable via an in-place update. Because
// Delete is a no-op, ForceNew ensures Terraform plans a new resource at the
// new location rather than silently leaving the old entry orphaned.
//
// CustomizeDiff: validates provider_connector_type values against the node
// type declared in the parent template service graph ("other" allows only
// "none"/"redir"; "firewall" additionally allows "snat", "dnat",
// "snat_dnat"). Validation is skipped when schema_id is not yet known (e.g.
// when the schema resource is being created in the same plan).
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
				// ForceNew: schema_id is part of the resource identity and the
				// API path. Changing it targets a different schema document,
				// which requires destroy+recreate.
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				// ForceNew: template_name is part of the API path key
				// ({siteId}-{template}). Changing it targets a different
				// site-template association, which requires destroy+recreate.
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				// ForceNew: site_id is part of the API path key
				// ({siteId}-{template}). Changing it targets a different site,
				// which requires destroy+recreate.
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				// ForceNew: service_graph_name is the final segment of the API
				// path. Changing it targets a different graph entry, which
				// requires destroy+recreate.
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

			// When the schema_id is empty, it means the schema resource is being created in the same plan which means that the value is only known after apply.
			// In this case, we can skip the validation and validation is triggered in the create function.
			if schemaId == "" {
				return nil
			}

			cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
			if err != nil {
				return err
			}

			_, serviceNode := diff.GetChange("service_node")
			err = validateServiceNodeConfig(msoClient, serviceNode, cont, templateName.(string), graphName.(string))
			if err != nil {
				return err
			}

			return err
		},
	}
}

func validateServiceNodeConfig(msoClient *client.Client, serviceNode interface{}, cont *container.Container, templateName, graphName string) error {
	if len(serviceNode.([]interface{})) != 0 {
		for _, node := range serviceNode.([]interface{}) {

			serviceNodeMap := node.(map[string]interface{})
			if !valueInSliceofStrings(serviceNodeMap["provider_connector_type"].(string), []string{"none", "redir"}) { // If provider_connector_type is not none, then validate the user input.
				sgCont, _, err := getTemplateServiceGraphCont(cont, templateName, graphName)
				if strings.Contains(fmt.Sprint(err), "No Template found") {
					// The function getTemplateServiceGraphCont() is not required when the template is attached to physical site.
					return nil
				} else if err != nil {
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

					/* Loop through the templateServiceNodeList and validate the site level user input(provider_connector_type)
					to verify it's value for nodetype 'other' and 'firewall'. */
					if len(serviceNode.([]interface{})) != len(templateServiceNodeList) {
						return fmt.Errorf("service graph has %d service node(s) in the template but %d service node(s) were provided", len(templateServiceNodeList), len(serviceNode.([]interface{})))
					}
					for i, val := range serviceNode.([]interface{}) {
						serviceNode := val.(map[string]interface{})
						if templateServiceNodeList[i] == "other" && !valueInSliceofStrings(serviceNode["provider_connector_type"].(string), []string{"none", "redir"}) {
							return fmt.Errorf("The expected value for service_node.%d.provider_connector_type have to be one of [none, redir] when template's service node type is other, got %s.", i, serviceNode["provider_connector_type"])
						} else if templateServiceNodeList[i] == "firewall" && !valueInSliceofStrings(serviceNode["provider_connector_type"].(string), []string{"none", "redir", "snat", "dnat", "snat_dnat"}) {
							return fmt.Errorf("The expected value for service_node.%d.provider_connector_type have to be one of [none, redir, snat, dnat, snat_dnat] when template's service node type is firewall, got %s.", i, serviceNode["provider_connector_type"])
						}
					}
					return nil
				}
			}
		}
	}
	return nil
}

func resourceMSOSchemaSiteServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	if len(get_attribute) < 7 {
		return nil, fmt.Errorf("invalid import ID %q: expected format {schema_id}/sites/{site_id}/template/{template_name}/serviceGraphs/{graph_name}", d.Id())
	}
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
	if err != nil {
		return nil, err
	}
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(fmt.Sprintf("%s/sites/%s/template/%s/serviceGraphs/%s", schemaId, siteId, templateName, graphName))
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Creation Site Service Graph")
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
		// Validate here because when schema_id input is unknown during plan, the validation is skipped from in the CustomizeDiff function.
		// Downside is that the validation is done twice.
		err = validateServiceNodeConfig(msoClient, siteServiceNodes, cont, templateName, graphName)
		if err != nil {
			return err
		}
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
	if err != nil {
		return err
	}
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(nodeIdSt)
	return nil
}

func resourceMSOSchemaSiteServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Update Site Service Graph")
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
			// Validate here because when schema_id input is unknown during plan, the validation is skipped from in the CustomizeDiff function.
			// Downside is that the validation is done twice.
			err = validateServiceNodeConfig(msoClient, siteServiceNodes, cont, templateName, graphName)
			if err != nil {
				return err
			}

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
	templateNodes := graphCont.S("serviceNodes").Data().([]interface{})
	siteNodes := siteServiceNodes.([]interface{})
	if len(siteNodes) != len(templateNodes) {
		return nil, fmt.Errorf("service graph has %d service node(s) in the template but %d service node(s) were provided", len(templateNodes), len(siteNodes))
	}
	siteServiceNodeList := make([]interface{}, 0, len(templateNodes))
	for index, serviceNode := range templateNodes {
		siteServiceNodeMap := siteNodes[index].(map[string]interface{})

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

func getSiteServiceNodeCont(graphCont *container.Container, schemaId, templateName, graphName, nodeName string) (*container.Container, int, error) {
	nodesCount, err := graphCont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load site service node count")
	}
	for i := 0; i < nodesCount; i++ {
		nodeCont, err := graphCont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load site service node element")
		}

		nodeRef := models.StripQuotes(nodeCont.S("serviceNodeRef").String())

		nodeSplit := strings.Split(nodeRef, "/")
		if len(nodeSplit) == 9 {
			if nodeSplit[2] == schemaId && nodeSplit[4] == templateName && nodeSplit[6] == graphName && nodeSplit[8] == nodeName {
				return nodeCont, i, nil
			}
		} else {
			return nil, -1, fmt.Errorf("Split on nodeRef failed")
		}
	}
	return nil, -1, fmt.Errorf("Unable to find site service node")
}

func getSiteServiceGraphCont(cont *container.Container, schemaId, templateName, siteId, graphName string) (*container.Container, int, error) {
	sitesCount, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to find sites")
	}

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load site element")
		}

		siteTemplate := models.StripQuotes(siteCont.S("templateName").String())
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if siteTemplate == templateName && siteId == apiSiteId {
			sgCount, err := siteCont.ArrayCount("serviceGraphs")
			if err != nil {
				return nil, -1, fmt.Errorf("Unable to load site service graphs")
			}
			for j := 0; j < sgCount; j++ {
				sgCont, err := siteCont.ArrayElement(j, "serviceGraphs")
				if err != nil {
					return nil, -1, fmt.Errorf("Unable to load site service graph element")
				}

				graphRef := models.StripQuotes(sgCont.S("serviceGraphRef").String())
				graphEle := strings.Split(graphRef, "/")

				if len(graphEle) != 7 {
					return nil, -1, fmt.Errorf("Invalid site service graph")
				}

				if schemaId == graphEle[2] && templateName == graphEle[4] && graphName == graphEle[6] {
					return sgCont, j, nil
				}
			}
		}
	}

	return nil, -1, fmt.Errorf("Unable to find site service graph")
}

func getTemplateServiceNodeCont(cont *container.Container, nodeName, nodeType string) (*container.Container, int, error) {
	nodeCount, err := cont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load node count")
	}

	for i := 0; i < nodeCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load node element")
		}

		apiNodeName := models.StripQuotes(nodeCont.S("name").String())
		apiNodeType := models.StripQuotes(nodeCont.S("serviceNodeTypeId").String())

		if apiNodeName == nodeName && apiNodeType == nodeType {
			return nodeCont, i, nil
		}
	}
	return nil, -1, fmt.Errorf("Unable to find the service node")
}
