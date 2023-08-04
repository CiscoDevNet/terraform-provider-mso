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

func dataSourceMSOTemplateExternalEpgContract() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOTemplateExternalEpgContractRead,

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
			"external_epg_name": &schema.Schema{
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

func dataSourceMSOTemplateExternalEpgContractRead(d *schema.ResourceData, m interface{}) error {
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
	epg := d.Get("external_epg_name").(string)
	contractName := d.Get("contract_name").(string)
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
			epgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get External Epg list")
			}
			for j := 0; j < epgCount && !found; j++ {
				epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				currentEpg := models.StripQuotes(epgCont.S("name").String())
				if currentEpg == epg {
					contractCount, err := epgCont.ArrayCount("contractRelationships")
					if err != nil {
						return fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
						if err != nil {
							return err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if contractName == split[3] && contractSchemaId == split[1] && contractTemplateName == split[2] {
							d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s-%s-%s", schemaId, template, epg, contractSchemaId, contractTemplateName, contractName))
							d.Set("contract_name", contractName)
							d.Set("contract_schema_id", contractSchemaId)
							d.Set("contract_template_name", contractTemplateName)
							d.Set("relationship_type", models.StripQuotes(contractCont.S("relationshipType").String()))
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the External EPG Contract %s in Template %s of Schema Id %s", contractName, contractTemplateName, contractSchemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
