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

func resourceMSOVirtualPortChannelInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOVirtualPortChannelInterfaceCreate,
		Read:   resourceMSOVirtualPortChannelInterfaceRead,
		Update: resourceMSOVirtualPortChannelInterfaceUpdate,
		Delete: resourceMSOVirtualPortChannelInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOVirtualPortChannelInterfaceImport,
		},

		SchemaVersion: version,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Fabric Resource template ID.",
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
				Description:  "Virtual Port Channel Interface name.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual Port Channel Interface UUID.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Virtual Port Channel Interface description.",
			},
			"node_1": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"node_2": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"node_1_interfaces": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of interface IDs (or ranges) for node 1.",
			},
			"node_2_interfaces": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of interface IDs (or ranges) for node 2.",
			},
			"interface_policy_group_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "UUID of the Port Channel Interface Policy Group.",
			},
			"interface_descriptions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of interface descriptions; provided list replaces the existing list in NDO.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"interface": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func setVPCInterfaceData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	vpcCont, err := GetPolicyByName(response, policyName, "fabricResourceTemplate", "template", "virtualPortChannels")
	if err != nil {
		return err
	}
	name := models.StripQuotes(vpcCont.S("name").String())

	d.SetId(fmt.Sprintf("templateId/%s/virtualPortChannelInterface/%s", templateId, name))
	d.Set("template_id", templateId)

	d.Set("name", models.StripQuotes(vpcCont.S("name").String()))
	d.Set("description", models.StripQuotes(vpcCont.S("description").String()))
	d.Set("uuid", models.StripQuotes(vpcCont.S("uuid").String()))

	d.Set("interface_policy_group_uuid", models.StripQuotes(vpcCont.S("policy").String()))

	d.Set("node_1", models.StripQuotes(vpcCont.S("node1Details", "node").String()))
	d.Set("node_2", models.StripQuotes(vpcCont.S("node2Details", "node").String()))

	d.Set("node_1_interfaces", splitCommaString(models.StripQuotes(vpcCont.S("node1Details", "memberInterfaces").String())))
	d.Set("node_2_interfaces", splitCommaString(models.StripQuotes(vpcCont.S("node2Details", "memberInterfaces").String())))

	if vpcCont.Exists("interfaceDescriptions") {
		count, _ := vpcCont.ArrayCount("interfaceDescriptions")
		out := make([]any, 0, count)
		for i := 0; i < count; i++ {
			descCont, err := vpcCont.ArrayElement(i, "interfaceDescriptions")
			if err != nil {
				return err
			}
			entry := make(map[string]any)
			entry["node"] = models.StripQuotes(descCont.S("nodeID").String())
			entry["interface"] = models.StripQuotes(descCont.S("interfaceID").String())
			entry["description"] = models.StripQuotes(descCont.S("description").String())
			out = append(out, entry)
		}
		d.Set("interface_descriptions", out)
	}

	return nil
}

func resourceMSOVirtualPortChannelInterfaceImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO VPC Interface - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}
	name, err := GetPolicyNameFromResourceId(d.Id(), "virtualPortChannelInterface")
	if err != nil {
		return nil, err
	}

	err = setVPCInterfaceData(d, msoClient, templateId, name)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] MSO VPC Interface - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOVirtualPortChannelInterfaceCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VPC Interface - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	interfaces1 := getListOfStringsFromSchemaList(d, "node_1_interfaces")
	interfaces2 := getListOfStringsFromSchemaList(d, "node_2_interfaces")

	payload := map[string]any{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"node1Details": map[string]any{
			"node":             d.Get("node_1").(string),
			"memberInterfaces": strings.Join(interfaces1, ","),
		},
		"node2Details": map[string]any{
			"node":             d.Get("node_2").(string),
			"memberInterfaces": strings.Join(interfaces2, ","),
		},
		"policy":                d.Get("interface_policy_group_uuid").(string),
		"interfaceDescriptions": buildInterfaceDescriptionsPayload(d),
	}
	payloadModel := models.GetPatchPayload("add", "/fabricResourceTemplate/template/virtualPortChannels/-", payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	d.SetId(fmt.Sprintf("templateId/%s/virtualPortChannelInterface/%s", templateId, name))

	log.Printf("[DEBUG] MSO VPC Interface - Create Complete: %v", d.Id())
	return resourceMSOVirtualPortChannelInterfaceRead(d, m)
}

func resourceMSOVirtualPortChannelInterfaceRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VPC Interface - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setVPCInterfaceData(d, msoClient, templateId, policyName)

	log.Printf("[DEBUG] MSO VPC Interface - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOVirtualPortChannelInterfaceUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VPC Interface - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	index, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "virtualPortChannels")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricResourceTemplate/template/virtualPortChannels/%d", index)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string)); err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string)); err != nil {
			return err
		}
	}

	if d.HasChange("node_1") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/node1Details/node", updatePath), d.Get("node_1").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("node_1_interfaces") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/node1Details/memberInterfaces", updatePath), strings.Join(getListOfStringsFromSchemaList(d, "node_1_interfaces"), ",")); err != nil {
			return err
		}
	}

	if d.HasChange("node_2") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/node2Details/node", updatePath), d.Get("node_2").(string)); err != nil {
			return err
		}
	}

	if d.HasChange("node_2_interfaces") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/node2Details/memberInterfaces", updatePath), strings.Join(getListOfStringsFromSchemaList(d, "node_2_interfaces"), ",")); err != nil {
			return err
		}
	}

	if d.HasChange("interface_policy_group_uuid") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/policy", updatePath), d.Get("interface_policy_group_uuid").(string)); err != nil {
			return err
		}
	}

	if d.HasChange("interface_descriptions") {
		if err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaceDescriptions", updatePath), buildInterfaceDescriptionsPayload(d)); err != nil {
			return err
		}
	}

	if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/virtualPortChannelInterface/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO VPC Interface - Update Complete: %v", d.Id())
	return resourceMSOVirtualPortChannelInterfaceRead(d, m)
}

func resourceMSOVirtualPortChannelInterfaceDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VPC Interface - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	index, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricResourceTemplate", "template", "virtualPortChannels")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricResourceTemplate/template/virtualPortChannels/%d", index))
	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO VPC Interface - Delete Complete: %v", d.Id())
	return nil
}

func buildInterfaceDescriptionsPayload(d *schema.ResourceData) []map[string]any {
	raw := d.Get("interface_descriptions").([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, v := range raw {
		m := v.(map[string]any)
		out = append(out, map[string]any{
			"nodeID":      m["node"].(string),
			"interfaceID": m["interface"].(string),
			"description": m["description"].(string),
		})
	}
	return out
}
