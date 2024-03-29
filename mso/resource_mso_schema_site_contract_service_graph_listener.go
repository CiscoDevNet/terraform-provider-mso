package mso

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var listenerProtocols = []string{"http", "https", "tcp", "udp", "tls", "inherit"}

var listenerSecurityPolicyMap = map[string]string{
	"default":                     "default",
	"elb_sec_2016_18":             "eLBSecurityPolicy-2016-08",
	"elb_sec_fs_2018_06":          "eLBSecurityPolicy-FS-2018-06",
	"elb_sec_tls_1_2_2017_01":     "eLBSecurityPolicy-TLS-1-2-2017-01",
	"elb_sec_tls_1_2_ext_2018_06": "eLBSecurityPolicy-TLS-1-2-Ext-2018-06",
	"elb_sec_tls_1_1_2017_01":     "eLBSecurityPolicy-TLS-1-1-2017-01",
	"elb_sec_2015_05":             "eLBSecurityPolicy-2015-05",
	"elb_sec_tls_1_0_2015_04":     "eLBSecurityPolicy-TLS-1-0-2015-04",
	"app_gw_ssl_default":          "AppGwSslPolicyDefault",
	"app_gw_ssl_2015_501":         "AppGwSslPolicy20150501",
	"app_gw_ssl_2017_401":         "AppGwSslPolicy20170401",
	"app_gw_ssl_2017_401s":        "AppGwSslPolicy20170401S",
}

var listenerActionTypeMap = map[string]string{
	"fixed_response": "fixedResponse",
	"forward":        "forward",
	"redirect":       "redirect",
	"ha_port":        "haPort",
}

var listenerContentTypeMap = map[string]string{
	"text_plain": "textPlain",
	"text_css":   "textCSS",
	"text_html":  "textHtml",
	"app_js":     "appJS",
	"app_json":   "appJson",
}

var listenerRedirectCodeMap = map[string]string{
	"unknown":            "unknown",
	"permanently_moved":  "permMoved",
	"found":              "found",
	"see_other":          "seeOther",
	"temporary_redirect": "temporary",
}

var listenerSecurityPolicyKeys = getMapKeys(listenerSecurityPolicyMap)

var listenerActionTypeKeys = getMapKeys(listenerActionTypeMap)

var listenerContentTypeKeys = getMapKeys(listenerContentTypeMap)

