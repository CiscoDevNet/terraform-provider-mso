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

func resourceMSOMcastRouteMapPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOMcastRouteMapPolicyCreate,
		Read:   resourceMSOMcastRouteMapPolicyRead,
		Update: resourceMSOMcastRouteMapPolicyUpdate,
		Delete: resourceMSOMcastRouteMapPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOMcastRouteMapPolicyImport,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_map_entries_multicast": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"group_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"rp_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"permit", "deny",
							}, false),
						},
					},
				},
			},
		},
	}
}

func setMcastRouteMapEntryList(mcastRouteMapEntries *schema.Set) []map[string]any {
	multicastRouteMapEntryList := mcastRouteMapEntries.List()
	mcastRouteMapEntryList := make([]map[string]any, len(multicastRouteMapEntryList))

	for i, val := range multicastRouteMapEntryList {
		multicastRouteMapEntry := val.(map[string]any)
		mcastRouteMapEntryList[i] = map[string]any{
			"order":  multicastRouteMapEntry["order"].(int),
			"group":  multicastRouteMapEntry["group_ip"].(string),
			"source": multicastRouteMapEntry["source_ip"].(string),
			"rp":     multicastRouteMapEntry["rp_ip"].(string),
			"action": multicastRouteMapEntry["action"].(string),
		}
	}

	return mcastRouteMapEntryList
}

func setMcastRouteMapPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicyMulticast/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	count, _ := response.ArrayCount("mcastRtMapEntryList")
	mcastRouteMapEntryList := make([]any, 0)
	for i := range count {
		mcastRouteMapEntryCont, err := response.ArrayElement(i, "mcastRtMapEntryList")
		if err != nil {
			return fmt.Errorf("unable to parse the multicast route map entries list")
		}
		mcastRouteMapEntry := make(map[string]any)
		mcastRouteMapEntry["order"] = mcastRouteMapEntryCont.S("order").Data().(float64)
		mcastRouteMapEntry["group_ip"] = models.StripQuotes(mcastRouteMapEntryCont.S("group").String())
		mcastRouteMapEntry["source_ip"] = models.StripQuotes(mcastRouteMapEntryCont.S("source").String())
		mcastRouteMapEntry["rp_ip"] = models.StripQuotes(mcastRouteMapEntryCont.S("rp").String())
		mcastRouteMapEntry["action"] = models.StripQuotes(mcastRouteMapEntryCont.S("action").String())
		mcastRouteMapEntryList = append(mcastRouteMapEntryList, mcastRouteMapEntry)
	}
	d.Set("route_map_entries_multicast", mcastRouteMapEntryList)

	return nil

}

func resourceMSOMcastRouteMapPolicyImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Beginning Import: %v", d.Id())
	resourceMSOMcastRouteMapPolicyRead(d, m)
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOMcastRouteMapPolicyCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if mcastRouteMapEntries, ok := d.GetOk("route_map_entries_multicast"); ok {
		payload["mcastRtMapEntryList"] = setMcastRouteMapEntryList(mcastRouteMapEntries.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/mcastRouteMapPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicyMulticast/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Create Complete: %v", d.Id())
	return resourceMSOMcastRouteMapPolicyRead(d, m)
}

func resourceMSOMcastRouteMapPolicyRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "RouteMapPolicyMulticast")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "mcastRouteMapPolicies")
	if err != nil {
		return err
	}

	setMcastRouteMapPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOMcastRouteMapPolicyUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "mcastRouteMapPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/mcastRouteMapPolicies/%d", policyIndex)

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

	if d.HasChange("route_map_entries_multicast") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/mcastRtMapEntryList", updatePath), setMcastRouteMapEntryList(d.Get("route_map_entries_multicast").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicyMulticast/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Update Complete: %v", d.Id())
	return resourceMSOMcastRouteMapPolicyRead(d, m)
}

func resourceMSOMcastRouteMapPolicyDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "mcastRouteMapPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/mcastRouteMapPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Resource - Delete Complete: %v", d.Id())
	return nil
}
