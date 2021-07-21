package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateAnpEpgSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgSubnetCreate,
		Update: resourceMSOSchemaTemplateAnpEpgSubnetUpdate,
		Read:   resourceMSOSchemaTemplateAnpEpgSubnetRead,
		Delete: resourceMSOSchemaTemplateAnpEpgSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpEpgSubnetImport,
		},

		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template": &schema.Schema{
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
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
		}),
	}
}

func resourceMSOSchemaTemplateAnpEpgSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	import_attribute := regexp.MustCompile("(.*)/ip/(.*)")
	import_split := import_attribute.FindStringSubmatch(d.Id())
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	d.Set("schema_id", schemaId)
	templateName := get_attribute[2]
	anpName := get_attribute[4]
	epgName := get_attribute[6]
	ip := import_split[2]
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return nil, fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return nil, err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
					d.Set("anp_name", currentAnpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						log.Println("currentepgname", currentEpgName)
						if currentEpgName == epgName {
							log.Println("found correct epgname")
							d.Set("epg_name", currentEpgName)
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return nil, fmt.Errorf("No Subnets found")
							}
							for s := 0; s < subnetCount; s++ {
								subnetCont, err := epgCont.ArrayElement(s, "subnets")
								if err != nil {
									return nil, err
								}
								currentIp := models.StripQuotes(subnetCont.S("ip").String())
								log.Println("currentip", currentIp)
								if currentIp == ip {
									log.Println("found correct ip")
									d.SetId(currentIp)
									d.Set("ip", currentIp)
									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										shared, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("shared").String()))
										d.Set("shared", shared)
									}

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
		}
		if found {
			break
		}
	}

	if !found {
		d.SetId("")
		return nil, fmt.Errorf("The ANP EPG Subnet is not found")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateAnpEpgSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg Subnet: Beginning Creation")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var anpName string
	if name, ok := d.GetOk("anp_name"); ok {
		anpName = name.(string)
	}

	var epgName string
	if name, ok := d.GetOk("epg_name"); ok {
		epgName = name.(string)
	}

	var ip string
	if tempVar, ok := d.GetOk("ip"); ok {
		ip = tempVar.(string)
	}

	scope := "private"
	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}

	shared := false
	if tempVar, ok := d.GetOk("shared"); ok {
		shared = tempVar.(bool)
	}

	schemaTemplateAnpEpgSubnetApp := models.NewSchemaTemplateAnpEpgSubnet("add", "/templates/"+templateName+"/anps/"+anpName+"/epgs/"+epgName+"/subnets/-", ip, scope, shared)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpEpgSubnetApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", ip))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgSubnetRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg Subnet: Beginning Updating")
	msoClient := m.(*client.Client)

	ips := d.Id()

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var anpName string
	if name, ok := d.GetOk("anp_name"); ok {
		anpName = name.(string)
	}

	var epgName string
	if name, ok := d.GetOk("epg_name"); ok {
		epgName = name.(string)
	}

	var ip string
	if tempVar, ok := d.GetOk("ip"); ok {
		ip = tempVar.(string)
	}

	scope := "private"
	if tempVar, ok := d.GetOk("scope"); ok {
		scope = tempVar.(string)
	}

	shared := false
	if tempVar, ok := d.GetOk("shared"); ok {
		shared = tempVar.(bool)
	}

	conts, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	index, err := fetchIndex(conts, templateName, anpName, epgName, ips)
	if err != nil {
		return err
	}

	if index == -1 {
		fmt.Errorf("The given subnet ip is not found")
	}

	indexs := strconv.Itoa(index)

	schemaTemplateAnpEpgSubnetApp := models.NewSchemaTemplateAnpEpgSubnet("replace", "/templates/"+templateName+"/anps/"+anpName+"/epgs/"+epgName+"/subnets/"+indexs, ip, scope, shared)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpEpgSubnetApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", ip))
	log.Printf("[DEBUG] %s: Updating finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgSubnetRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgSubnetRead(d *schema.ResourceData, m interface{}) error {
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

	templateName := d.Get("template").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	ip := d.Get("ip").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
					d.Set("anp_name", currentAnpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						log.Println("currentepgname", currentEpgName)
						if currentEpgName == epgName {
							log.Println("found correct epgname")
							d.Set("epg_name", currentEpgName)
							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("No Subnets found")
							}
							for s := 0; s < subnetCount; s++ {
								subnetCont, err := epgCont.ArrayElement(s, "subnets")
								if err != nil {
									return err
								}
								currentIp := models.StripQuotes(subnetCont.S("ip").String())
								log.Println("currentip", currentIp)
								if currentIp == ip {
									log.Println("found correct ip")
									d.SetId(currentIp)
									d.Set("ip", currentIp)
									if subnetCont.Exists("scope") {
										d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
									}
									if subnetCont.Exists("shared") {
										shared, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("shared").String()))
										d.Set("shared", shared)
									}

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

func resourceMSOSchemaTemplateAnpEpgSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	template := d.Get("template").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := fetchIndex(cont, template, anpName, epgName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}
	indexs := strconv.Itoa(index)
	schemaTemplateAnpEpgSubnetApp := models.NewSchemaTemplateAnpEpgSubnet("remove", "/templates/"+template+"/anps/"+anpName+"/epgs/"+epgName+"/subnets/"+indexs, "", "", false)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpEpgSubnetApp)
	if errs != nil {
		return errs
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}

func fetchIndex(cont *container.Container, templateName, anpName, epgName, ip string) (int, error) {
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
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {

			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return index, fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return index, err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")

					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return index, fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return index, err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						log.Println("currentepgname", currentEpgName)
						if currentEpgName == epgName {
							log.Println("found correct epgname")

							subnetCount, err := epgCont.ArrayCount("subnets")
							if err != nil {
								return index, fmt.Errorf("No Subnets found")
							}
							for s := 0; s < subnetCount; s++ {
								subnetCont, err := epgCont.ArrayElement(s, "subnets")
								if err != nil {
									return index, err
								}
								currentIp := models.StripQuotes(subnetCont.S("ip").String())
								log.Println("currentip", currentIp)
								if currentIp == ip {
									log.Println("found correct ip")
									index = s
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

		}
		if found {
			break
		}
	}
	return index, nil

}
