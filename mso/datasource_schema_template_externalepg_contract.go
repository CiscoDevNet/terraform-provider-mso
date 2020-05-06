package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceMSOTemplateExternalEpgContract() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOTemplateExternalEpgContractRead,

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
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"relationship_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
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
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateEPG := d.Get("external_epg_name").(string)
	stateContract := d.Get("contract_name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			epgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get External Epg list")
			}
			for j := 0; j < epgCount; j++ {
				epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				apiEpg := models.StripQuotes(epgCont.S("name").String())
				if apiEpg == stateEPG {
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
						if stateContract == (fmt.Sprintf("%s", split[3])) {
							d.SetId(fmt.Sprintf("%s", split[3]))
							d.Set("contract_name", fmt.Sprintf("%s", split[3]))
							d.Set("contract_schema_id", fmt.Sprintf("%s", split[1]))
							d.Set("contract_template_name", fmt.Sprintf("%s", split[2]))
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
		return fmt.Errorf("External Epg Contract Not Found")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
