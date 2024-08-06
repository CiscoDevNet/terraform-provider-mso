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

func resourceMSOTemplateContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateContractCreate,
		Read:   resourceMSOTemplateContractRead,
		Update: resourceMSOTemplateContractUpdate,
		Delete: resourceMSOTemplateContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateContractImport,
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
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"application-profile",
					"tenant",
					"context",
					"global",
				}, false),
			},
			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"af11",
					"af12",
					"af13",
					"af21",
					"af22",
					"af23",
					"af31",
					"af32",
					"af33",
					"af41",
					"af42",
					"af43",
					"cs0",
					"cs1",
					"cs2",
					"cs3",
					"cs4",
					"cs5",
					"cs6",
					"cs7",
					"expeditedForwarding",
					"voiceAdmit",
					"unspecified",
				}, false),
			},
			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"level6",
					"level5",
					"level4",
					"level3",
					"level2",
					"level1",
					"unspecified",
				}, false),
			},
			"filter_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bothWay",
					"oneWay",
				}, false),
			},
			"filter_relationship": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_schema_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"filter_template_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"filter_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"filter_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"bothWay",
								"provider_to_consumer",
								"consumer_to_provider",
							}, false),
						},
						"directives": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"none",
									"no_stats",
									"log",
								}, false),
							},
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"deny",
								"permit",
							}, false),
						},
						"priority": {
							Type:     schema.TypeString,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"level3",
								"level2",
								"level1",
								"default",
							}, false),
						},
					},
				},
			},
			"filter_relationships": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				ConflictsWith: []string{"filter_relationship"},
				Deprecated:    "use filter_relationship instead",
			},
			"directives": {
				Type:       schema.TypeList,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Optional:   true,
				Computed:   true,
				Deprecated: "use directives in filter_relationship instead",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			stateFilterType, configFilterType := diff.GetChange("filter_type")
			if configFilterType != stateFilterType && stateFilterType != "" {
				return fmt.Errorf("The filter_type cannot be changed. Change detected from '%s' to '%s'.", stateFilterType, configFilterType)
			}

			return nil
		},
	}
}

func createMSOTemplateContractPath(templateName, contractName string) string {
	return fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
}

// TODO remove this deprecated function when filter_relationships is removed
func getDeprecatedFilterRelationshipsFromConfig(schemaId, templateName string, filterRelationshipsConfig map[string]interface{}, directives []interface{}) []interface{} {

	var filterRelationships []interface{}

	relationshipMap := make(map[string]interface{})

	if filterRelationshipsConfig["filter_schema_id"] != nil {
		schemaId = filterRelationshipsConfig["filter_schema_id"].(string)
	}
	if filterRelationshipsConfig["filter_template_name"] != nil {
		templateName = filterRelationshipsConfig["filter_template_name"].(string)
	}
	relationshipMap["filterRef"] = map[string]interface{}{
		"schemaId":     schemaId,
		"templateName": templateName,
		"filterName":   filterRelationshipsConfig["filter_name"].(string),
	}

	relationshipMap["directives"] = directives

	return append(filterRelationships, relationshipMap)

}

func getFilterRelationshipsFromConfig(schemaId, templateName string, filterRelationshipsConfig, directives []interface{}) ([]interface{}, []interface{}, []interface{}) {

	var filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider []interface{}

	for _, relationshipConfig := range filterRelationshipsConfig {
		relationshipConfigMap := relationshipConfig.(map[string]interface{})
		relationshipMap := make(map[string]interface{})

		relationshipSchemaId := schemaId
		if relationshipConfigMap["filter_schema_id"] != "" {
			relationshipSchemaId = relationshipConfigMap["filter_schema_id"].(string)
		}
		relationshipTemplateName := templateName
		if relationshipConfigMap["filter_template_name"] != "" {
			relationshipTemplateName = relationshipConfigMap["filter_template_name"].(string)
		}
		relationshipMap["filterRef"] = map[string]interface{}{
			"schemaId":     relationshipSchemaId,
			"templateName": relationshipTemplateName,
			"filterName":   relationshipConfigMap["filter_name"].(string),
		}

		if len(relationshipConfigMap["directives"].(*schema.Set).List()) > 0 {
			relationshipMap["directives"] = relationshipConfigMap["directives"].(*schema.Set).List()
		} else {
			relationshipMap["directives"] = directives
		}

		if relationshipConfigMap["action"].(string) != "" {
			relationshipMap["action"] = relationshipConfigMap["action"].(string)
		}

		if relationshipConfigMap["priority"].(string) != "" {
			relationshipMap["priorityOverride"] = relationshipConfigMap["priority"].(string)
		}

		if relationshipConfigMap["filter_type"].(string) == "bothWay" {
			filterRelationships = append(filterRelationships, relationshipMap)
		} else if relationshipConfigMap["filter_type"].(string) == "provider_to_consumer" {
			filterRelationshipsProviderToConsumer = append(filterRelationshipsProviderToConsumer, relationshipMap)
		} else if relationshipConfigMap["filter_type"].(string) == "consumer_to_provider" {
			filterRelationshipsConsumerToProvider = append(filterRelationshipsConsumerToProvider, relationshipMap)
		}
	}

	return filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider

}

