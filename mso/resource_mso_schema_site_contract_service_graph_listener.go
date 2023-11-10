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

func resourceMSOSchemaSiteContractServiceGraphListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteContractServiceGraphListenerCreate,
		Update: resourceMSOSchemaSiteContractServiceGraphListenerUpdate,
		Read:   resourceMSOSchemaSiteContractServiceGraphListenerRead,
		Delete: resourceMSOSchemaSiteContractServiceGraphListenerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteContractServiceGraphListenerImport,
		},

		Schema: map[string]*schema.Schema{
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
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node_index": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"listener_name": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"http",
					"https",
					"tcp",
					"udp",
					"tls",
					"inherit",
				}, false),
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
					"eLBSecurityPolicy-2016-08",
					"eLBSecurityPolicy-FS-2018-06",
					"eLBSecurityPolicy-TLS-1-2-2017-01",
					"eLBSecurityPolicy-TLS-1-2-Ext-2018-06",
					"eLBSecurityPolicy-TLS-1-1-2017-01",
					"eLBSecurityPolicy-2015-05",
					"eLBSecurityPolicy-TLS-1-0-2015-04",
					"AppGwSslPolicyDefault",
					"AppGwSslPolicy20150501",
					"AppGwSslPolicy20170401",
					"AppGwSslPolicy20170401S",
				}, false),
			},
			"ssl_certificates": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"target_dn": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"default": &schema.Schema{ // default should be true for a valid SSL Certificate
							Type:     schema.TypeBool,
							Computed: true,
						},
						"certificate_store": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"default",
								"iam",
								"acm",
							}, false),
						},
					},
				},
			},
			"frontend_ip_dn": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"rules": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"floating_ip": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"host": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"path": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"action": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"condition": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"action_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"fixedResponse",
								"forward",
								"redirect",
								"haPort",
							}, false),
						},
						"content_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"textPlain",
								"textCSS",
								"textHtml",
								"appJS",
								"appJson",
							}, false),
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
								"tcp",
								"udp",
								"tls",
								"inherit",
							}, false),
						},
						"provider_epg_ref": { // Only one object allowed
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schema_id": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"template_name": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
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
								},
							},
						},
						"url_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"original",
								"custom",
							}, false),
						},
						"custom_url": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"redirect_host_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"redirect_path": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"redirect_query": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"response_code": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"response_body": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"redirect_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
								"tcp",
								"udp",
								"tls",
								"inherit",
							}, false),
						},
						"redirect_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"redirect_code": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"unknown",
								"permMoved",
								"found",
								"seeOther",
								"temporary",
							}, false),
						},
						"health_check": &schema.Schema{ // Only one object allowed
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"http",
											"https",
											"tcp",
											"udp",
											"tls",
											"inherit",
										}, false),
									},
									"path": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"interval": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"timeout": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"unhealthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"use_host_from_rule": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"yes",
											"no",
										}, false),
									},
									"success_code": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"host": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
								},
							},
						},
						"target_ip_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"unspecified",
								"primary",
								"secondary",
							}, false),
						},
					},
				},
			},
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			oldRules, newRules := diff.GetChange("rules")
			oldRulesList := oldRules.(*schema.Set).List()
			newRulesList := newRules.(*schema.Set).List()

			if len(oldRulesList) != len(newRulesList) {
				return nil
			}

			for _, newRule := range newRulesList {

				for _, oldRules := range oldRulesList {

					newRuleMap := newRule.(map[string]interface{})
					newRuleName := newRuleMap["name"].(string)
					newRulePriority := newRuleMap["priority"].(int)
					oldRuleMap := oldRules.(map[string]interface{})
					oldRuleName := oldRuleMap["name"].(string)
					oldRulePriority := oldRuleMap["priority"].(int)

					if newRuleName == oldRuleName && newRulePriority == oldRulePriority {

						newRuleHealthCheckList := newRuleMap["health_check"].(*schema.Set).List()
						oldRuleHealthCheckList := oldRuleMap["health_check"].(*schema.Set).List()

						healCheckDiffBoolList := make([]bool, 0)
						if len(newRuleHealthCheckList) != 0 && len(oldRuleHealthCheckList) != 0 {
							newRuleHealthCheckMap := newRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})
							oldRuleHealthCheckMap := oldRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})

							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["host"].(string) == "" && oldRuleHealthCheckMap["host"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["interval"].(int) == 0 && oldRuleHealthCheckMap["interval"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["path"].(string) == "" && oldRuleHealthCheckMap["path"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["port"].(int) == 0 && oldRuleHealthCheckMap["port"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["protocol"].(string) == "" && oldRuleHealthCheckMap["protocol"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, (newRuleHealthCheckMap["success_code"].(string) == "200-399" && oldRuleHealthCheckMap["success_code"].(string) == "200-399") || (newRuleHealthCheckMap["success_code"].(string) == "" && oldRuleHealthCheckMap["success_code"].(string) == ""))
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["timeout"].(int) == 0 && oldRuleHealthCheckMap["timeout"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, newRuleHealthCheckMap["unhealthy_threshold"].(int) == 0 && oldRuleHealthCheckMap["unhealthy_threshold"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, (newRuleHealthCheckMap["use_host_from_rule"].(string) == "no" && oldRuleHealthCheckMap["use_host_from_rule"].(string) == "no") || (newRuleHealthCheckMap["use_host_from_rule"].(string) == "" && oldRuleHealthCheckMap["use_host_from_rule"].(string) == ""))

						} else if len(newRuleHealthCheckList) == 0 && len(oldRuleHealthCheckList) != 0 {
							// When the newRuleHealthCheckList length is 0 and oldRuleHealthCheckList length is not 0
							oldRuleHealthCheckMap := oldRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})

							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["host"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["interval"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["path"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["port"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["protocol"].(string) == "")
							healCheckDiffBoolList = append(healCheckDiffBoolList, (oldRuleHealthCheckMap["success_code"].(string) == "200-399" || oldRuleHealthCheckMap["success_code"].(string) == ""))
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["timeout"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, oldRuleHealthCheckMap["unhealthy_threshold"].(int) == 0)
							healCheckDiffBoolList = append(healCheckDiffBoolList, (oldRuleHealthCheckMap["use_host_from_rule"].(string) == "no" || oldRuleHealthCheckMap["use_host_from_rule"].(string) == ""))
						}

						healCheckDiffFalse := false
						for _, healCheckDiffBool := range healCheckDiffBoolList {
							healCheckDiffFalse = healCheckDiffFalse || !healCheckDiffBool
						}
						if !healCheckDiffFalse {
							// Clearing the rules diff
							diff.Clear("rules")
						}
					}
				}
			}
			return nil
		},
	}
}

