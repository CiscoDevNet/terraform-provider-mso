package models

type TemplateBD struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateBD(ops, path, name, displayName, layer2Unicast, unkMcastAct, multiDstPktAct, v6unkMcastAct, vmac, description string, intersiteBumTrafficAllow, optimizeWanBandwidth, l2Stretch, l3MCast, arpFlood, unicastRouting bool, vrfRef, dhcpLabel map[string]interface{}, dhcpLabels []interface{}) *PatchPayload {
	var bdMap map[string]interface{}
	bdMap = map[string]interface{}{
		"name":                     name,
		"displayName":              displayName,
		"l2UnknownUnicast":         layer2Unicast,
		"unkMcastAct":              unkMcastAct,
		"multiDstPktAct":           multiDstPktAct,
		"v6unkMcastAct":            v6unkMcastAct,
		"vmac":                     vmac,
		"arpFlood":                 arpFlood,
		"unicastRouting":           unicastRouting,
		"intersiteBumTrafficAllow": intersiteBumTrafficAllow,
		"optimizeWanBandwidth":     optimizeWanBandwidth,
		"l2Stretch":                l2Stretch,
		"l3MCast":                  l3MCast,
		"vrfRef":                   vrfRef,
		"dhcpLabel":                dhcpLabel,
		"dhcpLabels":               dhcpLabels,
		"subnets":                  []interface{}{},
		"description":              description,
	}

	if bdMap["l2UnknownUnicast"] == "" {
		bdMap["l2UnknownUnicast"] = "flood"
	}

	if bdMap["unkMcastAct"] == "optimized_flooding" {
		bdMap["unkMcastAct"] = "opt-flood"
	} else {
		bdMap["unkMcastAct"] = "flood"
	}

	if bdMap["multiDstPktAct"] == "flood_in_bd" || bdMap["multiDstPktAct"] == "" {
		bdMap["multiDstPktAct"] = "bd-flood"
	}

	if bdMap["multiDstPktAct"] == "flood_in_encap" {
		bdMap["multiDstPktAct"] = "encap-flood"
	}

	if bdMap["v6unkMcastAct"] == "optimized_flooding" {
		bdMap["v6unkMcastAct"] = "opt-flood"
	} else {
		bdMap["v6unkMcastAct"] = "flood"
	}

	if bdMap["vmac"] == "" {
		delete(bdMap, "vmac")
	}

	if len(dhcpLabel) == 0 {
		delete(bdMap, "dhcpLabel")
	}

	return &PatchPayload{
		Ops:   ops,
		Path:  path,
		Value: bdMap,
	}
}
