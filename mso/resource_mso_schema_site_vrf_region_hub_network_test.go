package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaSiteVRFRegionHubNetwork_Basic(t *testing.T) {
	var hubNetwork1 models.InterSchemaSiteVrfRegionHubNetork
	var hubNetwork2 models.InterSchemaSiteVrfRegionHubNetork
	resourceName := "mso_schema_site_vrf_region_hub_network.test"
	nameUpdated := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	regionName := makeTestVariable(acctest.RandString(5))
	vrfName := makeTestVariable(acctest.RandString(5))
	vrfNameUpdated := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVRFRegionHubNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "vrf_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "region_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], "tenant_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVRFRegionHubNetworkExists(resourceName, &hubNetwork1),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrfName),
					resource.TestCheckResourceAttr(resourceName, "region_name", regionName),
					resource.TestCheckResourceAttrSet(resourceName, "template_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfName, regionName, nameUpdated, tenantNames[0]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVRFRegionHubNetworkExists(resourceName, &hubNetwork2),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrfName),
					resource.TestCheckResourceAttr(resourceName, "region_name", regionName),
					resource.TestCheckResourceAttrSet(resourceName, "template_name"),
					testAccCheckMSOSchemaSiteVRFRegionHubNetworkIdNotEqual(&hubNetwork1, &hubNetwork2),
				),
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0]),
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfNameUpdated, regionName, name, tenantNames[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrfNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "region_name", regionName),
					resource.TestCheckResourceAttrSet(resourceName, "template_name"),
					testAccCheckMSOSchemaSiteVRFRegionHubNetworkIdNotEqual(&hubNetwork1, &hubNetwork2),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteVRFRegionHubNetwork_Negative(t *testing.T) {
	name := makeTestVariable(acctest.RandString(5))
	regionName := makeTestVariable(acctest.RandString(5))
	vrfName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVRFRegionHubNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfName, regionName, acctest.RandString(1001), tenantNames[0]),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSOSchemaSiteVRFRegionHubNetworkAttr(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0], randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteNames[0], schemaName, templateName, vrfName, regionName, name, tenantNames[0]),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteVRFRegionHubNetworkIdNotEqual(h1, h2 *models.InterSchemaSiteVrfRegionHubNetork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := hubNetworkModeltohubNetworkID(h1)
		id2 := hubNetworkModeltohubNetworkID(h2)
		if id1 == id2 {
			return fmt.Errorf("Schema Site VRF Region Hub Network Ids are equal")
		}
		return nil
	}
}

func testAccCheckMSOSchemaSiteVRFRegionHubNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_vrf_region_hub_network" {
			id := rs.Primary.ID
			hubNetwork, _ := hubNetworkIDtohubNetwork(id)
			_, err := client.ReadInterSchemaSiteVrfRegionHubNetwork(hubNetwork)
			if err == nil {
				return fmt.Errorf("Schema Site VRF Region Hub Network still exist")
			}
		}
	}
	return nil
}

func testAccCheckMSOSchemaSiteVRFRegionHubNetworkExists(hubNetwork string, m *models.InterSchemaSiteVrfRegionHubNetork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, ok := s.RootModule().Resources[hubNetwork]
		if !ok {
			return fmt.Errorf("Schema Site VRF Region Hub Network %s not found", hubNetwork)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schema Site VRF Region Hub Network Id was set")
		}
		hubNetworkModel, err := hubNetworkIDtohubNetwork(rs.Primary.ID)
		if err != nil {
			return err
		}
		var read *models.InterSchemaSiteVrfRegionHubNetork
		read, err = client.ReadInterSchemaSiteVrfRegionHubNetwork(hubNetworkModel)
		if err != nil {
			return err
		}
		*m = *read
		return nil
	}
}

func MSOSchemaSiteVrfRegionHubNetworkWithoutRequired(siteName, schemaName, templateName, vrfName, regionName, name, tenantName, attribute string) string {
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
	`
	switch attribute {
	case "schema_id":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
		#	schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "site_id":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
		#	site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "template_name":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
		#	template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "vrf_name":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
		#	vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "region_name":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
		#	region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "name":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
		#	name = "%s"
			tenant_name = data.mso_tenant.test.id
		}`
	case "tenant_name":
		rBlock += `resource "mso_schema_site_vrf_region_hub_network" "test"{
			schema_id = mso_schema_site.test.schema_id
			template_name = mso_schema_site.test.template_name
			site_id = mso_schema_site.test.site_id
			vrf_name = mso_schema_site_vrf.test.vrf_name
			region_name = mso_schema_site_vrf_region.test.region_name
			name = "%s"
		#	tenant_name = data.mso_tenant.test.id
		}`
	}
	return fmt.Sprintf(rBlock, vrfName, regionName, name)
}

func MSOSchemaSiteVRFRegionHubNetworkWithRequired(siteName, schemaName, templateName, vrfName, regionName, name, tenantName string) string {
	rBlock := CreatSchemaSiteConfig(siteName, tenantName, schemaName)
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
	}`
	return fmt.Sprintf(rBlock, vrfName, regionName, name)
}

func MSOSchemaSiteVRFRegionHubNetworkAttr(siteName, schemaName, templateName, vrfName, regionName, name, tenantName, key, value string) string {
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
		%s = "%s"
		tenant_name = data.mso_tenant.test.id
	}`
	return fmt.Sprintf(rBlock, vrfName, regionName, name, key, value)
}
