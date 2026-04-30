package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTemplateExternalEpgContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateExternalEpgContractCreate,
		Read:   resourceMSOTemplateExternalEpgContractRead,
		Update: resourceMSOTemplateExternalEpgContractUpdate,
		Delete: resourceMSOTemplateExternalEpgContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateExternalEpgContractImport,
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
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"relationship_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		}),
	}
}

// getSchemaTemplateExtEpgContract finds the contractRelationships entry on an external EPG that matches
// contract name + ref schema/template + relationship_type. Returns (index, container, error).
// index == -1 means not found.
func getSchemaTemplateExtEpgContract(cont *container.Container, templateName, epgName, contractName, contractSchemaId, contractTemplateName, relationshipType string) (int, *container.Container, error) {
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, nil, fmt.Errorf("No Template found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, nil, err
		}
		if models.StripQuotes(tempCont.S("name").String()) != templateName {
			continue
		}
		epgCount, err := tempCont.ArrayCount("externalEpgs")
		if err != nil {
			return index, nil, fmt.Errorf("Unable to get External Epg list")
		}
		for j := 0; j < epgCount; j++ {
			epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
			if err != nil {
				return index, nil, err
			}
			if models.StripQuotes(epgCont.S("name").String()) != epgName {
				continue
			}
			contractCount, err := epgCont.ArrayCount("contractRelationships")
			if err != nil {
				return index, nil, fmt.Errorf("Unable to get contract Relationships list")
			}
			for k := 0; k < contractCount; k++ {
				contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
				if err != nil {
					return index, nil, err
				}
				contractRef := models.StripQuotes(contractCont.S("contractRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
				match := re.FindStringSubmatch(contractRef)
				if match == nil {
					continue
				}
				apiRelationshipType := models.StripQuotes(contractCont.S("relationshipType").String())
				if match[3] == contractName &&
					match[1] == contractSchemaId &&
					match[2] == contractTemplateName &&
					apiRelationshipType == relationshipType {
					return k, contractCont, nil
				}
			}
		}
	}
	return index, nil, nil
}

func resourceMSOTemplateExternalEpgContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	if len(get_attribute) < 8 {
		return nil, fmt.Errorf("Invalid import id format, expected <schema_id>/templates/<template>/externalEpgs/<epg>/contractRelationships/<contract>/<relationship_type>")
	}
	schemaId := get_attribute[0]
	stateTemplate := get_attribute[2]
	stateEPG := get_attribute[4]
	stateContract := get_attribute[6]
	stateType := get_attribute[7]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	// Without explicit contract_schema_id / contract_template_name, default to the same schema/template.
	index, crefCont, err := getSchemaTemplateExtEpgContract(cont, stateTemplate, stateEPG, stateContract, schemaId, stateTemplate, stateType)
	if err != nil {
		return nil, err
	}
	if index == -1 {
		d.SetId("")
		return nil, fmt.Errorf("External Epg Contract Not Found")
	}

	contractRef := models.StripQuotes(crefCont.S("contractRef").String())
	re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
	match := re.FindStringSubmatch(contractRef)

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("external_epg_name", stateEPG)
	d.Set("contract_name", match[3])
	d.Set("contract_schema_id", match[1])
	d.Set("contract_template_name", match[2])
	d.Set("relationship_type", models.StripQuotes(crefCont.S("relationshipType").String()))
	d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s/%s", schemaId, stateTemplate, stateEPG, stateContract, stateType))

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTemplateExternalEpgContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External Epg Contract: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	epgName := d.Get("external_epg_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaID
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = templateName
	}

	contractRefMap := map[string]interface{}{
		"schemaId":     contractSchemaId,
		"templateName": contractTemplateName,
		"contractName": contractName,
	}

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/-", templateName, epgName)
	contractStruct := models.NewTemplateExternalEpgContract("add", path, relationshipType, contractRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s/%s", schemaID, templateName, epgName, contractName, relationshipType))
	return resourceMSOTemplateExternalEpgContractRead(d, m)
}

func resourceMSOTemplateExternalEpgContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	templateName := d.Get("template_name").(string)
	epgName := d.Get("external_epg_name").(string)
	contractName := d.Get("contract_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaId
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = templateName
	}

	index, crefCont, err := getSchemaTemplateExtEpgContract(cont, templateName, epgName, contractName, contractSchemaId, contractTemplateName, relationshipType)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
	} else {
		contractRef := models.StripQuotes(crefCont.S("contractRef").String())
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
		match := re.FindStringSubmatch(contractRef)
		apiRelationshipType := models.StripQuotes(crefCont.S("relationshipType").String())
		d.Set("contract_name", match[3])
		d.Set("contract_schema_id", match[1])
		d.Set("contract_template_name", match[2])
		d.Set("relationship_type", apiRelationshipType)
		d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s/%s", schemaId, templateName, epgName, match[3], apiRelationshipType))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTemplateExternalEpgContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External Epg Contract: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	epgName := d.Get("external_epg_name").(string)
	contractName := d.Get("contract_name").(string)

	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaID
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = templateName
	}

	// Use the old relationship_type to locate the existing contract; relationship_type is the only updatable field.
	oldRelationshipType, newRelationshipType := d.GetChange("relationship_type")

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	index, _, err := getSchemaTemplateExtEpgContract(cont, templateName, epgName, contractName, contractSchemaId, contractTemplateName, oldRelationshipType.(string))
	if err != nil {
		return err
	}
	if index == -1 {
		return fmt.Errorf("Unable to find the External Epg Contract %s with relationship type %s", contractName, oldRelationshipType.(string))
	}

	updatePath := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/%d", templateName, epgName, index)
	payloadCon := container.New()
	payloadCon.Array()

	if d.HasChange("relationship_type") {
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/relationshipType", updatePath), newRelationshipType.(string))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaID), payloadCon)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s/%s", schemaID, templateName, epgName, contractName, newRelationshipType.(string)))
	return resourceMSOTemplateExternalEpgContractRead(d, m)
}

func resourceMSOTemplateExternalEpgContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template ExternalEpg Contract: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	epgName := d.Get("external_epg_name").(string)
	contractName := d.Get("contract_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaID
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = templateName
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	index, _, err := getSchemaTemplateExtEpgContract(cont, templateName, epgName, contractName, contractSchemaId, contractTemplateName, relationshipType)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}

	indexs := strconv.Itoa(index)
	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/%s", templateName, epgName, indexs)
	contractStruct := models.NewTemplateExternalEpgContract("remove", path, "", nil)

	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return errs
	}
	d.SetId("")
	return nil
}
