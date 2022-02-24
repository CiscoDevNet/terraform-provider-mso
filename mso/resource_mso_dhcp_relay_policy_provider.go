package mso

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var dhcpRelayPolicyMut sync.Mutex

func resourceMSODHCPRelayPolicyProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSODHCPRelayPolicyProviderCreate,
		Read:   resourceMSODHCPRelayPolicyProviderRead,
		Update: resourceMSODHCPRelayPolicyProviderUpdate,
		Delete: resourceMSODHCPRelayPolicyProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSODHCPRelayPolicyProviderImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"dhcp_relay_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"dhcp_server_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"epg_ref": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"external_epg_ref"},
				AtLeastOneOf:  []string{"epg_ref", "external_epg_ref"},
			},
			"external_epg_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMSODHCPRelayPolicyProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] DHCP Relay Policy Provider: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	dhcpRelayPolicyProviderModel, err := DHCPRelayPolicyProviderIdtoModel(id)
	if err != nil {
		return nil, err
	}
	dhcpRelayPolicyRemote, err := msoClient.ReadDHCPRelayPolicyProvider(dhcpRelayPolicyProviderModel)
	if err != nil {
		return nil, err
	}
	setDHCPRelayPolicyProviderAttr(d, dhcpRelayPolicyRemote)
	d.SetId(id)
	log.Println("[DEBUG] DHCP Relay Policy Provider: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSODHCPRelayPolicyProviderCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCP Relay Policy Provider: Beginning Creation")
	msoClient := m.(*client.Client)
	DHCPRelayPolicyProviderModel := models.DHCPRelayPolicyProvider{
		PolicyName: d.Get("dhcp_relay_policy_name").(string),
		Addr:       d.Get("dhcp_server_address").(string),
	}
	if epgRef, ok := d.GetOk("epg_ref"); ok {
		DHCPRelayPolicyProviderModel.EpgRef = epgRef.(string)
	}
	if extEpgRef, ok := d.GetOk("external_epg_ref"); ok {
		DHCPRelayPolicyProviderModel.ExternalEpgRef = extEpgRef.(string)
	}
	dhcpRelayPolicyMut.Lock()
	err := msoClient.CreateDHCPRelayPolicyProvider(&DHCPRelayPolicyProviderModel)
	id := DHCPRelayPolicyProviderModeltoId(&DHCPRelayPolicyProviderModel)
	if err != nil {
		return err
	}
	d.SetId(id)
	log.Printf("[DEBUG] DHCP Relay Policy Provider: Creation Completed %s", d.Id())
	err = resourceMSODHCPRelayPolicyProviderRead(d, m)
	dhcpRelayPolicyMut.Unlock()
	return err
}

func resourceMSODHCPRelayPolicyProviderRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] DHCP Relay Policy Provider: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	dhcpRelayPolicyProviderModel, err := DHCPRelayPolicyProviderIdtoModel(id)
	if err != nil {
		return err
	}
	var dhcpRelayPolicyProviderRemote *models.DHCPRelayPolicyProvider
	dhcpRelayPolicyProviderRemote, err = msoClient.ReadDHCPRelayPolicyProvider(dhcpRelayPolicyProviderModel)
	if err != nil {
		d.SetId("")
		return nil
	}
	setDHCPRelayPolicyProviderAttr(d, dhcpRelayPolicyProviderRemote)
	d.SetId(DHCPRelayPolicyProviderModeltoId(dhcpRelayPolicyProviderRemote))
	log.Println("[DEBUG] DHCP Relay Policy Provider: Reading Completed", d.Id())
	return nil
}

func resourceMSODHCPRelayPolicyProviderUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] DHCP Relay Policy Provider: Beginning Update", d.Id())
	msoClient := m.(*client.Client)
	oldPolicyName, newPolicyName := d.GetChange("dhcp_relay_policy_name")
	oldAddress, newAddress := d.GetChange("dhcp_server_address")
	oldEpg, newEpg := d.GetChange("epg_ref")
	oldExternalEpg, newExternalEpg := d.GetChange("external_epg_ref")
	oldPolicy := models.DHCPRelayPolicyProvider{
		PolicyName:     oldPolicyName.(string),
		Addr:           oldAddress.(string),
		EpgRef:         oldEpg.(string),
		ExternalEpgRef: oldExternalEpg.(string),
	}
	newPolicy := models.DHCPRelayPolicyProvider{
		PolicyName:     newPolicyName.(string),
		Addr:           newAddress.(string),
		EpgRef:         newEpg.(string),
		ExternalEpgRef: newExternalEpg.(string),
	}
	dhcpRelayPolicyMut.Lock()
	err := msoClient.UpdateDHCPRelayPolicyProvider(&newPolicy, &oldPolicy)
	if err != nil {
		return err
	}
	d.SetId(DHCPRelayPolicyProviderModeltoId(&newPolicy))
	log.Println("[DEBUG] DHCP Relay Policy Provider: Update Completed", d.Id())
	err = resourceMSODHCPRelayPolicyProviderRead(d, m)
	dhcpRelayPolicyMut.Unlock()
	return err
}

func resourceMSODHCPRelayPolicyProviderDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site L3out: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	dhcpRelayPolicyModel, err := DHCPRelayPolicyProviderIdtoModel(id)
	if err != nil {
		return err
	}
	dhcpRelayPolicyMut.Lock()
	err = msoClient.DeleteDHCPRelayPolicyProvider(dhcpRelayPolicyModel)
	if err != nil {
		return err
	}
	dhcpRelayPolicyMut.Unlock()
	log.Println("[DEBUG] Schema Site L3out: Beginning Destroy", d.Id())
	d.SetId("")
	return nil
}

func setDHCPRelayPolicyProviderAttr(d *schema.ResourceData, m *models.DHCPRelayPolicyProvider) {
	d.Set("dhcp_relay_policy_name", m.PolicyName)
	d.Set("dhcp_server_address", m.Addr)
	d.Set("epg_ref", m.EpgRef)
	d.Set("external_epg_ref", m.ExternalEpgRef)
}

func DHCPRelayPolicyProviderModeltoId(m *models.DHCPRelayPolicyProvider) string {
	if m.EpgRef != "" {
		return fmt.Sprintf("%s%s/%s", m.PolicyName, m.EpgRef, m.Addr)
	} else {
		return fmt.Sprintf("%s%s/%s", m.PolicyName, m.ExternalEpgRef, m.Addr)
	}
}

func DHCPRelayPolicyProviderIdtoModel(id string) (*models.DHCPRelayPolicyProvider, error) {
	idSplitted := strings.Split(id, "/")
	if len(idSplitted) == 10 && idSplitted[1] == "schemas" && idSplitted[3] == "templates" && idSplitted[5] == "anps" && idSplitted[7] == "epgs" {
		epgRef := "/"
		for i := 1; i <= 7; i++ {
			epgRef += idSplitted[i]
			epgRef += "/"
		}
		epgRef += idSplitted[8]
		provider := models.DHCPRelayPolicyProvider{
			Addr:       idSplitted[9],
			PolicyName: idSplitted[0],
			EpgRef:     epgRef,
		}
		return &provider, nil
	} else if len(idSplitted) == 8 && idSplitted[1] == "schemas" && idSplitted[3] == "templates" && idSplitted[5] == "externalEpgs" {
		externalEpgRef := "/"
		for i := 1; i <= 5; i++ {
			externalEpgRef += idSplitted[i]
			externalEpgRef += "/"
		}
		externalEpgRef += idSplitted[6]
		provider := models.DHCPRelayPolicyProvider{
			Addr:           idSplitted[7],
			PolicyName:     idSplitted[0],
			ExternalEpgRef: externalEpgRef,
		}
		return &provider, nil
	} else {
		return nil, fmt.Errorf("invalid DHCP Relay Policy Provider id format")
	}
}
