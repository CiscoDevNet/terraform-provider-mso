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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
	template := d.Get("template_name").(string)
	anp := d.Get("anp_name")
	epg := d.Get("epg_name")
	contract := d.Get("contract_name")
	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaId
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = template
	}
	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplate := models.StripQuotes(tempCont.S("name").String())

		if currentTemplate == template {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount && !found; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				currentAnp := models.StripQuotes(anpCont.S("name").String())
				if currentAnp == anp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount && !found; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpg := models.StripQuotes(epgCont.S("name").String())
						if currentEpg == epg {
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
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								if match[3] == contract && match[1] == contractSchemaId && match[2] == contractTemplateName {
									d.SetId(fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/contracts/%s-%s-%s", schemaId, template, anp, epg, contractSchemaId, contractTemplateName, contract))
									d.Set("contract_name", contract)
									d.Set("contract_schema_id", contractSchemaId)
									d.Set("contract_template_name", contractTemplateName)
									d.Set("relationship_type", models.StripQuotes(crefCont.S("relationshipType").String()))
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
		return fmt.Errorf("Unable to find the ANP EPG Contract %s in Template %s of Schema Id %s ", contract, contractTemplateName, contractSchemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
