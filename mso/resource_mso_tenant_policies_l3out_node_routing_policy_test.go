package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesL3OutNodeRoutingPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy with BFD and BGP") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_l3out_node_routing_policy.node_policy", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "as_path_multipath_relax", "false"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "300"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Modify BGP Settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateBGP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Updated BGP settings"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.graceful_restart_helper", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "30"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "90"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "180"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Enable AS Path Multipath Relax") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEnableASPath(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "as_path_multipath_relax", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove BFD Multi-Hop") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBFD(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove All Settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Maximum Values") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy_max"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.detection_multiplier", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_receive_interval", "999"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "999"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "3599"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "3600"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "3600"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "2000"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Minimum Values") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy_min"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.detection_multiplier", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy Name") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Renamed Policy"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import L3Out Node Routing Policy") },
				ResourceName:      "mso_tenant_policies_l3out_node_routing_policy.node_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_l3out_node_routing_policy", "l3OutNodePolGroup"),
	})
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"
        as_path_multipath_relax = false
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = 180
            stale_interval          = 300
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateBGP() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Updated BGP settings"
        as_path_multipath_relax = false
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 30
            hold_interval           = 90
            stale_interval          = 180
            max_as_limit            = 10
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEnableASPath() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "AS Path Multipath Relax Enabled"
        as_path_multipath_relax = true
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 30
            hold_interval           = 90
            stale_interval          = 180
            max_as_limit            = 10
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBFD() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "BFD Multi-Hop Removed"
        as_path_multipath_relax = true
        
        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 30
            hold_interval           = 90
            stale_interval          = 180
            max_as_limit            = 10
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveAll() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "All Settings Removed"
        as_path_multipath_relax = true
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_max"
        description            = "Maximum Values Test"
        as_path_multipath_relax = true
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 50
            min_receive_interval  = 999
            min_transmit_interval = 999
        }
        
        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 3599
            hold_interval           = 3600
            stale_interval          = 3600
            max_as_limit            = 2000
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_min"
        description            = "Minimum Values Test"
        as_path_multipath_relax = false
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 1
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 1
            hold_interval           = 3
            stale_interval          = 1
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_renamed"
        description            = "Renamed Policy"
        as_path_multipath_relax = true
        
        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = 180
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}
