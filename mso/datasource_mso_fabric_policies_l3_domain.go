package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	domain, err := GetPolicyByName(response, domainName, "fabricPolicyTemplate", "template", "l3Domains")
	if err != nil {
		return err
	}

	setL3DomainData(d, domain, templateId)
	log.Printf("[DEBUG] MSO L3 Domain Data Source - Read Complete: %v", d.Id())
	return nil
}
