package models

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type SchemaSiteVrfRegionHubNetork struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

type InterSchemaSiteVrfRegionHubNetork struct {
	Name         string
	TenantName   string
	SiteID       string
	TemplateName string
	VrfName      string
	Region       string
	SchemaID     string
}

func CreateInterSchemaSiteVrfRegionNetworkModel(hubNetwork *InterSchemaSiteVrfRegionHubNetork, cont *container.Container) (*SchemaSiteVrfRegionHubNetork, error) {
	vrfHubNetwork := SchemaSiteVrfRegionHubNetork{
		Ops:  "replace",
		Path: fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", hubNetwork.SiteID, hubNetwork.TemplateName, hubNetwork.VrfName, hubNetwork.Region),
	}
	vrfRegionMap, err := InterSchemaSiteVrfRegionFromContainer(cont, hubNetwork)
	if err != nil {
		return nil, fmt.Errorf("No VRF Region found")
	}
	vrfRegionMap["cloudRsCtxProfileToGatewayRouterP"] = map[string]string{
		"name":       hubNetwork.Name,
		"tenantName": hubNetwork.TenantName,
	}
	vrfHubNetwork.Value = vrfRegionMap
	return &vrfHubNetwork, nil
}

func DeleteInterSchemaSiteVrfRegionNetworkModel(hubNetwork *InterSchemaSiteVrfRegionHubNetork, cont *container.Container) (*SchemaSiteVrfRegionHubNetork, error) {
	vrfHubNetwork := SchemaSiteVrfRegionHubNetork{
		Ops:  "replace",
		Path: fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", hubNetwork.SiteID, hubNetwork.TemplateName, hubNetwork.VrfName, hubNetwork.Region),
	}
	vrfHubNetworkMap := make(map[string]interface{})
	vrfHubNetworkMap["name"] = hubNetwork.Region
	vrfRegionMap, err := InterSchemaSiteVrfRegionFromContainer(cont, hubNetwork)
	if err != nil {
		return nil, fmt.Errorf("No VRF Region found")
	}
	vrfRegionMap["cloudRsCtxProfileToGatewayRouterP"] = nil
	vrfHubNetwork.Value = vrfRegionMap
	return &vrfHubNetwork, nil
}

func InterSchemaSiteVrfRegionFromContainer(cont *container.Container, regionHubNetwork *InterSchemaSiteVrfRegionHubNetork) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	siteCont, err := cont.S("sites").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "siteId") == regionHubNetwork.SiteID && G(cont, "templateName") == regionHubNetwork.TemplateName
		},
	)
	if err != nil {
		return nil, err
	}
	vrfCont, err := siteCont.S("vrfs").SearchInObjectList(
		func(cont *container.Container) bool {
			vrfRef := G(cont, "vrfRef")
			re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
			match := re.FindStringSubmatch(vrfRef)
			vrfName := match[3]
			return vrfName == regionHubNetwork.VrfName
		},
	)
	if err != nil {
		return nil, err
	}
	regionCont, err := vrfCont.S("regions").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "name") == regionHubNetwork.Region
		},
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(regionCont.EncodeJSON(), &regionMap)
	if err != nil {
		return nil, err
	}
	return regionMap, nil
}

func InterSchemaSiteVrfRegionHubNetworkFromContainer(cont *container.Container, regionHubNetwork *InterSchemaSiteVrfRegionHubNetork) (*InterSchemaSiteVrfRegionHubNetork, error) {
	hubNetwork := InterSchemaSiteVrfRegionHubNetork{}
	hubNetwork.SiteID = regionHubNetwork.SiteID
	hubNetwork.TemplateName = regionHubNetwork.TemplateName
	hubNetwork.SchemaID = regionHubNetwork.SchemaID
	siteCont, err := cont.S("sites").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "siteId") == regionHubNetwork.SiteID && G(cont, "templateName") == regionHubNetwork.TemplateName
		},
	)
	if err != nil {
		return nil, err
	}
	hubNetwork.VrfName = regionHubNetwork.VrfName
	vrfCont, err := siteCont.S("vrfs").SearchInObjectList(
		func(cont *container.Container) bool {
			vrfRef := G(cont, "vrfRef")
			re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
			match := re.FindStringSubmatch(vrfRef)
			vrfName := match[3]
			return vrfName == regionHubNetwork.VrfName
		},
	)
	if err != nil {
		return nil, err
	}
	hubNetwork.Region = regionHubNetwork.Region
	regionCont, err := vrfCont.S("regions").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "name") == regionHubNetwork.Region
		},
	)
	if err != nil {
		return nil, err
	}
	hubNetwork.Name = regionHubNetwork.Name
	hubNetwork.TenantName = regionHubNetwork.TenantName
	if regionCont.Exists("cloudRsCtxProfileToGatewayRouterP") {
		hubNetworkCont := regionCont.S("cloudRsCtxProfileToGatewayRouterP")
		if G(hubNetworkCont, "name") == hubNetwork.Name && G(hubNetworkCont, "tenantName") == hubNetwork.TenantName {
			return &hubNetwork, nil
		}
	}
	return nil, fmt.Errorf("No Schema Site VRF Region Hub Network Found")
}

func (hubNetwork *SchemaSiteVrfRegionHubNetork) ToMap() (map[string]interface{}, error) {
	hubNetworkMap := make(map[string]interface{})
	A(hubNetworkMap, "op", hubNetwork.Ops)
	A(hubNetworkMap, "path", hubNetwork.Path)
	if hubNetwork.Value != nil {
		A(hubNetworkMap, "value", hubNetwork.Value)
	}
	return hubNetworkMap, nil
}
