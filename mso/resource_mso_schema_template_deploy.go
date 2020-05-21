package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateDeploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateDeployCreate,
		Read:   resourceMSOSchemaTemplateDeployRead,
		Delete: resourceMSOSchemaTemplateDeployDelete,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"undeploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaTemplateDeployCreate(d *schema.ResourceData, m interface{}) error {
	templateName := d.Get("template_name").(string)
	schemaID := d.Get("schema_id").(string)

	var siteId, queryString string
	var undeploy bool

	if tempVar, ok := d.GetOk("undeploy"); ok {
		undeploy = tempVar.(bool)
		if undeploy {
			if siteTemp, ok := d.GetOk("site_id"); ok {
				siteId = siteTemp.(string)
				queryString = fmt.Sprintf("?undeploy=%s", siteId)
			} else {
				return fmt.Errorf("SiteID must be provided with undeploy = true")
			}
		} else {
			queryString = ""
		}
	}
	path := fmt.Sprintf("/api/v1/execute/schema/%s/template/%s%s", schemaID, templateName, queryString)
	msoClient := m.(*client.Client)
	_, err := msoClient.GetViaURL(path)
	if err != nil {
		return err
	}
	d.SetId(schemaID)
	return resourceMSOSchemaTemplateDeployRead(d, m)
}

func resourceMSOSchemaTemplateDeployRead(d *schema.ResourceData, m interface{}) error {
	// We set this intentionally blank so that we execute this in every run.
	d.Set("template_name", "")
	return nil
}

func resourceMSOSchemaTemplateDeployDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
