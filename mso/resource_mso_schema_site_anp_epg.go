package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// resourceMSOSchemaSiteAnpEpg manages an mso_schema_site_anp_epg association
// (a template ANP EPG surfaced under a specific site on the schema).
//
// Unlike mso_schema_site_vrf and mso_schema_site_anp (which are pure
// deprecation candidates), this resource will gain a new "shutdown" attribute
// that controls the EPG admin state at the site level. Until that attribute
// is added, the resource provides no incremental value on newer NDO releases
// where the schema validation engine is always-on: once a template ANP EPG
// exists and the template is associated with a site, NDO automatically
// materializes the corresponding sites[].anps[].epgs[] entry on every schema
// update.
//
// private_link_label: Deprecation candidate – this is a cloud-only service
// EPG attribute that is not exercised by the current acceptance tests.
// It should be investigated further and potentially moved to a dedicated
// cloud-site resource.
//
// Create-error behaviour (intentionally not changed): if the user applies
// this resource against a template that has no mso_schema_site association,
// older NDO returns "Resource Not Found" and newer NDO silently drops the
// PATCH (the follow-up Read finds nothing and the Terraform SDK reports
// "Provider produced inconsistent result after apply"). Either surfaces the
// misconfiguration to the user; we deliberately leave the current
// replace-then-add Create flow alone to preserve backward compatibility on
// pre-validator NDO releases that still rely on the inject-via-PATCH path.
func resourceMSOSchemaSiteAnpEpg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgCreate,
		Update: resourceMSOSchemaSiteAnpEpgUpdate,
		Read:   resourceMSOSchemaSiteAnpEpgRead,
		Delete: resourceMSOSchemaSiteAnpEpgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgImport,
		},

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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"private_link_label": &schema.Schema{
				// cloud feature should be investigated further
				// see https://www.cisco.com/c/en/us/td/docs/dcn/mso/use-case/configuring-service-epgs-in-cisco-aci-multi-site-orchestrator.html
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", get_attribute[0]))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}

	stateSite := get_attribute[2]
	found := false
	stateAnp := get_attribute[6]
	stateEpg := get_attribute[8]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
				split := strings.Split(apiAnpRef, "/")
				apiAnp := split[6]
				if apiAnp == stateAnp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}
						apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
						split := strings.Split(apiEpgRef, "/")
						apiEPG := split[8]
						if apiEPG == stateEpg {
							d.SetId(apiEPG)
							d.Set("site_id", apiSite)
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)
							privatelinklabelsCont := epgCont.S("privateLinkLabel")
							if models.StripQuotes(privatelinklabelsCont.S("name").String()) == "{}" {
								d.Set("private_link_label", "")
							} else {
								d.Set("private_link_label", models.StripQuotes(privatelinklabelsCont.S("name").String()))
							}
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find the Site Anp Epg %s", stateEpg)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteAnpEpgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	anpEpgRefMap := make(map[string]interface{})
	anpEpgRefMap["schemaId"] = schemaId
	anpEpgRefMap["templateName"] = templateName
	anpEpgRefMap["anpName"] = anpName
	anpEpgRefMap["epgName"] = epgName

	privateLinkLabel := make(map[string]interface{})
	if val, ok := d.GetOk("private_link_label"); ok {
		map_private_link_label := make(map[string]interface{})
		map_private_link_label["name"] = val
		privateLinkLabel = map_private_link_label
	} else {
		privateLinkLabel = nil
	}

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	if versionInt != 1 {
		path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s", siteId, templateName, anpName, epgName)
		anpEpgStruct := models.NewSchemaSiteAnpEpg("replace", path, privateLinkLabel, anpEpgRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
	}

	if versionInt == 1 || err != nil {
		path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", siteId, templateName, anpName)
		anpEpgStruct := models.NewSchemaSiteAnpEpg("add", path, privateLinkLabel, anpEpgRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
	}

	if err != nil {
		return err
	}

	return resourceMSOSchemaSiteAnpEpgRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	stateSite := d.Get("site_id").(string)
	found := false
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
				split := strings.Split(apiAnpRef, "/")
				apiAnp := split[6]
				if apiAnp == stateAnp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
						split := strings.Split(apiEpgRef, "/")
						apiEPG := split[8]
						if apiEPG == stateEpg {
							d.SetId(apiEPG)
							d.Set("site_id", apiSite)
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
		d.Set("epg_name", "")
		d.Set("anp_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	updatePath := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s", siteId, templateName, anpName, epgName)
	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("private_link_label") {
		// when a epg is not of type service the privateLinkLabel can still be set in the site while being a service type attribute
		// this could cause a change trigger wiping site configuration even for application type epg when attribute is changed
		// unsupported in NaC https://github.com/netascode/terraform-mso-nac-ndo/blob/main/ndo_schemas.tf#L1024
		// changing the PATCH to avoid removal regardless of epg type
		oldPrivateLinkLabel, newPrivateLinkLabel := d.GetChange("private_link_label")
		operation := "replace"
		if oldPrivateLinkLabel == "" {
			operation = "add"
		}
		privateLinkLabelMap := make(map[string]interface{})
		if newPrivateLinkLabel != "" {
			privateLinkLabelMap["name"] = newPrivateLinkLabel
		}
		err := addPatchPayloadToContainer(payloadCont, operation, fmt.Sprintf("%s/privateLinkLabel", updatePath), privateLinkLabelMap)
		if err != nil {
			return err
		}
	}

	err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaId), payloadCont)
	if err != nil {
		return err
	}

	return resourceMSOSchemaSiteAnpEpgRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	anpEpgRefMap := make(map[string]interface{})
	anpEpgRefMap["schemaId"] = schemaId
	anpEpgRefMap["templateName"] = templateName
	anpEpgRefMap["anpName"] = anpName
	anpEpgRefMap["epgName"] = epgName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s", siteId, templateName, anpName, epgName)
	privateLinkLabel := make(map[string]interface{})
	if val, ok := d.GetOk("private_link_label"); ok {
		map_private_link_label := make(map[string]interface{})
		map_private_link_label["name"] = val
		privateLinkLabel = map_private_link_label
	} else {
		privateLinkLabel = nil
	}
	anpEpgStruct := models.NewSchemaSiteAnpEpg("remove", path, privateLinkLabel, anpEpgRefMap)
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	return nil
}
