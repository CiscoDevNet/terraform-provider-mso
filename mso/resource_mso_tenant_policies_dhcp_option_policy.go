package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSODHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSODHCPOptionPolicyCreate,
		Read:   resourceMSODHCPOptionPolicyRead,
		Update: resourceMSODHCPOptionPolicyUpdate,
		Delete: resourceMSODHCPOptionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSODHCPOptionPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the DHCP Option Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the DHCP Option Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the DHCP Option Policy.",
			},
			"options": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "A set of DHCP options. At least one option is required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the DHCP option.",
						},
						"id": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 255),
							Description:  "The ID of the DHCP option. Range: 0-255.",
						},
						"data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The data value of the DHCP option.",
						},
					},
				},
			},
		},
	}
}

func buildDHCPOptionsPayload(optionsRaw interface{}) []interface{} {
	optionsSet := optionsRaw.(*schema.Set)
	optionsList := optionsSet.List()
	options := make([]interface{}, 0, len(optionsList))

	for _, item := range optionsList {
		option := item.(map[string]interface{})
		optionPayload := map[string]interface{}{
			"name": option["name"].(string),
		}

		if optionId, ok := option["id"].(int); ok && optionId != 0 {
			optionPayload["id"] = optionId
		}

		if data, ok := option["data"].(string); ok && data != "" {
			optionPayload["data"] = data
		}

		options = append(options, optionPayload)
	}

	return options
}

func setDHCPOptionPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/DHCPOptionPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	if response.Exists("options") {
		optionsCount, _ := response.S("options").ArrayCount()
		options := make([]interface{}, 0, optionsCount)

		for i := 0; i < optionsCount; i++ {
			option := response.S("options").Index(i)
			optionMap := map[string]interface{}{
				"name": models.StripQuotes(option.S("name").String()),
			}

			if option.Exists("id") {
				if optionId, ok := option.S("id").Data().(float64); ok {
					optionMap["id"] = int(optionId)
				}
			}

			if option.Exists("data") {
				optionMap["data"] = models.StripQuotes(option.S("data").String())
			}

			options = append(options, optionMap)
		}
		d.Set("options", options)
	}

	return nil
}

func resourceMSODHCPOptionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Beginning Import: %v", d.Id())
	resourceMSODHCPOptionPolicyRead(d, m)
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSODHCPOptionPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if optionsRaw, ok := d.GetOk("options"); ok {
		payload["options"] = buildDHCPOptionsPayload(optionsRaw)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/dhcpOptionPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/DHCPOptionPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Create Complete: %v", d.Id())
	return resourceMSODHCPOptionPolicyRead(d, m)
}

func resourceMSODHCPOptionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "DHCPOptionPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "dhcpOptionPolicies")
	if err != nil {
		return err
	}

	setDHCPOptionPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSODHCPOptionPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "dhcpOptionPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/dhcpOptionPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
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

	if d.HasChange("options") {
		options := buildDHCPOptionsPayload(d.Get("options"))
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/options", updatePath), options)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/DHCPOptionPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Update Complete: %v", d.Id())
	return resourceMSODHCPOptionPolicyRead(d, m)
}

func resourceMSODHCPOptionPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "dhcpOptionPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/dhcpOptionPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO DHCP Option Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
