package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteL3out() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteL3outCreate,
		Read:   resourceMSOSchemaSiteL3outRead,
		Delete: resourceMSOSchemaSiteL3outDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteL3outImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"l3out_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func setMSOSchemaSiteL3outAttributes(l3outMap *models.IntersiteL3outs, d *schema.ResourceData) {
	d.Set("l3out_name", l3outMap.L3outName)
	d.Set("vrf_name", l3outMap.VRFName)
	d.Set("template_name", l3outMap.TemplateName)
	d.Set("site_id", l3outMap.SiteId)
	d.Set("schema_id", l3outMap.SchemaID)
}

func resourceMSOSchemaSiteL3outImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Schema Site L3out: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	L3out, err := L3outIdToL3outModel(id)
	if err != nil {
		return nil, err
	}
	l3outMapRemote, err := msoClient.ReadIntersiteL3outs(L3out)
	if err != nil {
		return nil, err
	}
	setMSOSchemaSiteL3outAttributes(l3outMapRemote, d)
	d.SetId(id)
	log.Println("[DEBUG] Schema Site L3out: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteL3outCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site L3out: Beginning Creation")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	l3outName := d.Get("l3out_name").(string)
	l3outMap := models.IntersiteL3outs{
		L3outName:    l3outName,
		VRFName:      vrfName,
		SiteId:       siteId,
		TemplateName: templateName,
		SchemaID:     schemaId,
	}
	err := msoClient.CreateIntersiteL3outs(&l3outMap)
	if err != nil {
		return err
	}
	l3outId := L3outModelToL3outId(&l3outMap)
	d.SetId(l3outId)
	log.Printf("[DEBUG] Schema Site L3out: Creation Completed")
	return resourceMSOSchemaSiteL3outRead(d, m)
}

func resourceMSOSchemaSiteL3outRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site L3out: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	l3out, err := L3outIdToL3outModel(id)
	if err != nil {
		return err
	}
	l3outMapRemote, err := msoClient.ReadIntersiteL3outs(l3out)
	if err != nil {
		d.SetId("")
		return nil
	}
	setMSOSchemaSiteL3outAttributes(l3outMapRemote, d)
	d.SetId(L3outModelToL3outId(l3out))
	log.Println("[DEBUG] Schema Site L3out: Reading Completed", d.Id())
	return nil
}

func resourceMSOSchemaSiteL3outDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site L3out: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	l3out, err := L3outIdToL3outModel(id)
	if err != nil {
		return err
	}
	err = msoClient.DeleteIntersiteL3outs(l3out)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Schema Site L3out: Beginning Destroy", d.Id())
	d.SetId("")
	return err
}

func L3outModelToL3outId(m *models.IntersiteL3outs) string {
	return fmt.Sprintf("%s/site/%s/template/%s/vrf/%s/l3out/%s", m.SchemaID, m.SiteId, m.TemplateName, m.VRFName, m.L3outName)
}

func L3outIdToL3outModel(id string) (*models.IntersiteL3outs, error) {
	getAttributes := strings.Split(id, "/")
	if len(getAttributes) != 9 || getAttributes[1] != "site" || getAttributes[3] != "template" || getAttributes[5] != "vrf" || getAttributes[7] != "l3out" {
		return nil, fmt.Errorf("invalid mso_schema_site_l3out id format")
	}
	l3outMap := models.IntersiteL3outs{
		SchemaID:     getAttributes[0],
		SiteId:       getAttributes[2],
		TemplateName: getAttributes[4],
		VRFName:      getAttributes[6],
		L3outName:    getAttributes[8],
	}
	return &l3outMap, nil
}
