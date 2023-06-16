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

func resourceMSOTemplateBDSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateBDSubnetCreate,
		Read:   resourceMSOTemplateBDSubnetRead,
		Update: resourceMSOTemplateBDSubnetUpdate,
		Delete: resourceMSOTemplateBDSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateBDSubnetImport,
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

func resourceMSOTemplateBDSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No TemplateSubnet found")
	}
	stateTemplateSubnet := get_attribute[2]
	found := false
	stateBD := get_attribute[4]
	stateIP := import_split[2]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplateSubnet := models.StripQuotes(tempCont.S("name").String())

		if apiTemplateSubnet == stateTemplateSubnet {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
				}

				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return nil, fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						dataCon, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return nil, fmt.Errorf("Unable to parse the subnets list")
						}

						apiIP := models.StripQuotes(dataCon.S("ip").String())
						if apiIP == stateIP {
							log.Println(dataCon)
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplateSubnet)
							d.Set("bd_name", apiBD)
							ip := models.StripQuotes(dataCon.S("ip").String())
							idSubnet := strings.Split(ip, "/")
							d.SetId(idSubnet[0])
							d.Set("ip", models.StripQuotes(dataCon.S("ip").String()))
							d.Set("scope", models.StripQuotes(dataCon.S("scope").String()))
							d.Set("description", models.StripQuotes(dataCon.S("description").String()))
							d.Set("shared", dataCon.S("shared").Data().(bool))
							if dataCon.Exists("noDefaultGateway") {
								d.Set("no_default_gateway", dataCon.S("noDefaultGateway").Data().(bool))
							}
							if dataCon.Exists("querier") {
								d.Set("querier", dataCon.S("querier").Data().(bool))
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
		return nil, fmt.Errorf("Unable to find the BD Subnet with IP: %s", stateIP)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOTemplateBDSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

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

	path := fmt.Sprintf("/templates/%s/bds/%s/subnets/-", templateName, bdName)
	bdSubnetStruct := models.NewTemplateBDSubnet("add", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdSubnetStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateBDSubnetRead(d, m)
}

func resourceMSOTemplateBDSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {

			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {

					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						subnetsCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subntes list")
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplate)
							d.Set("bd_name", apiBD)
							ip := models.StripQuotes(subnetsCont.S("ip").String())
							idSubnet := strings.Split(ip, "/")
							d.SetId(idSubnet[0])
							d.Set("ip", models.StripQuotes(subnetsCont.S("ip").String()))
							d.Set("scope", models.StripQuotes(subnetsCont.S("scope").String()))
							d.Set("description", models.StripQuotes(subnetsCont.S("description").String()))
							d.Set("shared", subnetsCont.S("shared").Data().(bool))
							if subnetsCont.Exists("noDefaultGateway") {
								d.Set("no_default_gateway", subnetsCont.S("noDefaultGateway").Data().(bool))
							}
							if subnetsCont.Exists("querier") {
								d.Set("querier", subnetsCont.S("querier").Data().(bool))
							}
							found = true
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

func resourceMSOTemplateBDSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

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
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {

			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {

					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						subnetsCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subntes list")
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							index := k
							path := fmt.Sprintf("/templates/%s/bds/%s/subnets/%v", apiTemplate, apiBD, index)
							bdSubnetStruct := models.NewTemplateBDSubnet("replace", path, apiIP, Desc, Scope, Shared, NoDefaultGateway, Querier)
							_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdSubnetStruct)
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
	if !found {
		return fmt.Errorf("The specified parameters not found for update operation")
	}
	return resourceMSOTemplateBDSubnetRead(d, m)
}

func resourceMSOTemplateBDSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

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
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {

			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {

					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						subnetsCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subntes list")
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							index := k
							path := fmt.Sprintf("/templates/%s/bds/%s/subnets/%v", apiTemplate, apiBD, index)
							bdSubnetStruct := models.NewTemplateBDSubnet("remove", path, apiIP, Desc, Scope, Shared, NoDefaultGateway, Querier)
							response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdSubnetStruct)

							// Ignoring Error with code 141: Resource Not Found when deleting
							if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
								return err
							}
							break
						}
					}
				}
			}
		}
	}
	d.SetId("")
	return nil
}
