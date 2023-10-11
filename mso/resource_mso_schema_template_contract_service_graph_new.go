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

func resourceTemplateContractServiceGraph_New() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateContractServiceGraphCreate_New,
		Update: resourceTemplateContractServiceGraphUpdate_New,
		Read:   resourceTemplateContractServiceGraphRead_New,
		Delete: resourceTemplateContractServiceGraphDelete_New,

		Importer: &schema.ResourceImporter{
			State: resourceTemplateContractServiceGraphImport_New,
		},

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:       schema.TypeString,
				Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
				Optional:   true,
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
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
			"service_graph_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_site_id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
			},
			"service_graph_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"node_relationship": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_connector_bd_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_bd_schema_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_bd_template_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_bd_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_bd_schema_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_bd_template_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_cluster_interface": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_connector_redirect_policy_tenant": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_connector_redirect_policy": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"provider_subnet_ips": &schema.Schema{
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_cluster_interface": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_redirect_policy_tenant": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_connector_redirect_policy": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
						"consumer_subnet_ips": &schema.Schema{
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Deprecated: "Use mso_site_contract_service_graph resource to configure the site",
						},
					},
				},
			},
		},
	}
}

func resourceTemplateContractServiceGraphImport_New(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	serviceGraphTokens := strings.Split(d.Id(), "/")
	d.Set("schema_id", serviceGraphTokens[2])
	d.Set("template_name", serviceGraphTokens[4])
	d.Set("contract_name", serviceGraphTokens[6])
	err := setTemplateContractServiceGraphAttrs(d, m, true)
	if err != nil {
		return nil, err
	}
	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceTemplateContractServiceGraphCreate_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Service Graph: Beginning Creation")
	err := PostTemplateContractServiceGraphConfig("add", d, m)
	if err != nil {
		return err
	}
	err = resourceTemplateContractServiceGraphRead_New(d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return nil
}

func resourceTemplateContractServiceGraphUpdate_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Service Graph: Beginning Update")
	err := PostTemplateContractServiceGraphConfig("add", d, m)
	if err != nil {
		return err
	}
	err = resourceTemplateContractServiceGraphRead_New(d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return nil
}

func resourceTemplateContractServiceGraphRead_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Template Contract Service Graph")
	serviceGraphTokens := strings.Split(d.Id(), "/")
	d.Set("schema_id", serviceGraphTokens[2])
	d.Set("template_name", serviceGraphTokens[4])
	d.Set("contract_name", serviceGraphTokens[6])
	err := setTemplateContractServiceGraphAttrs(d, m, false)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}

func resourceTemplateContractServiceGraphDelete_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Delete Template Contract Service Graph")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("remove", tempPath, nil, nil)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Template Contract Service Graph")
	return nil
}

// Returns the Service Graph object from the Template Service Graph list based on the templateName and graphName
// Return values: Template serviceGraph Object, Template serviceGraph index position, error
func getTemplateServiceGraph(cont *container.Container, templateName, graphName string) (*container.Container, int, error) {
	templateCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, -1, fmt.Errorf("No Template found")
	}

	for i := 0; i < templateCount; i++ {
		templateCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to get template element")
		}

		apiTemplate := models.StripQuotes(templateCont.S("name").String())

		if apiTemplate == templateName {
			log.Printf("[DEBUG] Template found")

			sgCount, err := templateCont.ArrayCount("serviceGraphs")

			if err != nil {
				return nil, -1, fmt.Errorf("No Service Graph found")
			}

			for j := 0; j < sgCount; j++ {
				sgCont, err := templateCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to get service graph element")
				}

				apiSgName := models.StripQuotes(sgCont.S("name").String())

				if apiSgName == graphName {
					return sgCont, j, nil
				}
			}

		}
	}
	return nil, -1, fmt.Errorf("unable to find service graph")
}

