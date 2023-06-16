package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOServiceNodeType() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOServiceNodeTypeCreate,
		Read:   resourceMSOServiceNodeTypeRead,
		Delete: resourceMSOServiceNodeTypeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOServiceNodeTypeImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOServiceNodeTypeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Beginning Import %s", d.Id())

	msoClient := m.(*client.Client)

	typeName := d.Id()

	found := false

	cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return nil, err
	}

	nodesCount, err := cont.ArrayCount("serviceNodeTypes")
	if err != nil {
		return nil, err
	}

	for i := 0; i < nodesCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
		if err != nil {
			return nil, err
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
		return nil, fmt.Errorf("Unable to find service node type %s", typeName)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOServiceNodeTypeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Service Node Type: Beginning Creation")
	msoClient := m.(*client.Client)
	typeAttr := models.ServiceNodeTypeAttributes{}

	if name, ok := d.GetOk("name"); ok {
		typeAttr.Name = name.(string)
	}

	if display_name, ok := d.GetOk("display_name"); ok {
		typeAttr.DisplayName = display_name.(string)
	} else {
		typeAttr.DisplayName = typeAttr.Name
	}

	nodeType := models.NewServiceNodeType(typeAttr)
	d.Partial(true)
	cont, err := msoClient.Save("api/v1/schemas/service-node-types", nodeType)

	if err != nil {
		return err
	}
	d.Partial(false)
	d.SetId(models.StripQuotes(cont.S("id").String()))
	log.Printf("[DEBUG] Creation finished successfully %s", d.Id())
	return resourceMSOServiceNodeTypeRead(d, m)
}

func resourceMSOServiceNodeTypeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Read %s", d.Id())

	msoClient := m.(*client.Client)

	typeId := d.Id()

	found := false

	cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
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

		apiId := models.StripQuotes(nodeCont.S("id").String())

		if apiId == typeId {
			d.SetId(apiId)
			d.Set("name", models.StripQuotes(nodeCont.S("name").String()))
			d.Set("display_name", models.StripQuotes(nodeCont.S("displayName").String()))
			found = true
		}
	}
	if !found {
		d.SetId("")
	}
	log.Printf("[DEBUG] Read Finished Successfully %s", d.Id())
	return nil
}

func resourceMSOServiceNodeTypeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Destroy %s", d.Id())

	msoClient := m.(*client.Client)

	typeId := d.Id()
	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/schemas/service-node-types/%s", typeId))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
