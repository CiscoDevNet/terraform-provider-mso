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

func resourceMSOMLDSnoopingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOMLDSnoopingPolicyCreate,
		Read:   resourceMSOMLDSnoopingPolicyRead,
		Update: resourceMSOMLDSnoopingPolicyUpdate,
		Delete: resourceMSOMLDSnoopingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOMLDSnoopingPolicyImport,
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
				Description:  "The name of the MLD Snooping Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the MLD Snooping Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the MLD Snooping Policy.",
			},
			"admin_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
				Description: "The administrative state of the MLD Snooping Policy. Default: disabled when unset during creation.",
			},
			"fast_leave_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable fast leave control. Default: false (disabled) when unset during creation.",
			},
			"querier_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable querier control. Default: false (disabled) when unset during creation.",
			},
			"querier_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false),
				Description: "The querier version. Default: v2 when unset during creation.",
			},
			"query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 18000),
				Description:  "The query interval in seconds. Default: 125 when unset during creation. Range: 1-18000.",
			},
			"query_response_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 25),
				Description:  "The query response interval in seconds. Default: 10 when unset during creation. Range: 1-25.",
			},
			"last_member_query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 25),
				Description:  "The last member query interval in seconds. Default: 1 when unset during creation. Range: 1-25.",
			},
			"start_query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 18000),
				Description:  "The start query interval in seconds. Default: 31 when unset during creation. Range: 1-18000.",
			},
			"start_query_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "The start query count. Default: 2 when unset during creation. Range: 1-10.",
			},
		},
	}
}

func setMLDSnoopingPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/MLDSnoopingPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	d.Set("admin_state", models.StripQuotes(response.S("enableAdminState").String()))

	if response.Exists("enableFastLeaveControl") {
		fastLeave, _ := response.S("enableFastLeaveControl").Data().(bool)
		d.Set("fast_leave_control", fastLeave)
	}

	if response.Exists("enableQuerierControl") {
		querierControl, _ := response.S("enableQuerierControl").Data().(bool)
		d.Set("querier_control", querierControl)
	}

	d.Set("querier_version", models.StripQuotes(response.S("mldQuerierVersion").String()))

	if response.Exists("queryInterval") {
		d.Set("query_interval", int(response.S("queryInterval").Data().(float64)))
	}
	if response.Exists("queryResponseInterval") {
		d.Set("query_response_interval", int(response.S("queryResponseInterval").Data().(float64)))
	}
	if response.Exists("lastMemberQueryInterval") {
		d.Set("last_member_query_interval", int(response.S("lastMemberQueryInterval").Data().(float64)))
	}
	if response.Exists("startQueryInterval") {
		d.Set("start_query_interval", int(response.S("startQueryInterval").Data().(float64)))
	}
	if response.Exists("startQueryCount") {
		d.Set("start_query_count", int(response.S("startQueryCount").Data().(float64)))
	}

	return nil
}

func resourceMSOMLDSnoopingPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOMLDSnoopingPolicyRead(d, m)
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOMLDSnoopingPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if adminState, ok := d.GetOk("admin_state"); ok {
		payload["enableAdminState"] = adminState.(string)
	}

	if fastLeave, ok := d.GetOk("fast_leave_control"); ok {
		payload["enableFastLeaveControl"] = fastLeave.(bool)
	}

	if querierControl, ok := d.GetOk("querier_control"); ok {
		payload["enableQuerierControl"] = querierControl.(bool)
	}

	if querierVersion, ok := d.GetOk("querier_version"); ok {
		payload["mldQuerierVersion"] = querierVersion.(string)
	}

	if queryInterval, ok := d.GetOk("query_interval"); ok {
		payload["queryInterval"] = queryInterval.(int)
	}

	if queryResponseInterval, ok := d.GetOk("query_response_interval"); ok {
		payload["queryResponseInterval"] = queryResponseInterval.(int)
	}

	if lastMemberQueryInterval, ok := d.GetOk("last_member_query_interval"); ok {
		payload["lastMemberQueryInterval"] = lastMemberQueryInterval.(int)
	}

	if startQueryInterval, ok := d.GetOk("start_query_interval"); ok {
		payload["startQueryInterval"] = startQueryInterval.(int)
	}

	if startQueryCount, ok := d.GetOk("start_query_count"); ok {
		payload["startQueryCount"] = startQueryCount.(int)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/mldSnoopPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/MLDSnoopingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOMLDSnoopingPolicyRead(d, m)
}

func resourceMSOMLDSnoopingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "MLDSnoopingPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "mldSnoopPolicies")
	if err != nil {
		return err
	}

	setMLDSnoopingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOMLDSnoopingPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "mldSnoopPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/mldSnoopPolicies/%d", policyIndex)

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

	if d.HasChange("admin_state") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableAdminState", updatePath), d.Get("admin_state").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fast_leave_control") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableFastLeaveControl", updatePath), d.Get("fast_leave_control").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("querier_control") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableQuerierControl", updatePath), d.Get("querier_control").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("querier_version") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/mldQuerierVersion", updatePath), d.Get("querier_version").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("query_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/queryInterval", updatePath), d.Get("query_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("query_response_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/queryResponseInterval", updatePath), d.Get("query_response_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("last_member_query_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/lastMemberQueryInterval", updatePath), d.Get("last_member_query_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("start_query_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/startQueryInterval", updatePath), d.Get("start_query_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("start_query_count") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/startQueryCount", updatePath), d.Get("start_query_count").(int))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/MLDSnoopingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOMLDSnoopingPolicyRead(d, m)
}

func resourceMSOMLDSnoopingPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "mldSnoopPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/mldSnoopPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO MLD Snooping Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
