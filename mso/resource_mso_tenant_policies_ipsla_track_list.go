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

var typeValues = map[string][2]string{
	"percentage": {"percentageUp", "percentageDown"},
	"weight":     {"weightUp", "weightDown"},
}

func resourceMSOIPSLATrackList() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOIPSLATrackListCreate,
		Read:   resourceMSOIPSLATrackListRead,
		Update: resourceMSOIPSLATrackListUpdate,
		Delete: resourceMSOIPSLATrackListDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOIPSLATrackListImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"percentage", "weight",
				}, false),
			},
			"threshold_up": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"threshold_down": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipsla_monitoring_policy_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scope_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"bd", "l3out",
							}, false),
						},
						"scope_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func setIPSLATrackListData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/IPSLATrackLists/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	d.Set("type", models.StripQuotes(response.S("type").String()))
	trackListType := d.Get("type").(string)
	d.Set("threshold_up", response.S(typeValues[trackListType][0]).Data().(float64))
	d.Set("threshold_down", response.S(typeValues[trackListType][1]).Data().(float64))

	membersDataList, err := response.S("trackListMembers").Children()
	if err != nil {
		return err
	}
	membersList := make([]map[string]interface{}, 0)
	for _, member := range membersDataList {
		trackMember := member.S("trackMember")
		memberMap := map[string]interface{}{
			"destination_ip":               models.StripQuotes(trackMember.S("destIP").String()),
			"ipsla_monitoring_policy_uuid": models.StripQuotes(trackMember.S("ipslaMonitoringRef").String()),
			"scope_type":                   models.StripQuotes(trackMember.S("scopeType").String()),
			"scope_uuid":                   models.StripQuotes(trackMember.S("scope").String()),
			"weight":                       member.S("weight").Data().(float64),
		}
		membersList = append(membersList, memberMap)
	}

	d.Set("members", membersList)

	return nil
}

func resourceMSOIPSLATrackListImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Beginning Import: %v", d.Id())
	resourceMSOIPSLATrackListRead(d, m)
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOIPSLATrackListCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}
	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	trackListType := d.Get("type").(string)

	if trackListType != "" {
		payload["type"] = trackListType
		thresholdKey := typeValues[trackListType]

		if thresholdUp, ok := d.GetOk("threshold_up"); ok {
			payload[thresholdKey[0]] = thresholdUp.(int)
		}

		if thresholdDown, ok := d.GetOk("threshold_down"); ok {
			payload[thresholdKey[1]] = thresholdDown.(int)
		}
	}

	if members, ok := d.GetOk("members"); ok {
		payload["trackListMembers"] = getIPSLATrackListMembersPayload(members.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/ipslaTrackLists/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IPSLATrackLists/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Create Complete: %v", d.Id())
	return resourceMSOIPSLATrackListRead(d, m)
}

func resourceMSOIPSLATrackListRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "IPSLATrackLists")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "ipslaTrackLists")
	if err != nil {
		return err
	}

	setIPSLATrackListData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOIPSLATrackListUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "ipslaTrackLists")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/ipslaTrackLists/%d", policyIndex)

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

	if d.HasChange("type") || d.HasChange("threshold_up") || d.HasChange("threshold_up") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/type", updatePath), d.Get("type").(string))
		if err != nil {
			return err
		}

		trackListType := d.Get("type").(string)

		if trackListType != "" {
			thresholdKey := typeValues[trackListType]

			if thresholdUp, ok := d.GetOk("threshold_up"); ok {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/%s", updatePath, thresholdKey[0]), thresholdUp.(int))
				if err != nil {
					return err
				}
			}

			if thresholdDown, ok := d.GetOk("threshold_down"); ok {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/%s", updatePath, thresholdKey[1]), thresholdDown.(int))
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("members") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/trackListMembers", updatePath), getIPSLATrackListMembersPayload(d.Get("members").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IPSLATrackLists/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Update Complete: %v", d.Id())
	return resourceMSOIPSLATrackListRead(d, m)
}

func resourceMSOIPSLATrackListDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "ipslaTrackLists")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/ipslaTrackLists/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO IPSLA Track List Resource - Delete Complete: %v", d.Id())
	return nil
}

func getIPSLATrackListMembersPayload(members *schema.Set) []interface{} {
	membersList := members.List()
	payloadMembersList := make([]interface{}, 0)
	for _, member := range membersList {
		trackMember := member.(map[string]interface{})
		payloadMembersList = append(
			payloadMembersList,
			map[string]interface{}{
				"trackMember": map[string]interface{}{
					"destIP":             trackMember["destination_ip"].(string),
					"scopeType":          trackMember["scope_type"].(string),
					"scope":              trackMember["scope_uuid"].(string),
					"ipslaMonitoringRef": trackMember["ipsla_monitoring_policy_uuid"].(string),
				},
			},
		)
	}
	return payloadMembersList
}
