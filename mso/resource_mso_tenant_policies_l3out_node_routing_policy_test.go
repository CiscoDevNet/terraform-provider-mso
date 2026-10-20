package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesL3OutNodeRoutingPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_l3out_node_routing_policy.node_policy"
	zeroHoldResourceName := "mso_tenant_policies_l3out_node_routing_policy.zero_hold_interval"
	asPathResourceName := "mso_tenant_policies_l3out_node_routing_policy.as_path_only"
	bestPathDefaultsResourceName := "mso_tenant_policies_l3out_node_routing_policy.best_path_defaults"
	bfdDefaultsResourceName := "mso_tenant_policies_l3out_node_routing_policy.bfd_defaults"
	bgpDefaultsResourceName := "mso_tenant_policies_l3out_node_routing_policy.bgp_defaults"
	var templateID, uuid string
	var originalTemplateID, originalUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create L3Out Node Routing Policy - Default BGP Best Path Control without Description")
				},
				Config: testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBGPBestPathDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "name", "test_node_routing_policy_best_path_defaults"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(bestPathDefaultsResourceName, "uuid"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Set Description") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigBGPBestPathDefaultsWithDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove Description") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBGPBestPathDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "description", ""),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(bestPathDefaultsResourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy - Default BFD Multi-Hop Settings Only") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBFDDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "name", "test_node_routing_policy_bfd_defaults"),
					resource.TestCheckResourceAttrSet(bfdDefaultsResourceName, "uuid"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bgp_best_path_control.#", "0"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(bfdDefaultsResourceName, "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy - Default BGP Node Settings Only") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBGPDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "name", "test_node_routing_policy_bgp_defaults"),
					resource.TestCheckResourceAttrSet(bgpDefaultsResourceName, "uuid"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_best_path_control.#", "0"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.0.stale_interval", "300"),
					resource.TestCheckResourceAttr(bgpDefaultsResourceName, "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy - Zero BGP Hold Interval") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateZeroHoldInterval(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(zeroHoldResourceName, "name", "test_node_routing_policy_zero_hold_interval"),
					resource.TestCheckResourceAttrSet(zeroHoldResourceName, "uuid"),
					resource.TestCheckResourceAttr(zeroHoldResourceName, "bgp_node_settings.0.graceful_restart_helper", "false"),
					resource.TestCheckResourceAttr(zeroHoldResourceName, "bgp_node_settings.0.keep_alive_interval", "0"),
					resource.TestCheckResourceAttr(zeroHoldResourceName, "bgp_node_settings.0.hold_interval", "0"),
					resource.TestCheckResourceAttr(zeroHoldResourceName, "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Create L3Out Node Routing Policy - AS Path Multipath Relax Disabled Only")
				},
				Config: testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateASPathOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(asPathResourceName, "name", "test_node_routing_policy_as_path_only"),
					resource.TestCheckResourceAttr(asPathResourceName, "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttrSet(asPathResourceName, "uuid"),
					resource.TestCheckResourceAttr(asPathResourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(asPathResourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
					resource.TestCheckResourceAttr(asPathResourceName, "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr(asPathResourceName, "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_l3out_node_routing_policy.node_policy", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_best_path_control.0.as_path_multipath_relax", "false"),

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
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						templateID = rs.Primary.Attributes["template_id"]
						uuid = rs.Primary.Attributes["uuid"]
						if templateID == "" || uuid == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate L3Out Node Routing Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "l3OutNodePolGroups"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "0"),
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
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove Disabled BGP Best Path Control") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBGPBestPathControl(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated BGP settings"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.graceful_restart_helper", "false"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.keep_alive_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.hold_interval", "90"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "180"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Add Default BGP Best Path Control") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEmptyBGPBestPathControl(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Enable AS Path Multipath Relax") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEnableASPath(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_best_path_control.0.as_path_multipath_relax", "true"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update L3Out Node Routing Policy - Retain Computed AS Path Multipath Relax")
				},
				Config: testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEmptyBGPBestPathControl(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Disable AS Path Multipath Relax") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDisableASPath(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Re-enable AS Path Multipath Relax") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEnableASPath(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove Enabled BGP Best Path Control") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBGPBestPathControl(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated BGP settings"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.graceful_restart_helper", "false"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.keep_alive_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.hold_interval", "90"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "180"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Add Disabled BGP Best Path Control") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDisableASPath(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "false"),
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
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Remove BFD and BGP Settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBFDAndBGP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Add Default BFD Multi-Hop Settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigAddBFDDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Add Default BGP Node Settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigAddBGPDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.0.min_transmit_interval", "250"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "0"),
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

					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Minimum Positive BGP Timers") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinPositiveBGPTimers(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "3"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out Node Routing Policy - Name") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Renamed Policy"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Prepare L3Out Node Routing Policy Template Replacement") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigTemplateChangeSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet("mso_template.template_tenant_replacement", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						originalTemplateID = rs.Primary.Attributes["template_id"]
						originalUUID = rs.Primary.Attributes["uuid"]
						if originalTemplateID == "" || originalUUID == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace L3Out Node Routing Policy - Template ID") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigTemplateChangeDestination(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant_replacement", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_node_routing_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Renamed Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "bgp_best_path_control.0.as_path_multipath_relax", "true"),
					resource.TestCheckResourceAttr(resourceName, "bfd_multi_hop_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.stale_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "bgp_node_settings.0.max_as_limit", "0"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						if rs.Primary.Attributes["uuid"] == originalUUID {
							return fmt.Errorf("resource %s retained its UUID after changing template ID", resourceName)
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", originalTemplateID))
						if err != nil {
							return fmt.Errorf("read original template after replacement: %w", err)
						}
						if _, err := GetPolicyIndexByKeyAndValue(response, "uuid", originalUUID, "tenantPolicyTemplate", "template", "l3OutNodePolGroups"); !isPolicyNotFound(err) {
							if err != nil {
								return fmt.Errorf("verify original policy was deleted during replacement: %w", err)
							}
							return fmt.Errorf("original policy with UUID %s still exists after replacement", originalUUID)
						}
						return nil
					},
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import L3Out Node Routing Policy") },
				ResourceName:      "mso_tenant_policies_l3out_node_routing_policy.node_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing L3Out Node Routing Policy") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_l3out_node_routing_policy", "l3OutNodePolGroup"),
	})
}

