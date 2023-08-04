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

func dataSourceMSOTemplateL3out() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateL3outRead,

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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateL3outRead(d *schema.ResourceData, m interface{}) error {
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
	stateL3out := d.Get("l3out_name")

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())
		if apiTemplate == stateTemplate {
			l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
			if err != nil {
				return fmt.Errorf("Unable to get L3out list")
			}
			for j := 0; j < l3outCount; j++ {
				l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
				if err != nil {
					return err
				}
				apiL3out := models.StripQuotes(l3outCont.S("name").String())
				if apiL3out == stateL3out {
					d.SetId(fmt.Sprintf("%s/templates/%s/intersiteL3outs/%s", schemaId, stateTemplate, stateL3out))
					d.Set("l3out_name", apiL3out)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(l3outCont.S("displayName").String()))

					vrfRef := models.StripQuotes(l3outCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					found = true
					break
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the L3out %s in Template %s of Schema Id %s", stateL3out, stateTemplate, schemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
