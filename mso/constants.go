package mso

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

var enabledDisabledMap = map[interface{}]interface{}{
	"enabled":  true,
	"disabled": false,
	true:       "enabled",
	false:      "disabled",
}

var portChannelModeMap = map[string]string{
	"lacp_active":                   "active",
	"lacp_passive":                  "passive",
	"static_channel_mode_on":        "off",
	"mac_pinning":                   "mac-pin",
	"mac_pinning_physical_nic_load": "mac-pin-nicload",
	"use_explicit_failover_order":   "explicit-failover",
}

var controlMap = map[string]string{
	"fast_sel_hot_stdby": "fast-sel-hot-stdby",
	"graceful_conv":      "graceful-conv",
	"susp_individual":    "susp-individual",
	"load_defer":         "load-defer",
	"symmetric_hash":     "symmetric-hash",
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

var l2InterfaceQinqMap = map[string]string{
	"double_q_tag_port": "doubleQtagPort",
	"core_port":         "corePort",
	"edge_port":         "edgePort",
	"disabled":          "disabled",
}

var loadBalanceHashingMap = map[string]string{
	"destination_ip":         "dst-ip",
	"layer_4_destination_ip": "l4-dst-port",
	"layer_4_source_ip":      "l4-src-port",
	"source_ip":              "src-ip",
}
