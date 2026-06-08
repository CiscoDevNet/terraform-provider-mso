package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMSONetflowRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSONetflowRecordCreate,
		Read:   resourceMSONetflowRecordRead,
		Update: resourceMSONetflowRecordUpdate,
		Delete: resourceMSONetflowRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSONetflowRecordImport,
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
				Description:  "The name of the NetFlow Record.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Record.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NetFlow Record.",
			},
			"match_parameters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The match parameters of the NetFlow Record.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"destination_ip",
						"destination_ipv4",
						"destination_ipv6",
						"destination_mac",
						"destination_port",
						"ethertype",
						"ip_protocol",
						"source_ip",
						"source_ipv4",
						"source_ipv6",
						"source_mac",
						"source_port",
					}, false),
				},
			},
		},
	}
}

func buildMatchParametersPayload(matchParamsRaw interface{}) []string {
	matchParamsList := matchParamsRaw.(*schema.Set).List()
	matchParams := make([]string, 0, len(matchParamsList))
	for _, item := range matchParamsList {
		if param, ok := matchParameterMap[item.(string)]; ok {
			matchParams = append(matchParams, param)
		}
	}
	return matchParams
}

func setNetflowRecordData(d *schema.ResourceData, response *container.Container, templateId string) error {
	name := models.StripQuotes(response.S("name").String())
	d.SetId(fmt.Sprintf("templateId/%s/NetflowRecord/%s", templateId, name))
	d.Set("template_id", templateId)
	d.Set("name", name)
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	if response.Exists("description") {
		d.Set("description", models.StripQuotes(response.S("description").String()))
	}

	if response.Exists("match") {
		matchCount, _ := response.ArrayCount("match")
		matchParams := make([]string, matchCount)
		for i := range matchCount {
			matchParams[i] = convertValueWithMap(models.StripQuotes(response.S("match").Index(i).String()), matchParameterMap)
		}
		d.Set("match_parameters", matchParams)
	}

	return nil
}

func resourceMSONetflowRecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Beginning Import: %v", d.Id())
	err := resourceMSONetflowRecordRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSONetflowRecordCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}

	if matchParams, ok := d.GetOk("match_parameters"); ok {
		payload["match"] = buildMatchParametersPayload(matchParams)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/netFlowRecords/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowRecord/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Create Complete: %v", d.Id())
	return resourceMSONetflowRecordRead(d, m)
}

func resourceMSONetflowRecordRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "NetflowRecord")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "netFlowRecords")
	if err != nil {
		return err
	}

	err = setNetflowRecordData(d, policy, templateId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSONetflowRecordUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowRecords")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/netFlowRecords/%d", policyIndex)

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

	if d.HasChange("match_parameters") {
		matchParams := buildMatchParametersPayload(d.Get("match_parameters"))
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/match", updatePath), matchParams)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowRecord/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Update Complete: %v", d.Id())
	return resourceMSONetflowRecordRead(d, m)
}

func resourceMSONetflowRecordDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowRecords")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/netFlowRecords/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO NetFlow Record Resource - Delete Complete: %v", d.Id())
	return nil
}