var listenerRedirectCodeKeys = getMapKeys(listenerRedirectCodeMap)

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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(listenerProtocols, false),
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(listenerSecurityPolicyKeys, false),
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// NOTE: rules is a required attribute but to utilize the ResourceDiff functionality marked as optional and computed
			"rules": &schema.Schema{
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
						"floating_ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"host": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"action": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// TODO: Should be uncommented once condition is configured through UI
						// "condition": &schema.Schema{
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// 	Computed: true,
						// },
						"action_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(listenerActionTypeKeys, false),
						},
						"content_type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(listenerContentTypeKeys, false),
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(listenerProtocols, false),
						},
						"provider_epg_ref": { // Only one object allowed
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schema_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"template_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"redirect_host_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"redirect_path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"redirect_query": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"response_code": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"response_body": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"redirect_protocol": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(listenerProtocols, false),
						},
						"redirect_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"redirect_code": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(listenerRedirectCodeKeys, false),
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
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(listenerProtocols, false),
									},
									"path": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
									"use_host_from_rule": &schema.Schema{ // By default the TypeBool returns false
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"success_code": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"host": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
		// Clear the "rules" attribute diff when its not a valid change
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// When the listener Protocol is https
			_, listenerProtocol := diff.GetChange("protocol")
			if listenerProtocol.(string) == "https" {
				_, newSecurityPolicy := diff.GetChange("security_policy")
				_, newSSLCertificates := diff.GetChange("ssl_certificates")
				if newSecurityPolicy.(string) == "" || len(newSSLCertificates.(*schema.Set).List()) == 0 {
					return fmt.Errorf("When the listener Protocol is https, the security_policy and ssl_certificates attributes must be set")
				}
			}

			// Get the diff between old and new rules
			oldRules, newRules := diff.GetChange("rules")

			// oldRules - holds the API response content
			oldRulesList := oldRules.(*schema.Set).List()

			// newRules - holds the user config content
			newRulesList := newRules.(*schema.Set).List()

			if len(newRulesList) == 0 {
				return fmt.Errorf("Rules cannot be empty, minimum one item is required to perform 'create/update' operation")
			}

			for _, newRule := range newRulesList {
				newRuleMap := newRule.(map[string]interface{})

				// When the rule action_type is "redirect"
				if newRuleMap["action_type"].(string) == "redirect" {
					newRedirectProtocol := newRuleMap["redirect_protocol"].(string)
					newRedirectPort := newRuleMap["redirect_port"].(int)
					newUrlType := newRuleMap["url_type"].(string)
					newResponseCode := newRuleMap["redirect_code"].(string)
					if newRedirectProtocol == "" || newRedirectPort == 0 || newUrlType == "" || newResponseCode == "" {
						return fmt.Errorf("When the 'action_type' is 'redirect', the 'redirect_protocol', 'redirect_port', 'url_type', 'redirect_code' attributes must be set")
					}
				}

				// When the rule action_type is "forward"
				if newRuleMap["action_type"].(string) == "forward" {
					newRuleProtocol := newRuleMap["protocol"].(string)
					newRulePort := newRuleMap["port"].(int)
					newHealthCheck := newRuleMap["health_check"].(*schema.Set).List()
					if newRuleProtocol == "" || newRulePort == 0 || len(newHealthCheck) == 0 {
						return fmt.Errorf("When the 'action_type' is 'forward', the 'protocol', 'port', 'health_check' attributes must be set")
					}
				}

				// When the rule url_type is "custom"
				if newRuleMap["url_type"].(string) == "custom" {
					newRuleRedirectHostName := newRuleMap["redirect_host_name"].(string)
					newRuleRedirectPath := newRuleMap["redirect_path"].(string)
					newRuleRedirectQuery := newRuleMap["redirect_query"].(string)
					newRuleResponseCode := newRuleMap["response_code"].(string)
					if newRuleRedirectHostName == "" || newRuleRedirectPath == "" || newRuleRedirectQuery == "" || newRuleResponseCode == "" {
						return fmt.Errorf("When the 'url_type' is 'custom', the 'redirect_host_name', 'redirect_path', 'redirect_query', 'response_code' attributes must be set")
					}
				}

				// When the Health Checks protocol is "http/https"
				newHealthCheckList := newRuleMap["health_check"].(*schema.Set).List()
				if len(newHealthCheckList) > 0 {
					newHealthCheckMap := newHealthCheckList[0].(interface{}).(map[string]interface{})
					newHealthCheckProtocol := newHealthCheckMap["protocol"].(string)
					if newHealthCheckProtocol == "http" || newHealthCheckProtocol == "https" {
						newHealthCheckPort := newHealthCheckMap["port"].(int)
						newHealthCheckPath := newHealthCheckMap["path"].(string)
						newHealthCheckUnhealthyThreshold := newHealthCheckMap["unhealthy_threshold"].(int)
						newHealthCheckTimeout := newHealthCheckMap["timeout"].(int)
						newHealthCheckInterval := newHealthCheckMap["interval"].(int)
						newHealthCheckSuccessCode := newHealthCheckMap["success_code"].(string)
						if newHealthCheckPort == 0 || newHealthCheckPath == "" || newHealthCheckUnhealthyThreshold == 0 || newHealthCheckTimeout == 0 || newHealthCheckInterval == 0 || newHealthCheckSuccessCode == "" {
							return fmt.Errorf("When the 'health_check' protocol is 'http/https', the 'port', 'path', 'unhealthy_threshold', 'timeout', 'interval', 'success_code' attributes must be set")
						}
						// when the listener protocol is "http/https", the Health Checks protocol is "http/https" and use_host_from_rule is true then the host attribute must be set
						if listenerProtocol.(string) == "http" || listenerProtocol.(string) == "https" {
							newHealthCheckUseHostFromRule := newHealthCheckMap["use_host_from_rule"].(bool)
							newHealthCheckHost := newHealthCheckMap["host"].(string)
							if !newHealthCheckUseHostFromRule && newHealthCheckHost == "" {
								return fmt.Errorf("When the 'health_check' protocol is 'http/https', the 'use_host_from_rule' and 'host' attributes must be set")
							} else if newHealthCheckUseHostFromRule && newHealthCheckHost != "" {
								return fmt.Errorf("When the 'use_host_from_rule' is true, the 'host' should be empty")
							}
						}
					}

					// when the listener protocol is "tcp/udp" Provider EPG is required
					if listenerProtocol.(string) == "tcp" || listenerProtocol.(string) == "udp" {
						newProviderEPGRefList := newRuleMap["provider_epg_ref"].(*schema.Set).List()
						if len(newProviderEPGRefList) == 0 || len(newRuleMap["health_check"].(*schema.Set).List()) == 0 {
							return fmt.Errorf("When the 'listener_protocol' is 'tcp/udp', the 'provider_epg_ref', 'health_check' attributes must be set")
						}
						// When the listener protocol is "tcp/udp" and the Health Checks protocol is "tcp"
						if newHealthCheckProtocol == "tcp" {
							newHealthCheckPort := newHealthCheckMap["port"].(int)
							newHealthCheckUnhealthyThreshold := newHealthCheckMap["unhealthy_threshold"].(int)
							newHealthCheckInterval := newHealthCheckMap["interval"].(int)
							if newHealthCheckPort == 0 || newHealthCheckUnhealthyThreshold == 0 || newHealthCheckInterval == 0 {
								return fmt.Errorf("When the 'health_check' protocol is 'tcp', the 'port', 'unhealthy_threshold', 'interval' attributes must be set")
							}
						} else if newHealthCheckProtocol == "http" || newHealthCheckProtocol == "https" {
							// When the listener protocol is "tcp/udp" and the Health Checks protocol is "http/https"
							newHealthCheckPort := newHealthCheckMap["port"].(int)
							newHealthCheckPath := newHealthCheckMap["path"].(string)
							newHealthCheckUnhealthyThreshold := newHealthCheckMap["unhealthy_threshold"].(int)
							newHealthCheckInterval := newHealthCheckMap["interval"].(int)
							if newHealthCheckPort == 0 || newHealthCheckPath == "" || newHealthCheckUnhealthyThreshold == 0 || newHealthCheckInterval == 0 {
								return fmt.Errorf("When the 'health_check' protocol is 'tcp', the 'port', 'newHealthCheckPath', 'unhealthy_threshold', 'interval' attributes must be set")
							}
						}
					}
				}
			}

			sort.Slice(oldRulesList, func(i, j int) bool {
				return oldRulesList[i].(interface{}).(map[string]interface{})["priority"].(int) < oldRulesList[j].(interface{}).(map[string]interface{})["priority"].(int)
			})

			sort.Slice(newRulesList, func(i, j int) bool {
				return newRulesList[i].(interface{}).(map[string]interface{})["priority"].(int) < newRulesList[j].(interface{}).(map[string]interface{})["priority"].(int)
			})

			if len(oldRulesList) != len(newRulesList) { // Valid change, no need to clear rules attribute diff
				return nil
			}

			rulesDiffBool := true // true - is not a valid change and false - is a valid change
			for _, newRule := range newRulesList {
				for _, oldRule := range oldRulesList {

					newRuleMap := newRule.(map[string]interface{})
					newRuleName := newRuleMap["name"].(string)
					newRulePriority := newRuleMap["priority"].(int)
					oldRuleMap := oldRule.(map[string]interface{})
					oldRuleName := oldRuleMap["name"].(string)
					oldRulePriority := oldRuleMap["priority"].(int)

					if newRuleName == oldRuleName && newRulePriority == oldRulePriority {

						newRuleHealthCheckList := newRuleMap["health_check"].(*schema.Set).List()
						oldRuleHealthCheckList := oldRuleMap["health_check"].(*schema.Set).List()

						// When the newRuleHealthCheckList (user config) and oldRuleHealthCheckList (state file content) not empty
						if len(newRuleHealthCheckList) != 0 && len(oldRuleHealthCheckList) != 0 {

							newRuleHealthCheckMap := newRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})
							oldRuleHealthCheckMap := oldRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})

							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["host"].(string) == oldRuleHealthCheckMap["host"].(string)
							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["interval"].(int) == oldRuleHealthCheckMap["interval"].(int)
							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["path"].(string) == oldRuleHealthCheckMap["path"].(string)
							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["port"].(int) == oldRuleHealthCheckMap["port"].(int)
							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["protocol"].(string) == oldRuleHealthCheckMap["protocol"].(string)

							newRuleHealthCheckSuccessCode := newRuleHealthCheckMap["success_code"].(string)
							oldRuleHealthCheckSuccessCode := oldRuleHealthCheckMap["success_code"].(string)
							if newRuleHealthCheckSuccessCode == "" && oldRuleHealthCheckSuccessCode == "200-399" {
								rulesDiffBool = rulesDiffBool && true // Not a valid change because 200-399 is the default value
							} else if newRuleHealthCheckSuccessCode == oldRuleHealthCheckSuccessCode {
								rulesDiffBool = rulesDiffBool && true // Not a valid change - both are equal
							} else if newRuleHealthCheckSuccessCode == "" && oldRuleHealthCheckSuccessCode != "" {
								rulesDiffBool = rulesDiffBool && true // Not a valid change - because the success_code is optional and computed, so the when the success_code is empty we will not consider it as a valid change
							} else if (newRuleHealthCheckSuccessCode != "" && oldRuleHealthCheckSuccessCode != "") && (newRuleHealthCheckSuccessCode != oldRuleHealthCheckSuccessCode) {
								rulesDiffBool = rulesDiffBool && false // Valid change - both are not equal and also not empty
							} else {
								rulesDiffBool = rulesDiffBool && false // By default set false when the above conditions are not met
							}

							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["timeout"].(int) == oldRuleHealthCheckMap["timeout"].(int)
							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["unhealthy_threshold"].(int) == oldRuleHealthCheckMap["unhealthy_threshold"].(int)

							rulesDiffBool = rulesDiffBool && newRuleHealthCheckMap["use_host_from_rule"] == oldRuleHealthCheckMap["use_host_from_rule"]
						} else if len(newRuleHealthCheckList) == 0 && len(oldRuleHealthCheckList) != 0 {
							// When the newRuleHealthCheckList (user config) is empty and oldRuleHealthCheckList (state file content) not empty
							oldRuleHealthCheckMap := oldRuleMap["health_check"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})

							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["host"].(string) == ""
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["interval"].(int) == 0
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["path"].(string) == ""
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["port"].(int) == 0
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["protocol"].(string) == ""
							rulesDiffBool = rulesDiffBool && (oldRuleHealthCheckMap["success_code"].(string) == "200-399" || oldRuleHealthCheckMap["success_code"].(string) == "")
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["timeout"].(int) == 0
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["unhealthy_threshold"].(int) == 0
							rulesDiffBool = rulesDiffBool && oldRuleHealthCheckMap["use_host_from_rule"] == false
						}

						newProviderEPGRefList := newRuleMap["provider_epg_ref"].(*schema.Set).List()
						oldProviderEPGRefList := oldRuleMap["provider_epg_ref"].(*schema.Set).List()

						if len(newProviderEPGRefList) == 1 && len(oldProviderEPGRefList) == 1 {
							newProviderEPGRefMap := newRuleMap["provider_epg_ref"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})
							oldProviderEPGRefMap := oldRuleMap["provider_epg_ref"].(*schema.Set).List()[0].(interface{}).(map[string]interface{})

							newSchemaID := convertInterfaceToString(newProviderEPGRefMap["schema_id"].(string))
							newTemplateName := convertInterfaceToString(newProviderEPGRefMap["template_name"].(string))

							if newSchemaID == "" {
								newSchemaID = diff.Get("schema_id").(string)
							}
							if newTemplateName == "" {
								newTemplateName = diff.Get("template_name").(string)
							}

							rulesDiffBool = rulesDiffBool && newProviderEPGRefMap["anp_name"].(string) == oldProviderEPGRefMap["anp_name"].(string)
							rulesDiffBool = rulesDiffBool && newProviderEPGRefMap["epg_name"].(string) == oldProviderEPGRefMap["epg_name"].(string)
							rulesDiffBool = rulesDiffBool && newSchemaID == oldProviderEPGRefMap["schema_id"].(string)
							rulesDiffBool = rulesDiffBool && newTemplateName == oldProviderEPGRefMap["template_name"].(string)

						} else if (len(newProviderEPGRefList) == 0 && len(oldProviderEPGRefList) == 1) || (len(newProviderEPGRefList) == 1 && len(oldProviderEPGRefList) == 0) {
							rulesDiffBool = rulesDiffBool && false // its a valid change
						}

						newFloatingIp := newRuleMap["floating_ip"].(string)
						newHost := newRuleMap["host"].(string)
						newPath := newRuleMap["path"].(string)
						newAction := newRuleMap["action"].(string)
						// TODO: Should be uncommented once condition is configured through UI
						// newCondition := newRuleMap["condition"].(string)
						newActionType := newRuleMap["action_type"].(string)
						newContentType := newRuleMap["content_type"].(string)
						newProtocol := newRuleMap["protocol"].(string)
						newUrlType := newRuleMap["url_type"].(string)
						newCustomUrl := newRuleMap["custom_url"].(string)
						newRedirectHostName := newRuleMap["redirect_host_name"].(string)
						newRedirectPath := newRuleMap["redirect_path"].(string)
						newRedirectQuery := newRuleMap["redirect_query"].(string)
						newResponseCode := newRuleMap["response_code"].(string)
						newResponseBody := newRuleMap["response_body"].(string)
						newRedirectProtocol := newRuleMap["redirect_protocol"].(string)
						newRedirectCode := newRuleMap["redirect_code"].(string)
						newTargetIpType := newRuleMap["target_ip_type"].(string)

						newPort := newRuleMap["port"].(int)
						newRedirectPort := newRuleMap["redirect_port"].(int)

						oldFloatingIp := oldRuleMap["floating_ip"].(string)
						oldHost := oldRuleMap["host"].(string)
						oldPath := oldRuleMap["path"].(string)
						oldAction := oldRuleMap["action"].(string)
						// TODO: Should be uncommented once condition is configured through UI
						// oldCondition := oldRuleMap["condition"].(string)
						oldActionType := oldRuleMap["action_type"].(string)
						oldContentType := oldRuleMap["content_type"].(string)
						oldProtocol := oldRuleMap["protocol"].(string)
						oldUrlType := oldRuleMap["url_type"].(string)
						oldCustomUrl := oldRuleMap["custom_url"].(string)
						oldRedirectHostName := oldRuleMap["redirect_host_name"].(string)
						oldRedirectPath := oldRuleMap["redirect_path"].(string)
						oldRedirectQuery := oldRuleMap["redirect_query"].(string)
						oldResponseCode := oldRuleMap["response_code"].(string)
						oldResponseBody := oldRuleMap["response_body"].(string)
						oldRedirectProtocol := oldRuleMap["redirect_protocol"].(string)
						oldRedirectCode := oldRuleMap["redirect_code"].(string)
						oldTargetIpType := oldRuleMap["target_ip_type"].(string)

						oldPort := oldRuleMap["port"].(int)
						oldRedirectPort := oldRuleMap["redirect_port"].(int)

						rulesDiffBool = rulesDiffBool && !(newPort != 0 && newPort != oldPort)
						rulesDiffBool = rulesDiffBool && !(newRedirectPort != 0 && newRedirectPort != oldRedirectPort)
						rulesDiffBool = rulesDiffBool && !(newFloatingIp != "" && newFloatingIp != oldFloatingIp)
						rulesDiffBool = rulesDiffBool && !(newHost != "" && newHost != oldHost)
						rulesDiffBool = rulesDiffBool && !(newPath != "" && newPath != oldPath)
						rulesDiffBool = rulesDiffBool && !(newAction != "" && newAction != oldAction)
						// TODO: Should be uncommented once condition is configured through UI
						// rulesDiffBool = rulesDiffBool &&  !(newCondition != "" && newCondition != oldCondition)
						rulesDiffBool = rulesDiffBool && !(newActionType != "" && newActionType != oldActionType)
						rulesDiffBool = rulesDiffBool && !(newContentType != "" && newContentType != oldContentType)
						rulesDiffBool = rulesDiffBool && !(newProtocol != "" && newProtocol != oldProtocol)
						rulesDiffBool = rulesDiffBool && !(newUrlType != "" && newUrlType != oldUrlType)
						rulesDiffBool = rulesDiffBool && !(newCustomUrl != "" && newCustomUrl != oldCustomUrl)
						rulesDiffBool = rulesDiffBool && !(newRedirectHostName != "" && newRedirectHostName != oldRedirectHostName)
						rulesDiffBool = rulesDiffBool && !(newRedirectPath != "" && newRedirectPath != oldRedirectPath)
						rulesDiffBool = rulesDiffBool && !(newRedirectQuery != "" && newRedirectQuery != oldRedirectQuery)
						rulesDiffBool = rulesDiffBool && !(newResponseCode != "" && newResponseCode != oldResponseCode)
						rulesDiffBool = rulesDiffBool && !(newResponseBody != "" && newResponseBody != oldResponseBody)
						rulesDiffBool = rulesDiffBool && !(newRedirectProtocol != "" && newRedirectProtocol != oldRedirectProtocol)
						rulesDiffBool = rulesDiffBool && !(newRedirectCode != "" && newRedirectCode != oldRedirectCode)
						rulesDiffBool = rulesDiffBool && !(newTargetIpType != "" && newTargetIpType != oldTargetIpType)

						break
					} else if (newRuleName != oldRuleName && newRulePriority == oldRulePriority) || (newRuleName == oldRuleName && newRulePriority != oldRulePriority) {
						rulesDiffBool = rulesDiffBool && false // its a valid change
					}
				}
				// When the rulesDiffBool is false (valid change), stop the diff checking
				if !rulesDiffBool {
					return nil
				}
			}

			// rulesDiffBool = true - is not a valid change and also the initial value, so clear the diff
			// rulesDiffBool = false - is a valid change
			if rulesDiffBool {
				diff.Clear("rules")
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
								d.Set("security_policy", getKeyByValue(listenerSecurityPolicyMap, convertInterfaceToString(listenerMap["secPolicy"])))

								if listenerMap["certificates"] != nil {
									sslCertificates := make([]map[string]interface{}, 0)
									for _, certificate := range listenerMap["certificates"].([]interface{}) {
										certificateMap := certificate.(map[string]interface{})
										sslCertMap := map[string]interface{}{
											"name":              convertInterfaceToString(certificateMap["name"]),
											"target_dn":         convertInterfaceToString(certificateMap["tDn"]),
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
										"name":        convertInterfaceToString(apiRuleMap["name"]),
										"floating_ip": convertInterfaceToString(apiRuleMap["floatingIp"]),
										"priority":    convertInterfaceToInt(apiRuleMap["index"]),
										"host":        convertInterfaceToString(apiRuleMap["host"]),
										"path":        convertInterfaceToString(apiRuleMap["path"]),
										"action":      convertInterfaceToString(apiRuleMap["action"]),
										// TODO: Should be uncommented once condition is configured through UI
										// "condition":          convertInterfaceToString(apiRuleMap["condition"]),
										"action_type":        getKeyByValue(listenerActionTypeMap, convertInterfaceToString(apiRuleMap["actionType"])),
										"content_type":       getKeyByValue(listenerContentTypeMap, convertInterfaceToString(apiRuleMap["contentType"])),
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
										"redirect_code":      getKeyByValue(listenerRedirectCodeMap, convertInterfaceToString(apiRuleMap["redirectCode"])),
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
												"use_host_from_rule":  yesNoToBool(convertInterfaceToString(apiHealthCheckMap["useHostFromRule"])),
												"success_code":        convertInterfaceToString(apiHealthCheckMap["successCode"]),
												"host":                convertInterfaceToString(apiHealthCheckMap["host"]),
											},
										}
									}

									if apiRuleMap["providerEpgRef"] != nil && apiRuleMap["providerEpgRef"] != "" {
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
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)
	siteID := d.Get("site_id").(string)
	serviceNodeIndex := d.Get("service_node_index")
	listenerName := d.Get("listener_name").(string)
	protocol := d.Get("protocol").(string)
	port := d.Get("port").(int)

	securityPolicy := listenerSecurityPolicyMap[d.Get("security_policy").(string)]
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
			"name":       ruleMap["name"].(string),
			"floatingIp": ruleMap["floating_ip"].(string),
			"index":      ruleMap["priority"].(int),
			"host":       ruleMap["host"].(string),
			"path":       ruleMap["path"].(string),
			"action":     ruleMap["action"].(string),
			// TODO: Should be uncommented once condition is configured through UI
			// "condition":        ruleMap["condition"].(string),
			"actionType":       listenerActionTypeMap[ruleMap["action_type"].(string)],
			"contentType":      listenerContentTypeMap[ruleMap["content_type"].(string)],
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
			"redirectCode":     listenerRedirectCodeMap[ruleMap["redirect_code"].(string)],
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
			healthCheckPayloadMap["successCode"] = healthCheckMap["success_code"].(string)

			healthCheckPayloadMap["useHostFromRule"] = boolToYesNo(healthCheckMap["use_host_from_rule"].(bool))
			if healthCheckPayloadMap["useHostFromRule"] == "no" {
				healthCheckPayloadMap["host"] = healthCheckMap["host"].(string)
			}
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

// removeEmptyValuesFromListenerData removes empty values from the given listener data.
//
// value: the value to check for empty values.
// Returns: true if the value is empty, false otherwise.
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
