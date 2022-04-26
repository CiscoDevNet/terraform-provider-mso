package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgUsegAttrCreate,
		Read:   resourceMSOSchemaSiteAnpEpgUsegAttrRead,
		Delete: resourceMSOSchemaSiteAnpEpgUsegAttrDelete,
		Update: resourceMSOSchemaSiteAnpEpgUsegAttrUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgUsegAttrImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
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
			"anp_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"mac",
					"vm-name",      // Vm Name
					"rootContName", // VM data center
					"hv",           // Hypervisor
					"guest-os",     // Operating System
					"tag",
					"vm",     // Identifier
					"domain", // VMM domain
					"vnic",
				}, false),
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"equals",
					"startsWith",
					"endsWith",
					"contains",
				}, false),
			},
			"fv_subnet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.If(
				// if useg type is ip it's value should be validated
				// 1. condition method
				func(d *schema.ResourceDiff, m interface{}) bool {
					useg_type := d.Get("useg_type").(string)
					return useg_type == "ip"
				},
				// 2. validation method
				customdiff.ValidateValue("value",
					func(v, m interface{}) error {
						_, val_errors := validation.IsIPAddress(v, "value")
						var err_str error
						for _, err := range val_errors {
							if err_str == nil {
								err_str = fmt.Errorf("%s\n", err)
							} else {
								err_str = fmt.Errorf("%s\n%s", err_str, err)
							}
						}
						return err_str
					},
				),
			),
			customdiff.If(
				// if useg type is mac it's value should be validated
				// 1. condition method
				func(d *schema.ResourceDiff, m interface{}) bool {
					useg_type := d.Get("useg_type").(string)
					return useg_type == "mac"
				},
				// 2. validation method
				customdiff.ValidateValue("value",
					func(v, m interface{}) error {
						_, val_errors := validation.IsMACAddress(v, "value")
						var err_str error
						for _, err := range val_errors {
							if err_str == nil {
								err_str = fmt.Errorf("%s\n", err)
							} else {
								err_str = fmt.Errorf("%s\n%s", err_str, err)
							}
						}
						return err_str
					},
				),
			),
		),
	}
}

func setMSOSchemaSiteAnpEpgUsegAttributes(usegMap *models.SiteUsegAttr, d *schema.ResourceData) {
	d.Set("template_name", usegMap.TemplateName)
	d.Set("site_id", usegMap.SiteID)
	d.Set("schema_id", usegMap.SchemaID)
	d.Set("anp_name", usegMap.AnpName)
	d.Set("epg_name", usegMap.EpgName)
	d.Set("useg_name", usegMap.UsegName)
	d.Set("useg_type", usegMap.Type)
	d.Set("value", usegMap.Value)
	d.Set("description", usegMap.Description)
	d.Set("operator", usegMap.Operator)
	d.Set("category", usegMap.Category)
	d.Set("fv_subnet", usegMap.FvSubnet)
}

func resourceMSOSchemaSiteAnpEpgUsegAttrImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	useg, err := UsegIdToUsegAttrModel(id)
	if err != nil {
		return nil, err
	}
	usegMapRemote, _, err := msoClient.ReadAnpEpgUsegAttr(useg)
	if err != nil {
		return nil, err
	}
	setMSOSchemaSiteAnpEpgUsegAttributes(usegMapRemote, d)
	d.SetId(id)
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteAnpEpgUsegAttrCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Creation")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	usegName := d.Get("useg_name").(string)
	usegType := d.Get("useg_type").(string)
	value := d.Get("value").(string)

	usegMap := models.SiteUsegAttr{
		SiteID:       siteId,
		TemplateName: templateName,
		SchemaID:     schemaId,
		AnpName:      anpName,
		EpgName:      epgName,
		UsegName:     usegName,
		Type:         usegType,
		Value:        value,
	}

	if operator, ok := d.GetOk("operator"); ok {
		usegMap.Operator = operator.(string)
	}
	if category, ok := d.GetOk("category"); ok {
		usegMap.Category = category.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		usegMap.Description = description.(string)
	}
	if fv_subnet, ok := d.GetOk("fv_subnet"); ok {
		usegMap.FvSubnet = fv_subnet.(bool)
	}

	err := msoClient.CreateAnpEpgUsegAttr(&usegMap)
	if err != nil {
		return err
	}
	usegId := UsegAttrModelToUsegId(&usegMap)
	d.SetId(usegId)
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Creation Completed", d.Id())
	return resourceMSOSchemaSiteAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgUsegAttrUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Update", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	usegName := d.Get("useg_name").(string)
	usegType := d.Get("useg_type").(string)
	value := d.Get("value").(string)

	usegMap := models.SiteUsegAttr{
		SiteID:       siteId,
		TemplateName: templateName,
		SchemaID:     schemaId,
		AnpName:      anpName,
		EpgName:      epgName,
		UsegName:     usegName,
		Type:         usegType,
		Value:        value,
	}

	if operator, ok := d.GetOk("operator"); ok {
		usegMap.Operator = operator.(string)
	}
	if category, ok := d.GetOk("category"); ok {
		usegMap.Category = category.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		usegMap.Description = description.(string)
	}
	if fv_subnet, ok := d.GetOk("fv_subnet"); ok {
		usegMap.FvSubnet = fv_subnet.(bool)
	}

	err := msoClient.UpdateAnpEpgUsegAttr(&usegMap)
	if err != nil {
		return err
	}
	usegId := UsegAttrModelToUsegId(&usegMap)
	d.SetId(usegId)
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Update Completed", d.Id())
	return resourceMSOSchemaSiteAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	useg, err := UsegIdToUsegAttrModel(id)
	if err != nil {
		return err
	}
	usegMapRemote, _, err := msoClient.ReadAnpEpgUsegAttr(useg)
	if err != nil {
		// make empty to remove entry from tfstate file as resource has been removed from the server or not found on the server
		d.SetId("")
		return nil
	}
	setMSOSchemaSiteAnpEpgUsegAttributes(usegMapRemote, d)
	d.SetId(id)
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Reading Completed", d.Id())
	return nil
}

func resourceMSOSchemaSiteAnpEpgUsegAttrDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	useg, err := UsegIdToUsegAttrModel(id)
	if err != nil {
		return err
	}
	err = msoClient.DeleteAnpEpgUsegAttr(useg)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Schema Site Anp Epg Useg Attr: Beginning Destroy", d.Id())
	// make empty to remove entry from tfstate file as resource has been destroyed
	d.SetId("")
	return err
}

func UsegAttrModelToUsegId(m *models.SiteUsegAttr) string {
	return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s/uSegAttr/%s", m.SchemaID, m.SiteID, m.TemplateName, m.AnpName, m.EpgName, m.UsegName)
}

func UsegIdToUsegAttrModel(id string) (*models.SiteUsegAttr, error) {
	re := regexp.MustCompile("(.*)/site/(.*)/template/(.*)/anp/(.*)/epg/(.*)/uSegAttr/(.*)")
	match := re.FindStringSubmatch(id)

	if len(match) <= 0 {
		return nil, fmt.Errorf("invalid mso_schema_site_anp_epg_useg_attr id format")
	}
	usegMap := models.SiteUsegAttr{
		SchemaID:     match[1],
		SiteID:       match[2],
		TemplateName: match[3],
		AnpName:      match[4],
		EpgName:      match[5],
		UsegName:     match[6],
	}
	return &usegMap, nil
}
