package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateDeploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateDeployCreate,
		Read:   resourceMSOSchemaTemplateDeployRead,
		Update: resourceMSOSchemaTemplateDeployCreate,
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

			"force_apply": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "re-deploy",
			},

			"undeploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaTemplateDeployCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Template Deploy", d.Id())
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
	log.Printf("[DEBUG] %s: Template deployed successfully", d.Id())
	return resourceMSOSchemaTemplateDeployRead(d, m)
}

func resourceMSOSchemaTemplateDeployRead(d *schema.ResourceData, m interface{}) error {
	// We set this intentionally blank so that we execute this in every run.
	d.Set("force_apply", "")
	return nil
}

func resourceMSOSchemaTemplateDeployDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Template Undeploy", d.Id())
	msoClient := m.(*client.Client)
	templateName := d.Get("template_name").(string)
	schemaID := d.Get("schema_id").(string)
	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	siteCount, err := schemaCont.ArrayCount("sites")
	if err != nil {
		return err
	}

	for i := 0; i < siteCount; i++ {
		siteCont, err := schemaCont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		currentSiteId := models.StripQuotes(siteCont.S("siteId").String())
		currentTemplateName := models.StripQuotes(siteCont.S("templateName").String())

		if currentTemplateName == templateName {
			log.Printf("[DEBUG] %s: Undeploying site: %s for Template: %s", d.Id(), currentSiteId, currentTemplateName)
			queryString := fmt.Sprintf("?undeploy=%s", currentSiteId)
			path := fmt.Sprintf("/api/v1/execute/schema/%s/template/%s%s", schemaID, templateName, queryString)
			_, err := msoClient.GetViaURL(path)
			if err != nil {
				return err
			}
		}
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Template undeployed successfully", d.Id())
	return nil
}