func resourceMSOSchemaSiteContractServiceGraphListenerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	ListenerTokens := strings.Split(d.Id(), "/")
	d.Set("schema_id", ListenerTokens[0])
	d.Set("site_id", ListenerTokens[2])
	d.Set("template_name", ListenerTokens[4])
	d.Set("contract_name", ListenerTokens[6])

	serviceNodeIndex, err := strconv.Atoi(ListenerTokens[8])
	if err == nil {
		d.Set("service_node_index", serviceNodeIndex)
	}

	d.Set("listener_name", ListenerTokens[10])

	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", ListenerTokens[0]))
	if err != nil {
		return nil, err
	}
	err = setSchemaSiteContractServiceGraphListenerAttrs(cont, d)
	if err != nil {
		return nil, err
	}
	d.SetId(d.Id())
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteContractServiceGraphListenerCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Contract Service Graph Listener: Beginning Creation")
	err := postSchemaSiteContractServiceGraphListenerConfig("add", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceMSOSchemaSiteContractServiceGraphListenerRead(d, m)
}

func resourceMSOSchemaSiteContractServiceGraphListenerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Contract Service Graph Listener: Beginning Update")
	err := postSchemaSiteContractServiceGraphListenerConfig("replace", d, m)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteContractServiceGraphListenerRead(d, m)
}

