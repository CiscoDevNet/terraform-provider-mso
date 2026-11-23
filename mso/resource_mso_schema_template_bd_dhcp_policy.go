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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The DHCP Relay Policy name.",
			},
			"dhcp_option_policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DHCP Option Policy name.",
			},
		},
	}
}

func getTemplateBDDHCPPolicyRef(c *client.Client, schemaID, templateName, relayName, optionName string) (map[string]interface{}, error) {
	tenantID, err := c.GetTenantIDFromSchemaTemplate(schemaID, templateName)
	if err != nil {
		return nil, err
	}
	refs, err := c.GetDHCPPoliciesUUIDByName(tenantID, []interface{}{map[string]interface{}{
		"relayName":  relayName,
		"optionName": optionName,
	}})
	if err != nil {
		return nil, err
	}
	if len(refs) != 1 {
		return nil, fmt.Errorf("expected one DHCP policy reference, got %d", len(refs))
	}
	ref := refs[0].(map[string]interface{})
	if option, ok := ref["dhcpOptionLabel"].(map[string]interface{}); ok && option["ref"] == "" {
		delete(ref, "dhcpOptionLabel")
	}
	return ref, nil
}

func templateBDDHCPPolicyPatch(c *client.Client, d *schema.ResourceData, op string) error {
	versionCheck, err := c.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}
	if versionCheck != -1 {
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
			d.Get("name").(string),
			d.Get("dhcp_option_policy_name").(string),
		)
		if err != nil {
			return err
		}
		value = ref
		if op == "replace" {
			index, err := findTemplateBDDHCPPolicyIndex(c, d, ref["ref"].(string))
			if err != nil {
				return err
			}
			path = fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%d", d.Get("template_name"), d.Get("bd_name"), index)
		}
	} else {
		ref, err := getTemplateBDDHCPPolicyRef(c,
			d.Get("schema_id").(string),
			d.Get("template_name").(string),
			d.Get("name").(string),
			"",
		)
		if err != nil {
			return err
		}
		index, err := findTemplateBDDHCPPolicyIndex(c, d, ref["ref"].(string))
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

func findTemplateBDDHCPPolicyIndex(c *client.Client, d *schema.ResourceData, relayRef string) (int, error) {
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
	if err := templateBDDHCPPolicyPatch(c, d, "add"); err != nil {
		return err
	}
	d.SetId(templateBDDHCPPolicyID(d))
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	return readTemplateBDDHCPPolicy(d, m)
}

func readTemplateBDDHCPPolicy(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	response, err := c.ReadTemplateBDDHCPPolicy(d.Get("schema_id").(string))
	if err != nil {
		return err
	}
	versionCheck, err := c.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}
	if versionCheck != -1 {
		return fmt.Errorf("mso_schema_template_bd_dhcp_policy requires NDO 4.0 or later")
	}
	template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("template_name").(string)
	})
	if err != nil {
		return err
	}
	bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
		return models.G(item, "name") == d.Get("bd_name").(string)
	})
	if err != nil {
		return err
	}
	count, err := bd.ArrayCount("dhcpLabels")
	if err != nil {
		return err
	}
	nameOptionList, err := getDHCPPolicesNameByRef(count, d.Get("schema_id").(string), d.Get("template_name").(string), bd, c)
	if err != nil {
		return err
	}
	found := false
	for _, item := range nameOptionList {
		policy := item.(map[string]interface{})
		if policy["name"] != d.Get("name").(string) {
			continue
		}
		found = true
		if err := d.Set("dhcp_option_policy_name", policy["dhcp_option_policy_name"]); err != nil {
			return err
		}
		break
	}
	if !found {
		d.SetId("")
		return nil
	}
	d.SetId(templateBDDHCPPolicyID(d))
	return nil
}

func resourceMSOTemplateBDDHCPPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD DHCP Policy: Beginning Update")
	if err := templateBDDHCPPolicyPatch(m.(*client.Client), d, "replace"); err != nil {
		return err
	}
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD DHCP Policy: Beginning Deletion")
	if err := templateBDDHCPPolicyPatch(m.(*client.Client), d, "remove"); err != nil {
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

func templateBDDHCPPolicyID(d *schema.ResourceData) string {
	return fmt.Sprintf("%s/%s/%s/%s", d.Get("schema_id"), d.Get("template_name"), d.Get("bd_name"), d.Get("name"))
}
