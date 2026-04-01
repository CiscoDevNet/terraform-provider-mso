package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTenantPoliciesNetflowMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTenantPoliciesNetflowMonitorCreate,
		Read:   resourceMSOTenantPoliciesNetflowMonitorRead,
		Update: resourceMSOTenantPoliciesNetflowMonitorUpdate,
		Delete: resourceMSOTenantPoliciesNetflowMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOTenantPoliciesNetflowMonitorImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the NetFlow Monitor.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NetFlow Monitor.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Monitor.",
			},
			"netflow_record_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the NetFlow Record.",
			},
			"netflow_exporter_uuids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of NetFlow Exporter UUIDs.",
			},
		},
	}
}

func buildExporterUUIDsPayload(exportersRaw interface{}) []string {
	exportersList := exportersRaw.(*schema.Set).List()
	exporters := make([]string, 0, len(exportersList))
	for _, item := range exportersList {
		exporters = append(exporters, item.(string))
	}
	return exporters
}

func setNetflowMonitorData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/NetflowMonitor/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	if response.Exists("description") {
		d.Set("description", models.StripQuotes(response.S("description").String()))
	}
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	if response.Exists("recordRef") {
		d.Set("netflow_record_uuid", models.StripQuotes(response.S("recordRef").String()))
	}

	if response.Exists("exporterRefs") {
		exporterCount, _ := response.ArrayCount("exporterRefs")
		exporters := make([]string, exporterCount)
		for i := range exporterCount {
			exporters[i] = models.StripQuotes(response.S("exporterRefs").Index(i).String())
		}
		d.Set("netflow_exporter_uuids", exporters)
	}

	return nil
}

func resourceMSOTenantPoliciesNetflowMonitorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Beginning Import: %v", d.Id())
	err := resourceMSOTenantPoliciesNetflowMonitorRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTenantPoliciesNetflowMonitorCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Beginning Create: %v", d.Id())

	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"exporterRefs": buildExporterUUIDsPayload(d.Get("netflow_exporter_uuids")),
	}

	if recordUUID, ok := d.GetOk("netflow_record_uuid"); ok {
		payload["recordRef"] = recordUUID.(string)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/netFlowMonitors/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowMonitor/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Create Complete: %v", d.Id())
	return resourceMSOTenantPoliciesNetflowMonitorRead(d, m)
}

func resourceMSOTenantPoliciesNetflowMonitorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	monitorName, err := GetPolicyNameFromResourceId(d.Id(), "NetflowMonitor")
	if err != nil {
		return err
	}

	monitor, err := GetPolicyByName(response, monitorName, "tenantPolicyTemplate", "template", "netFlowMonitors")
	if err != nil {
		return err
	}

	setNetflowMonitorData(d, monitor, templateId)
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOTenantPoliciesNetflowMonitorUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Beginning Update: %v", d.Id())

	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	monitorIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowMonitors")
	if err != nil {
		return err
	}

	payloadCont := container.New()
	payloadCont.Array()

	monitorPath := fmt.Sprintf("/tenantPolicyTemplate/template/netFlowMonitors/%d", monitorIndex)

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", monitorPath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", monitorPath), d.Get("description").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("netflow_record_uuid") {
		recordUUID := d.Get("netflow_record_uuid").(string)
		if recordUUID != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/recordRef", monitorPath), recordUUID)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/recordRef", monitorPath), nil)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("netflow_exporter_uuids") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/exporterRefs", monitorPath), buildExporterUUIDsPayload(d.Get("netflow_exporter_uuids")))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowMonitor/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Update Complete: %v", d.Id())
	return resourceMSOTenantPoliciesNetflowMonitorRead(d, m)
}

func resourceMSOTenantPoliciesNetflowMonitorDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	monitorIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowMonitors")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/netFlowMonitors/%d", monitorIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO NetFlow Monitor Resource - Delete Complete: %v", d.Id())
	return nil
}
