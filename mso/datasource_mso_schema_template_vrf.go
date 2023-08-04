package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateVrf() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaTemplateVrfRead,

		SchemaVersion: version,
		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vzany": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ip_data_plane_learning": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
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
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
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
					d.SetId(fmt.Sprintf("%s/templates/%s/vrfs/%s", schemaId, templateName, vrfName))
					d.Set("name", currentVrfName)
					d.Set("template", currentTemplateName)
					d.Set("display_name", models.StripQuotes(vrfCont.S("displayName").String()))
					if vrfCont.Exists("l3MCast") {
						l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
						d.Set("layer3_multicast", l3Mcast)
					}
					if vrfCont.Exists("vzAnyEnabled") {
						vzAnyEnabled, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("vzAnyEnabled").String()))
						d.Set("vzany", vzAnyEnabled)
					}
					if vrfCont.Exists("ipDataPlaneLearning") {
						d.Set("ip_data_plane_learning", models.StripQuotes(vrfCont.S("ipDataPlaneLearning").String()))
					}
					if vrfCont.Exists("preferredGroup") {
						preferredGroup, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("preferredGroup").String()))
						d.Set("preferred_group", preferredGroup)
					}
					found = true
					break
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
