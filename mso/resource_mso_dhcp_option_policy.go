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

func resourceMSODHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSODHCPOptionPolicyCreate,
		Update: resourceMSODHCPOptionPolicyUpdate,
		Read:   resourceMSODHCPOptionPolicyRead,
		Delete: resourceMSODHCPOptionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSODHCPOptionPolicyImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"option": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashResource(DHCPOptionSchema),
				Elem:     DHCPOptionSchema,
			},
		}),
	}
}

var DHCPOptionSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]+$`),
				"value should be alphanumeric",
			),
		},
		"id": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9]+$`),
				"value should be numeric",
			),
		},
		"data": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	},
}

func resourceMSODHCPOptionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	DHCPOptionPolicy, err := getDHCPOptionPolicy(msoClient, id)
	if err != nil {
		return nil, err
	}
	setDHCPOptionPolicy(DHCPOptionPolicy, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func getDHCPOptionPolicy(client *client.Client, id string) (*models.DHCPOptionPolicy, error) {
	cont, err := client.ReadDHCPOptionPolicy(id)
	if err != nil {
		return nil, err
	}

	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(cont)
	if err != nil {
		return nil, err
	}

	return DHCPOptionPolicy, nil
}

func setDHCPOptionPolicy(DHCPOptionPolicy *models.DHCPOptionPolicy, d *schema.ResourceData) {
	d.Set("description", DHCPOptionPolicy.Desc)
	d.Set("name", DHCPOptionPolicy.Name)
	tfOptionList := make([]map[string]string, 0)
	if _, ok := d.GetOk("tenant_id"); ok {
		optionList := d.Get("option").(*schema.Set).List()
		for _, option := range optionList {
			optionMap := option.(map[string]interface{})
			for _, remoteOption := range DHCPOptionPolicy.DHCPOption {
				if optionMap["name"].(string) != remoteOption.Name {
					continue
				}
				tfOptionList = append(tfOptionList, map[string]string{
					"name": remoteOption.Name,
					"data": remoteOption.Data,
					"id":   remoteOption.ID,
				})
			}
		}
	} else {
		for _, remoteOption := range DHCPOptionPolicy.DHCPOption {
			tfOptionList = append(tfOptionList, map[string]string{
				"name": remoteOption.Name,
				"data": remoteOption.Data,
				"id":   remoteOption.ID,
			})
		}
	}
	d.Set("tenant_id", DHCPOptionPolicy.TenantID)
	d.Set("option", tfOptionList)
	d.SetId(DHCPOptionPolicy.ID)
}

func resourceMSODHCPOptionPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())

	msoClient := m.(*client.Client)

	DHCPOptionPolicy := models.DHCPOptionPolicy{
		TenantID: d.Get("tenant_id").(string),
		Name:     d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		DHCPOptionPolicy.Desc = desc.(string)
	}

	if optionList, ok := d.GetOk("option"); ok {
		optionModelList := make([]models.DHCPOption, 0)
		for _, option := range optionList.(*schema.Set).List() {
			optionMap := option.(map[string]interface{})
			optionModelList = append(optionModelList, models.DHCPOption{
				Name: optionMap["name"].(string),
				ID:   optionMap["id"].(string),
				Data: optionMap["data"].(string),
			})
		}
		DHCPOptionPolicy.DHCPOption = optionModelList
	}

	cont, err := msoClient.CreateDHCPOptionPolicy(&DHCPOptionPolicy)
	if err != nil {
		return err
	}
	d.SetId(models.StripQuotes(cont.S("id").String()))

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSODHCPOptionPolicyRead(d, m)
}

func resourceMSODHCPOptionPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Policy Update", d.Id())

	msoClient := m.(*client.Client)

	DHCPOptionPolicy := models.DHCPOptionPolicy{
		TenantID: d.Get("tenant_id").(string),
		Name:     d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		DHCPOptionPolicy.Desc = desc.(string)
	}

	if d.HasChange("option") {
		oldOptions, newOptions := d.GetChange("option")
		oldOptionSet := oldOptions.(*schema.Set)
		newOptionSet := newOptions.(*schema.Set)

		optionsToDel := oldOptionSet.Difference(newOptionSet).List()
		optionsToCreate := newOptionSet.Difference(oldOptionSet).List()

		for i, delOption := range optionsToDel {
			for _, crOption := range optionsToCreate {
				if crOption.(map[string]interface{})["name"].(string) == delOption.(map[string]interface{})["name"].(string) {
					fmt.Printf("optionsToDel bfr: %v\n", optionsToDel)
					optionsToDel = append(optionsToDel[:i], optionsToDel[i+1:]...)
					fmt.Printf("optionsToDel aft: %v\n", optionsToDel)
				}
			}
		}

		optionModelList := make([]models.DHCPOption, 0)

		for _, option := range optionsToDel {
			optionMap := option.(map[string]interface{})
			optionModelList = append(optionModelList, models.DHCPOption{
				Name: optionMap["name"].(string),
				Data: optionMap["data"].(string),
				ID:   "remove",
			})
		}

		for _, option := range newOptionSet.List() {
			optionMap := option.(map[string]interface{})
			optionModelList = append(optionModelList, models.DHCPOption{
				Name: optionMap["name"].(string),
				Data: optionMap["data"].(string),
				ID:   optionMap["id"].(string),
			})
		}

		fmt.Printf("optionModelList: %v\n", optionModelList)
		DHCPOptionPolicy.DHCPOption = optionModelList
	}

	_, err := msoClient.UpdateDHCPOptionPolicy(d.Id(), &DHCPOptionPolicy)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Policy Update finished successfully: %s", d.Id())

	return resourceMSODHCPOptionPolicyRead(d, m)
}

func resourceMSODHCPOptionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	id := d.Id()
	log.Printf("id: %v\n", id)
	DHCPOptionPolicy, err := getDHCPOptionPolicy(msoClient, id)
	if err != nil {
		d.SetId("")
		return err
	}
	setDHCPOptionPolicy(DHCPOptionPolicy, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSODHCPOptionPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()

	err := msoClient.DeleteDHCPOptionPolicy(id)
	if err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
