package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

func CreateInterSchemaSiteVrfRegionNetworkModel(hubNetwork *InterSchemaSiteVrfRegionHubNetork, cont *container.Container) *SchemaSiteVrfRegionHubNetork {
	vrfHubNetwork := SchemaSiteVrfRegionHubNetork{
		Ops:  "replace",
		Path: fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", hubNetwork.SiteID, hubNetwork.TemplateName, hubNetwork.VrfName, hubNetwork.Region),
	}
	vrfRegionMap, err := InterSchemaSiteVrfRegionFromContainer(cont, hubNetwork)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	vrfRegionMap["cloudRsCtxProfileToGatewayRouterP"] = map[string]string{
		"name":       hubNetwork.Name,
		"tenantName": hubNetwork.TenantName,
	}
	vrfHubNetwork.Value = vrfRegionMap
	return &vrfHubNetwork
}

func DeleteInterSchemaSiteVrfRegionNetworkModel(hubNetwork *InterSchemaSiteVrfRegionHubNetork, cont *container.Container) *SchemaSiteVrfRegionHubNetork {
	vrfHubNetwork := SchemaSiteVrfRegionHubNetork{
		Ops:  "replace",
		Path: fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", hubNetwork.SiteID, hubNetwork.TemplateName, hubNetwork.VrfName, hubNetwork.Region),
	}
	vrfHubNetworkMap := make(map[string]interface{})
	vrfHubNetworkMap["name"] = hubNetwork.Region
	vrfRegionMap, err := InterSchemaSiteVrfRegionFromContainer(cont, hubNetwork)
	if err != nil {
		return nil
	}
	vrfRegionMap["cloudRsCtxProfileToGatewayRouterP"] = nil
	vrfHubNetwork.Value = vrfRegionMap
	return &vrfHubNetwork
}

func InterSchemaSiteVrfRegionFromContainer(cont *container.Container, regionHubNetwork *InterSchemaSiteVrfRegionHubNetork) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	var found bool = false
	count, err := cont.ArrayCount("sites")
	log.Printf("Count: %v\n", count)
	if err != nil {
		return nil, fmt.Errorf("no sites found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		log.Printf("tempCont: %v\n", tempCont)
		if err != nil {
			log.Printf("error in tempCont: %v", err)
			return nil, err
		}
		siteId := StripQuotes(tempCont.S("siteId").String())
		log.Printf("siteId: %v\n", siteId)
		templateName := StripQuotes((tempCont.S("templateName")).String())
		log.Printf("templateName: %v\n", templateName)
		log.Printf("regionHubNetwork.SiteID: %v\n", regionHubNetwork.SiteID)
		log.Printf("regionHubNetwork.TemplateName: %v\n", regionHubNetwork.TemplateName)
		if (siteId == regionHubNetwork.SiteID) && (templateName == regionHubNetwork.TemplateName) {
			vrfCount, err := tempCont.ArrayCount("vrfs")
			log.Printf("vrfCount: %v\n", vrfCount)
			if err != nil {
				return nil, fmt.Errorf("unable to get VRFs List")
			}
			vrfCont := tempCont.S("vrfs")
			for j := 0; j < vrfCount; j++ {
				vrfTempCont := vrfCont.Index(j)
				log.Printf("vrfTempCont: %v\n", vrfTempCont)
				log.Printf("%v", StripQuotes(vrfTempCont.S("vrfRef").String()))
				vrfRef := strings.Split(StripQuotes(vrfTempCont.S("vrfRef").String()), "/")
				log.Printf("vrfRef: %v", vrfRef)
				vrfName := vrfRef[len(vrfRef)-1]
				log.Printf("vrfName: %v", vrfName)
				if vrfName == regionHubNetwork.VrfName {
					regionCount, err := vrfTempCont.ArrayCount("regions")
					log.Printf("regionCount: %v\n", regionCount)
					if err != nil {
						return nil, fmt.Errorf("unable to Regions List")
					}
					regionCont := vrfTempCont.S("regions")
					for k := 0; k < regionCount; k++ {
						regionTempCont := regionCont.Index(k)
						log.Printf("regionTempCont: %v\n", regionTempCont)
						regionName := G(regionTempCont, "name")
						if regionName == regionHubNetwork.Region {
							log.Printf("in if block of regionName: %v", regionTempCont)
							err := json.Unmarshal(regionTempCont.EncodeJSON(), &regionMap)
							if err != nil {
								log.Printf("error is: %v", err)
								return nil, err
							}
							found = true
							return regionMap, nil
						}
					}
				}
			}
		}
	}
	log.Printf("region Map holds value: %v", regionMap)
	if !found {
		return nil, fmt.Errorf("unable to find siteVrfRegionHubNetwork %s", regionHubNetwork.Name)
	}
	return regionMap, nil
}

