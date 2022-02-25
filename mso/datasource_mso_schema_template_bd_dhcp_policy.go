package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOTemplateBDDHCPPolicy() *schema.Resource {
	return &schema.Resource{
		Read: datasourceMSOTemplateBDDHCPPolicyRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"bd_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dhcp_option_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dhcp_option_version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	bdDHCPPolicy := models.TemplateBDDHCPPolicy{
		Name:         d.Get("name").(string),
		SchemaID:     d.Get("schema_id").(string),
		TemplateName: d.Get("template_name").(string),
		BDName:       d.Get("bd_name").(string),
	}

	remoteBDDHCPPolicy, err := getMSOTemplateBDDHCPPolicy(msoClient, &bdDHCPPolicy)
	if err != nil {
		return err
	}
	setMSOTemplateBDDHCPPolicy(d, remoteBDDHCPPolicy)

	d.SetId(createMSOTemplateBDDHCPPolicyId(&bdDHCPPolicy))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