func TestAccMSOTenantPoliciesL3OutNodeRoutingPolicyResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject L3Out Node Routing Policy without template ID") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutTemplateID(),
				ExpectError: regexp.MustCompile(`The argument "template_id" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject L3Out Node Routing Policy without name") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject empty L3Out Node Routing Policy name") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorName(""),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject L3Out Node Routing Policy name longer than 64 characters") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorName(strings.Repeat("a", 65)),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject multiple BFD multi-hop settings blocks") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBFDSettings(),
				ExpectError: regexp.MustCompile(`Too many bfd_multi_hop_settings blocks`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject multiple BGP node settings blocks") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBGPSettings(),
				ExpectError: regexp.MustCompile(`Too many bgp_node_settings blocks`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject invalid BFD admin state") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorInvalidBFDAdminState(),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.admin_state to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD detection multiplier below 1") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDetectionMultiplier(0),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.detection_multiplier to be in the range \(1 - 50\), got 0`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD detection multiplier above 50") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDetectionMultiplier(51),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.detection_multiplier to be in the range \(1 - 50\), got 51`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD minimum receive interval below 250") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinReceiveInterval(249),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.min_receive_interval to be in the range \(250 - 999\), got 249`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD minimum receive interval above 999") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinReceiveInterval(1000),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.min_receive_interval to be in the range \(250 - 999\), got 1000`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD minimum transmit interval below 250") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinTransmitInterval(249),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.min_transmit_interval to be in the range \(250 - 999\), got 249`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BFD minimum transmit interval above 999") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinTransmitInterval(1000),
				ExpectError: regexp.MustCompile(`expected bfd_multi_hop_settings\.0\.min_transmit_interval to be in the range \(250 - 999\), got 1000`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject negative BGP keep-alive interval") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigKeepAliveInterval(-1),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.keep_alive_interval to be in the range \(0 - 3600\), got -1`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP keep-alive interval above 3600") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigKeepAliveInterval(3601),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.keep_alive_interval to be in the range \(0 - 3600\), got 3601`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject negative BGP hold interval") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigHoldInterval(-1),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.hold_interval to be in the range \(3 - 3600\), got -1`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP hold interval of 1") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigHoldInterval(1),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.hold_interval to be in the range \(3 - 3600\), got 1`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP hold interval of 2") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigHoldInterval(2),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.hold_interval to be in the range \(3 - 3600\), got 2`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP hold interval above 3600") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigHoldInterval(3601),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.hold_interval to be in the range \(3 - 3600\), got 3601`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP stale interval below 1") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigStaleInterval(0),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.stale_interval to be in the range \(1 - 3600\), got 0`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP stale interval above 3600") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigStaleInterval(3601),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.stale_interval to be in the range \(1 - 3600\), got 3601`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject negative BGP maximum AS limit") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxASLimit(-1),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.max_as_limit to be in the range \(0 - 2000\), got -1`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject BGP maximum AS limit above 2000") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxASLimit(2001),
				ExpectError: regexp.MustCompile(`expected bgp_node_settings\.0\.max_as_limit to be in the range \(0 - 2000\), got 2001`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject multiple BGP best path control blocks") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBGPBestPathControl(),
				ExpectError: regexp.MustCompile(`Too many bgp_best_path_control blocks`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out Node Routing Policy before removing all settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttrSet("mso_tenant_policies_l3out_node_routing_policy.node_policy", "uuid"),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject removing all L3Out Node Routing Policy settings") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveAll(),
				ExpectError: regexp.MustCompile("one of\\s+`bfd_multi_hop_settings,bgp_best_path_control,bgp_node_settings`\\s+must be\\s+specified"),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject L3Out Node Routing Policy without settings") },
				Config:      testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutSettings(),
				ExpectError: regexp.MustCompile("one of\\s+`bfd_multi_hop_settings,bgp_best_path_control,bgp_node_settings`\\s+must be\\s+specified"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Restore valid L3Out Node Routing Policy settings") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttrSet("mso_tenant_policies_l3out_node_routing_policy.node_policy", "uuid"),
			},
		},
	})
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBGPBestPathDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "best_path_defaults" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_best_path_defaults"
        bgp_best_path_control {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigBGPBestPathDefaultsWithDescription() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "best_path_defaults" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_best_path_defaults"
        description = "Test L3Out Node Routing Policy"

        bgp_best_path_control {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBFDDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "bfd_defaults" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_bfd_defaults"
        description = "Test L3Out Node Routing Policy"

        bfd_multi_hop_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateBGPDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "bgp_defaults" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_bgp_defaults"
        description = "Test L3Out Node Routing Policy"

        bgp_node_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateZeroHoldInterval() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "zero_hold_interval" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_zero_hold_interval"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }

        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 0
            hold_interval           = 0
            stale_interval          = 300
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreateASPathOnly() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "as_path_only" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_as_path_only"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutTemplateID() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorName(name string) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = %q
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }
    }`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBFDSettings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {}
        bfd_multi_hop_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBGPSettings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bgp_node_settings {}
        bgp_node_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorInvalidBFDAdminState() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            admin_state = "invalid"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDetectionMultiplier(detectionMultiplier int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            detection_multiplier = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), detectionMultiplier)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinReceiveInterval(minReceiveInterval int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            min_receive_interval = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), minReceiveInterval)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinTransmitInterval(minTransmitInterval int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            min_transmit_interval = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), minTransmitInterval)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigKeepAliveInterval(keepAliveInterval int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bgp_node_settings {
            keep_alive_interval = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), keepAliveInterval)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigStaleInterval(staleInterval int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bgp_node_settings {
            stale_interval = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), staleInterval)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxASLimit(maxASLimit int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bgp_node_settings {
            max_as_limit = %d
        }
    }`, testAccMSOTemplateResourceTenantConfig(), maxASLimit)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorWithoutSettings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "as_path_only" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_settings"
        description = "Test L3Out Node Routing Policy"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigErrorMultipleBGPBestPathControl() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bgp_best_path_control {
            as_path_multipath_relax = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {}
        bgp_node_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigHoldInterval(holdInterval int) string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Test L3Out Node Routing Policy"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 3
            min_receive_interval  = 250
            min_transmit_interval = 250
        }

        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = %d
            stale_interval          = 300
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig(), holdInterval)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateBGP() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "Updated BGP settings"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

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
        description            = "Updated BGP settings"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }
        
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

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigDisableASPath() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Updated BGP settings"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

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

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigEmptyBGPBestPathControl() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Updated BGP settings"

        bgp_best_path_control {}

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

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBGPBestPathControl() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Updated BGP settings"

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

        bgp_best_path_control {
            as_path_multipath_relax = true
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

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveBFDAndBGP() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "BFD and BGP Settings Removed"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigAddBFDDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Default BFD Multi-Hop Settings Added"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }

        bfd_multi_hop_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigAddBGPDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy"
        description = "Default BGP Node Settings Added"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }

        bfd_multi_hop_settings {}
        bgp_node_settings {}
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigRemoveAll() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy"
        description            = "All Settings Removed"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_max"
        description            = "Maximum Values Test"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }
        
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

        bgp_best_path_control {
            as_path_multipath_relax = false
        }
        
        bfd_multi_hop_settings {
            admin_state           = "enabled"
            detection_multiplier  = 1
            min_receive_interval  = 250
            min_transmit_interval = 250
        }
        
        bgp_node_settings {
            graceful_restart_helper = false
            keep_alive_interval     = 0
            hold_interval           = 0
            stale_interval          = 1
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigMinPositiveBGPTimers() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_min"
        description            = "Minimum Values Test"

        bgp_best_path_control {
            as_path_multipath_relax = false
        }

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
            max_as_limit            = 10
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id            = mso_template.template_tenant.id
        name                   = "test_node_routing_policy_renamed"
        description            = "Renamed Policy"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }
        
        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = 180
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigTemplateChangeSource() string {
	return fmt.Sprintf(`%s
    resource "mso_template" "template_tenant_replacement" {
        template_name = "test_template_tenant_replacement"
        template_type = "tenant"
        tenant_id     = mso_tenant.%s.id
        sites         = [data.mso_site.%s.id]
    }

    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_node_routing_policy_renamed"
        description = "Renamed Policy"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }

        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = 180
            stale_interval          = 1
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigTemplateChangeDestination() string {
	return fmt.Sprintf(`%s
    resource "mso_template" "template_tenant_replacement" {
        template_name = "test_template_tenant_replacement"
        template_type = "tenant"
        tenant_id     = mso_tenant.%s.id
        sites         = [data.mso_site.%s.id]
    }

    resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_template.template_tenant_replacement.id
        name        = "test_node_routing_policy_renamed"
        description = "Renamed Policy"

        bgp_best_path_control {
            as_path_multipath_relax = true
        }

        bgp_node_settings {
            graceful_restart_helper = true
            keep_alive_interval     = 60
            hold_interval           = 180
            stale_interval          = 1
            max_as_limit            = 0
        }
    }`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}
