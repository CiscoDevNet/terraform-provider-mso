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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"filter_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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

	var filter_type string
	if tempVar, ok := d.GetOk("filter_type"); ok {
		filter_type = tempVar.(string)
	}

	var filterId, filterTemplate, filterName string
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterId = tempVar.(string)
	} else {
		filterId = schemaID
	}

	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplate = tempVar.(string)
	} else {
		filterTemplate = templateName
	}

	if tempVar, ok := d.GetOk("filter_name"); ok {
		filterName = tempVar.(string)
	}

	filterRefMap := make(map[string]interface{})
	filterRefMap["schemaId"] = filterId
	filterRefMap["templateName"] = filterTemplate
	filterRefMap["filterName"] = filterName

	var directives []interface{}
	if tempVar, ok := d.GetOk("directives"); ok {
		directives = tempVar.([]interface{})
	}

	if filter_type == "bothWay" {
		path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
		contractStruct := models.NewTemplateContractFilter("replace", path, filter_type)
		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
		if err != nil {
			return err
		}
		paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationships/-", templateName, contractName)
		contract := models.NewTemplateContractFilterRelationShip("add", paths, filterRefMap, directives)
		_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
		if err1 != nil {
			return err1
		}
	} else if filter_type == "provider_to_consumer" {
		path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
		contractStruct := models.NewTemplateContractFilter("replace", path, "oneWay")
		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
		if err != nil {
			return err
		}
		paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsProviderToConsumer/-", templateName, contractName)
		contract := models.NewTemplateContractFilterRelationShip("add", paths, filterRefMap, directives)
		_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
		if err1 != nil {
			return err1
		}
	} else if filter_type == "consumer_to_provider" {
		path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
		contractStruct := models.NewTemplateContractFilter("replace", path, "oneWay")
		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
		if err != nil {
			return err
		}
		paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsConsumerToProvider/-", templateName, contractName)
		contract := models.NewTemplateContractFilterRelationShip("add", paths, filterRefMap, directives)
		_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
		if err1 != nil {
			return err1
		}
	}else{
		return fmt.Errorf("Filter Type is not valid")
	}

	return resourceMSOTemplateContractFilterRead(d, m)
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
	stateFilterType := d.Get("filter_type")

	var stateFilterSchema, stateFilterTemplate string
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		stateFilterSchema = tempVar.(string)
	} else {
		stateFilterSchema = schemaId
	}

	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		stateFilterTemplate = tempVar.(string)
	} else {
		stateFilterTemplate = stateTemplate
	}

	stateName := d.Get("filter_name")
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
					d.Set("filter_type", stateFilterType)

					contractRef := models.StripQuotes(contractCont.S("contractRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
					match := re.FindStringSubmatch(contractRef)
					d.Set("contract_name", match[3])

					if stateFilterType == "bothWay" {
						if contractCont.Exists("filterRelationships") {
							filtercount, _ := contractCont.ArrayCount("filterRelationships")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationships")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == stateName && split[4] == stateFilterTemplate && split[2] == stateFilterSchema {
										d.Set("filter_name", split[6])
										d.Set("filter_template_name", split[4])
										d.Set("filter_schema_id", split[2])
										d.Set("directives", filterCont.S("directives").Data().([]interface{}))
										found = true
										break
									}
								}
							}
						}
					} else if stateFilterType == "provider_to_consumer" {
						if contractCont.Exists("filterRelationshipsProviderToConsumer") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsProviderToConsumer")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Provider to Consumer list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")

									if split[6] == stateName && split[4] == stateFilterTemplate && split[2] == stateFilterSchema {

										d.Set("filter_name", split[6])
										d.Set("filter_template_name", split[4])
										d.Set("filter_schema_id", split[2])
										d.Set("directives", filterCont.S("directives").Data().([]interface{}))
										found = true
										break
									}
								}
							}
						}
					} else if stateFilterType == "consumer_to_provider" {
						if contractCont.Exists("filterRelationshipsConsumerToProvider") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsConsumerToProvider")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsConsumerToProvider")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Consumer To Provider list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == stateName && split[4] == stateFilterTemplate && split[2] == stateFilterSchema {
										d.Set("filter_name", split[6])
										d.Set("filter_template_name", split[4])
										d.Set("filter_schema_id", split[2])
										d.Set("directives", filterCont.S("directives").Data().([]interface{}))
										found = true
										break
									}
								}
							}
						}
					}else{
						return fmt.Errorf("Filter Type is not valid")
					}
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

func resourceMSOTemplateContractFilterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Filter: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	var stateFilterType string

	if tempVar, ok := d.GetOk("filter_type"); ok {
		stateFilterType = tempVar.(string)
	}

	var filterId, filterTemplate, filterName string
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterId = tempVar.(string)
	} else {
		filterId = schemaID
	}

	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplate = tempVar.(string)
	} else {
		filterTemplate = templateName
	}

	if tempVar, ok := d.GetOk("filter_name"); ok {
		filterName = tempVar.(string)
	}

	filterRefMap := make(map[string]interface{})
	filterRefMap["schemaId"] = filterId
	filterRefMap["templateName"] = filterTemplate
	filterRefMap["filterName"] = filterName

	var directives []interface{}
	if tempVar, ok := d.GetOk("directives"); ok {
		directives = tempVar.([]interface{})
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	found := false
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == templateName {
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
				if apiContract == contractName {
					if stateFilterType == "bothWay" {
						if contractCont.Exists("filterRelationships") {
							filtercount, _ := contractCont.ArrayCount("filterRelationships")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationships")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true
										path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
										contractStruct := models.NewTemplateContractFilter("replace", path, stateFilterType)
										_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
										if err != nil {
											return err
										}
										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationships/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("replace", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}
										break
									}
								}
							}
						}
					} else if stateFilterType == "provider_to_consumer" {
						if contractCont.Exists("filterRelationshipsProviderToConsumer") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsProviderToConsumer")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Provider to Consumer list")
								}
								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true
										path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
										contractStruct := models.NewTemplateContractFilter("replace", path, "oneWay")
										_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
										if err != nil {
											return err
										}
										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsProviderToConsumer/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("replace", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}
										break
									}
								}
							}

						}

					} else if stateFilterType == "consumer_to_provider" {
						if contractCont.Exists("filterRelationshipsConsumerToProvider") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsConsumerToProvider")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsConsumerToProvider")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Consumer To Provider list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true
										path := fmt.Sprintf("/templates/%s/contracts/%s/filterType", templateName, contractName)
										contractStruct := models.NewTemplateContractFilter("replace", path, "oneWay")
										_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
										if err != nil {
											return err
										}
										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsConsumerToProvider/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("replace", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}

										break
									}
								}
							}
						}
					}else{
						return fmt.Errorf("Filter Type is not valid")
					}
				
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the given contract filter")
	}
	return resourceMSOTemplateContractFilterRead(d, m)
}

func resourceMSOTemplateContractFilterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Contract Filter: Beginning Delete")

	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	var stateFilterType string
	stateFilterType = d.Get("filter_type").(string)

	var filterId, filterTemplate, filterName string
	if tempVar, ok := d.GetOk("filter_schema_id"); ok {
		filterId = tempVar.(string)
	} else {
		filterId = schemaID
	}

	if tempVar, ok := d.GetOk("filter_template_name"); ok {
		filterTemplate = tempVar.(string)
	} else {
		filterTemplate = templateName
	}

	if tempVar, ok := d.GetOk("filter_name"); ok {
		filterName = tempVar.(string)
	}

	filterRefMap := make(map[string]interface{})
	filterRefMap["schemaId"] = filterId
	filterRefMap["templateName"] = filterTemplate
	filterRefMap["filterName"] = filterName

	var directives []interface{}
	if tempVar, ok := d.GetOk("directives"); ok {
		directives = tempVar.([]interface{})
	}

	indexf := 0
	indexfp := 0
	indexfc := 0
	found := false
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == templateName {
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
				if apiContract == contractName {
					if contractCont.Exists("filterRelationships") {
						countf, _ := contractCont.ArrayCount("filterRelationships")
						indexf = countf
					}
					if contractCont.Exists("filterRelationshipsProviderToConsumer") {
						countf, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")
						indexfp = countf
					}
					if contractCont.Exists("filterRelationshipsConsumerToProvider") {
						countf, _ := contractCont.ArrayCount("filterRelationshipsConsumerToProvider")
						indexfc = countf
					}

					if indexf+indexfp+indexfc <= 1 {
						path := fmt.Sprintf("/templates/%s/contracts/%s", templateName, contractName)
						contractStruct := models.NewTemplateContractFilter("remove", path, stateFilterType)

						_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
						if err != nil {
							return err
						}
						d.SetId("")
						return resourceMSOTemplateContractFilterRead(d, m)
					}

					if stateFilterType == "bothWay" {
						if contractCont.Exists("filterRelationships") {
							filtercount, _ := contractCont.ArrayCount("filterRelationships")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationships")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true
										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationships/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("remove", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}
										break
									}
								}
							}
						}
					} else if stateFilterType == "provider_to_consumer" {
						if contractCont.Exists("filterRelationshipsProviderToConsumer") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsProviderToConsumer")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Provider to Consumer list")
								}
								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true

										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsProviderToConsumer/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("remove", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}
										break
									}
								}
							}

						}

					} else if stateFilterType == "consumer_to_provider" {
						if contractCont.Exists("filterRelationshipsConsumerToProvider") {
							filtercount, _ := contractCont.ArrayCount("filterRelationshipsConsumerToProvider")
							for k := 0; k < filtercount; k++ {
								filterCont, err := contractCont.ArrayElement(k, "filterRelationshipsConsumerToProvider")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Consumer To Provider list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									if split[6] == filterName && split[4] == filterTemplate && split[2] == filterId {
										found = true
										paths := fmt.Sprintf("/templates/%s/contracts/%s/filterRelationshipsConsumerToProvider/%v", templateName, contractName, k)
										contract := models.NewTemplateContractFilterRelationShip("remove", paths, filterRefMap, directives)
										_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contract)
										if err1 != nil {
											return err1
										}

										break
									}
								}
							}
						}
					}else{
						return fmt.Errorf("Filter Type is not valid")
					}
				

				}
			}
		}
	}
	if !found {
		return fmt.Errorf("Unable to find given Contract Filter")
	}
	return resourceMSOTemplateContractFilterRead(d, m)

}
