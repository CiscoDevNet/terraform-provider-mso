package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateServiceGraphs() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateServiceGraphCreate,
		Read:   resourceMSOSchemaTemplateServiceGraphRead,
		Update: resourceMSOSchemaTemplateServiceGraphUpdate,
		Delete: resourceMSOSchemaTemplateServiceGraphDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateServiceGraphImport,
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
			"service_graph_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				// ConflictsWith: []string{"service_node"},
				Deprecated: "Use service_node to configure service nodes.",
			},
			"service_node": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configure service nodes for the service graph.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaTemplateServiceGraphImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	stateTemplate := get_attribute[2]
	var graphName string
	graphName = get_attribute[4]

	sgCont, _, err := getTemplateServiceGraphCont(cont, stateTemplate, graphName)

	if err != nil {
		d.SetId("")
		log.Printf("graphcont err %v", err)
		return nil, err
	}

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)
	nodeInd, err := strconv.Atoi(get_attribute[6])
	if err != nil {
		return nil, err
	}
	tempNodeCont, _, err := getTemplateServiceNodeContFromIndex(sgCont, nodeInd)

	if err != nil {
		d.SetId("")
		return nil, err
	}

	d.Set("node_index", nodeInd)

	nodeId := models.StripQuotes(tempNodeCont.S("serviceNodeTypeId").String())

	// nodeName := models.StripQuotes(tempNodeCont.S("name").String())

	nodeIdHuman, err := getNodeNameFromId(msoClient, nodeId)

	if err != nil {
		return nil, err
	}

	d.Set("service_node_type", nodeIdHuman)

	d.SetId(graphName)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateServiceGraphCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Creation Template Service Graph")
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

	var desc string
	if tempVar, ok := d.GetOk("description"); ok {
		desc = tempVar.(string)
	}

	templatePayload := make(map[string]interface{})
	templatePayload["name"] = graphName
	templatePayload["displayName"] = graphName
	templatePayload["description"] = desc

	serviceNodes := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("service_node"); ok {
		log.Printf("CHECK val %v", val)
		service_node_list := val.([]interface{})
		log.Printf("CHECK service_node_list %v", service_node_list)
		for i, val := range service_node_list {
			log.Printf("CHECK FOR  vals %v", val)
			log.Printf("CHECK FOR  i %v", i)

			service_node_map := make(map[string]interface{})
			service_node_values := val.(map[string]interface{})
			if service_node_values["type"] != "" {
				nodeId, err := getNodeIdFromName(msoClient, fmt.Sprintf("%v", service_node_values["type"]))
				if err != nil {
					return err
				}
				log.Printf("CHECK nodeId %v", nodeId)
				service_node_map["serviceNodeTypeId"] = nodeId
				index := i + 1
				service_node_map["index"] = index
				service_node_map["name"] = fmt.Sprintf("node%v", index)

			}
			serviceNodes = append(serviceNodes, service_node_map)
			log.Printf("CHECK serviceNodes %v", serviceNodes)
		}
	}
	templatePayload["serviceNodes"] = serviceNodes

	templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/-", templateName)
	templatePatchStruct := models.NewTemplateServiceGraph("add", templatePath, templatePayload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), templatePatchStruct)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", graphName))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateServiceGraphRead(d, m)

}

func resourceMSOSchemaTemplateServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	log.Printf("CHECK READ ")

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	log.Printf("CHECK READ cont %v", cont)

	stateTemplate := d.Get("template_name").(string)
	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	log.Printf("CHECK READ  graphName %v", graphName)
	sgCont, _, err := getTemplateServiceGraphCont(cont, stateTemplate, graphName)
	if err != nil {
		d.SetId("")
		log.Printf("graphcont err %v", err)
		return nil
	}
	log.Printf("CHECK READ sgCont %v", sgCont)

	serviceNodes := sgCont.S("serviceNodes").Data().([]interface{})
	log.Printf("CHECK READ serviceNodes %v", serviceNodes)

	serviceNodeList := make([]interface{}, 0, 1)
	for _, val := range serviceNodes {
		serviceNodeValues := val.(map[string]interface{})
		log.Printf("CHECK READ serviceNodeValues %v", serviceNodeValues)
		serviceNodeMap := make(map[string]interface{})
		nodeId := models.StripQuotes(serviceNodeValues["serviceNodeTypeId"].(string))
		log.Printf("CHECK READ nodeId %v", nodeId)

		nodeType, err := getNodeNameFromId(msoClient, nodeId)
		if err != nil {
			return err
		}
		log.Printf("CHECK READ nodeType %v", nodeType)
		serviceNodeMap["type"] = nodeType

		serviceNodeList = append(serviceNodeList, serviceNodeMap)
	}
	log.Printf("CHECK READ serviceNodeList %v", serviceNodeList)
	d.Set("service_node", serviceNodeList)

	d.Set("schema_id", schemaId)
	d.Set("template_name", stateTemplate)
	d.Set("service_graph_name", graphName)
	d.Set("description", models.StripQuotes(sgCont.S("description").String()))

	count2, _ := sgCont.ArrayCount("serviceNodes")
	log.Printf("CHECK READ count2 %v", count2)
	if err != nil {
		d.Set("service_node", make([]interface{}, 0))
	}

	service_node := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		snCont, err := sgCont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return fmt.Errorf("Unable to parse the user associations list")
		}
		log.Printf("CHECK READ snCont %v", snCont)

		mapUser := make(map[string]interface{})
		mapUser["user_id"] = models.StripQuotes(snCont.S("userId").String())
		service_node = append(service_node, mapUser)
	}

	d.Set("service_node", service_node)

	d.SetId(graphName)
	return nil
}

