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

func resourceMSOSchemaTemplateContractServiceGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateContractServiceGraphCreate,
		Update: resourceMSOSchemaTemplateContractServiceGraphUpdate,
		Read:   resourceMSOSchemaTemplateContractServiceGraphRead,
		Delete: resourceMSOSchemaTemplateContractServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateContractServiceGraphImport,
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
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_template_name": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"node_relationship": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The order of the node_relationship object should match the node types in the Service Graph",
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
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"provider_connector_bd_template_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
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
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"consumer_connector_bd_template_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		},
	}
}

func resourceMSOSchemaTemplateContractServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	serviceGraphTokens := strings.Split(d.Id(), "/")
	d.Set("schema_id", serviceGraphTokens[0])
	d.Set("template_name", serviceGraphTokens[2])
	d.Set("contract_name", serviceGraphTokens[4])

	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", serviceGraphTokens[0]))
	if err != nil {
		return nil, err
	}
	err = setSchemaTemplateContractServiceGraphAttrs(cont, d)
	if err != nil {
		return nil, err
	}
	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateContractServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Service Graph: Beginning Creation")
	err := postSchemaTemplateContractServiceGraphConfig("add", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaTemplateContractServiceGraphRead(d, m)
}

func resourceMSOSchemaTemplateContractServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Service Graph: Beginning Update")
	err := postSchemaTemplateContractServiceGraphConfig("replace", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaTemplateContractServiceGraphRead(d, m)
}

func resourceMSOSchemaTemplateContractServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Template Contract Service Graph")
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSchemaTemplateContractServiceGraphAttrs(cont, d)

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Completed Read Template Contract Service Graph")
	return nil
}

func resourceMSOSchemaTemplateContractServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Delete Template Contract Service Graph")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	tempPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", templateName, contractName)
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), models.GetRemovePatchPayload(tempPath))

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Template Contract Service Graph")
	return nil
}

// Returns the List of Service Graph Node map object with Consumer and Provider BD values
func getSchemaTemplateContractServiceGraphNodes(cont *container.Container, schemaID, templateName string, apiServiceGraphNodeList, nodeRelationshipList []interface{}, serviceGraphRef map[string]interface{}) ([]interface{}, error) {
	templateNodesList := make([]interface{}, 0)

	for i := 0; i < len(nodeRelationshipList); i++ {
		node := nodeRelationshipList[i].(map[string]interface{})

		serviceGraphNodeRef := make(map[string]interface{})

		serviceGraphNodeRef["schemaId"] = serviceGraphRef["schemaId"]
		serviceGraphNodeRef["serviceGraphName"] = serviceGraphRef["serviceGraphName"]
		serviceGraphNodeRef["templateName"] = serviceGraphRef["templateName"]
		serviceGraphNodeRef["serviceNodeName"] = apiServiceGraphNodeList[i].(string)

		providerConnectorMap := make(map[string]interface{})
		providerConnectorMap["connectorType"] = "general"
		providerConnectorBDRef := make(map[string]interface{})
		if node["provider_connector_bd_schema_id"] != "" {
			providerConnectorBDRef["schemaId"] = node["provider_connector_bd_schema_id"].(string)
		} else {
			providerConnectorBDRef["schemaId"] = schemaID
		}

		if node["provider_connector_bd_template_name"] != "" {
			providerConnectorBDRef["templateName"] = node["provider_connector_bd_template_name"].(string)
		} else {
			providerConnectorBDRef["templateName"] = templateName
		}
		providerConnectorBDRef["bdName"] = node["provider_connector_bd_name"].(string)
		providerConnectorMap["bdRef"] = providerConnectorBDRef

		consumerConnectorMap := make(map[string]interface{})
		consumerConnectorMap["connectorType"] = "general"
		consumerConnectorBDRef := make(map[string]interface{})
		if node["consumer_connector_bd_schema_id"] != "" {
			consumerConnectorBDRef["schemaId"] = node["consumer_connector_bd_schema_id"].(string)
		} else {
			consumerConnectorBDRef["schemaId"] = schemaID
		}

		if node["consumer_connector_bd_template_name"] != "" {
			consumerConnectorBDRef["templateName"] = node["consumer_connector_bd_template_name"].(string)
		} else {
			consumerConnectorBDRef["templateName"] = templateName
		}
		consumerConnectorBDRef["bdName"] = node["consumer_connector_bd_name"].(string)
		consumerConnectorMap["bdRef"] = consumerConnectorBDRef

		templateNodeMap := make(map[string]interface{})
		templateNodeMap["serviceNodeRef"] = serviceGraphNodeRef
		templateNodeMap["providerConnector"] = providerConnectorMap
		templateNodeMap["consumerConnector"] = consumerConnectorMap

		templateNodesList = append(templateNodesList, templateNodeMap)
	}
	return templateNodesList, nil
}

