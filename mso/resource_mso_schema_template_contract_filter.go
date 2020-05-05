package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOTemplateContractFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateContractFilterCreate,
		Read:   resourceMSOTemplateContractFilterRead,
		Update: resourceMSOTemplateContractFilterUpdate,
		Delete: resourceMSOTemplateContractFilterDelete,

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
				Type:     schema.TypeString,
				Required: true,

				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,

				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
				Optional: true,
			},
			"directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"filter_relationships_procon": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"procon_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"procon_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"procon_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"procon_directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"filter_relationships_conpro": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"conpro_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"conpro_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"conpro_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"conpro_directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		}),
	}
}

func resourceMSOTemplateContractFilterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template ContractFilter: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	err := resourceMSOTemplateContractFilterRead(d, m)
	
	if err == nil {

		err1 := resourceMSOTemplateContractFilterUpdate(d, m)
		if err1 != nil {
			return err1
		}
	} else {

		var scope, filter_type, contract_schema_id, contract_template_name string

		if tempVar, ok := d.GetOk("scope"); ok {
			scope = tempVar.(string)
		}
		if tempVar, ok := d.GetOk("filter_type"); ok {
			filter_type = tempVar.(string)
		}

		if tempVar, ok := d.GetOk("contract_schema_id"); ok {
			contract_schema_id = tempVar.(string)
		} else {
			contract_schema_id = schemaID
		}
		if tempVar, ok := d.GetOk("contract_template_name"); ok {
			contract_template_name = tempVar.(string)
		} else {
			contract_template_name = templateName
		}

		contractRefMap := make(map[string]interface{})
		contractRefMap["schemaId"] = contract_schema_id
		contractRefMap["templateName"] = contract_template_name
		contractRefMap["contractName"] = contractName

		filter := make([]interface{}, 0)
		filterMap := make(map[string]interface{})
		if tempVar, ok := d.GetOk("filter_relationships"); ok {
			filter_relationships := tempVar.(map[string]interface{})

			filterRefMap := make(map[string]interface{})

			if filter_relationships["filter_schema_id"] != nil {
				filterRefMap["schemaId"] = filter_relationships["filter_schema_id"].(string)
			} else {
				filterRefMap["schemaId"] = schemaID
			}

			if filter_relationships["filter_template_name"] != nil {
				filterRefMap["templateName"] = filter_relationships["filter_template_name"].(string)
			} else {
				filterRefMap["templateName"] = templateName
			}

			filterRefMap["filterName"] = filter_relationships["filter_name"].(string)

			filterMap["filterRef"] = filterRefMap
			if tempVar1, ok := d.GetOk("directives"); ok {
				filterMap["directives"] = tempVar1.([]interface{})
			}
			filter = append(filter, filterMap)
		}

	
		filterProCon := make([]interface{}, 0)
		filterProConMap := make(map[string]interface{})
		if tempVar, ok := d.GetOk("filter_relationships_procon"); ok {
			filter_relationships := tempVar.(map[string]interface{})

			filterRefMap := make(map[string]interface{})
		
			if filter_relationships["procon_schema_id"] != nil {
				filterRefMap["schemaId"] = fmt.Sprintf("%v", filter_relationships["procon_schema_id"])
			} else {
				filterRefMap["schemaId"] = schemaID
			}
		
			if filter_relationships["procon_template_name"] != nil {
				filterRefMap["templateName"] = fmt.Sprintf("%v", filter_relationships["procon_template_name"])
			} else {
				filterRefMap["templateName"] = templateName
			}

			filterRefMap["filterName"] = fmt.Sprintf("%v", filter_relationships["procon_name"])
			
			filterProConMap["filterRef"] = filterRefMap
			if tempVar1, ok := d.GetOk("procon_directives"); ok {
				filterProConMap["directives"] = tempVar1
			}

		} else {
			filterProConMap = nil
		}
		
		filterProCon = append(filterProCon, filterProConMap)

		filterConPro := make([]interface{}, 0)
		filterConProMap := make(map[string]interface{})
		if tempVar, ok := d.GetOk("filter_relationships_conpro"); ok {
			filter_relationships := tempVar.(map[string]interface{})

			filterRefMap := make(map[string]interface{})

			if filter_relationships["conpro_schema_id"] != nil {
				filterRefMap["schemaId"] = fmt.Sprintf("%v", filter_relationships["conpro_schema_id"])
			} else {
				filterRefMap["schemaId"] = schemaID
			}

			if filter_relationships["conpro_template_name"] != nil {
				filterRefMap["templateName"] = fmt.Sprintf("%v", filter_relationships["conpro_template_name"])
			} else {
				filterRefMap["templateName"] = templateName
			}

			filterRefMap["filterName"] = fmt.Sprintf("%v", filter_relationships["conpro_name"])

			filterConProMap["filterRef"] = filterRefMap
			if tempVar1, ok := d.GetOk("conpro_directives"); ok {
				filterConProMap["directives"] = tempVar1
			}

		} else {
			filterConProMap = nil
		}
		filterConPro = append(filterConPro, filterConProMap)

		path := fmt.Sprintf("/templates/%s/contracts/-", templateName)
		contractStruct := models.NewTemplateContractFilter("add", path, contractName, displayName, scope, filter_type, contractRefMap, filter, filterProCon, filterConPro)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

		if err != nil {
			return err
		}
	}
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractFilterRead(d *schema.ResourceData, m interface{}) error {
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

					contractRef := models.StripQuotes(contractCont.S("contractRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
					match := re.FindStringSubmatch(contractRef)
					d.Set("contract_name", match[3])
					d.Set("contract_schema_id", match[1])
					d.Set("contract_template_name", match[2])

					if contractCont.Exists("filterRelationships") {
						count, _ := contractCont.ArrayCount("filterRelationships")

						filterMap := make(map[string]interface{})
						for i := 0; i < count; i++ {
							filterCont, err := contractCont.ArrayElement(i, "filterRelationships")
							if err != nil {
								return fmt.Errorf("Unable to parse the filter Relationships list")
							}
							if filterCont.Exists("directives") {
								d.Set("directives", filterCont.S("directives").Data().([]interface{}))
							}
							if filterCont.Exists("filterRef") {
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")

								filterMap["filter_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["filter_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["filter_name"] = fmt.Sprintf("%s", split[6])
							}
						}
						d.Set("filter_relationships", filterMap)
					}

					if contractCont.Exists("filterRelationshipsProviderToConsumer") {
						count, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")

						filterMap := make(map[string]interface{})
						for i := 0; i < count; i++ {
							filterCont, err := contractCont.ArrayElement(i, "filterRelationshipsProviderToConsumer")
							if err != nil {
								return fmt.Errorf("Unable to parse the filter Relationships Provider to Consumer list")
							}
							if filterCont.Exists("directives") {
								d.Set("procon_directives", filterCont.S("directives").Data().([]interface{}))
							}
							if filterCont.Exists("filterRef") {
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")
								
								filterMap["procon_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["procon_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["procon_name"] = fmt.Sprintf("%s", split[6])
							}
						}
						d.Set("filter_relationships_procon", filterMap)
					}

					if contractCont.Exists("filterRelationshipsConsumerToProvider") {
						count, _ := contractCont.ArrayCount("filterRelationshipsConsumerToProvider")

						filterMap := make(map[string]interface{})
						for i := 0; i < count; i++ {
							filterCont, err := contractCont.ArrayElement(i, "filterRelationshipsConsumerToProvider")
							if err != nil {
								return fmt.Errorf("Unable to parse the filter Relationships Consumer to Provider list")
							}
							if filterCont.Exists("directives") {
								d.Set("conpro_directives", filterCont.S("directives").Data().([]interface{}))
							}
							if filterCont.Exists("filterRef") {
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")

								filterMap["conpro_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["conpro_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["conpro_name"] = fmt.Sprintf("%s", split[6])
							}
						}
						d.Set("filter_relationships_conpro", filterMap)
					}

					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Not Found")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTemplateContractFilterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	var scope, filter_type, contract_schema_id, contract_template_name string

	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schema_id = tempVar.(string)
	} else {
		contract_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_template_name = tempVar.(string)
	} else {
		contract_template_name = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schema_id
	contractRefMap["templateName"] = contract_template_name
	contractRefMap["contractName"] = contractName

	filter := make([]interface{}, 0)
	filterMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})

		if filter_relationships["filter_schema_id"] != nil {
			filterRefMap["schemaId"] = filter_relationships["filter_schema_id"].(string)
		} else {
			filterRefMap["schemaId"] = schemaID
		}

		if filter_relationships["filter_template_name"] != nil {
			filterRefMap["templateName"] = filter_relationships["filter_template_name"].(string)
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = filter_relationships["filter_name"].(string)

		filterMap["filterRef"] = filterRefMap
		if tempVar1, ok := d.GetOk("directives"); ok {
			filterMap["directives"] = tempVar1.([]interface{})

		}
		filter = append(filter, filterMap)
	}
	

	filterProCon := make([]interface{}, 0)
	filterProConMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships_procon"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})
		
		if filter_relationships["procon_schema_id"] != nil {
			filterRefMap["schemaId"] = fmt.Sprintf("%v", filter_relationships["procon_schema_id"])
		} else {
			filterRefMap["schemaId"] = schemaID
		}
		
		if filter_relationships["procon_template_name"] != nil {
			filterRefMap["templateName"] = fmt.Sprintf("%v", filter_relationships["procon_template_name"])
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = fmt.Sprintf("%v", filter_relationships["procon_name"])

		filterProConMap["filterRef"] = filterRefMap
		if tempVar1, ok := d.GetOk("procon_directives"); ok {
			filterProConMap["directives"] = tempVar1
		}

	} else {
		filterProConMap = nil
	}
	
	filterProCon = append(filterProCon, filterProConMap)

	filterConPro := make([]interface{}, 0)
	filterConProMap := make(map[string]interface{})
	if tempVar, ok := d.GetOk("filter_relationships_conpro"); ok {
		filter_relationships := tempVar.(map[string]interface{})

		filterRefMap := make(map[string]interface{})

		if filter_relationships["conpro_schema_id"] != nil {
			filterRefMap["schemaId"] = fmt.Sprintf("%v", filter_relationships["conpro_schema_id"])
		} else {
			filterRefMap["schemaId"] = schemaID
		}

		if filter_relationships["conpro_template_name"] != nil {
			filterRefMap["templateName"] = fmt.Sprintf("%v", filter_relationships["conpro_template_name"])
		} else {
			filterRefMap["templateName"] = templateName
		}

		filterRefMap["filterName"] = fmt.Sprintf("%v", filter_relationships["conpro_name"])

		filterConProMap["filterRef"] = filterRefMap
		if tempVar1, ok := d.GetOk("conpro_directives"); ok {
			filterConProMap["directives"] = tempVar1
		}

	} else {
		filterConProMap = nil
	}
	filterConPro = append(filterConPro, filterConProMap)

	path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
	contractStruct := models.NewTemplateContractFilter("replace", path, contractName, displayName, scope, filter_type, contractRefMap, filter, filterProCon, filterConPro)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateContractRead(d, m)
}

func resourceMSOTemplateContractFilterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Filter: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	displayName := d.Get("display_name").(string)

	var scope, filter_type, contract_schema_id, contract_template_name string

	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schema_id = tempVar.(string)
	} else {
		contract_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_template_name = tempVar.(string)
	} else {
		contract_template_name = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schema_id
	contractRefMap["templateName"] = contract_template_name
	contractRefMap["contractName"] = contractName

	filter := make([]interface{}, 0)
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
		if tempVar1, ok := d.GetOk("directives"); ok {
			filterMap["directives"] = tempVar1
		}

	} else {
		filterMap = nil
	}
	filter = append(filter, filterMap)

	path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
	contractStruct := models.NewTemplateContractFilter("remove", path, contractName, displayName, scope, filter_type, contractRefMap, filter, nil, nil)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return resourceMSOTemplateContractRead(d, m)
}
