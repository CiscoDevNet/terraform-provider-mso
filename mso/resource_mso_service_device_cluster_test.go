package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOServiceDeviceClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Device Cluster with one interface") },
				Config:    testAccMSOServiceDeviceClusterConfigCreateOneInterface(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "name", "test_device_cluster"),
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "device_mode", "layer3"),
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "device_type", "firewall"),
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "interface_properties.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":                  "interface1",
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"threshold_down_action": "permit",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster to three interfaces") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateThreeInterfaces(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "name", "test_device_cluster"),
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "interface_properties.#", "3"),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":                  "interface1",
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"threshold_down_action": "permit",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":                 "interface2",
						"load_balance_hashing": "destinationIP",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":    "interface3",
						"anycast": "true",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster to two interfaces") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateTwoInterfaces(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "name", "test_device_cluster"),
					resource.TestCheckResourceAttr("mso_service_device_cluster.cluster", "interface_properties.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":                  "interface1",
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"threshold_down_action": "permit",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":    "interface3",
						"anycast": "true",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Device Cluster") },
				ResourceName:      "mso_service_device_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMSOServiceDeviceClusterDependencies() string {
	return fmt.Sprintf(`%s
    resource "mso_template" "device_template" {
      template_name = "test_device_template"
      template_type = "service_device"
	  tenant_id     = mso_tenant.%s.id
    }

    resource "mso_template" "tenant_template" {
	  template_name = "test_tenant_template_for_device"
      template_type = "tenant"
	  tenant_id     = mso_tenant.%s.id
    }

    resource "mso_schema" "schema_blocks" {
		name = "demo_schema_blocks"
		template {
			name          = "Template1"
			display_name  = "TEMP1"
			tenant_id     = mso_tenant.%s.id
			template_type = "aci_multi_site"
		}
	}

	resource "mso_schema_template_vrf" "vrf" {
		schema_id    = mso_schema.schema_blocks.id
		template     = "Template1"
		name         = "template_vrf"
		display_name = "template_vrf"
	  }

    resource "mso_schema_template_bd" "bd1" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = "Template1"
        name          = "test_bd_1"
		vrf_name      = mso_schema_template_vrf.vrf.name
		display_name  = "template_bd1"
		arp_flooding  = true
    }

    resource "mso_schema_template_bd" "bd2" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = "Template1"
        name          = "test_bd_2"
		vrf_name      = mso_schema_template_vrf.vrf.name
		display_name  = "template_bd2"
		arp_flooding  = true
    }

    resource "mso_schema_template_external_epg" "epg1" {
        schema_id          = mso_schema.schema_blocks.id
        template_name      = "Template1"
        external_epg_name  = "test_epg_1"
		vrf_name           = mso_schema_template_vrf.vrf.name
		display_name       = "template_epg"
    }

    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla1" {
        template_id = mso_template.tenant_template.id
        name        = "test_ipsla_for_device"
        sla_type    = "icmp"
    }
`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateTenantName, msoTemplateTenantName)
}

func testAccMSOServiceDeviceClusterConfigCreateOneInterface() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                         = "interface1"
            external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "sourceIP"
            min_threshold                = 10
            max_threshold                = 90
            threshold_down_action        = "permit"
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

func testAccMSOServiceDeviceClusterConfigUpdateThreeInterfaces() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                         = "interface1"
            external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "sourceIP"
            min_threshold                = 10
            max_threshold                = 90
            threshold_down_action        = "permit"
        }
        interface_properties {
            name                         = "interface2"
            bd_uuid                      = mso_schema_template_bd.bd1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "destinationIP"
        }
        interface_properties {
            name                         = "interface3"
            bd_uuid                      = mso_schema_template_bd.bd2.uuid
            anycast                      = true
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

func testAccMSOServiceDeviceClusterConfigUpdateTwoInterfaces() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                         = "interface1"
            external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "sourceIP"
            min_threshold                = 10
            max_threshold                = 90
            threshold_down_action        = "permit"
        }
        interface_properties {
            name                         = "interface3"
            bd_uuid                      = mso_schema_template_bd.bd2.uuid
            anycast                      = true
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}
