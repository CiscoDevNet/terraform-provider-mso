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

func resourceMSOIPSLAMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOIPSLAMonitoringPolicyCreate,
		Read:   resourceMSOIPSLAMonitoringPolicyRead,
		Update: resourceMSOIPSLAMonitoringPolicyUpdate,
		Delete: resourceMSOIPSLAMonitoringPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOIPSLAMonitoringPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sla_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"http", "tcp", "icmp", "l2ping",
				}, false),
			},
			"destination_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"http_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP10", "HTTP11",
				}, false),
			},
			"http_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_frequency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"detect_multiplier": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"request_data_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 17512),
			},
			"type_of_service": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"operation_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 604800000),
			},
			"threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 604800000),
			},
			"ipv6_traffic_class": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			slaType, _ := diff.GetOk("sla_type")
			if slaType == "icmp" || slaType == "l2ping" {
				diff.SetNew("destination_port", 0)
			} else if slaType == "http" {
				diff.SetNew("destination_port", 80)
			}
			return nil
		},
	}
}

func setIPSLAMonitoringPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {

	d.SetId(fmt.Sprintf("templateId/%s/IPSLAMonitoringPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("sla_type", models.StripQuotes(response.S("slaType").String()))
	d.Set("destination_port", response.S("slaPort").Data().(float64))
	d.Set("http_version", models.StripQuotes(response.S("httpVersion").String()))
	d.Set("http_uri", models.StripQuotes(response.S("httpUri").String()))
	d.Set("sla_frequency", response.S("slaFrequency").Data().(float64))
	d.Set("detect_multiplier", response.S("slaDetectMultiplier").Data().(float64))
	d.Set("request_data_size", response.S("reqDataSize").Data().(float64))
	d.Set("type_of_service", response.S("ipv4ToS").Data().(float64))
	d.Set("operation_timeout", response.S("timeout").Data().(float64))
	d.Set("threshold", response.S("threshold").Data().(float64))
	d.Set("ipv6_traffic_class", response.S("ipv6TrfClass").Data().(float64))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	return nil

}

func resourceMSOIPSLAMonitoringPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOIPSLAMonitoringPolicyRead(d, m)
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOIPSLAMonitoringPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if slaType, ok := d.GetOk("sla_type"); ok {
		payload["slaType"] = slaType.(string)
	}

	if destinationPort, ok := d.GetOk("destination_port"); ok {
		payload["slaPort"] = destinationPort.(int)
	}

	if httpVersion, ok := d.GetOk("http_version"); ok {
		payload["httpVersion"] = httpVersion.(string)
	}

	if httpUri, ok := d.GetOk("http_uri"); ok {
		payload["httpUri"] = httpUri.(string)
	}

	if slaFrequency, ok := d.GetOk("sla_frequency"); ok {
		payload["slaFrequency"] = slaFrequency.(int)
	}

	if detectMultiplier, ok := d.GetOk("detect_multiplier"); ok {
		payload["slaDetectMultiplier"] = detectMultiplier.(int)
	}

	if requestDataSize, ok := d.GetOk("request_data_size"); ok {
		payload["reqDataSize"] = requestDataSize.(int)
	}

	if typeOfService, ok := d.GetOk("type_of_service"); ok {
		payload["ipv4ToS"] = typeOfService.(int)
	}

	if operationTimeout, ok := d.GetOk("operation_timeout"); ok {
		payload["timeout"] = operationTimeout.(int)
	}

	if threshold, ok := d.GetOk("threshold"); ok {
		payload["threshold"] = threshold.(int)
	}

	if ipv6TrafficClass, ok := d.GetOk("ipv6_traffic_class"); ok {
		payload["ipv6TrfClass"] = ipv6TrafficClass.(int)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/ipslaMonitoringPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IPSLAMonitoringPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOIPSLAMonitoringPolicyRead(d, m)
}

func resourceMSOIPSLAMonitoringPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "IPSLAMonitoringPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "ipslaMonitoringPolicies")
	if err != nil {
		return err
	}

	setIPSLAMonitoringPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOIPSLAMonitoringPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "ipslaMonitoringPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/ipslaMonitoringPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()
	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("sla_type") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/slaType", updatePath), d.Get("sla_type").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("destination_port") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/slaPort", updatePath), d.Get("destination_port").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("http_version") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/httpVersion", updatePath), d.Get("http_version").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("http_uri") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/httpUri", updatePath), d.Get("http_uri").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("sla_frequency") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/slaFrequency", updatePath), d.Get("sla_frequency").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("detect_multiplier") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/slaDetectMultiplier", updatePath), d.Get("detect_multiplier").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("request_data_size") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/reqDataSize", updatePath), d.Get("request_data_size").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("type_of_service") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/ipv4ToS", updatePath), d.Get("type_of_service").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("operation_timeout") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/timeout", updatePath), d.Get("operation_timeout").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("threshold") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/threshold", updatePath), d.Get("threshold").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("ipv6_traffic_class") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/ipv6TrfClass", updatePath), d.Get("ipv6_traffic_class").(int))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IPSLAMonitoringPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOIPSLAMonitoringPolicyRead(d, m)
}

func resourceMSOIPSLAMonitoringPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "ipslaMonitoringPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/ipslaMonitoringPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
