package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getDHCPPolicyUUIDByName(c *client.Client, tenantID, policyType, policyName string) (string, error) {
	policies, err := c.GetPoliciesByTenantID(policyType, tenantID)
	if err != nil {
		return "", err
	}

	var commonPolicyUUID string
	for _, item := range policies.Data().([]interface{}) {
		policy, ok := item.(map[string]interface{})
		if !ok || policy["name"] != policyName {
			continue
		}

		policyUUID, _ := policy["uuid"].(string)
		if policy["tenantId"] == tenantID {
			if policyUUID == "" {
				return "", fmt.Errorf("%s policy %q has no UUID", policyType, policyName)
			}
			return policyUUID, nil
		}
		if policy["tenantName"] == "common" {
			commonPolicyUUID = policyUUID
		}
	}

	if commonPolicyUUID != "" {
		return commonPolicyUUID, nil
	}
	return "", fmt.Errorf("%s policy %q not found", policyType, policyName)
}

func setDHCPPolicyUUIDsByName(c *client.Client, schemaID, templateName string, policies []interface{}) error {
	tenantID, err := c.GetTenantIDFromSchemaTemplate(schemaID, templateName)
	if err != nil {
		return err
	}
	for _, item := range policies {
		policy := item.(map[string]interface{})
		name, _ := policy["name"].(string)
		uuid, err := getDHCPPolicyUUIDByName(c, tenantID, "dhcpRelay", name)
		if err != nil {
			return err
		}
		policy["uuid"] = uuid

		optionName, _ := policy["dhcp_option_policy_name"].(string)
		if optionName == "" || optionName == "{}" {
			policy["dhcp_option_policy_uuid"] = ""
			continue
		}
		optionUUID, err := getDHCPPolicyUUIDByName(c, tenantID, "dhcpOption", optionName)
		if err != nil {
			return err
		}
		policy["dhcp_option_policy_uuid"] = optionUUID
	}
	return nil
}

func resourceMSOTemplateBDDHCPPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateBDDHCPPolicyCreate,
		Read:   resourceMSOTemplateBDDHCPPolicyRead,
		Update: resourceMSOTemplateBDDHCPPolicyUpdate,
		Delete: resourceMSOTemplateBDDHCPPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateBDDHCPPolicyImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the schema containing the Bridge Domain.",
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the template containing the Bridge Domain.",
			},
			"bd_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Bridge Domain.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The DHCP Relay Policy name.",
				ExactlyOneOf: []string{"uuid"},
			},
			"uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The DHCP Relay Policy UUID.",
				ExactlyOneOf: []string{"name"},
			},
			"dhcp_option_policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The DHCP Option Policy name.",
				ConflictsWith: []string{"dhcp_option_policy_uuid"},
			},
			"dhcp_option_policy_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(1, 64),
				Description:   "The DHCP Option Policy UUID.",
				ConflictsWith: []string{"dhcp_option_policy_name"},
			},
		},
	}
}

func getTemplateBDDHCPPolicyInput(d *schema.ResourceData, key string) string {
	// Use raw config to distinguish an omitted Optional+Computed value from its
	// value previously populated in state.
	raw := d.GetRawConfig()
	if !raw.IsNull() && raw.IsKnown() {
		value := raw.GetAttr(key)
		if value.IsNull() || !value.IsKnown() {
			return ""
		}
		return value.AsString()
	}
	return d.Get(key).(string)
}

func getTemplateBDDHCPPolicyRef(c *client.Client, schemaID, templateName, relayName, relayUUID, optionName, optionUUID string) (map[string]interface{}, error) {
	tenantID, err := c.GetTenantIDFromSchemaTemplate(schemaID, templateName)
	if err != nil {
		return nil, err
	}
	if relayUUID == "" {
		relayUUID, err = getDHCPPolicyUUIDByName(c, tenantID, "dhcpRelay", relayName)
		if err != nil {
			return nil, err
		}
	}

	ref := map[string]interface{}{"ref": relayUUID}
	if optionUUID == "" && optionName != "" {
		optionUUID, err = getDHCPPolicyUUIDByName(c, tenantID, "dhcpOption", optionName)
		if err != nil {
			return nil, err
		}
	}
	if optionUUID != "" {
		ref["dhcpOptionLabel"] = map[string]interface{}{"ref": optionUUID}
	}
	return ref, nil
}

