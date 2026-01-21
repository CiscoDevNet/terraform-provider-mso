package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOL3OutInterfaceRoutingPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Interface Routing Policy with BFD") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigCreateWithBFD(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "description", "Test L3Out Interface Routing Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "uuid"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_receive_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_transmit_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_receive_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.interface_control", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Add OSPF") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigAddOSPF(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "description", "Updated with OSPF settings"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.network_type", "broadcast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.priority", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.cost_of_interface", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.hello_interval", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.dead_interval", "40"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.retransmit_interval", "5"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.transmit_delay", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.advertise_subnet", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.bfd", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.mtu_ignore", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.passive_participation", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Point-to-Point OSPF") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigPointToPoint(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.network_type", "point_to_point"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.priority", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.cost_of_interface", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.bfd", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Remove BFD Multi-Hop") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigRemoveBFDMultiHop(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.#", "0"),
					// BFD settings should still exist
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Remove All Settings") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigRemoveAllSettingsExceptOspf(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Maximum Values") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy_max"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.detection_multiplier", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_receive_interval", "999"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "999"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.detection_multiplier", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_receive_interval", "999"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_transmit_interval", "999"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_receive_interval", "999"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.priority", "255"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.cost_of_interface", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.hello_interval", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.dead_interval", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.retransmit_interval", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.transmit_delay", "450"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - Minimum Values") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy_min"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.detection_multiplier", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.detection_multiplier", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_receive_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_transmit_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_receive_interval", "50"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.priority", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.cost_of_interface", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.hello_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.dead_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.retransmit_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.transmit_delay", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy - OSPF with All Features") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigOSPFAllFeatures(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.network_type", "point_to_point"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.advertise_subnet", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.bfd", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.mtu_ignore", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.passive_participation", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Interface Routing Policy Name") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "description", "Renamed Policy"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import L3Out Interface Routing Policy") },
				ResourceName:      "mso_tenant_policies_l3out_interface_routing_policy.routing_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_l3out_interface_routing_policy", "l3OutIntfPolGroup"),
	})
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigCreateWithBFD() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "Test L3Out Interface Routing Policy"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigAddOSPF() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "Updated with OSPF settings"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = false
        }
        
        ospf_interface_settings {
            network_type          = "broadcast"
            priority              = 1
            cost_of_interface     = 0
            hello_interval        = 10
            dead_interval         = 40
            retransmit_interval   = 5
            transmit_delay        = 1
            advertise_subnet      = false
            bfd                   = false
            mtu_ignore            = false
            passive_participation = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigPointToPoint() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "Point-to-Point OSPF configuration"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = false
        }
        
        ospf_interface_settings {
            network_type      = "point_to_point"
            priority          = 10
            cost_of_interface = 100
            hello_interval    = 10
            dead_interval     = 40
            bfd               = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigRemoveBFDMultiHop() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "BFD Multi-Hop Removed"

        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = false
        }
        
        ospf_interface_settings {
            network_type      = "point_to_point"
            priority          = 10
            cost_of_interface = 100
            bfd               = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigRemoveAllSettingsExceptOspf() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "All Settings Removed"

		ospf_interface_settings {
            network_type      = "point_to_point"
            priority          = 10
            cost_of_interface = 100
            bfd               = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy_max"
        description = "Maximum Values Test"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 50
            min_receive_interval  = 999
            min_transmit_interval = 999
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 50
            min_receive_interval  = 999
            min_transmit_interval = 999
            echo_receive_interval = 999
            echo_admin_state      = "enabled"
            interface_control     = true
        }
        
        ospf_interface_settings {
            network_type          = "broadcast"
            priority              = 255
            cost_of_interface     = 65535
            hello_interval        = 65535
            dead_interval         = 65535
            retransmit_interval   = 65535
            transmit_delay        = 450
            advertise_subnet      = true
            bfd                   = true
            mtu_ignore            = true
            passive_participation = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigMinValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy_min"
        description = "Minimum Values Test"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 1
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 1
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = false
        }
        
        ospf_interface_settings {
            network_type          = "broadcast"
            priority              = 0
            cost_of_interface     = 0
            hello_interval        = 1
            dead_interval         = 1
            retransmit_interval   = 1
            transmit_delay        = 1
            advertise_subnet      = false
            bfd                   = false
            mtu_ignore            = false
            passive_participation = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigOSPFAllFeatures() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy"
        description = "OSPF with All Features Enabled"
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
            echo_receive_interval = 50
            echo_admin_state      = "enabled"
            interface_control     = true
        }
        
        ospf_interface_settings {
            network_type          = "point_to_point"
            priority              = 100
            cost_of_interface     = 10
            hello_interval        = 10
            dead_interval         = 40
            retransmit_interval   = 5
            transmit_delay        = 1
            advertise_subnet      = true
            bfd                   = true
            mtu_ignore            = true
            passive_participation = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOL3OutInterfaceRoutingPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_routing_policy_renamed"
        description = "Renamed Policy"
        
        bfd_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 50
            min_transmit_interval = 50
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}
