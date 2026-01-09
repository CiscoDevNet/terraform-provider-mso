package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSODHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSODHCPOptionPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DHCP Option Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the DHCP Option Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DHCP Option Policy.",
			},
			"options": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of DHCP options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the DHCP option.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the DHCP option.",
						},
						"data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data value of the DHCP option.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMSODHCPOptionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Option Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "dhcpOptionPolicies")
	if err != nil {
		return err
	}

	setDHCPOptionPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO DHCP Option Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
