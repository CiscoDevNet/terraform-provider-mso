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
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"filter_relationships_provider_to_consumer": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"provider_to_consumer_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"provider_to_consumer_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"provider_to_consumer_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"provider_to_consumer_directives": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"filter_relationships_consumer_to_provider": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"consumer_to_provider_schema_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"consumer_to_provider_template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"consumer_to_provider_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"consumer_to_provider_directives": {
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
								d.Set("provider_to_consumer_directives", filterCont.S("directives").Data().([]interface{}))
							}
							if filterCont.Exists("filterRef") {
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")

								filterMap["provider_to_consumer_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["provider_to_consumer_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["provider_to_consumer_name"] = fmt.Sprintf("%s", split[6])
							}
						}
						d.Set("filter_relationships_provider_to_consumer", filterMap)
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
								d.Set("consumer_to_provider_directives", filterCont.S("directives").Data().([]interface{}))
							}
							if filterCont.Exists("filterRef") {
								filRef := filterCont.S("filterRef").Data()
								split := strings.Split(filRef.(string), "/")
								
								filterMap["consumer_to_provider_schema_id"] = fmt.Sprintf("%s", split[2])
								filterMap["consumer_to_provider_template_name"] = fmt.Sprintf("%s", split[4])
								filterMap["consumer_to_provider_name"] = fmt.Sprintf("%s", split[6])
							}
						}
						d.Set("filter_relationships_consumer_to_provider", filterMap)
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
