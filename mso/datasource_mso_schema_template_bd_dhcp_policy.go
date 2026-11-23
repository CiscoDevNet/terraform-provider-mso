package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMSOTemplateBDDHCPPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateBDDHCPPolicyRead,
		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": {
				Type: schema.TypeString, Optional: true, Computed: true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				ExactlyOneOf: []string{"uuid"},
			},
			"uuid": {
				Type: schema.TypeString, Optional: true, Computed: true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				ExactlyOneOf: []string{"name"},
			},
			"dhcp_option_policy_name": {Type: schema.TypeString, Computed: true},
			"dhcp_option_policy_uuid": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	policy, err := getTemplateBDDHCPPolicy(d, m.(*client.Client))
	if err != nil {
		return err
	}
	if policy == nil {
		identifier := getTemplateBDDHCPPolicyInput(d, "name")
		if identifier == "" {
			identifier = getTemplateBDDHCPPolicyInput(d, "uuid")
		}
		return fmt.Errorf("DHCP policy %q not found on Bridge Domain %q", identifier, d.Get("bd_name"))
	}
	for key, value := range policy {
		if err := d.Set(key, value); err != nil {
			return err
		}
	}
	d.SetId(getTemplateBDDHCPPolicyID(d))
	return nil
}
