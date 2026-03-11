package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOVirtualPortChannelInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOFabricResourcePoliciesVirtualPortChannelInterfaceRead,

		SchemaVersion: version,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Fabric Resource template ID.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Virtual Port Channel Interface name.",
			},
			"uuid":        {Type: schema.TypeString, Computed: true},
			"description": {Type: schema.TypeString, Computed: true},

			"node_1": {Type: schema.TypeString, Computed: true},
			"node_2": {Type: schema.TypeString, Computed: true},

			"node_1_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_2_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"interface_policy_group_uuid": {Type: schema.TypeString, Computed: true},

			"interface_descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node":        {Type: schema.TypeString, Computed: true},
						"interface":   {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceMSOFabricResourcePoliciesVirtualPortChannelInterfaceRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VPC Interface Data Source - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	name := d.Get("name").(string)

	err := setVPCInterfaceData(d, msoClient, templateId, name)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO VPC Interface Data Source - Read Complete: %v", d.Id())
	return nil
}
