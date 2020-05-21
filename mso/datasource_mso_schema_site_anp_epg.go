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

func dataSourceMSOSchemaSiteAnpEpg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteAnpEpgRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOSchemaSiteAnpEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	stateSite := d.Get("site_id").(string)
	found := false
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
				split := strings.Split(apiAnpRef, "/")
				apiAnp := split[6]
				if apiAnp == stateAnp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
						split := strings.Split(apiEpgRef, "/")
						apiEPG := split[8]
						if apiEPG == stateEpg {
							d.SetId(apiEPG)
							d.Set("site_id", apiSite)
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)

							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site Anp Epg %s", stateEpg)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
