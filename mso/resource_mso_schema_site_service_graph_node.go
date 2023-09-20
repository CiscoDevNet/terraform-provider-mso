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
			"service_node_type": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"service_node"},
				Deprecated:    "Use service_node to configure service node devices.",
			},

			"site_nodes": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"service_node"},
				Deprecated:    "Use service_node to configure service nodes devices.",
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
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
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

	var siteId string
	if site_id, ok := d.GetOk("site_id"); ok {
		siteId = site_id.(string)
	}
	log.Printf("CHECK CREATE siteId %v", siteId)

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
	log.Printf("CHECK CREATE graphcont %v", grapCont)
	if err != nil {
		return err
	}

	siteServiceGraphPayload := make([]models.Model, 0, 1)

	if siteServiceNodes, ok := d.GetOk("service_node"); ok { // New code
		templateNodeCount, err := grapCont.ArrayCount("serviceNodes")
		if err != nil {
			return fmt.Errorf("Unable to count template service node")
		}
		log.Printf("CHECK CREATE templateNodeCount %v", templateNodeCount)
		log.Printf("CHECK CREATE siteServiceNodes %v", siteServiceNodes)

		templateServiceNodes := grapCont.S("serviceNodes").Data().([]interface{})
		log.Printf("CHECK CREATE templateServiceNodes %v", templateServiceNodes)
		siteServiceNodeList := make([]interface{}, 0, 1)
		for index, templateServiceNode := range templateServiceNodes {
			templateServiceNodeMap := templateServiceNode.(map[string]interface{})
			log.Printf("CHECK CREATE serviceNodeValues %v", templateServiceNodeMap)
			log.Printf("CHECK CREATE serviceNodeRef %v", templateServiceNodeMap["serviceNodeRef"])
			siteServiceNodeDeviceList := siteServiceNodes.([]interface{})
			siteServiceNodeMap := siteServiceNodeDeviceList[index].(map[string]interface{})
			log.Printf("CHECK CREATE siteServiceNodeMap %v", siteServiceNodeMap)

			serviceNodeMap := map[string]interface{}{
				"serviceNodeRef": templateServiceNodeMap["serviceNodeRef"],
				"device": map[string]interface{}{
					"dn": siteServiceNodeMap["device_dn"],
				},
			}
			log.Printf("CHECK CREATE serviceNodeMap %v", serviceNodeMap)
			siteServiceNodeList = append(siteServiceNodeList, serviceNodeMap)

		}
		log.Printf("CHECK CREATE siteServiceNodeList %v", siteServiceNodeList)

		serviceNodePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%s/serviceNodes", siteId, templateName, graphName)
		log.Printf("CHECK CREATE serviceNodePath %v", serviceNodePath)

		siteServiceGraphPayload := models.NewSchemaSiteServiceGraph("add", serviceNodePath, siteServiceNodeList)
		_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload)

		if err != nil {
			return err
		}

	} else { // To be removed in future
		nodeInd := getTemplateNodeIndex(grapCont)
		log.Printf("CHECK CREATE nodeInd %v", nodeInd)

		if nodeInd == -1 {
			return fmt.Errorf("Unable to get Template node index list")
		}

		templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/%s/serviceNodes/-", templateName, graphName)

		templatePayload := map[string]interface{}{
			"name":              fmt.Sprintf("node%d", nodeInd),
			"index":             nodeInd,
			"serviceNodeTypeId": nodeId,
		}

		templatePatchStruct := models.NewTemplateServiceGraph("add", templatePath, templatePayload)

		siteServiceGraphPayload = append(siteServiceGraphPayload, templatePatchStruct)

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

				// <---- Begin site payload creation -
				siteVarMap := map[string]interface{}{
					"serviceNodeRef": map[string]interface{}{
						"schemaId":         schemaId,
						"templateName":     templateName,
						"serviceGraphName": graphName,
						"serviceNodeName":  fmt.Sprintf("node%d", nodeInd),
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

				siteServiceGraphPayload = append(siteServiceGraphPayload, models.NewTemplateServiceGraph("add", sitePath, siteVarMap))
				_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload...)

				if err != nil {
					return err
				}
				log.Printf("CHECK IGNORE %v %v", sitePath, siteVarMap)
			}
		}

		d.SetId(fmt.Sprintf("node%d", nodeInd))
		log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	}

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
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	log.Printf("CHECK READ cont %v", cont)

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", templateName)
	d.Set("site_id", siteId)
	d.Set("service_graph_name", graphName)

	// To be removed in future.
	if tempVar, ok := d.GetOk("service_node_type"); ok {
		nodeType := tempVar.(string)
		nodeId, err := getNodeIdFromName(msoClient, nodeType)
		if err != nil {
			return err
		}

		sgCont, _, err := getTemplateServiceGraphCont(cont, templateName, graphName)
		if err != nil {
			d.SetId("")
			log.Printf("graphcont err %v", err)
			return nil
		}

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
					templateName,
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
					templateName,
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
		}
	} else { // New code
		log.Printf("---------- CHECK READ ----------")
		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			schemaId,
			templateName,
			siteId,
			graphName,
		)
		log.Printf("CHECK READ graphCont %v", graphCont)
		if err != nil {
			d.SetId("")
			log.Printf("sitegraphcont err %v", err)
			return nil
		}

		serviceNodeList := make([]interface{}, 0, 1)
		serviceNodes := graphCont.S("serviceNodes").Data().([]interface{})
		log.Printf("CHECK READ serviceNodes %v", serviceNodes)
		for _, val := range serviceNodes {
			serviceNodeValues := val.(map[string]interface{})
			serviceNodeMap := make(map[string]interface{})
			log.Printf("CHECK READ serviceNodeValues %v", serviceNodeValues)
			device := serviceNodeValues["device"].(map[string]interface{})["dn"]
			log.Printf("CHECK READ device %v", device)
			serviceNodeMap["device_dn"] = device
			serviceNodeList = append(serviceNodeList, serviceNodeMap)
			log.Printf("CHECK READ serviceNodeList %v", serviceNodeList)
		}
		d.Set("service_node", serviceNodeList)
		d.Set("site_nodes", nil)
		d.Set("service_node_type", nil)
	}

	d.SetId(nodeIdSt)
	log.Printf("CHECK READ dddddddd %v", d)
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
		// siteServiceGraphPayload := make([]models.Model, 0, 1)

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

				// siteServiceGraphPayload = append(siteServiceGraphPayload, models.NewTemplateServiceGraph("replace", sitePath, siteVarMap))
				log.Printf("CHECK IGNORE %v %v", sitePath, siteVarMap)

			}
			// _, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload...)

			// if err != nil {
			// 	return err
			// }

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
	siteServiceGraphPayload := make([]models.Model, 0, 1)

	siteServiceGraphPayload = append(siteServiceGraphPayload, templatePatchStruct)
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
					// siteServiceGraphPayload = append(siteServiceGraphPayload, models.NewTemplateServiceGraph("remove", sitePath, nil))
					log.Printf("CHECK IGNORE %v ", sitePath)
				}
			}
		}
	}

	response, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), siteServiceGraphPayload...)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")

	return nil
}
