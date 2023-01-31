package models

import (
	"fmt"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type SiteAnpEpgUsegAttr struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

type SiteUsegAttr struct {
	SchemaID     string
	TemplateName string
	SiteID       string
	AnpName      string
	EpgName      string
	UsegName     string
	Description  string
	Type         string
	Operator     string
	Category     string
	Value        string
	FvSubnet     bool
}

func SiteAnpEpgUsegAttrForCreation(useg *SiteUsegAttr) *SiteAnpEpgUsegAttr {
	siteAnpEpgUsegAttr := SiteAnpEpgUsegAttr{
		Ops:  "add",
		Path: fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/uSegAttrs/-", useg.SiteID, useg.TemplateName, useg.AnpName, useg.EpgName),
	}

	usegAttr := map[string]interface{}{
		"name":        useg.UsegName,
		"displayName": useg.UsegName,
		"type":        useg.Type,
		"value":       useg.Value,
	}

	if StringInSlice(useg.Type, []string{"tag", "domain", "guest-os", "hv", "rootContName", "vm", "vm-name", "vnic"}) {
		usegAttr["operator"] = useg.Operator
	}

	if useg.Type == "tag" {
		usegAttr["category"] = useg.Category
	}

	if useg.Description != "" {
		usegAttr["description"] = useg.Description
	}

	if useg.Type == "ip" && useg.FvSubnet == true {
		usegAttr["fvSubnet"] = useg.FvSubnet
		usegAttr["value"] = "0.0.0.0"
	} else if useg.Type == "ip" {
		usegAttr["fvSubnet"] = useg.FvSubnet
	}

	siteAnpEpgUsegAttr.Value = usegAttr
	return &siteAnpEpgUsegAttr
}

func SiteAnpEpgUsegAttrforDeletion(useg *SiteUsegAttr, index int) *SiteAnpEpgUsegAttr {
	siteAnpEpgUsegAttr := SiteAnpEpgUsegAttr{
		Ops:  "remove",
		Path: fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/uSegAttrs/%d", useg.SiteID, useg.TemplateName, useg.AnpName, useg.EpgName, index),
	}
	return &siteAnpEpgUsegAttr
}

func SiteAnpEpgUsegAttrforUpdate(useg *SiteUsegAttr, index int) *SiteAnpEpgUsegAttr {
	siteAnpEpgUsegAttr := SiteAnpEpgUsegAttr{
		Ops:  "replace",
		Path: fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/uSegAttrs/%d", useg.SiteID, useg.TemplateName, useg.AnpName, useg.EpgName, index),
	}

	usegAttr := map[string]interface{}{
		"name":        useg.UsegName,
		"displayName": useg.UsegName,
		"type":        useg.Type,
		"value":       useg.Value,
	}

	if StringInSlice(useg.Type, []string{"tag", "domain", "guest-os", "hv", "rootContName", "vm", "vm-name"}) {
		usegAttr["operator"] = useg.Operator
	}

	if useg.Type == "tag" {
		usegAttr["category"] = useg.Category
	}

	if useg.Description != "" {
		usegAttr["description"] = useg.Description
	}

	if useg.Type == "ip" && useg.FvSubnet == true {
		usegAttr["fvSubnet"] = useg.FvSubnet
		usegAttr["value"] = "0.0.0.0"
	} else if useg.Type == "ip" {
		usegAttr["fvSubnet"] = useg.FvSubnet
	}

	siteAnpEpgUsegAttr.Value = usegAttr
	return &siteAnpEpgUsegAttr
}

func SiteAnpEpgUsegAttrFromContainer(cont *container.Container, tf *SiteUsegAttr) (*SiteUsegAttr, int, error) {
	siteUsegAttr := SiteUsegAttr{}
	siteUsegAttr.SchemaID = tf.SchemaID
	siteUsegAttr.SiteID = tf.SiteID
	siteUsegAttr.TemplateName = tf.TemplateName
	siteCont, err := cont.S("sites").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "siteId") == tf.SiteID && G(cont, "templateName") == tf.TemplateName
		},
	)
	if err != nil {
		return nil, -1, err
	}

	siteUsegAttr.AnpName = tf.AnpName
	anpCont, err := siteCont.S("anps").SearchInObjectList(
		func(cont *container.Container) bool {
			anpRef := G(cont, "anpRef")
			re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
			match := re.FindStringSubmatch(anpRef)
			anpName := match[3]
			return anpName == tf.AnpName
		},
	)
	if err != nil {
		return nil, -1, err
	}

	siteUsegAttr.EpgName = tf.EpgName
	epgCont, err := anpCont.S("epgs").SearchInObjectList(
		func(cont *container.Container) bool {
			epgRef := G(cont, "epgRef")
			re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
			match := re.FindStringSubmatch(epgRef)
			epgName := match[3]
			return epgName == tf.EpgName
		},
	)
	if err != nil {
		return nil, -1, err
	}

	siteUsegAttr.UsegName = tf.UsegName
	usegCont, useg_index, err := epgCont.S("uSegAttrs").SearchInObjectListWithIndex(
		func(cont *container.Container) bool {
			return G(cont, "name") == tf.UsegName
		},
	)
	if err != nil {
		return nil, -1, err
	}

	siteUsegAttr.Type = G(usegCont, "type")
	siteUsegAttr.Value = G(usegCont, "value")

	if StringInSlice(siteUsegAttr.Type, []string{"tag", "domain", "guest-os", "hv", "rootContName", "vm", "vm-name"}) {
		siteUsegAttr.Operator = G(usegCont, "operator")
	}

	if siteUsegAttr.Type == "tag" {
		siteUsegAttr.Category = G(usegCont, "category")
	}

	if usegCont.Exists("description") {
		siteUsegAttr.Description = G(usegCont, "description")
	}

	if siteUsegAttr.Type == "ip" && G(usegCont, "fvSubnet") == "true" {
		siteUsegAttr.FvSubnet = true
	}

	return &siteUsegAttr, useg_index, nil
}

func (useg *SiteAnpEpgUsegAttr) ToMap() (map[string]interface{}, error) {
	usegMap := make(map[string]interface{})
	A(usegMap, "op", useg.Ops)
	A(usegMap, "path", useg.Path)
	if useg.Value != nil {
		A(usegMap, "value", useg.Value)
	}
	return usegMap, nil
}
