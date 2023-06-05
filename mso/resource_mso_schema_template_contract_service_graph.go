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

func resourceTemplateContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateContractServiceGraphCreate_new,
		Update: resourceTemplateContractServiceGraphUpdate_new,
		Read:   resourceTemplateContractServiceGraphRead_New,
		Delete: resourceTemplateContractServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceTemplateContractServiceGraphImport,
		},

		Schema: map[string]*schema.Schema{
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
			"service_graph_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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
					},
				},
			},
		},
	}
}

// To import the Template Contract Service Graph
// 1. Schema ID
// 2. Template Name
// 3. Contract Name

// So the ID format should be: /schemas/6475c55904affb217a4cb2b2/templates/T1/contracts/C1
// But directly you can not query Template/Contract/ServiceGraph.

func resourceTemplateContractServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	serviceGraphTokens := strings.Split(d.Id(), "/")
	log.Printf("[DEBUG] %s: ###### tokens Beginning Import, 2nd element: %s", serviceGraphTokens, serviceGraphTokens[1])

	d.Set("schema_id", serviceGraphTokens[2])
	d.Set("template_name", serviceGraphTokens[4])
	d.Set("contract_name", serviceGraphTokens[6])
	// d.Set("service_graph_name", "T2SG")
	setTemplateContractServiceGraphAttrs(d, m, true)
	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceTemplateContractServiceGraphCreate_new(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Template Contract Service Graph")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceGraph := d.Get("service_graph_name").(string)

	// Note
	// if not service_graph_template_name then use template_name
	// if not service_graph_schema_id then use schema_id
	//
	// Sample - serviceGraphRef object
	// "serviceGraphRef": {
	// 	"schemaId": "646301ee04affb217a4aaa84",
	// 	"serviceGraphName": "SG1",
	// 	"templateName": "Template1"
	// }

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

	// Sample - serviceNodesRelationship object
	// "serviceNodesRelationship": [
	//     {
	//         "consumerConnector": {
	//             "bdRef": {
	//                 "bdName": "BD1",
	//                 "schemaId": "6463040f1d0000c0e0f94392",
	//                 "templateName": "Template1"
	//             },
	//             "connectorType": "general"
	//         },
	//         "providerConnector": {
	//             "bdRef": {
	//                 "bdName": "BD2",
	//                 "schemaId": "6463040f1d0000c0e0f94392",
	//                 "templateName": "Template1"
	//             },
	//             "connectorType": "general"
	//         },
	//         "serviceNodeRef": {
	//             "schemaId": "6463040f1d0000c0e0f94392",
	//             "serviceGraphName": "SG1",
	//             "serviceNodeName": "firewall",
	//             "templateName": "Template1"
	//         }
	//     }
	// ]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	if err != nil {
		return err
	}

	graphCont, _, err := getTemplateServiceGraph(cont, serviceGraphRef["templateName"].(string), serviceGraph)

	log.Printf("[DEBUG] ########## after the getTemplateServiceGraph call: %v", graphCont)
	if err != nil {
		return err
	}
	nodeList, err := extractNodes(graphCont)
	if err != nil {
		return err
	}

	tempNodeList := d.Get("node_relationship").([]interface{})

	log.Printf("[DEBUG] ########## before length check tempNodeList: %v, nodeList: %v", tempNodeList, nodeList)

	if len(tempNodeList) != len(nodeList) {
		return fmt.Errorf("Length mismatch between total nodes and total node relation")
	}

	templateNodes := getServiceNodesRelationshipObject(schemaID, TemplateName, nodeList, tempNodeList, serviceGraphRef)

	// Sample path for the template contract service graph: "/templates/Template1/contracts/Contract1/serviceGraphRelationship"
	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("add", tempPath, serviceGraphRef, templateNodes)

	log.Printf("[DEBUG]: ########## Template Contract Service Graph: %v", tempConGraph)
	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Creation Template Contract Service Graph")
	return setTemplateContractServiceGraphAttrs(d, m, false)
}

func resourceTemplateContractServiceGraphUpdate_new(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Template Contract Service Graph")

	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceGraph := d.Get("service_graph_name").(string)
	// Note
	// if not service_graph_template_name then use template_name
	// if not service_graph_schema_id then use schema_id
	//
	// Sample - serviceGraphRef object
	// "serviceGraphRef": {
	// 	"schemaId": "646301ee04affb217a4aaa84",
	// 	"serviceGraphName": "SG1",
	// 	"templateName": "Template1"
	// }

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

	// Sample - serviceNodesRelationship object
	// "serviceNodesRelationship": [
	//     {
	//         "consumerConnector": {
	//             "bdRef": {
	//                 "bdName": "BD1",
	//                 "schemaId": "6463040f1d0000c0e0f94392",
	//                 "templateName": "Template1"
	//             },
	//             "connectorType": "general"
	//         },
	//         "providerConnector": {
	//             "bdRef": {
	//                 "bdName": "BD2",
	//                 "schemaId": "6463040f1d0000c0e0f94392",
	//                 "templateName": "Template1"
	//             },
	//             "connectorType": "general"
	//         },
	//         "serviceNodeRef": {
	//             "schemaId": "6463040f1d0000c0e0f94392",
	//             "serviceGraphName": "SG1",
	//             "serviceNodeName": "firewall",
	//             "templateName": "Template1"
	//         }
	//     }
	// ]

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

	log.Printf("[DEBUG] ########## before length check tempNodeList: %v, nodeList: %v", tempNodeList, nodeList)

	if len(tempNodeList) != len(nodeList) {
		return fmt.Errorf("Length mismatch between total nodes and total node relation")
	}

	templateNodes := getServiceNodesRelationshipObject(schemaID, TemplateName, nodeList, tempNodeList, serviceGraphRef)

	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("replace", tempPath, serviceGraphRef, templateNodes)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Update Template Contract Service Graph")
	return setTemplateContractServiceGraphAttrs(d, m, false)
}

