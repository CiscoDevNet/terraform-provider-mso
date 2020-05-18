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

func dataSourceMSOTemplateContractFilter() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOTemplateContractFilterRead,

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
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"filter_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"filter_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"filter_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		}),
	}
}

func dataSourceMSOTemplateContractFilterRead(d *schema.ResourceData, m interface{}) error {
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
					}

				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the given Contract Filter")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