func resourceMSOSchemaTemplateServiceGraphUpdate(d *schema.ResourceData, m interface{}) error {
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

	var graphName string
	if tempVar, ok := d.GetOk("service_graph_name"); ok {
		graphName = tempVar.(string)
	}

	if d.HasChange("description") {
		var desc string
		if tempVar, ok := d.GetOk("description"); ok {
			desc = tempVar.(string)

		} else {
			desc = ""
		}

		templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/%s/description", templateName, graphName)

		graphUpdate := models.NewTemplateServiceGraphUpdate("replace", templatePath, desc)

		_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), graphUpdate)

		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("%v", graphName))

	}

	// if d.HasChange("site_nodes") {
	// 	var siteParams []interface{}
	// 	sitePayload := make([]models.Model, 0, 1)

	// 	if tempVar, ok := d.GetOk("site_nodes"); ok {
	// 		siteParams = tempVar.([]interface{})

	// 		for _, site := range siteParams {
	// 			siteMap := site.(map[string]interface{})
	// 			if siteMap["site_id"] == "" {
	// 				return fmt.Errorf("site_id is required in site_nodes list")
	// 			}

	// 			if siteMap["tenant_name"] == "" {
	// 				return fmt.Errorf("tenant_name is required in site_nodes list")
	// 			}

	// 			if siteMap["node_name"] == "" {
	// 				return fmt.Errorf("node_name is required in site_nodes list")
	// 			}

	// 			// <---- Begin site payload creation
	// 			siteVarMap := map[string]interface{}{
	// 				"serviceNodeRef": map[string]interface{}{
	// 					"schemaId":         schemaId,
	// 					"templateName":     templateName,
	// 					"serviceGraphName": graphName,
	// 					"serviceNodeName":  "tfnode1",
	// 				},
	// 				"device": map[string]interface{}{
	// 					"dn": fmt.Sprintf("uni/tn-%s/lDevVip-%s", siteMap["tenant_name"].(string), siteMap["node_name"].(string)),
	// 				},
	// 			}

	// 			// ----> site payload creation ends

	// 			cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	// 			if err != nil {
	// 				return err
	// 			}

	// 			graphCont, graphind, err := getSiteServiceGraphCont(
	// 				cont,
	// 				schemaId,
	// 				templateName,
	// 				siteMap["site_id"].(string),
	// 				graphName,
	// 			)

	// 			if err != nil {
	// 				return err
	// 			}

	// 			_, nodeind, err := getSiteServiceNodeCont(
	// 				graphCont,
	// 				schemaId,
	// 				templateName,
	// 				graphName,
	// 				"tfnode1",
	// 			)

	// 			if err != nil {
	// 				return err
	// 			}

	// 			sitePath := fmt.Sprintf(
	// 				"/sites/%s-%s/serviceGraphs/%d/serviceNodes/%d",
	// 				siteMap["site_id"].(string),
	// 				templateName,
	// 				graphind,
	// 				nodeind,
	// 			)

	// 			sitePayload = append(sitePayload, models.NewTemplateServiceGraph("replace", sitePath, siteVarMap))

	// 		}
	// 		_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), sitePayload...)

	// 		if err != nil {
	// 			return err
	// 		}

	// 		d.SetId(fmt.Sprintf("%v", graphName))
	// 	}
	// }
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaTemplateServiceGraphRead(d, m)
}

func resourceMSOSchemaTemplateServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	_, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
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

	path := fmt.Sprintf("/templates/%s/serviceGraphs/%s", templateName, graphName)
	log.Printf("CHECK DELETE path  %v", path)

	templatePatchStruct := models.NewTemplateServiceGraph("remove", path, nil)
	log.Printf("CHECK DELETE templatePatchStruct  %v", templatePatchStruct)

	response, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), templatePatchStruct)
	log.Printf("CHECK DELETE templatePatchStruct  %v", response)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")

	return nil
}

