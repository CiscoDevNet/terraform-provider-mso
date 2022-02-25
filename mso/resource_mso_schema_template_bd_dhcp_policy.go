package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dhcp_option_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				RequiredWith: []string{"dhcp_option_version"},
			},
			"dhcp_option_version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getMSOTemplateBDDHCPPolicy(msoClient *client.Client, obj *models.TemplateBDDHCPPolicy) (*models.TemplateBDDHCPPolicy, error) {
	cont, err := msoClient.ReadTemplateBDDHCPPolicy(obj.SchemaID)
	if err != nil {
		return nil, err
	}
	remotePolicy, err := models.TemplateBDDHCPPolicyFromContainer(cont, obj)
	if err != nil {
		return nil, err
	}
	return remotePolicy, nil
}

func setMSOTemplateBDDHCPPolicy(d *schema.ResourceData, obj *models.TemplateBDDHCPPolicy) {
	d.Set("schema_id", obj.SchemaID)
	d.Set("template_name", obj.TemplateName)
	d.Set("bd_name", obj.BDName)
	d.Set("name", obj.Name)
	d.Set("version", obj.Version)
	d.Set("dhcp_option_name", obj.DHCPOptionName)
	d.Set("dhcp_option_version", obj.DHCPOptionVersion)
}

func resourceMSOTemplateBDDHCPPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	bdDHCPPolicy := modelFromMSOTemplateBDDHCPPolicyId(id)

	remoteBDDHCPPolicy, err := getMSOTemplateBDDHCPPolicy(msoClient, bdDHCPPolicy)
	if err != nil {
		return nil, err
	}
	setMSOTemplateBDDHCPPolicy(d, remoteBDDHCPPolicy)

	d.SetId(createMSOTemplateBDDHCPPolicyId(bdDHCPPolicy))

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOTemplateBDDHCPPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Creation")
	msoClient := m.(*client.Client)

	bdDHCPPolicy := models.TemplateBDDHCPPolicy{
		Name:         d.Get("name").(string),
		SchemaID:     d.Get("schema_id").(string),
		TemplateName: d.Get("template_name").(string),
		BDName:       d.Get("bd_name").(string),
	}

	if version, ok := d.GetOk("version"); ok {
		bdDHCPPolicy.Version = version.(int)
	}

	if dhcpName, ok := d.GetOk("dhcp_option_name"); ok {
		bdDHCPPolicy.DHCPOptionName = dhcpName.(string)
	}

	if dhcpVersion, ok := d.GetOk("dhcp_option_version"); ok {
		bdDHCPPolicy.DHCPOptionVersion = dhcpVersion.(int)
	}

	_, err := msoClient.CreateTemplateBDDHCPPolicy(&bdDHCPPolicy)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Creation Complete")
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	bdDHCPPolicy := models.TemplateBDDHCPPolicy{
		Name:         d.Get("name").(string),
		SchemaID:     d.Get("schema_id").(string),
		TemplateName: d.Get("template_name").(string),
		BDName:       d.Get("bd_name").(string),
	}

	if version, ok := d.GetOk("version"); ok {
		bdDHCPPolicy.Version = version.(int)
	}

	if dhcpName, ok := d.GetOk("dhcp_option_name"); ok {
		bdDHCPPolicy.DHCPOptionName = dhcpName.(string)
	}

	if dhcpVersion, ok := d.GetOk("dhcp_option_version"); ok {
		bdDHCPPolicy.DHCPOptionVersion = dhcpVersion.(int)
	}

	_, err := msoClient.UpdateTemplateBDDHCPPolicy(&bdDHCPPolicy)
	if err != nil {
		return err
	}
	return resourceMSOTemplateBDDHCPPolicyRead(d, m)
}

func resourceMSOTemplateBDDHCPPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	bdDHCPPolicy := models.TemplateBDDHCPPolicy{
		Name:         d.Get("name").(string),
		SchemaID:     d.Get("schema_id").(string),
		TemplateName: d.Get("template_name").(string),
		BDName:       d.Get("bd_name").(string),
	}

	remoteBDDHCPPolicy, err := getMSOTemplateBDDHCPPolicy(msoClient, &bdDHCPPolicy)
	if err != nil {
		d.SetId("")
		return nil
	}
	setMSOTemplateBDDHCPPolicy(d, remoteBDDHCPPolicy)

	d.SetId(
		createMSOTemplateBDDHCPPolicyId(&bdDHCPPolicy),
	)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateBDDHCPPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)
	id := d.Id()
	model := modelFromMSOTemplateBDDHCPPolicyId(id)

	_, err := msoClient.DeleteTemplateBDDHCPPolicy(model)
	if err != nil {
		return err
	}

	return nil
}

func createMSOTemplateBDDHCPPolicyId(obj *models.TemplateBDDHCPPolicy) string {
	return fmt.Sprintf("/schemas/%s/templates/%s/bds/%s/dhcpLabels/%s",
		obj.SchemaID,
		obj.TemplateName,
		obj.BDName,
		obj.Name,
	)
}

func modelFromMSOTemplateBDDHCPPolicyId(id string) *models.TemplateBDDHCPPolicy {
	re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)/dhcpLabels/(.*)")
	match := re.FindStringSubmatch(id)
	return &models.TemplateBDDHCPPolicy{
		SchemaID:     match[1],
		TemplateName: match[2],
		BDName:       match[3],
		Name:         match[4],
	}
}