func setTemplateBDDHCPPolicy(c *client.Client, d *schema.ResourceData, op string) error {
	versionCheck, err := c.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}
	if versionCheck > 0 {
		return fmt.Errorf("mso_schema_template_bd_dhcp_policy requires NDO 4.0 or later")
	}
	path := fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/-", d.Get("template_name"), d.Get("bd_name"))
	payload := container.New()
	payload.Array()
	var value interface{}
	if op != "remove" {
		ref, err := getTemplateBDDHCPPolicyRef(c,
			d.Get("schema_id").(string),
			d.Get("template_name").(string),
			getTemplateBDDHCPPolicyInput(d, "name"),
			getTemplateBDDHCPPolicyInput(d, "uuid"),
			getTemplateBDDHCPPolicyInput(d, "dhcp_option_policy_name"),
			getTemplateBDDHCPPolicyInput(d, "dhcp_option_policy_uuid"),
		)
		if err != nil {
			return err
		}
		value = ref
		if op == "replace" {
			index, err := getTemplateBDDHCPPolicyIndex(c, d, ref["ref"].(string))
			if err != nil {
				return err
			}
			path = fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%d", d.Get("template_name"), d.Get("bd_name"), index)
		}
	} else {
		ref, err := getTemplateBDDHCPPolicyRef(c,
			d.Get("schema_id").(string),
			d.Get("template_name").(string),
			getTemplateBDDHCPPolicyInput(d, "name"),
			getTemplateBDDHCPPolicyInput(d, "uuid"),
			"",
			"",
		)
		if err != nil {
			return err
		}
		index, err := getTemplateBDDHCPPolicyIndex(c, d, ref["ref"].(string))
		if err != nil {
			return err
		}
		path = fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%d", d.Get("template_name"), d.Get("bd_name"), index)
	}
	if err := addPatchPayloadToContainer(payload, op, path, value); err != nil {
		return err
	}
	return doPatchRequest(c, "api/v1/schemas/"+d.Get("schema_id").(string), payload)
}

func getTemplateBDDHCPPolicyIndex(c *client.Client, d *schema.ResourceData, relayRef string) (int, error) {
	response, err := c.ReadTemplateBDDHCPPolicy(d.Get("schema_id").(string))
	if err != nil {
		return -1, err
	}
	template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("template_name").(string)
	})
	if err != nil {
		return -1, err
	}
	bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("bd_name").(string)
	})
	if err != nil {
		return -1, err
	}
	count, err := bd.ArrayCount("dhcpLabels")
	if err != nil {
		return -1, err
	}
	for index := 0; index < count; index++ {
		label, err := bd.ArrayElement(index, "dhcpLabels")
		if err != nil {
			return -1, err
		}
		if models.StripQuotes(label.S("ref").String()) == relayRef {
			return index, nil
		}
	}
	return -1, fmt.Errorf("DHCP policy %q not found on Bridge Domain %q", d.Get("name"), d.Get("bd_name"))
}

func resourceMSOTemplateBDDHCPPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD DHCP Policy: Beginning Creation")
	c := m.(*client.Client)
	if err := setTemplateBDDHCPPolicy(c, d, "add"); err != nil {
		return err
	}
	d.SetId(getTemplateBDDHCPPolicyID(d))
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	policy, err := getTemplateBDDHCPPolicy(d, c)
	if err != nil {
		return err
	}
	if policy == nil {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", policy["name"]); err != nil {
		return err
	}
	if err := d.Set("uuid", policy["uuid"]); err != nil {
		return err
	}
	for _, key := range []string{"dhcp_option_policy_name", "dhcp_option_policy_uuid"} {
		if err := setTemplateBDDHCPPolicyOptionState(d, key, policy); err != nil {
			return err
		}
	}
	d.SetId(getTemplateBDDHCPPolicyID(d))
	return nil
}

