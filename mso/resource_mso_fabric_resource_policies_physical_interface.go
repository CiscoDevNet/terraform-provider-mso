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

func resourceMSOPhysicalInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOPhysicalInterfaceCreate,
		Read:   resourceMSOPhysicalInterfaceRead,
		Update: resourceMSOPhysicalInterfaceUpdate,
		Delete: resourceMSOPhysicalInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOPhysicalInterfaceImport,
		},

		CustomizeDiff: customizeDiffPhysicalInterface,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"interfaces": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"interface_policy_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"breakout_mode"},
			},
			"breakout_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"4x100G",
					"4x25G",
					"4x10G",
				}, false),
				ConflictsWith: []string{"interface_policy_uuid"},
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
							Computed: true,
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customizeDiffPhysicalInterface(d *schema.ResourceDiff, m interface{}) error {
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

func getInterfaceDescriptionsPayloadPhysical(interfaceDescriptions *schema.Set) []map[string]interface{} {
	interfaceDescriptionsList := interfaceDescriptions.List()
	payload := make([]map[string]interface{}, 0, len(interfaceDescriptionsList))
	for _, interfaceDescription := range interfaceDescriptionsList {
		interfaceDescriptionMap := interfaceDescription.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"interfaceID": interfaceDescriptionMap["interface"].(string),
			"description": interfaceDescriptionMap["description"].(string),
		})
	}
	return payload
}

func nodeSetToStringList(nodeSet *schema.Set) []interface{} {
	nodeList := make([]interface{}, 0)
	if nodeSet != nil && nodeSet.Len() > 0 {
		nodeList = make([]interface{}, 0, nodeSet.Len())
		for _, node := range nodeSet.List() {
			nodeList = append(nodeList, node.(string))
		}
	}
	return nodeList
}

func physicalInterfaceSetToString(interfaceSet *schema.Set) string {
	interfaceList := make([]string, 0)
	if interfaceSet != nil && interfaceSet.Len() > 0 {
		interfaceList = make([]string, 0, interfaceSet.Len())
		for _, memberInterface := range interfaceSet.List() {
			interfaceList = append(interfaceList, memberInterface.(string))
		}
	}
	return strings.Join(interfaceList, ",")
}

func physicalInterfaceStringToSet(interfaceString string) []interface{} {
	interfaceSet := make([]interface{}, 0)
	if interfaceString != "" {
		interfaceList := strings.Split(interfaceString, ",")
		interfaceSet = make([]interface{}, 0, len(interfaceList))
		for _, interfaceID := range interfaceList {
			cleanedInterfaceID := strings.TrimSpace(interfaceID)
			if cleanedInterfaceID != "" {
				interfaceSet = append(interfaceSet, cleanedInterfaceID)
			}
		}
	}
	return interfaceSet
}

func setPhysicalInterfaceData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/PhysicalInterface/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	d.Set("interfaces", physicalInterfaceStringToSet(models.StripQuotes(response.S("interfaces").String())))
	d.Set("policy_group_type", models.StripQuotes(response.S("policyGroupType").String()))

	if response.Exists("nodes") {
		nodeCount, err := response.ArrayCount("nodes")
		if err == nil {
			nodeList := make([]interface{}, 0, nodeCount)
			for i := 0; i < nodeCount; i++ {
				nodeElement, err := response.ArrayElement(i, "nodes")
				if err == nil {
					nodeList = append(nodeList, models.StripQuotes(nodeElement.String()))
				}
			}
			d.Set("nodes", nodeList)
		}
	}

	if response.Exists("policy") {
		d.Set("interface_policy_uuid", models.StripQuotes(response.S("policy").String()))
	}

	if response.Exists("breakoutMode") {
		d.Set("breakout_mode", models.StripQuotes(response.S("breakoutMode").String()))
	}

	interfaceDescriptionsList := make([]map[string]interface{}, 0)
	count, err := response.ArrayCount("interfaceDescriptions")
	if err == nil {
		for i := 0; i < count; i++ {
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

func resourceMSOPhysicalInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Physical Interface Resource - Beginning Import: %v", d.Id())
	resourceMSOPhysicalInterfaceRead(d, m)
	log.Printf("[DEBUG] MSO Physical Interface Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOPhysicalInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Interface Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	payload := map[string]interface{}{
		"name":       d.Get("name").(string),
		"templateId": templateId,
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if nodes, ok := d.GetOk("nodes"); ok {
		payload["nodes"] = nodeSetToStringList(nodes.(*schema.Set))
	}

	if memberInterfaces, ok := d.GetOk("interfaces"); ok {
		payload["interfaces"] = physicalInterfaceSetToString(memberInterfaces.(*schema.Set))
	}

	interfacePolicyUUID := d.Get("interface_policy_uuid").(string)
	breakoutMode := d.Get("breakout_mode").(string)

	if interfacePolicyUUID == "" && breakoutMode == "" {
		return fmt.Errorf("Either 'interface_policy_uuid' or 'breakout_mode' must be specified for creating a Physical Interface")
	}

	if interfacePolicyUUID != "" {
		payload["policy"] = interfacePolicyUUID
		payload["policyGroupType"] = "physical"
	} else {
		payload["breakoutMode"] = breakoutMode
		payload["policyGroupType"] = "breakout"
	}

	if interfaceDescriptions, ok := d.GetOk("interface_descriptions"); ok {
		payload["interfaceDescriptions"] = getInterfaceDescriptionsPayloadPhysical(interfaceDescriptions.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/fabricResourceTemplate/template/interfaceProfiles/-", payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/PhysicalInterface/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Physical Interface Resource - Create Complete: %v", d.Id())
	return resourceMSOPhysicalInterfaceRead(d, m)
}

func resourceMSOPhysicalInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Interface Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "PhysicalInterface")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricResourceTemplate", "template", "interfaceProfiles")
	if err != nil {
		return err
	}

	err = setPhysicalInterfaceData(d, policy, templateId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Physical Interface Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOPhysicalInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Interface Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	if d.Get("interface_policy_uuid").(string) == "" && d.Get("breakout_mode").(string) == "" {
		return fmt.Errorf("Either 'interface_policy_uuid' or 'breakout_mode' must be specified for creating a Physical Interface")
	}

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "interfaceProfiles")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricResourceTemplate/template/interfaceProfiles/%d", policyIndex)

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

	if d.HasChange("nodes") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/nodes", updatePath), nodeSetToStringList(d.Get("nodes").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	if d.HasChange("interfaces") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaces", updatePath), physicalInterfaceSetToString(d.Get("interfaces").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	if d.HasChange("interface_policy_uuid") {
		interfacePolicyUUID := d.Get("interface_policy_uuid").(string)
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/policy", updatePath), interfacePolicyUUID)
		if err != nil {
			return err
		}
	}

	if d.HasChange("breakout_mode") {
		breakoutMode := d.Get("breakout_mode").(string)
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/breakoutMode", updatePath), breakoutMode)
		if err != nil {
			return err
		}
	}

	if d.HasChange("interface_descriptions") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaceDescriptions", updatePath), getInterfaceDescriptionsPayloadPhysical(d.Get("interface_descriptions").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/PhysicalInterface/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Physical Interface Resource - Update Complete: %v", d.Id())
	return resourceMSOPhysicalInterfaceRead(d, m)
}

func resourceMSOPhysicalInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Interface Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateContainer, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateContainer, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "interfaceProfiles")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricResourceTemplate/template/interfaceProfiles/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Physical Interface Resource - Delete Complete: %v", d.Id())
	return nil
}
