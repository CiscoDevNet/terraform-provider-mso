package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgStaticPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgStaticPortCreate,
		Read:   resourceMSOSchemaSiteAnpEpgStaticPortRead,
		Update: resourceMSOSchemaSiteAnpEpgStaticPortUpdate,
		Delete: resourceMSOSchemaSiteAnpEpgStaticPortDelete,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
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
			"path_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"pod": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"leaf": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"deployment_immediacy": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"micro_segvlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgStaticPortCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	statesiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	var pathType, pod, leaf, path, deploymentImmediacy, mode string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_segvlan"); ok {
		microsegvlan = tempVar.(int)
	}

	portpath := fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/-", statesiteId, stateTemplateName, stateANPName, stateEpgName)
	staticStruct := models.NewSchemaSiteAnpEpgStaticPort("add", pathsp, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
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
	statepod := d.Get("pod").(string)
	stateleaf := d.Get("leaf").(string)
	statepath := d.Get("path").(string)
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
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								portpath := fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", statepod, stateleaf, statepath)
								apiportpath := models.StripQuotes(portCont.S("path").String())
								if portpath == apiportpath {
									d.SetId(apiportpath)
									if portCont.Exists("type") {
										d.Set("type", models.StripQuotes(portCont.S("type").String()))
									}
									if portCont.Exists("path") {
										d.Set("pod", statepod)
										d.Set("leaf", stateleaf)
										d.Set("path", statepath)
									}
									if portCont.Exists("portEncapVlan") {
										tempvar, _ := strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
										d.Set("vlan", tempvar)
									}
									if portCont.Exists("deploymentImmediacy") {
										d.Set("deployment_immediacy", models.StripQuotes(portCont.S("deploymentImmediacy").String()))
									}
									if portCont.Exists("microSegVlan") {
										tempvar1, _ := strconv.Atoi(fmt.Sprintf("%v", portCont.S("microSegVlan")))
										d.Set("micro_segvlan", tempvar1)
									}
									if portCont.Exists("mode") {
										d.Set("mode", models.StripQuotes(portCont.S("mode").String()))
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

func resourceMSOSchemaSiteAnpEpgStaticPortUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	found := false

	var pathType, pod, leaf, path, deploymentImmediacy, mode string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_segvlan"); ok {
		microsegvlan = tempVar.(int)
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
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								portpath := fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
								apiportpath := models.StripQuotes(portCont.S("path").String())
								if portpath == apiportpath {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/%v", stateSite, stateTemplate, stateAnp, stateEpg, index)
									anpStruct := models.NewSchemaSiteAnpEpgStaticPort("replace", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
									_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

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
		return fmt.Errorf("The specified parameters to update static port entry not found")
	}

	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticPortDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	found := false

	var pathType, pod, leaf, path, deploymentImmediacy, mode string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_segvlan"); ok {
		microsegvlan = tempVar.(int)
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
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								portpath := fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
								apiportpath := models.StripQuotes(portCont.S("path").String())
								if portpath == apiportpath {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/%v", stateSite, stateTemplate, stateAnp, stateEpg, index)
									anpStruct := models.NewSchemaSiteAnpEpgStaticPort("remove", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
									_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

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
		return fmt.Errorf("The specified parameters to delete the static port entry not found")
	}
	d.SetId("")
	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}
