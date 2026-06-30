package mso

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Scope of the Service Device Cluster Site acceptance tests in this file.
//
// What IS exercised:
//   - domain_type = "physicalDomain" with a real fabric-policy physical
//     domain (mso_fabric_policies_physical_domain) and a vlan pool.
//   - device_mode = "layer3" on the cluster. high_availability_mode,
//     promiscuous_mode and trunking_port are NOT exercised in the live
//     lifecycle steps: high_availability_mode is rejected by NDO on
//     Layer 3 devices, and promiscuous_mode / trunking_port are
//     vmmDomain-specific and rejected on physicalDomain sites. Coverage
//     for those three flags is left to a future vmmDomain (and
//     Layer-1/2 for HA mode) site scenario; the SDK-validation error
//     test below still asserts the schema enum on high_availability_mode
//     against a placeholder template_id without hitting NDO.
//   - Cluster interfaces bound to BD UUIDs (bd_uuid) by default, with
//     interface3 in every multi-interface step bound to the schema
//     external EPG (external_epg_uuid) so both binding flavours are
//     covered against a physicalDomain. L3out (external_epg-bound)
//     interfaces omit the site-bucket `vlan` field because NDO rejects a
//     VLAN on L3out interfaces. Each step's mso_service_device_cluster
//     declares interface_properties names that match the site's interface
//     blocks 1:1, because NDO silently drops the site-bucket device-add
//     when the name sets differ.
//   - Cluster-level PBR/IPSLA wiring on the primary interface (entry 0,
//     "interface1") in every live step. Each step pairs a cluster
//     interface that has ipsla_monitoring_policy_uuid / thresholds /
//     load_balance_hashing / the PBR boolean toggles (config_static_mac,
//     is_backup_redirect_ip, resilient_hashing, tag_based_sorting) with a
//     site-bucket pbr_destination on interface1. The two must stay aligned
//     at every apply boundary because NDO server-side validation rejects
//     both directions of the cross-check: "At least one PBR destination
//     should be configured when redirect is enabled" and "PBR destination
//     cannot be configured when redirect is disabled". Terraform updates
//     the cluster before the site bucket, so introducing redirect or PBR
//     in an update apply leaves an intermediate inconsistent state that
//     fails validation; the create step therefore already provisions
//     both pieces together and the update step only flips attribute
//     values within the rich shape. Additional cluster interfaces stay
//     as bare bd1 bindings, mirroring the "Internal" + "External*"
//     docs example.
//   - fabric_to_device_connectivity with all three port_type values:
//     "port" (plain interfaceDn), "dpc" (single-node port-channel
//     policy-group, name written into pathep-[<name>]) and "vpc"
//     (two-node virtual port-channel, comma-joined nodeID in JSON +
//     hyphen-joined in protpaths URL). The create step uses "port"
//     only; the update step grows the list to all three so the in-place
//     update path is exercised across every NDO path encoding. The
//     supporting policy-group resources are wired into the deps stack
//     (mso_fabric_policies_interface_setting "portchannel",
//     mso_fabric_resource_policies_port_channel_interface and
//     mso_fabric_resource_policies_virtual_port_channel_interface on a
//     dedicated fabric_resource template + deploy_ndo).
//   - Lifecycle: create, in-place update (attributes + pbr_destination),
//     expand to three interfaces, shrink to two, ForceNew rename of the
//     cluster, drift recovery after manual NDO-side deletion, and import.
//
// What is NOT exercised (intentionally skipped — needs separate scenarios):
//   - domain_type = "vmmDomain" plus vmm_domain_type and vm_information /
//     enhanced_lag_policy on the interface.
//   - domain_dn variant (only the domain_type + domain_name form is tested).

// Captured from the create step's Check so the drift-recovery PreConfig can
// reach into the NDO API and remove the site-bucket device entry without
// re-parsing state.
var (
	testServiceDeviceClusterSiteTemplateID  string
	testServiceDeviceClusterSiteSiteID      string
	testServiceDeviceClusterSiteCurrentName string
)

const (
	testServiceDeviceClusterSiteClusterName        = "test_device_cluster"
	testServiceDeviceClusterSiteClusterNameRenamed = "test_device_cluster_renamed"
)

