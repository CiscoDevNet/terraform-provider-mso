package models

import (
	"fmt"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type SiteL3Out struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

type IntersiteL3outs struct {
	L3outName    string
	VRFName      string
	SchemaID     string
	TemplateName string
	SiteId       string
}

func CreateIntersiteL3outsModel(l3out *IntersiteL3outs) *SiteL3Out {
	site := SiteL3Out{
		Ops:  "add",
		Path: fmt.Sprintf("/sites/%s-%s/intersiteL3outs/-", l3out.SiteId, l3out.TemplateName),
	}
	sitemap := make(map[string]interface{})
	sitemap["l3outRef"] = map[string]string{
		"l3outName":    l3out.L3outName,
		"schemaId":     l3out.SchemaID,
		"templateName": l3out.TemplateName,
	}
	sitemap["vrfRef"] = map[string]string{
		"vrfName":      l3out.VRFName,
		"schemaId":     l3out.SchemaID,
		"templateName": l3out.TemplateName,
	}
	site.Value = sitemap
	return &site
}

func DeleteIntersiteL3outsModel(l3out *IntersiteL3outs) *SiteL3Out {
	site := SiteL3Out{
		Ops:  "remove",
		Path: fmt.Sprintf("/sites/%s-%s/intersiteL3outs/%s", l3out.SiteId, l3out.TemplateName, l3out.L3outName),
	}
	sitemap := make(map[string]interface{})
	sitemap["l3outRef"] = map[string]string{
		"l3outName":    l3out.L3outName,
		"schemaId":     l3out.SchemaID,
		"templateName": l3out.TemplateName,
	}
	site.Value = sitemap
	return &site
}

func IntersiteL3outsFromContainer(cont *container.Container, tf *IntersiteL3outs) (*IntersiteL3outs, error) {
	remoteL3out := IntersiteL3outs{}
	var found bool = false
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("no Sites found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := StripQuotes(tempCont.S("siteId").String())
		templateName := StripQuotes(tempCont.S("templateName").String())
		if apiSite == tf.SiteId && templateName == tf.TemplateName {
			l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
			if err != nil {
				return nil, fmt.Errorf("unable to get l3out list")
			}
			l3outCont := tempCont.S("intersiteL3outs")
			for j := 0; j < l3outCount; j++ {
				l3outTempCont := l3outCont.Index(j)
				l3outRef := strings.Split(StripQuotes(l3outTempCont.S("l3outRef").String()), "/")
				l3outName := l3outRef[len(l3outRef)-1]
				vrfRef := strings.Split(StripQuotes(l3outTempCont.S("vrfRef").String()), "/")
				vrfName := vrfRef[len(vrfRef)-1]
				if l3outName == tf.L3outName && vrfName == tf.VRFName {
					remoteL3out.L3outName = l3outName
					remoteL3out.VRFName = vrfName
					remoteL3out.SchemaID = l3outRef[2]
					remoteL3out.SiteId = apiSite
					remoteL3out.TemplateName = templateName
					found = true
					break
				}
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("unable to find siteL3out %s", tf.L3outName)
	}
	return &remoteL3out, nil
}

func (l3out *SiteL3Out) ToMap() (map[string]interface{}, error) {
	l3outMap := make(map[string]interface{})
	A(l3outMap, "op", l3out.Ops)
	A(l3outMap, "path", l3out.Path)
	if l3out.Value != nil {
		A(l3outMap, "value", l3out.Value)
	}
	return l3outMap, nil
}