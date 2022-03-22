package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
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
			"filter_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter_relationships": {
				Type: schema.TypeMap,
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
				Optional:      true,
				ConflictsWith: []string{"filter_relationship"},
				Deprecated:    "use filter_relationship instead",
			},
			"filter_relationship": {
				Type:     schema.TypeList,
				Optional: true,
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
					},
				},
			},
			"directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		}),
	}
}

func resourceMSOTemplateContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)

	get_attribute := strings.Split(d.Id(), "/")
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", get_attribute[0]))
	if err != nil {
		return nil, err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateContract := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			contractCount, err := tempCont.ArrayCount("contracts")
			if err != nil {
				return nil, fmt.Errorf("Unable to get contract list")
			}
			for j := 0; j < contractCount; j++ {
				contractCont, err := tempCont.ArrayElement(j, "contracts")
				if err != nil {
					return nil, err
				}
				apiContract := models.StripQuotes(contractCont.S("name").String())
				if apiContract == stateContract {
					d.SetId(get_attribute[4])
					d.Set("contract_name", apiContract)
					d.Set("schema_id", get_attribute[0])
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(contractCont.S("displayName").String()))
					d.Set("filter_type", models.StripQuotes(contractCont.S("filterType").String()))
					d.Set("scope", models.StripQuotes(contractCont.S("scope").String()))

					count, _ := contractCont.ArrayCount("filterRelationships")
					filterMap := make(map[string]interface{})
					for i := 0; i < count; i++ {
						filterCont, err := contractCont.ArrayElement(i, "filterRelationships")
						if err != nil {
							return nil, fmt.Errorf("Unable to parse the filter Relationships list")
						}

						d.Set("directives", filterCont.S("directives").Data().([]interface{}))
						filRef := filterCont.S("filterRef").Data()
						split := strings.Split(filRef.(string), "/")

						filterMap["filter_schema_id"] = fmt.Sprintf("%s", split[2])
						filterMap["filter_template_name"] = fmt.Sprintf("%s", split[4])
						filterMap["filter_name"] = fmt.Sprintf("%s", split[6])
					}
					d.Set("filter_relationships", filterMap)

					found = true
					break
				}
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("Unable to find the Contract %s", stateContract)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOTemplateContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	var scope, filter_type string

	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	filter := make([]interface{}, 0, 1)
	filterMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})

		if filter_relationships["filter_schema_id"] != nil {
			filterRefMap["schemaId"] = filter_relationships["filter_schema_id"]
		} else {
			filterRefMap["schemaId"] = schemaID
		}

		if filter_relationships["filter_template_name"] != nil {
			filterRefMap["templateName"] = filter_relationships["filter_template_name"]
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = filter_relationships["filter_name"]

		filterMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterMap["directives"] = tempVar
		}

	} else {
		filterMap = nil
	}
	filter = append(filter, filterMap)

	var filterList []interface{}
	filter_relationship := d.Get("filter_relationship").([]interface{})
	for _, tempFilter := range filter_relationship {
		filterRelationship := tempFilter.(map[string]interface{})
		filterRelMap := make(map[string]interface{})
		filterRefMap := make(map[string]interface{})

		filterRefMap["schemaId"] = filterRelationship["filter_schema_id"]
		filterRefMap["templateName"] = filterRelationship["filter_template_name"]
		filterRefMap["filterName"] = filterRelationship["filter_name"]

		filterRelMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterRelMap["directives"] = tempVar
		}

		filterList = append(filterList, filterRelMap)
	}

	filter_check := filter[0].(map[string]interface{})
	if filter_check == nil {
		filter = filterList
	}

	path := fmt.Sprintf("/templates/%s/contracts/-", templateName)
	contractStruct := models.NewTemplateContract("add", path, contractName, displayName, scope, filter_type, filter)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	templateCount := 0
	if cont.Exists("templates") {
		templateCount = len(cont.S("templates").Data().([]interface{}))
	} else {
		return fmt.Errorf("No Template found")
	}

	templateState := d.Get("template_name").(string)
	contractNameState := d.Get("contract_name").(string)
	foundContract := false

	for i := 0; i < templateCount; i++ {
		templateCont := cont.S("templates").Index(i)
		templateNameCont := models.StripQuotes(templateCont.S("name").String())

		if templateState == templateNameCont {
			contractCount := 0
			if templateCont.Exists("contracts") {
				contractCount = len(templateCont.S("contracts").Data().([]interface{}))
			} else {
				return fmt.Errorf("Unable to get contract list")
			}
			for j := 0; j < contractCount; j++ {
				contractCont := templateCont.S("contracts").Index(j)
				contractNameCont := models.StripQuotes(contractCont.S("name").String())

				if contractNameCont == contractNameState {
					foundContract = true
					d.SetId(contractNameCont)
					d.Set("contract_name", contractNameCont)
					d.Set("schema_id", schemaId)
					d.Set("template_name", templateNameCont)
					d.Set("display_name", models.StripQuotes(contractCont.S("displayName").String()))
					d.Set("filter_type", models.StripQuotes(contractCont.S("filterType").String()))
					d.Set("scope", models.StripQuotes(contractCont.S("scope").String()))
					filterRelationshipCount := 0
					if contractCont.Exists("filterRelationships") {
						filterRelationshipCount = len(contractCont.S("filterRelationships").Data().([]interface{}))
						filterRelationships := make([]map[string]string, 0)

						for k := 0; k < filterRelationshipCount; k++ {
							filterCont := contractCont.S("filterRelationships").Index(k)
							d.Set("directives", filterCont.S("directives").Data().([]interface{}))
							filterRelation := make(map[string]string)
							filterRef := models.G(filterCont, "filterRef")
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/filters/(.*)")
							match := re.FindStringSubmatch(filterRef)
							filterRelation["filter_schema_id"] = match[1]
							filterRelation["filter_template_name"] = match[2]
							filterRelation["filter_name"] = match[3]
							filterRelationships = append(filterRelationships, filterRelation)
						}
						if _, ok := d.GetOk("filter_relationships"); !ok {
							d.Set("filter_relationship", filterRelationships)
						}
					}
				}
			}
		}
	}
	if !foundContract {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTemplateContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	var scope, filter_type string

	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	filter := make([]interface{}, 0, 1)
	filterMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})

		if filter_relationships["filter_schema_id"] != nil {
			filterRefMap["schemaId"] = filter_relationships["filter_schema_id"]
		} else {
			filterRefMap["schemaId"] = schemaID
		}

		if filter_relationships["filter_template_name"] != nil {
			filterRefMap["templateName"] = filter_relationships["filter_template_name"]
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = filter_relationships["filter_name"]

		filterMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterMap["directives"] = tempVar
		}

	} else {
		filterMap = nil
	}
	filter = append(filter, filterMap)

	var filterList []interface{}
	filter_relationship := d.Get("filter_relationship").([]interface{})
	for _, tempFilter := range filter_relationship {
		filterRelationship := tempFilter.(map[string]interface{})
		filterRelMap := make(map[string]interface{})
		filterRefMap := make(map[string]interface{})

		filterRefMap["schemaId"] = filterRelationship["filter_schema_id"]
		filterRefMap["templateName"] = filterRelationship["filter_template_name"]
		filterRefMap["filterName"] = filterRelationship["filter_name"]

		filterRelMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterRelMap["directives"] = tempVar
		}

		filterList = append(filterList, filterRelMap)
	}

	filter_check := filter[0].(map[string]interface{})
	if filter_check == nil {
		filter = filterList
	}

	path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
	contractStruct := models.NewTemplateContract("replace", path, contractName, displayName, scope, filter_type, filter)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	var scope, filter_type string

	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	filter := make([]interface{}, 0, 1)
	filterMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})

		if filter_relationships["filter_schema_id"] != nil {
			filterRefMap["schemaId"] = filter_relationships["filter_schema_id"]
		} else {
			filterRefMap["schemaId"] = schemaID
		}

		if filter_relationships["filter_template_name"] != nil {
			filterRefMap["templateName"] = filter_relationships["filter_template_name"]
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = filter_relationships["filter_name"]

		filterMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterMap["directives"] = tempVar
		}

	} else {
		filterMap = nil
	}
	filter = append(filter, filterMap)

	var filterList []interface{}
	filter_relationship := d.Get("filter_relationship").([]interface{})
	for _, tempFilter := range filter_relationship {
		filterRelationship := tempFilter.(map[string]interface{})
		filterRelMap := make(map[string]interface{})
		filterRefMap := make(map[string]interface{})

		filterRefMap["schemaId"] = filterRelationship["filter_schema_id"]
		filterRefMap["templateName"] = filterRelationship["filter_template_name"]
		filterRefMap["filterName"] = filterRelationship["filter_name"]

		filterRelMap["filterRef"] = filterRefMap
		if tempVar, ok := d.GetOk("directives"); ok {
			filterRelMap["directives"] = tempVar
		}

		filterList = append(filterList, filterRelMap)
	}

	filter_check := filter[0].(map[string]interface{})
	if filter_check == nil {
		filter = filterList
	}

	path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
	contractStruct := models.NewTemplateContract("remove", path, contractName, displayName, scope, filter_type, filter)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return resourceMSOTemplateContractRead(d, m)
}