func TestAccMSOServiceDeviceClusterSiteResource(t *testing.T) {
	// Several steps in this test grow/shrink the cluster interface list or
	// rename the cluster, which fans out into two NDO calls in a single apply:
	// Destroy/Update on mso_service_device_cluster_site (the per-site bucket)
	// and Update on mso_service_device_cluster (the cluster-template scope).
	// NDO cross-validates the interface counts between those two scopes and
	// returns 400 "Mismatch in number of interfaces" if it sees them while
	// they disagree.
	//
	// Modern Terraform CLI (>= 0.13) adds DestroyEdgeTransformer2 which
	// inserts a stored-dependency edge so Destroy(cluster_site) is ordered
	// strictly before Update(cluster). A real `terraform apply` against this
	// configuration is therefore race-free.
	//
	// The vendored terraform-plugin-sdk v1 in-process driver predates that
	// transformer and walks the graph in parallel by default (Parallelism=10),
	// racing the two NDO calls and intermittently producing the 400 above.
	//
	// Workaround: a small LOCAL PATCH in
	// vendor/github.com/hashicorp/terraform-plugin-sdk/helper/resource/testing.go
	// LOCAL PATCH (not upstream): honor TF_ACC_PARALLELISM so a test can force
	// the in-process graph walker to run serially (=1). Used to work around
	// the SDK v1 missing DestroyEdgeTransformer2 race in mso_service_device_*.
	//
	//  if v := os.Getenv("TF_ACC_PARALLELISM"); v != "" {
	// 	    if p, perr := strconv.Atoi(v); perr == nil && p > 0 {
	// 		    opts.Parallelism = p
	// 	    }
	//  }
	//
	// honours TF_ACC_PARALLELISM so this test can pin the walker to 1. Run with:
	//
	//   TF_ACC=1 TF_ACC_PARALLELISM=1 \
	//     go test -run TestAccMSOServiceDeviceClusterSiteResource ./mso
	//
	// Without TF_ACC_PARALLELISM=1 the test is racy and will flake on the
	// interface-count grow/shrink/rename steps.

	resourceName := "mso_service_device_cluster_site.cluster_site"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// The deps-only intermediate step below uses waitForAPICMOs
			// (apic_test_helper_test.go) to gate the cluster_site apply on APIC
			// convergence, so the APIC env vars are required for this test.
			for _, env := range []string{envAPICURL, envAPICUsername, envAPICPassword} {
				if os.Getenv(env) == "" {
					t.Skipf("Skipping: %s must be set for the APIC settle Check used between the deploy chain and the cluster_site apply", env)
				}
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Intermediate "prereqs only" step: apply the full deps stack
			// (templates + deploys + cluster) but NOT cluster_site, then poll
			// APIC for every deployed MO that the cluster_site apply references
			// before letting the next step's apply fire. NDO's deploy task
			// returns "Complete" when the template is dispatched to APIC, not
			// when APIC has finished applying it; without this gate the
			// cluster_site PATCH races APIC and fails with "X does not exist
			// on the fabric ...".
			//
			// The DNs we wait for cover everything the update step refers to
			// by name:
			//   - uni/phys-test_physical_domain_for_device — the
			//     physical_domain referenced by the cluster_site's domain_name
			//     (used by every step, including the create step that runs
			//     right after this one).
			//   - uni/phys-test_physical_domain_for_device_alt — the second
			//     physical_domain referenced by the later "switch physical
			//     domain via domain_dn" step.
			//   - uni/infra/funcprof/accbundle-test_dpc_for_device — the PC
			//     policy group referenced by the dpc fabric_to_device_connectivity
			//     entry that the update step adds.
			//   - uni/infra/funcprof/accbundle-test_vpc_for_device — the VPC
			//     policy group referenced by the vpc fabric_to_device_connectivity
			//     entry that the update step adds.
			{
				PreConfig: func() {
					fmt.Println("Test: Apply prerequisites and wait for APIC deployed MOs to settle")
				},
				Config: testAccMSOServiceDeviceClusterSiteDependencies(testServiceDeviceClusterSiteClusterName, []siteClusterInterface{siteClusterRichInterface1}),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						return waitForAPICMOs(
							2*time.Minute,
							"uni/phys-test_physical_domain_for_device",
							"uni/phys-test_physical_domain_for_device_alt",
							"uni/infra/funcprof/accbundle-test_dpc_for_device",
							"uni/infra/funcprof/accbundle-test_vpc_for_device",
						)
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Service Device Cluster Site with one interface") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigOneInterface(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterName),
					resource.TestCheckResourceAttrSet(resourceName, "template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "1"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                      "210",
						"fabric_to_device_connectivity.#":           "1",
						"fabric_to_device_connectivity.0.pod_id":    "1",
						"fabric_to_device_connectivity.0.node_id.#": "1",
						"fabric_to_device_connectivity.0.node_id.0": "101",
						"fabric_to_device_connectivity.0.path":      "eth1/10",
						"fabric_to_device_connectivity.0.port_type": "port",
						"pbr_destinations.#":                        "1",
						"pbr_destinations.0.ip":                     "10.10.10.10",
						"pbr_destinations.0.mac":                    "00:11:22:33:44:55",
						"pbr_destinations.0.pod_id":                 "1",
						"pbr_destinations.0.additional_tracking_ip": "10.10.10.11",
						"pbr_destinations.0.weight":                 "5",
						"pbr_destinations.0.tag":                    "10",
					}),
					// Capture template_id, site_id and current name so the drift-recovery
					// PreConfig can target the device entry on the site bucket directly.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testServiceDeviceClusterSiteTemplateID = rs.Primary.Attributes["template_id"]
						testServiceDeviceClusterSiteSiteID = rs.Primary.Attributes["site_id"]
						testServiceDeviceClusterSiteCurrentName = rs.Primary.Attributes["name"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster Site attributes") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigUpdateAttributes(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "1"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                      "215",
						"fabric_to_device_connectivity.#":           "3",
						"fabric_to_device_connectivity.0.pod_id":    "1",
						"fabric_to_device_connectivity.0.node_id.#": "1",
						"fabric_to_device_connectivity.0.node_id.0": "101",
						"fabric_to_device_connectivity.0.path":      "eth1/10",
						"fabric_to_device_connectivity.0.port_type": "port",
						"fabric_to_device_connectivity.1.pod_id":    "1",
						"fabric_to_device_connectivity.1.node_id.#": "1",
						"fabric_to_device_connectivity.1.node_id.0": "101",
						"fabric_to_device_connectivity.1.path":      "test_dpc_for_device",
						"fabric_to_device_connectivity.1.port_type": "dpc",
						"fabric_to_device_connectivity.2.pod_id":    "1",
						"fabric_to_device_connectivity.2.node_id.#": "2",
						"fabric_to_device_connectivity.2.node_id.0": "101",
						"fabric_to_device_connectivity.2.node_id.1": "102",
						"fabric_to_device_connectivity.2.path":      "test_vpc_for_device",
						"fabric_to_device_connectivity.2.port_type": "vpc",
						"pbr_destinations.#":                        "1",
						"pbr_destinations.0.ip":                     "10.10.10.10",
						"pbr_destinations.0.mac":                    "00:11:22:33:44:66",
						"pbr_destinations.0.pod_id":                 "1",
						"pbr_destinations.0.additional_tracking_ip": "10.10.10.12",
						"pbr_destinations.0.weight":                 "7",
						"pbr_destinations.0.tag":                    "20",
					}),
				),
			},
			// This step grows the interface list, which is ForceNew on the site
			// resource and an in-place Update on the cluster in the same apply.
			// See the comment at the top of this function for why TF_ACC_PARALLELISM=1
			// is required to avoid the SDK v1 in-process walker racing the two NDO
			// calls on this step.
			{
				PreConfig: func() { fmt.Println("Test: Expand Service Device Cluster Site to three interfaces") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigThreeInterfaces(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "3"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                 "215",
						"fabric_to_device_connectivity.0.path": "eth1/10",
						"fabric_to_device_connectivity.0.port_type": "port",
						"pbr_destinations.0.ip":                     "10.10.10.10",
						"pbr_destinations.0.weight":                 "5",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface2",
					}, map[string]string{
						"vlan":                                 "216",
						"fabric_to_device_connectivity.0.path": "eth1/11",
						"fabric_to_device_connectivity.0.port_type": "port",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"fabric_to_device_connectivity.0.path":      "eth1/12",
						"fabric_to_device_connectivity.0.port_type": "port",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Service Device Cluster Site three-interface attribute changes")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigThreeInterfacesAttrs(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "3"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                      "220",
						"fabric_to_device_connectivity.0.path":      "eth1/10",
						"pbr_destinations.0.weight":                 "7",
						"pbr_destinations.0.tag":                    "20",
						"pbr_destinations.0.additional_tracking_ip": "10.10.10.12",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface2",
					}, map[string]string{
						"vlan":                                 "226",
						"fabric_to_device_connectivity.0.path": "eth1/21",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"fabric_to_device_connectivity.0.path": "eth1/22",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Shrink Service Device Cluster Site to two interfaces") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigTwoInterfaces(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                 "215",
						"fabric_to_device_connectivity.0.path": "eth1/10",
						"pbr_destinations.0.ip":                "10.10.10.10",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"fabric_to_device_connectivity.0.path": "eth1/12",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Service Device Cluster Site two-interface attribute changes")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesAttrs(testServiceDeviceClusterSiteClusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                 "230",
						"fabric_to_device_connectivity.0.path": "eth1/10",
						"pbr_destinations.0.mac":               "00:11:22:33:44:66",
						"pbr_destinations.0.weight":            "3",
					}),
					CustomTestCheckTypeSetElemAttrsByKeys(resourceName, "interfaces", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"fabric_to_device_connectivity.0.path": "eth1/32",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Rename cluster (site ForceNew destroy + recreate)") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigTwoInterfaces(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					// Refresh the captured name so the drift-recovery PreConfig
					// targets the device entry that lives on NDO after the rename.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testServiceDeviceClusterSiteCurrentName = rs.Primary.Attributes["name"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Service Device Cluster Site after manual deletion from NDO")
					if err := manuallyDeleteServiceDeviceClusterSiteFromTemplate(
						testServiceDeviceClusterSiteTemplateID,
						testServiceDeviceClusterSiteSiteID,
						testServiceDeviceClusterSiteCurrentName,
					); err != nil {
						panic(err)
					}
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigTwoInterfaces(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Device Cluster Site") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Apply the same two-interface set in swapped HCL order while leaving
			// the cluster's interface_properties order unchanged. The outer
			// interfaces block is ForceNew so this triggers a destroy+recreate of
			// the site bucket; if the site-to-cluster wiring were positional
			// instead of name-based the create PATCH would either be rejected by
			// NDO cross-validation or bind the rich PBR/IPSLA settings to the
			// wrong cluster interface_properties entry.
			{
				PreConfig: func() {
					fmt.Println("Test: Swap interfaces order in site config (name-based binding)")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesSwapped(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					// Positional state checks: HCL position 0 is the bare
					// external_epg-bound interface3, HCL position 1 is the rich
					// PBR/IPSLA interface1. The cluster's interface_properties
					// stay in [interface1, interface3] order.
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "interface3"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.fabric_to_device_connectivity.0.path", "eth1/12"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.pbr_destinations.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "interface1"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.vlan", "215"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.fabric_to_device_connectivity.0.path", "eth1/10"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.pbr_destinations.0.ip", "10.10.10.10"),
				),
			},
			// Switch the site-bucket domain reference from the (type, name)
			// pair to a literal domain_dn pointing at the alternate physical
			// domain. The site bucket's domain_type / domain_name attributes
			// are Computed, so NDO's response back-fills them from the DN;
			// asserting all three confirms the bucket landed on the alt
			// domain via every accessor.
			{
				PreConfig: func() {
					fmt.Println("Test: Switch physical domain via domain_dn input")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigSwitchDomainViaDN(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "domain_dn", "uni/phys-test_physical_domain_for_device_alt"),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test_physical_domain_for_device_alt"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
				),
			},
			// Swap back to (type, name) on the primary physical domain.
			// domain_dn is Computed and back-filled by NDO, so we assert the
			// full triple is consistent with the primary domain.
			{
				PreConfig: func() {
					fmt.Println("Test: Switch physical domain back via domain_type + domain_name input")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigSwitchDomainBackViaTypeName(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "domain_dn", "uni/phys-test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
				),
			},
			// Switch the cluster_site `name` to the layer-1 HA cluster and
			// flip the bucket into activeActive mode. Domain moves to
			// interface scope, so we assert each interface carries its own
			// domain attrs along with the pbr_destination mac/tag from the
			// JSON sample.
			{
				PreConfig: func() {
					fmt.Println("Test: activeActive HA mode with per-interface domain and pbr_destination")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigActiveActive(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteHaClusterName),
					resource.TestCheckResourceAttr(resourceName, "high_availability_mode", "activeActive"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "Internal"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.domain_name", "test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.fabric_to_device_connectivity.0.path", "eth1/1"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.fabric_to_device_connectivity.0.vlan", "210"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.fabric_to_device_connectivity.0.tag", "aaa"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.pbr_destinations.0.mac", "aa:bb:ee:cc:aa:dd"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.pbr_destinations.0.tag", "aaa"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "External"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.domain_dn", "uni/phys-test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.fabric_to_device_connectivity.0.path", "eth1/2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.fabric_to_device_connectivity.0.vlan", "220"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.pbr_destinations.0.mac", "cc:aa:cc:aa:cc:aa"),
				),
			},
			// Flip the bucket into activeStandby. Domain returns to device
			// scope (top-level domain_type + domain_name + vlan) and each
			// interface carries two fabric_to_device_connectivity entries
			// and two pbr_destinations (one per tag).
			{
				PreConfig: func() {
					fmt.Println("Test: activeStandby HA mode with device-level domain and per-tag interface entries")
				},
				Config: testAccMSOServiceDeviceClusterSiteConfigActiveStandby(testServiceDeviceClusterSiteClusterNameRenamed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterSiteHaClusterName),
					resource.TestCheckResourceAttr(resourceName, "high_availability_mode", "activeStandby"),
					resource.TestCheckResourceAttr(resourceName, "vlan", "230"),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "Internal"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.fabric_to_device_connectivity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.pbr_destinations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "External"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.fabric_to_device_connectivity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.pbr_destinations.#", "2"),
				),
			},
		},
	})
}

