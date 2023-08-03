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

func dataSourceMSOSchemaTemplateAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaTemplateAnpEpgUsegAttrRead,

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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"useg_subnet": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaTemplateAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
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

	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template_name", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					d.Set("anp_name", currentAnpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							d.Set("epg_name", currentEpgName)
							usegCount, err := epgCont.ArrayCount("uSegAttrs")
							if err != nil {
								return fmt.Errorf("No usegAttrs found")
							}
							for s := 0; s < usegCount; s++ {
								usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
								if err != nil {
									return err
								}
								currentName := models.StripQuotes(usegCont.S("name").String())
								if currentName == name {
									d.SetId(currentName)
									d.Set("name", currentName)
									d.Set("useg_type", models.StripQuotes(usegCont.S("type").String()))
									d.Set("value", models.StripQuotes(usegCont.S("value").String()))

									if usegCont.Exists("operator") {
										d.Set("operator", models.StripQuotes(usegCont.S("operator").String()))
									} else {
										d.Set("operator", "")
									}
									if usegCont.Exists("category") {
										d.Set("category", models.StripQuotes(usegCont.S("category").String()))
									} else {
										d.Set("category", "")
									}
									if usegCont.Exists("description") {
										d.Set("description", models.StripQuotes(usegCont.S("description").String()))
									} else {
										d.Set("description", "")
									}
									if usegCont.Exists("fvSubnet") {
										usegSubnet, _ := strconv.ParseBool(models.StripQuotes(usegCont.S("fvSubnet").String()))
										d.Set("useg_subnet", usegSubnet)
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
				}
				if found {
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
		d.Set("operator", "")
		d.Set("useg_type", "")
		d.Set("value", "")
		return fmt.Errorf("Unable to find Schema template anp epg useg attribute %s", name)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
