package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

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
				Computed: true,
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOTemplateExternalEpgContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External Epg Contract: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	epgName := d.Get("external_epg_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	var contract_schema_id, contract_template_name string

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schema_id = tempVar.(string)
	} else {
		contract_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_template_name = tempVar.(string)
	} else {
		contract_template_name = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schema_id
	contractRefMap["templateName"] = contract_template_name
	contractRefMap["contractName"] = contractName

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/-", templateName, epgName)
	contractStruct := models.NewTemplateExternalEpgContract("add", path, relationshipType, contractRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateExternalEpgContractRead(d, m)
}

func resourceMSOTemplateExternalEpgContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateEPG := d.Get("external_epg_name").(string)
	stateContract := d.Get("contract_name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			epgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get External Epg list")
			}
			for j := 0; j < epgCount; j++ {
				epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				apiEpg := models.StripQuotes(epgCont.S("name").String())
				if apiEpg == stateEPG {
					d.Set("external_epg_name", apiEpg)
					contractCount, err := epgCont.ArrayCount("contractRelationships")
					if err != nil {
						return fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
						if err != nil {
							return err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if stateContract == fmt.Sprintf("%s", split[3]) {
							d.SetId(fmt.Sprintf("%s", split[3]))
							d.Set("contract_name", fmt.Sprintf("%s", split[3]))
							d.Set("contract_schema_id", fmt.Sprintf("%s", split[1]))
							d.Set("contract_template_name", fmt.Sprintf("%s", split[2]))
							d.Set("relationship_type", models.StripQuotes(contractCont.S("relationshipType").String()))
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateExternalEpgContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External Epg Contract: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	epgName := d.Get("external_epg_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	var contract_schema_id, contract_template_name string

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schema_id = tempVar.(string)
	} else {
		contract_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_template_name = tempVar.(string)
	} else {
		contract_template_name = templateName
	}

	contractRefMap := make(map[string]interface{})
	contractRefMap["schemaId"] = contract_schema_id
	contractRefMap["templateName"] = contract_template_name
	contractRefMap["contractName"] = contractName
	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := fetchIndexs(cont, templateName, epgName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		return fmt.Errorf("The given External Epg Contract is not found")
	}

	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/%s", templateName, epgName, indexs)
	contractStruct := models.NewTemplateExternalEpgContract("replace", path, relationshipType, contractRefMap)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if errs != nil {
		return errs
	}
	return resourceMSOTemplateExternalEpgContractRead(d, m)
}

func resourceMSOTemplateExternalEpgContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template ExternalEpg Contract: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	//contractName := d.Get("contract_name").(string)
	templateName := d.Get("template_name").(string)
	epgName := d.Get("external_epg_name").(string)
	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := fetchIndexs(cont, templateName, epgName, id)
	if err != nil {
		return err
	}

	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/%s", templateName, epgName, indexs)
	contractStruct := models.NewTemplateExternalEpgContract("remove", path, "", nil)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
	if errs != nil {
		return errs
	}
	d.SetId("")
	return nil
}
func fetchIndexs(cont *container.Container, templateName, epgName, contractName string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, fmt.Errorf("No Template found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == templateName {
			epgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return index, fmt.Errorf("Unable to get External Epg list")
			}
			for j := 0; j < epgCount; j++ {
				epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return index, err
				}
				apiEpg := models.StripQuotes(epgCont.S("name").String())
				if apiEpg == epgName {
					contractCount, err := epgCont.ArrayCount("contractRelationships")
					if err != nil {
						return index, fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
						if err != nil {
							return index, err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if contractName == fmt.Sprintf("%s", split[3]) {
							found = true
							index = k
							break
						}
					}
				}
			}
		}
	}
	if !found {
		return index, fmt.Errorf("Unable to Find the external Epg Contract")
	}

	return index, nil
}
