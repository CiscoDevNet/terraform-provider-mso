package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOMCPGlobalPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOMCPGlobalPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_mcp_pdu_per_vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"loop_detect_multiplication_factor": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_disable_protection": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_delay_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transmission_frequency_sec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transmission_frequency_msec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOMCPGlobalPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MCP Global Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "mcpGlobalPolicy")
	if err != nil {
		return err
	}

	setMCPGlobalPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO MCP Global Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