func resourceMSOSchemaSiteContractServiceGraphListenerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Read Site Contract Service Graph Listener")
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSchemaSiteContractServiceGraphListenerAttrs(cont, d)

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Completed Read Site Contract Service Graph Listener")
	return nil
}

func resourceMSOSchemaSiteContractServiceGraphListenerDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Begining Delete Site Contract Service Graph Listener")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceNodeIndex := d.Get("service_node_index").(int)
	listenerName := d.Get("listener_name").(string)

	listenerPath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship/serviceNodesRelationship/%d/deviceConfiguration/cloudLoadBalancer/listeners/%s",
		siteID, templateName, contractName, serviceNodeIndex, listenerName,
	)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), models.GetRemovePatchPayload(listenerPath))

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] Completed Delete Site Contract Service Graph Listener")
	return nil
}

func setSchemaSiteContractServiceGraphListenerAttrs(cont *container.Container, d *schema.ResourceData) error {
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	serviceNodeIndex := d.Get("service_node_index").(int)
	listenerName := d.Get("listener_name").(string)

	sitesCont := cont.S("sites")
	for _, siteCont := range sitesCont.Data().([]interface{}) {

		apiSiteID := models.StripQuotes(siteCont.(map[string]interface{})["siteId"].(string))
		apiTemplateName := models.StripQuotes(siteCont.(map[string]interface{})["templateName"].(string))
		if siteID == apiSiteID && templateName == apiTemplateName {

			siteContractsCont := siteCont.(map[string]interface{})["contracts"]
			for _, contractCont := range siteContractsCont.([]interface{}) {

				contractRefTokens := strings.Split(models.StripQuotes(contractCont.(map[string]interface{})["contractRef"].(string)), "/")
				apiContractName := contractRefTokens[len(contractRefTokens)-1]
				if contractName == apiContractName {

					serviceGraphRelationship := contractCont.(map[string]interface{})["serviceGraphRelationship"]
					deviceConfigurationMap := serviceGraphRelationship.(map[string]interface{})["serviceNodesRelationship"].([]interface{})[serviceNodeIndex].(map[string]interface{})["deviceConfiguration"]
					if deviceConfigurationMap != nil {

						listenersMap := deviceConfigurationMap.(map[string]interface{})["cloudLoadBalancer"].(map[string]interface{})["listeners"]
						for _, listener := range listenersMap.([]interface{}) {
							listenerMap := listener.(map[string]interface{})
							if listenerName == listenerMap["name"].(string) {

								d.Set("protocol", listenerMap["protocol"].(string))
								d.Set("port", listenerMap["port"])
								d.Set("security_policy", convertInterfaceToString(listenerMap["secPolicy"]))

								if listenerMap["certificates"] != nil {
									sslCertificates := make([]map[string]interface{}, 0)
									for _, certificate := range listenerMap["certificates"].([]interface{}) {
										certificateMap := certificate.(map[string]interface{})
										sslCertMap := map[string]interface{}{
											"name":              convertInterfaceToString(certificateMap["name"]),
											"target_dn":         convertInterfaceToString(certificateMap["tDn"]),
											"default":           certificateMap["default"].(bool),
											"certificate_store": convertInterfaceToString(certificateMap["store"]),
										}
										sslCertificates = append(sslCertificates, sslCertMap)
									}
									d.Set("ssl_certificates", sslCertificates)
								}

								if listenerMap["nlbDevIp"] != nil {
									frontendIpDn := convertInterfaceToString(listenerMap["nlbDevIp"].(map[string]interface{})["dn"])
									d.Set("frontend_ip_dn", frontendIpDn)
								}

								rules := make([]map[string]interface{}, 0)
								for _, apiRule := range listenerMap["rules"].([]interface{}) {
									apiRuleMap := apiRule.(map[string]interface{})

									ruleMap := map[string]interface{}{
										"name":               convertInterfaceToString(apiRuleMap["name"]),
										"floating_ip":        convertInterfaceToString(apiRuleMap["floatingIp"]),
										"priority":           convertInterfaceToInt(apiRuleMap["index"]),
										"host":               convertInterfaceToString(apiRuleMap["host"]),
										"path":               convertInterfaceToString(apiRuleMap["path"]),
										"action":             convertInterfaceToString(apiRuleMap["action"]),
										"condition":          convertInterfaceToString(apiRuleMap["condition"]),
										"action_type":        convertInterfaceToString(apiRuleMap["actionType"]),
										"content_type":       convertInterfaceToString(apiRuleMap["contentType"]),
										"port":               convertInterfaceToInt(apiRuleMap["port"]),
										"protocol":           convertInterfaceToString(apiRuleMap["protocol"]),
										"url_type":           convertInterfaceToString(apiRuleMap["urlType"]),
										"custom_url":         convertInterfaceToString(apiRuleMap["customURL"]),
										"redirect_host_name": convertInterfaceToString(apiRuleMap["redirectHostName"]),
										"redirect_path":      convertInterfaceToString(apiRuleMap["redirectPath"]),
										"redirect_query":     convertInterfaceToString(apiRuleMap["redirectQuery"]),
										"response_code":      convertInterfaceToString(apiRuleMap["responseCode"]),
										"response_body":      convertInterfaceToString(apiRuleMap["responseBody"]),
										"redirect_protocol":  convertInterfaceToString(apiRuleMap["redirectProtocol"]),
										"redirect_port":      convertInterfaceToInt(apiRuleMap["redirectPort"]),
										"redirect_code":      convertInterfaceToString(apiRuleMap["redirectCode"]),
										"target_ip_type":     convertInterfaceToString(apiRuleMap["targetIpType"]),
									}

									if apiRuleMap["healthCheck"] != nil {
										apiHealthCheckMap := apiRuleMap["healthCheck"].(map[string]interface{})
										ruleMap["health_check"] = []interface{}{
											map[string]interface{}{
												"port":                convertInterfaceToInt(apiHealthCheckMap["port"]),
												"protocol":            convertInterfaceToString(apiHealthCheckMap["protocol"]),
												"path":                convertInterfaceToString(apiHealthCheckMap["path"]),
												"interval":            convertInterfaceToInt(apiHealthCheckMap["interval"]),
												"timeout":             convertInterfaceToInt(apiHealthCheckMap["timeout"]),
												"unhealthy_threshold": convertInterfaceToInt(apiHealthCheckMap["unhealthyThreshold"]),
												"use_host_from_rule":  convertInterfaceToString(apiHealthCheckMap["useHostFromRule"]),
												"success_code":        convertInterfaceToString(apiHealthCheckMap["successCode"]),
												"host":                convertInterfaceToString(apiHealthCheckMap["host"]),
											},
										}
									}

									if apiRuleMap["providerEpgRef"] != nil {
										apiProviderEpgRefTokens := strings.Split(apiRuleMap["providerEpgRef"].(string), "/")
										ruleMap["provider_epg_ref"] = []interface{}{
											map[string]string{
												"schema_id":     apiProviderEpgRefTokens[2],
												"template_name": apiProviderEpgRefTokens[4],
												"anp_name":      apiProviderEpgRefTokens[6],
												"epg_name":      apiProviderEpgRefTokens[8],
											},
										}
									}
									rules = append(rules, ruleMap)
								}
								d.Set("rules", rules)
								d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s/serviceNodes/%d/listeners/%s", d.Get("schema_id"), siteID, templateName, contractName, serviceNodeIndex, listenerName))
								return nil
							}
						}
					}
				}
			}
		}
	}
	d.SetId("")
	return fmt.Errorf("Unable to find Site Contract Service Graph Listener: %s", listenerName)
}

