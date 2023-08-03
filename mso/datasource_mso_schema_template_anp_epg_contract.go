package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateAnpEpgContract() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateAnpEpgContractRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relationship_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}

}

func dataSourceMSOTemplateAnpEpgContractRead(d *schema.ResourceData, m interface{}) error {
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
	stateANP := d.Get("anp_name")
	stateEPG := d.Get("epg_name")
	stateContract := d.Get("contract_name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							crefCount, err := epgCont.ArrayCount("contractRelationships")
							if err != nil {
								return fmt.Errorf("Unable to get the contract relationships list")
							}
							for l := 0; l < crefCount; l++ {
								crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
								if err != nil {
									return err
								}
								contractRef := models.StripQuotes(crefCont.S("contractRef").String())
								apiRelationshipType := models.StripQuotes(crefCont.S("relationshipType").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								apiContract := match[3]
								if apiContract == stateContract {
									d.SetId(apiContract)
									d.Set("contract_name", match[3])
									d.Set("contract_schema_id", match[1])
									d.Set("contract_template_name", match[2])
									d.Set("relationship_type", apiRelationshipType)
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

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the Contract %s", stateContract)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