// Returns the List of Service Graph Node map object with Consumer and Provider BD values
func getServiceNodesRelationshipObject(cont *container.Container, schemaID, graphSiteID, templateName string, nodeList, tempNodeList []interface{}, serviceGraphRef map[string]interface{}) ([]interface{}, []interface{}, error) {
	templateNodes := make([]interface{}, 0)
	siteNodes := make([]interface{}, 0, 1)

	for i := 0; i < len(tempNodeList); i++ {
		node := tempNodeList[i].(map[string]interface{})

		tempnodeRef := make(map[string]interface{})

		tempnodeRef["schemaId"] = serviceGraphRef["schemaId"]
		tempnodeRef["serviceGraphName"] = serviceGraphRef["serviceGraphName"]
		tempnodeRef["templateName"] = serviceGraphRef["templateName"]
		tempnodeRef["serviceNodeName"] = nodeList[i].(string)

		tempproConnector := make(map[string]interface{})
		tempproConnector["connectorType"] = "general"
		bdRef := make(map[string]interface{})
		if node["provider_connector_bd_schema_id"] != "" {
			bdRef["schemaId"] = node["provider_connector_bd_schema_id"].(string)
		} else {
			bdRef["schemaId"] = schemaID
		}

		if node["provider_connector_bd_template_name"] != "" {
			bdRef["templateName"] = node["provider_connector_bd_template_name"].(string)
		} else {
			bdRef["templateName"] = templateName
		}
		bdRef["bdName"] = node["provider_connector_bd_name"].(string)
		tempproConnector["bdRef"] = bdRef

		tempconConnector := make(map[string]interface{})
		tempconConnector["connectorType"] = "general"
		conbdRef := make(map[string]interface{})
		if node["consumer_connector_bd_schema_id"] != "" {
			conbdRef["schemaId"] = node["consumer_connector_bd_schema_id"].(string)
		} else {
			conbdRef["schemaId"] = schemaID
		}

		if node["consumer_connector_bd_template_name"] != "" {
			conbdRef["templateName"] = node["consumer_connector_bd_template_name"].(string)
		} else {
			conbdRef["templateName"] = templateName
		}
		conbdRef["bdName"] = node["consumer_connector_bd_name"].(string)
		tempconConnector["bdRef"] = conbdRef

		tempnodeMap := make(map[string]interface{})
		tempnodeMap["serviceNodeRef"] = tempnodeRef
		tempnodeMap["providerConnector"] = tempproConnector
		tempnodeMap["consumerConnector"] = tempconConnector

		templateNodes = append(templateNodes, tempnodeMap)
	}
	return templateNodes, siteNodes, nil
}

