package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateContractServiceChainingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Chaining with two nodes") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigCreateTwoNodes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "name", "contract1"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "node_filter", "allow-all"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface1",
						"provider_connector.0.is_redirect":    "false",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node2",
						"device_type":                         "firewall",
						"consumer_connector.0.interface_name": "interface",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface",
						"provider_connector.0.is_redirect":    "false",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Chaining to one node") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigUpdateOneNode(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "name", "contract1"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"provider_connector.0.interface_name": "interface1",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Chaining back to two nodes (reordered)") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigUpdateTwoNodesReordered(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "name", "contract1"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "firewall",
						"consumer_connector.0.interface_name": "interface",
						"provider_connector.0.interface_name": "interface",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node2",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"provider_connector.0.interface_name": "interface1",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Chaining") },
				ResourceName:      "mso_schema_template_contract_service_chaining.chain1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMSOSchemaTemplateContractServiceChainingDependencies() string {
	return fmt.Sprintf(`%s
    resource "mso_template" "device_template" {
      template_name = "test_device_template_sc"
      template_type = "service_device"
      tenant_id     = mso_tenant.%s.id
    }

    resource "mso_template" "tenant_template" {
      template_name = "test_tenant_template_sc"
      template_type = "tenant"
      tenant_id     = mso_tenant.%s.id
    }

    resource "mso_schema" "schema_blocks" {
      name = "demo_schema_sc"
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
      name         = "template_vrf_sc"
      display_name = "template_vrf_sc"
    }

    resource "mso_schema_template_bd" "bd1" {
      schema_id     = mso_schema.schema_blocks.id
      template_name = "Template1"
      name          = "test_bd_sc_1"
      vrf_name      = mso_schema_template_vrf.vrf.name
      display_name  = "template_bd_sc_1"
      arp_flooding  = true
    }

    resource "mso_schema_template_bd" "bd2" {
      schema_id     = mso_schema.schema_blocks.id
      template_name = "Template1"
      name          = "test_bd_sc_2"
      vrf_name      = mso_schema_template_vrf.vrf.name
      display_name  = "template_bd_sc_2"
      arp_flooding  = true
    }

    resource "mso_service_device_cluster" "device1" {
      template_id = mso_template.device_template.id
      name        = "device_cluster_lb_sc"
      device_mode = "layer3"
      device_type = "load_balancer"

      interface_properties {
        name    = "interface1"
        bd_uuid = mso_schema_template_bd.bd1.uuid
      }

      interface_properties {
        name    = "interface2"
        bd_uuid = mso_schema_template_bd.bd2.uuid
      }
    }

    resource "mso_service_device_cluster" "device2" {
      template_id = mso_service_device_cluster.device1.template_id
      name        = "device_cluster_fw_sc"
      device_mode = "layer3"
      device_type = "firewall"

      interface_properties {
        name    = "interface"
        bd_uuid = mso_schema_template_bd.bd1.uuid
      }
    }

	resource "mso_schema_template_filter_entry" "filter_entry" {
		schema_id            = mso_schema.schema_blocks.id
		template_name        = "Template1"
		name                 = "Filter1"
		display_name         = "Filter1"
		entry_name           = "entry1"
		entry_display_name   = "entry1"
		entry_description    = "DemoEntry"
		ether_type           = "arp"
		destination_from     = "unspecified"
		destination_to       = "unspecified"
		source_from          = "unspecified"
		source_to            = "unspecified"
		arp_flag             = "unspecified"
		stateful             = false
		match_only_fragments = false
	  }
	  
	  resource "mso_schema_template_contract" "contract1" {
		schema_id     = mso_schema.schema_blocks.id
		template_name = "Template1"
		contract_name = "contract1"
		display_name  = "contract1"
		filter_type   = "bothWay"
		scope         = "context"
		filter_relationship {
		  filter_schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
		  filter_template_name = mso_schema_template_filter_entry.filter_entry.template_name
		  filter_name          = mso_schema_template_filter_entry.filter_entry.name
		  filter_type          = "bothWay"
		}
	  }

	resource "mso_schema_template_external_epg" "template_externalepg" {
		schema_id                  = mso_schema.schema_blocks.id
		template_name              = "Template1"
		external_epg_name          = "external_epg1"
		display_name               = "external_epg1"
		vrf_name                   = mso_schema_template_vrf.vrf.name
		vrf_schema_id              = mso_schema_template_vrf.vrf.schema_id
		vrf_template_name          = "Template1"
		external_epg_type          = "on-premise"
		include_in_preferred_group = false
	  }

	  resource "mso_schema_template_external_epg_contract" "provider_contract" {
		schema_id                 = mso_schema.schema_blocks.id
		template_name             = "Template1"
		contract_name             = mso_schema_template_contract.contract1.contract_name
		external_epg_name         = mso_schema_template_external_epg.template_externalepg.external_epg_name
		relationship_type         = "provider"
		contract_schema_id        = mso_schema_template_contract.contract1.schema_id
		contract_template_name    = mso_schema_template_contract.contract1.template_name
	  }
`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateTenantName, msoTemplateTenantName)
}

func testAccMSOSchemaTemplateContractServiceChainingConfigCreateTwoNodes() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.schema_blocks.id
	template_name = "Template1"
	contract_name = mso_schema_template_contract.contract1.contract_name

	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.device1.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}

	service_nodes {
		name        = "node2"
		device_type = "firewall"
		device_ref  = mso_service_device_cluster.device2.uuid

		consumer_connector {
			interface_name = "interface"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface"
			is_redirect    = false
		}
	}
	}
`, testAccMSOSchemaTemplateContractServiceChainingDependencies())
}

func testAccMSOSchemaTemplateContractServiceChainingConfigUpdateOneNode() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.schema_blocks.id
	template_name = "Template1"
	contract_name = mso_schema_template_contract.contract1.contract_name
	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.device1.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}
	}
`, testAccMSOSchemaTemplateContractServiceChainingDependencies())
}

func testAccMSOSchemaTemplateContractServiceChainingConfigUpdateTwoNodesReordered() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.schema_blocks.id
	template_name = "Template1"
	contract_name = mso_schema_template_contract.contract1.contract_name

	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "firewall"
		device_ref  = mso_service_device_cluster.device2.uuid

		consumer_connector {
			interface_name = "interface"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface"
			is_redirect    = false
		}
	}

	service_nodes {
		name        = "node2"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.device1.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}
	}
`, testAccMSOSchemaTemplateContractServiceChainingDependencies())
}
