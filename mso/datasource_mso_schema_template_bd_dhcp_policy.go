package mso

import (
	"fmt"

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
				Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"dhcp_option_policy_name": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	if err := readTemplateBDDHCPPolicy(d, m); err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("DHCP policy %q not found on Bridge Domain %q", d.Get("name"), d.Get("bd_name"))
	}
	return nil
}
