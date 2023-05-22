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
		Read:   resourceTemplateContractServiceGraphRead,
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

			"site_id": &schema.Schema{ // Not required - marked by sabari
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

			"service_graph_site_id": &schema.Schema{ // Not required - marked by sabari
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

			// Note
			// if not service_graph_template_name then use template_name
			// if not service_graph_schema_id then use schema_id
			//

			"node_relationship": &schema.Schema{ // Not required - marked by sabari
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

						// "provider_connector_cluster_interface": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Required:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "provider_connector_redirect_policy_tenant": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// 	Computed:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "provider_connector_redirect_policy": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// 	Computed:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "provider_subnet_ips": &schema.Schema{
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Computed: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// },

						// "consumer_connector_cluster_interface": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Required:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "consumer_connector_redirect_policy_tenant": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// 	Computed:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "consumer_connector_redirect_policy": &schema.Schema{
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// 	Computed:     true,
						// 	ValidateFunc: validation.StringLenBetween(1, 1000),
						// },

						// "consumer_subnet_ips": &schema.Schema{
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Computed: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// },
					},
				},
			},
			// CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// 	// workaround because ValidateFunc is not allowed on TypeSet
			// 	// if diff.HasChange("relation_l3ext_rs_subnet_to_profile") {
			// 	// 	_, new := diff.GetChange("relation_l3ext_rs_subnet_to_profile")
			// 	// 	// validate that a direction type ( import/export ) is not defined more than once
			// 	// 	err := validateDirection(new.(*schema.Set).List())
			// 	// 	if err != nil {
			// 	// 		return err
			// 	// 	}
			// 	// 	// validate that dn and name are not both defined
			// 	// 	err2 := validateDnAndName(new.(*schema.Set).List())
			// 	// 	if err2 != nil {
			// 	// 		return err2
			// 	// 	}
			// 	// }
			// 	return nil
			// },
		},
	}
}

func resourceTemplateContractServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	foundTemp := false
	foundSite := false

	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	// siteID := get_attribute[2]
	templateName := get_attribute[2]
	contractName := get_attribute[4]
	serviceGraph := get_attribute[6]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No templates found")
	}

	temprelationList := make([]interface{}, 0, 1)
	for i := 0; i < tempCount; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, fmt.Errorf("Error in fetch of template")
		}
		template := models.StripQuotes(tempCont.S("name").String())
		if templateName == template {
			contractCount, err := tempCont.ArrayCount("contracts")
			if err != nil {
				return nil, fmt.Errorf("No contracts found")
			}

			for j := 0; j < contractCount; j++ {
				contractCont, err := tempCont.ArrayElement(j, "contracts")
				if err != nil {
					return nil, fmt.Errorf("Error fetching contract")
				}
				conName := models.StripQuotes(contractCont.S("name").String())
				if conName == contractName {
					if !contractCont.Exists("serviceGraphRelationship") {
						return nil, fmt.Errorf("No service graph found")
					} else {

						graphRelation := contractCont.S("serviceGraphRelationship")

						graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
						tokens := strings.Split(graphRef, "/")
						d.Set("service_graph_name", tokens[len(tokens)-1])

						nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
						if err != nil {
							return nil, err
						}
						for k := 0; k < nodeCount; k++ {
							relationMap := make(map[string]interface{})
							node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
							if err != nil {
								return nil, err
							}
							nodeRef := models.StripQuotes(node.S("serviceNodeRef").String())
							tokensNode := strings.Split(nodeRef, "/")
							relationMap["node_name"] = tokensNode[len(tokensNode)-1]

							probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
							probdRefTokens := strings.Split(probdRef, "/")
							relationMap["provider_connector_bd_name"] = probdRefTokens[len(probdRefTokens)-1]

							conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
							conbdRefTokens := strings.Split(conbdRef, "/")
							relationMap["consumer_connector_bd_name"] = conbdRefTokens[len(conbdRefTokens)-1]

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

	// siterelationList := make([]interface{}, 0, 1)
	// siteCount, err := cont.ArrayCount("sites")
	// if err != nil {
	// 	return nil, fmt.Errorf("No sites found")
	// }
	// for i := 0; i < siteCount; i++ {
	// 	siteCont, err := cont.ArrayElement(i, "sites")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Error fetching site")
	// 	}

	// 	site := models.StripQuotes(siteCont.S("siteId").String())
	// 	temp := models.StripQuotes(siteCont.S("templateName").String())
	// 	if siteID == site && temp == templateName {
	// 		contractCount, err := siteCont.ArrayCount("contracts")
	// 		if err != nil {
	// 			return nil, fmt.Errorf("No contracts found in site")
	// 		}

	// 		for j := 0; j < contractCount; j++ {
	// 			contractCont, err := siteCont.ArrayElement(j, "contracts")
	// 			if err != nil {
	// 				return nil, fmt.Errorf("Error fetching contract from site")
	// 			}

	// 			conRef := models.StripQuotes(contractCont.S("contractRef").String())
	// 			conTokens := strings.Split(conRef, "/")
	// 			conName := conTokens[len(conTokens)-1]
	// 			if conName == contractName {
	// 				if !contractCont.Exists("serviceGraphRelationship") {
	// 					return nil, fmt.Errorf("No service graph found")
	// 				} else {
	// 					graphRelation := contractCont.S("serviceGraphRelationship")

	// 					nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
	// 					if err != nil {
	// 						return nil, err
	// 					}
	// 					for k := 0; k < nodeCount; k++ {
	// 						relationMap := make(map[string]interface{})
	// 						node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
	// 						if err != nil {
	// 							return nil, err
	// 						}

	// 						relationMap["provider_connector_cluster_interface"] = models.StripQuotes(node.S("providerConnector", "clusterInterface", "dn").String())

	// 						if node.Exists("providerConnector", "redirectPolicy", "dn") {
	// 							relationMap["provider_connector_redirect_policy"] = models.StripQuotes(node.S("providerConnector", "redirectPolicy", "dn").String())
	// 						}

	// 						if node.Exists("providerConnector", "subnets") {
	// 							subCounts, err := node.ArrayCount("providerConnector", "subnets")
	// 							if err != nil {
	// 								return nil, err
	// 							}
	// 							subList := make([]interface{}, 0, 1)
	// 							for l := 0; l < subCounts; l++ {
	// 								subnet, err := node.ArrayElement(l, "providerConnector", "subnets", "ip")
	// 								if err != nil {
	// 									return nil, err
	// 								}
	// 								subList = append(subList, models.StripQuotes(subnet.String()))
	// 							}
	// 							relationMap["provider_subnet_ips"] = subList
	// 						}

	// 						relationMap["consumer_connector_cluster_interface"] = models.StripQuotes(node.S("consumerConnector", "clusterInterface", "dn").String())

	// 						if node.Exists("consumerConnector", "redirectPolicy", "dn") {
	// 							relationMap["consumer_connector_redirect_policy"] = models.StripQuotes(node.S("consumerConnector", "redirectPolicy", "dn").String())
	// 						}

	// 						if node.Exists("consumerConnector", "subnets") {
	// 							subCounts, err := node.ArrayCount("consumerConnector", "subnets")
	// 							if err != nil {
	// 								return nil, err
	// 							}
	// 							subList := make([]interface{}, 0, 1)
	// 							for l := 0; l < subCounts; l++ {
	// 								subnet, err := node.ArrayElement(l, "consumerConnector", "subnets", "ip")
	// 								if err != nil {
	// 									return nil, err
	// 								}
	// 								subList = append(subList, models.StripQuotes(subnet.String()))
	// 							}
	// 							relationMap["consumer_subnet_ips"] = subList
	// 						}

	// 						siterelationList = append(siterelationList, relationMap)
	// 					}
	// 					foundSite = true
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if foundSite {
	// 		break
	// 	}
	// }

	foundSite = true
	if foundSite && foundTemp {
		length := len(temprelationList)
		nodeList := make([]interface{}, 0, 1)
		for i := 0; i < length; i++ {
			tempMap := temprelationList[i].(map[string]interface{})
			// siteMap := siterelationList[i].(map[string]interface{})

			allMap := make(map[string]interface{})
			allMap["provider_connector_bd_name"] = tempMap["provider_connector_bd_name"]
			allMap["consumer_connector_bd_name"] = tempMap["consumer_connector_bd_name"]

			// tp := strings.Split(siteMap["provider_connector_cluster_interface"].(string), "/")
			// token := strings.Split(tp[len(tp)-1], "-")
			// allMap["provider_connector_cluster_interface"] = token[1]

			// tp = strings.Split(siteMap["consumer_connector_cluster_interface"].(string), "/")
			// token = strings.Split(tp[len(tp)-1], "-")
			// allMap["consumer_connector_cluster_interface"] = token[1]

			// if siteMap["provider_connector_redirect_policy"] != nil {
			// 	tp := strings.Split(siteMap["provider_connector_redirect_policy"].(string), "/")
			// 	token1 := strings.Split(tp[1], "-")
			// 	allMap["provider_connector_redirect_policy_tenant"] = token1[1]

			// 	token2 := strings.Split(tp[len(tp)-1], "-")
			// 	allMap["provider_connector_redirect_policy"] = token2[1]
			// }
			// if siteMap["consumer_connector_redirect_policy"] != nil {
			// 	tp := strings.Split(siteMap["consumer_connector_redirect_policy"].(string), "/")
			// 	token1 := strings.Split(tp[1], "-")
			// 	allMap["consumer_connector_redirect_policy_tenant"] = token1[1]

			// 	token2 := strings.Split(tp[len(tp)-1], "-")
			// 	allMap["consumer_connector_redirect_policy"] = token2[1]
			// }

			// if siteMap["provider_subnet_ips"] != nil {
			// 	allMap["provider_subnet_ips"] = siteMap["provider_subnet_ips"]
			// }
			// if siteMap["consumer_subnet_ips"] != nil {
			// 	allMap["consumer_subnet_ips"] = siteMap["consumer_subnet_ips"]
			// }

			nodeList = append(nodeList, allMap)
		}
		d.Set("schema_id", schemaId)
		// d.Set("site_id", siteID)
		d.Set("template_name", templateName)
		d.Set("node_relationship", nodeList)
		d.Set("contract_name", contractName)

		if d.Get("service_graph_name") == serviceGraph {
			d.SetId(serviceGraph)
		} else {
			d.SetId("")
			return nil, fmt.Errorf("No service graph found for given name")
		}
	} else {
		d.SetId("")
		return nil, fmt.Errorf("No service graph found for given name")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceTemplateContractServiceGraphCreate_new(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Template Contract Service Graph")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
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

	// var graphSiteID string
	// if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
	// 	graphSiteID = graphSite.(string)
	// } else {
	// 	graphSiteID = siteID
	// }

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

	templateNodes := make([]interface{}, 0, 1)
	// siteNodes := make([]interface{}, 0, 1)

	tempNodeList := d.Get("node_relationship").([]interface{})

	log.Printf("[DEBUG] ########## before length check tempNodeList: %v, nodeList: %v", tempNodeList, nodeList)

	if len(tempNodeList) != len(nodeList) {
		return fmt.Errorf("Length mismatch between total nodes and total node relation")
	}

	// log.Printf("[DEBUG]: ####### inside create tempNodeList: %v", tempNodeList)
	for i := 0; i < len(tempNodeList); i++ {
		node := tempNodeList[i].(map[string]interface{})

		log.Printf("[DEBUG]: ######### Before assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)
		// tempnodeRef := serviceGraphRef

		tempnodeRef := make(map[string]interface{})

		tempnodeRef["schemaId"] = serviceGraphRef["schemaId"]
		tempnodeRef["serviceGraphName"] = serviceGraphRef["serviceGraphName"]
		tempnodeRef["templateName"] = serviceGraphRef["templateName"]
		tempnodeRef["serviceNodeName"] = nodeList[i].(string)

		log.Printf("[DEBUG]: ######### After assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)

		// 	sitegraphCont, _, err := getSiteServiceGraph(cont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), graphSiteID, tempnodeRef["serviceGraphName"].(string))
		// 	if err != nil {
		// 		return err
		// 	}

		// 	siteNodeCont, _, err := getSiteServiceNode(sitegraphCont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), tempnodeRef["serviceGraphName"].(string), nodeList[i].(string))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	dn := models.StripQuotes(siteNodeCont.S("device", "dn").String())

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
			bdRef["templateName"] = TemplateName
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
			conbdRef["templateName"] = TemplateName
		}
		conbdRef["bdName"] = node["consumer_connector_bd_name"].(string)
		tempconConnector["bdRef"] = conbdRef

		tempnodeMap := make(map[string]interface{})
		tempnodeMap["serviceNodeRef"] = tempnodeRef
		tempnodeMap["providerConnector"] = tempproConnector
		tempnodeMap["consumerConnector"] = tempconConnector

		templateNodes = append(templateNodes, tempnodeMap)

		// Sabari end

		// 	/*providerConnector*/
		// 	proConnector := make(map[string]interface{})
		// 	proClusterInterface := make(map[string]interface{})
		// 	proClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["provider_connector_cluster_interface"].(string))
		// 	proConnector["clusterInterface"] = proClusterInterface

		// 	if node["provider_connector_redirect_policy"] != "" {
		// 		if node["provider_connector_redirect_policy_tenant"] == "" {
		// 			return fmt.Errorf("provider redirect policy tenant is required")
		// 		}
		// 		proRedPolicy := make(map[string]interface{})
		// 		proRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["provider_connector_redirect_policy_tenant"].(string), node["provider_connector_redirect_policy"].(string))
		// 		proConnector["redirectPolicy"] = proRedPolicy
		// 	}

		// 	if node["provider_subnet_ips"] != nil {
		// 		ips := node["provider_subnet_ips"].([]interface{})
		// 		prosubnets := make([]interface{}, 0, 1)
		// 		for _, ip := range ips {
		// 			subnet := make(map[string]interface{})
		// 			subnet["ip"] = ip.(string)

		// 			prosubnets = append(prosubnets, subnet)
		// 		}
		// 		proConnector["subnets"] = prosubnets
		// 	}

		// 	/*consumerConnector*/
		// 	conConnector := make(map[string]interface{})
		// 	conClusterInterface := make(map[string]interface{})
		// 	conClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["consumer_connector_cluster_interface"].(string))
		// 	conConnector["clusterInterface"] = conClusterInterface

		// 	if node["consumer_connector_redirect_policy"] != "" {
		// 		if node["consumer_connector_redirect_policy_tenant"] == "" {
		// 			return fmt.Errorf("consumer redirect policy tenant is required")
		// 		}
		// 		conRedPolicy := make(map[string]interface{})
		// 		conRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["consumer_connector_redirect_policy_tenant"].(string), node["consumer_connector_redirect_policy"].(string))
		// 		conConnector["redirectPolicy"] = conRedPolicy
		// 	}

		// 	if node["consumer_subnet_ips"] != nil {
		// 		ips := node["consumer_subnet_ips"].([]interface{})
		// 		consubnets := make([]interface{}, 0, 1)
		// 		for _, ip := range ips {
		// 			subnet := make(map[string]interface{})
		// 			subnet["ip"] = ip.(string)

		// 			consubnets = append(consubnets, subnet)
		// 		}
		// 		conConnector["subnets"] = consubnets
		// 	}

		// 	nodeMap := make(map[string]interface{})
		// 	nodeMap["serviceNodeRef"] = tempnodeRef
		// 	nodeMap["providerConnector"] = proConnector
		// 	nodeMap["consumerConnector"] = conConnector

		// 	siteNodes = append(siteNodes, nodeMap)
	}

	// sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, TemplateName, contractName)
	// siteConGraph := models.NewSiteContractServiceGraph("add", sitePath, serviceGraphRef, siteNodes)

	// Sample path for the template contract service graph: "/templates/Template1/contracts/Contract1/serviceGraphRelationship"
	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("add", tempPath, serviceGraphRef, templateNodes)

	// _, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph, siteConGraph)

	log.Printf("[DEBUG]: ########## Template Contract Service Graph: %v", tempConGraph)
	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Creation Template Contract Service Graph")
	return resourceTemplateContractServiceGraphRead(d, m)
}

func resourceTemplateContractServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Template Contract Service Graph")
	// msoClient := m.(*client.Client)

	// schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
	// TemplateName := d.Get("template_name").(string)
	// contractName := d.Get("contract_name").(string)
	// serviceGraph := d.Get("service_graph_name").(string)

	// // Note
	// // if not service_graph_template_name then use template_name
	// // if not service_graph_schema_id then use schema_id
	// //

	// serviceGraphRef := make(map[string]interface{})
	// if graphSchema, ok := d.GetOk("service_graph_schema_id"); ok {
	// 	serviceGraphRef["schemaId"] = graphSchema.(string)
	// } else {
	// 	serviceGraphRef["schemaId"] = schemaID
	// }

	// if graphTemp, ok := d.GetOk("service_graph_template_name"); ok {
	// 	serviceGraphRef["templateName"] = graphTemp.(string)
	// } else {
	// 	serviceGraphRef["templateName"] = TemplateName
	// }
	// serviceGraphRef["serviceGraphName"] = serviceGraph

	// var graphSiteID string
	// if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
	// 	graphSiteID = graphSite.(string)
	// } else {
	// 	graphSiteID = siteID
	// }

	// cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	// if err != nil {
	// 	return err
	// }

	// graphCont, _, err := getTemplateServiceGraph(cont, serviceGraphRef["templateName"].(string), serviceGraph)

	// log.Printf("[DEBUG] ########## after the getTemplateServiceGraph call: %v", graphCont)
	// if err != nil {
	// 	return err
	// }
	// nodeList, err := extractNodes(graphCont)
	// if err != nil {
	// 	return err
	// }

	// templateNodes := make([]interface{}, 0, 1)
	// siteNodes := make([]interface{}, 0, 1)

	// tempNodeList := d.Get("node_relationship").([]interface{})

	// log.Printf("[DEBUG] ########## before length check tempNodeList: %v, nodeList: %v", tempNodeList, nodeList)

	// if len(tempNodeList) != len(nodeList) {
	// 	return fmt.Errorf("Length mismatch between total nodes and total node relation")
	// }

	// log.Printf("[DEBUG]: ####### inside create tempNodeList: %v", tempNodeList)
	// for i := 0; i < len(tempNodeList); i++ {
	// 	node := tempNodeList[i].(map[string]interface{})

	// 	tempnodeRef := serviceGraphRef
	// 	tempnodeRef["serviceNodeName"] = nodeList[i].(string)

	// 	sitegraphCont, _, err := getSiteServiceGraph(cont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), graphSiteID, tempnodeRef["serviceGraphName"].(string))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	siteNodeCont, _, err := getSiteServiceNode(sitegraphCont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), tempnodeRef["serviceGraphName"].(string), nodeList[i].(string))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	dn := models.StripQuotes(siteNodeCont.S("device", "dn").String())

	// 	tempproConnector := make(map[string]interface{})
	// 	tempproConnector["connectorType"] = "general"
	// 	bdRef := make(map[string]interface{})
	// 	if node["provider_connector_bd_schema_id"] != "" {
	// 		bdRef["schemaId"] = node["provider_connector_bd_schema_id"].(string)
	// 	} else {
	// 		bdRef["schemaId"] = schemaID
	// 	}

	// 	if node["provider_connector_bd_template_name"] != "" {
	// 		bdRef["templateName"] = node["provider_connector_bd_template_name"].(string)
	// 	} else {
	// 		bdRef["templateName"] = TemplateName
	// 	}
	// 	bdRef["bdName"] = node["provider_connector_bd_name"].(string)
	// 	tempproConnector["bdRef"] = bdRef

	// 	tempconConnector := make(map[string]interface{})
	// 	tempconConnector["connectorType"] = "general"
	// 	conbdRef := make(map[string]interface{})
	// 	if node["consumer_connector_bd_schema_id"] != "" {
	// 		conbdRef["schemaId"] = node["consumer_connector_bd_schema_id"].(string)
	// 	} else {
	// 		conbdRef["schemaId"] = schemaID
	// 	}

	// 	if node["consumer_connector_bd_template_name"] != "" {
	// 		conbdRef["templateName"] = node["consumer_connector_bd_template_name"].(string)
	// 	} else {
	// 		conbdRef["templateName"] = TemplateName
	// 	}
	// 	conbdRef["bdName"] = node["consumer_connector_bd_name"].(string)
	// 	tempconConnector["bdRef"] = conbdRef

	// 	tempnodeMap := make(map[string]interface{})
	// 	tempnodeMap["serviceNodeRef"] = tempnodeRef
	// 	tempnodeMap["providerConnector"] = tempproConnector
	// 	tempnodeMap["consumerConnector"] = tempconConnector

	// 	templateNodes = append(templateNodes, tempnodeMap)

	// 	/*providerConnector*/
	// 	proConnector := make(map[string]interface{})
	// 	proClusterInterface := make(map[string]interface{})
	// 	proClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["provider_connector_cluster_interface"].(string))
	// 	proConnector["clusterInterface"] = proClusterInterface

	// 	if node["provider_connector_redirect_policy"] != "" {
	// 		if node["provider_connector_redirect_policy_tenant"] == "" {
	// 			return fmt.Errorf("provider redirect policy tenant is required")
	// 		}
	// 		proRedPolicy := make(map[string]interface{})
	// 		proRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["provider_connector_redirect_policy_tenant"].(string), node["provider_connector_redirect_policy"].(string))
	// 		proConnector["redirectPolicy"] = proRedPolicy
	// 	}

	// 	if node["provider_subnet_ips"] != nil {
	// 		ips := node["provider_subnet_ips"].([]interface{})
	// 		prosubnets := make([]interface{}, 0, 1)
	// 		for _, ip := range ips {
	// 			subnet := make(map[string]interface{})
	// 			subnet["ip"] = ip.(string)

	// 			prosubnets = append(prosubnets, subnet)
	// 		}
	// 		proConnector["subnets"] = prosubnets
	// 	}

	// 	/*consumerConnector*/
	// 	conConnector := make(map[string]interface{})
	// 	conClusterInterface := make(map[string]interface{})
	// 	conClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["consumer_connector_cluster_interface"].(string))
	// 	conConnector["clusterInterface"] = conClusterInterface

	// 	if node["consumer_connector_redirect_policy"] != "" {
	// 		if node["consumer_connector_redirect_policy_tenant"] == "" {
	// 			return fmt.Errorf("consumer redirect policy tenant is required")
	// 		}
	// 		conRedPolicy := make(map[string]interface{})
	// 		conRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["consumer_connector_redirect_policy_tenant"].(string), node["consumer_connector_redirect_policy"].(string))
	// 		conConnector["redirectPolicy"] = conRedPolicy
	// 	}

	// 	if node["consumer_subnet_ips"] != nil {
	// 		ips := node["consumer_subnet_ips"].([]interface{})
	// 		consubnets := make([]interface{}, 0, 1)
	// 		for _, ip := range ips {
	// 			subnet := make(map[string]interface{})
	// 			subnet["ip"] = ip.(string)

	// 			consubnets = append(consubnets, subnet)
	// 		}
	// 		conConnector["subnets"] = consubnets
	// 	}

	// 	nodeMap := make(map[string]interface{})
	// 	nodeMap["serviceNodeRef"] = tempnodeRef
	// 	nodeMap["providerConnector"] = proConnector
	// 	nodeMap["consumerConnector"] = conConnector

	// 	siteNodes = append(siteNodes, nodeMap)
	// }

	// tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	// tempConGraph := models.NewTemplateContractServiceGraph("add", tempPath, serviceGraphRef, templateNodes)

	// sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, TemplateName, contractName)
	// siteConGraph := models.NewSiteContractServiceGraph("add", sitePath, serviceGraphRef, siteNodes)

	// _, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph, siteConGraph)
	// if err != nil {
	// 	return err
	// }

	log.Printf("[DEBUG] Completed Creation Template Contract Service Graph")
	return resourceTemplateContractServiceGraphRead(d, m)
}

func resourceTemplateContractServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Template Contract Service Graph")

	// msoClient := m.(*client.Client)

	// schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
	// TemplateName := d.Get("template_name").(string)
	// contractName := d.Get("contract_name").(string)
	// serviceGraph := d.Get("service_graph_name").(string)

	// serviceGraphRef := make(map[string]interface{})
	// if graphSchema, ok := d.GetOk("service_graph_schema_id"); ok {
	// 	serviceGraphRef["schemaId"] = graphSchema.(string)
	// } else {
	// 	serviceGraphRef["schemaId"] = schemaID
	// }

	// if graphTemp, ok := d.GetOk("service_graph_template_name"); ok {
	// 	serviceGraphRef["templateName"] = graphTemp.(string)
	// } else {
	// 	serviceGraphRef["templateName"] = TemplateName
	// }
	// serviceGraphRef["serviceGraphName"] = serviceGraph

	// var graphSiteID string
	// if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
	// 	graphSiteID = graphSite.(string)
	// } else {
	// 	graphSiteID = siteID
	// }

	// cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	// if err != nil {
	// 	return err
	// }

	// graphCont, _, err := getTemplateServiceGraph(cont, serviceGraphRef["templateName"].(string), serviceGraph)
	// if err != nil {
	// 	return err
	// }
	// nodeList, err := extractNodes(graphCont)
	// if err != nil {
	// 	return err
	// }

	// templateNodes := make([]interface{}, 0, 1)
	// siteNodes := make([]interface{}, 0, 1)

	// tempNodeList := d.Get("node_relationship").([]interface{})
	// if len(tempNodeList) != len(nodeList) {
	// 	return fmt.Errorf("Length mismatch between total nodes and total node relation")
	// }
	// for i := 0; i < len(tempNodeList); i++ {
	// 	node := tempNodeList[i].(map[string]interface{})

	// 	tempnodeRef := serviceGraphRef
	// 	tempnodeRef["serviceNodeName"] = nodeList[i].(string)

	// 	sitegraphCont, _, err := getSiteServiceGraph(cont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), graphSiteID, tempnodeRef["serviceGraphName"].(string))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	siteNodeCont, _, err := getSiteServiceNode(sitegraphCont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), tempnodeRef["serviceGraphName"].(string), nodeList[i].(string))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	dn := models.StripQuotes(siteNodeCont.S("device", "dn").String())

	// 	tempproConnector := make(map[string]interface{})
	// 	tempproConnector["connectorType"] = "general"
	// 	bdRef := make(map[string]interface{})
	// 	if node["provider_connector_bd_schema_id"] != "" {
	// 		bdRef["schemaId"] = node["provider_connector_bd_schema_id"].(string)
	// 	} else {
	// 		bdRef["schemaId"] = schemaID
	// 	}

	// 	if node["provider_connector_bd_template_name"] != "" {
	// 		bdRef["templateName"] = node["provider_connector_bd_template_name"].(string)
	// 	} else {
	// 		bdRef["templateName"] = TemplateName
	// 	}
	// 	bdRef["bdName"] = node["provider_connector_bd_name"].(string)
	// 	tempproConnector["bdRef"] = bdRef

	// 	tempconConnector := make(map[string]interface{})
	// 	tempconConnector["connectorType"] = "general"
	// 	conbdRef := make(map[string]interface{})
	// 	if node["consumer_connector_bd_schema_id"] != "" {
	// 		conbdRef["schemaId"] = node["consumer_connector_bd_schema_id"].(string)
	// 	} else {
	// 		conbdRef["schemaId"] = schemaID
	// 	}

	// 	if node["consumer_connector_bd_template_name"] != "" {
	// 		conbdRef["templateName"] = node["consumer_connector_bd_template_name"].(string)
	// 	} else {
	// 		conbdRef["templateName"] = TemplateName
	// 	}
	// 	conbdRef["bdName"] = node["consumer_connector_bd_name"].(string)
	// 	tempconConnector["bdRef"] = conbdRef

	// 	tempnodeMap := make(map[string]interface{})
	// 	tempnodeMap["serviceNodeRef"] = tempnodeRef
	// 	tempnodeMap["providerConnector"] = tempproConnector
	// 	tempnodeMap["consumerConnector"] = tempconConnector

	// 	templateNodes = append(templateNodes, tempnodeMap)

	// 	/*providerConnector*/
	// 	proConnector := make(map[string]interface{})
	// 	proClusterInterface := make(map[string]interface{})
	// 	proClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["provider_connector_cluster_interface"].(string))
	// 	proConnector["clusterInterface"] = proClusterInterface

	// 	if node["provider_connector_redirect_policy"] != "" {
	// 		if node["provider_connector_redirect_policy_tenant"] == "" {
	// 			return fmt.Errorf("provider redirect policy tenant is required")
	// 		}
	// 		proRedPolicy := make(map[string]interface{})
	// 		proRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["provider_connector_redirect_policy_tenant"].(string), node["provider_connector_redirect_policy"].(string))
	// 		proConnector["redirectPolicy"] = proRedPolicy
	// 	}

	// 	if node["provider_subnet_ips"] != nil {
	// 		ips := node["provider_subnet_ips"].([]interface{})
	// 		prosubnets := make([]interface{}, 0, 1)
	// 		for _, ip := range ips {
	// 			subnet := make(map[string]interface{})
	// 			subnet["ip"] = ip.(string)

	// 			prosubnets = append(prosubnets, subnet)
	// 		}
	// 		proConnector["subnets"] = prosubnets
	// 	}

	// 	/*consumerConnector*/
	// 	conConnector := make(map[string]interface{})
	// 	conClusterInterface := make(map[string]interface{})
	// 	conClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["consumer_connector_cluster_interface"].(string))
	// 	conConnector["clusterInterface"] = conClusterInterface

	// 	if node["consumer_connector_redirect_policy"] != "" {
	// 		if node["consumer_connector_redirect_policy_tenant"] == "" {
	// 			return fmt.Errorf("consumer redirect policy tenant is required")
	// 		}
	// 		conRedPolicy := make(map[string]interface{})
	// 		conRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["consumer_connector_redirect_policy_tenant"].(string), node["consumer_connector_redirect_policy"].(string))
	// 		conConnector["redirectPolicy"] = conRedPolicy
	// 	}

	// 	if node["consumer_subnet_ips"] != nil {
	// 		ips := node["consumer_subnet_ips"].([]interface{})
	// 		consubnets := make([]interface{}, 0, 1)
	// 		for _, ip := range ips {
	// 			subnet := make(map[string]interface{})
	// 			subnet["ip"] = ip.(string)

	// 			consubnets = append(consubnets, subnet)
	// 		}
	// 		conConnector["subnets"] = consubnets
	// 	}

	// 	nodeMap := make(map[string]interface{})
	// 	nodeMap["serviceNodeRef"] = tempnodeRef
	// 	nodeMap["providerConnector"] = proConnector
	// 	nodeMap["consumerConnector"] = conConnector

	// 	siteNodes = append(siteNodes, nodeMap)
	// }

	// tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	// tempConGraph := models.NewTemplateContractServiceGraph("replace", tempPath, serviceGraphRef, templateNodes)

	// sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, TemplateName, contractName)
	// siteConGraph := models.NewSiteContractServiceGraph("replace", sitePath, serviceGraphRef, siteNodes)

	// _, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph, siteConGraph)
	// if err != nil {
	// 	return err
	// }

	log.Printf("[DEBUG] Completed Update Template Contract Service Graph")
	return resourceTemplateContractServiceGraphRead(d, m)
}

func resourceTemplateContractServiceGraphUpdate_new(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Template Contract Service Graph")

	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceGraph := d.Get("service_graph_name").(string)

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

	// var graphSiteID string
	// if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
	// 	graphSiteID = graphSite.(string)
	// } else {
	// 	graphSiteID = siteID
	// }

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

	templateNodes := make([]interface{}, 0, 1)
	// siteNodes := make([]interface{}, 0, 1)

	for i := 0; i < len(tempNodeList); i++ {
		node := tempNodeList[i].(map[string]interface{})

		log.Printf("[DEBUG]: ######### Before assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)
		// tempnodeRef := serviceGraphRef

		tempnodeRef := make(map[string]interface{})

		tempnodeRef["schemaId"] = serviceGraphRef["schemaId"]
		tempnodeRef["serviceGraphName"] = serviceGraphRef["serviceGraphName"]
		tempnodeRef["templateName"] = serviceGraphRef["templateName"]
		tempnodeRef["serviceNodeName"] = nodeList[i].(string)

		log.Printf("[DEBUG]: ######### After assigning to tempnodeRef the serviceGraphRef value: %v", serviceGraphRef)

		// sitegraphCont, _, err := getSiteServiceGraph(cont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), graphSiteID, tempnodeRef["serviceGraphName"].(string))
		// if err != nil {
		// 	return err
		// }

		// siteNodeCont, _, err := getSiteServiceNode(sitegraphCont, tempnodeRef["schemaId"].(string), tempnodeRef["templateName"].(string), tempnodeRef["serviceGraphName"].(string), nodeList[i].(string))
		// if err != nil {
		// 	return err
		// }
		// dn := models.StripQuotes(siteNodeCont.S("device", "dn").String())

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
			bdRef["templateName"] = TemplateName
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
			conbdRef["templateName"] = TemplateName
		}
		conbdRef["bdName"] = node["consumer_connector_bd_name"].(string)
		tempconConnector["bdRef"] = conbdRef

		tempnodeMap := make(map[string]interface{})
		tempnodeMap["serviceNodeRef"] = tempnodeRef
		tempnodeMap["providerConnector"] = tempproConnector
		tempnodeMap["consumerConnector"] = tempconConnector

		templateNodes = append(templateNodes, tempnodeMap)

		// 	/*providerConnector*/
		// 	proConnector := make(map[string]interface{})
		// 	proClusterInterface := make(map[string]interface{})
		// 	proClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["provider_connector_cluster_interface"].(string))
		// 	proConnector["clusterInterface"] = proClusterInterface

		// 	if node["provider_connector_redirect_policy"] != "" {
		// 		if node["provider_connector_redirect_policy_tenant"] == "" {
		// 			return fmt.Errorf("provider redirect policy tenant is required")
		// 		}
		// 		proRedPolicy := make(map[string]interface{})
		// 		proRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["provider_connector_redirect_policy_tenant"].(string), node["provider_connector_redirect_policy"].(string))
		// 		proConnector["redirectPolicy"] = proRedPolicy
		// 	}

		// 	if node["provider_subnet_ips"] != nil {
		// 		ips := node["provider_subnet_ips"].([]interface{})
		// 		prosubnets := make([]interface{}, 0, 1)
		// 		for _, ip := range ips {
		// 			subnet := make(map[string]interface{})
		// 			subnet["ip"] = ip.(string)

		// 			prosubnets = append(prosubnets, subnet)
		// 		}
		// 		proConnector["subnets"] = prosubnets
		// 	}

		// 	/*consumerConnector*/
		// 	conConnector := make(map[string]interface{})
		// 	conClusterInterface := make(map[string]interface{})
		// 	conClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["consumer_connector_cluster_interface"].(string))
		// 	conConnector["clusterInterface"] = conClusterInterface

		// 	if node["consumer_connector_redirect_policy"] != "" {
		// 		if node["consumer_connector_redirect_policy_tenant"] == "" {
		// 			return fmt.Errorf("consumer redirect policy tenant is required")
		// 		}
		// 		conRedPolicy := make(map[string]interface{})
		// 		conRedPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["consumer_connector_redirect_policy_tenant"].(string), node["consumer_connector_redirect_policy"].(string))
		// 		conConnector["redirectPolicy"] = conRedPolicy
		// 	}

		// 	if node["consumer_subnet_ips"] != nil {
		// 		ips := node["consumer_subnet_ips"].([]interface{})
		// 		consubnets := make([]interface{}, 0, 1)
		// 		for _, ip := range ips {
		// 			subnet := make(map[string]interface{})
		// 			subnet["ip"] = ip.(string)

		// 			consubnets = append(consubnets, subnet)
		// 		}
		// 		conConnector["subnets"] = consubnets
		// 	}

		// 	nodeMap := make(map[string]interface{})
		// 	nodeMap["serviceNodeRef"] = tempnodeRef
		// 	nodeMap["providerConnector"] = proConnector
		// 	nodeMap["consumerConnector"] = conConnector

		// 	siteNodes = append(siteNodes, nodeMap)
	}

	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("replace", tempPath, serviceGraphRef, templateNodes)

	// sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, TemplateName, contractName)
	// siteConGraph := models.NewSiteContractServiceGraph("replace", sitePath, serviceGraphRef, siteNodes)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Update Template Contract Service Graph")
	return resourceTemplateContractServiceGraphRead(d, m)
}

func resourceTemplateContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Template Contract Service Graph")
	msoClient := m.(*client.Client)
	foundTemp := false
	foundSite := false

	schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
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
		serviceGraphRef["templateName"] = templateName
	}

	// var graphSiteID string
	// if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
	// 	graphSiteID = graphSite.(string)
	// } else {
	// 	graphSiteID = siteID
	// }

	// contSite, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	// if err != nil {
	// 	return err
	// }

	// _, _, err = getSiteServiceGraph(contSite, serviceGraphRef["schemaId"].(string), serviceGraphRef["templateName"].(string), graphSiteID, serviceGraphRef["serviceGraphName"].(string))
	// if err == nil {
	// 	if _, ok := d.GetOk("service_graph_site_id"); ok {
	// 		d.Set("service_graph_site_id", graphSiteID)
	// 	} else {
	// 		d.Set("service_graph_site_id", "")
	// 	}
	// } else {
	// 	d.Set("service_graph_site_id", "")
	// }

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

						graphRelation := contractCont.S("serviceGraphRelationship")

						graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
						tokens := strings.Split(graphRef, "/")
						d.Set("service_graph_name", tokens[len(tokens)-1])
						if _, ok := d.GetOk("service_graph_schema_id"); ok {
							d.Set("service_graph_schema_id", tokens[len(tokens)-5])
						} else {
							d.Set("service_graph_schema_id", "")
						}
						if _, ok := d.GetOk("service_graph_template_name"); ok {
							d.Set("service_graph_template_name", tokens[len(tokens)-3])
						} else {
							d.Set("service_graph_template_name", "")
						}

						nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
						if err != nil {
							return err
						}
						for k := 0; k < nodeCount; k++ {
							relationMap := make(map[string]interface{})
							node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
							if err != nil {
								return err
							}
							nodeRef := models.StripQuotes(node.S("serviceNodeRef").String())
							tokensNode := strings.Split(nodeRef, "/")
							relationMap["node_name"] = tokensNode[len(tokensNode)-1]

							nodeInterface := d.Get("node_relationship")

							probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
							probdRefTokens := strings.Split(probdRef, "/")
							relationMap["provider_connector_bd_name"] = probdRefTokens[len(probdRefTokens)-1]

							if checkNodeAttr(nodeInterface, "provider_connector_bd_schema_id", k) {
								relationMap["provider_connector_bd_schema_id"] = probdRefTokens[len(probdRefTokens)-5]
							} else {
								relationMap["provider_connector_bd_schema_id"] = ""
							}
							if checkNodeAttr(nodeInterface, "provider_connector_bd_template_name", k) {
								relationMap["provider_connector_bd_template_name"] = probdRefTokens[len(probdRefTokens)-3]
							} else {
								relationMap["provider_connector_bd_template_name"] = ""
							}

							conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
							conbdRefTokens := strings.Split(conbdRef, "/")
							relationMap["consumer_connector_bd_name"] = conbdRefTokens[len(conbdRefTokens)-1]

							if checkNodeAttr(nodeInterface, "consumer_connector_bd_schema_id", k) {
								relationMap["consumer_connector_bd_schema_id"] = conbdRefTokens[len(conbdRefTokens)-5]
							} else {
								relationMap["consumer_connector_bd_schema_id"] = ""
							}

							if checkNodeAttr(nodeInterface, "consumer_connector_bd_template_name", k) {
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

	// siterelationList := make([]interface{}, 0, 1)
	// siteCount, err := cont.ArrayCount("sites")
	// if err != nil {
	// 	return fmt.Errorf("No sites found")
	// }
	// for i := 0; i < siteCount; i++ {
	// 	siteCont, err := cont.ArrayElement(i, "sites")
	// 	if err != nil {
	// 		return fmt.Errorf("Error fetching site")
	// 	}

	// 	site := models.StripQuotes(siteCont.S("siteId").String())
	// 	temp := models.StripQuotes(siteCont.S("templateName").String())
	// 	if siteID == site && temp == templateName {
	// 		contractCount, err := siteCont.ArrayCount("contracts")
	// 		if err != nil {
	// 			return fmt.Errorf("No contracts found in site")
	// 		}

	// 		for j := 0; j < contractCount; j++ {
	// 			contractCont, err := siteCont.ArrayElement(j, "contracts")
	// 			if err != nil {
	// 				return fmt.Errorf("Error fetching contract from site")
	// 			}

	// 			conRef := models.StripQuotes(contractCont.S("contractRef").String())
	// 			conTokens := strings.Split(conRef, "/")
	// 			conName := conTokens[len(conTokens)-1]
	// 			if conName == contractName {
	// 				if !contractCont.Exists("serviceGraphRelationship") {
	// 					return fmt.Errorf("No service graph found")
	// 				} else {
	// 					graphRelation := contractCont.S("serviceGraphRelationship")

	// 					nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
	// 					if err != nil {
	// 						return err
	// 					}
	// 					for k := 0; k < nodeCount; k++ {
	// 						relationMap := make(map[string]interface{})
	// 						node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
	// 						if err != nil {
	// 							return err
	// 						}

	// 						relationMap["provider_connector_cluster_interface"] = models.StripQuotes(node.S("providerConnector", "clusterInterface", "dn").String())

	// 						if node.Exists("providerConnector", "redirectPolicy", "dn") {
	// 							relationMap["provider_connector_redirect_policy"] = models.StripQuotes(node.S("providerConnector", "redirectPolicy", "dn").String())
	// 						}

	// 						if node.Exists("providerConnector", "subnets") {
	// 							subCounts, err := node.ArrayCount("providerConnector", "subnets")
	// 							if err != nil {
	// 								return err
	// 							}
	// 							subList := make([]interface{}, 0, 1)
	// 							for l := 0; l < subCounts; l++ {
	// 								subnet, err := node.ArrayElement(l, "providerConnector", "subnets", "ip")
	// 								if err != nil {
	// 									return err
	// 								}
	// 								subList = append(subList, models.StripQuotes(subnet.String()))
	// 							}
	// 							relationMap["provider_subnet_ips"] = subList
	// 						}

	// 						relationMap["consumer_connector_cluster_interface"] = models.StripQuotes(node.S("consumerConnector", "clusterInterface", "dn").String())

	// 						if node.Exists("consumerConnector", "redirectPolicy", "dn") {
	// 							relationMap["consumer_connector_redirect_policy"] = models.StripQuotes(node.S("consumerConnector", "redirectPolicy", "dn").String())
	// 						}

	// 						if node.Exists("consumerConnector", "subnets") {
	// 							subCounts, err := node.ArrayCount("consumerConnector", "subnets")
	// 							if err != nil {
	// 								return err
	// 							}
	// 							subList := make([]interface{}, 0, 1)
	// 							for l := 0; l < subCounts; l++ {
	// 								subnet, err := node.ArrayElement(l, "consumerConnector", "subnets", "ip")
	// 								if err != nil {
	// 									return err
	// 								}
	// 								subList = append(subList, models.StripQuotes(subnet.String()))
	// 							}
	// 							relationMap["consumer_subnet_ips"] = subList
	// 						}

	// 						siterelationList = append(siterelationList, relationMap)
	// 					}
	// 					foundSite = true
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if foundSite {
	// 		break
	// 	}
	// }

	foundSite = true

	if foundSite && foundTemp {
		length := len(temprelationList)
		nodeList := make([]interface{}, 0, 1)
		for i := 0; i < length; i++ {
			tempMap := temprelationList[i].(map[string]interface{})
			// siteMap := siterelationList[i].(map[string]interface{})

			allMap := make(map[string]interface{})
			allMap["provider_connector_bd_name"] = tempMap["provider_connector_bd_name"]
			allMap["provider_connector_bd_schema_id"] = tempMap["provider_connector_bd_schema_id"]
			allMap["provider_connector_bd_template_name"] = tempMap["provider_connector_bd_template_name"]
			allMap["consumer_connector_bd_name"] = tempMap["consumer_connector_bd_name"]
			allMap["consumer_connector_bd_schema_id"] = tempMap["consumer_connector_bd_schema_id"]
			allMap["consumer_connector_bd_template_name"] = tempMap["consumer_connector_bd_template_name"]

			// tp := strings.Split(siteMap["provider_connector_cluster_interface"].(string), "/")
			// token := strings.Split(tp[len(tp)-1], "-")
			// allMap["provider_connector_cluster_interface"] = token[1]

			// tp = strings.Split(siteMap["consumer_connector_cluster_interface"].(string), "/")
			// token = strings.Split(tp[len(tp)-1], "-")
			// allMap["consumer_connector_cluster_interface"] = token[1]

			// if siteMap["provider_connector_redirect_policy"] != nil {
			// 	tp := strings.Split(siteMap["provider_connector_redirect_policy"].(string), "/")
			// 	token1 := strings.Split(tp[1], "-")
			// 	allMap["provider_connector_redirect_policy_tenant"] = token1[1]

			// 	token2 := strings.Split(tp[len(tp)-1], "-")
			// 	allMap["provider_connector_redirect_policy"] = token2[1]
			// }
			// if siteMap["consumer_connector_redirect_policy"] != nil {
			// 	tp := strings.Split(siteMap["consumer_connector_redirect_policy"].(string), "/")
			// 	token1 := strings.Split(tp[1], "-")
			// 	allMap["consumer_connector_redirect_policy_tenant"] = token1[1]

			// 	token2 := strings.Split(tp[len(tp)-1], "-")
			// 	allMap["consumer_connector_redirect_policy"] = token2[1]
			// }

			// if siteMap["provider_subnet_ips"] != nil {
			// 	allMap["provider_subnet_ips"] = siteMap["provider_subnet_ips"]
			// }
			// if siteMap["consumer_subnet_ips"] != nil {
			// 	allMap["consumer_subnet_ips"] = siteMap["consumer_subnet_ips"]
			// }

			nodeList = append(nodeList, allMap)
		}
		d.Set("schema_id", schemaID)
		// d.Set("site_id", siteID)
		d.Set("template_name", templateName)
		d.Set("node_relationship", nodeList)
		d.Set("contract_name", contractName)

		if d.Get("service_graph_name") == serviceGraph {
			d.SetId(serviceGraph)
		} else {
			d.SetId("")
		}
	} else {
		d.SetId("")
	}

	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}

// func resourceTemplateContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
// 	log.Printf("[DEBUG] Begining Read Template Contract Service Graph")
// 	msoClient := m.(*client.Client)
// 	foundTemp := false
// 	foundSite := false

// 	schemaID := d.Get("schema_id").(string)
// 	siteID := d.Get("site_id").(string)
// 	templateName := d.Get("template_name").(string)
// 	contractName := d.Get("contract_name").(string)
// 	serviceGraph := d.Get("service_graph_name").(string)

// 	serviceGraphRef := make(map[string]interface{})
// 	if graphSchema, ok := d.GetOk("service_graph_schema_id"); ok {
// 		serviceGraphRef["schemaId"] = graphSchema.(string)
// 	} else {
// 		serviceGraphRef["schemaId"] = schemaID
// 	}

// 	if graphTemp, ok := d.GetOk("service_graph_template_name"); ok {
// 		serviceGraphRef["templateName"] = graphTemp.(string)
// 	} else {
// 		serviceGraphRef["templateName"] = templateName
// 	}
// 	serviceGraphRef["serviceGraphName"] = serviceGraph

// 	var graphSiteID string
// 	if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
// 		graphSiteID = graphSite.(string)
// 	} else {
// 		graphSiteID = siteID
// 	}

// 	contSite, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
// 	if err != nil {
// 		return err
// 	}

// 	_, _, err = getSiteServiceGraph(contSite, serviceGraphRef["schemaId"].(string), serviceGraphRef["templateName"].(string), graphSiteID, serviceGraphRef["serviceGraphName"].(string))
// 	if err == nil {
// 		if _, ok := d.GetOk("service_graph_site_id"); ok {
// 			d.Set("service_graph_site_id", graphSiteID)
// 		} else {
// 			d.Set("service_graph_site_id", "")
// 		}
// 	} else {
// 		d.Set("service_graph_site_id", "")
// 	}

// 	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
// 	if err != nil {
// 		return err
// 	}
// 	tempCount, err := cont.ArrayCount("templates")
// 	if err != nil {
// 		return fmt.Errorf("No templates found")
// 	}

// 	temprelationList := make([]interface{}, 0, 1)
// 	for i := 0; i < tempCount; i++ {
// 		tempCont, err := cont.ArrayElement(i, "templates")
// 		if err != nil {
// 			return fmt.Errorf("Error in fetch of template")
// 		}
// 		template := models.StripQuotes(tempCont.S("name").String())
// 		if templateName == template {
// 			contractCount, err := tempCont.ArrayCount("contracts")
// 			if err != nil {
// 				return fmt.Errorf("No contracts found")
// 			}

// 			for j := 0; j < contractCount; j++ {
// 				contractCont, err := tempCont.ArrayElement(j, "contracts")
// 				if err != nil {
// 					return fmt.Errorf("Error fetching contract")
// 				}
// 				conName := models.StripQuotes(contractCont.S("name").String())
// 				if conName == contractName {
// 					if !contractCont.Exists("serviceGraphRelationship") {
// 						return fmt.Errorf("No service graph found")
// 					} else {

// 						graphRelation := contractCont.S("serviceGraphRelationship")

// 						graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
// 						tokens := strings.Split(graphRef, "/")
// 						d.Set("service_graph_name", tokens[len(tokens)-1])
// 						if _, ok := d.GetOk("service_graph_schema_id"); ok {
// 							d.Set("service_graph_schema_id", tokens[len(tokens)-5])
// 						} else {
// 							d.Set("service_graph_schema_id", "")
// 						}
// 						if _, ok := d.GetOk("service_graph_template_name"); ok {
// 							d.Set("service_graph_template_name", tokens[len(tokens)-3])
// 						} else {
// 							d.Set("service_graph_template_name", "")
// 						}

// 						nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
// 						if err != nil {
// 							return err
// 						}
// 						for k := 0; k < nodeCount; k++ {
// 							relationMap := make(map[string]interface{})
// 							node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
// 							if err != nil {
// 								return err
// 							}
// 							nodeRef := models.StripQuotes(node.S("serviceNodeRef").String())
// 							tokensNode := strings.Split(nodeRef, "/")
// 							relationMap["node_name"] = tokensNode[len(tokensNode)-1]

// 							nodeInterface := d.Get("node_relationship")

// 							probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
// 							probdRefTokens := strings.Split(probdRef, "/")
// 							relationMap["provider_connector_bd_name"] = probdRefTokens[len(probdRefTokens)-1]

// 							if checkNodeAttr(nodeInterface, "provider_connector_bd_schema_id", k) {
// 								relationMap["provider_connector_bd_schema_id"] = probdRefTokens[len(probdRefTokens)-5]
// 							} else {
// 								relationMap["provider_connector_bd_schema_id"] = ""
// 							}
// 							if checkNodeAttr(nodeInterface, "provider_connector_bd_template_name", k) {
// 								relationMap["provider_connector_bd_template_name"] = probdRefTokens[len(probdRefTokens)-3]
// 							} else {
// 								relationMap["provider_connector_bd_template_name"] = ""
// 							}

// 							conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
// 							conbdRefTokens := strings.Split(conbdRef, "/")
// 							relationMap["consumer_connector_bd_name"] = conbdRefTokens[len(conbdRefTokens)-1]

// 							if checkNodeAttr(nodeInterface, "consumer_connector_bd_schema_id", k) {
// 								relationMap["consumer_connector_bd_schema_id"] = conbdRefTokens[len(conbdRefTokens)-5]
// 							} else {
// 								relationMap["consumer_connector_bd_schema_id"] = ""
// 							}

// 							if checkNodeAttr(nodeInterface, "consumer_connector_bd_template_name", k) {
// 								relationMap["consumer_connector_bd_template_name"] = conbdRefTokens[len(conbdRefTokens)-3]
// 							} else {
// 								relationMap["consumer_connector_bd_template_name"] = ""
// 							}

// 							temprelationList = append(temprelationList, relationMap)
// 						}
// 						foundTemp = true
// 						break
// 					}
// 				}
// 			}
// 		}
// 		if foundTemp {
// 			break
// 		}
// 	}

// 	siterelationList := make([]interface{}, 0, 1)
// 	siteCount, err := cont.ArrayCount("sites")
// 	if err != nil {
// 		return fmt.Errorf("No sites found")
// 	}
// 	for i := 0; i < siteCount; i++ {
// 		siteCont, err := cont.ArrayElement(i, "sites")
// 		if err != nil {
// 			return fmt.Errorf("Error fetching site")
// 		}

// 		site := models.StripQuotes(siteCont.S("siteId").String())
// 		temp := models.StripQuotes(siteCont.S("templateName").String())
// 		if siteID == site && temp == templateName {
// 			contractCount, err := siteCont.ArrayCount("contracts")
// 			if err != nil {
// 				return fmt.Errorf("No contracts found in site")
// 			}

// 			for j := 0; j < contractCount; j++ {
// 				contractCont, err := siteCont.ArrayElement(j, "contracts")
// 				if err != nil {
// 					return fmt.Errorf("Error fetching contract from site")
// 				}

// 				conRef := models.StripQuotes(contractCont.S("contractRef").String())
// 				conTokens := strings.Split(conRef, "/")
// 				conName := conTokens[len(conTokens)-1]
// 				if conName == contractName {
// 					if !contractCont.Exists("serviceGraphRelationship") {
// 						return fmt.Errorf("No service graph found")
// 					} else {
// 						graphRelation := contractCont.S("serviceGraphRelationship")

// 						nodeCount, err := graphRelation.ArrayCount("serviceNodesRelationship")
// 						if err != nil {
// 							return err
// 						}
// 						for k := 0; k < nodeCount; k++ {
// 							relationMap := make(map[string]interface{})
// 							node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
// 							if err != nil {
// 								return err
// 							}

// 							relationMap["provider_connector_cluster_interface"] = models.StripQuotes(node.S("providerConnector", "clusterInterface", "dn").String())

// 							if node.Exists("providerConnector", "redirectPolicy", "dn") {
// 								relationMap["provider_connector_redirect_policy"] = models.StripQuotes(node.S("providerConnector", "redirectPolicy", "dn").String())
// 							}

// 							if node.Exists("providerConnector", "subnets") {
// 								subCounts, err := node.ArrayCount("providerConnector", "subnets")
// 								if err != nil {
// 									return err
// 								}
// 								subList := make([]interface{}, 0, 1)
// 								for l := 0; l < subCounts; l++ {
// 									subnet, err := node.ArrayElement(l, "providerConnector", "subnets", "ip")
// 									if err != nil {
// 										return err
// 									}
// 									subList = append(subList, models.StripQuotes(subnet.String()))
// 								}
// 								relationMap["provider_subnet_ips"] = subList
// 							}

// 							relationMap["consumer_connector_cluster_interface"] = models.StripQuotes(node.S("consumerConnector", "clusterInterface", "dn").String())

// 							if node.Exists("consumerConnector", "redirectPolicy", "dn") {
// 								relationMap["consumer_connector_redirect_policy"] = models.StripQuotes(node.S("consumerConnector", "redirectPolicy", "dn").String())
// 							}

// 							if node.Exists("consumerConnector", "subnets") {
// 								subCounts, err := node.ArrayCount("consumerConnector", "subnets")
// 								if err != nil {
// 									return err
// 								}
// 								subList := make([]interface{}, 0, 1)
// 								for l := 0; l < subCounts; l++ {
// 									subnet, err := node.ArrayElement(l, "consumerConnector", "subnets", "ip")
// 									if err != nil {
// 										return err
// 									}
// 									subList = append(subList, models.StripQuotes(subnet.String()))
// 								}
// 								relationMap["consumer_subnet_ips"] = subList
// 							}

// 							siterelationList = append(siterelationList, relationMap)
// 						}
// 						foundSite = true
// 					}
// 				}
// 			}
// 		}
// 		if foundSite {
// 			break
// 		}
// 	}

// 	if foundSite && foundTemp {
// 		length := len(temprelationList)
// 		nodeList := make([]interface{}, 0, 1)
// 		for i := 0; i < length; i++ {
// 			tempMap := temprelationList[i].(map[string]interface{})
// 			siteMap := siterelationList[i].(map[string]interface{})

// 			allMap := make(map[string]interface{})
// 			allMap["provider_connector_bd_name"] = tempMap["provider_connector_bd_name"]
// 			allMap["provider_connector_bd_schema_id"] = tempMap["provider_connector_bd_schema_id"]
// 			allMap["provider_connector_bd_template_name"] = tempMap["provider_connector_bd_template_name"]
// 			allMap["consumer_connector_bd_name"] = tempMap["consumer_connector_bd_name"]
// 			allMap["consumer_connector_bd_schema_id"] = tempMap["consumer_connector_bd_schema_id"]
// 			allMap["consumer_connector_bd_template_name"] = tempMap["consumer_connector_bd_template_name"]

// 			tp := strings.Split(siteMap["provider_connector_cluster_interface"].(string), "/")
// 			token := strings.Split(tp[len(tp)-1], "-")
// 			allMap["provider_connector_cluster_interface"] = token[1]

// 			tp = strings.Split(siteMap["consumer_connector_cluster_interface"].(string), "/")
// 			token = strings.Split(tp[len(tp)-1], "-")
// 			allMap["consumer_connector_cluster_interface"] = token[1]

// 			if siteMap["provider_connector_redirect_policy"] != nil {
// 				tp := strings.Split(siteMap["provider_connector_redirect_policy"].(string), "/")
// 				token1 := strings.Split(tp[1], "-")
// 				allMap["provider_connector_redirect_policy_tenant"] = token1[1]

// 				token2 := strings.Split(tp[len(tp)-1], "-")
// 				allMap["provider_connector_redirect_policy"] = token2[1]
// 			}
// 			if siteMap["consumer_connector_redirect_policy"] != nil {
// 				tp := strings.Split(siteMap["consumer_connector_redirect_policy"].(string), "/")
// 				token1 := strings.Split(tp[1], "-")
// 				allMap["consumer_connector_redirect_policy_tenant"] = token1[1]

// 				token2 := strings.Split(tp[len(tp)-1], "-")
// 				allMap["consumer_connector_redirect_policy"] = token2[1]
// 			}

// 			if siteMap["provider_subnet_ips"] != nil {
// 				allMap["provider_subnet_ips"] = siteMap["provider_subnet_ips"]
// 			}
// 			if siteMap["consumer_subnet_ips"] != nil {
// 				allMap["consumer_subnet_ips"] = siteMap["consumer_subnet_ips"]
// 			}

// 			nodeList = append(nodeList, allMap)
// 		}
// 		d.Set("schema_id", schemaID)
// 		d.Set("site_id", siteID)
// 		d.Set("template_name", templateName)
// 		d.Set("node_relationship", nodeList)
// 		d.Set("contract_name", contractName)

// 		if d.Get("service_graph_name") == serviceGraph {
// 			d.SetId(serviceGraph)
// 		} else {
// 			d.SetId("")
// 		}
// 	} else {
// 		d.SetId("")
// 	}

// 	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
// 	return nil
// }

func resourceTemplateContractServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Delete Template Contract Service Graph")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	// siteID := d.Get("site_id").(string)
	TemplateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", TemplateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph("remove", tempPath, nil, nil)

	// sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, TemplateName, contractName)
	// siteConGraph := models.NewSiteContractServiceGraph("remove", sitePath, nil, nil)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Template Contract Service Graph")
	return nil
}

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

func getSiteServiceNode(graphCont *container.Container, schemaId, templateName, graphName, nodeName string) (*container.Container, int, error) {

	nodesCount, err := graphCont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load count site service node")
	}
	for i := 0; i < nodesCount; i++ {
		nodeCont, err := graphCont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to site service node element")
		}

		nodeRef := models.StripQuotes(nodeCont.S("serviceNodeRef").String())

		nodeSplit := strings.Split(nodeRef, "/")
		if len(nodeSplit) == 9 {
			if nodeSplit[2] == schemaId && nodeSplit[4] == templateName && nodeSplit[6] == graphName && nodeSplit[8] == nodeName {
				return nodeCont, i, nil

			}
		} else {
			return nil, -1, fmt.Errorf("Spilt on nodeRef failed")
		}
	}
	return nil, -1, fmt.Errorf("Unable to find site service node")
}

func getSiteServiceGraph(cont *container.Container, schemaId, templateName, siteId, graphName string) (*container.Container, int, error) {
	log.Printf("[DEBUG] : ######## getSiteServiceGraph - cont: %v", cont)

	sitesCount, err := cont.ArrayCount("sites")

	if err != nil {
		return nil, -1, fmt.Errorf("Unable to find sites")
	}
	log.Printf("[DEBUG] : ######## sitesCount: %v", sitesCount)

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load site element")
		}

		siteTemplate := models.StripQuotes(siteCont.S("templateName").String())
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		log.Printf("[DEBUG] : siteTemplate : %s, templateName: %s, siteId: %s, apiSiteId: %s", siteTemplate, templateName, siteId, apiSiteId)

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
					return nil, -1, fmt.Errorf("Inavlid site service graph")
				}

				if schemaId == graphEle[2] && templateName == graphEle[4] && graphName == graphEle[6] {
					return sgCont, j, nil
				}

			}
		}
	}

	return nil, -1, fmt.Errorf("Unable to find site service graph")
}

func checkNodeAttr(object interface{}, attrName string, index int) bool {
	objList := object.([]interface{})

	instance := objList[index].(map[string]interface{})

	if instance[attrName] != "" {
		return true
	}
	return false
}
