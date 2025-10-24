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

func resourceMSOSyncEInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSyncEInterfacePolicyCreate,
		Read:   resourceMSOSyncEInterfacePolicyRead,
		Update: resourceMSOSyncEInterfacePolicyUpdate,
		Delete: resourceMSOSyncEInterfacePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOSyncEInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"sync_state_msg": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"selection_input": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"src_priority": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 254),
			},
			"wait_to_restore": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 12),
			},
		},
	}
}

func setSyncEInterfacePolicyData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "syncEthIntfPolicies")
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/SyncEInterfacePolicy/%s", templateId, models.StripQuotes(policy.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(policy.S("name").String()))
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("uuid").String()))
	d.Set("admin_state", models.StripQuotes(policy.S("adminState").String()))
	syncStateMsg := "disabled"
	if policy.S("syncStateMsgEnabled").Data().(bool) {
		syncStateMsg = "enabled"
	}
	d.Set("sync_state_msg", syncStateMsg)
	selectionInput := "disabled"
	if policy.S("selectionInputEnabled").Data().(bool) {
		selectionInput = "enabled"
	}
	d.Set("selection_input", selectionInput)
	d.Set("src_priority", policy.S("srcPriority").Data().(float64))
	d.Set("wait_to_restore", policy.S("waitToRestore").Data().(float64))

	return nil

}

func resourceMSOSyncEInterfacePolicyImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "SyncEInterfacePolicy")
	if err != nil {
		return nil, err
	}

	setSyncEInterfacePolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSyncEInterfacePolicyCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if adminState, ok := d.GetOk("admin_state"); ok {
		payload["adminState"] = adminState.(string)
	}

	if syncStateMsg, ok := d.GetOk("sync_state_msg"); ok {
		if syncStateMsg.(string) == "enabled" {
			payload["syncStateMsgEnabled"] = true
		} else {
			payload["syncStateMsgEnabled"] = false
		}
	}

	if selectionInput, ok := d.GetOk("selection_input"); ok {
		if selectionInput.(string) == "enabled" {
			payload["selectionInputEnabled"] = true
		} else {
			payload["selectionInputEnabled"] = false
		}
	}

	if srcPriority, ok := d.GetOk("src_priority"); ok {
		payload["srcPriority"] = srcPriority.(int)
	}

	if waitToRestore, ok := d.GetOk("wait_to_restore"); ok {
		payload["waitToRestore"] = waitToRestore.(int)
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/syncEthIntfPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/SyncEInterfacePolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOSyncEInterfacePolicyRead(d, m)
}

func resourceMSOSyncEInterfacePolicyRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setSyncEInterfacePolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOSyncEInterfacePolicyUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "syncEthIntfPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/syncEthIntfPolicies/%d", policyIndex)

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
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/adminState", updatePath), d.Get("admin_state").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("sync_state_msg") {
		syncStateMsg := false
		if d.Get("sync_state_msg").(string) == "enabled" {
			syncStateMsg = true
		}
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/syncStateMsgEnabled", updatePath), syncStateMsg)
		if err != nil {
			return err
		}
	}

	if d.HasChange("selection_input") {
		selectionInput := false
		if d.Get("selection_input").(string) == "enabled" {
			selectionInput = true
		}
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/selectionInputEnabled", updatePath), selectionInput)
		if err != nil {
			return err
		}
	}

	if d.HasChange("src_priority") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/srcPriority", updatePath), d.Get("src_priority").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("wait_to_restore") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/waitToRestore", updatePath), d.Get("wait_to_restore").(int))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/SyncEInterfacePolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOSyncEInterfacePolicyRead(d, m)
}

func resourceMSOSyncEInterfacePolicyDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "syncEthIntfPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/syncEthIntfPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO SyncE Interface Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
