package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMSOServiceNodeType() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOServiceNodeTypeRead,

		SchemaVersion: version,

		DeprecationMessage: "mso_service_node_type data source is deprecated: it remains functional on Nexus Dashboard (ND) 4.3+ / NDO 5.3+, but the corresponding mso_service_node_type resource is no longer functional. The data source will be removed once ND 4.2 / NDO 5.2 is no longer supported.",

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOServiceNodeTypeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Read %s", d.Id())

	msoClient := m.(*client.Client)

	typeName := d.Get("name").(string)

	found := false

	cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return err
	}

	nodesCount, err := cont.ArrayCount("serviceNodeTypes")
	if err != nil {
		return err
	}

	for i := 0; i < nodesCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
		if err != nil {
			return err
		}

		apiName := models.StripQuotes(nodeCont.S("name").String())

		if apiName == typeName {
			d.SetId(models.StripQuotes(nodeCont.S("id").String()))
			d.Set("name", models.StripQuotes(nodeCont.S("name").String()))
			d.Set("display_name", models.StripQuotes(nodeCont.S("displayName").String()))
			found = true
		}
	}
	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find service node type %s", typeName)
	}
	log.Printf("[DEBUG] Read Finished Successfully %s", d.Id())
	return nil
}
