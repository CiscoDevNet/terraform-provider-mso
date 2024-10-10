package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteContractServiceGraphCreate,
		Update: resourceMSOSchemaSiteContractServiceGraphUpdate,
		Read:   resourceMSOSchemaSiteContractServiceGraphRead,
		Delete: resourceMSOSchemaSiteContractServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteContractServiceGraphImport,
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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"node_relationship": &schema.Schema{ // Only for non-cloud sites
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_connector_cluster_interface": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_redirect_policy_tenant": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_redirect_policy": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_cluster_interface": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_redirect_policy_tenant": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_redirect_policy": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_subnet_ips": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceMSOSchemaSiteContractServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	serviceGraphTokens := strings.Split(d.Id(), "/")
	d.Set("schema_id", serviceGraphTokens[0])
	d.Set("site_id", serviceGraphTokens[2])
	d.Set("template_name", serviceGraphTokens[4])
	d.Set("contract_name", serviceGraphTokens[6])
	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphTokens[0]))
	if err != nil {
		return nil, err
	}
	err = setSiteContractServiceGraphAttrs(cont, d)
	if err != nil {
		return nil, err
	}
	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteContractServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Template Contract Service Graph: Beginning Creation")
	err := postSiteContractServiceGraphConfig("add", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaSiteContractServiceGraphRead(d, m)

}

func resourceMSOSchemaSiteContractServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Template Contract Service Graph: Beginning Update")
	err := postSiteContractServiceGraphConfig("replace", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteContractServiceGraphRead(d, m)
}

func resourceMSOSchemaSiteContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Site Template Contract Service Graph")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSiteContractServiceGraphAttrs(cont, d)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Completed Read Site Template Contract Service Graph")
	return nil
}

func resourceMSOSchemaSiteContractServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Delete Site Template Contract Service Graph")
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, templateName, contractName)
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), models.GetRemovePatchPayload(sitePath))

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Site Template Contract Service Graph")
	return nil
}

