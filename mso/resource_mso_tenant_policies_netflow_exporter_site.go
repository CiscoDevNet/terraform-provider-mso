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

func resourceMSONetflowExporterSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSONetflowExporterSiteCreate,
		Read:   resourceMSONetflowExporterSiteRead,
		Update: resourceMSONetflowExporterSiteUpdate,
		Delete: resourceMSONetflowExporterSiteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSONetflowExporterSiteImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the tenant policy template.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the site on which to configure the NetFlow Exporter.",
			},
			"netflow_exporter_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the template-level NetFlow Exporter.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The site-specific description of the NetFlow Exporter.",
			},
			"source_ip_address_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"inbandMgmtIP", "oobMgmtIP", "ptep", "customSrcIP",
				}, false),
				Description: "Source IP address type.",
			},
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The custom source IP address or prefix.",
			},
			"destination_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The destination port of the NetFlow Exporter.",
			},
			"destination_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IPv4 or IPv6 destination address.",
			},
			"dscp": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
					"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6", "cs7",
					"expedited_forwarding", "unspecified", "voice_admit",
				}, false),
				Description: "The QoS DSCP value.",
			},
			"application_epg": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"external_epg"},
				Description:   "The application EPG associated with the NetFlow Exporter.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"tenant_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the tenant.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"anp_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the application profile.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"epg_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the EPG.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"vrf_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the VRF.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"vrf_tenant_name": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						Description:  "The name of the tenant that contains the VRF. If omitted, the EPG tenant is used.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
				}},
			},
			"external_epg": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"application_epg"},
				Description:   "The external EPG associated with the NetFlow Exporter.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"tenant_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the tenant.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"l3out_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the L3Out.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"external_epg_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the external EPG.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"vrf_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The name of the VRF.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
					"vrf_tenant_name": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						Description:  "The name of the tenant that contains the VRF. If omitted, the external EPG tenant is used.",
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
				}},
			},
		},
	}
}

func netflowExporterSitePayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{
		"ref":             d.Get("netflow_exporter_uuid").(string),
		"destinationPort": d.Get("destination_port").(string),
		"destinationIP":   d.Get("destination_ip").(string),
		"dscp":            d.Get("dscp").(string),
		"versionFormat":   "v9",
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}
	if sourceIPType, ok := d.GetOk("source_ip_address_type"); ok {
		payload["srcIPType"] = sourceIPType.(string)
	}
	if srcIP, ok := d.GetOk("source_ip"); ok {
		payload["srcIP"] = srcIP.(string)
	}

	if application, ok := d.GetOk("application_epg"); ok {
		block := application.([]interface{})[0].(map[string]interface{})
		tenant := block["tenant_name"].(string)
		vrfTenant := netflowExporterSiteVrfTenant(d, tenant, "application_epg")
		payload["epgDn"] = fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, block["anp_name"].(string), block["epg_name"].(string))
		payload["vrfDn"] = fmt.Sprintf("uni/tn-%s/ctx-%s", vrfTenant, block["vrf_name"].(string))
	}

	if external, ok := d.GetOk("external_epg"); ok {
		block := external.([]interface{})[0].(map[string]interface{})
		tenant := block["tenant_name"].(string)
		vrfTenant := netflowExporterSiteVrfTenant(d, tenant, "external_epg")
		payload["externalEpgDn"] = fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", tenant, block["l3out_name"].(string), block["external_epg_name"].(string))
		payload["vrfDn"] = fmt.Sprintf("uni/tn-%s/ctx-%s", vrfTenant, block["vrf_name"].(string))
	}

	return payload
}

func netflowExporterSiteVrfTenant(d *schema.ResourceData, tenant, blockName string) string {
	vrfTenant := tenant
	// vrf_tenant_name is Optional and Computed. ResourceData.GetOk can therefore
	// return a previously computed state value even when the attribute was
	// removed from configuration. Inspect the raw configuration so omission
	// correctly falls back to the EPG or external EPG tenant.
	rawConfig := d.GetRawConfig()
	if rawConfig.IsKnown() && !rawConfig.IsNull() {
		rawBlock := rawConfig.GetAttr(blockName)
		if rawBlock.IsKnown() && !rawBlock.IsNull() && rawBlock.LengthInt() > 0 {
			rawVrfTenant := rawBlock.AsValueSlice()[0].GetAttr("vrf_tenant_name")
			if rawVrfTenant.IsKnown() && !rawVrfTenant.IsNull() && rawVrfTenant.AsString() != "" {
				vrfTenant = rawVrfTenant.AsString()
			}
		}
	}
	return vrfTenant
}

func resourceMSONetflowExporterSiteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Beginning Import: %v", d.Id())
	if err := resourceMSONetflowExporterSiteRead(d, m); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSONetflowExporterSiteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Beginning Create: %v", d.Id())
	if err := patchNetflowExporterSite(d, m); err != nil {
		return err
	}
	if err := resourceMSONetflowExporterSiteRead(d, m); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Create Complete: %v", d.Id())
	return nil
}

func resourceMSONetflowExporterSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID, siteID, exporterUUID, err := parseNetflowExporterSiteID(d.Id())
	if err != nil {
		return err
	}

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		d.SetId("")
		return nil
	}
	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteID, "tenantPolicyTemplate", "sites")
	if err != nil {
		d.SetId("")
		return nil
	}
	siteCont := templateCont.S("tenantPolicyTemplate", "sites").Index(siteIndex)
	exporterIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "ref", exporterUUID, "netFlowExporters")
	if err != nil {
		d.SetId("")
		return nil
	}

	if err := setNetflowExporterSiteData(d, siteCont.S("netFlowExporters").Index(exporterIndex), templateID, siteID, exporterUUID); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSONetflowExporterSiteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Beginning Update: %v", d.Id())
	if err := patchNetflowExporterSite(d, m); err != nil {
		return err
	}
	if err := resourceMSONetflowExporterSiteRead(d, m); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Update Complete: %v", d.Id())
	return nil
}

func patchNetflowExporterSite(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	siteID := d.Get("site_id").(string)
	exporterUUID := d.Get("netflow_exporter_uuid").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteID, "tenantPolicyTemplate", "sites")
	if err != nil {
		return fmt.Errorf("site %q is not present in template %q: %v", siteID, templateID, err)
	}

	exporterIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", exporterUUID, "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return fmt.Errorf("NetFlow Exporter %q does not exist in template %q: %v", exporterUUID, templateID, err)
	}

	siteCont := templateCont.S("tenantPolicyTemplate", "sites").Index(siteIndex)
	patchOperation := "add"
	path := fmt.Sprintf("/tenantPolicyTemplate/sites/%d/netFlowExporters/-", siteIndex)
	if _, err := GetPolicyIndexByKeyAndValue(siteCont, "ref", exporterUUID, "netFlowExporters"); err == nil {
		patchOperation = "replace"
		path = fmt.Sprintf("/tenantPolicyTemplate/sites/%d/netFlowExporters/%d", siteIndex, exporterIndex)
	}

	payloadCont := container.New()
	payloadCont.Array()
	if err := addPatchPayloadToContainer(payloadCont, patchOperation, path, netflowExporterSitePayload(d)); err != nil {
		return err
	}
	if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payloadCont); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/site/%s/NetflowExporter/%s", templateID, siteID, exporterUUID))
	log.Printf("[DEBUG] MSO NetFlow Exporter Site - Patch Complete: %v", d.Id())
	return nil
}

func resourceMSONetflowExporterSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	siteID := d.Get("site_id").(string)
	exporterUUID := d.Get("netflow_exporter_uuid").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		d.SetId("")
		return nil
	}
	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteID, "tenantPolicyTemplate", "sites")
	if err != nil {
		d.SetId("")
		return nil
	}
	siteCont := templateCont.S("tenantPolicyTemplate", "sites").Index(siteIndex)
	exporterIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "ref", exporterUUID, "netFlowExporters")
	if err != nil {
		d.SetId("")
		return nil
	}

	path := fmt.Sprintf("/tenantPolicyTemplate/sites/%d/netFlowExporters/%d", siteIndex, exporterIndex)
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Resource - Delete Complete")
	return nil
}

