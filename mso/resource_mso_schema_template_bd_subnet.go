package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
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
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public",
					"private",
				}, false),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Set minimal length to 0 to allow removal of description
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"virtual": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

func setBDSubnetResourceData(d *schema.ResourceData, schemaId, templateName, bdName string, subnetCont *container.Container) {
	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("bd_name", bdName)
	ip := models.StripQuotes(subnetCont.S("ip").String())
	idSubnet := strings.Split(ip, "/")
	d.SetId(idSubnet[0])
	d.Set("ip", ip)
	d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
	if subnetCont.Exists("description") {
		d.Set("description", models.StripQuotes(subnetCont.S("description").String()))
	} else {
		d.Set("description", "")
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
	if subnetCont.Exists("primary") {
		d.Set("primary", subnetCont.S("primary").Data().(bool))
	}
	if subnetCont.Exists("virtual") {
		d.Set("virtual", subnetCont.S("virtual").Data().(bool))
	}
}

func resourceMSOTemplateBDSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	getAttribute := strings.Split(d.Id(), "/")
	importAttribute := regexp.MustCompile("(.*)/ip/(.*)")
	importSplit := importAttribute.FindStringSubmatch(d.Id())
	schemaId := getAttribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := getAttribute[2]
	stateBD := getAttribute[4]
	stateIP := importSplit[2]
	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount && !found; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return nil, fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < subnetCount && !found; k++ {
						subnetCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return nil, fmt.Errorf("Unable to parse the subnets list")
						}
						apiIP := models.StripQuotes(subnetCont.S("ip").String())
						if apiIP == stateIP {
							setBDSubnetResourceData(d, schemaId, apiTemplate, apiBD, subnetCont)
							found = true
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
	log.Printf("[DEBUG] Template BD Subnet: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

	IP := d.Get("ip").(string)
	Scope := d.Get("scope").(string)
	Shared := d.Get("shared").(bool)
	NoDefaultGateway := d.Get("no_default_gateway").(bool)
	Querier := d.Get("querier").(bool)
	Desc := d.Get("description").(string)
	Primary := d.Get("primary").(bool)
	Virtual := d.Get("virtual").(bool)

	path := fmt.Sprintf("/templates/%s/bds/%s/subnets/-", templateName, bdName)
	bdSubnetStruct := models.NewTemplateBDSubnet("add", path, IP, Desc, Scope, Shared, NoDefaultGateway, Querier, Primary, Virtual)

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
	stateBD := d.Get("bd_name").(string)
	stateIP := d.Get("ip").(string)
	for i := 0; i < count && !found; i++ {
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
			for j := 0; j < bdCount && !found; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {

					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < subnetCount && !found; k++ {
						subnetCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subnets list")
						}
						apiIP := models.StripQuotes(subnetCont.S("ip").String())
						if apiIP == stateIP {
							setBDSubnetResourceData(d, schemaId, apiTemplate, apiBD, subnetCont)
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
	log.Printf("[DEBUG] Template BD Subnet: Beginning Update")
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
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")
	for i := 0; i < count && !found; i++ {
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
			for j := 0; j < bdCount && !found; j++ {
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
					for k := 0; k < count1 && !found; k++ {
						subnetsCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subnets list")
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							updatePath := fmt.Sprintf("/templates/%s/bds/%s/subnets/%d", apiTemplate, apiBD, k)
							payloadCont := container.New()
							payloadCont.Array()

							if d.HasChange("scope") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/scope", updatePath), d.Get("scope").(string))
								if err != nil {
									return err
								}
							}

							if d.HasChange("description") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
								if err != nil {
									return err
								}
							}

							if d.HasChange("shared") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/shared", updatePath), d.Get("shared").(bool))
								if err != nil {
									return err
								}
							}

							if d.HasChange("no_default_gateway") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/noDefaultGateway", updatePath), d.Get("no_default_gateway").(bool))
								if err != nil {
									return err
								}
							}

							if d.HasChange("querier") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/querier", updatePath), d.Get("querier").(bool))
								if err != nil {
									return err
								}
							}

							if d.HasChange("primary") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/primary", updatePath), d.Get("primary").(bool))
								if err != nil {
									return err
								}
							}

							if d.HasChange("virtual") {
								err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/virtual", updatePath), d.Get("virtual").(bool))
								if err != nil {
									return err
								}
							}

							err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaId), payloadCont)
							if err != nil {
								return err
							}
							found = true
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
	log.Printf("[DEBUG] Template BD Subnet: Beginning Delete")
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
	stateBD := d.Get("bd_name").(string)
	stateIP := d.Get("ip").(string)
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
							return fmt.Errorf("Unable to parse the subnets list")
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							path := fmt.Sprintf("/templates/%s/bds/%s/subnets/%v", apiTemplate, apiBD, k)
							response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), models.GetRemovePatchPayload(path))

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