// Sets the resource attribute values
func setSiteContractServiceGraphAttrs(cont *container.Container, d *schema.ResourceData) error {
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	siteID := d.Get("site_id").(string)
	includeNodesRelationship := false

	siteRelationList := make([]interface{}, 0, 1)

	siteCount, err := cont.ArrayCount("sites")

	if err != nil {
		return fmt.Errorf("No sites found")
	}
	for i := 0; i < siteCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return fmt.Errorf("Error fetching site")
		}

		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())
		apiTemplateName := models.StripQuotes(siteCont.S("templateName").String())

		if siteID == apiSiteId && apiTemplateName == templateName {
			if siteID == "" {
				siteID = apiSiteId
			}
			contractCount, err := siteCont.ArrayCount("contracts")
			if err != nil {
				return fmt.Errorf("No contracts found in site")
			}

			for j := 0; j < contractCount; j++ {
				contractCont, err := siteCont.ArrayElement(j, "contracts")
				if err != nil {
					return fmt.Errorf("Error fetching contract from site")
				}

				contractRef := models.StripQuotes(contractCont.S("contractRef").String())
				contractTokens := strings.Split(contractRef, "/")
				apiContractName := contractTokens[len(contractTokens)-1]
				if apiContractName == contractName {
					if !contractCont.Exists("serviceGraphRelationship") {
						return fmt.Errorf("No service graph found")
					} else {
						siteServiceGraphRef := models.StripQuotes(contractCont.S("serviceGraphRelationship", "serviceGraphRef").String())
						consumerConnectorPresent := models.StripQuotes(contractCont.S("serviceGraphRelationship", "serviceNodesRelationship", "consumerConnector").String())
						providerConnectorPresent := models.StripQuotes(contractCont.S("serviceGraphRelationship", "serviceNodesRelationship", "providerConnector").String())

						// Only for non-cloud sites
						if consumerConnectorPresent != "{}" && providerConnectorPresent != "{}" {
							includeNodesRelationship = true
							serviceGraphRelationship := contractCont.S("serviceGraphRelationship")
							nodeCount, err := serviceGraphRelationship.ArrayCount("serviceNodesRelationship")
							if err != nil {
								return err
							}
							for k := 0; k < nodeCount; k++ {
								relationMap := make(map[string]interface{})
								node, err := serviceGraphRelationship.ArrayElement(k, "serviceNodesRelationship")
								if err != nil {
									return err
								}

								relationMap["provider_connector_cluster_interface"] = models.StripQuotes(node.S("providerConnector", "clusterInterface", "dn").String())

								if node.Exists("providerConnector", "redirectPolicy", "dn") {
									relationMap["provider_connector_redirect_policy"] = models.StripQuotes(node.S("providerConnector", "redirectPolicy", "dn").String())
								}

								relationMap["consumer_connector_cluster_interface"] = models.StripQuotes(node.S("consumerConnector", "clusterInterface", "dn").String())

								if node.Exists("consumerConnector", "redirectPolicy", "dn") {
									relationMap["consumer_connector_redirect_policy"] = models.StripQuotes(node.S("consumerConnector", "redirectPolicy", "dn").String())
								}

								if node.Exists("consumerConnector", "subnets") {
									subnetsCount, err := node.ArrayCount("consumerConnector", "subnets")
									if err != nil {
										return err
									}
									subnetList := make([]interface{}, 0, 1)
									for l := 0; l < subnetsCount; l++ {
										subnet, err := node.ArrayElement(l, "consumerConnector", "subnets", "ip")
										if err != nil {
											return err
										}
										subnetList = append(subnetList, models.StripQuotes(subnet.String()))
									}
									relationMap["consumer_subnet_ips"] = subnetList
								}
								siteRelationList = append(siteRelationList, relationMap)
							}
						}
						siteServiceGraphRefValues := strings.Split(siteServiceGraphRef, "/")
						d.Set("service_graph_name", siteServiceGraphRefValues[6])
						d.Set("service_graph_schema_id", siteServiceGraphRefValues[2])
						d.Set("service_graph_template_name", siteServiceGraphRefValues[4])

						nodeList := make([]interface{}, 0, 1)
						// Only for non-cloud sites
						if includeNodesRelationship {
							length := len(siteRelationList)
							for i := 0; i < length; i++ {
								siteMap := siteRelationList[i].(map[string]interface{})

								allMap := make(map[string]interface{})

								re := regexp.MustCompile("uni/tn-(.*)/lDevVip-(.*)/lIf-(.*)")
								provSplit := re.FindStringSubmatch(siteMap["provider_connector_cluster_interface"].(string))
								allMap["provider_connector_cluster_interface"] = provSplit[3]
								consSplit := re.FindStringSubmatch(siteMap["consumer_connector_cluster_interface"].(string))
								allMap["consumer_connector_cluster_interface"] = consSplit[3]

								re = regexp.MustCompile("uni/tn-(.*)/svcCont/svcRedirectPol-(.*)")
								if siteMap["provider_connector_redirect_policy"] != nil {
									split := re.FindStringSubmatch(siteMap["provider_connector_redirect_policy"].(string))
									allMap["provider_connector_redirect_policy_tenant"] = split[1]
									allMap["provider_connector_redirect_policy"] = split[2]
								}
								if siteMap["consumer_connector_redirect_policy"] != nil {
									split := re.FindStringSubmatch(siteMap["consumer_connector_redirect_policy"].(string))
									allMap["consumer_connector_redirect_policy_tenant"] = split[1]
									allMap["consumer_connector_redirect_policy"] = split[2]
								}
								if siteMap["consumer_subnet_ips"] != nil {
									allMap["consumer_subnet_ips"] = siteMap["consumer_subnet_ips"]
								}
								nodeList = append(nodeList, allMap)
							}
						}

						d.Set("site_id", siteID)
						d.Set("node_relationship", nodeList)

						var graphSiteID string
						if graphSite, ok := d.GetOk("service_graph_site_id"); ok {
							graphSiteID = graphSite.(string)
						} else {
							graphSiteID = siteID
						}
						d.Set("service_graph_site_id", graphSiteID)
						d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s", schemaID, siteID, templateName, contractName))
						return nil
					}
				}
			}
		}
	}

	d.SetId("")
	return nil
}

