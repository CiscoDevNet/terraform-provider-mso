package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOSchemaTemplateVrf() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaTemplateVrfRead,

		SchemaVersion: version,
		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"vzany": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		}),
	}
}
func datasourceMSOSchemaTemplateVrfRead(d *schema.ResourceData, m interface{}) error {
	schemaId := d.Get("schema_id").(string)
	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	templateName := d.Get("template").(string)
	vrfName := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return fmt.Errorf("No Vrf found")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				currentVrfName := models.StripQuotes(vrfCont.S("name").String())
				if currentVrfName == vrfName {
					d.SetId(currentVrfName)
					d.Set("name", currentVrfName)
					d.Set("display_name", models.StripQuotes(vrfCont.S("displayName").String()))
					if vrfCont.Exists("l3MCast") {
						l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
						d.Set("layer3_multicast", l3Mcast)
					}
					if vrfCont.Exists("vzAnyEnabled") {
						vzAnyEnabled, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("vzAnyEnabled").String()))
						d.Set("vzany", vzAnyEnabled)
					}
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		d.SetId("")
		d.Set("name", "")
		d.Set("display_name", "")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
