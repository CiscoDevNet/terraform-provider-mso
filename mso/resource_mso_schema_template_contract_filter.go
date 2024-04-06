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

func resourceMSOTemplateContractFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateContractFilterCreate,
		Read:   resourceMSOTemplateContractFilterRead,
		Update: resourceMSOTemplateContractFilterUpdate,
		Delete: resourceMSOTemplateContractFilterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateContractFilterImport,
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
			"filter_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bothWay",
					"provider_to_consumer",
					"consumer_to_provider",
				}, false),
			},
			"filter_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"filter_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"filter_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deny",
					"permit",
				}, false),
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"level3",
					"level2",
					"level1",
					"default",
				}, false),
			},
		}),
	}
}

func getFilterRef(filterSchemaId, filterTemplateName, filterName string) map[string]interface{} {
	return map[string]interface{}{"schemaId": filterSchemaId, "templateName": filterTemplateName, "filterName": filterName}
}

func getFilterRelationshipTypeMap() map[string]string {
	return map[string]string{
		"bothWay":              "filterRelationships",
		"provider_to_consumer": "filterRelationshipsProviderToConsumer",
		"consumer_to_provider": "filterRelationshipsConsumerToProvider",
	}
}

func createMSOTemplateContractFilterPath(templateName, contractName, filterRelationshipType, name string) string {
	return fmt.Sprintf("/templates/%s/contracts/%s/%s/%s", templateName, contractName, filterRelationshipType, name)
}

func setContractFilterFromSchema(d *schema.ResourceData, schemaCont *container.Container, schemaId, templateName, contractName, filterType, filterSchemaId, filterTemplateName, filterName string) error {
	log.Printf("[DEBUG] %s: Beginning set contract filter", d.Id())
	templates := schemaCont.Search("templates").Data()
	if templates == nil || len(templates.([]interface{})) == 0 {
		return fmt.Errorf("no templates found")
	}

	filterRelationshipType := getFilterRelationshipTypeMap()[filterType]

	for _, template := range templates.([]interface{}) {
		templateDetails := template.(map[string]interface{})
		if templateDetails["name"].(string) == templateName {
			for _, contract := range templateDetails["contracts"].([]interface{}) {
				contractDetails := contract.(map[string]interface{})
				if contractDetails["name"].(string) == contractName {
					if filterRelationships, ok := contractDetails[filterRelationshipType]; filterRelationships != nil && ok {
						for _, filterRelationship := range filterRelationships.([]interface{}) {
							filterRelationshipMap := filterRelationship.(map[string]interface{})
							if val, ok := filterRelationshipMap["filterRef"]; ok {
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/filters/(.*)")
								match := re.FindStringSubmatch(val.(string))
								if match[1] == filterSchemaId && match[2] == filterTemplateName && match[3] == filterName {
									d.SetId(fmt.Sprintf("%s/templates/%s/contracts/%s/%s/%s/%s/%s", schemaId, templateName, contractName, filterType, filterSchemaId, filterTemplateName, filterName))
									d.Set("schema_id", schemaId)
									d.Set("template_name", templateName)
									d.Set("contract_name", contractName)
									d.Set("directives", filterRelationshipMap["directives"])
									d.Set("action", filterRelationshipMap["action"])
									d.Set("priority", filterRelationshipMap["priority"])
									d.Set("filter_type", filterType)
									d.Set("filter_schema_id", filterSchemaId)
									d.Set("filter_template_name", filterTemplateName)
									d.Set("filter_name", filterName)
									log.Printf("[DEBUG] %s: Finished set contract filter", d.Id())
									return nil
								}
							}
						}
					}
				}
			}
		}
	}
	d.SetId("")
	return fmt.Errorf("unable to find contract: %s", contractName)
}

func resourceMSOTemplateContractFilterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	splitImport := strings.Split(d.Id(), "/")

	schemaId := splitImport[0]
	templateName := splitImport[2]
	contractName := splitImport[4]
	filterType := splitImport[5]
	filterSchemaId := splitImport[6]
	filterTemplateName := splitImport[7]
	filterName := splitImport[8]

	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	err = setContractFilterFromSchema(d, schemaCont, schemaId, templateName, contractName, filterType, filterSchemaId, filterTemplateName, filterName)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTemplateContractFilterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Filter: Beginning Creation")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	filterType := d.Get("filter_type").(string)
	filterName := d.Get("filter_name").(string)

	filterSchemaId := schemaId
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterSchemaId = tempVar.(string)
	}
	filterTemplateName := templateName
	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplateName = tempVar.(string)
	}
	filterRefMap := getFilterRef(filterSchemaId, filterTemplateName, filterName)

	var directives []interface{}
	var action, priority string
	if tempVar, ok := d.GetOk("directives"); ok {
		directives = tempVar.(*schema.Set).List()
	}
	if tempVar, ok := d.GetOk("action"); ok {
		action = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("priority"); ok {
		priority = tempVar.(string)
	}

	path := createMSOTemplateContractFilterPath(templateName, contractName, getFilterRelationshipTypeMap()[filterType], "-")
	filterStruct := models.NewTemplateContractFilterRelationShip("add", path, action, priority, "", filterRefMap, directives)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceMSOTemplateContractFilterRead(d, m)
}

func resourceMSOTemplateContractFilterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	filterType := d.Get("filter_type").(string)
	filterName := d.Get("filter_name").(string)

	filterSchemaId := schemaId
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterSchemaId = tempVar.(string)
	}
	filterTemplateName := templateName
	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplateName = tempVar.(string)
	}

	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), schemaCont, d)
	}
	setContractFilterFromSchema(d, schemaCont, schemaId, templateName, contractName, filterType, filterSchemaId, filterTemplateName, filterName)
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTemplateContractFilterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	filterName := d.Get("filter_name").(string)

	filterSchemaId := schemaId
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterSchemaId = tempVar.(string)
	}
	filterTemplateName := templateName
	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplateName = tempVar.(string)
	}
	filterRefMap := getFilterRef(filterSchemaId, filterTemplateName, filterName)

	var directives []interface{}
	var action, priority string
	if tempVar, ok := d.GetOk("directives"); ok {
		directives = tempVar.(*schema.Set).List()
	}
	if tempVar, ok := d.GetOk("action"); ok {
		action = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("priority"); ok {
		priority = tempVar.(string)
	}

	path := createMSOTemplateContractFilterPath(templateName, d.Get("contract_name").(string), getFilterRelationshipTypeMap()[d.Get("filter_type").(string)], filterName)
	filterStruct := models.NewTemplateContractFilterRelationShip("replace", path, action, priority, "", filterRefMap, directives)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOTemplateContractFilterRead(d, m)
}

func resourceMSOTemplateContractFilterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Delete", d.Id())
	msoClient := m.(*client.Client)
	templateName := d.Get("template_name").(string)
	filterName := d.Get("filter_name").(string)
	path := createMSOTemplateContractFilterPath(templateName, d.Get("contract_name").(string), getFilterRelationshipTypeMap()[d.Get("filter_type").(string)], filterName)
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Get("schema_id").(string)), models.GetRemovePatchPayload(path))
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")

	log.Printf("[DEBUG] Delete finished successfully")
	return nil

}
