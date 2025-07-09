package mso

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateCreate,
		Read:   resourceMSOTemplateRead,
		Update: resourceMSOTemplateUpdate,
		Delete: resourceMSOTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateImport,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tenant",
					"l3out",
					"fabric_policy",
					"fabric_resource",
					"monitoring_tenant",
					"monitoring_access",
					"service_device",
				}, false),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sites": {
				Type:     schema.TypeList, // Set cannot not be used because the order of sites is important for deletion, since full list replacement is not supported in API
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// ValidateFunc not supported on type list / set else duplication could be caught with duplication.ListOfUniqueStrings
				// │ Error: Internal validation of the provider failed! This is always a bug
				// │ with the provider itself, and not a user issue. Please report
				// │ this bug:
				// │
				// │ 1 error occurred:
				// │       * resource mso_template: sites: ValidateFunc is not yet supported on lists or sets.
			},
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			oldSites, newSites := diff.GetChange("sites")

			sites := convertToListOfStrings(newSites.([]interface{}))

			if len(oldSites.([]interface{})) != len(newSites.([]interface{})) {
				return nil
			}

			for _, oldSite := range oldSites.([]interface{}) {
				if !valueInSliceofStrings(oldSite.(string), sites) {
					return nil
				}
			}
			diff.SetNew("sites", oldSites)

			return nil
		},
	}
}

type ndoTemplateType struct {
	templateType          string // templateType in payload
	templateTypeContainer string // templateType container in payload
	tenant                bool   // tenant required
	siteAmount            int    // 1 = 1 site, 2 = multiple sites
	templateContainer     bool   // configuration is set in template container in payload
}

var ndoTemplateTypes = map[string]ndoTemplateType{
	"tenant": ndoTemplateType{
		templateType:          "tenantPolicy",
		templateTypeContainer: "tenantPolicyTemplate",
		tenant:                true,
		siteAmount:            2,
		templateContainer:     true,
	},
	"l3out": ndoTemplateType{
		templateType:          "l3out",
		templateTypeContainer: "l3outTemplate",
		tenant:                true,
		siteAmount:            1,
		templateContainer:     false,
	},
	"fabric_policy": ndoTemplateType{
		templateType:          "fabricPolicy",
		templateTypeContainer: "fabricPolicyTemplate",
		tenant:                false,
		siteAmount:            2,
		templateContainer:     true,
	},
	"fabric_resource": ndoTemplateType{
		templateType:          "fabricResource",
		templateTypeContainer: "fabricResourceTemplate",
		tenant:                false,
		siteAmount:            2,
		templateContainer:     true,
	},
	"monitoring_tenant": ndoTemplateType{
		templateType:          "monitoring",
		templateTypeContainer: "monitoringTemplate",
		tenant:                true,
		siteAmount:            1,
		templateContainer:     true,
	},
	"monitoring_access": ndoTemplateType{
		templateType:          "monitoring",
		templateTypeContainer: "monitoringTemplate",
		tenant:                false,
		siteAmount:            1,
		templateContainer:     true,
	},
	"service_device": ndoTemplateType{
		templateType:          "serviceDevice",
		templateTypeContainer: "deviceTemplate",
		tenant:                true,
		siteAmount:            2,
		templateContainer:     true,
	},
}

type ndoTemplate struct {
	id           string
	templateName string
	templateType string
	tenantId     string
	sites        []string
	msoClient    *client.Client
}

func (ndoTemplate *ndoTemplate) SetSchemaResourceData(d *schema.ResourceData) {
	d.SetId(ndoTemplate.id)
	d.Set("template_name", ndoTemplate.templateName)
	d.Set("template_type", ndoTemplate.templateType)
	d.Set("tenant_id", ndoTemplate.tenantId)
	d.Set("sites", ndoTemplate.sites)
}

func (ndoTemplate *ndoTemplate) validateConfig() []error {
	errors := []error{}

	if ndoTemplate.tenantId != "" && !ndoTemplateTypes[ndoTemplate.templateType].tenant {
		errors = append(errors, fmt.Errorf("A Tenant cannot be attached to a template of type %s.", ndoTemplate.templateType))
	}
	if ndoTemplate.tenantId == "" && ndoTemplateTypes[ndoTemplate.templateType].tenant {
		errors = append(errors, fmt.Errorf("A Tenant is required for a template of type %s. Use the `tenant_id` attribute to specify the Tenant to associate with this template.", ndoTemplate.templateType))
	}
	if len(ndoTemplate.sites) == 0 && ndoTemplateTypes[ndoTemplate.templateType].siteAmount == 1 {
		errors = append(errors, fmt.Errorf("At least one site is required for a template of type %s.", ndoTemplate.templateType))
	}
	if len(ndoTemplate.sites) > 1 && ndoTemplateTypes[ndoTemplate.templateType].siteAmount == 1 {
		errors = append(errors, fmt.Errorf("Only one site is allowed for a template of type %s.", ndoTemplate.templateType))
	}
	duplicates := duplicatesInList(ndoTemplate.sites)
	if len(duplicates) > 0 {
		duplicatesErrors := []error{fmt.Errorf("Duplication found in the sites list")}
		for _, site := range duplicates {
			duplicatesErrors = append(duplicatesErrors, fmt.Errorf("Site %s is duplicated", site))
		}
		return duplicatesErrors
	}

	return errors
}

