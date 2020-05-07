package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOTemplateExtenalepgSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateExtenalepgSubnetCreate,
		Read:   resourceMSOTemplateExtenalepgSubnetRead,
		Update: resourceMSOTemplateExtenalepgSubnetUpdate,
		Delete: resourceMSOTemplateExtenalepgSubnetDelete,

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
			"externalepg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"aggregate": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func resourceMSOTemplateExtenalepgSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg Subnet: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	extenalepgName := d.Get("externalepg_name").(string)
	templateName := d.Get("template_name").(string)

	var IP, Name string
	var Scope, Aggregate []interface{}

	if tempVar, ok := d.GetOk("ip"); ok {
		IP = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("name"); ok {
		Name = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("scope"); ok {
		Scope = tempVar.([]interface{})
	}
	if tempVar, ok := d.GetOk("aggregate"); ok {
		Aggregate = tempVar.([]interface{})
	}

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/subnets/-", templateName, extenalepgName)
	externalepgStruct := models.NewTemplateExternalEpgSubnet("add", path, IP, Name, Scope, Aggregate)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateExtenalepgSubnetRead(d, m)
}

func resourceMSOTemplateExtenalepgSubnetRead(d *schema.ResourceData, m interface{}) error {
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
	stateExternalepg := d.Get("externalepg_name")
	stateIP := d.Get("ip")

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			externalepgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalepgCount; j++ {
				externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
				if apiExternalepg == stateExternalepg {
					subnetCount, err := externalepgCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get subnets list")
					}
					for k := 0; k < subnetCount; k++ {
						subnetsCont, err := externalepgCont.ArrayElement(k, "subnets")
						if err != nil {
							return err
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplate)
							d.Set("externalepg_name", apiExternalepg)
							ip := models.StripQuotes(subnetsCont.S("ip").String())
							idSubnet := strings.Split(ip, "/")
							d.SetId(idSubnet[0])
							d.Set("ip", models.StripQuotes(subnetsCont.S("ip").String()))
							d.Set("name", models.StripQuotes(subnetsCont.S("name").String()))
							d.Set("scope", subnetsCont.S("scope").Data().([]interface{}))
							d.Set("aggregate", subnetsCont.S("aggregate").Data().([]interface{}))

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
		d.Set("ip", "")
		d.Set("scope", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateExtenalepgSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg Subnet: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	extenalepgName := d.Get("externalepg_name").(string)
	templateName := d.Get("template_name").(string)

	var IP, Name string
	var Scope, Aggregate []interface{}

	if tempVar, ok := d.GetOk("ip"); ok {
		IP = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("name"); ok {
		Name = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("scope"); ok {
		Scope = tempVar.([]interface{})
	}
	if tempVar, ok := d.GetOk("aggregate"); ok {
		Aggregate = tempVar.([]interface{})
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := getIndex(cont, templateName, extenalepgName, IP)
	if err != nil {
		return err
	}
	if index == -1 {
		fmt.Errorf("The given subnet ip is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/subnets/%s", templateName, extenalepgName, indexs)
	externalepgStruct := models.NewTemplateExternalEpgSubnet("replace", path, IP, Name, Scope, Aggregate)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)

	if errs != nil {
		return errs
	}
	return resourceMSOTemplateExtenalepgSubnetRead(d, m)
}

func resourceMSOTemplateExtenalepgSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg Subnet: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	extenalepgName := d.Get("externalepg_name").(string)
	templateName := d.Get("template_name").(string)

	var IP, Name string
	var Scope, Aggregate []interface{}

	if tempVar, ok := d.GetOk("ip"); ok {
		IP = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("name"); ok {
		Name = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("scope"); ok {
		Scope = tempVar.([]interface{})
	}
	if tempVar, ok := d.GetOk("aggregate"); ok {
		Aggregate = tempVar.([]interface{})
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := getIndex(cont, templateName, extenalepgName, IP)
	if err != nil {
		return err
	}
	if index == -1 {
		fmt.Errorf("The given subnet ip is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/subnets/%s", templateName, extenalepgName, indexs)
	externalepgStruct := models.NewTemplateExternalEpgSubnet("remove", path, IP, Name, Scope, Aggregate)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)

	if errs != nil {
		return errs
	}
	d.SetId("")
	return nil
}

func getIndex(cont *container.Container, templateName, extenalepgName, ip string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, fmt.Errorf("No Template found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == templateName {
			externalepgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return index, fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalepgCount; j++ {
				externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return index, err
				}
				apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
				if apiExternalepg == extenalepgName {
					subnetCount, err := externalepgCont.ArrayCount("subnets")
					if err != nil {
						return index, fmt.Errorf("Unable to get subnets list")
					}
					for k := 0; k < subnetCount; k++ {
						subnetsCont, err := externalepgCont.ArrayElement(k, "subnets")
						if err != nil {
							return index, err
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == ip {
							log.Println("Correct IP")
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
