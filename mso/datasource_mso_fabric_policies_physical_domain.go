package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOPhysicalDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOPhysicalDomainRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_pool_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOPhysicalDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Domain Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	name := d.Get("name").(string)

	setPhysicalDomainData(d, msoClient, templateId, name)
	log.Printf("[DEBUG] MSO Physical Domain Data Source - Read Complete : %v", d.Id())
	return nil
}