func InterSchemaSiteVrfRegionHubNetworkFromContainer(cont *container.Container, regionHubNetwork *InterSchemaSiteVrfRegionHubNetork) (*InterSchemaSiteVrfRegionHubNetork, error) {
	hubNetwork := InterSchemaSiteVrfRegionHubNetork{}
	var found bool = false
	count, err := cont.ArrayCount("sites")
	log.Printf("Count: %v\n", count)
	if err != nil {
		return nil, fmt.Errorf("no sites found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		log.Printf("tempCont: %v\n", tempCont)
		if err != nil {
			return nil, err
		}
		siteId := StripQuotes(tempCont.S("siteId").String())
		log.Printf("siteId: %v\n", siteId)
		templateName := StripQuotes((tempCont.S("templateName")).String())
		log.Printf("templateName: %v\n", templateName)
		log.Printf("regionHubNetwork.SiteID: %v\n", regionHubNetwork.SiteID)
		log.Printf("regionHubNetwork.TemplateName: %v\n", regionHubNetwork.TemplateName)
		if (siteId == regionHubNetwork.SiteID) && (templateName == regionHubNetwork.TemplateName) {
			vrfCount, err := tempCont.ArrayCount("vrfs")
			log.Printf("vrfCount: %v\n", vrfCount)
			if err != nil {
				return nil, fmt.Errorf("unable to get VRFs List")
			}
			vrfCont := tempCont.S("vrfs")
			for j := 0; j < vrfCount; j++ {
				vrfTempCont := vrfCont.Index(j)
				log.Printf("vrfTempCont: %v\n", vrfTempCont)
				vrfRef := strings.Split(StripQuotes(vrfTempCont.S("vrfRef").String()), "/")
				vrfName := vrfRef[len(vrfRef)-1]
				if vrfName == regionHubNetwork.VrfName {
					regionCount, err := vrfTempCont.ArrayCount("regions")
					log.Printf("regionCount: %v\n", regionCount)
					if err != nil {
						return nil, fmt.Errorf("unable to Regions List")
					}
					regionCont := vrfTempCont.S("regions")
					for k := 0; k < regionCount; k++ {
						regionTempCont := regionCont.Index(k)
						log.Printf("regionTempCont: %v\n", regionTempCont)
						regionName := G(regionTempCont, "name")
						if regionName == regionHubNetwork.Region {
							routePCont := regionTempCont.S("cloudRsCtxProfileToGatewayRouterP")
							log.Printf("routePCont: %v\n", routePCont)
							hubName := StripQuotes(routePCont.S("name").String())
							tenantName := StripQuotes(routePCont.S("tenantName").String())
							log.Printf("hubName: %v\n", hubName)
							log.Printf("tenantName: %v\n", tenantName)
							log.Printf("regionHubNetwork.Name: %v\n", regionHubNetwork.Name)
							log.Printf("regionHubNetwork.TenantName: %v\n", regionHubNetwork.TenantName)
							if hubName == regionHubNetwork.Name && tenantName == regionHubNetwork.TenantName {
								log.Printf("in if container")
								hubNetwork.Name = hubName
								hubNetwork.TenantName = tenantName
								hubNetwork.Region = regionName
								hubNetwork.VrfName = vrfName
								hubNetwork.TemplateName = templateName
								hubNetwork.SiteID = siteId
								hubNetwork.SchemaID = vrfRef[2]
								found = true
								break
							}
						}
					}
				}
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("unable to find siteVrfRegionHubNetwork %s", regionHubNetwork.Name)
	}
	return &hubNetwork, nil
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
