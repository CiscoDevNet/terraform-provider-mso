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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"site_nodes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

	nodeName := models.StripQuotes(tempNodeCont.S("name").String())

	nodeIdHuman, err := getNodeNameFromId(msoClient, nodeId)

	if err != nil {
		return nil, err
	}

	d.Set("service_node_type", nodeIdHuman)

	siteParams := make([]interface{}, 0, 1)

	sitesCount, err := cont.ArrayCount("sites")

	if err != nil {
		d.SetId(graphName)
		d.Set("site_nodes", nil)
		log.Printf("Unable to find sites")
		return nil, err

	}

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, fmt.Errorf("Unable to load site element")
		}

		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			schemaId,
			stateTemplate,
			apiSiteId,
			graphName,
		)

		if err == nil {
			nodeCont, _, nodeerr := getSiteServiceNodeCont(
				graphCont,
				schemaId,
				stateTemplate,
				graphName,
				nodeName,
			)

			if nodeerr == nil {
				siteMap := make(map[string]interface{})

				deviceDn := models.StripQuotes(nodeCont.S("device", "dn").String())

				dnSplit := strings.Split(deviceDn, "/")

				tnName := strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")
				siteMap["tenant_name"] = tnName
				siteMap["node_name"] = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")
				siteMap["site_id"] = apiSiteId

				siteParams = append(siteParams, siteMap)
			}

		}

	}

	d.Set("site_nodes", siteParams)
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
	var nodeType string
	if tempVar, ok := d.GetOk("service_node_type"); ok {
		nodeType = tempVar.(string)
	}

	var desc string
	if tempVar, ok := d.GetOk("description"); ok {
		desc = tempVar.(string)
	}

	nodeId, err := getNodeIdFromName(msoClient, nodeType)
	if err != nil {
		return err
	}

	templatePayload := make(map[string]interface{})
	templatePayload["name"] = graphName
	templatePayload["displayName"] = graphName
	templatePayload["description"] = desc

	serviceNode := map[string]interface{}{
		"name":              "tfnode1",
		"index":             1,
		"serviceNodeTypeId": nodeId,
	}
	templatePayload["serviceNodes"] = []interface{}{serviceNode}
	templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/-", templateName)
	templatePatchStruct := models.NewTemplateServiceGraph("add", templatePath, templatePayload)

	var siteParams []interface{}
	sitePayload := make([]models.Model, 0, 1)

	sitePayload = append(sitePayload, templatePatchStruct)

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
				"serviceGraphRef": map[string]interface{}{
					"schemaId":         schemaId,
					"templateName":     templateName,
					"serviceGraphName": graphName,
				},
				"serviceNodes": []interface{}{
					map[string]interface{}{
						"serviceNodeRef": map[string]interface{}{
							"schemaId":         schemaId,
							"templateName":     templateName,
							"serviceGraphName": graphName,
							"serviceNodeName":  "tfnode1",
						},
						"device": map[string]interface{}{
							"dn": fmt.Sprintf("uni/tn-%s/lDevVip-%s", siteMap["tenant_name"].(string), siteMap["node_name"].(string)),
						},
					},
				},
			}
			// ----> site payload creation ends
			sitePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/-", siteMap["site_id"].(string), templateName)

			sitePayload = append(sitePayload, models.NewTemplateServiceGraph("add", sitePath, siteVarMap))

		}
	}

	// PATCH the payload

	_, err = msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaId), sitePayload...)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", graphName))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateServiceGraphRead(d, m)

}

func resourceMSOSchemaTemplateServiceGraphRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

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
	d.Set("description", models.StripQuotes(sgCont.S("description").String()))

	_, _, err = getTemplateServiceNodeCont(sgCont, "tfnode1", nodeId)

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
				"tfnode1",
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
						"serviceNodeName":  "tfnode1",
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
					"tfnode1",
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

			d.SetId(fmt.Sprintf("%v", graphName))
		}
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaTemplateServiceGraphRead(d, m)
}

func resourceMSOSchemaTemplateServiceGraphDelete(d *schema.ResourceData, m interface{}) error {
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

	templatePath := fmt.Sprintf("/templates/%s/serviceGraphs/%s", templateName, graphName)
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

			_, ind, err := getSiteServiceGraphCont(
				cont,
				schemaId,
				templateName,
				siteMap["site_id"].(string),
				graphName,
			)

			if err == nil {
				sitePath := fmt.Sprintf("/sites/%s-%s/serviceGraphs/%d", siteMap["site_id"].(string), templateName, ind)
				sitePayload = append(sitePayload, models.NewTemplateServiceGraph("remove", sitePath, nil))
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
					log.Printf("tppppp %v len %d", graphEle, len(graphEle))
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
