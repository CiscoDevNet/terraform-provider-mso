package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const systemConfigUrl = "api/v1/platform/systemConfig"

func resourceMSOSystemConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSystemConfigCreate,
		Update: resourceMSOSystemConfigUpdate,
		Read:   resourceMSOSystemConfigRead,
		Delete: resourceMSOSystemConfigDelete,

		// Import is not defined because the create function can behave as an import when no config is provided

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"alias": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"banner": {
				Type: schema.TypeList,
				// TypeList chosen because api returns a list of banners (even though there is only one)
				// To avoid behaviour change in future decided to create list with max-elements 1
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"active",
								"inactive",
							}, false),
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"critical",
								"warning",
								"informational",
							}, false),
						},
						"message": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
			"change_control": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"enabled",
								"disabled",
							}, false),
						},
						"number_of_approvers": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		}),
	}
}

func getSystemConfig(d *schema.ResourceData, m interface{}) error {

	msoClient := m.(*client.Client)
	con, err := msoClient.GetViaURL(systemConfigUrl)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.Search("systemConfigs").Search("id").String()))

	bannerConfig := con.Search("systemConfigs").Search("bannerConfig").Data().([]interface{})[0].(map[string]interface{})
	banner := bannerConfig["banner"].(map[string]interface{})
	d.Set("alias", bannerConfig["alias"].(string))
	bannerMap := map[string]interface{}{
		"state":   banner["bannerState"].(string),
		"type":    banner["bannerType"].(string),
		"message": banner["message"].(string),
	}
	d.Set("banner", []interface{}{bannerMap})

	changeControl := con.Search("systemConfigs").Search("changeControl").Data().(map[string]interface{})
	enable := "disabled"
	if changeControl["enable"].(bool) {
		enable = "enabled"
	}
	changeControlMap := map[string]interface{}{
		"enable":              enable,
		"number_of_approvers": strconv.Itoa(int(changeControl["numOfApprovers"].(float64))),
	}
	d.Set("change_control", changeControlMap)

	return nil
}

func patchSystemConfig(d *schema.ResourceData, msoClient *client.Client, systemConfigId string) error {

	var patchPayloads []models.Model

	changeControl := d.Get("change_control").(interface{})
	if changeControl != nil {
		changeControlMap := changeControl.(map[string]interface{})
		enable := false
		if changeControlMap["enable"].(string) == "enabled" {
			enable = true
		}
		approvers, err := strconv.Atoi(changeControlMap["number_of_approvers"].(string))
		if err != nil {
			return err
		}
		patchPayloads = append(patchPayloads, models.NewSystemConfigChangeControl(enable, approvers))
	}

	alias := d.Get("alias").(string)
	banner := d.Get("banner").(interface{})
	if (banner != nil && len(banner.([]interface{})) > 0) || alias != "" {
		bannerMap := banner.([]interface{})[0].(map[string]interface{})
		patchPayloads = append(patchPayloads, models.NewSystemConfigBanner(
			alias,
			bannerMap["state"].(string),
			bannerMap["type"].(string),
			bannerMap["message"].(string),
		))
	}

	if len(patchPayloads) > 0 {
		_, err := msoClient.PatchbyID(fmt.Sprintf("%s/%s", systemConfigUrl, systemConfigId), patchPayloads...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	d.SetId(systemConfigId)

	return nil
}

func resourceMSOSystemConfigCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] System Config Beginning Creation")
	msoClient := m.(*client.Client)

	con, err := msoClient.GetViaURL(systemConfigUrl)
	if err != nil {
		return err
	}
	systemConfigId := models.StripQuotes(con.Search("systemConfigs").Search("id").String())

	err = patchSystemConfig(d, msoClient, systemConfigId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: System Config Creation finished successfully", d.Id())

	return resourceMSOSystemConfigRead(d, m)
}

func resourceMSOSystemConfigUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Update")

	err := patchSystemConfig(d, m.(*client.Client), d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Schema Update finished successfully", d.Id())
	return resourceMSOSystemConfigRead(d, m)
}

func resourceMSOSystemConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	err := getSystemConfig(d, m)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSystemConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	log.Printf("[DEBUG] Destroy finished successfully")
	return nil
}
