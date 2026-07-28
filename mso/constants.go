package mso

import "fmt"

func cloudDeprecationMessage(resourceName string) string {
	return fmt.Sprintf("%s is deprecated: cloud-specific features are no longer supported in Nexus Dashboard 4.x (NDO 5.x) releases. Ensure the related configuration is removed from the NDO schema before upgrading Nexus Dashboard to 4.x (NDO 5.x).", resourceName)
}

func nd32DeprecationMessage(resourceName string) string {
	return fmt.Sprintf("%s is deprecated: no longer functional on Nexus Dashboard (ND) 3.2+ / NDO 4.4+ and will be removed in the next major release.", resourceName)
}

func nd4DeprecationMessage(resourceName string) string {
	return fmt.Sprintf("%s is deprecated: no longer functional on Nexus Dashboard (ND) 4.0+ / NDO 5.0+ and will be removed once ND 3.x / NDO 4.x is no longer supported.", resourceName)
}

func nd4FunctionalDataSourceDeprecationMessage(dataSourceName string) string {
	return fmt.Sprintf("%s data source is deprecated: it remains functional on Nexus Dashboard (ND) 4.0+ / NDO 5.0+, but the corresponding %s resource is no longer functional. The data source will be removed once ND 3.x / NDO 4.x is no longer supported.", dataSourceName, dataSourceName)
}

var controlMap = map[string]string{
	"fast_sel_hot_stdby": "fast-sel-hot-stdby",
	"graceful_conv":      "graceful-conv",
	"susp_individual":    "susp-individual",
	"load_defer":         "load-defer",
	"symmetric_hash":     "symmetric-hash",
}

var enabledDisabledMap = map[interface{}]interface{}{
	"enabled":  true,
	"disabled": false,
	true:       "enabled",
	false:      "disabled",
}

var l2InterfaceQinqMap = map[string]string{
	"double_q_tag_port": "doubleQtagPort",
	"core_port":         "corePort",
	"edge_port":         "edgePort",
	"disabled":          "disabled",
}

var linkLevelFecMap = map[string]string{
	"inherit":       "inherit",
	"cl74_fc_fec":   "cl74-fc-fec",
	"cl91_rs_fec":   "cl91-rs-fec",
	"cons16_rs_fec": "cons16-rs-fec",
	"ieee_rs_fec":   "ieee-rs-fec",
	"kp_fec":        "kp-fec",
	"disable_fec":   "disable-fec",
}

var loadBalanceHashingMap = map[string]string{
	"destination_ip":         "dst-ip",
	"layer_4_destination_ip": "l4-dst-port",
	"layer_4_source_ip":      "l4-src-port",
	"source_ip":              "src-ip",
}

var matchParameterMap = map[string]string{
	"destination_ip":   "dstIP",
	"destination_ipv4": "dstIPv4",
	"destination_ipv6": "dstIPv6",
	"destination_mac":  "dstMac",
	"destination_port": "dstPort",
	"ethertype":        "ethertype",
	"ip_protocol":      "proto",
	"source_ip":        "srcIP",
	"source_ipv4":      "srcIPv4",
	"source_ipv6":      "srcIPv6",
	"source_mac":       "srcMac",
	"source_port":      "srcPort",
	"dstIP":            "destination_ip",
	"dstIPv4":          "destination_ipv4",
	"dstIPv6":          "destination_ipv6",
	"dstMac":           "destination_mac",
	"dstPort":          "destination_port",
	"proto":            "ip_protocol",
	"srcIP":            "source_ip",
	"srcIPv4":          "source_ipv4",
	"srcIPv6":          "source_ipv6",
	"srcMac":           "source_mac",
	"srcPort":          "source_port",
}

var portChannelModeMap = map[string]string{
	"lacp_active":                   "active",
	"lacp_passive":                  "passive",
	"static_channel_mode_on":        "off",
	"mac_pinning":                   "mac-pin",
	"mac_pinning_physical_nic_load": "mac-pin-nicload",
	"use_explicit_failover_order":   "explicit-failover",
}

var ptpDestinationMacMap = map[string]string{
	"forwardable":     "forwardable",
	"non_forwardable": "nonForwardable",
	"nonForwardable":  "non_forwardable",
}

var ptpMismatchedMacHandlingMap = map[string]string{
	"drop":                    "drop",
	"reply_with_config_mac":   "replyWithCfgMac",
	"replyWithCfgMac":         "reply_with_config_mac",
	"reply_with_received_mac": "replyWithRxMac",
	"replyWithRxMac":          "reply_with_received_mac",
}

var ptpProfileTemplateMap = map[string]string{
	"aes67":           "aes67",
	"default":         "default",
	"smpte":           "smpte",
	"telecom":         "telecomFullPath",
	"telecomFullPath": "telecom",
}

var synceQualityLevelOptionsMap = map[string]string{
	"op1":                   "option_1",
	"op2g1":                 "option_2_generation_1",
	"op2g2":                 "option_2_generation_2",
	"option_1":              "op1",
	"option_2_generation_1": "op2g1",
	"option_2_generation_2": "op2g2",
}

var targetCosMap = map[string]string{
	"background":            "cos0",
	"best_effort":           "cos1",
	"excellent_effort":      "cos2",
	"critical_applications": "cos3",
	"video":                 "cos4",
	"voice":                 "cos5",
	"internetwork_control":  "cos6",
	"network_control":       "cos7",
	"unspecified":           "unspecified",
	"cos0":                  "background",
	"cos1":                  "best_effort",
	"cos2":                  "excellent_effort",
	"cos3":                  "critical_applications",
	"cos4":                  "video",
	"cos5":                  "voice",
	"cos6":                  "internetwork_control",
	"cos7":                  "network_control",
}

var targetDscpMap = map[string]string{
	"af11":                 "af11",
	"af12":                 "af12",
	"af13":                 "af13",
	"af21":                 "af21",
	"af22":                 "af22",
	"af23":                 "af23",
	"af31":                 "af31",
	"af32":                 "af32",
	"af33":                 "af33",
	"af41":                 "af41",
	"af42":                 "af42",
	"af43":                 "af43",
	"cs0":                  "cs0",
	"cs1":                  "cs1",
	"cs2":                  "cs2",
	"cs3":                  "cs3",
	"cs4":                  "cs4",
	"cs5":                  "cs5",
	"cs6":                  "cs6",
	"cs7":                  "cs7",
	"expedited_forwarding": "expeditedForwarding",
	"voice_admit":          "voiceAdmit",
	"unspecified":          "unspecified",
	"expeditedForwarding":  "expedited_forwarding",
	"voiceAdmit":           "voice_admit",
}
