package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgUsegAttrCreate,
		Update: resourceMSOSchemaTemplateAnpEpgUsegAttrUpdate,
		Read:   resourceMSOSchemaTemplateAnpEpgUsegAttrRead,
		Delete: resourceMSOSchemaTemplateAnpEpgUsegAttrDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpEpgUsegAttrImport,
		},

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

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"useg_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"mac",
					"dns",
					"vm-name",      // Vm Name
					"rootContName", // VM data center
					"hv",           // Hypervisor
					"guest-os",     // Operating System
					"tag",
					"vm",     // Identifier
					"domain", // VMM domain
					"vnic",   // Vnic DN
				}, false),
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Computed removed to allow setting description to empty string
				// Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},

			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Computed removed to allow setting operator to empty string
				// Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"equals",
					"startsWith",
					"endsWith",
					"contains",
				}, false),
			},

			"category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Computed removed to allow setting category to empty string
				// Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"value": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				// Note: NDO normalizes the `value` for several `useg_type`s by uppercasing
				// the stored value (observed for `vm-name`, `dns`, `hv`, `guest-os`, `vnic`).
				// Supplying a non-normalized value will result in a perpetual plan diff.
				// Provide the value in the form NDO will store (uppercase) to avoid drift.
			},

			"useg_subnet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	templateName := get_attribute[2]
	anpName := get_attribute[4]
	epgName := get_attribute[6]
	name := get_attribute[8]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	usegCont, err := findUsegAttrContainer(cont, templateName, anpName, epgName, name)
	if err != nil {
		return nil, err
	}
	if usegCont == nil {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find Schema template anp epg useg attribute %s", name)
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("anp_name", anpName)
	d.Set("epg_name", epgName)
	setUsegAttrAttributes(d, usegCont)
	d.SetId(name)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg UsegAttr: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	usegAttrMap := getUsegAttrPayload(d)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/-", templateName, anpName, epgName)
	usegAttrApp := models.NewSchemaTemplateAnpEpgUsegAttr("add", path, usegAttrMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), usegAttrApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", name))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg UsegAttr: Beginning Updating")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	usegAttrMap := getUsegAttrPayload(d)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s", templateName, anpName, epgName, name)
	usegAttrApp := models.NewSchemaTemplateAnpEpgUsegAttr("replace", path, usegAttrMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), usegAttrApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", name))
	log.Printf("[DEBUG] %s: Updating finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	usegCont, err := findUsegAttrContainer(cont, templateName, anpName, epgName, name)
	if err != nil {
		return err
	}
	if usegCont == nil {
		d.SetId("")
		log.Printf("[DEBUG] Read finished successfully (object not found)")
		return nil
	}

	setUsegAttrAttributes(d, usegCont)
	d.SetId(name)
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s", templateName, anpName, epgName, name)
	removePayload := models.GetRemovePatchPayload(path)
	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), removePayload)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return errs
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}

// getUsegAttrPayload builds the uSegAttr patch payload from resource data,
// honoring NDO (4.x+) semantics: only attributes relevant to the configured
// useg_type are included so create/replace payloads match what the server stores.
func getUsegAttrPayload(d *schema.ResourceData) map[string]interface{} {
	name := d.Get("name").(string)
	usegType := d.Get("useg_type").(string)
	val := d.Get("value").(string)

	operator, operatorSet := "", false
	if v, ok := d.GetOk("operator"); ok {
		operator = v.(string)
		operatorSet = true
	}

	category := ""
	if v, ok := d.GetOk("category"); ok {
		category = v.(string)
	}

	desc := ""
	if v, ok := d.GetOk("description"); ok {
		desc = v.(string)
	}

	usegSubnet := false
	if v, ok := d.GetOk("useg_subnet"); ok {
		usegSubnet = v.(bool)
	}

	payload := map[string]interface{}{
		"name":        name,
		"displayName": name,
		"type":        usegType,
		"value":       val,
	}

	// operator is forwarded whenever the user explicitly sets it. NDO 4.x discards
	// operator for ip/mac/dns (treated as implicit "equals"); users on 4.x should
	// omit operator from config for those types to avoid drift.
	if operatorSet {
		payload["operator"] = operator
	}
	if category != "" {
		payload["category"] = category
	}
	if desc != "" {
		payload["description"] = desc
	}
	// fvSubnet is only applicable to type "ip" and only when true.
	if usegType == "ip" && usegSubnet {
		payload["fvSubnet"] = true
	}

	return payload
}

// findUsegAttrContainer walks the schema container looking for a uSegAttr entry
// matching the given template/anp/epg/name. Returns nil (without error) when
// the entry is not found so callers can decide how to handle a missing object.
func findUsegAttrContainer(cont *container.Container, templateName, anpName, epgName, name string) (*container.Container, error) {
	templateCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	for i := 0; i < templateCount; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		if models.StripQuotes(tempCont.S("name").String()) != templateName {
			continue
		}

		anpCount, err := tempCont.ArrayCount("anps")
		if err != nil {
			return nil, fmt.Errorf("No Anp found")
		}
		for j := 0; j < anpCount; j++ {
			anpCont, err := tempCont.ArrayElement(j, "anps")
			if err != nil {
				return nil, err
			}
			if models.StripQuotes(anpCont.S("name").String()) != anpName {
				continue
			}

			epgCount, err := anpCont.ArrayCount("epgs")
			if err != nil {
				return nil, fmt.Errorf("No Epg found")
			}
			for k := 0; k < epgCount; k++ {
				epgCont, err := anpCont.ArrayElement(k, "epgs")
				if err != nil {
					return nil, err
				}
				if models.StripQuotes(epgCont.S("name").String()) != epgName {
					continue
				}

				usegCount, err := epgCont.ArrayCount("uSegAttrs")
				if err != nil {
					// No uSegAttrs array on this EPG -> entry not found.
					return nil, nil
				}
				for s := 0; s < usegCount; s++ {
					usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
					if err != nil {
						return nil, err
					}
					if models.StripQuotes(usegCont.S("name").String()) == name {
						return usegCont, nil
					}
				}
			}
		}
	}

	return nil, nil
}

// setUsegAttrAttributes populates resource data from a uSegAttr container.
// NDO (4.x+) omits attributes that aren't applicable to the useg_type
// (e.g. operator for ip/mac/dns, fvSubnet when false). Guard with Exists so
// state isn't polluted with the literal "{}" returned for missing keys.
func setUsegAttrAttributes(d *schema.ResourceData, usegCont *container.Container) {
	d.Set("name", models.StripQuotes(usegCont.S("name").String()))
	d.Set("useg_type", models.StripQuotes(usegCont.S("type").String()))
	d.Set("value", models.StripQuotes(usegCont.S("value").String()))

	if usegCont.Exists("operator") {
		d.Set("operator", models.StripQuotes(usegCont.S("operator").String()))
	} else {
		d.Set("operator", "")
	}

	if usegCont.Exists("category") {
		d.Set("category", models.StripQuotes(usegCont.S("category").String()))
	} else {
		d.Set("category", "")
	}

	if usegCont.Exists("description") {
		d.Set("description", models.StripQuotes(usegCont.S("description").String()))
	} else {
		d.Set("description", "")
	}

	if usegCont.Exists("fvSubnet") {
		usegSubnet, _ := strconv.ParseBool(models.StripQuotes(usegCont.S("fvSubnet").String()))
		d.Set("useg_subnet", usegSubnet)
	} else {
		d.Set("useg_subnet", false)
	}
}
