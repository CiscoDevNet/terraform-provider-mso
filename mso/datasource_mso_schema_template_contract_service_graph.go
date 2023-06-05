package mso

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceTemplateContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTemplateContractServiceGraphRead,

		Schema: map[string]*schema.Schema{
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
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_graph_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_relationship": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_connector_bd_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_bd_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_connector_bd_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector_bd_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceTemplateContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {

	// msoClient := m.(*client.Client)
	// foundTemp := false
	// foundSite := false

	// schemaID := d.Get("schema_id").(string)
	// // siteID := d.Get("site_id").(string)
	// templateName := d.Get("template_name").(string)
	// contractName := d.Get("contract_name").(string)
	// serviceGraph := d.Get("service_graph_name").(string)

	// serviceGraphTokens := strings.Split(d.Id(), "/")
	// log.Printf("[DEBUG] %s: ###### tokens Beginning Import, 2nd element: %s", serviceGraphTokens, serviceGraphTokens[1])

	// d.Set("schema_id", serviceGraphTokens[2])
	// d.Set("template_name", serviceGraphTokens[4])
	// d.Set("contract_name", serviceGraphTokens[6])
	// d.Set("service_graph_name", "TSG")
	setTemplateContractServiceGraphAttrs(d, m, true)

	// cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	// if err != nil {
	// 	return err
	// }
	// tempCount, err := cont.ArrayCount("templates")
	// if err != nil {
	// 	return fmt.Errorf("No templates found")
	// }

	// temprelationList := make([]interface{}, 0, 1)
	// for i := 0; i < tempCount; i++ {
	// 	tempCont, err := cont.ArrayElement(i, "templates")
	// 	if err != nil {
	// 		return fmt.Errorf("Error in fetch of template")
	// 	}
	// 	template := models.StripQuotes(tempCont.S("name").String())
	// 	if templateName == template {
	// 		contractCount, err := tempCont.ArrayCount("contracts")
	// 		if err != nil {
	// 			return fmt.Errorf("No contracts found")
	// 		}

	// 		for j := 0; j < contractCount; j++ {
	// 			contractCont, err := tempCont.ArrayElement(j, "contracts")
	// 			if err != nil {
	// 				return fmt.Errorf("Error fetching contract")
	// 			}
	// 			conName := models.StripQuotes(contractCont.S("name").String())
	// 			if conName == contractName {
	// 				if !contractCont.Exists("serviceGraphRelationship") {
	// 					return fmt.Errorf("No service graph found")
	// 				} else {

	// 					graphRelation := contractCont.S("serviceGraphRelationship")

	// 					graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
	// 					tokens := strings.Split(graphRef, "/")
	// 					d.Set("service_graph_name", tokens[len(tokens)-1])

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
	// 						nodeRef := models.StripQuotes(node.S("serviceNodeRef").String())
	// 						tokensNode := strings.Split(nodeRef, "/")
	// 						relationMap["node_name"] = tokensNode[len(tokensNode)-1]

	// 						probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
	// 						probdRefTokens := strings.Split(probdRef, "/")
	// 						relationMap["provider_connector_bd_name"] = probdRefTokens[len(probdRefTokens)-1]

	// 						conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
	// 						conbdRefTokens := strings.Split(conbdRef, "/")
	// 						relationMap["consumer_connector_bd_name"] = conbdRefTokens[len(conbdRefTokens)-1]

	// 						temprelationList = append(temprelationList, relationMap)
	// 					}
	// 					foundTemp = true
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if foundTemp {
	// 		break
	// 	}
	// }

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

	// if foundSite && foundTemp {
	// 	length := len(temprelationList)
	// 	nodeList := make([]interface{}, 0, 1)
	// 	for i := 0; i < length; i++ {
	// 		tempMap := temprelationList[i].(map[string]interface{})
	// 		siteMap := siterelationList[i].(map[string]interface{})

	// 		allMap := make(map[string]interface{})
	// 		allMap["provider_connector_bd_name"] = tempMap["provider_connector_bd_name"]
	// 		allMap["consumer_connector_bd_name"] = tempMap["consumer_connector_bd_name"]

	// 		tp := strings.Split(siteMap["provider_connector_cluster_interface"].(string), "/")
	// 		token := strings.Split(tp[len(tp)-1], "-")
	// 		allMap["provider_connector_cluster_interface"] = token[1]

	// 		tp = strings.Split(siteMap["consumer_connector_cluster_interface"].(string), "/")
	// 		token = strings.Split(tp[len(tp)-1], "-")
	// 		allMap["consumer_connector_cluster_interface"] = token[1]

	// 		if siteMap["provider_connector_redirect_policy"] != nil {
	// 			tp := strings.Split(siteMap["provider_connector_redirect_policy"].(string), "/")
	// 			token1 := strings.Split(tp[1], "-")
	// 			allMap["provider_connector_redirect_policy_tenant"] = token1[1]

	// 			token2 := strings.Split(tp[len(tp)-1], "-")
	// 			allMap["provider_connector_redirect_policy"] = token2[1]
	// 		}
	// 		if siteMap["consumer_connector_redirect_policy"] != nil {
	// 			tp := strings.Split(siteMap["consumer_connector_redirect_policy"].(string), "/")
	// 			token1 := strings.Split(tp[1], "-")
	// 			allMap["consumer_connector_redirect_policy_tenant"] = token1[1]

	// 			token2 := strings.Split(tp[len(tp)-1], "-")
	// 			allMap["consumer_connector_redirect_policy"] = token2[1]
	// 		}

	// 		if siteMap["provider_subnet_ips"] != nil {
	// 			allMap["provider_subnet_ips"] = siteMap["provider_subnet_ips"]
	// 		}
	// 		if siteMap["consumer_subnet_ips"] != nil {
	// 			allMap["consumer_subnet_ips"] = siteMap["consumer_subnet_ips"]
	// 		}

	// 		nodeList = append(nodeList, allMap)
	// 	}
	// 	d.Set("schema_id", schemaID)
	// 	d.Set("site_id", siteID)
	// 	d.Set("template_name", templateName)
	// 	d.Set("node_relationship", nodeList)
	// 	d.Set("contract_name", contractName)

	// 	if d.Get("service_graph_name") == serviceGraph {
	// 		d.SetId(serviceGraph)
	// 	} else {
	// 		d.SetId("")
	// 		return fmt.Errorf("No service graph found for given name")
	// 	}
	// } else {
	// 	d.SetId("")
	// 	return fmt.Errorf("No service graph found for given name")
	// }

	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}
