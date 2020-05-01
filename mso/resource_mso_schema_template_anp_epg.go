package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaTemplateAnpEpg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgCreate,
		Read:   resourceMSOSchemaTemplateAnpEpgRead,
		Update: resourceMSOSchemaTemplateAnpEpgUpdate,
		Delete: resourceMSOSchemaTemplateAnpEpgDelete,

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
			"anp_name": &schema.Schema{
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
			"epg_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"epg_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"epg_anp_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"bd_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bd_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"useg_epg": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"intra_epg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"intersite_multicaste_source": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaTemplateAnpEpgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Anp Epg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epg_schema_id, epg_template_name, epg_anp_name string
	var uSegEpg, intersiteMulticasteSource, preferredGroup bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicaste_source, ok := d.GetOk("intersite_multicaste_source"); ok {
		intersiteMulticasteSource = intersite_multicaste_source.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}

	if tempVar, ok := d.GetOk("epg_schema_id"); ok {
		epg_schema_id = tempVar.(string)
	} else {
		epg_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("epg_template_name"); ok {
		epg_template_name = tempVar.(string)
	} else {
		epg_template_name = templateName
	}
	if tempVar, ok := d.GetOk("epg_anp_name"); ok {
		epg_anp_name = tempVar.(string)
	} else {
		epg_anp_name = anpName
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	epgRefMap := make(map[string]interface{})
	epgRefMap["schemaId"] = epg_schema_id
	epgRefMap["templateName"] = epg_template_name
	epgRefMap["anpName"] = epg_anp_name
	epgRefMap["epgName"] = Name

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/-", templateName, anpName)
	anpEpgStruct := models.NewTemplateAnpEpg("add", path, Name, displayName, intraEpg, uSegEpg, intersiteMulticasteSource, preferredGroup, vrfRefMap, bdRefMap, epgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaTemplateAnpEpgRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateANP := d.Get("anp_name")
	stateEPG := d.Get("name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							d.SetId(apiEPG)
							d.Set("schema_id", schemaId)
							d.Set("name", apiEPG)
							d.Set("template_name", apiTemplate)
							d.Set("display_name", models.StripQuotes(epgCont.S("displayName").String()))
							d.Set("intra_epg", models.StripQuotes(epgCont.S("intraEpg").String()))
							d.Set("useg_epg", epgCont.S("uSegEpg").Data().(bool))
							d.Set("intersite_multicaste_source", epgCont.S("proxyArp").Data().(bool))
							d.Set("preferred_group", epgCont.S("preferredGroup").Data().(bool))

							vrfRef := models.StripQuotes(epgCont.S("vrfRef").String())
							re_vrf := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
							match_vrf := re_vrf.FindStringSubmatch(vrfRef)
							d.Set("vrf_name", match_vrf[3])
							d.Set("vrf_schema_id", match_vrf[1])
							d.Set("vrf_template_name", match_vrf[2])

							bdRef := models.StripQuotes(epgCont.S("bdRef").String())
							re_bd := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match_bd := re_bd.FindStringSubmatch(bdRef)
							d.Set("bd_name", match_bd[3])
							d.Set("bd_schema_id", match_bd[1])
							d.Set("bd_template_name", match_bd[2])

							epgRef := models.StripQuotes(epgCont.S("epgRef").String())
							re_epg := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)/epgs/(.*)")
							match_epg := re_epg.FindStringSubmatch(epgRef)
							d.Set("epg_name", match_epg[4])
							d.Set("epg_schema_id", match_epg[1])
							d.Set("epg_template_name", match_epg[2])
							d.Set("epg_anp_name", match_epg[3])

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
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaTemplateAnpEpgUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epg_schema_id, epg_template_name, epg_anp_name string
	var uSegEpg, intersiteMulticasteSource, preferredGroup bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicaste_source, ok := d.GetOk("intersite_multicaste_source"); ok {
		intersiteMulticasteSource = intersite_multicaste_source.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}

	if tempVar, ok := d.GetOk("epg_schema_id"); ok {
		epg_schema_id = tempVar.(string)
	} else {
		epg_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("epg_template_name"); ok {
		epg_template_name = tempVar.(string)
	} else {
		epg_template_name = templateName
	}
	if tempVar, ok := d.GetOk("epg_anp_name"); ok {
		epg_anp_name = tempVar.(string)
	} else {
		epg_anp_name = anpName
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	epgRefMap := make(map[string]interface{})
	epgRefMap["schemaId"] = epg_schema_id
	epgRefMap["templateName"] = epg_template_name
	epgRefMap["epgName"] = Name
	epgRefMap["anpName"] = epg_anp_name

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s", templateName, anpName, d.Id())
	anpEpgStruct := models.NewTemplateAnpEpg("replace", path, Name, displayName, intraEpg, uSegEpg, intersiteMulticasteSource, preferredGroup, vrfRefMap, bdRefMap, epgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaTemplateAnpEpgRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epg_schema_id, epg_template_name, epg_anp_name string
	var uSegEpg, intersiteMulticasteSource, preferredGroup bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicaste_source, ok := d.GetOk("intersite_multicaste_source"); ok {
		intersiteMulticasteSource = intersite_multicaste_source.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}

	if tempVar, ok := d.GetOk("epg_schema_id"); ok {
		epg_schema_id = tempVar.(string)
	} else {
		epg_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("epg_template_name"); ok {
		epg_template_name = tempVar.(string)
	} else {
		epg_template_name = templateName
	}
	if tempVar, ok := d.GetOk("epg_anp_name"); ok {
		epg_anp_name = tempVar.(string)
	} else {
		epg_anp_name = anpName
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	epgRefMap := make(map[string]interface{})
	epgRefMap["schemaId"] = epg_schema_id
	epgRefMap["templateName"] = epg_template_name
	epgRefMap["anpName"] = epg_anp_name
	epgRefMap["epgName"] = Name

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s", templateName, anpName, d.Id())
	anpEpgStruct := models.NewTemplateAnpEpg("remove", path, Name, displayName, intraEpg, uSegEpg, intersiteMulticasteSource, preferredGroup, vrfRefMap, bdRefMap, epgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
