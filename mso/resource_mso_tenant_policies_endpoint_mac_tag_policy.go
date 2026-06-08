package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMSOEndpointMACTagPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOEndpointMACTagPolicyCreate,
		Read:   resourceMSOEndpointMACTagPolicyRead,
		Update: resourceMSOEndpointMACTagPolicyUpdate,
		Delete: resourceMSOEndpointMACTagPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOEndpointMACTagPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bd_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vrf_uuid"},
			},
			"vrf_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bd_uuid"},
			},
			"tag_annotations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"policy_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func setEndpointMACTagPolicyId(m interface{}, templateId, name, mac, bdUUID, vrfUUID string) (string, string, error) {
	if name == "" {
		msoClient := m.(*client.Client)
		if bdUUID != "" {
			bridgeDomainObject, err := GetTemplateObjectByUUID(msoClient, "bd", bdUUID)
			if err != nil {
				return "", "", err
			}
			name = fmt.Sprintf("%s-[%s]", mac, models.StripQuotes(bridgeDomainObject.S("name").String()))
		} else if vrfUUID != "" {
			name = fmt.Sprintf("%s-[*]", mac)
		}
	}
	return name, fmt.Sprintf("templateId/%s/EndpointMACTagPolicy/%s", templateId, name), nil
}

func getTagAnnotationsPayload(tagAnnotationsSet *schema.Set) []map[string]interface{} {
	if tagAnnotationsSet == nil || tagAnnotationsSet.Len() == 0 {
		return []map[string]interface{}{}
	}
	list := tagAnnotationsSet.List()
	payload := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"tagKey":   m["key"].(string),
			"tagValue": m["value"].(string),
		})
	}
	return payload
}

func getPolicyTagsPayload(policyTagsSet *schema.Set) []map[string]interface{} {
	if policyTagsSet == nil || policyTagsSet.Len() == 0 {
		return []map[string]interface{}{}
	}
	list := policyTagsSet.List()
	payload := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"key":   m["key"].(string),
			"value": m["value"].(string),
		})
	}
	return payload
}

func setEndpointMACTagPolicyData(d *schema.ResourceData, response *container.Container, templateId string, m interface{}) error {
	d.Set("template_id", templateId)
	d.Set("mac", models.StripQuotes(response.S("mac").String()))
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	if response.Exists("bdRef") {
		d.Set("bd_uuid", models.StripQuotes(response.S("bdRef").String()))
	}
	if response.Exists("vrfRef") {
		d.Set("vrf_uuid", models.StripQuotes(response.S("vrfRef").String()))
	}

	tagAnnotationsList := make([]map[string]interface{}, 0)
	if count, err := response.ArrayCount("tagAnnotations"); err == nil {
		for i := range count {
			c, err := response.ArrayElement(i, "tagAnnotations")
			if err != nil {
				return err
			}
			tagAnnotationsList = append(tagAnnotationsList, map[string]interface{}{
				"key":   models.StripQuotes(c.S("tagKey").String()),
				"value": models.StripQuotes(c.S("tagValue").String()),
			})
		}
	}
	d.Set("tag_annotations", tagAnnotationsList)

	policyTagsList := make([]map[string]interface{}, 0)
	if count, err := response.ArrayCount("policyTags"); err == nil {
		for i := range count {
			c, err := response.ArrayElement(i, "policyTags")
			if err != nil {
				return err
			}
			policyTagsList = append(policyTagsList, map[string]interface{}{
				"key":   models.StripQuotes(c.S("key").String()),
				"value": models.StripQuotes(c.S("value").String()),
			})
		}
	}
	d.Set("policy_tags", policyTagsList)

	return nil
}

func resourceMSOEndpointMACTagPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Beginning Import: %v", d.Id())
	err := resourceMSOEndpointMACTagPolicyRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOEndpointMACTagPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	// Ignoring the input validation since the API returns a clear error when bd_uuid and vrf_uuid are not provided.
	bdUUID := d.Get("bd_uuid").(string)
	vrfUUID := d.Get("vrf_uuid").(string)

	payload := map[string]interface{}{
		"mac":        d.Get("mac").(string),
		"templateId": templateId,
	}

	if bdUUID != "" {
		payload["bdRef"] = bdUUID
	} else if vrfUUID != "" {
		payload["vrfRef"] = vrfUUID
	}

	if tagAnnotations, ok := d.GetOk("tag_annotations"); ok {
		payload["tagAnnotations"] = getTagAnnotationsPayload(tagAnnotations.(*schema.Set))
	}

	if policyTags, ok := d.GetOk("policy_tags"); ok {
		payload["policyTags"] = getPolicyTagsPayload(policyTags.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/endpointMacTagPolicies/-", payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	_, resourceId, err := setEndpointMACTagPolicyId(m, templateId, "", d.Get("mac").(string), d.Get("bd_uuid").(string), d.Get("vrf_uuid").(string))
	if err != nil {
		return err
	}

	d.SetId(resourceId)
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOEndpointMACTagPolicyRead(d, m)
}

func resourceMSOEndpointMACTagPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "EndpointMACTagPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "endpointMacTagPolicies")
	if err != nil {
		return err
	}

	err = setEndpointMACTagPolicyData(d, policy, templateId, m)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOEndpointMACTagPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "endpointMacTagPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/endpointMacTagPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("mac") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/mac", updatePath), d.Get("mac").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("bd_uuid") {
		if bdUUID, ok := d.GetOk("bd_uuid"); ok && bdUUID.(string) != "" {
			// Remove vrfRef if bd_uuid is being set, since both cannot coexist and NDO does not remove it automatically when bdRef is added
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/vrfRef", updatePath), nil)
			if err != nil {
				return err
			}

			err = addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/bdRef", updatePath), bdUUID.(string))
			if err != nil {
				return err
			}
			d.Set("vrf_uuid", "")
		}
	}

	if d.HasChange("vrf_uuid") {
		if vrfUUID, ok := d.GetOk("vrf_uuid"); ok && vrfUUID.(string) != "" {
			// Remove bdRef if vrf_uuid is being set, since both cannot coexist and NDO does not remove it automatically when vrfRef is added
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/bdRef", updatePath), nil)
			if err != nil {
				return err
			}

			err = addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/vrfRef", updatePath), vrfUUID.(string))
			if err != nil {
				return err
			}
			d.Set("bd_uuid", "")
		}
	}

	if d.HasChange("tag_annotations") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/tagAnnotations", updatePath), getTagAnnotationsPayload(d.Get("tag_annotations").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	if d.HasChange("policy_tags") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/policyTags", updatePath), getPolicyTagsPayload(d.Get("policy_tags").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	_, resourceId, err := setEndpointMACTagPolicyId(m, templateId, "", d.Get("mac").(string), d.Get("bd_uuid").(string), d.Get("vrf_uuid").(string))
	if err != nil {
		return err
	}

	d.SetId(resourceId)
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOEndpointMACTagPolicyRead(d, m)
}

func resourceMSOEndpointMACTagPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "endpointMacTagPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/endpointMacTagPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
