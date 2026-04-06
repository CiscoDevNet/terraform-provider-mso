package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOPortChannelInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOPortChannelInterfaceCreate,
		Read:   resourceMSOPortChannelInterfaceRead,
		Update: resourceMSOPortChannelInterfaceUpdate,
		Delete: resourceMSOPortChannelInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOPortChannelInterfaceImport,
		},

		CustomizeDiff: customizeDiffPortChannelInterface,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interfaces": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"interface_policy_group_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interface_descriptions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customizeDiffPortChannelInterface(d *schema.ResourceDiff, m interface{}) error {
	interfacesRaw := d.Get("interfaces")
	if interfacesRaw == nil {
		return nil
	}
	interfaces := interfacesRaw.(*schema.Set)
	if interfaces.Len() == 0 {
		return nil
	}

	// Create a map for quick lookup
	interfaceMap := make(map[string]bool)
	for _, memberInterface := range interfaces.List() {
		interfaceMap[memberInterface.(string)] = true
	}

	interfaceDescriptionsRaw := d.Get("interface_descriptions")
	if interfaceDescriptionsRaw == nil {
		return nil
	}
	interfaceDescriptions := interfaceDescriptionsRaw.(*schema.Set)
	if interfaceDescriptions.Len() == 0 {
		return nil
	}

	// Validate each interface_description
	for _, interfaceDescriptionRaw := range interfaceDescriptions.List() {
		interfaceDescription := interfaceDescriptionRaw.(map[string]interface{})
		interfaceID := interfaceDescription["interface"].(string)

		if !interfaceMap[interfaceID] {
			return fmt.Errorf("Interface '%s' specified in 'interface_descriptions' must be defined in the 'interfaces' list", interfaceID)
		}
	}

	return nil
}

func getInterfaceDescriptionsPayload(node string, interfaceDescriptions *schema.Set) []map[string]interface{} {
	interfaceDescriptionsList := interfaceDescriptions.List()
	payload := make([]map[string]interface{}, 0)
	for _, interfaceDescription := range interfaceDescriptionsList {
		interfaceDescriptionMap := interfaceDescription.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"nodeID":      node,
			"interfaceID": interfaceDescriptionMap["interface"].(string),
			"description": interfaceDescriptionMap["description"].(string),
		})
	}
	return payload
}

func setPortChannelInterfaceData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/PortChannelInterface/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	d.Set("node", models.StripQuotes(response.S("node").String()))
	d.Set("interfaces", splitCommaString(models.StripQuotes(response.S("memberInterfaces").String())))

	if response.Exists("policy") {
		d.Set("interface_policy_group_uuid", models.StripQuotes(response.S("policy").String()))
	}

	interfaceDescriptionsList := make([]map[string]interface{}, 0)
	count, err := response.ArrayCount("interfaceDescriptions")
	if err == nil {
		for i := range count {
			descriptionContainer, err := response.ArrayElement(i, "interfaceDescriptions")
			if err != nil {
				return err
			}
			interfaceDescriptionsList = append(interfaceDescriptionsList, map[string]interface{}{
				"interface":   models.StripQuotes(descriptionContainer.S("interfaceID").String()),
				"description": models.StripQuotes(descriptionContainer.S("description").String()),
			})
		}
	}

	d.Set("interface_descriptions", interfaceDescriptionsList)

	return nil
}

func resourceMSOPortChannelInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Beginning Import: %v", d.Id())
	err := resourceMSOPortChannelInterfaceRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOPortChannelInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	payload := map[string]interface{}{
		"name":       d.Get("name").(string),
		"templateId": templateId,
	}

	payload["description"] = d.Get("description").(string)

	payload["node"] = d.Get("node").(string)
	payload["memberInterfaces"] = strings.Join(getListOfStringsFromSchemaSet(d, "interfaces"), ",")
	payload["policy"] = d.Get("interface_policy_group_uuid").(string)

	if interfaceDescriptions, ok := d.GetOk("interface_descriptions"); ok {
		payload["interfaceDescriptions"] = getInterfaceDescriptionsPayload(payload["node"].(string), interfaceDescriptions.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/fabricResourceTemplate/template/portChannels/-", payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/PortChannelInterface/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Create Complete: %v", d.Id())
	return resourceMSOPortChannelInterfaceRead(d, m)
}

func resourceMSOPortChannelInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "PortChannelInterface")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricResourceTemplate", "template", "portChannels")
	if err != nil {
		return err
	}

	err = setPortChannelInterfaceData(d, policy, templateId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOPortChannelInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)
	node := d.Get("node").(string)

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "portChannels")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricResourceTemplate/template/portChannels/%d", policyIndex)

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

	if d.HasChange("node") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/node", updatePath), d.Get("node").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("interfaces") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/memberInterfaces", updatePath), strings.Join(getListOfStringsFromSchemaSet(d, "interfaces"), ","))
		if err != nil {
			return err
		}
	}

	if d.HasChange("interface_policy_group_uuid") {
		interfacePolicyUUID := d.Get("interface_policy_group_uuid").(string)
		if interfacePolicyUUID != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/policy", updatePath), interfacePolicyUUID)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("interface_descriptions") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaceDescriptions", updatePath), getInterfaceDescriptionsPayload(node, d.Get("interface_descriptions").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/PortChannelInterface/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Update Complete: %v", d.Id())
	return resourceMSOPortChannelInterfaceRead(d, m)
}

func resourceMSOPortChannelInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "portChannels")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricResourceTemplate/template/portChannels/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Port Channel Interface Resource - Delete Complete: %v", d.Id())
	return nil
}
