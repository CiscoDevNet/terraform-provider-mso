package mso

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workflow": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"number_of_approvers": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		}),
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
