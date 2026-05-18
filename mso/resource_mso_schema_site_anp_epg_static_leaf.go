package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// resourceMSOSchemaSiteAnpEpgStaticleaf manages an
// mso_schema_site_anp_epg_static_leaf entry (a static leaf binding attached
// to a site ANP EPG's staticLeafs list on the schema).
//
// Create-error behaviour (intentionally not changed): if the user applies
// this resource against a template that has no mso_schema_site association,
// older NDO returns "Resource Not Found" and newer NDO silently drops the
// PATCH (the follow-up Read finds nothing and the Terraform SDK reports
// "Provider produced inconsistent result after apply"). Either surfaces the
// misconfiguration to the user.
//
// Delete implementation note: the static leaf entry is stored as an object
// in the sites[].anps[].epgs[].staticLeafs[] array and is removed by its
// array index (path: /sites/{siteId}-{template}/anps/{anp}/epgs/{epg}/staticLeafs/{index}).
// The index is resolved at delete time by scanning the array for the matching
// path. If the entry is not found the delete is treated as a no-op.
//
// All schema attributes are ForceNew because the resource has no Update
// function. The NDO PATCH API could support in-place updates of
// port_encap_vlan (addressed by array index), but without an Update handler
// Terraform would silently ignore any changes. ForceNew ensures a
// destroy+recreate instead.
//
// Future improvement: this resource manages one static leaf per instance,
// requiring N separate schema GETs and PATCHes for N static leafs. A
// potential replacement is a static_leaf TypeSet block on mso_schema_site_anp_epg,
// which already owns the parent EPG and will gain further site-level EPG
// attributes. That would reduce N static leaf changes to a single PATCH on
// the EPG resource. This resource would then be deprecated in favour of
// that block.
func resourceMSOSchemaSiteAnpEpgStaticleaf() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgStaticleafCreate,
		Read:   resourceMSOSchemaSiteAnpEpgStaticleafRead,
		Delete: resourceMSOSchemaSiteAnpEpgStaticleafDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgStaticleafImport,
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
			"path": &schema.Schema{
				// ForceNew: the topology path (e.g. "topology/pod-1/node-101")
				// identifies the static leaf binding. Changing it requires
				// removing the old entry and adding a new one, which is
				// equivalent to destroy+recreate.
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"port_encap_vlan": &schema.Schema{
				// ForceNew: this resource has no Update function. Without
				// ForceNew, a port_encap_vlan change would be silently ignored
				// (the plan would show a diff but no PATCH would be sent).
				// ForceNew ensures Terraform destroys and recreates the entry
				// instead. An Update function could be added in the future to
				// allow in-place changes.
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgStaticleafImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	import_attribute := regexp.MustCompile("(.*)/path/(.*)")
	import_split := import_attribute.FindStringSubmatch(d.Id())
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}

	stateSite := get_attribute[2]
	stateTemplate := get_attribute[4]
	found := false
	stateAnp := get_attribute[6]
	stateEpg := get_attribute[8]
	statePath := import_split[2]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
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
							staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
							if err != nil {
								return nil, fmt.Errorf("Unable to get Static Leaf list")
							}
							for s := 0; s < staticLeafCount; s++ {
								staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
								if err != nil {
									return nil, err
								}
								apiPath := models.StripQuotes(staticLeafCont.S("path").String())
								if apiPath == statePath {
									d.SetId(apiPath)
									d.Set("path", apiPath)
									d.Set("site_id", apiSite)
									d.Set("schema_id", schemaId)
									d.Set("template_name", stateTemplate)
									d.Set("anp_name", stateAnp)
									d.Set("epg_name", apiEPG)
									apiPort, _ := strconv.Atoi(staticLeafCont.S("portEncapVlan").String())
									d.Set("port_encap_vlan", apiPort)
									found = true
									break
								}
							}
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
		d.Set("path", "")
		return nil, fmt.Errorf("Unable to find the given Anp Epg StaticLeaf")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteAnpEpgStaticleafCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg StaticLeaf: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	paths := d.Get("path").(string)
	portEncapVlan := d.Get("port_encap_vlan").(int)

	foundAnp := false
	foundEpg := false
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == siteId && apiTemplate == templateName {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}

				anpRef := models.StripQuotes(anpCont.S("anpRef").String())

				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)

				if match[3] == anpName {

					foundAnp = true
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

						if apiEPG == epgName {
							foundEpg = true
							break
						}
					}

					if !foundEpg {
						log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
						anpEpgRefMap := make(map[string]interface{})
						anpEpgRefMap["schemaId"] = schemaId
						anpEpgRefMap["templateName"] = apiTemplate
						anpEpgRefMap["anpName"] = anpName
						anpEpgRefMap["epgName"] = epgName

						pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, anpName)
						//private_link_label argument used in resource site_anp_epg is set to nil here
						anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

						_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
						if ers != nil {
							return ers
						}
						break

					}
				}
			}

			if !foundAnp {
				log.Printf("[DEBUG] Site Anp: Beginning Creation")

				anpRefMap := make(map[string]interface{})
				anpRefMap["schemaId"] = schemaId
				anpRefMap["templateName"] = apiTemplate
				anpRefMap["anpName"] = anpName

				pathAnp := fmt.Sprintf("/sites/%s-%s/anps/-", apiSite, apiTemplate)
				anpStruct := models.NewSchemaSiteAnp("add", pathAnp, anpRefMap)

				_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
				anpEpgRefMap := make(map[string]interface{})
				anpEpgRefMap["schemaId"] = schemaId
				anpEpgRefMap["templateName"] = apiTemplate
				anpEpgRefMap["anpName"] = anpName
				anpEpgRefMap["epgName"] = epgName

				pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, anpName)
				//private_link_label argument used in resource site_anp_epg is set to nil here
				anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

				_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
				if ers != nil {
					return ers
				}

			}

		}
	}

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticLeafs/-", siteId, templateName, anpName, epgName)
	anpEpgStaticStruct := models.NewSchemaSiteAnpEpgStaticleaf("add", path, paths, portEncapVlan)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStaticStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOSchemaSiteAnpEpgStaticleafRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticleafRead(d *schema.ResourceData, m interface{}) error {
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
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	statePath := d.Get("path").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
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
							staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
							if err != nil {
								return fmt.Errorf("Unable to get Static Leaf list")
							}
							for s := 0; s < staticLeafCount; s++ {
								staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
								if err != nil {
									return err
								}
								apiPath := models.StripQuotes(staticLeafCont.S("path").String())
								if apiPath == statePath {
									d.SetId(apiPath)
									d.Set("path", apiPath)
									d.Set("site_id", apiSite)
									d.Set("schema_id", schemaId)
									d.Set("template_name", stateTemplate)
									d.Set("anp_name", stateAnp)
									d.Set("epg_name", apiEPG)
									apiPort, _ := strconv.Atoi(staticLeafCont.S("portEncapVlan").String())
									d.Set("port_encap_vlan", apiPort)
									found = true
									break
								}
							}
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
		d.Set("path", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgStaticleafDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg StaticLeaf: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	paths := d.Get("path").(string)
	portEncapVlan := d.Get("port_encap_vlan").(int)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	index := -1
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	statePath := d.Get("path").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
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
							staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
							if err != nil {
								return fmt.Errorf("Unable to get Static Leaf list")
							}
							for s := 0; s < staticLeafCount; s++ {
								staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
								if err != nil {
									return err
								}
								apiPath := models.StripQuotes(staticLeafCont.S("path").String())
								if apiPath == statePath {
									index = s
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if index == -1 {
		d.SetId("")
		return nil
	}

	indexs := strconv.Itoa(index)
	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticLeafs/%s", siteId, templateName, anpName, epgName, indexs)
	anpEpgStaticStruct := models.NewSchemaSiteAnpEpgStaticleaf("remove", path, paths, portEncapVlan)
	response, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStaticStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err1 != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err1
	}
	d.SetId("")
	return nil
}
