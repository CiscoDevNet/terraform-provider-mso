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
		Update: resourceMSOTemplateAnpEpgContractUpdate,
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
			},
			"contract_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relationship_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
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
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateANP := get_attribute[4]
	stateEPG := get_attribute[6]
	stateContract := get_attribute[8]
	stateRelationshipType := get_attribute[10]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					d.Set("anp_name", apiANP)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							d.Set("epg_name", apiEPG)
							crefCount, err := epgCont.ArrayCount("contractRelationships")
							if err != nil {
								return nil, fmt.Errorf("Unable to get the contract relationships list")
							}
							for l := 0; l < crefCount; l++ {
								crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
								if err != nil {
									return nil, err
								}
								contractRef := models.StripQuotes(crefCont.S("contractRef").String())
								apiRelationshipType := models.StripQuotes(crefCont.S("relationshipType").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								apiContract := match[3]
								if apiContract == stateContract && apiRelationshipType == stateRelationshipType {
									d.SetId(apiContract)
									d.Set("contract_name", match[3])
									d.Set("contract_schema_id", match[1])
									d.Set("contract_template_name", match[2])
									d.Set("relationship_type", apiRelationshipType)
									found = true
									break
								}
							}

						}

					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the Contract %s", stateContract)
	}

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
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateANP := d.Get("anp_name")
	stateEPG := d.Get("epg_name")
	stateContract := d.Get("contract_name")
	stateRelationshipType := d.Get("relationship_type")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							crefCount, err := epgCont.ArrayCount("contractRelationships")
							if err != nil {
								return fmt.Errorf("Unable to get the contract relationships list")
							}
							for l := 0; l < crefCount; l++ {
								crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
								if err != nil {
									return err
								}
								contractRef := models.StripQuotes(crefCont.S("contractRef").String())
								apiRelationshipType := models.StripQuotes(crefCont.S("relationshipType").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								apiContract := match[3]
								if apiContract == stateContract && apiRelationshipType == stateRelationshipType {
									d.SetId(apiContract)
									d.Set("contract_name", match[3])
									d.Set("contract_schema_id", match[1])
									d.Set("contract_template_name", match[2])
									d.Set("relationship_type", apiRelationshipType)
									found = true
									break
								}
							}
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

func resourceMSOTemplateAnpEpgContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
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

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := fetchindex(cont, templateName, anpName, epgName, id, relationship_type)
	if err != nil {
		return err
	}
	if index == -1 {
		fmt.Errorf("The given contract id is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/contractRelationships/%s", templateName, anpName, epgName, indexs)
	crefStruct := models.NewTemplateAnpEpgContract("replace", path, contractRefMap, relationship_type)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), crefStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOTemplateAnpEpgContractRead(d, m)
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

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := fetchindex(cont, templateName, anpName, epgName, id, relationship_type)
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

func fetchindex(cont *container.Container, templateName, anpName, epgName, contractName, relationship_type string) (int, error) {
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
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return index, fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return index, err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")

					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return index, fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return index, err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						log.Println("currentepgname", currentEpgName)
						if currentEpgName == epgName {
							log.Println("found correct epgname")

							contractCount, err := epgCont.ArrayCount("contractRelationships")
							if err != nil {
								return index, fmt.Errorf("No contractRelationships found")
							}
							for s := 0; s < contractCount; s++ {
								contractCont, err := epgCont.ArrayElement(s, "contractRelationships")
								if err != nil {
									return index, err
								}
								contractRef := models.StripQuotes(contractCont.S("contractRef").String())
								apiRelationshipType := models.StripQuotes(contractCont.S("relationshipType").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								match := re.FindStringSubmatch(contractRef)
								apiContract := match[3]
								if apiContract == contractName && apiRelationshipType == relationship_type {
									log.Println("found correct name")
									index = s
									found = true
									break
								}
							}
						}
						if found {
							break
						}
					}
				}
				if found {
					break
				}
			}

		}
		if found {
			break
		}
	}
	return index, nil

}
