package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOL3Domain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOL3DomainRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the fabric policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the L3 Domain.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3 Domain.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the L3 Domain.",
			},
			"vlan_pool_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the VLAN Pool associated with this L3 Domain.",
			},
		},
	}
}

func dataSourceMSOL3DomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3 Domain Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	domainName := d.Get("name").(string)

	domain, err := getL3Domain(msoClient, templateId, domainName)
	if err != nil {
		return err
	}

	if err := setL3DomainData(d, domain, templateId); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO L3 Domain Data Source - Read Complete: %v", d.Id())
	return nil
}
