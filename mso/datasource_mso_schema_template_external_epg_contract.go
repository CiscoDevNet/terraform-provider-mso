package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMSOTemplateExternalEpgContract() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOTemplateExternalEpgContractRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"relationship_type": &schema.Schema{
				Type: schema.TypeString,
				// Required because relationship_type is part of the identifier:
				// the same contract can be added with different relationship_type values.
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOTemplateExternalEpgContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	template := d.Get("template_name").(string)
	epg := d.Get("external_epg_name").(string)
	contractName := d.Get("contract_name").(string)
	relationshipType := d.Get("relationship_type").(string)

	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaId
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = template
	}

	index, crefCont, err := getSchemaTemplateExtEpgContract(cont, template, epg, contractName, contractSchemaId, contractTemplateName, relationshipType)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
		return fmt.Errorf("Unable to find the External EPG Contract %s in Template %s of Schema Id %s", contractName, contractTemplateName, contractSchemaId)
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s-%s-%s", schemaId, template, epg, contractSchemaId, contractTemplateName, contractName))
	d.Set("contract_name", contractName)
	d.Set("contract_schema_id", contractSchemaId)
	d.Set("contract_template_name", contractTemplateName)
	d.Set("relationship_type", models.StripQuotes(crefCont.S("relationshipType").String()))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
