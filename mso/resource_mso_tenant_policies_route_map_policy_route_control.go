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

func resourceMSORouteMapPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORouteMapPolicyCreate,
		Read:   resourceMSORouteMapPolicyRead,
		Update: resourceMSORouteMapPolicyUpdate,
		Delete: resourceMSORouteMapPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSORouteMapPolicyImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Route Map Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Route Map Policy.",
			},
			"match_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "routing_policy_only",
				ValidateFunc: validation.StringInSlice([]string{"prefix_and_routing_policy", "routing_policy_only"}, false),
				Description:  "The match type of the Route Map Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Route Map Policy.",
			},
		},
	}
}

func setRouteMapPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	name, ok := response.S("name").Data().(string)
	if !ok || name == "" {
		return fmt.Errorf("Route Map Policy name is missing or is not a string")
	}
	uuid, ok := response.S("uuid").Data().(string)
	if !ok || uuid == "" {
		return fmt.Errorf("Route Map Policy UUID is missing or is not a string")
	}
	description, _ := response.S("description").Data().(string)
	matchType := "routing_policy_only"
	if value := response.S("routeMapType").Data(); value != nil {
		routeMapType, ok := value.(string)
		if !ok {
			return fmt.Errorf("Route Map Policy routeMapType is not a string")
		}
		switch routeMapType {
		case "", "global":
		case "combinable":
			matchType = "prefix_and_routing_policy"
		default:
			return fmt.Errorf("unsupported Route Map Policy routeMapType %q", routeMapType)
		}
	}
	if err := d.Set("template_id", templateId); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("description", description); err != nil {
		return err
	}
	if err := d.Set("match_type", matchType); err != nil {
		return err
	}
	if err := d.Set("uuid", uuid); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, name))

	return nil
}

func resourceMSORouteMapPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Import: %v", d.Id())
	importID := d.Id()
	if err := resourceMSORouteMapPolicyRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("cannot import mso_tenant_policies_route_map_policy_route_control %q: resource not found", importID)
	}
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSORouteMapPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)
	routeMapType := "global"
	if d.Get("match_type").(string) == "prefix_and_routing_policy" {
		routeMapType = "combinable"
	}

	payload := map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"routeMapType": routeMapType,
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/routeMapPolicies/-", payload)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Create Complete: %v", d.Id())
	return resourceMSORouteMapPolicyRead(d, m)
}

func resourceMSORouteMapPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "RouteMapPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if err := setRouteMapPolicyData(d, policy, templateId); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSORouteMapPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("match_type") {
		routeMapType := "global"
		if d.Get("match_type").(string) == "prefix_and_routing_policy" {
			routeMapType = "combinable"
		}
		err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/routeMapType", updatePath), routeMapType)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Update Complete: %v", d.Id())
	return resourceMSORouteMapPolicyRead(d, m)
}

func resourceMSORouteMapPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
