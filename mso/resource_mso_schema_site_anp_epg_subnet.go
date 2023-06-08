package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgSubnetCreate,
		Read:   resourceMSOSchemaSiteAnpEpgSubnetRead,
		Update: resourceMSOSchemaSiteAnpEpgSubnetUpdate,
		Delete: resourceMSOSchemaSiteAnpEpgSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgSubnetImport,
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
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	import_attribute := regexp.MustCompile("(.*)/ip/(.*)")
	import_split := import_attribute.FindStringSubmatch(d.Id())
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}
	stateSite := get_attribute[2]
	found := false
	stateTemplate := get_attribute[4]
	stateAnp := get_attribute[6]
	stateEpg := get_attribute[8]
	stateIp := import_split[2]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			d.Set("site_id", apiSite)
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
					d.Set("anp_name", match[3])
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
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
						match := re.FindStringSubmatch(apiEpgRef)
						apiEPG := match[3]
						if apiEPG == stateEpg {
							d.Set("epg_name", apiEPG)
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return nil, fmt.Errorf("Unable to get Subnet list")
							}
							for l := 0; l < subnetCount; l++ {
								subnetCont, err := epgCont.ArrayElement(l, "subnets")
								if err != nil {
									return nil, err
								}
								apiIP := models.StripQuotes(subnetCont.S("ip").String())
								if stateIp == apiIP {
									d.SetId(apiIP)
									if subnetCont.Exists("ip") {
										d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
									}
									if subnetCont.Exists("description") {
										d.Set("description", models.StripQuotes(subnetCont.S("description").String()))
									}
									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										d.Set("shared", subnetCont.S("shared").Data().(bool))
									}
									if subnetCont.Exists("noDefaultGateway") {
										d.Set("no_default_gateway", subnetCont.S("noDefaultGateway").Data().(bool))
									}
									if subnetCont.Exists("querier") {
										d.Set("querier", subnetCont.S("querier").Data().(bool))
									}
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
		return nil, fmt.Errorf("The subnet entry with specified ip %s not found", stateIp)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteAnpEpgSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site ANP EPG Subnet: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateSiteId := d.Get("site_id").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	var IP string
	if ip, ok := d.GetOk("ip"); ok {
		IP = ip.(string)
	}
	var Scope string
	if scope, ok := d.GetOk("scope"); ok {
		Scope = scope.(string)
	}
	var Shared bool
	if shared, ok := d.GetOk("shared"); ok {
		Shared = shared.(bool)
	}
	var NoDefaultGateway bool
	if ndg, ok := d.GetOk("no_default_gateway"); ok {
		NoDefaultGateway = ndg.(bool)
	}
	var Querier bool
	if qr, ok := d.GetOk("querier"); ok {
		Querier = qr.(bool)
	}
	var Desc string
	if d, ok := d.GetOk("description"); ok {
		Desc = d.(string)
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	foundEpg := false
	foundAnp := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSiteId && apiTemplate == stateTemplateName {
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
				if match[3] == stateANPName {
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
						if apiEPG == stateEpgName {
							foundEpg = true
							break

						}
					}
					if !foundEpg {
						log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
						anpEpgRefMap := make(map[string]interface{})
						anpEpgRefMap["schemaId"] = schemaId
						anpEpgRefMap["templateName"] = stateTemplateName
						anpEpgRefMap["anpName"] = stateANPName
						anpEpgRefMap["epgName"] = stateEpgName

						pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, stateANPName)
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
				anpRefMap["templateName"] = stateTemplateName
				anpRefMap["anpName"] = stateANPName

				pathAnp := fmt.Sprintf("/sites/%s-%s/anps/-", stateSiteId, stateTemplateName)
				anpStruct := models.NewSchemaSiteAnp("add", pathAnp, anpRefMap)

				_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")

				anpEpgRefMap := make(map[string]interface{})
				anpEpgRefMap["schemaId"] = schemaId
				anpEpgRefMap["templateName"] = stateTemplateName
				anpEpgRefMap["anpName"] = stateANPName
				anpEpgRefMap["epgName"] = stateEpgName

				pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", stateSiteId, stateTemplateName, stateANPName)
				//private_link_label argument used in resource site_anp_epg is set to nil here
				anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

				_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
				if ers != nil {
					return ers
				}
			}
		}
	}

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/subnets/-", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	AnpEpgSubnetStruct := models.NewSchemaSiteAnpEpgSubnet("add", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), AnpEpgSubnetStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOSchemaSiteAnpEpgSubnetRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgSubnetRead(d *schema.ResourceData, m interface{}) error {
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
	stateIp := d.Get("ip").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			d.Set("site_id", apiSite)
			d.Set("template_name", apiTemplate)
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
				if match[3] == stateAnp {
					d.Set("anp_name", match[3])
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
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
						match := re.FindStringSubmatch(apiEpgRef)
						apiEPG := match[3]
						if apiEPG == stateEpg {
							d.Set("epg_name", apiEPG)
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("Unable to get Subnet list")
							}
							for l := 0; l < subnetCount; l++ {
								subnetCont, err := epgCont.ArrayElement(l, "subnets")
								if err != nil {
									return err
								}
								apiIP := models.StripQuotes(subnetCont.S("ip").String())
								if stateIp == apiIP {
									d.SetId(apiIP)
									if subnetCont.Exists("ip") {
										d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
									}
									if subnetCont.Exists("description") {
										d.Set("description", models.StripQuotes(subnetCont.S("description").String()))
									}
									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										d.Set("shared", subnetCont.S("shared").Data().(bool))
									}
									if subnetCont.Exists("noDefaultGateway") {
										d.Set("no_default_gateway", subnetCont.S("noDefaultGateway").Data().(bool))
									}
									if subnetCont.Exists("querier") {
										d.Set("querier", subnetCont.S("querier").Data().(bool))
									}
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
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site ANP EPG Subnet: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	statesiteId := d.Get("site_id").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)
	found := false

	var IP string
	if ip, ok := d.GetOk("ip"); ok {
		IP = ip.(string)
	}

	var Scope string
	if scope, ok := d.GetOk("scope"); ok {
		Scope = scope.(string)
	}

	var Shared bool
	if shared, ok := d.GetOk("shared"); ok {
		Shared = shared.(bool)
	}

	var NoDefaultGateway bool
	if ndg, ok := d.GetOk("no_default_gateway"); ok {
		NoDefaultGateway = ndg.(bool)
	}

	var Querier bool
	if qr, ok := d.GetOk("querier"); ok {
		Querier = qr.(bool)
	}
	var Desc string
	if d, ok := d.GetOk("description"); ok {
		Desc = d.(string)
	}
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")

	if err != nil {
		return fmt.Errorf("No Site found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplateName := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == statesiteId && apiTemplateName == stateTemplateName {

			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)

				apiANPName := match[3]

				if apiANPName == stateANPName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return err
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err1 := anpCont.ArrayElement(k, "epgs")
						if err1 != nil {
							return err1
						}
						epgRef := models.StripQuotes(epgCont.S("epgRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
						match := re.FindStringSubmatch(epgRef)
						apiEpgName := match[3]

						if apiEpgName == stateEpgName {
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("Unable to get Subnet list")
							}
							for l := 0; l < subnetCount; l++ {
								subnetCont, err := epgCont.ArrayElement(l, "subnets")
								if err != nil {
									return err
								}
								apiIP := models.StripQuotes(subnetCont.S("ip").String())
								if IP == apiIP {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/subnets/%v", statesiteId, stateTemplateName, stateANPName, stateEpgName, index)
									AnpEpgSubnetStruct := models.NewSchemaSiteAnpEpgSubnet("replace", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
									_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), AnpEpgSubnetStruct)
									if err != nil {
										return err
									}
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
		return fmt.Errorf("Subnet entry cannot be updated as specified parameters not found")
	}

	return resourceMSOSchemaSiteAnpEpgSubnetRead(d, m)

}

func resourceMSOSchemaSiteAnpEpgSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site ANP EPG Subnet: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)

	var IP string
	if ip, ok := d.GetOk("ip"); ok {
		IP = ip.(string)
	}

	var Scope string
	if scope, ok := d.GetOk("scope"); ok {
		Scope = scope.(string)
	}

	var Shared bool
	if shared, ok := d.GetOk("shared"); ok {
		Shared = shared.(bool)
	}

	var NoDefaultGateway bool
	if ndg, ok := d.GetOk("no_default_gateway"); ok {
		NoDefaultGateway = ndg.(bool)
	}

	var Querier bool
	if qr, ok := d.GetOk("querier"); ok {
		Querier = qr.(bool)
	}
	var Desc string
	if d, ok := d.GetOk("description"); ok {
		Desc = d.(string)
	}

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
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
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
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
						match := re.FindStringSubmatch(apiEpgRef)
						apiEPG := match[3]
						if apiEPG == stateEpg {
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("Unable to get Subnet list")
							}
							for l := 0; l < subnetCount; l++ {
								subnetCont, err := epgCont.ArrayElement(l, "subnets")
								if err != nil {
									return err
								}
								apiIP := models.StripQuotes(subnetCont.S("ip").String())
								if IP == apiIP {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/subnets/%v", stateSite, stateTemplate, stateAnp, stateEpg, index)
									AnpEpgSubnetStruct := models.NewSchemaSiteAnpEpgSubnet("remove", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
									response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), AnpEpgSubnetStruct)

									// Ignoring Error with code 141: Resource Not Found when deleting
									if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
										return err
									}
								}
							}
						}
					}

				}
			}
		}
	}

	d.SetId("")
	return resourceMSOSchemaSiteAnpEpgSubnetRead(d, m)
}