// TestAccMSOServiceDeviceClusterSiteResourceErrors exercises every
// input-validation path on the site resource. The step order is deliberate:
// SDK schema validation failures (StringInSlice / IntBetween / IsIPAddress /
// StringLenBetween) come first, and CustomizeDiff failures come last.
// terraform-plugin-sdk v1's post-test destroy step always re-runs
// ctx.Validate() against the last step's config (vendored
// helper/resource/testing.go around the destroy step block), so the last
// config must pass SDK schema validation; otherwise the destroy walk errors
// with "Error destroying resource ... config is invalid". CustomizeDiff
// failure configs have HCL that validates cleanly, so they are safe
// terminators.
func TestAccMSOServiceDeviceClusterSiteResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: invalid port_type enum value is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidPortType(),
				ExpectError: regexp.MustCompile(`expected interfaces\.\d+\.fabric_to_device_connectivity\.\d+\.port_type to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: vlan out of range is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrVlanOutOfRange(),
				ExpectError: regexp.MustCompile(`expected interfaces\.\d+\.vlan to be in the range \(1 - 4094\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: pbr_destination invalid IP is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidPbrIp(),
				ExpectError: regexp.MustCompile(`expected interfaces\.\d+\.pbr_destinations\.\d+\.ip to (be a valid IP|contain a valid IP)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid high_availability_mode enum value is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidHaMode(),
				ExpectError: regexp.MustCompile(`expected high_availability_mode to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: name longer than 64 characters is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrNameTooLong(),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid domain_type enum value is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidDomainType(),
				ExpectError: regexp.MustCompile(`expected domain_type to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid vmm_domain_type enum value is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidVmmDomainType(),
				ExpectError: regexp.MustCompile(`expected vmm_domain_type to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: pbr_destination weight out of range is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrPbrWeightOutOfRange(),
				ExpectError: regexp.MustCompile(`expected interfaces\.\d+\.pbr_destinations\.\d+\.weight to be in the range \(1 - 10\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: fabric_to_device_connectivity vlan out of range is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrFabricPathVlanOutOfRange(),
				ExpectError: regexp.MustCompile(`expected interfaces\.\d+\.fabric_to_device_connectivity\.\d+\.vlan to be in the range \(1 - 4094\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: site without domain_type or domain_dn is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrNoDomain(),
				ExpectError: regexp.MustCompile(`one of domain_type or domain_dn must be set`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: vmmDomain without vmm_domain_type is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrVmmDomainMissingVmmType(),
				ExpectError: regexp.MustCompile(`vmm_domain_type is required when domain_type is "vmmDomain"`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: physicalDomain combined with vmm_domain_type is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithVmmType(),
				ExpectError: regexp.MustCompile(`vmm_domain_type must not be set when domain_type is "physicalDomain"`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid domain_dn prefix is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInvalidDomainDn(),
				ExpectError: regexp.MustCompile(`domain_dn .* must start with "uni/phys-" or "uni/vmmp-"`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: interface with both fabric and vm blocks is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInterfaceBothFabricAndVm(),
				ExpectError: regexp.MustCompile(`fabric_to_device_connectivity and vm_information are mutually exclusive`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: interface with neither fabric nor vm block is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInterfaceNoConnectivity(),
				ExpectError: regexp.MustCompile(`one of fabric_to_device_connectivity or vm_information must be set`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: vm_information on a physical domain is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithVmInformation(),
				ExpectError: regexp.MustCompile(`vm_information is not allowed when the device uses a physicalDomain`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: enhanced_lag_policy on a physical domain is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithEnhancedLag(),
				ExpectError: regexp.MustCompile(`enhanced_lag_policy is not allowed when the device uses a physicalDomain`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: fabric_to_device_connectivity on a vmm domain is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrVmmWithFabricConnectivity(),
				ExpectError: regexp.MustCompile(`fabric_to_device_connectivity is not allowed when the device uses a vmmDomain`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: vpc port_type with a single node_id entry is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrVpcWithSingleNode(),
				ExpectError: regexp.MustCompile(`port_type "vpc" requires exactly two node_id entries`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: port port_type with two node_id entries is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrPortWithVpcStyleNode(),
				ExpectError: regexp.MustCompile(`port_type "port" requires exactly one node_id entry`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: dpc port_type with two node_id entries is rejected") },
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrDpcWithTwoNodes(),
				ExpectError: regexp.MustCompile(`port_type "dpc" requires exactly one node_id entry`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: activeActive without per-interface domain is rejected")
				},
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrActiveActiveMissingInterfaceDomain(),
				ExpectError: regexp.MustCompile(`domain must be configured on every interface when high_availability_mode is "activeActive"`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: activeActive interface-scoped physicalDomain with vm_information is rejected")
				},
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedPhysicalWithVm(),
				ExpectError: regexp.MustCompile(`vm_information is not allowed when the interface uses a physicalDomain`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: activeActive interface-scoped physicalDomain with enhanced_lag_policy is rejected")
				},
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedPhysicalWithEnhancedLag(),
				ExpectError: regexp.MustCompile(`enhanced_lag_policy is not allowed when the interface uses a physicalDomain`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: activeActive interface-scoped vmmDomain with fabric_to_device_connectivity is rejected")
				},
				Config:      testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedVmmWithFabric(),
				ExpectError: regexp.MustCompile(`fabric_to_device_connectivity is not allowed when the interface uses a vmmDomain`),
			},
		},
	})
}

// manuallyDeleteServiceDeviceClusterSiteFromTemplate removes the site-bucket
// device entry from a service device template via the NDO API. Used by the
// drift-recovery test step to simulate an out-of-band deletion before the same
// Terraform configuration is re-applied.
func manuallyDeleteServiceDeviceClusterSiteFromTemplate(templateID, siteID, deviceName string) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return fmt.Errorf("manual delete: get template %s: %w", templateID, err)
	}
	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteID, "deviceTemplate", "sites")
	if err != nil {
		return fmt.Errorf("manual delete: locate site %q: %w", siteID, err)
	}
	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)
	deviceIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "name", deviceName, "devices")
	if err != nil {
		return fmt.Errorf("manual delete: locate device %q on site %q: %w", deviceName, siteID, err)
	}
	removePayload := models.GetRemovePatchPayload(fmt.Sprintf("/deviceTemplate/sites/%d/devices/%d", siteIndex, deviceIndex))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), removePayload); err != nil {
		return fmt.Errorf("manual delete: patch remove site %q device %q: %w", siteID, deviceName, err)
	}
	return nil
}

// siteClusterInterface describes a single mso_service_device_cluster
// interface_properties block emitted by
// testAccMSOServiceDeviceClusterSiteDependencies. Only Name is required;
// every other field is opt-in so a caller can pair a rich PBR/IPSLA entry
// (mirroring the "Internal" interface from the resource docs) with bare
// entries (mirroring "External1"/"External2") in the same cluster. Zero
// values are skipped from the rendered HCL so NDO defaults apply for unset
// fields. Entries default to binding bd1.uuid; set WithExternalEPG to bind
// the interface to mso_schema_template_external_epg.epg1.uuid instead.
type siteClusterInterface struct {
	Name                string
	WithIPSLA           bool
	LoadBalanceHashing  string
	MinThreshold        int
	MaxThreshold        int
	ThresholdDownAction string
	ConfigStaticMAC     bool
	IsBackupRedirectIP  bool
	ResilientHashing    bool
	TagBasedSorting     bool
	WithExternalEPG     bool
	// WithRedirect emits `redirect = true` on the cluster interface_properties
	// block, setting NDO's interface-level redirect flag directly. NDO
	// requires this before it accepts a matching site-bucket
	// pbr_destination on the same interface. Used by the layer-1 HA cluster
	// (cluster_ha) for activeActive / activeStandby site tests where the
	// cluster has no other signal that would auto-derive redirect.
	WithRedirect bool
}

// renderClusterInterfaceProperties returns the interface_properties HCL
// blocks for one mso_service_device_cluster resource, applying the same
// opt-in emission rules as the inline loop it replaces. Used twice in
// testAccMSOServiceDeviceClusterSiteDependencies: once for the layer-3
// `cluster` (parameterised per step) and once for the layer-1 HA
// `cluster_ha` (fixed Internal/External shape from siteClusterHaInterfaces).
func renderClusterInterfaceProperties(interfaces []siteClusterInterface) string {
	var b strings.Builder
	for _, iface := range interfaces {
		fmt.Fprintf(&b, "        interface_properties {\n")
		fmt.Fprintf(&b, "            name    = %q\n", iface.Name)
		if iface.WithExternalEPG {
			fmt.Fprintf(&b, "            external_epg_uuid = mso_schema_template_external_epg.epg1.uuid\n")
		} else {
			fmt.Fprintf(&b, "            bd_uuid = mso_schema_template_bd.bd1.uuid\n")
		}
		if iface.WithIPSLA {
			fmt.Fprintf(&b, "            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid\n")
		}
		if iface.LoadBalanceHashing != "" {
			fmt.Fprintf(&b, "            load_balance_hashing = %q\n", iface.LoadBalanceHashing)
		}
		if iface.MinThreshold != 0 {
			fmt.Fprintf(&b, "            min_threshold = %d\n", iface.MinThreshold)
		}
		if iface.MaxThreshold != 0 {
			fmt.Fprintf(&b, "            max_threshold = %d\n", iface.MaxThreshold)
		}
		if iface.ThresholdDownAction != "" {
			fmt.Fprintf(&b, "            threshold_down_action = %q\n", iface.ThresholdDownAction)
		}
		if iface.ConfigStaticMAC {
			fmt.Fprintf(&b, "            config_static_mac = true\n")
		}
		if iface.IsBackupRedirectIP {
			fmt.Fprintf(&b, "            is_backup_redirect_ip = true\n")
		}
		if iface.ResilientHashing {
			fmt.Fprintf(&b, "            resilient_hashing = true\n")
		}
		if iface.TagBasedSorting {
			fmt.Fprintf(&b, "            tag_based_sorting = true\n")
		}
		if iface.WithRedirect {
			fmt.Fprintf(&b, "            redirect = true\n")
		}
		fmt.Fprintf(&b, "        }\n")
	}
	return b.String()
}

// siteClusterRichInterface1 is the canonical "primary" cluster interface
// used by every site-bucket config step in this file. The IPSLA / threshold
// / PBR-boolean settings enable redirect on the cluster interface, and NDO
// server-side validation requires a matching site-bucket pbr_destination on
// interface1 whenever redirect is enabled (and forbids one when it is not).
// Because terraform updates the cluster before the site bucket, introducing
// or removing either piece during an update leaves a transient inconsistent
// state that fails validation; the create step therefore provisions both
// the rich cluster interface and the site-bucket pbr_destination together,
// and every subsequent step keeps them aligned.
var siteClusterRichInterface1 = siteClusterInterface{
	Name:                "interface1",
	WithIPSLA:           true,
	LoadBalanceHashing:  "sourceIP",
	MinThreshold:        3,
	MaxThreshold:        100,
	ThresholdDownAction: "permit",
	ConfigStaticMAC:     true,
	IsBackupRedirectIP:  true,
	ResilientHashing:    true,
	TagBasedSorting:     true,
}

// testServiceDeviceClusterSiteHaClusterName is the name of the layer-1 HA
// cluster (cluster_ha) provisioned unconditionally by
// testAccMSOServiceDeviceClusterSiteDependencies. The activeActive and
// activeStandby site-bucket steps switch the cluster_site `name` from the
// layer-3 `cluster` to this layer-1 cluster.
const testServiceDeviceClusterSiteHaClusterName = "test_device_cluster_ha"

// siteClusterHaInterfaces is the fixed interface_properties shape on the
// layer-1 HA cluster (cluster_ha) used by the activeActive and
// activeStandby site-bucket steps. Names "Internal" and "External" match
// the JSON sample's device interfaces. WithRedirect sets `redirect = true`
// on each cluster interface_properties block, which NDO requires before
// it accepts the matching site-bucket pbr_destinations carried by those
// HA steps.
var siteClusterHaInterfaces = []siteClusterInterface{
	{Name: "Internal", WithRedirect: true},
	{Name: "External", WithRedirect: true},
}

// testAccMSOServiceDeviceClusterSiteDependencies is the self-contained
// prerequisite stack for the site-bucket tests.
//
// The stack provisioned here covers a physicalDomain site with both
// bd_uuid- and external_epg_uuid-bound cluster interfaces:
//   - tenant (via testAccTenantConfig)
//   - fabric_policy_template + vlan_pool + physical_domain
//   - tenant_template + ipsla monitoring policy
//   - schema with vrf, bd1, bd2, l3out and external_epg (the L3Out and
//     extEPG share the template VRF; any siteClusterInterface entry
//     flagged with WithExternalEPG binds to extEPG.uuid instead of
//     bd1.uuid)
//   - service_device template, schema_site, three sequential deploys
//   - mso_service_device_cluster.cluster with interface_properties built
//     from the interfaces []siteClusterInterface argument; each entry
//     binds to bd1.uuid by default (or to extEPG.uuid when
//     WithExternalEPG is set) and any opt-in PBR/IPSLA fields are
//     emitted unconditionally.
//
// Every NDO-named object uses a "_site" suffix to avoid colliding with the
// cluster acceptance test, which provisions its own service-device template
// against the same shared tenant.
func testAccMSOServiceDeviceClusterSiteDependencies(clusterName string, interfaces []siteClusterInterface) string {
	interfaceAttributes := renderClusterInterfaceProperties(interfaces)
	haInterfaceAttributes := renderClusterInterfaceProperties(siteClusterHaInterfaces)

	return fmt.Sprintf(`%[1]s
    resource "mso_template" "fabric_policy_template" {
        template_name = "test_fabric_policy_for_site_device"
        template_type = "fabric_policy"
        sites         = [data.mso_site.%[2]s.id]
        depends_on    = [mso_tenant.%[3]s]
    }

    resource "mso_template" "tenant_template" {
        template_name = "test_tenant_template_for_site_device"
        template_type = "tenant"
        tenant_id     = mso_tenant.%[3]s.id
        sites         = [data.mso_site.%[2]s.id]
        # Explicit tenant dep (in addition to the implicit tenant_id ref) so the
        # tenant is guaranteed to be destroyed after this template.
        depends_on    = [mso_tenant.%[3]s, mso_template.fabric_policy_template]
    }

    resource "mso_schema" "schema_blocks" {
        name = "demo_schema_blocks_site"
        template {
            name          = "Template1"
            display_name  = "TEMP1"
            tenant_id     = mso_tenant.%[3]s.id
            template_type = "aci_multi_site"
        }
        # Explicit tenant dep so the tenant is guaranteed to be destroyed after
        # this schema (the tenant_id reference inside the template block alone
        # is not always honoured for destroy ordering in the legacy SDK).
        depends_on = [mso_tenant.%[3]s, mso_template.tenant_template]
    }

    resource "mso_template" "device_template" {
        template_name = "test_device_template_for_site"
        template_type = "service_device"
        tenant_id     = mso_tenant.%[3]s.id
        sites         = [data.mso_site.%[2]s.id]
        # Explicit tenant dep (in addition to the implicit tenant_id ref) so the
        # tenant is guaranteed to be destroyed after this template.
        depends_on    = [mso_tenant.%[3]s, mso_schema.schema_blocks]
    }

    resource "mso_schema_template_vrf" "vrf" {
        schema_id    = mso_schema.schema_blocks.id
        template     = tolist(mso_schema.schema_blocks.template)[0].name
        name         = "template_vrf"
        display_name = "template_vrf"
    }

    resource "mso_schema_template_bd" "bd1" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = tolist(mso_schema.schema_blocks.template)[0].name
        name          = "test_bd_1"
        vrf_name      = mso_schema_template_vrf.vrf.name
        display_name  = "template_bd1"
        arp_flooding  = true
    }

    resource "mso_schema_template_bd" "bd2" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = mso_schema_template_bd.bd1.template_name
        name          = "test_bd_2"
        vrf_name      = mso_schema_template_vrf.vrf.name
        display_name  = "template_bd2"
        arp_flooding  = true
    }

    resource "mso_schema_template_l3out" "l3out1" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = mso_schema_template_bd.bd1.template_name
        l3out_name    = "test_l3out_for_site_device"
        display_name  = "template_l3out_site"
        vrf_name      = mso_schema_template_vrf.vrf.name
    }

    resource "mso_schema_template_external_epg" "epg1" {
        schema_id           = mso_schema.schema_blocks.id
        template_name       = mso_schema_template_bd.bd1.template_name
        external_epg_name   = "test_extepg_for_site_device"
        vrf_name            = mso_schema_template_vrf.vrf.name
        display_name        = "template_extepg_site"
        external_epg_type   = "on-premise"
        l3out_name          = mso_schema_template_l3out.l3out1.l3out_name
        l3out_schema_id     = mso_schema_template_l3out.l3out1.schema_id
        l3out_template_name = mso_schema_template_l3out.l3out1.template_name
    }

    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla1" {
        template_id = mso_template.tenant_template.id
        name        = "test_ipsla_for_site_device"
        sla_type    = "icmp"
    }

    resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
        template_id = mso_template.fabric_policy_template.id
        name        = "test_vlan_pool_for_device"
        vlan_range {
            from = 200
            to   = 250
        }
    }

    resource "mso_fabric_policies_physical_domain" "physical_domain" {
        template_id    = mso_template.fabric_policy_template.id
        name           = "test_physical_domain_for_device"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }

    # Second physical domain used by the "switch physical domain" test
    # steps. Shares the existing vlan pool. The site-bucket steps reach
    # this domain via either the type/name pair or the literal domain_dn
    # "uni/phys-test_physical_domain_for_device_alt".
    resource "mso_fabric_policies_physical_domain" "physical_domain_alt" {
        template_id    = mso_template.fabric_policy_template.id
        name           = "test_physical_domain_for_device_alt"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }

    resource "mso_fabric_policies_interface_setting" "portchannel_setting" {
        template_id = mso_template.fabric_policy_template.id
        name        = "test_portchannel_setting_for_device"
        type        = "portchannel"
    }

    resource "mso_template" "fabric_resource_template" {
        template_name = "test_fabric_resource_for_site_device"
        template_type = "fabric_resource"
        sites         = [data.mso_site.%[2]s.id]
        depends_on    = [mso_template.fabric_policy_template]
    }

    resource "mso_fabric_resource_policies_port_channel_interface" "pc_interface" {
        template_id                 = mso_template.fabric_resource_template.id
        name                        = "test_dpc_for_device"
        node                        = "101"
        interfaces                  = ["1/15", "1/16"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.portchannel_setting.uuid
    }

    resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_interface" {
        template_id                 = mso_template.fabric_resource_template.id
        name                        = "test_vpc_for_device"
        node_1                      = "101"
        node_2                      = "102"
        node_1_interfaces           = ["1/17"]
        node_2_interfaces           = ["1/17"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.portchannel_setting.uuid
    }

    resource "mso_schema_site" "schema_site" {
        schema_id           = mso_schema.schema_blocks.id
        template_name       = mso_schema_template_bd.bd1.template_name
        site_id             = data.mso_site.%[2]s.id
        undeploy_on_destroy = true
    }

    resource "mso_schema_template_deploy_ndo" "fabric_policy_deploy" {
        template_id         = mso_fabric_policies_physical_domain.physical_domain.template_id
        template_type       = "fabric_policy"
        force_apply         = ""
        undeploy_on_destroy = true
        depends_on          = [
            mso_fabric_policies_interface_setting.portchannel_setting,
            mso_fabric_policies_physical_domain.physical_domain_alt,
        ]
    }

    resource "mso_schema_template_deploy_ndo" "fabric_resource_deploy" {
        template_id         = mso_template.fabric_resource_template.id
        template_type       = "fabric_resource"
        force_apply         = ""
        undeploy_on_destroy = true
        depends_on          = [
            mso_schema_template_deploy_ndo.fabric_policy_deploy,
            mso_fabric_resource_policies_port_channel_interface.pc_interface,
            mso_fabric_resource_policies_virtual_port_channel_interface.vpc_interface,
        ]
    }

    resource "mso_schema_template_deploy_ndo" "tenant_template_deploy" {
        template_id         = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.template_id
        template_type       = "tenant"
        force_apply         = ""
        undeploy_on_destroy = true
        depends_on          = [mso_schema_template_deploy_ndo.fabric_resource_deploy]
    }

    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "%[4]s"
        device_mode = "layer3"
        device_type = "firewall"
%[5]s
        depends_on = [
            mso_schema_template_bd.bd2,
            mso_schema_template_external_epg.epg1,
            mso_schema_template_deploy_ndo.tenant_template_deploy,
        ]
    }

    # Layer-1 HA cluster used by the activeActive and activeStandby site
    # tests. Internal/External interfaces are bd1-bound with
    # redirect = true set directly so NDO accepts the matching site-bucket
    # pbr_destinations. device_type = "other" mirrors the JSON sample's
    # deviceType "other" on L1 HA devices.
    resource "mso_service_device_cluster" "cluster_ha" {
        template_id = mso_template.device_template.id
        name        = "%[6]s"
        device_mode = "layer1"
        device_type = "other"
%[7]s
        depends_on = [
            mso_schema_template_bd.bd2,
            mso_schema_template_external_epg.epg1,
            mso_schema_template_deploy_ndo.tenant_template_deploy,
        ]
    }
`, testAccTenantConfig(), msoTemplateSiteName1, msoTemplateTenantName,
		clusterName, interfaceAttributes,
		testServiceDeviceClusterSiteHaClusterName, haInterfaceAttributes,
	)
}

// testAccMSOServiceDeviceClusterSiteConfigOneInterface creates a single-
// interface site with both the cluster (rich: IPSLA + thresholds + PBR
// booleans) and the site-bucket pbr_destination present from the start.
// Aligning the two pieces at create time is required because NDO
// cross-validates redirect (cluster) with PBR destination (site) and
// terraform's resource ordering (cluster PATCH before site PATCH) cannot
// flip either side independently in a later apply.
func testAccMSOServiceDeviceClusterSiteConfigOneInterface(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 210
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{siteClusterRichInterface1}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigUpdateAttributes flips vlan and
// pbr_destination attribute values on the single-interface site AND grows
// the interface's fabric_to_device_connectivity list from 1 to 3 entries
// (port + dpc + vpc). The list grow triggers ForceNew on the outer
// interface block (declared on the resource schema), so this step
// exercises the destroy+recreate path on the site bucket rather than an
// in-place update. The cluster shape does NOT change (same single rich
// interface1 bound to bd1) so terraform produces no cluster diff in this
// apply, which keeps the cluster-redirect + site-PBR cross-check aligned
// at every apply boundary: pbr_destination stays present on interface1
// throughout the recreate. high_availability_mode, promiscuous_mode and
// trunking_port are intentionally NOT set: NDO rejects HA mode on Layer
// 3 devices, and promiscuous_mode / trunking_port are vmmDomain-only and
// rejected on physicalDomain sites.
//
// The dpc entry references the policy-group name of
// mso_fabric_resource_policies_port_channel_interface.pc_interface, and
// the vpc entry references
// mso_fabric_resource_policies_virtual_port_channel_interface.vpc_interface;
// NDO encodes those names verbatim into the pathep-[<name>] suffix of
// the interfaceDn. The vpc entry exercises the two-element node_id list
// form (comma-joined nodeID in JSON, hyphen-joined protpaths URL).
func testAccMSOServiceDeviceClusterSiteConfigUpdateAttributes(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = mso_fabric_resource_policies_port_channel_interface.pc_interface.name
                port_type = "dpc"
            }
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101", "102"]
                path      = mso_fabric_resource_policies_virtual_port_channel_interface.vpc_interface.name
                port_type = "vpc"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:66"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.12"
                weight                 = 7
                tag                    = "20"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{siteClusterRichInterface1}), msoTemplateSiteName1)
}

func testAccMSOServiceDeviceClusterSiteConfigThreeInterfaces(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
        interfaces {
            name = "interface2"
            vlan = 216
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/11"
                port_type = "port"
            }
        }
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/12"
                port_type = "port"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface2"},
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigThreeInterfacesAttrs keeps the
// three-interface shape from the previous step but flips per-interface
// attributes in place (vlan on every interface, pbr_destination
// weight/tag/additional_tracking_ip on interface1, fabric path on
// interface2 and interface3) so the in-place update path is exercised
// on existing slots without changing the slot count. This is the
// site-bucket regression case for the original TypeSet hash slot-swap
// bug, mirroring testAccMSOServiceDeviceClusterConfigUpdateThreeInterfacesAttrs.
func testAccMSOServiceDeviceClusterSiteConfigThreeInterfacesAttrs(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 220
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.12"
                weight                 = 7
                tag                    = "20"
            }
        }
        interfaces {
            name = "interface2"
            vlan = 226
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/21"
                port_type = "port"
            }
        }
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/22"
                port_type = "port"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface2"},
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

func testAccMSOServiceDeviceClusterSiteConfigTwoInterfaces(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/12"
                port_type = "port"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesAttrs keeps the
// two-interface shape from the previous step but flips per-interface
// attributes in place (vlan + pbr_destination mac/weight on interface1,
// vlan + fabric path on interface3), mirroring
// testAccMSOServiceDeviceClusterConfigUpdateTwoInterfacesAttrs.
func testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesAttrs(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface1"
            vlan = 230
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:66"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 3
                tag                    = "10"
            }
        }
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/32"
                port_type = "port"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesSwapped emits the
// same interface set as testAccMSOServiceDeviceClusterSiteConfigTwoInterfaces
// but with the HCL interface blocks in [interface3, interface1] order while
// the cluster's interface_properties stay in [interface1, interface3] order.
// Used by the post-import step to verify that site-to-cluster wiring is
// name-based rather than positional.
func testAccMSOServiceDeviceClusterSiteConfigTwoInterfacesSwapped(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/12"
                port_type = "port"
            }
        }
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigSwitchDomainViaDN swaps the
// site-bucket domain reference from the (type, name) pair used by every
// preceding step to a literal domain_dn pointing at the alternate physical
// domain provisioned by the dependencies stack. Same two-interface set as
// the swap step (rich PBR/IPSLA on interface1 + bare external_epg-bound
// interface3) so this exercises only the domain-attr surface.
func testAccMSOServiceDeviceClusterSiteConfigSwitchDomainViaDN(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_dn   = "uni/phys-${mso_fabric_policies_physical_domain.physical_domain_alt.name}"
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/12"
                port_type = "port"
            }
        }
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigSwitchDomainBackViaTypeName flips
// the site-bucket back to the original (domain_type, domain_name) pair
// pointing at the primary physical domain. Matches the domain attrs of the
// swap step so the only diff from the preceding step is the domain.
func testAccMSOServiceDeviceClusterSiteConfigSwitchDomainBackViaTypeName(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = mso_service_device_cluster.cluster.name
        domain_type = "physicalDomain"
        domain_name = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "interface3"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/12"
                port_type = "port"
            }
        }
        interfaces {
            name = "interface1"
            vlan = 215
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/10"
                port_type = "port"
            }
            pbr_destinations {
                ip                     = "10.10.10.10"
                mac                    = "00:11:22:33:44:55"
                pod_id                 = "1"
                additional_tracking_ip = "10.10.10.11"
                weight                 = 5
                tag                    = "10"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{
		siteClusterRichInterface1,
		{Name: "interface3", WithExternalEPG: true},
	}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigActiveActive switches the
// cluster_site to the layer-1 HA cluster (cluster_ha) and configures
// high_availability_mode = "activeActive". activeActive moves the domain
// down to interface scope (no top-level domain attrs allowed), so each of
// the two interfaces ("Internal" and "External") carries its own domain
// reference. The cluster's Internal/External interface_properties already
// have redirect = true set directly so NDO accepts the per-interface
// pbr_destination. mac/tag values mirror the JSON sample for the
// "ActiveActive" template.
func testAccMSOServiceDeviceClusterSiteConfigActiveActive(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id            = mso_template.device_template.id
        site_id                = data.mso_site.%[2]s.id
        name                   = mso_service_device_cluster.cluster_ha.name
        high_availability_mode = "activeActive"
        interfaces {
            name        = "Internal"
            domain_type = "physicalDomain"
            domain_name = mso_fabric_policies_physical_domain.physical_domain.name
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/1"
                port_type = "port"
                vlan      = 210
                tag       = "aaa"
            }
            pbr_destinations {
                mac = "aa:bb:ee:cc:aa:dd"
                tag = "aaa"
            }
        }
        interfaces {
            name      = "External"
            domain_dn = "uni/phys-${mso_fabric_policies_physical_domain.physical_domain.name}"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/2"
                port_type = "port"
                vlan      = 220
                tag       = "aaa"
            }
            pbr_destinations {
                mac = "cc:aa:cc:aa:cc:aa"
                tag = "aaa"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{siteClusterRichInterface1}), msoTemplateSiteName1)
}