func (ndoTemplate *ndoTemplate) ToMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"displayName":  ndoTemplate.templateName,
		"templateType": ndoTemplateTypes[ndoTemplate.templateType].templateType,
		ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer: ndoTemplate.createTypeSpecificPayload(),
	}, nil
}

func (ndoTemplate *ndoTemplate) createTypeSpecificPayload() map[string]interface{} {
	if ndoTemplate.templateType == "tenant" {
		return map[string]interface{}{"template": map[string]interface{}{"tenantId": ndoTemplate.tenantId}, "sites": ndoTemplate.createSitesPayload()}
	} else if ndoTemplate.templateType == "l3out" {
		return map[string]interface{}{"tenantId": ndoTemplate.tenantId, "siteId": ndoTemplate.createSitesPayload()[0]["siteId"]}
	} else if ndoTemplate.templateType == "fabric_policy" {
		return map[string]interface{}{"sites": ndoTemplate.createSitesPayload()}
	} else if ndoTemplate.templateType == "fabric_resource" {
		return map[string]interface{}{"sites": ndoTemplate.createSitesPayload()}
	} else if ndoTemplate.templateType == "monitoring_tenant" {
		return map[string]interface{}{"template": map[string]interface{}{"mtType": "tenant", "tenant": ndoTemplate.tenantId}, "sites": ndoTemplate.createSitesPayload()}
	} else if ndoTemplate.templateType == "monitoring_access" {
		return map[string]interface{}{"template": map[string]interface{}{"mtType": "access"}, "sites": ndoTemplate.createSitesPayload()}
	} else if ndoTemplate.templateType == "service_device" {
		return map[string]interface{}{"template": map[string]interface{}{"tenantId": ndoTemplate.tenantId}, "sites": ndoTemplate.createSitesPayload()}
	}
	return nil
}

func (ndoTemplate *ndoTemplate) createSitesPayload() []map[string]interface{} {
	siteIds := []map[string]interface{}{}
	for _, siteId := range ndoTemplate.sites {
		siteIds = append(siteIds, ndoTemplate.createSitePayload(siteId))
	}
	return siteIds
}

func (ndoTemplate *ndoTemplate) createSitePayload(siteId string) map[string]interface{} {
	return map[string]interface{}{"siteId": siteId}
}

func (ndoTemplate *ndoTemplate) getTemplate(errorNotFound bool) error {

	if ndoTemplate.id == "" {
		err := ndoTemplate.setTemplateId()
		if err != nil {
			return err
		}
	}

	cont, err := ndoTemplate.msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", ndoTemplate.id))
	if err != nil {

		// 404 scenario where the json response is not a valid json and cannot be parsed overwriting the response from mso in the client
		if err.Error() == "invalid character 'p' after top-level value" {
			return fmt.Errorf("Template ID %s invalid", ndoTemplate.id)
		}

		// If template is not found, we can remove the id and try to find the template by name
		if !errorNotFound && (cont.S("code").String() == "400" && strings.Contains(cont.S("message").String(), fmt.Sprintf("Template ID %s invalid", ndoTemplate.id))) {
			ndoTemplate.id = ""
			return nil
		}
		return err
	}

	ndoTemplate.sites = []string{}

	ndoTemplate.templateName = models.StripQuotes(cont.S("displayName").String())
	templateType := models.StripQuotes(cont.S("templateType").String())

	if templateType == "tenantPolicy" {
		ndoTemplate.templateType = "tenant"
		ndoTemplate.tenantId = models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("template").S("tenantId").String())

	} else if templateType == "l3out" {
		ndoTemplate.templateType = "l3out"
		ndoTemplate.tenantId = models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("tenantId").String())
		ndoTemplate.sites = append(ndoTemplate.sites, models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("siteId").String()))

	} else if templateType == "fabricPolicy" {
		ndoTemplate.templateType = "fabric_policy"

	} else if templateType == "fabricResource" {
		ndoTemplate.templateType = "fabric_resource"

	} else if templateType == "monitoring" {
		ndoTemplate.templateType = "monitoring_access"
		if models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("template").S("mtType").String()) == "tenant" {
			ndoTemplate.templateType = "monitoring_tenant"
			ndoTemplate.tenantId = models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("template").S("tenant").String())
		}

	} else if templateType == "serviceDevice" {
		ndoTemplate.templateType = "service_device"
		ndoTemplate.tenantId = models.StripQuotes(cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("template").S("tenantId").String())
	}

	if ndoTemplate.templateType != "l3out" {
		if cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).Exists("sites") {
			for _, site := range cont.S(ndoTemplateTypes[ndoTemplate.templateType].templateTypeContainer).S("sites").Data().([]interface{}) {
				siteId := models.StripQuotes(site.(map[string]interface{})["siteId"].(string))
				ndoTemplate.sites = append(ndoTemplate.sites, siteId)
			}
		}
	}

	return nil

}

