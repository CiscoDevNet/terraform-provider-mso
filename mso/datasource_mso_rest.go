package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSORest() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSORestRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func datasourceMSORestRead(d *schema.ResourceData, m interface{}) error {
	path := d.Get("path").(string)
	log.Printf("[DEBUG] %s: Beginning Read", path)

	msoClient := m.(*client.Client)
	content, err := MakeRestRequest(msoClient, path, "GET", "{}")
	if err != nil {
		return err
	}
	d.SetId(path)
	d.Set("content", content.String())

	log.Printf("[DEBUG] %s: Read finished successfully", path)
	return nil
}