// Sets the resource attributes
func setTemplateContractServiceGraphAttrs(d *schema.ResourceData, m interface{}, importFlag bool) error {
	msoClient := m.(*client.Client)
	foundTemp := false
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No templates found")
	}

	temprelationList := make([]interface{}, 0, 1)
	for i := 0; i < tempCount; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return fmt.Errorf("Error in fetch of template")
		}
		template := models.StripQuotes(tempCont.S("name").String())
		if templateName == template {
			contractCount, err := tempCont.ArrayCount("contracts")
			if err != nil {
				return fmt.Errorf("No contracts found")
			}

			for j := 0; j < contractCount; j++ {
				contractCont, err := tempCont.ArrayElement(j, "contracts")
				if err != nil {
					return fmt.Errorf("Error fetching contract")
				}
				conName := models.StripQuotes(contractCont.S("name").String())
				if conName == contractName {
					if !contractCont.Exists("serviceGraphRelationship") {
						return fmt.Errorf("No service graph found")
					} else {
						// Template Contract Service Graph configurations
						graphRelation := contractCont.S("serviceGraphRelationship")
						graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
						tokens := strings.Split(graphRef, "/")

						d.Set("service_graph_name", tokens[6])

						if _, ok := d.GetOk("service_graph_schema_id"); !ok {
							d.Set("service_graph_schema_id", tokens[2])
						} else {
							d.Set("service_graph_schema_id", d.Get("service_graph_schema_id"))
						}

						if _, ok := d.GetOk("service_graph_template_name"); !ok {
							d.Set("service_graph_template_name", tokens[4])
						} else {
							d.Set("service_graph_template_name", d.Get("service_graph_template_name"))
						}

						// Template Contract Service Graph Node configurations
						nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
						if err != nil {
							return err
						}
						nodeInterface := d.Get("node_relationship")

						for k := 0; k < nodeCount; k++ {
							relationMap := make(map[string]interface{})
							node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
							if err != nil {
								return err
							}

							probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
							probdRefTokens := strings.Split(probdRef, "/")
							relationMap["provider_connector_bd_name"] = probdRefTokens[len(probdRefTokens)-1]

							if importFlag {
								relationMap["provider_connector_bd_schema_id"] = probdRefTokens[len(probdRefTokens)-5]
							} else if checkNodeAttr(nodeInterface, "provider_connector_bd_schema_id", k) {
								relationMap["provider_connector_bd_schema_id"] = probdRefTokens[len(probdRefTokens)-5]
							} else {
								relationMap["provider_connector_bd_schema_id"] = ""
							}

							if importFlag {
								relationMap["provider_connector_bd_template_name"] = probdRefTokens[len(probdRefTokens)-3]
							} else if checkNodeAttr(nodeInterface, "provider_connector_bd_template_name", k) {
								relationMap["provider_connector_bd_template_name"] = probdRefTokens[len(probdRefTokens)-3]
							} else {
								relationMap["provider_connector_bd_template_name"] = ""
							}

							conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
							conbdRefTokens := strings.Split(conbdRef, "/")
							relationMap["consumer_connector_bd_name"] = conbdRefTokens[len(conbdRefTokens)-1]

							if importFlag {
								relationMap["consumer_connector_bd_schema_id"] = conbdRefTokens[len(conbdRefTokens)-5]
							} else if checkNodeAttr(nodeInterface, "consumer_connector_bd_schema_id", k) {
								relationMap["consumer_connector_bd_schema_id"] = conbdRefTokens[len(conbdRefTokens)-5]
							} else {
								relationMap["consumer_connector_bd_schema_id"] = ""
							}

							if importFlag {
								relationMap["consumer_connector_bd_template_name"] = conbdRefTokens[len(conbdRefTokens)-3]
							} else if checkNodeAttr(nodeInterface, "consumer_connector_bd_template_name", k) {
								relationMap["consumer_connector_bd_template_name"] = conbdRefTokens[len(conbdRefTokens)-3]
							} else {
								relationMap["consumer_connector_bd_template_name"] = ""
							}

							temprelationList = append(temprelationList, relationMap)
						}

						foundTemp = true
						break
					}
				}
			}
		}
		if foundTemp {
			break
		}
	}

	if foundTemp {
		length := len(temprelationList)
		nodeList := make([]interface{}, 0, 1)
		for i := 0; i < length; i++ {
			tempMap := temprelationList[i].(map[string]interface{})

			allMap := make(map[string]interface{})
			allMap["provider_connector_bd_name"] = tempMap["provider_connector_bd_name"]
			allMap["provider_connector_bd_schema_id"] = tempMap["provider_connector_bd_schema_id"]
			allMap["provider_connector_bd_template_name"] = tempMap["provider_connector_bd_template_name"]
			allMap["consumer_connector_bd_name"] = tempMap["consumer_connector_bd_name"]
			allMap["consumer_connector_bd_schema_id"] = tempMap["consumer_connector_bd_schema_id"]
			allMap["consumer_connector_bd_template_name"] = tempMap["consumer_connector_bd_template_name"]

			nodeList = append(nodeList, allMap)
		}
		d.Set("node_relationship", nodeList)
		d.SetId(fmt.Sprintf("/schemas/%s/templates/%s/contracts/%s", schemaID, templateName, contractName))

	} else {
		d.SetId("")
	}
	return nil
}

func PostTemplateContractServiceGraphConfig(action string, d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceGraph := d.Get("service_graph_name").(string)

	serviceGraphRef := make(map[string]interface{})
	serviceGraphRef["serviceGraphName"] = serviceGraph

	if graphSchema, ok := d.GetOk("service_graph_schema_id"); ok {
		serviceGraphRef["schemaId"] = graphSchema.(string)
	} else {
		serviceGraphRef["schemaId"] = schemaID
	}

	if graphTemp, ok := d.GetOk("service_graph_template_name"); ok {
		serviceGraphRef["templateName"] = graphTemp.(string)
	} else {
		serviceGraphRef["templateName"] = TemplateName
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	if err != nil {
		return err
	}

	graphCont, _, err := getTemplateServiceGraph(cont, serviceGraphRef["templateName"].(string), serviceGraph)

	if err != nil {
		return err
	}
	nodeList, err := extractNodes(graphCont)
	if err != nil {
		return err
	}

	tempNodeList := d.Get("node_relationship").([]interface{})

	if len(tempNodeList) != len(nodeList) {
		return fmt.Errorf("Length mismatch between total nodes and total node relation")
	}

	templateNodes, _, err := getServiceNodesRelationshipObject(cont, schemaID, "graphSiteID", TemplateName, nodeList, tempNodeList, serviceGraphRef)

	if err != nil {
		return err
	}
	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)

	tempConGraph := models.NewTemplateContractServiceGraph("replace", tempPath, serviceGraphRef, templateNodes)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	if err != nil {
		return err
	}
	return nil
}
