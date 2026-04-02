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

func resourceMSOTemplateAnpEpgContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateAnpEpgContractCreate,
		Read:   resourceMSOTemplateAnpEpgContractRead,
		Delete: resourceMSOTemplateAnpEpgContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateAnpEpgContractImport,
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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
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
				Computed: true,
				ForceNew: true,
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"relationship_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}

}

func resourceMSOTemplateAnpEpgContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	stateTemplate := get_attribute[2]
	stateANP := get_attribute[4]
	stateEPG := get_attribute[6]
	stateContract := get_attribute[8]
	stateRelationshipType := get_attribute[10]

	index, crefCont, err := findEpgContractRelationship(cont, stateTemplate, stateANP, stateEPG, stateContract, schemaId, stateTemplate, stateRelationshipType)
	if err != nil {
		return nil, err
	}
	if index == -1 {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the Contract %s", stateContract)
	}

	contractRef := models.StripQuotes(crefCont.S("contractRef").String())
	re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
	match := re.FindStringSubmatch(contractRef)

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("anp_name", stateANP)
	d.Set("epg_name", stateEPG)
	d.SetId(match[3])
	d.Set("contract_name", match[3])
	d.Set("contract_schema_id", match[1])
	d.Set("contract_template_name", match[2])
	d.Set("relationship_type", models.StripQuotes(crefCont.S("relationshipType").String()))

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTemplateAnpEpgContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	contractName := d.Get("contract_name").(string)

	var relationship_type, contract_schemaid, contract_templatename string
	if tempVar, ok := d.GetOk("relationship_type"); ok {
		relationship_type = tempVar.(string)
	}

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schemaid = tempVar.(string)
	} else {
		contract_schemaid = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_templatename = tempVar.(string)
	} else {
		contract_templatename = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schemaid
	contractRefMap["templateName"] = contract_templatename
	contractRefMap["contractName"] = contractName

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/contractRelationships/-", templateName, anpName, epgName)
	bdStruct := models.NewTemplateAnpEpgContract("add", path, contractRefMap, relationship_type)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), bdStruct)
	if err != nil {
		return err
	}
	return resourceMSOTemplateAnpEpgContractRead(d, m)
}

func resourceMSOTemplateAnpEpgContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	contractName := d.Get("contract_name").(string)
	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaId
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = templateName
	}
	relationshipType := d.Get("relationship_type").(string)

	index, crefCont, err := findEpgContractRelationship(cont, templateName, anpName, epgName, contractName, contractSchemaId, contractTemplateName, relationshipType)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
	} else {
		contractRef := models.StripQuotes(crefCont.S("contractRef").String())
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
		match := re.FindStringSubmatch(contractRef)
		d.SetId(match[3])
		d.Set("contract_name", match[3])
		d.Set("contract_schema_id", match[1])
		d.Set("contract_template_name", match[2])
		d.Set("relationship_type", models.StripQuotes(crefCont.S("relationshipType").String()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateAnpEpgContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template ANP EPG Contract: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	contractName := d.Get("contract_name").(string)

	var relationship_type, contract_schemaid, contract_templatename string
	if tempVar, ok := d.GetOk("relationship_type"); ok {
		relationship_type = tempVar.(string)
	}

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schemaid = tempVar.(string)
	} else {
		contract_schemaid = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_templatename = tempVar.(string)
	} else {
		contract_templatename = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schemaid
	contractRefMap["templateName"] = contract_templatename
	contractRefMap["contractName"] = contractName

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, _, err := findEpgContractRelationship(cont, templateName, anpName, epgName, contractName, contract_schemaid, contract_templatename, relationship_type)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/contractRelationships/%s", templateName, anpName, epgName, indexs)
	crefStruct := models.NewTemplateAnpEpgContract("remove", path, contractRefMap, relationship_type)

	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), crefStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return errs
	}
	d.SetId("")
	return resourceMSOTemplateAnpEpgContractRead(d, m)
}

func findEpgContractRelationship(cont *container.Container, templateName, anpName, epgName, contractName, contractSchemaId, contractTemplateName, relationshipType string) (int, *container.Container, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, nil, fmt.Errorf("No Template found")
	}
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, nil, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return index, nil, fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount && !found; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return index, nil, err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return index, nil, fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount && !found; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return index, nil, err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							contractCount, err := epgCont.ArrayCount("contractRelationships")
							if err != nil {
								return index, nil, fmt.Errorf("No contractRelationships found")
							}
							for s := 0; s < contractCount; s++ {
								contractCont, err := epgCont.ArrayElement(s, "contractRelationships")
								if err != nil {
									return index, nil, err
								}
								contractRef := models.StripQuotes(contractCont.S("contractRef").String())
								apiRelationshipType := models.StripQuotes(contractCont.S("relationshipType").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								apiContractSchemaId := match[1]
								apiContractTemplateName := match[2]
								apiContract := match[3]
								if apiContract == contractName &&
									apiRelationshipType == relationshipType &&
									apiContractSchemaId == contractSchemaId &&
									apiContractTemplateName == contractTemplateName {
									index = s
									found = true
									return index, contractCont, nil
								}
							}
						}
					}
				}
			}
		}
	}
	return index, nil, nil
}
