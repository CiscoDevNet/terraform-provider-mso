package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteServiceGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteServiceGraphCreate,
		Read:   resourceMSOSchemaSiteServiceGraphRead,
		Update: resourceMSOSchemaSiteServiceGraphUpdate,
		Delete: resourceMSOSchemaSiteServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteServiceGraphImport,
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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configure service nodes for the service graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_dn": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		}),
	}
}

func resourceMSOSchemaSiteServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	siteId := get_attribute[2]
	templateName := get_attribute[4]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	var graphName string
	graphName = get_attribute[6]

	graphCont, _, err := getSiteServiceGraphCont(
		cont,
		schemaId,
		templateName,
		siteId,
		graphName,
	)
	if err != nil {
		d.SetId("")
		log.Printf("sitegraphcont err %v", err)
		return nil, err
	}

	serviceNodeList, err := setServiceNodeList(graphCont)
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(fmt.Sprintf("%s/templates/%s/serviceGraphs/%s", schemaId, templateName, graphName))
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Site Service Node")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if tempVar, ok := d.GetOk("template_name"); ok {
		templateName = tempVar.(string)
	}

	var siteId string
	if site_id, ok := d.GetOk("site_id"); ok {
		siteId = site_id.(string)
	}

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	graphCont, _, err := getTemplateServiceGraphCont(
		cont,
		templateName,
		graphName,
	)
	if err != nil {
		return err
	}

	var siteServiceNodeList []interface{}
	if siteServiceNodes, ok := d.GetOk("service_node"); ok {
		siteServiceNodeList = getServiceNodeList(siteServiceNodes, graphCont)
	}
	serviceNodePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%s/serviceNodes", siteId, templateName, graphName)
	siteServiceGraphPayload := models.NewSchemaSiteServiceGraph("add", serviceNodePath, siteServiceNodeList)
	_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s/template/%s/serviceGraphs/%s", schemaId, siteId, templateName, graphName))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaSiteServiceGraphRead(d, m)
}

func resourceMSOSchemaSiteServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	nodeIdSt := d.Id()
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	graphCont, _, err := getSiteServiceGraphCont(
		cont,
		schemaId,
		templateName,
		siteId,
		graphName,
	)
	if err != nil {
		d.SetId("")
		log.Printf("sitegraphcont err %v", err)
		return nil
	}

	serviceNodeList, err := setServiceNodeList(graphCont)
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	d.SetId(nodeIdSt)
	return nil
}

func resourceMSOSchemaSiteServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Template Service Graph")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if tempVar, ok := d.GetOk("template_name"); ok {
		templateName = tempVar.(string)
	}

	var siteId string
	if site_id, ok := d.GetOk("site_id"); ok {
		siteId = site_id.(string)
	}

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	if d.HasChange("service_node") {
		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			schemaId,
			templateName,
			siteId,
			graphName,
		)
		if err != nil {
			return err
		}

		if siteServiceNodes, ok := d.GetOk("service_node"); ok {
			siteServiceNodeList := getServiceNodeList(siteServiceNodes, graphCont)
			serviceNodePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%s/serviceNodes", siteId, templateName, graphName)
			siteServiceGraphPayload := models.NewSchemaSiteServiceGraph("replace", serviceNodePath, siteServiceNodeList)
			_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload)
			if err != nil {
				return err
			}
		}
	}

	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteServiceGraphRead(d, m)
}

func resourceMSOSchemaSiteServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[NOTE]: Deletion of site Service Graph is not supported by the API.  Site Service Graph will be removed when site is disassociated from the template or when Service Graph is removed at the template level.")
	return nil
}

func getServiceNodeList(siteServiceNodes interface{}, graphCont *container.Container) []interface{} {
	siteServiceNodeList := make([]interface{}, 0, 1)
	templateServiceNodes := graphCont.S("serviceNodes").Data().([]interface{})
	for index, templateServiceNode := range templateServiceNodes {
		templateServiceNodeMap := templateServiceNode.(map[string]interface{})
		siteServiceNodeDeviceList := siteServiceNodes.([]interface{})
		siteServiceNodeMap := siteServiceNodeDeviceList[index].(map[string]interface{})
		serviceNodeMap := map[string]interface{}{
			"serviceNodeRef": templateServiceNodeMap["serviceNodeRef"],
			"device": map[string]interface{}{
				"dn": siteServiceNodeMap["device_dn"],
			},
		}
		siteServiceNodeList = append(siteServiceNodeList, serviceNodeMap)
	}
	return siteServiceNodeList
}

func setServiceNodeList(graphCont *container.Container) ([]interface{}, error) {
	serviceNodeList := make([]interface{}, 0, 1)
	serviceNodes := graphCont.S("serviceNodes").Data().([]interface{})
	for _, val := range serviceNodes {
		serviceNodeValues := val.(map[string]interface{})
		device := serviceNodeValues["device"].(map[string]interface{})["dn"]
		serviceNodeMap := map[string]interface{}{
			"device_dn": device,
		}
		serviceNodeList = append(serviceNodeList, serviceNodeMap)
	}
	return serviceNodeList, nil
}
