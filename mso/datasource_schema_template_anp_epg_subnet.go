package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceMSOSchemaTemplateAnpEpgSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaTemplateAnpEpgSubnetRead,

		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template": &schema.Schema{
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

			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"scope": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
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
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
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
						log.Println("currentepgname", currentEpgName)
						if currentEpgName == epgName {
							log.Println("found correct epgname")
							d.Set("epg_name", currentEpgName)
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
								log.Println("currentip", currentIp)
								if currentIp == ip {
									log.Println("found correct ip")
									d.SetId(currentIp)
									d.Set("ip", currentIp)
									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										shared, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("shared").String()))
										d.Set("shared", shared)
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
		return fmt.Errorf("The parameters inserted are not valid")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