// setSchemaTemplateContractServiceGraphAttrs sets the resource attributes of the service graph
// for a template contract in the schema.
//
// Parameters:
// - d: *schema.ResourceData: The resource data containing the schema attributes.
// - m: interface{}: The client interface.
//
// Returns:
// - error: The error encountered, if any.
func setSchemaTemplateContractServiceGraphAttrs(cont *container.Container, d *schema.ResourceData) error {
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	templatesCount, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No templates found")
	}

	serviceGraphRelationshipList := make([]interface{}, 0, 1)
	for i := 0; i < templatesCount; i++ {
		templatesCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return fmt.Errorf("Error in fetch of template")
		}
		apiTemplateName := models.StripQuotes(templatesCont.S("name").String())
		if templateName == apiTemplateName {
			contractCount, err := templatesCont.ArrayCount("contracts")
			if err != nil {
				return fmt.Errorf("No contracts found")
			}

			for j := 0; j < contractCount; j++ {
				contractCont, err := templatesCont.ArrayElement(j, "contracts")
				if err != nil {
					return fmt.Errorf("Error fetching contract")
				}
				apiContractName := models.StripQuotes(contractCont.S("name").String())
				if apiContractName == contractName {
					if !contractCont.Exists("serviceGraphRelationship") {
						return fmt.Errorf("No service graph found")
					} else {
						// Template Contract Service Graph configurations
						serviceGraphRelationship := contractCont.S("serviceGraphRelationship")
						serviceGraphRef := models.StripQuotes(serviceGraphRelationship.S("serviceGraphRef").String())
						serviceGraphReftokens := strings.Split(serviceGraphRef, "/")

						d.Set("service_graph_name", serviceGraphReftokens[6])

						if _, ok := d.GetOk("service_graph_schema_id"); !ok {
							d.Set("service_graph_schema_id", serviceGraphReftokens[2])
						} else {
							d.Set("service_graph_schema_id", d.Get("service_graph_schema_id"))
						}

						if _, ok := d.GetOk("service_graph_template_name"); !ok {
							d.Set("service_graph_template_name", serviceGraphReftokens[4])
						} else {
							d.Set("service_graph_template_name", d.Get("service_graph_template_name"))
						}

						// Template Contract Service Graph Node configurations
						serviceNodesRelationshipCount, err := serviceGraphRelationship.ArrayCount("serviceNodesRelationship")
						if err != nil {
							return err
						}

						for k := 0; k < serviceNodesRelationshipCount; k++ {
							relationMap := make(map[string]interface{})
							node, err := serviceGraphRelationship.ArrayElement(k, "serviceNodesRelationship")
							if err != nil {
								return err
							}

							providerConnectorBDRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
							providerConnectorBDRefTokens := strings.Split(providerConnectorBDRef, "/")

							relationMap["provider_connector_bd_name"] = providerConnectorBDRefTokens[len(providerConnectorBDRefTokens)-1]
							relationMap["provider_connector_bd_schema_id"] = providerConnectorBDRefTokens[len(providerConnectorBDRefTokens)-5]
							relationMap["provider_connector_bd_template_name"] = providerConnectorBDRefTokens[len(providerConnectorBDRefTokens)-3]

							consumerConnectorBDRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
							consumerConnectorBDRefTokens := strings.Split(consumerConnectorBDRef, "/")

							relationMap["consumer_connector_bd_name"] = consumerConnectorBDRefTokens[len(consumerConnectorBDRefTokens)-1]
							relationMap["consumer_connector_bd_schema_id"] = consumerConnectorBDRefTokens[len(consumerConnectorBDRefTokens)-5]
							relationMap["consumer_connector_bd_template_name"] = consumerConnectorBDRefTokens[len(consumerConnectorBDRefTokens)-3]

							serviceGraphRelationshipList = append(serviceGraphRelationshipList, relationMap)
						}

						nodeList := make([]interface{}, 0, 1)
						for i := 0; i < len(serviceGraphRelationshipList); i++ {
							tempMap := serviceGraphRelationshipList[i].(map[string]interface{})

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
						d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s", d.Get("schema_id").(string), templateName, contractName))
						return nil
					}
				}
			}
		}
	}
	d.SetId("")
	return nil
}

// postSchemaTemplateContractServiceGraphConfig create/update a service graph configuration for a template contract.
//
// Parameters:
//   - ops: The ops to perform create(add)/update(replace) operations.
//   - d: The schema resource data.
//   - m: The client interface.
//
// Returns:
//   - An error if there was a problem creating the service graph configuration.
func postSchemaTemplateContractServiceGraphConfig(ops string, d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceGraphName := d.Get("service_graph_name").(string)

	serviceGraphRef := make(map[string]interface{})
	serviceGraphRef["serviceGraphName"] = serviceGraphName

	if serviceGraphSchemaID, ok := d.GetOk("service_graph_schema_id"); ok {
		serviceGraphRef["schemaId"] = serviceGraphSchemaID.(string)
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

	templateServiceGraphCont, _, err := getSchemaTemplateServiceGraphFromContainer(cont, serviceGraphRef["templateName"].(string), serviceGraphName)

	if err != nil {
		return err
	}
	apiServiceGraphNodeList := extractServiceGraphNodesFromContainer(templateServiceGraphCont)

	nodeRelationshipList := d.Get("node_relationship").([]interface{})

	if len(nodeRelationshipList) != len(apiServiceGraphNodeList) {
		return fmt.Errorf("Length mismatch between total nodes and total node relation")
	}

	contractServiceGraphNodes, err := getSchemaTemplateContractServiceGraphNodes(cont, schemaID, templateName, apiServiceGraphNodeList, nodeRelationshipList, serviceGraphRef)

	if err != nil {
		return err
	}

	contractServiceGraphPath := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", templateName, contractName)
	tempConGraph := models.NewTemplateContractServiceGraph(ops, contractServiceGraphPath, serviceGraphRef, contractServiceGraphNodes)
	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), tempConGraph)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s", schemaID, templateName, contractName))
	return nil
}
