package mso

import (
	"fmt"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName string, msoClient *client.Client) (*container.Container, error) {
	schemaObject, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	siteCount, err := schemaObject.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}

	for i := 0; i < siteCount; i++ {
		site, err := schemaObject.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}

		if models.G(site, "siteId") == siteId && models.G(site, "templateName") == templateName {
			return site, nil
		}
	}
	return nil, fmt.Errorf("Site-Template association for %v-%v is not found.", siteId, templateName)
}

func getSiteAnp(anpName string, site *container.Container) (*container.Container, error) {

	anpCount, err := site.ArrayCount("anps")
	if err != nil {
		return nil, fmt.Errorf("Unable to get ANP list")
	}
	for i := 0; i < anpCount; i++ {
		anpCont, err := site.ArrayElement(i, "anps")
		if err != nil {
			return nil, err
		}
		anpRef := models.G(anpCont, "anpRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
		match := re.FindStringSubmatch(anpRef)
		if match[3] == anpName {
			return anpCont, nil
		}
	}
	return nil, fmt.Errorf("ANP %v is not found in Site.", anpName)
}

func getSiteEpg(epgName string, anpCont *container.Container) (*container.Container, error) {

	epgCount, err := anpCont.ArrayCount("epgs")
	if err != nil {
		return nil, fmt.Errorf("Unable to get EPG list")
	}
	for i := 0; i < epgCount; i++ {
		epgCont, err := anpCont.ArrayElement(i, "epgs")
		if err != nil {
			return nil, err
		}
		epgRef := models.G(epgCont, "epgRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
		match := re.FindStringSubmatch(epgRef)
		if match[3] == epgName {
			return epgCont, nil
		}
	}
	return nil, fmt.Errorf("EPG %v is not found in Site.", epgName)
}

func getSiteBd(bdName string, site *container.Container) (*container.Container, error) {

	bdCount, err := site.ArrayCount("bds")
	if err != nil {
		return nil, fmt.Errorf("Unable to get BD list")
	}
	for i := 0; i < bdCount; i++ {
		bdCont, err := site.ArrayElement(i, "bds")
		if err != nil {
			return nil, err
		}
		bdRef := models.G(bdCont, "bdRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
		match := re.FindStringSubmatch(bdRef)
		if match[3] == bdName {
			return bdCont, nil
		}
	}
	return nil, fmt.Errorf("BD %v is not found in Site.", bdName)
}

func getSiteExternalEpg(externalEpgName string, site *container.Container) (*container.Container, error) {

	externalEpgCount, err := site.ArrayCount("externalEpgs")
	if err != nil {
		return nil, fmt.Errorf("Unable to get External EPG list")
	}
	for i := 0; i < externalEpgCount; i++ {
		externalEpgCont, err := site.ArrayElement(i, "externalEpgs")
		if err != nil {
			return nil, err
		}
		externalEpgRef := models.G(externalEpgCont, "externalEpgRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/externalEpgs/(.*)")
		match := re.FindStringSubmatch(externalEpgRef)
		if match[3] == externalEpgName {
			return externalEpgCont, nil
		}
	}
	return nil, fmt.Errorf("External EPG %v is not found in Site.", externalEpgName)
}

func getSiteVrf(vrfName string, site *container.Container) (*container.Container, error) {

	vrfCount, err := site.ArrayCount("vrfs")
	if err != nil {
		return nil, fmt.Errorf("Unable to get VRF list")
	}
	for i := 0; i < vrfCount; i++ {
		vrfCont, err := site.ArrayElement(i, "vrfs")
		if err != nil {
			return nil, err
		}
		vrfRef := models.G(vrfCont, "vrfRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
		match := re.FindStringSubmatch(vrfRef)
		if match[3] == vrfName {
			return vrfCont, nil
		}
	}
	return nil, fmt.Errorf("VRF %v is not found in Site.", vrfName)
}

func getSiteVrfRegion(regionName string, vrfCont *container.Container) (*container.Container, error) {

	regionCount, err := vrfCont.ArrayCount("regions")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Region list")
	}
	for i := 0; i < regionCount; i++ {
		regionCont, err := vrfCont.ArrayElement(i, "regions")
		if err != nil {
			return nil, err
		}
		matchRegion := models.StripQuotes(regionCont.S("name").String())
		if matchRegion == regionName {
			return regionCont, nil
		}
	}
	return nil, fmt.Errorf("VRF Region %v is not found in Site.", regionName)
}

func getSiteVrfRegionCIDR(ip string, regionCont *container.Container) (*container.Container, error) {

	cidrCount, err := regionCont.ArrayCount("cidrs")
	if err != nil {
		return nil, fmt.Errorf("Unable to get CIDR list")
	}

	for l := 0; l < cidrCount; l++ {
		cidrCont, err := regionCont.ArrayElement(l, "cidrs")
		if err != nil {
			return nil, err
		}
		matchIp := models.StripQuotes(cidrCont.S("ip").String())
		if matchIp == ip {
			return cidrCont, nil
		}
	}
	return nil, fmt.Errorf("VRF Region CIDR %v is not found in Site.", ip)
}
