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

func datasourceMSOSchemaSiteExternalEpgSelector() *schema.Resource {
	return &schema.Resource{
		Read: datasourceMSOSchemaSiteExternalEpgSelectorRead,

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
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

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func datasourceMSOSchemaSiteExternalEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Data source Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Get("name").(string)
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	found := false

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		currSite := models.StripQuotes(siteCont.S("siteId").String())
		currTemplate := models.StripQuotes(siteCont.S("templateName").String())

		if currSite == siteID && currTemplate == templateName {
			extEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return err
			}

			for j := 0; j < extEpgCount; j++ {
				extEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}

				extEpgRef := models.StripQuotes(extEpgCont.S("externalEpgRef").String())
				tokens := strings.Split(extEpgRef, "/")
				extEpgName := tokens[len(tokens)-1]
				if extEpgName == externalEpgName {
					subnetCount, err := extEpgCont.ArrayCount("subnets")
					if err != nil {
						return err
					}

					for k := 0; k < subnetCount; k++ {
						subnetCont, err := extEpgCont.ArrayElement(k, "subnets")
						if err != nil {
							return err
						}

						subnetName := models.StripQuotes(subnetCont.S("name").String())
						if subnetName == dn {
							found = true
							d.SetId(dn)
							d.Set("name", subnetName)
							d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
							break
						}
					}
				}
				if found {
					d.Set("external_epg_name", extEpgName)
					break
				}
			}
		}
		if found {
			d.Set("site_id", siteID)
			d.Set("template_name", templateName)
			break
		}
	}

	if found {
		d.Set("schema_id", schemaID)
	} else {
		d.SetId("")
		return fmt.Errorf("Selector of specified name not found")
	}
	log.Printf("[DEBUG] %s: Datasource Read finished successfully", d.Id())
	return nil
}
