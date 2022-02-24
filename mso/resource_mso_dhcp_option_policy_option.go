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

func resourceMSODHCPOptionPolicyOption() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSODHCPOptionPolicyOptionCreate,
		Update: resourceMSODHCPOptionPolicyOptionUpdate,
		Read:   resourceMSODHCPOptionPolicyOptionRead,
		Delete: resourceMSODHCPOptionPolicyOptionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSODHCPOptionPolicyOptionImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"option_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"option_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "value should be alphanumeric"),
			},
			"option_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "value should be alphanumeric"),
			},
			"option_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMSODHCPOptionPolicyOptionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] DHCP Option Policy Option: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	dhcpOptionPolicyRemote, err := msoClient.ReadDHCPOptionPolicyOption(id)
	if err != nil {
		return nil, err
	}
	setDHCPOptionPolicyOption(d, dhcpOptionPolicyRemote)
	d.SetId(id)
	log.Println("[DEBUG] DHCP Option Policy Option: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}

func setDHCPOptionPolicyOption(d *schema.ResourceData, m *models.DHCPOptionPolicyOption) {
	d.Set("option_name", m.Name)
	d.Set("option_id", m.ID)
	d.Set("option_data", m.Data)
	d.Set("option_policy_name", m.PolicyName)
}

func resourceMSODHCPOptionPolicyOptionRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] DHCP Option Policy Option: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	var dhcpOptionPolicyOptionRemote *models.DHCPOptionPolicyOption
	dhcpOptionPolicyOptionRemote, err := msoClient.ReadDHCPOptionPolicyOption(id)
	if err != nil {
		log.Printf("err: %v", err)
		d.SetId("")
		return nil
	}
	setDHCPOptionPolicyOption(d, dhcpOptionPolicyOptionRemote)
	d.SetId(DHCPOptionPolicyOptionModeltoId(dhcpOptionPolicyOptionRemote))
	log.Println("[DEBUG] DHCP Option Policy Option: Reading Completed", d.Id())
	return nil
}

func resourceMSODHCPOptionPolicyOptionCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCP Option Policy Option: Beginning Creation")
	msoClient := m.(*client.Client)
	DHCPOptionPolicyOptionModel := models.DHCPOptionPolicyOption{
		PolicyName: d.Get("option_policy_name").(string),
		Name:       d.Get("option_name").(string),
	}
	if optionData, ok := d.GetOk("option_data"); ok {
		DHCPOptionPolicyOptionModel.Data = optionData.(string)
	}
	if optionId, ok := d.GetOk("option_id"); ok {
		DHCPOptionPolicyOptionModel.ID = optionId.(string)
	}
	log.Printf("DHCPOptionPolicyOptionModel.PolicyName %v", DHCPOptionPolicyOptionModel.PolicyName)
	log.Printf("DHCPOptionPolicyOptionModel.Name %v", DHCPOptionPolicyOptionModel.Name)
	err := msoClient.CreateDHCPOptionPolicyOption(&DHCPOptionPolicyOptionModel)
	id := DHCPOptionPolicyOptionModeltoId(&DHCPOptionPolicyOptionModel)
	if err != nil {
		return err
	}
	d.SetId(id)
	log.Printf("[DEBUG] DHCP Option Policy Option: Creation Completed %s", d.Id())
	return resourceMSODHCPOptionPolicyOptionRead(d, m)
}

func resourceMSODHCPOptionPolicyOptionUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] DHCP Option Policy Option : Beginning Update", d.Id())
	msoClient := m.(*client.Client)
	optionPolicyName := d.Get("option_policy_name")
	optionName := d.Get("option_name")
	optionId := d.Get("option_id")
	optionData := d.Get("option_data")

	newPolicy := models.DHCPOptionPolicyOption{
		Name:       optionName.(string),
		ID:         optionId.(string),
		Data:       optionData.(string),
		PolicyName: optionPolicyName.(string),
	}
	err := msoClient.UpdateDHCPOptionPolicyOption(&newPolicy)
	if err != nil {
		return err
	}
	d.SetId(DHCPOptionPolicyOptionModeltoId(&newPolicy))
	log.Println("[DEBUG] DHCP Option Policy Option: Update Completed", d.Id())
	return resourceMSODHCPOptionPolicyOptionRead(d, m)
}

func resourceMSODHCPOptionPolicyOptionDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] DHCP Option Policy Option: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	err := msoClient.DeleteDHCPOptionPolicyOption(id)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] DHCP Option Policy Option: Destroy Completed", d.Id())
	d.SetId("")
	return nil
}

func DHCPOptionPolicyOptionModeltoId(m *models.DHCPOptionPolicyOption) string {
	return fmt.Sprintf("%s/%s", m.PolicyName, m.Name)
}
