package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateAnpEpgContract() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateAnpEpgContractRead,

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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
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
				// Computed: true,
				// Changed to required because relationship_type is an identifier value
				// A contract could be added with the same reference but different relationship_type
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}

}

func dataSourceMSOTemplateAnpEpgContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	template := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	contract := d.Get("contract_name").(string)
	contractSchemaId := d.Get("contract_schema_id").(string)
	if contractSchemaId == "" {
		contractSchemaId = schemaId
	}
	contractTemplateName := d.Get("contract_template_name").(string)
	if contractTemplateName == "" {
		contractTemplateName = template
	}
	relationshipType := d.Get("relationship_type").(string)

	index, _, err := getSchemaTemplateEPGContract(cont, template, anp, epg, contract, contractSchemaId, contractTemplateName, relationshipType)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
		return fmt.Errorf("Unable to find the ANP EPG Contract %s in Template %s of Schema Id %s ", contract, contractTemplateName, contractSchemaId)
	}

	d.SetId(fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/contracts/%s-%s-%s", schemaId, template, anp, epg, contractSchemaId, contractTemplateName, contract))
	d.Set("contract_name", contract)
	d.Set("contract_schema_id", contractSchemaId)
	d.Set("contract_template_name", contractTemplateName)
	d.Set("relationship_type", relationshipType)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
