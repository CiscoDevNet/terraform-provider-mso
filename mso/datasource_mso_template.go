package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOTemplate() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOTemplateRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tenant",
					"l3out",
					"fabric_policy",
					"fabric_resource",
					"monitoring_tenant",
					"monitoring_access",
					"service_device",
				}, false),
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sites": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func datasourceMSOTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] MSO Template Datasource: Beginning Read")

	msoClient := m.(*client.Client)
	id := d.Get("template_id").(string)
	name := d.Get("template_name").(string)
	templateType := d.Get("template_type").(string)

	if id == "" && name == "" {
		return fmt.Errorf("either `template_id` or `template_name` must be provided")
	} else if id != "" && name != "" {
		return fmt.Errorf("only one of `template_id` or `template_name` must be provided")
	} else if name != "" && templateType == "" {
		return fmt.Errorf("`template_type` must be provided when `template_name` is provided")
	}

	ndoTemplate := ndoTemplate{msoClient: msoClient, id: id, templateName: name, templateType: templateType}
	err := ndoTemplate.getTemplate(true)
	if err != nil {
		return err
	}
	ndoTemplate.SetToSchema(d)
	d.Set("template_id", d.Id())
	log.Println("[DEBUG] MSO Template Datasource: Read Completed", d.Id())
	return nil

}
