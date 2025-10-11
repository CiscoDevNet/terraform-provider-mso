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

func resourceMSOSchemaTemplateContractServiceChaining() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaServiceChainingCreate,
		Read:   resourceMSOSchemaServiceChainingRead,
		Update: resourceMSOSchemaServiceChainingUpdate,
		Delete: resourceMSOSchemaServiceChainingDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaServiceChainingImport,
		},

		SchemaVersion: version,

		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"allow-all", "filters-from-contract"}, false),
			},
			"service_nodes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of service nodes in the service chaining graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"device_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"firewall", "loadBalancer", "other"}, false),
						},
						"device_ref": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_connector": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 100),
									},
									"is_redirect": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"provider_connector": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 100),
									},
									"is_redirect": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func contractServiceChainingPath(templateName, contractName string) string {
	return fmt.Sprintf("/templates/%s/contracts/%s/serviceChaining", templateName, contractName)
}

func parseServiceChainingId(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")

	if parts[1] != "templates" || parts[3] != "contracts" || parts[5] != "serviceChaining" {
		expectedFormat := "<schema_id>/templates/<template_name>/contracts/<contract_name>/serviceChaining"
		return "", "", "", fmt.Errorf("invalid ID structure: got '%s', expected format: %s", id, expectedFormat)
	}

	schemaId := parts[0]
	templateName := parts[2]
	contractName := parts[4]

	return schemaId, templateName, contractName, nil
}

func buildServiceChainingPayload(name, nodeFilter string, serviceNodes []interface{}) map[string]interface{} {
	serviceNodesPayload := make([]interface{}, 0, len(serviceNodes))
	for i, serviceNode := range serviceNodes {
		node := serviceNode.(map[string]interface{})

		nodePayload := map[string]interface{}{
			"name":       node["name"].(string),
			"deviceType": node["device_type"].(string),
			"deviceRef":  node["device_ref"].(string),
			"index":      i + 1,
		}

		if consConnectorList, ok := node["consumer_connector"].([]interface{}); ok && len(consConnectorList) > 0 {
			consConnector := consConnectorList[0].(map[string]interface{})
			nodePayload["consumerConnector"] = map[string]interface{}{
				"interfaceName": consConnector["interface_name"].(string),
				"isRedirect":    consConnector["is_redirect"].(bool),
			}
		}
		if provConnectorList, ok := node["provider_connector"].([]interface{}); ok && len(provConnectorList) > 0 {
			provConnector := provConnectorList[0].(map[string]interface{})
			nodePayload["providerConnector"] = map[string]interface{}{
				"interfaceName": provConnector["interface_name"].(string),
				"isRedirect":    provConnector["is_redirect"].(bool),
			}
		}

		serviceNodesPayload = append(serviceNodesPayload, nodePayload)
	}

	return map[string]interface{}{
		"name":         name,
		"nodeFilter":   nodeFilter,
		"serviceNodes": serviceNodesPayload,
	}
}

func resourceMSOSchemaServiceChainingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	resourceMSOSchemaServiceChainingRead(d, m)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaServiceChainingCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Service Chaining: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	nodeFilter := d.Get("node_filter").(string)
	serviceNodes := d.Get("service_nodes").([]interface{})

	payload := buildServiceChainingPayload(contractName, nodeFilter, serviceNodes)

	path := contractServiceChainingPath(templateName, contractName)
	payloadModel := models.GetPatchPayload("add", path, payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s/serviceChaining", schemaId, templateName, contractName))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaServiceChainingRead(d, m)
}

func resourceMSOSchemaServiceChainingUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Service Chaining: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	nodeFilter := d.Get("node_filter").(string)
	serviceNodes := d.Get("service_nodes").([]interface{})

	payload := buildServiceChainingPayload(contractName, nodeFilter, serviceNodes)
	path := contractServiceChainingPath(templateName, contractName)

	payloadCont := container.New()
	payloadCont.Array()
	if d.HasChange("service_nodes") {
		err := addPatchPayloadToContainer(payloadCont, "replace", path, payload)
		if err != nil {
			return err
		}
	}

	err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s/serviceChaining", schemaId, templateName, contractName))
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaServiceChainingRead(d, m)
}

func resourceMSOSchemaServiceChainingRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Read Service Chaining")
	msoClient := m.(*client.Client)

	schemaId, templateName, contractName, err := parseServiceChainingId(d.Id())
	if err != nil {
		return err
	}

	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), schemaCont, d)
	}

	if err := setServiceChainingFromSchema(d, schemaCont, schemaId, templateName, contractName); err != nil {
		return err
	}

	log.Printf("[DEBUG] Completed Read Service Chaining")
	return nil
}

func resourceMSOSchemaServiceChainingDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Delete Service Chaining")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	path := contractServiceChainingPath(templateName, contractName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), models.GetRemovePatchPayload(path))
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Service Chaining")
	return nil
}

func setServiceChainingFromSchema(d *schema.ResourceData, schemaCont *container.Container, schemaId, templateName, contractName string) error {
	log.Printf("[DEBUG] %s: Beginning set Service Chaining from schema", d.Id())

	templates := schemaCont.Search("templates").Data()
	if templates == nil || len(templates.([]interface{})) == 0 {
		return fmt.Errorf("no templates found")
	}

	var contractDetails map[string]interface{}
	for _, template := range templates.([]interface{}) {
		templateMap := template.(map[string]interface{})
		if templateMap["name"].(string) == templateName {
			d.Set("template_name", templateName)
			contracts, ok := templateMap["contracts"].([]interface{})
			if !ok || len(contracts) == 0 {
				return fmt.Errorf("no contracts found in template %s", templateName)
			}
			for _, contract := range contracts {
				contractMap := contract.(map[string]interface{})
				if contractMap["name"].(string) == contractName {
					d.Set("contract_name", contractName)
					contractDetails = contractMap
					break
				}
			}
			break
		}
	}

	if contractDetails == nil {
		d.SetId("")
		return fmt.Errorf("contract %s not found in template %s", contractName, templateName)
	}

	serviceChainingIface, ok := contractDetails["serviceChaining"]
	if !ok || serviceChainingIface == nil {
		d.SetId("")
		return fmt.Errorf("serviceChaining not found in contract %s", contractName)
	}

	serviceChain, ok := serviceChainingIface.(map[string]interface{})
	if !ok {
		d.SetId("")
		return fmt.Errorf("invalid serviceChaining structure in contract %s", contractName)
	}

	if nameVal, ok := serviceChain["name"].(string); ok {
		d.Set("name", nameVal)
	}

	if nodeFilter, ok := serviceChain["nodeFilter"].(string); ok {
		d.Set("node_filter", nodeFilter)
	}

	if serviceNodes, ok := serviceChain["serviceNodes"].([]interface{}); ok {
		out := make([]interface{}, 0, len(serviceNodes))
		for _, serviceNode := range serviceNodes {
			nodeMap := serviceNode.(map[string]interface{})
			item := map[string]interface{}{
				"name":        nodeMap["name"],
				"device_type": nodeMap["deviceType"],
				"device_ref":  nodeMap["deviceRef"],
				"index":       nodeMap["index"],
				"uuid":        nodeMap["uuid"],
			}

			if consConnector, ok := nodeMap["consumerConnector"].(map[string]interface{}); ok {
				item["consumer_connector"] = []interface{}{
					map[string]interface{}{
						"interface_name": consConnector["interfaceName"],
						"is_redirect":    consConnector["isRedirect"],
					},
				}
			}
			if provConnector, ok := nodeMap["providerConnector"].(map[string]interface{}); ok {
				item["provider_connector"] = []interface{}{
					map[string]interface{}{
						"interface_name": provConnector["interfaceName"],
						"is_redirect":    provConnector["isRedirect"],
					},
				}
			}

			out = append(out, item)
		}
		d.Set("service_nodes", out)
	}
	d.Set("schema_id", schemaId)
	d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s/serviceChaining", schemaId, templateName, contractName))

	log.Printf("[DEBUG] %s: Finished set Service Chaining from schema", d.Id())
	return nil
}
