package mso

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMSOSystemConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSystemConfigRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"banner": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"change_control": {
				Type:     schema.TypeMap,
				Computed: true,
				// SDKv2 does not support Elem with schema.Resource on TypeMap fields.
				// Expected keys: "workflow" (string: "enabled"/"disabled"), "number_of_approvers" (integer ≥ 1). Validation skipped - resource is deprecated.
			},
		}),
		DeprecationMessage: nd4DeprecationMessage("mso_system_config"),
	}
}

func dataSourceMSOSystemConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	err := getAndSetSystemConfig(d, m)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
