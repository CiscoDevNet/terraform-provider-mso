package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var humanToApiType = map[string]string{
	"provider": "vzAnyProviderContracts",
	"consumer": "vzAnyConsumerContracts",
}

func resourceMSOTemplateVRFContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateVRFContractCreate,
		Read:   resourceMSOTemplateVRFContractRead,
		Delete: resourceMSOTemplateVRFContractDelete,

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
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"relationship_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"provider",
					"consumer",
				}, false),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOTemplateVRFContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template VRF Contract: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	vrfName := d.Get("vrf_name").(string)
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

	vrfConRef := make(map[string]interface{})
	vrfConRef["contractRef"] = contractRefMap
	path := fmt.Sprintf("/templates/%s/vrfs/%s/%s/-", templateName, vrfName, humanToApiType[relationshipType])
	contractStruct := models.NewTemplateVRFContract("add", path, vrfConRef)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateVRFContractRead(d, m)
}

func resourceMSOTemplateVRFContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	relationshipType := d.Get("relationship_type").(string)

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
	stateVRF := d.Get("vrf_name").(string)
	stateContract := d.Get("contract_name").(string)

	var contract_schema_id, contract_template_name string

	if tempVar, ok := d.GetOk("contract_schema_id"); ok {
		contract_schema_id = tempVar.(string)
	} else {
		contract_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("contract_template_name"); ok {
		contract_template_name = tempVar.(string)
	} else {
		contract_template_name = stateTemplate
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return fmt.Errorf("Unable to get VRF list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				apiVRF := models.StripQuotes(vrfCont.S("name").String())
				if apiVRF == stateVRF {
					d.Set("vrf_name", apiVRF)
					log.Printf("uniiii %v", vrfCont)
					contractCount, err := vrfCont.ArrayCount(humanToApiType[relationshipType])
					if err != nil {
						return fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := vrfCont.ArrayElement(k, humanToApiType[relationshipType])
						if err != nil {
							return err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if contractRef != "{}" && contractRef != "" {
							if stateContract == fmt.Sprintf("%s", split[3]) && contract_schema_id == fmt.Sprintf("%s", split[1]) && contract_template_name == fmt.Sprintf("%s", split[2]) {
								d.SetId(fmt.Sprintf("%s", split[3]))
								d.Set("contract_name", fmt.Sprintf("%s", split[3]))
								d.Set("contract_schema_id", fmt.Sprintf("%s", split[1]))
								d.Set("contract_template_name", fmt.Sprintf("%s", split[2]))
								d.Set("relationship_type", relationshipType)
								found = true
								break
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

func resourceMSOTemplateVRFContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template ExternalEpg Contract: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	contractName := d.Get("contract_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
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
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}
	index, err := fetchVRFContractIndex(cont, templateName, vrfName, contractName, contract_schema_id, contract_template_name, humanToApiType[relationshipType])
	if err != nil {
		d.SetId("")
		return nil
	}

	path := fmt.Sprintf("/templates/%s/vrfs/%s/%s/%d", templateName, vrfName, humanToApiType[relationshipType], index)
	contractStruct := models.NewTemplateVRFContract("remove", path, nil)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), contractStruct)
	if errs != nil {
		return errs
	}
	d.SetId("")
	return nil
}
func fetchVRFContractIndex(cont *container.Container, templateName, vrfName, contractName, contract_schema_id, contract_template_name, relationshipType string) (int, error) {
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
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return index, fmt.Errorf("Unable to get VRF list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return index, err
				}
				apiVRF := models.StripQuotes(vrfCont.S("name").String())
				if apiVRF == vrfName {
					contractCount, err := vrfCont.ArrayCount(relationshipType)
					if err != nil {
						return index, fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := vrfCont.ArrayElement(k, relationshipType)
						if err != nil {
							return index, err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if contractRef != "{}" && contractRef != "" {
							if contractName == fmt.Sprintf("%s", split[3]) && contract_schema_id == fmt.Sprintf("%s", split[1]) && contract_template_name == fmt.Sprintf("%s", split[2]) {
								found = true
								index = k
								break
							}
						}
					}
				}
			}
		}
	}
	if !found {
		return index, fmt.Errorf("Unable to Find the VRF Contract")
	}

	return index, nil
}
