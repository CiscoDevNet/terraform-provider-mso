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

func dataSourceMSOSchemaSiteBdL3out() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteBdL3outRead,

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
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOSchemaSiteBdL3outRead(d *schema.ResourceData, m interface{}) error {
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
	stateBd := d.Get("bd_name").(string)
	stateL3out := d.Get("l3out_name").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
				split := strings.Split(apiBdRef, "/")
				apiBd := split[6]
				if apiBd == stateBd {
					d.Set("site_id", apiSite)
					d.Set("schema_id", split[2])
					d.Set("template_name", split[4])
					d.Set("bd_name", split[6])
					l3outCount, err := bdCont.ArrayCount("l3Outs")
					if err != nil {
						return fmt.Errorf("Unable to get l3Outs list")
					}
					for k := 0; k < l3outCount; k++ {
						l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
						if err != nil {
							return err
						}
						tempVar := l3outCont.String()
						apiL3out := strings.Trim(tempVar, "\"")
						if apiL3out == stateL3out {
							d.SetId(stateL3out)
							d.Set("l3out_name", apiL3out)
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site L3out %s", stateL3out)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