// testAccMSOServiceDeviceClusterSiteConfigActiveStandby switches the
// cluster_site to the layer-1 HA cluster with
// high_availability_mode = "activeStandby". activeStandby keeps the domain
// at device scope (top-level domain_type/domain_name + vlan) and each
// interface carries multiple fabric_to_device_connectivity entries (one per
// tag) plus the matching pbr_destinations. mac/tag/path values mirror the
// JSON sample for the "activeStandby" template.
func testAccMSOServiceDeviceClusterSiteConfigActiveStandby(clusterName string) string {
	return fmt.Sprintf(`%[1]s
    resource "mso_service_device_cluster_site" "cluster_site" {
        template_id            = mso_template.device_template.id
        site_id                = data.mso_site.%[2]s.id
        name                   = mso_service_device_cluster.cluster_ha.name
        high_availability_mode = "activeStandby"
        vlan                   = 230
        domain_type            = "physicalDomain"
        domain_name            = mso_fabric_policies_physical_domain.physical_domain.name
        interfaces {
            name = "Internal"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/1"
                port_type = "port"
                tag       = "www"
            }
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/2"
                port_type = "port"
                tag       = "aaa"
            }
            pbr_destinations {
                mac = "aa:bb:cc:dd:ee:ff"
                tag = "www"
            }
            pbr_destinations {
                mac = "aa:bb:ee:cc:aa:dd"
                tag = "aaa"
            }
        }
        interfaces {
            name = "External"
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/3"
                port_type = "port"
                tag       = "www"
            }
            fabric_to_device_connectivity {
                pod_id    = "1"
                node_id   = ["101"]
                path      = "eth1/4"
                port_type = "port"
                tag       = "aaa"
            }
            pbr_destinations {
                mac = "bb:cc:dd:ee:ff:aa"
                tag = "www"
            }
            pbr_destinations {
                mac = "cc:aa:cc:aa:cc:aa"
                tag = "aaa"
            }
        }
    }
`, testAccMSOServiceDeviceClusterSiteDependencies(clusterName, []siteClusterInterface{siteClusterRichInterface1}), msoTemplateSiteName1)
}

