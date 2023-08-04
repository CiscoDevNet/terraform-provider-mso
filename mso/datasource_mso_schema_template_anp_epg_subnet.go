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

func dataSourceMSOSchemaTemplateAnpEpgSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaTemplateAnpEpgSubnetRead,

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
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaTemplateAnpEpgSubnetRead(d *schema.ResourceData, m interface{}) error {
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

	templateName := d.Get("template").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	ip := d.Get("ip").(string)

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount && !found; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount && !found; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("No Subnets found")
							}
							for s := 0; s < subnetCount; s++ {
								subnetCont, err := epgCont.ArrayElement(s, "subnets")
								if err != nil {
									return err
								}
								currentIp := models.StripQuotes(subnetCont.S("ip").String())
								if currentIp == ip {
									d.SetId(fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/subnets/%s", schemaId, templateName, anpName, epgName, ip))
									d.Set("ip", currentIp)
									d.Set("template", currentTemplateName)
									d.Set("anp_name", currentAnpName)
									d.Set("epg_name", currentEpgName)
									d.Set("description", models.StripQuotes(subnetCont.S("description").String()))

									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										shared, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("shared").String()))
										d.Set("shared", shared)
									}
									if subnetCont.Exists("primary") {
										primary, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("primary").String()))
										d.Set("primary", primary)
									}
									if subnetCont.Exists("noDefaultGateway") {
										noDefaultGateway, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("noDefaultGateway").String()))
										d.Set("no_default_gateway", noDefaultGateway)
									}
									if subnetCont.Exists("querier") {
										querier, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("querier").String()))
										d.Set("querier", querier)
									}

									found = true
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the ANP EPG Subnet %s in Template %s of Schema Id %s ", ip, templateName, schemaId)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
