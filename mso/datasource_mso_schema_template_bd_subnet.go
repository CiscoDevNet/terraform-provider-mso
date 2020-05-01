package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceMSOTemplateSubnetBD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateSubnetBDRead,

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

			"bd_name": &schema.Schema{
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateSubnetBDRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No TemplateSubnet found")
	}
	stateTemplateSubnet := d.Get("template_name").(string)
	found := false
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplateSubnet := models.StripQuotes(tempCont.S("name").String())

		if apiTemplateSubnet == stateTemplateSubnet {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}

				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						dataCon, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subnets list")
						}

						apiIP := models.StripQuotes(dataCon.S("ip").String())
						if apiIP == stateIP {
							log.Println(dataCon)
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplateSubnet)
							d.Set("bd_name", apiBD)
							ip := models.StripQuotes(dataCon.S("ip").String())
							idSubnet := strings.Split(ip, "/")
							d.SetId(idSubnet[0])
							d.Set("ip", models.StripQuotes(dataCon.S("ip").String()))
							d.Set("scope", models.StripQuotes(dataCon.S("scope").String()))
							d.Set("description", models.StripQuotes(dataCon.S("description").String()))
							d.Set("shared", dataCon.S("shared").Data().(bool))
							d.Set("no_default_gateway", dataCon.S("noDefaultGateway").Data().(bool))
							d.Set("querier", dataCon.S("querier").Data().(bool))
							found = true
							break
						}

					}
				}

			}
		}

	}

	if !found {
		return fmt.Errorf("Unable to find the BD Subnet with IP: %s", stateIP)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