func postSchemaSiteContractServiceGraphListenerConfig(ops string, d *schema.ResourceData, m interface{}) error {
	log.Printf("Inside postconfig")
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	siteID := d.Get("site_id").(string)
	serviceNodeIndex := d.Get("service_node_index")
	listenerName := d.Get("listener_name").(string)
	protocol := d.Get("protocol").(string)
	port := d.Get("port").(int)

	securityPolicy := d.Get("security_policy").(string)
	sslCertificates := d.Get("ssl_certificates").(*schema.Set).List()
	sslCertsPayloadMap := make([]interface{}, 0)
	for _, sslCert := range sslCertificates {
		sslCertMap := sslCert.(map[string]interface{})
		sslCertPayloadMap := map[string]interface{}{
			"name":    sslCertMap["name"].(string),
			"tDn":     sslCertMap["target_dn"].(string),
			"default": true,
			"store":   sslCertMap["certificate_store"].(string),
		}
		sslCertsPayloadMap = append(sslCertsPayloadMap, sslCertPayloadMap)
	}

	rulesList := d.Get("rules").(*schema.Set).List()
	rulesPayloadMap := make([]interface{}, 0)
	for _, rule := range rulesList {

		ruleMap := rule.(map[string]interface{})
		rulePayloadMap := map[string]interface{}{
			"name":             ruleMap["name"].(string),
			"floatingIp":       ruleMap["floating_ip"].(string),
			"index":            ruleMap["priority"].(int),
			"host":             ruleMap["host"].(string),
			"path":             ruleMap["path"].(string),
			"action":           ruleMap["action"].(string),
			"condition":        ruleMap["condition"].(string),
			"actionType":       ruleMap["action_type"].(string),
			"contentType":      ruleMap["content_type"].(string),
			"port":             ruleMap["port"].(int),
			"protocol":         ruleMap["protocol"].(string),
			"urlType":          ruleMap["url_type"].(string),
			"customURL":        ruleMap["custom_url"].(string),
			"redirectHostName": ruleMap["redirect_host_name"].(string),
			"redirectPath":     ruleMap["redirect_path"].(string),
			"redirectQuery":    ruleMap["redirect_query"].(string),
			"responseCode":     ruleMap["response_code"].(string),
			"responseBody":     ruleMap["response_body"].(string),
			"redirectProtocol": ruleMap["redirect_protocol"].(string),
			"redirectPort":     ruleMap["redirect_port"].(int),
			"redirectCode":     ruleMap["redirect_code"].(string),
			"targetIpType":     ruleMap["target_ip_type"].(string),
		}

		providerEpgRefDn := ""
		providerEpgRefList := ruleMap["provider_epg_ref"].(*schema.Set).List()
		if len(providerEpgRefList) >= 1 {
			providerEpgRef := providerEpgRefList[0].(interface{}).(map[string]interface{})

			epgRefSchemaID := providerEpgRef["schema_id"].(string)
			if epgRefSchemaID == "" {
				epgRefSchemaID = schemaID
			}
			epgRefTemplateName := providerEpgRef["template_name"].(string)
			if epgRefTemplateName == "" {
				epgRefTemplateName = templateName
			}
			providerEpgRefDn = fmt.Sprintf("/schemas/%s/templates/%s/anps/%s/epgs/%s", epgRefSchemaID, epgRefTemplateName, providerEpgRef["anp_name"].(string), providerEpgRef["epg_name"].(string))
		}
		rulePayloadMap["providerEpgRef"] = providerEpgRefDn

		healthCheckList := ruleMap["health_check"].(*schema.Set).List()
		healthCheckPayloadMap := make(map[string]interface{})

		if len(healthCheckList) >= 1 {
			healthCheckMap := healthCheckList[0].(interface{}).(map[string]interface{})
			healthCheckPayloadMap["port"] = healthCheckMap["port"].(int)
			healthCheckPayloadMap["protocol"] = healthCheckMap["protocol"].(string)
			healthCheckPayloadMap["path"] = healthCheckMap["path"].(string)
			healthCheckPayloadMap["interval"] = healthCheckMap["interval"].(int)
			healthCheckPayloadMap["timeout"] = healthCheckMap["timeout"].(int)
			healthCheckPayloadMap["unhealthyThreshold"] = healthCheckMap["unhealthy_threshold"].(int)
			healthCheckPayloadMap["useHostFromRule"] = healthCheckMap["use_host_from_rule"].(string)
			healthCheckPayloadMap["successCode"] = healthCheckMap["success_code"].(string)
			healthCheckPayloadMap["host"] = healthCheckMap["host"].(string)
		}
		rulePayloadMap["healthCheck"] = healthCheckPayloadMap
		rulesPayloadMap = append(rulesPayloadMap, rulePayloadMap)
	}

	frontendIpDnMap := make(map[string]string)
	if frontendIpDn, ok := d.GetOk("frontend_ip_dn"); ok {
		frontendIpDnMap["name"] = strings.Split(frontendIpDn.(string), "/vip-")[1]
		frontendIpDnMap["dn"] = frontendIpDn.(string)
	}

	pathListenerName := "-"
	if ops == "replace" {
		pathListenerName = listenerName
	}

	listenerPath := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship/serviceNodesRelationship/%d/deviceConfiguration/cloudLoadBalancer/listeners/%s",
		siteID, templateName, contractName, serviceNodeIndex, pathListenerName,
	)

	// Removing empty string and zeros to build the valid payload
	removeEmptyValuesFromListenerData(sslCertsPayloadMap)
	removeEmptyValuesFromListenerData(rulesPayloadMap)

	listenerPayload := models.NewSiteContractServiceGraphListener(ops, listenerPath, listenerName, protocol, securityPolicy, port, sslCertsPayloadMap, rulesPayloadMap, frontendIpDnMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), listenerPayload)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s/serviceNodes/%d/listeners/%s", schemaID, siteID, templateName, contractName, serviceNodeIndex, listenerName))
	return nil
}