func (ndoTemplate *ndoTemplate) setTemplateId() error {
	cont, err := ndoTemplate.msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/summaries"))
	if err != nil {
		return err
	}

	templates, err := cont.Children()
	if err != nil {
		return err
	}

	for _, template := range templates {
		if ndoTemplate.templateName == models.StripQuotes(template.S("templateName").String()) && ndoTemplateTypes[ndoTemplate.templateType].templateType == models.StripQuotes(template.S("templateType").String()) {
			ndoTemplate.id = models.StripQuotes(template.S("templateId").String())
			return nil
		}
	}

	return fmt.Errorf("Template with name '%s' not found.", ndoTemplate.templateName)
}

func resourceMSOTemplateCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] MSO Template Resource: Beginning Create", d.Id())
	msoClient := m.(*client.Client)

	ndoTemplate := ndoTemplate{
		msoClient:    msoClient,
		templateName: d.Get("template_name").(string),
		templateType: d.Get("template_type").(string),
		sites:        getListOfStringsFromSchemaList(d, "sites"),
	}

	if tenantId, ok := d.GetOk("tenant_id"); ok {
		ndoTemplate.tenantId = tenantId.(string)
	}

	validationErrors := ndoTemplate.validateConfig()
	if len(validationErrors) > 0 {
		d.SetId("")
		return errors.Join(validationErrors...)
	}

	response, err := msoClient.Save("api/v1/templates", &ndoTemplate)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(response.S("templateId").String()))
	log.Println("[DEBUG] MSO Template Resource: Create Completed", d.Id())
	return nil
}

func resourceMSOTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] MSO Template Resource: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	ndoTemplate := ndoTemplate{msoClient: msoClient, id: d.Id()}
	err := ndoTemplate.getTemplate(false)
	if err != nil {
		return err
	}
	ndoTemplate.SetSchemaResourceData(d)
	log.Println("[DEBUG] MSO Template Resource: Read Completed", d.Id())
	return nil
}

func resourceMSOTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] MSO Template Resource: Beginning Update", d.Id())
	msoClient := m.(*client.Client)

	templateType := d.Get("template_type").(string)
	templateName := d.Get("template_name").(string)
	sites := getListOfStringsFromSchemaList(d, "sites")

	if ndoTemplateTypes[templateType].siteAmount == 1 && d.HasChange("sites") {
		return fmt.Errorf("Cannot change site for template of type %s.", templateType)
	}

	ndoTemplate := ndoTemplate{
		msoClient:    msoClient,
		templateName: templateName,
		templateType: templateType,
		sites:        sites,
	}

	if tenantId, ok := d.GetOk("tenant_id"); ok {
		ndoTemplate.tenantId = tenantId.(string)
	}

	validationErrors := ndoTemplate.validateConfig()
	if len(validationErrors) > 0 {
		return errors.Join(validationErrors...)
	}

	payloadCon := container.New()
	payloadCon.Array()

	if d.HasChange("template_name") {
		err := addPatchPayloadToContainer(payloadCon, "replace", "/displayName", templateName)
		if err != nil {
			return err
		}
	}

	if d.HasChange("sites") {
		// Replace operation is not supported for sites, so we need to remove and add sites individually

		oldSites, _ := d.GetChange("sites")

		// Reversed loop to remove sites from the end of the list first, to prevent index shifts with wrong deletes
		for index := len(oldSites.([]interface{})) - 1; index >= 0; index-- {
			if !valueInSliceofStrings(oldSites.([]interface{})[index].(string), sites) {
				err := addPatchPayloadToContainer(payloadCon, "remove", fmt.Sprintf("/%s/sites/%d", ndoTemplateTypes[templateType].templateTypeContainer, index), nil)
				if err != nil {
					return err
				}
			}
		}

		for _, site := range sites {
			if !valueInSliceofStrings(site, convertToListOfStrings(oldSites.([]interface{}))) {
				err := addPatchPayloadToContainer(payloadCon, "add", fmt.Sprintf("/%s/sites/-", ndoTemplateTypes[templateType].templateTypeContainer), ndoTemplate.createSitePayload(site))
				if err != nil {
					return err
				}
			}
		}
	}

	err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", d.Id()), payloadCon)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] MSO Template Resource: Updating Completed", d.Id())
	return nil
}

func resourceMSOTemplateDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] MSO Template Resource: Beginning Delete", d.Id())
	msoClient := m.(*client.Client)

	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/templates/%s", d.Id()))
	if err != nil {
		return err
	}
	log.Println("[DEBUG] MSO Template Resource: Delete Completed", d.Id())
	d.SetId("")
	return nil
}

func resourceMSOTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] MSO Template Resource: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	ndoTemplate := ndoTemplate{msoClient: msoClient, id: d.Id()}
	err := ndoTemplate.getTemplate(true)
	if err != nil {
		return nil, err
	}
	ndoTemplate.SetSchemaResourceData(d)
	log.Println("[DEBUG] MSO Template Resource: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}
