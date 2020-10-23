package mso

import (
	"fmt"
	"log"
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
							Required: true,
						},
					},
				},
				Required: true,
			},
			"directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		}),
	}
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

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateContract := d.Get("contract_name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			contractCount, err := tempCont.ArrayCount("contracts")
			if err != nil {
				return fmt.Errorf("Unable to get contract list")
			}
			for j := 0; j < contractCount; j++ {
				contractCont, err := tempCont.ArrayElement(j, "contracts")
				if err != nil {
					return err
				}
				apiContract := models.StripQuotes(contractCont.S("name").String())
				if apiContract == stateContract {
					d.SetId(apiContract)
					d.Set("contract_name", apiContract)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(contractCont.S("displayName").String()))
					d.Set("filter_type", models.StripQuotes(contractCont.S("filterType").String()))
					d.Set("scope", models.StripQuotes(contractCont.S("scope").String()))

					count, _ := contractCont.ArrayCount("filterRelationships")

					var filterSchema string
					var filterTemplate string
					var filterName string

					if tempVar, ok := d.GetOk("filter_relationships"); ok {
						filter_relationships := tempVar.(map[string]interface{})

						if filter_relationships["filter_schema_id"] != nil {
							filterSchema = filter_relationships["filter_schema_id"].(string)
						} else {
							filterSchema = schemaId
						}

						if filter_relationships["filter_template_name"] != nil {
							filterTemplate = filter_relationships["filter_template_name"].(string)
						} else {
							filterTemplate = apiTemplate
						}

						filterName = filter_relationships["filter_name"].(string)
					}

					flag := false
					for i := 0; i < count; i++ {
						filterCont, err := contractCont.ArrayElement(i, "filterRelationships")
						if err != nil {
							return fmt.Errorf("Unable to parse the filter Relationships list")
						}

						d.Set("directives", filterCont.S("directives").Data().([]interface{}))
						filRef := filterCont.S("filterRef").Data()
						split := strings.Split(filRef.(string), "/")

						if split[2] == filterSchema && split[4] == filterTemplate && split[6] == filterName {
							flag = true
						}
					}

					if flag {
						d.Set("filter_relationships", d.Get("filter_relationships"))
					} else {
						d.Set("filter_relationships", make(map[string]interface{}))
					}

					found = true
					break
				}
			}
		}
	}

	if !found {
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

	path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
	contractStruct := models.NewTemplateContract("remove", path, contractName, displayName, scope, filter_type, filter)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return resourceMSOTemplateContractRead(d, m)
}