func removeEmptyValuesFromListenerData(value interface{}) bool {
	switch v := value.(type) {
	case int:
		if v == 0 {
			return true
		}
	case string:
		if v == "" {
			return true
		}
	case float64:
		if v == 0.0 {
			return true
		}
	case nil:
		if v == nil {
			return true
		}
	case map[string]string:
		if len(v) == 0 {
			return true
		}
	case map[string]interface{}:
		for key, subValue := range value.(map[string]interface{}) {
			checkEmptyValue := removeEmptyValuesFromListenerData(subValue)
			if checkEmptyValue && key != "index" {
				delete(value.(map[string]interface{}), key)
			}
		}
	case []map[string]interface{}:
		if len(v) == 0 {
			return true
		} else {
			for _, mapObject := range value.([]map[string]interface{}) {
				for key, subValue := range mapObject {
					checkEmptyValue := removeEmptyValuesFromListenerData(subValue)
					if checkEmptyValue && key != "index" {
						delete(mapObject, key)
					}
				}
			}
		}
	case []interface{}:
		if len(v) == 0 {
			return true
		} else {
			for _, mapObject := range value.([]interface{}) {
				removeEmptyValuesFromListenerData(mapObject)
			}
		}
	case []map[string]string:
		if len(v) == 0 {
			return true
		} else {
			for _, mapObject := range value.([]map[string]string) {
				for key, subValue := range mapObject {
					checkEmptyValue := removeEmptyValuesFromListenerData(subValue)
					if checkEmptyValue && key != "index" {
						delete(mapObject, key)
					}
				}
			}
		}
	default:
		log.Printf("[DEBUG]: Unknown type, Object Value: %v", v)
	}
	return false
}

func convertInterfaceToString(interfaceObject interface{}) string {
	switch v := interfaceObject.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return ""
	}
}

func convertInterfaceToInt(interfaceObject interface{}) int {
	switch v := interfaceObject.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case nil:
		return 0
	default:
		return 0
	}
}
