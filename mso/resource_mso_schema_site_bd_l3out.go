package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceMSOSchemaSiteBdL3out manages an mso_schema_site_bd_l3out entry
// (an L3out reference attached to a site BD's l3Outs list on the schema).
//
// Create-error behaviour (intentionally not changed): if the user applies
// this resource against a template that has no mso_schema_site_bd parent,
// older NDO returns "Resource Not Found" and newer NDO silently drops the
// PATCH (the follow-up Read finds nothing and the Terraform SDK reports
// "Provider produced inconsistent result after apply"). Either surfaces the
// misconfiguration to the user.
//
// Delete implementation note: the l3out entry is stored as a plain string
// in the sites[].bds[].l3Outs[] array and is removed by its array index
// (path: /sites/{siteId}-{template}/bds/{bd_name}/l3Outs/{index}). The
// index is resolved at delete time by scanning the array for the matching
// name. If the entry is not found the delete is treated as a no-op.
//
// Future improvement: this resource manages one L3out per instance, meaning
// N L3outs require N separate schema GETs and PATCHes (plus the NDO
// validation engine running on each). A potential replacement is an l3out
// TypeSet block on mso_schema_site_bd, which already owns the parent BD and
// its site-level attributes (host_route, svi_mac). That would reduce N L3out
// changes to a single PATCH on the BD resource, and this resource would then
// be deprecated in favour of that block.
func resourceMSOSchemaSiteBdL3out() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteBdL3outCreate,
		Read:   resourceMSOSchemaSiteBdL3outRead,
		Delete: resourceMSOSchemaSiteBdL3outDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteBdL3outImport,
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
			"bd_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew is required: l3out entries are stored under the
				// site BD at /sites/{siteId}-{template}/bds/{bd_name}/l3Outs/.
				// Moving an l3out to a different BD would require removing it
				// from one BD's l3Outs array and adding it to another, which
				// is equivalent to destroy+recreate.
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew is required: l3out entries are plain strings in the
				// l3Outs array and are identified by their array index for
				// removal (see the Delete path: /l3Outs/{index}). There is no
				// API operation to rename an existing entry in-place; a rename
				// must go through destroy+recreate.
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaSiteBdL3outImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
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
	found := false
	stateBd := get_attribute[4]
	stateL3out := get_attribute[6]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
				}
				apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
				split := strings.Split(apiBdRef, "/")
				apiBd := split[6]
				if apiBd == stateBd {
					d.Set("site_id", apiSite)
					d.Set("schema_id", split[2])
					d.Set("template_name", split[4])
					d.Set("bd_name", split[6])
					l3outCount, err := bdCont.ArrayCount("l3Outs")
					if err != nil {
						return nil, fmt.Errorf("Unable to get l3Outs list")
					}
					for k := 0; k < l3outCount; k++ {
						l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
						if err != nil {
							return nil, err
						}
						tempVar := l3outCont.String()
						apiL3out := strings.Trim(tempVar, "\"")
						if apiL3out == stateL3out {
							d.SetId(stateL3out)
							d.Set("l3out_name", apiL3out)
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find the Site L3out %s", stateL3out)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteBdL3outCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd L3out: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)
	l3outName := d.Get("l3out_name").(string)

	path := fmt.Sprintf("/sites/%s-%s/bds/%s/l3Outs/-", siteId, templateName, bdName)
	BdL3outStruct := models.NewSchemaSiteBdL3out("add", path, l3outName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdL3outStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteBdL3outRead(d, m)
}

func resourceMSOSchemaSiteBdL3outRead(d *schema.ResourceData, m interface{}) error {
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
	stateBd := d.Get("bd_name").(string)
	stateL3out := d.Get("l3out_name").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
				split := strings.Split(apiBdRef, "/")
				apiBd := split[6]
				if apiBd == stateBd {
					d.Set("site_id", apiSite)
					d.Set("schema_id", split[2])
					d.Set("template_name", split[4])
					d.Set("bd_name", split[6])
					l3outCount, err := bdCont.ArrayCount("l3Outs")
					if err != nil {
						return fmt.Errorf("Unable to get l3Outs list")
					}
					for k := 0; k < l3outCount; k++ {
						l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
						if err != nil {
							return err
						}
						tempVar := l3outCont.String()
						apiL3out := strings.Trim(tempVar, "\"")
						if apiL3out == stateL3out {
							d.SetId(stateL3out)
							d.Set("l3out_name", apiL3out)
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteBdL3outDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd L3out: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)
	l3outName := d.Get("l3out_name").(string)

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := countIndex(cont, siteId, bdName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/bds/%s/l3Outs/%s", siteId, templateName, bdName, indexs)
	BdL3outStruct := models.NewSchemaSiteBdL3out("remove", path, l3outName)

	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdL3outStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return errs
	}

	d.SetId("")
	return nil
}

func countIndex(cont *container.Container, stateSite, stateBd, stateL3out string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return index, fmt.Errorf("No Sites found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return index, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return index, fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return index, err
				}
				apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
				split := strings.Split(apiBdRef, "/")
				apiBd := split[6]
				if apiBd == stateBd {
					l3outCount, err := bdCont.ArrayCount("l3Outs")
					if err != nil {
						return index, fmt.Errorf("Unable to get l3Outs list")
					}
					for k := 0; k < l3outCount; k++ {
						l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
						if err != nil {
							return index, err
						}
						tempVar := l3outCont.String()
						apiL3out := strings.Trim(tempVar, "\"")
						if apiL3out == stateL3out {
							log.Println("found correct L3out")
							index = k
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}

	return index, nil

}