func setFilterRelationshipList(relationships []interface{}, filterList []map[string]interface{}, filterType string) []map[string]interface{} {
	for _, relationship := range relationships {
		relationshipMap := relationship.(map[string]interface{})
		filterMap := map[string]interface{}{"filter_type": filterType}
		if val, ok := relationshipMap["filterRef"]; ok {
			re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/filters/(.*)")
			match := re.FindStringSubmatch(val.(string))
			filterMap["filter_schema_id"] = match[1]
			filterMap["filter_template_name"] = match[2]
			filterMap["filter_name"] = match[3]
		}
		if val, ok := relationshipMap["directives"]; val != nil && ok {
			filterMap["directives"] = val
		}
		if val, ok := relationshipMap["action"]; val != nil && ok {
			filterMap["action"] = val
		}
		if val, ok := relationshipMap["priorityOverride"]; val != nil && ok {
			filterMap["priority"] = val
		}
		filterList = append(filterList, filterMap)
	}
	return filterList
}

func setContractFromSchema(d *schema.ResourceData, schemaCont *container.Container, schemaId, templateName, contractName string) error {
	templates := schemaCont.Search("templates").Data()
	if templates == nil || len(templates.([]interface{})) == 0 {
		return fmt.Errorf("no templates found")
	}

	for _, template := range templates.([]interface{}) {
		templateDetails := template.(map[string]interface{})
		if templateDetails["name"].(string) == templateName {
			for _, contract := range templateDetails["contracts"].([]interface{}) {
				contractDetails := contract.(map[string]interface{})
				if contractDetails["name"].(string) == contractName {
					d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s", schemaId, templateName, contractName))
					d.Set("schema_id", schemaId)
					d.Set("template_name", templateName)
					d.Set("contract_name", contractName)

					if val, ok := contractDetails["displayName"]; val != nil && ok {
						d.Set("display_name", val.(string))
					}
					if val, ok := contractDetails["scope"]; val != nil && ok {
						d.Set("scope", val.(string))
					}
					if val, ok := contractDetails["filterType"]; val != nil && ok {
						d.Set("filter_type", val.(string))
					}
					if val, ok := contractDetails["prio"]; val != nil && ok {
						d.Set("priority", val.(string))
					}
					if val, ok := contractDetails["targetDscp"]; val != nil && ok {
						d.Set("target_dscp", val.(string))
					}
					if val, ok := contractDetails["description"]; val != nil && ok {
						d.Set("description", val.(string))
					}

					filterList := []map[string]interface{}{}
					if val, ok := contractDetails["filterRelationships"]; val != nil && ok {
						filterList = setFilterRelationshipList(val.([]interface{}), filterList, "bothWay")

						// TODO Remove below block of code once the filterRelationships is deprecated
						// Start of block
						// Reason for adding this logic is backworth compatibility with previous release

						if val, ok := d.GetOk("filter_relationships"); val != nil && ok {
							filterRelationshipsMap := val.(map[string]interface{})
							filterSchemaId := schemaId
							if filterRelationshipsMap["filter_schema_id"] != nil {
								filterSchemaId = filterRelationshipsMap["filter_schema_id"].(string)
							}
							filterTemplateName := templateName
							if filterRelationshipsMap["filter_template_name"] != nil {
								filterTemplateName = filterRelationshipsMap["filter_template_name"].(string)
							}
							filterName := filterRelationshipsMap["filter_name"].(string)

							for _, fiterMap := range filterList {
								if fiterMap["filter_schema_id"].(string) == filterSchemaId &&
									fiterMap["filter_template_name"].(string) == filterTemplateName &&
									fiterMap["filter_name"].(string) == filterName {
									d.Set("filter_relationships", filterRelationshipsMap)
									d.Set("directives", fiterMap["directives"])
								}
							}
						} else {
							/* When filterRelationships is not provided provide the last entry in the filter list
							Below was implemented for the datasource where it loops through and overwrites the value in map

							filterMap := make(map[string]interface{})
							for i := 0; i < count; i++ {
								filterCont, err := contractCont.ArrayElement(i, "filterRelationships")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships list")
								}

								d.Set("directives", filterCont.S("directives").Data().([]interface{}))
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")

								filterMap["filter_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["filter_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["filter_name"] = fmt.Sprintf("%s", split[6])
							}
							d.Set("filter_relationships", filterMap)

							*/
							filterRelationshipsMap := make(map[string]interface{})
							if len(filterList) != 0 {
								filterMap := filterList[len(filterList)-1]
								filterRelationshipsMap["filter_schema_id"] = filterMap["filter_schema_id"]
								filterRelationshipsMap["filter_template_name"] = filterMap["filter_template_name"]
								filterRelationshipsMap["filter_name"] = filterMap["filter_name"]
								d.Set("filter_relationships", filterRelationshipsMap)
								d.Set("directives", filterMap["directives"])
							} else {
								d.Set("filter_relationships", filterRelationshipsMap)
								d.Set("directives", []string{})
							}
						}

						// End of block

					}
					if val, ok := contractDetails["filterRelationshipsProviderToConsumer"]; val != nil && ok {
						filterList = setFilterRelationshipList(val.([]interface{}), filterList, "provider_to_consumer")
					}
					if val, ok := contractDetails["filterRelationshipsConsumerToProvider"]; val != nil && ok {
						filterList = setFilterRelationshipList(val.([]interface{}), filterList, "consumer_to_provider")
					}

					d.Set("filter_relationship", filterList)

					return nil
				}
			}
		}
	}

	return fmt.Errorf("Unable to find the Contract: %s", contractName)
}

func resourceMSOTemplateContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	splitImport := strings.Split(d.Id(), "/")
	schemaId := splitImport[0]
	templateName := splitImport[2]
	contractName := splitImport[4]
	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	err = setContractFromSchema(d, schemaCont, schemaId, templateName, contractName)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTemplateContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)
	scope := d.Get("scope").(string)
	priority := d.Get("priority").(string)
	targetDscp := d.Get("target_dscp").(string)
	filterType := d.Get("filter_type").(string)
	filterRelationship := d.Get("filter_relationship").([]interface{})
	// TODO remove when filter_relationships and directives are deprecated on next mayor version
	deprecatedFilterRelationship := d.Get("filter_relationships").(map[string]interface{})
	directives := d.Get("directives").([]interface{})
	var filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider []interface{}
	if len(deprecatedFilterRelationship) > 0 {
		filterRelationships = getDeprecatedFilterRelationshipsFromConfig(schemaId, templateName, deprecatedFilterRelationship, directives)
	} else {
		filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider = getFilterRelationshipsFromConfig(schemaId, templateName, filterRelationship, directives)
	}
	// TODO uncomment line below when filter_relationships and directives are deprecated on next major version
	// filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider := getFilterRelationshipsFromConfig(schemaId, templateName, filterRelationship)
	var description string
	if Description, ok := d.GetOk("description"); ok {
		description = Description.(string)
	}
	path := createMSOTemplateContractPath(templateName, "-")
	contractStruct := models.NewTemplateContract("add", path, contractName, displayName, scope, filterType, targetDscp, priority, description, filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), contractStruct)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), schemaCont, d)
	}
	err = setContractFromSchema(d, schemaCont, schemaId, templateName, contractName)
	if err != nil {
		d.SetId("")
		log.Printf("[DEBUG] Resetting Id due to setContractFromSchema returning '%s'", err)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTemplateContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)
	scope := d.Get("scope").(string)
	priority := d.Get("priority").(string)
	targetDscp := d.Get("target_dscp").(string)
	filterType := d.Get("filter_type").(string)
	filterRelationship := d.Get("filter_relationship").([]interface{})

	// TODO remove when filter_relationships and directives are deprecated on next mayor version
	directives := d.Get("directives").([]interface{})
	var filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider []interface{}
	if d.HasChange("filter_relationships") {
		deprecatedFilterRelationship := d.Get("filter_relationships").(map[string]interface{})
		filterRelationships = getDeprecatedFilterRelationshipsFromConfig(schemaId, templateName, deprecatedFilterRelationship, directives)
	} else {
		filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider = getFilterRelationshipsFromConfig(schemaId, templateName, filterRelationship, directives)
	}
	// TODO uncomment line below when filter_relationships and directives are deprecated on next mayor version
	// filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider := getFilterRelationshipsFromConfig(schemaId, templateName, filterRelationship)

	var description string
	if Description, ok := d.GetOk("description"); ok {
		description = Description.(string)
	}

	path := createMSOTemplateContractPath(templateName, contractName)
	contractStruct := models.NewTemplateContract("replace", path, contractName, displayName, scope, filterType, targetDscp, priority, description, filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), contractStruct)
	if err != nil {
		return err
	}
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Delete", d.Id())
	msoClient := m.(*client.Client)
	path := createMSOTemplateContractPath(d.Get("template_name").(string), d.Get("contract_name").(string))
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Get("schema_id").(string)), models.GetRemovePatchPayload(path))
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	log.Printf("[DEBUG] Delete finished successfully")
	return nil
}