func getSiteServiceNodeCont(graphCont *container.Container, schemaId, templateName, graphName, nodeName string) (*container.Container, int, error) {

	nodesCount, err := graphCont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load count site service node")
	}
	for i := 0; i < nodesCount; i++ {
		nodeCont, err := graphCont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to site service node element")
		}

		nodeRef := models.StripQuotes(nodeCont.S("serviceNodeRef").String())

		nodeSplit := strings.Split(nodeRef, "/")
		if len(nodeSplit) == 9 {
			if nodeSplit[2] == schemaId && nodeSplit[4] == templateName && nodeSplit[6] == graphName && nodeSplit[8] == nodeName {
				return nodeCont, i, nil

			}
		} else {
			return nil, -1, fmt.Errorf("Spilt on nodeRef failed")
		}
	}
	return nil, -1, fmt.Errorf("Unable to find site service node")
}
func getSiteServiceGraphCont(cont *container.Container, schemaId, templateName, siteId, graphName string) (*container.Container, int, error) {
	sitesCount, err := cont.ArrayCount("sites")

	if err != nil {
		return nil, -1, fmt.Errorf("Unable to find sites")
	}

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load site element")
		}

		siteTemplate := models.StripQuotes(siteCont.S("templateName").String())
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if siteTemplate == templateName && siteId == apiSiteId {
			sgCount, err := siteCont.ArrayCount("serviceGraphs")
			if err != nil {
				return nil, -1, fmt.Errorf("Unable to load site service graphs")
			}

			for j := 0; j < sgCount; j++ {
				sgCont, err := siteCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to load site service graph element")
				}

				graphRef := models.StripQuotes(sgCont.S("serviceGraphRef").String())

				graphEle := strings.Split(graphRef, "/")

				if len(graphEle) != 7 {
					// log.Printf("tppppp %v len %d", graphEle, len(graphEle))
					return nil, -1, fmt.Errorf("Inavlid site service graph")
				}

				if schemaId == graphEle[2] && templateName == graphEle[4] && graphName == graphEle[6] {
					return sgCont, j, nil
				}

			}
		}
	}

	return nil, -1, fmt.Errorf("Unable to find site service graph")
}

func getTemplateServiceNodeCont(cont *container.Container, nodeName, nodeType string) (*container.Container, int, error) {

	nodeCount, err := cont.ArrayCount("serviceNodes")
	if err != nil {
		return nil, -1, fmt.Errorf("Unable to load node count")
	}

	for i := 0; i < nodeCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodes")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to load node element")
		}

		apiNodeName := models.StripQuotes(nodeCont.S("name").String())
		apiNodeType := models.StripQuotes(nodeCont.S("serviceNodeTypeId").String())

		if apiNodeName == nodeName && apiNodeType == nodeType {
			return nodeCont, i, nil
		}

	}

	return nil, -1, fmt.Errorf("Unable to find the service node")
}

func getTemplateServiceGraphCont(cont *container.Container, templateName, graphName string) (*container.Container, int, error) {
	templateCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, -1, fmt.Errorf("No Template found")
	}

	for i := 0; i < templateCount; i++ {
		templateCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to get template element")
		}

		apiTemplate := models.StripQuotes(templateCont.S("name").String())

		if apiTemplate == templateName {
			log.Printf("[DEBUG] Template found")

			sgCount, err := templateCont.ArrayCount("serviceGraphs")

			if err != nil {
				return nil, -1, fmt.Errorf("No Service Graph found")
			}

			for j := 0; j < sgCount; j++ {
				sgCont, err := templateCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to get service graph element")
				}

				apiSgName := models.StripQuotes(sgCont.S("name").String())

				if apiSgName == graphName {
					return sgCont, j, nil
				}
			}

		}
	}

	return nil, -1, fmt.Errorf("unable to find service graph")
}
func getNodeIdFromName(msoClient *client.Client, nodeType string) (string, error) {
	cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return "", err
	}

	nodesCount, err := cont.ArrayCount("serviceNodeTypes")
	if err != nil {
		return "", err
	}

	for i := 0; i < nodesCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
		if err != nil {
			return "", err
		}

		apiName := models.StripQuotes(nodeCont.S("name").String())

		if apiName == nodeType {
			return models.StripQuotes(nodeCont.S("id").String()), nil
		}
	}

	return "", fmt.Errorf("Unable to find nodeid for nodetype %s", nodeType)
}

func getNodeNameFromId(msoClient *client.Client, nodeId string) (string, error) {
	cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return "", err
	}

	nodesCount, err := cont.ArrayCount("serviceNodeTypes")
	if err != nil {
		return "", err
	}

	for i := 0; i < nodesCount; i++ {
		nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
		if err != nil {
			return "", err
		}

		apiId := models.StripQuotes(nodeCont.S("id").String())

		if apiId == nodeId {
			return models.StripQuotes(nodeCont.S("name").String()), nil
		}
	}

	return "", fmt.Errorf("Unable to find nodeNamefor nodeid %s", nodeId)
}