// The ExpectError configurations below use literal placeholder values for
// template_id / site_id and any referenced UUIDs so each error case can fire
// at plan time (SDK validation or CustomizeDiff) without provisioning the
// full deploy chain. None of these configs should reach an API call.

func testAccMSOServiceDeviceClusterSiteConfigErrNoDomain() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrVmmDomainMissingVmmType() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "vmmDomain"
    domain_name = "some_vmm_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        vm_information {
            vm_name   = "vm1"
            vnic_name = "vnic1"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithVmmType() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id     = "00000000-0000-0000-0000-000000000000"
    site_id         = "site-placeholder"
    name            = "err_cluster_site"
    domain_type     = "physicalDomain"
    domain_name     = "some_phys_domain"
    vmm_domain_type = "VMware"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInvalidDomainDn() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_dn   = "uni/bogus-prefix"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInterfaceBothFabricAndVm() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id     = "00000000-0000-0000-0000-000000000000"
    site_id         = "site-placeholder"
    name            = "err_cluster_site"
    domain_type     = "vmmDomain"
    vmm_domain_type = "VMware"
    domain_name     = "some_vmm_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
        vm_information {
            vm_name   = "vm1"
            vnic_name = "vnic1"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInterfaceNoConnectivity() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithVmInformation() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        vm_information {
            vm_name   = "vm1"
            vnic_name = "vnic1"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrPhysicalWithEnhancedLag() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name                = "interface1"
        vlan                = 100
        enhanced_lag_policy = "some-uuid"
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrVmmWithFabricConnectivity() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id     = "00000000-0000-0000-0000-000000000000"
    site_id         = "site-placeholder"
    name            = "err_cluster_site"
    domain_type     = "vmmDomain"
    vmm_domain_type = "VMware"
    domain_name     = "some_vmm_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrVpcWithSingleNode() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "some_pc_policy"
            port_type = "vpc"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrPortWithVpcStyleNode() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101", "102"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInvalidPortType() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "wireless"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrVlanOutOfRange() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 5000
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInvalidPbrIp() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
        pbr_destinations {
            ip  = "not-an-ip"
            mac = "00:11:22:33:44:55"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrInvalidHaMode() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id            = "00000000-0000-0000-0000-000000000000"
    site_id                = "site-placeholder"
    name                   = "err_cluster_site"
    domain_type            = "physicalDomain"
    domain_name            = "some_phys_domain"
    high_availability_mode = "primarySecondary"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

func testAccMSOServiceDeviceClusterSiteConfigErrNameTooLong() string {
	return fmt.Sprintf(`
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "%s"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`, strings.Repeat("a", 65))
}

// Invalid value for the top-level domain_type enum. The SDK schema
// StringInSlice validator rejects this at plan time.
func testAccMSOServiceDeviceClusterSiteConfigErrInvalidDomainType() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "bogusDomain"
    domain_name = "some_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

// Invalid value for vmm_domain_type. Paired with domain_type "vmmDomain"
// so the StringInSlice validator on vmm_domain_type fires (and not the
// physicalDomain conflict).
func testAccMSOServiceDeviceClusterSiteConfigErrInvalidVmmDomainType() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id     = "00000000-0000-0000-0000-000000000000"
    site_id         = "site-placeholder"
    name            = "err_cluster_site"
    domain_type     = "vmmDomain"
    vmm_domain_type = "VirtualBox"
    domain_name     = "some_vmm_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        vm_information {
            vm_name   = "vm1"
            vnic_name = "vnic1"
        }
    }
}
`
}

// pbr_destinations.weight out of the IntBetween(1, 10) range; rejected at
// plan time by the SDK schema validator.
func testAccMSOServiceDeviceClusterSiteConfigErrPbrWeightOutOfRange() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
        pbr_destinations {
            ip     = "10.10.10.10"
            mac    = "00:11:22:33:44:55"
            weight = 99
        }
    }
}
`
}