// Returns the list of site service nodes relationship objects - serviceNodesRelationship
func getSiteServiceNodesRelationshipObject(cont *container.Container, schemaID, serviceGraphSiteID, templateName string, includeNodesRelationship bool, apiNodeRelationshipList, nodeRelationshipList []interface{}, serviceGraphRef map[string]interface{}) ([]interface{}, error) {
	siteNodes := make([]interface{}, 0, 1)

	for i := 0; i < len(apiNodeRelationshipList); i++ {
		nodeMap := make(map[string]interface{})
		serviceGraphSchemaID := serviceGraphRef["schemaId"].(string)
		serviceGraphTemplateName := serviceGraphRef["templateName"].(string)
		serviceGraphName := serviceGraphRef["serviceGraphName"].(string)
		serviceGraphNodeName := apiNodeRelationshipList[i].(string)

		siteGraphCont, _, err := getSiteTemplateServiceGraph(cont, serviceGraphSchemaID, serviceGraphTemplateName, serviceGraphSiteID, serviceGraphName)
		if err != nil {
			return nil, err
		}

		siteNodeCont, _, err := getSiteTemplateServiceGraphNode(siteGraphCont, serviceGraphSchemaID, serviceGraphTemplateName, serviceGraphName, serviceGraphNodeName)
		if err != nil {
			return nil, err
		}

		// Only for non-cloud sites
		if includeNodesRelationship {
			node := nodeRelationshipList[i].(map[string]interface{})
			dn := models.StripQuotes(siteNodeCont.S("device", "dn").String())
			providerConnector := make(map[string]interface{})
			providerClusterInterface := make(map[string]interface{})
			providerClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["provider_connector_cluster_interface"].(string))
			providerConnector["clusterInterface"] = providerClusterInterface
			if node["provider_connector_redirect_policy"] != "" {
				if node["provider_connector_redirect_policy_tenant"] == "" {
					return nil, fmt.Errorf("provider redirect policy tenant is required")
				}
				providerRedirectPolicy := make(map[string]interface{})
				providerRedirectPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["provider_connector_redirect_policy_tenant"].(string), node["provider_connector_redirect_policy"].(string))
				providerConnector["redirectPolicy"] = providerRedirectPolicy
			}

			consumerConnector := make(map[string]interface{})
			consumerClusterInterface := make(map[string]interface{})
			consumerClusterInterface["dn"] = fmt.Sprintf("%s/lIf-%s", dn, node["consumer_connector_cluster_interface"].(string))
			consumerConnector["clusterInterface"] = consumerClusterInterface
			if node["consumer_connector_redirect_policy"] != "" {
				if node["consumer_connector_redirect_policy_tenant"] == "" {
					return nil, fmt.Errorf("consumer redirect policy tenant is required")
				}
				consumerRedirectPolicy := make(map[string]interface{})
				consumerRedirectPolicy["dn"] = fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", node["consumer_connector_redirect_policy_tenant"].(string), node["consumer_connector_redirect_policy"].(string))
				consumerConnector["redirectPolicy"] = consumerRedirectPolicy
			}

			if node["consumer_subnet_ips"] != nil {
				ips := node["consumer_subnet_ips"].([]interface{})
				consumerSubnets := make([]interface{}, 0, 1)
				for _, ip := range ips {
					subnet := make(map[string]interface{})
					subnet["ip"] = ip.(string)
					consumerSubnets = append(consumerSubnets, subnet)
				}
				consumerConnector["subnets"] = consumerSubnets
			}
			nodeMap["providerConnector"] = providerConnector
			nodeMap["consumerConnector"] = consumerConnector
		}

		nodeMap["serviceNodeRef"] = map[string]interface{}{
			"schemaId":         serviceGraphSchemaID,
			"serviceGraphName": serviceGraphName,
			"templateName":     serviceGraphTemplateName,
			"serviceNodeName":  serviceGraphNodeName,
		}

		siteNodes = append(siteNodes, nodeMap)
	}
	return siteNodes, nil
}

// Returns the site service graph service node object and its index position
func getSiteTemplateServiceGraphNode(serviceGraphCont *container.Container, schemaID, templateName, serviceGraphName, serviceNodeName string) (*container.Container, int, error) {
	serviceNodesCount, err := serviceGraphCont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load count site service node")
	}
	for i := 0; i < serviceNodesCount; i++ {
		serviceNodeCont, err := serviceGraphCont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to site service node element")
		}

		serviceNodeRef := models.StripQuotes(serviceNodeCont.S("serviceNodeRef").String())
		nodeSplit := strings.Split(serviceNodeRef, "/")
		if len(nodeSplit) == 9 {
			if nodeSplit[2] == schemaID && nodeSplit[4] == templateName && nodeSplit[6] == serviceGraphName && nodeSplit[8] == serviceNodeName {
				return serviceNodeCont, i, nil
			}
		} else {
			return nil, -1, fmt.Errorf("Spilt on nodeRef failed")
		}
	}
	return nil, -1, fmt.Errorf("Unable to find site service node")
}

