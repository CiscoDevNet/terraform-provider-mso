package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSODHCPOptionPolicyOption() *schema.Resource {
	return &schema.Resource{
		Read:          datasourceMSODHCPOptionPolicyOptionRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"option_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"option_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"option_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"option_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func datasourceMSODHCPOptionPolicyOptionRead(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	dhcpOptionPolicyOption := models.DHCPOptionPolicyOption{
		Name:       d.Get("option_name").(string),
		PolicyName: d.Get("option_policy_name").(string),
	}
	id := fmt.Sprintf("%s/%s", dhcpOptionPolicyOption.PolicyName, dhcpOptionPolicyOption.Name)
	if data, ok := d.GetOk("option_data"); ok {
		dhcpOptionPolicyOption.Data = data.(string)
	}
	if optionId, ok := d.GetOk("option_id"); ok {
		dhcpOptionPolicyOption.ID = optionId.(string)
	}
	remoteDHCPOptionPolicyOption, err := msoClient.ReadDHCPOptionPolicyOption(id)
	if err != nil {
		return err
	}
	setDHCPOptionPolicyOption(d, remoteDHCPOptionPolicyOption)
	d.SetId(DHCPOptionPolicyOptionModeltoId(remoteDHCPOptionPolicyOption))
	return nil
}
