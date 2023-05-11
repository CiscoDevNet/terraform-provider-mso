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

func dataSourceMSOSchemaSiteExternalEpg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteExternalEpgRead,

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
			"site_id": &schema.Schema{
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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOSchemaSiteExternalEpgRead(d *schema.ResourceData, m interface{}) error {
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
	stateSiteId := d.Get("site_id").(string)
	found := false
	stateExternalEpg := d.Get("external_epg_name").(string)
	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if apiSiteId == stateSiteId {
			externalEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get External EPG list")
			}
			for j := 0; j < externalEpgCount; j++ {
				externalEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				externalEpgRef := models.StripQuotes(externalEpgCont.S("externalEpgRef").String())
				re := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/externalEpgs/(.*)")
				match := re.FindStringSubmatch(externalEpgRef)
				log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead externalEpgRef: %s match: %s", externalEpgRef, match)
				if len(match) >= 4 {
					if match[3] == stateExternalEpg {
						d.SetId(match[3])
						d.Set("external_epg_name", match[3])
						d.Set("schema_id", match[1])
						d.Set("template_name", match[2])
						d.Set("site_id", apiSiteId)

						l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
						if l3outRef != "{}" && l3outRef != "" {
							reL3out := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/l3outs/(.*)")
							matchL3out := reL3out.FindStringSubmatch(l3outRef)
							log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead l3outRef: %s matchL3out: %s", l3outRef, matchL3out)
							if len(matchL3out) >= 4 {
								d.Set("l3out_name", matchL3out[3])
							} else {
								return fmt.Errorf("Error in parsing l3outRef to get L3Out name")
							}
						}

						found = true
						break
					}
				} else {
					return fmt.Errorf("Error in parsing externalEpgRef to get External EPG name")
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
