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

func resourceMSOSchemaSiteBdSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteBdSubnetCreate,
		Read:   resourceMSOSchemaSiteBdSubnetRead,
		Update: resourceMSOSchemaSiteBdSubnetUpdate,
		Delete: resourceMSOSchemaSiteBdSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteBdSubnetImport,
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
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceMSOSchemaSiteBdSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	stateBd := get_attribute[4]
	stateIp := import_split[2]
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			d.Set("site_id", apiSite)
			d.Set("template_name", models.StripQuotes(tempCont.S("templateName").String()))
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == stateBd {
					d.Set("bd_name", match[3])
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return nil, fmt.Errorf("Unable to get Subnet list")
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := bdCont.ArrayElement(l, "subnets")
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
	if !found {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the Site BD Subnet")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteBdSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site Bd Subnet: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	statesiteId := d.Get("site_id").(string)
	stateBd := d.Get("bd_name").(string)

	var IP string
	if ip, ok := d.GetOk("ip"); ok {
		IP = ip.(string)
	}

	Scope := "private"
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

	path := fmt.Sprintf("/sites/%s-%s/bds/%s/subnets/-", statesiteId, stateTemplateName, stateBd)
	BdSubnetStruct := models.NewSchemaSiteBdSubnet("add", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdSubnetStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteBdSubnetRead(d, m)
}

func resourceMSOSchemaSiteBdSubnetRead(d *schema.ResourceData, m interface{}) error {
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
	stateBd := d.Get("bd_name").(string)
	stateIp := d.Get("ip").(string)
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			d.Set("site_id", apiSite)
			d.Set("template_name", apiTemplate)
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == stateBd {
					d.Set("bd_name", match[3])
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet list")
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := bdCont.ArrayElement(l, "subnets")
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
	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteBdSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site Bd Subnet: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	statesiteId := d.Get("site_id").(string)
	stateBd := d.Get("bd_name").(string)

	var IP string
	if ip, ok := d.GetOk("ip"); ok {
		IP = ip.(string)
	}

	Scope := "private"
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

	index := -1
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")

	if err != nil {
		return fmt.Errorf("No Site found")
	}
	updated := false

	for i := 0; i < count && !updated; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplateName := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == statesiteId && apiTemplateName == stateTemplateName {

			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}

			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)

				apiBdName := match[3]

				if apiBdName == stateBd {
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet list")
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := bdCont.ArrayElement(l, "subnets")
						if err != nil {
							return err
						}
						apiIP := models.StripQuotes(subnetCont.S("ip").String())
						if IP == apiIP {
							index = l
							path := fmt.Sprintf("/sites/%s-%s/bds/%s/subnets/%v", statesiteId, stateTemplateName, stateBd, index)
							BdSubnetStruct := models.NewSchemaSiteBdSubnet("replace", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
							_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdSubnetStruct)
							if err != nil {
								return err
							}
							updated = true
						}

					}
				}
			}
		}
	}
	if index == -1 {
		return fmt.Errorf("Unable to find the given subnet IP")
	}

	return resourceMSOSchemaSiteBdSubnetRead(d, m)

}

func resourceMSOSchemaSiteBdSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site Bd Subnet: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateBd := d.Get("bd_name").(string)
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

	deleted := false

	for i := 0; i < count && !deleted; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == stateBd {
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet list")
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := bdCont.ArrayElement(l, "subnets")
						if err != nil {
							return err
						}
						apiIP := models.StripQuotes(subnetCont.S("ip").String())
						if IP == apiIP {
							index := l
							path := fmt.Sprintf("/sites/%s-%s/bds/%s/subnets/%v", stateSite, stateTemplate, stateBd, index)
							BdSubnetStruct := models.NewSchemaSiteBdSubnet("remove", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)
							response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdSubnetStruct)

							// Ignoring Error with code 141: Resource Not Found when deleting
							if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
								return err
							}
							deleted = true
						}
					}
				}
			}

		}
	}
	d.SetId("")
	return resourceMSOSchemaSiteBdSubnetRead(d, m)
}
