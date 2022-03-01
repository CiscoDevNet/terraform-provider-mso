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

func resourceMSOSchemaSiteServiceGraphNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteServiceGraphNodeCreate,
		Read:   resourceMSOSchemaSiteServiceGraphNodeRead,
		Update: resourceMSOSchemaSiteServiceGraphNodeUpdate,
		Delete: resourceMSOSchemaSiteServiceGraphNodeDelete,

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
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"service_node_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"site_nodes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"tenant_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"node_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		}),
	}
}

func resourceMSOSchemaSiteServiceGraphNodeCreate(d *schema.ResourceData, m interface{}) error {
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

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}
	var nodeType string
	if tempVar, ok := d.GetOk("service_node_type"); ok {
		nodeType = tempVar.(string)
	}

	nodeId, err := getNodeIdFromName(msoClient, nodeType)
	if err != nil {
		return err
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	grapCont, _, err := getTemplateServiceGraphCont(
		cont,
		templateName,
		graphName,
	)
	if err != nil {
		return err
	}

	nodeInd := getTemplateNodeIndex(grapCont)

	if nodeInd == -1 {
		return fmt.Errorf("Unable to get Temlate node index list")
	}

	templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/%s/serviceNodes/-", templateName, graphName)

	templatePayload := map[string]interface{}{
		"name":              fmt.Sprintf("tfnode%d", nodeInd),
		"index":             nodeInd,
		"serviceNodeTypeId": nodeId,
	}

	templatePatchStruct := models.NewTemplateServiceGraph("add", templatePath, templatePayload)
	sitePayload := make([]models.Model, 0, 1)

	sitePayload = append(sitePayload, templatePatchStruct)

	var siteParams []interface{}

	if tempVar, ok := d.GetOk("site_nodes"); ok {
		siteParams = tempVar.([]interface{})

		for _, site := range siteParams {
			siteMap := site.(map[string]interface{})
			if siteMap["site_id"] == "" {
				return fmt.Errorf("site_id is required in site_nodes list")
			}

			if siteMap["tenant_name"] == "" {
				return fmt.Errorf("tenant_name is required in site_nodes list")
			}

			if siteMap["node_name"] == "" {
				return fmt.Errorf("node_name is required in site_nodes list")
			}

			// <---- Begin site payload creation
			siteVarMap := map[string]interface{}{
				"serviceNodeRef": map[string]interface{}{
					"schemaId":         schemaId,
					"templateName":     templateName,
					"serviceGraphName": graphName,
					"serviceNodeName":  fmt.Sprintf("tfnode%d", nodeInd),
				},
				"device": map[string]interface{}{
					"dn": fmt.Sprintf("uni/tn-%s/lDevVip-%s", siteMap["tenant_name"].(string), siteMap["node_name"].(string)),
				},
			}
			// ----> site payload creation ends

			_, graphind, err := getSiteServiceGraphCont(
				cont,
				schemaId,
				templateName,
				siteMap["site_id"].(string),
				graphName,
			)

			if err != nil {
				return err
			}

			sitePath := fmt.Sprintf(
				"/sites/%s-%s/serviceGraphs/%d/serviceNodes/-",
				siteMap["site_id"].(string),
				templateName,
				graphind,
			)

			sitePayload = append(sitePayload, models.NewTemplateServiceGraph("add", sitePath, siteVarMap))

		}
	}
	_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), sitePayload...)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("tfnode%d", nodeInd))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaSiteServiceGraphNodeRead(d, m)
}

func getTemplateNodeIndex(graphCont *container.Container) int {

	nodeCount, err := graphCont.ArrayCount("serviceNodes")
	if err != nil {
		return -1
	}

	return (nodeCount + 1)
}

func resourceMSOSchemaSiteServiceGraphNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	nodeIdSt := d.Id()
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	stateTemplate := d.Get("template_name").(string)
	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	sgCont, _, err := getTemplateServiceGraphCont(cont, stateTemplate, graphName)

	if err != nil {
		d.SetId("")
		log.Printf("graphcont err %v", err)
		return nil
	}
	var nodeType string
	if tempVar, ok := d.GetOk("service_node_type"); ok {
		nodeType = tempVar.(string)
	}
	nodeId, err := getNodeIdFromName(msoClient, nodeType)
	if err != nil {
		return err
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)

	_, _, err = getTemplateServiceNodeCont(sgCont, nodeIdSt, nodeId)

	if err != nil {
		d.Set("service_node_type", "")
		log.Printf("nodecont err %v", err)

	} else {
		d.Set("service_node_type", nodeType)
	}

	var siteParams []interface{}

	if tempVar, ok := d.GetOk("site_nodes"); ok {
		siteParams = tempVar.([]interface{})

		for ind, site := range siteParams {
			siteMap := site.(map[string]interface{})
			if siteMap["site_id"] == "" {
				return fmt.Errorf("site_id is required in site_nodes list")
			}

			graphCont, _, err := getSiteServiceGraphCont(
				cont,
				schemaId,
				stateTemplate,
				siteMap["site_id"].(string),
				graphName,
			)

			if err != nil {
				d.SetId("")
				log.Printf("sitegraphcont err %v", err)
				return nil
			}

			nodeCont, _, err := getSiteServiceNodeCont(
				graphCont,
				schemaId,
				stateTemplate,
				graphName,
				nodeIdSt,
			)

			if err != nil {
				d.SetId("")
				log.Printf("sitenodecont err %v", err)
				return nil
			}

			deviceDn := models.StripQuotes(nodeCont.S("device", "dn").String())

			dnSplit := strings.Split(deviceDn, "/")

			tnName := strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")
			siteMap["tenant_name"] = tnName
			siteMap["node_name"] = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")

			siteParams[ind] = siteMap

		}

		d.Set("site_nodes", siteParams)
	} else {
		d.Set("site_nodes", nil)
	}

	d.SetId(nodeIdSt)
	return nil
}
func resourceMSOSchemaSiteServiceGraphNodeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Update Template Service Graph")
	msoClient := m.(*client.Client)
	nodeId := d.Id()
	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if tempVar, ok := d.GetOk("template_name"); ok {
		templateName = tempVar.(string)
	}

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	if d.HasChange("site_nodes") {
		var siteParams []interface{}
		sitePayload := make([]models.Model, 0, 1)

		if tempVar, ok := d.GetOk("site_nodes"); ok {
			siteParams = tempVar.([]interface{})

			for _, site := range siteParams {
				siteMap := site.(map[string]interface{})
				if siteMap["site_id"] == "" {
					return fmt.Errorf("site_id is required in site_nodes list")
				}

				if siteMap["tenant_name"] == "" {
					return fmt.Errorf("tenant_name is required in site_nodes list")
				}

				if siteMap["node_name"] == "" {
					return fmt.Errorf("node_name is required in site_nodes list")
				}

				// <---- Begin site payload creation
				siteVarMap := map[string]interface{}{
					"serviceNodeRef": map[string]interface{}{
						"schemaId":         schemaId,
						"templateName":     templateName,
						"serviceGraphName": graphName,
						"serviceNodeName":  nodeId,
					},
					"device": map[string]interface{}{
						"dn": fmt.Sprintf("uni/tn-%s/lDevVip-%s", siteMap["tenant_name"].(string), siteMap["node_name"].(string)),
					},
				}

				// ----> site payload creation ends

				cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
				if err != nil {
					return err
				}

				graphCont, graphind, err := getSiteServiceGraphCont(
					cont,
					schemaId,
					templateName,
					siteMap["site_id"].(string),
					graphName,
				)

				if err != nil {
					return err
				}

				_, nodeind, err := getSiteServiceNodeCont(
					graphCont,
					schemaId,
					templateName,
					graphName,
					nodeId,
				)

				if err != nil {
					return err
				}

				sitePath := fmt.Sprintf(
					"/sites/%s-%s/serviceGraphs/%d/serviceNodes/%d",
					siteMap["site_id"].(string),
					templateName,
					graphind,
					nodeind,
				)

				sitePayload = append(sitePayload, models.NewTemplateServiceGraph("replace", sitePath, siteVarMap))

			}
			_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), sitePayload...)

			if err != nil {
				return err
			}

			d.SetId(nodeId)
		}
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteServiceGraphNodeRead(d, m)
}
func resourceMSOSchemaSiteServiceGraphNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	var templateName string
	if tempVar, ok := d.GetOk("template_name"); ok {
		templateName = tempVar.(string)
	}

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	nodeId := d.Id()

	templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/%s/serviceNodes/%s", templateName, graphName, nodeId)
	templatePatchStruct := models.NewTemplateServiceGraph("remove", templatePath, nil)
	sitePayload := make([]models.Model, 0, 1)

	sitePayload = append(sitePayload, templatePatchStruct)
	var siteParams []interface{}
	if tempVar, ok := d.GetOk("site_nodes"); ok {
		siteParams = tempVar.([]interface{})

		for _, site := range siteParams {
			siteMap := site.(map[string]interface{})
			if siteMap["site_id"] == "" {
				return fmt.Errorf("site_id is required in site_nodes list")
			}

			graphCont, ind, err := getSiteServiceGraphCont(
				cont,
				schemaId,
				templateName,
				siteMap["site_id"].(string),
				graphName,
			)

			if err == nil {
				_, nodeind, err := getSiteServiceNodeCont(
					graphCont,
					schemaId,
					templateName,
					graphName,
					nodeId,
				)

				if err == nil {
					sitePath := fmt.Sprintf(
						"/sites/%s-%s/serviceGraphs/%d/serviceNodes/%d",
						siteMap["site_id"].(string),
						templateName,
						ind,
						nodeind,
					)
					sitePayload = append(sitePayload, models.NewTemplateServiceGraph("remove", sitePath, nil))
				}
			}
		}
	}

	response, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), sitePayload...)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")

	return nil
}