// Returns the Site Template serviceGraphs container object and its index position
func getSiteTemplateServiceGraph(cont *container.Container, schemaID, templateName, siteID, serviceGraphName string) (*container.Container, int, error) {
	sitesCount, err := cont.ArrayCount("sites")

	if err != nil {
		return nil, -1, fmt.Errorf("Unable to find sites")
	}

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load site element")
		}

		apiSiteTemplateName := models.StripQuotes(siteCont.S("templateName").String())
		apiSiteID := models.StripQuotes(siteCont.S("siteId").String())

		if apiSiteTemplateName == templateName && apiSiteID == siteID {
			serviceGraphsCount, err := siteCont.ArrayCount("serviceGraphs")
			if err != nil {
				return nil, -1, fmt.Errorf("Unable to load site service graphs")
			}

			for j := 0; j < serviceGraphsCount; j++ {
				serviceGraphCont, err := siteCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to load site service graph element")
				}

				serviceGraphRef := models.StripQuotes(serviceGraphCont.S("serviceGraphRef").String())
				serviceGraphRefToken := strings.Split(serviceGraphRef, "/")

				if len(serviceGraphRefToken) != 7 {
					return nil, -1, fmt.Errorf("Invalid site service graph")
				}

				if schemaID == serviceGraphRefToken[2] && templateName == serviceGraphRefToken[4] && serviceGraphName == serviceGraphRefToken[6] {
					return serviceGraphCont, j, nil
				}

			}
		}
	}
	return nil, -1, fmt.Errorf("Unable to find site service graph")
}

// postSiteContractServiceGraphConfig create/update a service graph configuration for a site template contract.
//
// Parameters:
//   - ops: The ops to perform create(add)/update(replace) operations.
//   - d: The schema resource data.
//   - m: The client interface.
//
// Returns:
//   - An error if there was a problem creating the service graph configuration.
func postSiteContractServiceGraphConfig(ops string, d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	siteID := d.Get("site_id").(string)

	var serviceGraphSiteID string
	if tempServiceGraphSiteID, ok := d.GetOk("service_graph_site_id"); ok {
		serviceGraphSiteID = tempServiceGraphSiteID.(string)
	} else {
		serviceGraphSiteID = siteID
	}

	serviceGraphRef := make(map[string]interface{})
	serviceGraphName := d.Get("service_graph_name").(string)
	serviceGraphRef["serviceGraphName"] = serviceGraphName

	if serviceGraphSchemaId, ok := d.GetOk("service_graph_schema_id"); ok {
		serviceGraphRef["schemaId"] = serviceGraphSchemaId.(string)
	} else {
		serviceGraphRef["schemaId"] = schemaID
	}

	if serviceGraphTemplateName, ok := d.GetOk("service_graph_template_name"); ok {
		serviceGraphRef["templateName"] = serviceGraphTemplateName.(string)
	} else {
		serviceGraphRef["templateName"] = templateName
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphRef["schemaId"].(string)))
	if err != nil {
		return err
	}

	schemaTemplateServiceGraphCont, _, err := getSchemaTemplateServiceGraphFromContainer(cont, serviceGraphRef["templateName"].(string), serviceGraphName)
	if err != nil {
		return err
	}
	apiNodeRelationshipList := extractServiceGraphNodesFromContainer(schemaTemplateServiceGraphCont)
	if err != nil {
		return err
	}

	nodeRelationshipList := d.Get("node_relationship").([]interface{})
	includeNodesRelationship := false
	if len(nodeRelationshipList) == len(apiNodeRelationshipList) {
		includeNodesRelationship = true
	}

	siteNodes, err := getSiteServiceNodesRelationshipObject(cont, schemaID, serviceGraphSiteID, templateName, includeNodesRelationship, apiNodeRelationshipList, nodeRelationshipList, serviceGraphRef)
	if err != nil {
		return err
	}

	sitePath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship", siteID, templateName, contractName)
	siteContractServiceGraphObject := models.NewSiteContractServiceGraph(ops, sitePath, serviceGraphRef, siteNodes)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), siteContractServiceGraphObject)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s", schemaID, siteID, templateName, contractName))
	return nil
}