// Keep the option association returned by NDO in state so refresh can detect
// changes made outside Terraform. Preserve the configured identifier form;
// if an association exists but neither identifier is configured, record its
// name so Terraform can plan to reconcile it.
func setTemplateBDDHCPPolicyOptionState(d *schema.ResourceData, key string, policy map[string]interface{}) error {
	configuredUUID := getTemplateBDDHCPPolicyInput(d, "dhcp_option_policy_uuid")
	if key == "dhcp_option_policy_uuid" {
		if configuredUUID == "" {
			return d.Set(key, "")
		}
		value, _ := policy[key].(string)
		return d.Set(key, value)
	}
	if configuredUUID != "" {
		return d.Set(key, "")
	}
	value, _ := policy[key].(string)
	return d.Set(key, value)
}

// getTemplateBDDHCPPolicy resolves the configured name or UUID against the
// DHCP associations currently attached to the Bridge Domain.
func getTemplateBDDHCPPolicy(d *schema.ResourceData, c *client.Client) (map[string]interface{}, error) {
	response, err := c.ReadTemplateBDDHCPPolicy(d.Get("schema_id").(string))
	if err != nil {
		return nil, err
	}
	versionCheck, err := c.CompareVersion("4.0.0.0")
	if err != nil {
		return nil, err
	}
	if versionCheck > 0 {
		return nil, fmt.Errorf("mso_schema_template_bd_dhcp_policy requires NDO 4.0 or later")
	}
	template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("template_name").(string)
	})
	if err != nil {
		return nil, err
	}
	bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("bd_name").(string)
	})
	if err != nil {
		return nil, err
	}
	count, err := bd.ArrayCount("dhcpLabels")
	if err != nil {
		return nil, err
	}
	policies, err := getDHCPPolicesNameByRef(count, d.Get("schema_id").(string), d.Get("template_name").(string), bd, c)
	if err != nil {
		return nil, err
	}
	configuredUUID := getTemplateBDDHCPPolicyInput(d, "uuid")
	configuredName := getTemplateBDDHCPPolicyInput(d, "name")
	for _, item := range policies {
		policy := item.(map[string]interface{})
		if configuredUUID != "" && policy["uuid"] == configuredUUID || configuredUUID == "" && policy["name"] == configuredName {
			return policy, nil
		}
	}
	return nil, nil
}

func resourceMSOTemplateBDDHCPPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD DHCP Policy: Beginning Update")
	if err := setTemplateBDDHCPPolicy(m.(*client.Client), d, "replace"); err != nil {
		return err
	}
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD DHCP Policy: Beginning Deletion")
	if err := setTemplateBDDHCPPolicy(m.(*client.Client), d, "remove"); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceMSOTemplateBDDHCPPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	var schemaID, templateName, bdName, policyName string
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("expected import ID schema_id/template_name/bd_name/policy_name")
	}
	schemaID, templateName, bdName, policyName = parts[0], parts[1], parts[2], parts[3]
	for key, value := range map[string]string{"schema_id": schemaID, "template_name": templateName, "bd_name": bdName, "name": policyName} {
		if err := d.Set(key, value); err != nil {
			return nil, err
		}
	}
	if err := resourceMSOTemplateBDDHCPPolicyRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("DHCP policy %q not found", policyName)
	}
	return []*schema.ResourceData{d}, nil
}

func getTemplateBDDHCPPolicyID(d *schema.ResourceData) string {
	return fmt.Sprintf("%s/%s/%s/%s", d.Get("schema_id"), d.Get("template_name"), d.Get("bd_name"), d.Get("name"))
}