func parseNetflowExporterSiteID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 6 || parts[0] != "templateId" || parts[2] != "site" || parts[4] != "NetflowExporter" {
		return "", "", "", fmt.Errorf("invalid NetFlow Exporter site ID %q", id)
	}
	return parts[1], parts[3], parts[5], nil
}

func setNetflowExporterSiteData(d *schema.ResourceData, response *container.Container, templateID, siteID, exporterUUID string) error {
	d.SetId(fmt.Sprintf("templateId/%s/site/%s/NetflowExporter/%s", templateID, siteID, exporterUUID))
	d.Set("template_id", templateID)
	d.Set("site_id", siteID)
	d.Set("netflow_exporter_uuid", exporterUUID)
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("source_ip_address_type", models.StripQuotes(response.S("srcIPType").String()))
	d.Set("destination_port", models.StripQuotes(response.S("destinationPort").String()))
	d.Set("destination_ip", models.StripQuotes(response.S("destinationIP").String()))
	d.Set("dscp", models.StripQuotes(response.S("dscp").String()))
	if response.Exists("srcIP") {
		d.Set("source_ip", models.StripQuotes(response.S("srcIP").String()))
	} else {
		d.Set("source_ip", "")
	}

	if response.Exists("epgDn") {
		parts := strings.Split(models.StripQuotes(response.S("epgDn").String()), "/")
		if len(parts) != 4 || !strings.HasPrefix(parts[1], "tn-") || !strings.HasPrefix(parts[2], "ap-") || !strings.HasPrefix(parts[3], "epg-") {
			return fmt.Errorf("unexpected application EPG DN %q", response.S("epgDn").String())
		}
		vrfParts := strings.Split(models.StripQuotes(response.S("vrfDn").String()), "/")
		if len(vrfParts) != 3 || !strings.HasPrefix(vrfParts[1], "tn-") || !strings.HasPrefix(vrfParts[2], "ctx-") {
			return fmt.Errorf("unexpected VRF DN %q", response.S("vrfDn").String())
		}
		d.Set("application_epg", []interface{}{map[string]interface{}{
			"tenant_name":     strings.TrimPrefix(parts[1], "tn-"),
			"anp_name":        strings.TrimPrefix(parts[2], "ap-"),
			"epg_name":        strings.TrimPrefix(parts[3], "epg-"),
			"vrf_name":        strings.TrimPrefix(vrfParts[2], "ctx-"),
			"vrf_tenant_name": strings.TrimPrefix(vrfParts[1], "tn-"),
		}})
		d.Set("external_epg", nil)
	} else if response.Exists("externalEpgDn") {
		parts := strings.Split(models.StripQuotes(response.S("externalEpgDn").String()), "/")
		if len(parts) != 4 || !strings.HasPrefix(parts[1], "tn-") || !strings.HasPrefix(parts[2], "out-") || !strings.HasPrefix(parts[3], "instP-") {
			return fmt.Errorf("unexpected external EPG DN %q", response.S("externalEpgDn").String())
		}
		vrfParts := strings.Split(models.StripQuotes(response.S("vrfDn").String()), "/")
		if len(vrfParts) != 3 || !strings.HasPrefix(vrfParts[1], "tn-") || !strings.HasPrefix(vrfParts[2], "ctx-") {
			return fmt.Errorf("unexpected VRF DN %q", response.S("vrfDn").String())
		}
		d.Set("external_epg", []interface{}{map[string]interface{}{
			"tenant_name":       strings.TrimPrefix(parts[1], "tn-"),
			"l3out_name":        strings.TrimPrefix(parts[2], "out-"),
			"external_epg_name": strings.TrimPrefix(parts[3], "instP-"),
			"vrf_name":          strings.TrimPrefix(vrfParts[2], "ctx-"),
			"vrf_tenant_name":   strings.TrimPrefix(vrfParts[1], "tn-"),
		}})
		d.Set("application_epg", nil)
	} else {
		d.Set("application_epg", nil)
		d.Set("external_epg", nil)
	}
	return nil
}