func resourceTemplateContractServiceGraphRead_New(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Template Contract Service Graph")

	serviceGraphTokens := strings.Split(d.Id(), "/")
	log.Printf("[DEBUG] %s: ###### tokens Beginning Import, 2nd element: %s", serviceGraphTokens, serviceGraphTokens[1])

	d.Set("schema_id", serviceGraphTokens[2])
	d.Set("template_name", serviceGraphTokens[4])
	d.Set("contract_name", serviceGraphTokens[6])
	d.Set("service_graph_name", "TSG")
	setTemplateContractServiceGraphAttrs(d, m, false)
	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}

func resourceTemplateContractServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
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

// Used in the create and update
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
			log.Printf("[DEBUG] ############## Inside CSG: %v", templateCont)

			sgCount, err := templateCont.ArrayCount("serviceGraphs")

			log.Printf("[DEBUG] ############## Inside sgCount: %v", sgCount)

			if err != nil {
				return nil, -1, fmt.Errorf("No Service Graph found")
			}

			for j := 0; j < sgCount; j++ {
				sgCont, err := templateCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to get service graph element")
				}

				apiSgName := models.StripQuotes(sgCont.S("name").String())

				log.Printf("[DEBUG] ############## Inside CSG: apiSgName: %s, graphName; %s", apiSgName, graphName)

				if apiSgName == graphName {
					return sgCont, j, nil
				}
			}

		}
	}

	return nil, -1, fmt.Errorf("unable to find service graph")
}

// Used in the create and update
func extractNodes(cont *container.Container) ([]interface{}, error) {
	nodes := make([]interface{}, 0, 1)
	count, err := cont.ArrayCount("serviceNodes")
	if err != nil {
		return nodes, err
	}

	for i := 0; i < count; i++ {
		node, err := cont.ArrayElement(i, "serviceNodes", "name")
		if err != nil {
			return nodes, err
		}

		nodes = append(nodes, models.StripQuotes(node.String()))
	}
	return nodes, nil
}

func checkNodeAttr(object interface{}, attrName string, index int) bool {
	objList := object.([]interface{})
	instance := objList[index].(map[string]interface{})
	return instance[attrName] != ""
}

func getServiceNodesRelationshipObject(schemaID, templateName string, nodeList, tempNodeList []interface{}, serviceGraphRef map[string]interface{}) []interface{} {

	templateNodes := make([]interface{}, 0)

	for i := 0; i < len(tempNodeList); i++ {
		node := tempNodeList[i].(map[string]interface{})

		log.Printf("[DEBUG]: ######### Before assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)

		tempnodeRef := make(map[string]interface{})

		tempnodeRef["schemaId"] = serviceGraphRef["schemaId"]
		tempnodeRef["serviceGraphName"] = serviceGraphRef["serviceGraphName"]
		tempnodeRef["templateName"] = serviceGraphRef["templateName"]
		tempnodeRef["serviceNodeName"] = nodeList[i].(string)

		log.Printf("[DEBUG]: ######### After assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)

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
	return templateNodes
}

func setTemplateContractServiceGraphAttrs(d *schema.ResourceData, m interface{}, importFlag bool) error {
	log.Printf("[DEBUG] Begining of setTemplateContractServiceGraphAttrs - Read Template Contract Service Graph")
	msoClient := m.(*client.Client)
	foundTemp := false
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	// serviceGraph := d.Get("service_graph_name").(string)

	// serviceGraphRef := make(map[string]interface{})
	// serviceGraphRef["serviceGraphName"] = serviceGraph
	// if graphSchema, ok := d.GetOk("service_graph_schema_id"); ok {
	// 	serviceGraphRef["schemaId"] = graphSchema.(string)
	// } else {
	// 	serviceGraphRef["schemaId"] = schemaID
	// }

	// if graphTemp, ok := d.GetOk("service_graph_template_name"); ok {
	// 	serviceGraphRef["templateName"] = graphTemp.(string)
	// } else {
	// 	serviceGraphRef["templateName"] = templateName
	// }

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No templates found")
	}

	temprelationList := make([]interface{}, 0)
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
						// Template Contract Service Graph Details - begins
						graphRelation := contractCont.S("serviceGraphRelationship")
						graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
						tokens := strings.Split(graphRef, "/")
						log.Printf("[DEBUG]: serviceGraphRef values: %v, second element: %s", tokens, tokens[1])

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

						// if _, ok := d.GetOk("service_graph_name"); ok {
						// 	d.Set("service_graph_name", d.Get("service_graph_name"))
						// } else {
						// 	d.Set("service_graph_name", tokens[6])
						// }

						// Template Contract Service Graph Details - ends

						// Template Contract Service Graph Node Details - begins

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
							// nodeRef := models.StripQuotes(node.S("serviceNodeRef").String())
							// tokensNode := strings.Split(nodeRef, "/")
							// relationMap["node_name"] = tokensNode[len(tokensNode)-1]

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
						// Template Contract Service Graph Node Details - ends

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
	log.Printf("[DEBUG]: ####### Before setting node_relationship temprelationList: %v", temprelationList)

	if foundTemp {
		node_set_err := d.Set("node_relationship", temprelationList)
		log.Printf("[DEBUG]: ####### After setting node_relationship: %v", node_set_err)
		d.SetId(fmt.Sprintf("/schemas/%s/templates/%s/contracts/%s", schemaID, templateName, contractName))
	} else {
		d.SetId("")
	}

	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}
