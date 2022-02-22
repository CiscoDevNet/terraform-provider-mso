package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteVRFRegionHubNetworkDataSource_Basic(t *testing.T) {
	resourceName := "mso_schema_site_vrf_region_hub_network.test"
	dataSourceName := "data.mso_schema_site_vrf_region_hub_network.test"
	name := makeTestVariable(acctest.RandString(5))
	regionName := makeTestVariable(acctest.RandString(5))
	vrfName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	schemaID := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVRFRegionHubNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "vrf_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "region_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], "tenant_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVRFRegionHubNetworDataSourceAttr(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0], randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSOSchemaSiteVRFRegionHubNetworDataSourceInvalidName(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0]),
				ExpectError: regexp.MustCompile(`unable to find siteVrfRegionHubNetwork`),
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkDataSourceConfig(siteNames[0], schemaID, templateName, vrfName, regionName, name, tenantNames[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "tenant_name", dataSourceName, "tenant_name"),
					resource.TestCheckResourceAttrPair(resourceName, "template_name", dataSourceName, "template_name"),
					resource.TestCheckResourceAttrPair(resourceName, "region_name", dataSourceName, "region_name"),
					resource.TestCheckResourceAttrPair(resourceName, "vrf_name", dataSourceName, "vrf_name"),
					resource.TestCheckResourceAttrPair(resourceName, "site_id", dataSourceName, "site_id"),
					resource.TestCheckResourceAttrPair(resourceName, "schema_id", dataSourceName, "schema_id"),
				),
			},
		},
	})
}

func MSOSchemaSiteVrfRegionHubNetworkDataSourceWithoutRequired(siteName, schemaId, templateName, vrfName, regionName, name, tenantName, attribute string) string {
	rBlock := CreatSchemaSiteConfig(siteName, tenantName, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		schema_id = mso_schema_site.test.schema_id
		vrf_name = "%s"
	}
	
	resource "mso_schema_site_vrf_region" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = "%s"
		cidr {
     		cidr_ip = "2.2.2.2/10"
			primary = "true"
			subnet {
				ip = "1.20.30.4"
			}
		}
	}
	
	resource "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = "%s"
		tenant_name = data.mso_tenant.test.id
	}
	`
	switch attribute {
	case "schema_id":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
		#	schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "site_id":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
		#	site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "template_name":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
		#	template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "vrf_name":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
		#	vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "region_name":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
		#	region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "name":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
		#	name = mso_schema_site_vrf_region_hub_network.test.name
			tenant_name = data.mso_tenant.test.id
		}`
	case "tenant_name":
		rBlock += `data "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = mso_schema_site_vrf_region_hub_network.test.name
		#	tenant_name = data.mso_tenant.test.id
		}`
	}
	return fmt.Sprintf(rBlock, vrfName, regionName, name)
}

func MSOSchemaSiteVRFRegionHubNetworDataSourceAttr(siteName, schemaId, templateName, vrfName, regionName, name, tenantName, key, value string) string {
	rBlock := CreatSchemaSiteConfig(siteName, tenantName, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		schema_id = mso_schema_site.test.schema_id
		vrf_name = "%s"
	}
	
	resource "mso_schema_site_vrf_region" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = "%s"
		cidr {
     		cidr_ip = "2.2.2.2/10"
			primary = "true"
			subnet {
				ip = "1.20.30.4"
			}
		}
	}

	resource "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = "%s"
		tenant_name = data.mso_tenant.test.id
	}
	
	data "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = mso_schema_site_vrf_region_hub_network.test.name
		%s = "%s"
		tenant_name = data.mso_tenant.test.id
	}
	`
	return fmt.Sprintf(rBlock, vrfName, regionName, name, key, value)
}

func MSOSchemaSiteVRFRegionHubNetworDataSourceInvalidName(siteName, schemaId, templateName, vrfName, regionName, name, tenantName string) string {
	rBlock := CreatSchemaSiteConfig(siteName, tenantName, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		schema_id = mso_schema_site.test.schema_id
		vrf_name = "%s"
	}
	
	resource "mso_schema_site_vrf_region" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = "%s"
		cidr {
     		cidr_ip = "2.2.2.2/10"
			primary = "true"
			subnet {
				ip = "1.20.30.4"
			}
		}
	}

	resource "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = "%s"
		tenant_name = data.mso_tenant.test.id
	}
	
	data "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = "${mso_schema_site_vrf_region_hub_network.test.name}_invalid"
		tenant_name = data.mso_tenant.test.id
	}
	`
	return fmt.Sprintf(rBlock, vrfName, regionName, name)
}

func MSOSchemaSiteVRFRegionHubNetworkDataSourceConfig(siteName, schemaId, templateName, vrfName, regionName, name, tenantName string) string {
	rBlock := CreatSchemaSiteConfig(siteName, tenantName, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		schema_id = mso_schema_site.test.schema_id
		vrf_name = "%s"
	}
	
	resource "mso_schema_site_vrf_region" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = "%s"
		cidr {
     		cidr_ip = "2.2.2.2/10"
			primary = "true"
			subnet {
				ip = "1.20.30.4"
			}
		}
	}

	resource "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = "%s"
		tenant_name = data.mso_tenant.test.id
	}
	
	data "mso_schema_site_vrf_region_hub_network" "test"{
		schema_id = mso_schema_site.test.schema_id
		template_name = mso_schema_site.test.template_name
		site_id = mso_schema_site.test.site_id
		vrf_name = mso_schema_site_vrf.test.vrf_name
		region_name = mso_schema_site_vrf_region.test.region_name
		name = mso_schema_site_vrf_region_hub_network.test.name
		tenant_name = data.mso_tenant.test.id
	}
	`
	return fmt.Sprintf(rBlock, vrfName, regionName, name)
}
