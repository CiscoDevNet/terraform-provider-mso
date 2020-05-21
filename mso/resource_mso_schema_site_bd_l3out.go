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

func resourceMSOSchemaSiteBdL3out() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteBdL3outCreate,
		Read:   resourceMSOSchemaSiteBdL3outRead,
		Delete: resourceMSOSchemaSiteBdL3outDelete,

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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
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
		return err
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
		fmt.Errorf("The given Anp Epg Domain is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/bds/%s/l3Outs/%s", siteId, templateName, bdName, indexs)
	BdL3outStruct := models.NewSchemaSiteBdL3out("remove", path, l3outName)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), BdL3outStruct)
	if errs != nil {
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