// fabric_to_device_connectivity.vlan out of the IntBetween(1, 4094) range;
// distinct from the top-level interface vlan check. Rejected at plan time
// by the SDK schema validator.
func testAccMSOServiceDeviceClusterSiteConfigErrFabricPathVlanOutOfRange() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
            vlan      = 5000
        }
    }
}
`
}

// port_type "dpc" with two node_id entries violates the CustomizeDiff rule
// that port and dpc require exactly one node_id.
func testAccMSOServiceDeviceClusterSiteConfigErrDpcWithTwoNodes() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id = "00000000-0000-0000-0000-000000000000"
    site_id     = "site-placeholder"
    name        = "err_cluster_site"
    domain_type = "physicalDomain"
    domain_name = "some_phys_domain"
    interfaces {
        name = "interface1"
        vlan = 100
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101", "102"]
            path      = "test_dpc_for_device"
            port_type = "dpc"
        }
    }
}
`
}

// activeActive moves domain configuration to interface scope; an interface
// without any of domain_type / vmm_domain_type / domain_name / domain_dn
// trips the dedicated CustomizeDiff branch.
func testAccMSOServiceDeviceClusterSiteConfigErrActiveActiveMissingInterfaceDomain() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id            = "00000000-0000-0000-0000-000000000000"
    site_id                = "site-placeholder"
    name                   = "err_cluster_site"
    high_availability_mode = "activeActive"
    interfaces {
        name = "interface1"
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

// activeActive interface-scoped physicalDomain paired with vm_information
// trips the per-interface familyScope = "interface" branch. The wording
// uses "interface" instead of "device", so we assert against that to make
// sure the activeActive branch (not the device-scope branch) fired.
func testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedPhysicalWithVm() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id            = "00000000-0000-0000-0000-000000000000"
    site_id                = "site-placeholder"
    name                   = "err_cluster_site"
    high_availability_mode = "activeActive"
    interfaces {
        name        = "interface1"
        domain_type = "physicalDomain"
        domain_name = "some_phys_domain"
        vm_information {
            vm_name   = "vm1"
            vnic_name = "vnic1"
        }
    }
}
`
}

// activeActive interface-scoped physicalDomain paired with
// enhanced_lag_policy and fabric_to_device_connectivity (to keep the
// "neither fabric nor vm" check from firing first) trips the per-interface
// familyScope branch for enhanced_lag_policy.
func testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedPhysicalWithEnhancedLag() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id            = "00000000-0000-0000-0000-000000000000"
    site_id                = "site-placeholder"
    name                   = "err_cluster_site"
    high_availability_mode = "activeActive"
    interfaces {
        name                = "interface1"
        domain_type         = "physicalDomain"
        domain_name         = "some_phys_domain"
        enhanced_lag_policy = "some-uuid"
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}

// activeActive interface-scoped vmmDomain paired with
// fabric_to_device_connectivity trips the per-interface familyScope branch
// for the vmm-with-fabric check.
func testAccMSOServiceDeviceClusterSiteConfigErrInterfaceScopedVmmWithFabric() string {
	return `
resource "mso_service_device_cluster_site" "cluster_site" {
    template_id            = "00000000-0000-0000-0000-000000000000"
    site_id                = "site-placeholder"
    name                   = "err_cluster_site"
    high_availability_mode = "activeActive"
    interfaces {
        name            = "interface1"
        domain_type     = "vmmDomain"
        vmm_domain_type = "VMware"
        domain_name     = "some_vmm_domain"
        fabric_to_device_connectivity {
            pod_id    = "1"
            node_id   = ["101"]
            path      = "eth1/1"
            port_type = "port"
        }
    }
}
`
}
