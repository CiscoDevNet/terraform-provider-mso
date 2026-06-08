package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO: Reconsider the shape of this data source.
//
// Concerns with the current design:
//   - The only non-input output is `site_id`. The standalone `mso_site` data
//     source already resolves site name -> site_id without requiring a schema
//     or template, so this data source is largely redundant for that purpose.
//   - The schema/template -> sites relationship is 1:N, but this data source
//     hides that: it is keyed by a single site `name` and silently returns the
//     first match from `GET api/v1/sites`. It exposes no list of associated
//     sites and no template-side metadata (e.g. undeploy state, site type).
//   - The implicit existence check has no representation in state: a present
//     association yields `site_id`, an absent one yields a hard error. There
//     is no boolean attribute to branch on, and no way to ask "is this site
//     attached?" without aborting the plan. With multiple sites attached to a
//     template, callers must declare one data block per site they care about
//     and any absent one fails the whole run -- exactly the case where a
//     list-shaped data source (option 2 below) would be more ergonomic.

func datasourceMSOSchemaSite() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSchemaSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/sites"))
	if err != nil {
		return err
	}

	var flag bool
	var siteId string
	for _, info := range con.S("sites").Data().([]interface{}) {
		val := info.(map[string]interface{})
		if val["name"].(string) == name {
			flag = true
			siteId = val["id"].(string)
			break
		}
	}
	if flag != true {
		return fmt.Errorf("Site of specified name not found")
	}
	_, err = getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s-%s", schemaId, siteId, templateName))
	d.Set("schema_id", schemaId)
	d.Set("site_id", siteId)
	d.Set("template_name", templateName)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
