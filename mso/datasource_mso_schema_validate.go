package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaValidate() *schema.Resource {
	return &schema.Resource{
		Read:          datasourceMSOSchemaValidateRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"result": {
				Type:     schema.TypeString,
				Default:  "false",
				Optional: true,
			},
		},
	}
}

func datasourceMSOSchemaValidateRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Validate: Beginning Read")
	msoClient := m.(*client.Client)
	schemaValidate := models.SchemValidate{
		SchmaId: d.Get("schema_id").(string),
	}
	remoteSchemaValidate, err := msoClient.ReadSchemaValidate(&schemaValidate)
	if err != nil {
		return err
	}
	setSchemaValidateAttr(remoteSchemaValidate, d)
	d.SetId(fmt.Sprintf("schemas/%s/validate", remoteSchemaValidate.SchmaId))
	log.Println("[DEBUG] Schema Validate: Reading Completed", d.Id())
	return nil
}

func setSchemaValidateAttr(m *models.SchemValidate, d *schema.ResourceData) {
	d.Set("schema_id", m.SchmaId)
	d.Set("result", m.Result)
}
