package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Captured from the create step's Check so the drift-recovery PreConfig can
// delete the cluster directly via the NDO API without re-parsing state.
var (
	testServiceDeviceClusterTemplateID  string
	testServiceDeviceClusterCurrentName string
)

const (
	testServiceDeviceClusterName        = "test_device_cluster"
	testServiceDeviceClusterNameRenamed = "test_device_cluster_renamed"
)

func TestAccMSOServiceDeviceClusterResource(t *testing.T) {
	resourceName := "mso_service_device_cluster.cluster"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Device Cluster with one interface") },
				Config:    testAccMSOServiceDeviceClusterConfigCreateOneInterface(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterName),
					resource.TestCheckResourceAttr(resourceName, "device_mode", "layer3"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "firewall"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "template_id"),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "1"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":         "sourceIP",
						"min_threshold":                "10",
						"max_threshold":                "90",
						"threshold_down_action":        "permit",
						"external_epg_uuid":            "mso_schema_template_external_epg.epg1.uuid",
						"ipsla_monitoring_policy_uuid": "mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid",
					}),
					// Capture template_id and current name for the drift-recovery PreConfig step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testServiceDeviceClusterTemplateID = rs.Primary.Attributes["template_id"]
						testServiceDeviceClusterCurrentName = rs.Primary.Attributes["name"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster attributes on one interface") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated device cluster description"),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "1"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":  "destinationIP",
						"min_threshold":         "20",
						"max_threshold":         "80",
						"threshold_down_action": "deny",
						"preferred_group":       "true",
						"rewrite_source_mac":    "true",
						"config_static_mac":     "true",
						"is_backup_redirect_ip": "true",
						"pod_aware_redirection": "true",
						"resilient_hashing":     "true",
						"tag_based_sorting":     "true",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster to three interfaces") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateThreeInterfaces(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterName),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "3"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"threshold_down_action": "permit",
						"external_epg_uuid":     "mso_schema_template_external_epg.epg1.uuid",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface2",
					}, map[string]string{
						"load_balance_hashing":         "destinationIP",
						"bd_uuid":                      "mso_schema_template_bd.bd1.uuid",
						"ipsla_monitoring_policy_uuid": "mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"bd_uuid": "mso_schema_template_bd.bd2.uuid",
						"anycast": "true",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Service Device Cluster three-interface attribute changes")
				},
				Config: testAccMSOServiceDeviceClusterConfigUpdateThreeInterfacesAttrs(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "3"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":  "destinationIP",
						"min_threshold":         "25",
						"max_threshold":         "75",
						"threshold_down_action": "deny",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface2",
					}, map[string]string{
						"load_balance_hashing": "sourceDestinationAndProtocol",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"bd_uuid": "mso_schema_template_bd.bd2.uuid",
						"anycast": "false",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster to two interfaces") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateTwoInterfaces(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterName),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"threshold_down_action": "permit",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, nil),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Service Device Cluster two-interface attribute changes")
				},
				Config: testAccMSOServiceDeviceClusterConfigUpdateTwoInterfacesAttrs(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":  "destinationIP",
						"min_threshold":         "35",
						"max_threshold":         "65",
						"threshold_down_action": "deny",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"bd_uuid": "mso_schema_template_bd.bd2.uuid",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster reset thresholds to 0 and set QoS policy") },
				Config:    testAccMSOServiceDeviceClusterConfigUpdateThresholdsToZeroAndSetQoS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"min_threshold":                "0",
						"max_threshold":                "0",
						"qos_policy_uuid":              "mso_tenant_policies_custom_qos_policy.qos1.uuid",
						"ipsla_monitoring_policy_uuid": "mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, nil),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster clear IPSLA monitoring policy") },
				Config:    testAccMSOServiceDeviceClusterConfigClearIpslaPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"qos_policy_uuid": "mso_tenant_policies_custom_qos_policy.qos1.uuid",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, nil),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Device Cluster clear QoS policy") },
				Config:    testAccMSOServiceDeviceClusterConfigClearQosPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing": "sourceIP",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(resourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, nil),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Rename Service Device Cluster in-place") },
				Config:    testAccMSOServiceDeviceClusterConfigRenameTwoInterfaces(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterNameRenamed),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
					// Capture the renamed name so the drift PreConfig can locate the device.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testServiceDeviceClusterCurrentName = rs.Primary.Attributes["name"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Service Device Cluster after manual deletion from NDO")
					if err := manuallyDeleteServiceDeviceClusterFromTemplate(testServiceDeviceClusterTemplateID, testServiceDeviceClusterCurrentName); err != nil {
						panic(err)
					}
				},
				Config: testAccMSOServiceDeviceClusterConfigRenameTwoInterfaces(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testServiceDeviceClusterNameRenamed),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "interface_properties.#", "2"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Device Cluster") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccMSOServiceDeviceClusterResourceErrors exercises every input-validation
// path on the resource. The step order is deliberate: SDK schema validation
// failures (ValidateFunc / StringLenBetween / StringInSlice / IntBetween) come
// first, and CustomizeDiff failures come last. The acceptance-test destroy
// step revalidates the last step's configuration, so that configuration must
// pass SDK schema validation; otherwise the destroy walk reports an invalid
// configuration.
// Both CustomizeDiff steps below have HCL that validates cleanly, so they are
// safe terminators.
func TestAccMSOServiceDeviceClusterResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: invalid device_mode is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrInvalidDeviceMode(),
				ExpectError: regexp.MustCompile(`expected device_mode to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid device_type is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrInvalidDeviceType(),
				ExpectError: regexp.MustCompile(`expected device_type to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid load_balance_hashing is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrInvalidLoadBalanceHashing(),
				ExpectError: regexp.MustCompile(`expected interface_properties\.\d+\.load_balance_hashing to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: invalid threshold_down_action is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrInvalidThresholdDownAction(),
				ExpectError: regexp.MustCompile(`expected interface_properties\.\d+\.threshold_down_action to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: min_threshold above 100 is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrMinThresholdTooHigh(),
				ExpectError: regexp.MustCompile(`expected interface_properties\.\d+\.min_threshold to be in the range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: max_threshold above 100 is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrMaxThresholdTooHigh(),
				ExpectError: regexp.MustCompile(`expected interface_properties\.\d+\.max_threshold to be in the range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: cluster name longer than 64 characters is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrNameTooLong(),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: interface_properties with both bd_uuid and external_epg_uuid is rejected") },
				Config:      testAccMSOServiceDeviceClusterConfigErrBothInterfaceTargets(),
				ExpectError: regexp.MustCompile(`exactly one of bd_uuid or external_epg_uuid must be set`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: interface_properties with neither bd_uuid nor external_epg_uuid is rejected")
				},
				Config:      testAccMSOServiceDeviceClusterConfigErrNeitherInterfaceTarget(),
				ExpectError: regexp.MustCompile(`exactly one of bd_uuid or external_epg_uuid must be set`),
			},
		},
	})
}

// manuallyDeleteServiceDeviceClusterFromTemplate removes a device from the
// service device template via a JSON-patch remove on the index resolved from
// the live template. Used by drift-recovery test steps to simulate an
// out-of-band deletion on NDO before the same Terraform configuration is
// re-applied.
func manuallyDeleteServiceDeviceClusterFromTemplate(templateID, deviceName string) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return fmt.Errorf("manual delete: get template %s: %w", templateID, err)
	}
	deviceIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "name", deviceName, "deviceTemplate", "template", "devices")
	if err != nil {
		return fmt.Errorf("manual delete: locate device %q: %w", deviceName, err)
	}
	removePayload := models.GetRemovePatchPayload(fmt.Sprintf("/deviceTemplate/template/devices/%d", deviceIndex))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), removePayload); err != nil {
		return fmt.Errorf("manual delete: patch remove %q: %w", deviceName, err)
	}
	return nil
}

// testAccMSOServiceDeviceClusterDependencies provides the shared prerequisites
// for both the cluster and site Service Device Cluster acceptance tests:
// tenant, fabric-policy template, tenant template, schema with VRF/BDs/extEPG,
// an IPSLA monitoring policy and the service-device template. The four
// templates plus the schema are chained via depends_on so that creation order
// is tenant -> fabric_policy -> tenant_template -> schema -> device_template
// and destruction occurs in the reverse order. The service-device template is
// associated with the ansible_test site so the same configuration can be
// reused as the base for the site-bucket tests, which require the site to
// already be present on the device template.
func testAccMSOServiceDeviceClusterDependencies() string {
	return fmt.Sprintf(`%s
    resource "mso_template" "fabric_policy_template" {
      template_name = "test_fabric_policy_for_device"
      template_type = "fabric_policy"
      sites         = [data.mso_site.%[3]s.id]
      depends_on    = [mso_tenant.%[2]s]
    }

    resource "mso_template" "tenant_template" {
      template_name = "test_tenant_template_for_device"
      template_type = "tenant"
      tenant_id     = mso_tenant.%[2]s.id
      sites         = [data.mso_site.%[3]s.id]
      depends_on    = [mso_template.fabric_policy_template]
    }

    resource "mso_schema" "schema_blocks" {
        name = "demo_schema_blocks"
        template {
            name          = "Template1"
            display_name  = "TEMP1"
            tenant_id     = mso_tenant.%[2]s.id
            template_type = "aci_multi_site"
        }
        depends_on = [mso_template.tenant_template]
    }

    resource "mso_template" "device_template" {
      template_name = "test_device_template"
      template_type = "service_device"
      tenant_id     = mso_tenant.%[2]s.id
      sites         = [data.mso_site.%[3]s.id] 
	  depends_on    = [mso_schema.schema_blocks]
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

    resource "mso_schema_template_external_epg" "epg1" {
        schema_id           = mso_schema.schema_blocks.id
        template_name       = mso_schema_template_bd.bd2.template_name
        external_epg_name   = "test_epg_1"
        vrf_name            = mso_schema_template_vrf.vrf.name
        display_name        = "template_epg"
        external_epg_type   = "on-premise"
        l3out_name          = mso_schema_template_l3out.l3out1.l3out_name
        l3out_schema_id     = mso_schema_template_l3out.l3out1.schema_id
        l3out_template_name = mso_schema_template_l3out.l3out1.template_name
    }

    resource "mso_schema_template_l3out" "l3out1" {
        schema_id     = mso_schema.schema_blocks.id
        template_name = mso_schema_template_bd.bd2.template_name
        l3out_name    = "test_l3out_1"
        display_name  = "template_l3out"
        vrf_name      = mso_schema_template_vrf.vrf.name
    }

    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla1" {
        template_id = mso_template.tenant_template.id
        name        = "test_ipsla_for_device"
        sla_type    = "icmp"
    }

    resource "mso_tenant_policies_custom_qos_policy" "qos1" {
        template_id = mso_template.tenant_template.id
        name        = "test_qos_for_device"
    }
`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
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

// testAccMSOServiceDeviceClusterConfigUpdateThreeInterfacesAttrs keeps the
// three-interface shape from the previous step but flips per-interface
// attributes (load_balance_hashing on interface1 and interface2, the
// threshold trio on interface1, anycast on interface3) so the in-place
// update path is exercised on existing slots without changing the slot
// count. This is the regression case for the original TypeSet hash
// slot-swap bug, where attribute changes on already-present interfaces
// were silently masked in the diff.
func testAccMSOServiceDeviceClusterConfigUpdateThreeInterfacesAttrs() string {
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
            load_balance_hashing         = "destinationIP"
            min_threshold                = 25
            max_threshold                = 75
            threshold_down_action        = "deny"
        }
        interface_properties {
            name                         = "interface2"
            bd_uuid                      = mso_schema_template_bd.bd1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "sourceDestinationAndProtocol"
        }
        interface_properties {
            name                         = "interface3"
            bd_uuid                      = mso_schema_template_bd.bd2.uuid
            anycast                      = false
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
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

// testAccMSOServiceDeviceClusterConfigUpdateTwoInterfacesAttrs keeps the
// two-interface shape from the previous step but flips load_balance_hashing
// and the threshold trio on interface1 so the in-place update path is
// exercised on an existing slot without changing the slot count.
func testAccMSOServiceDeviceClusterConfigUpdateTwoInterfacesAttrs() string {
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
            load_balance_hashing         = "destinationIP"
            min_threshold                = 35
            max_threshold                = 65
            threshold_down_action        = "deny"
        }
        interface_properties {
            name                         = "interface3"
            bd_uuid                      = mso_schema_template_bd.bd2.uuid
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

func testAccMSOServiceDeviceClusterConfigUpdateThresholdsToZeroAndSetQoS() string {
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
            qos_policy_uuid              = mso_tenant_policies_custom_qos_policy.qos1.uuid
            load_balance_hashing         = "sourceIP"
            min_threshold                = 0
            max_threshold                = 0
        }
        interface_properties {
            name                         = "interface3"
            bd_uuid                      = mso_schema_template_bd.bd2.uuid
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

func testAccMSOServiceDeviceClusterConfigClearIpslaPolicy() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                    = "interface1"
            external_epg_uuid       = mso_schema_template_external_epg.epg1.uuid
            qos_policy_uuid         = mso_tenant_policies_custom_qos_policy.qos1.uuid
            load_balance_hashing    = "sourceIP"
        }
        interface_properties {
            name                    = "interface3"
            bd_uuid                 = mso_schema_template_bd.bd2.uuid
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

func testAccMSOServiceDeviceClusterConfigClearQosPolicy() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                    = "interface1"
            external_epg_uuid       = mso_schema_template_external_epg.epg1.uuid
            load_balance_hashing    = "sourceIP"
        }
        interface_properties {
            name                    = "interface3"
            bd_uuid                 = mso_schema_template_bd.bd2.uuid
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

// testAccMSOServiceDeviceClusterConfigUpdateAttributes flips every optional
// attribute on the single-interface variant so the in-place update path is
// exercised across description, all advanced booleans, the threshold trio and
// load_balance_hashing.
func testAccMSOServiceDeviceClusterConfigUpdateAttributes() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "test_device_cluster"
        description = "updated device cluster description"
        device_mode = "layer3"
        device_type = "firewall"
        interface_properties {
            name                         = "interface1"
            external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
            ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
            load_balance_hashing         = "destinationIP"
            min_threshold                = 20
            max_threshold                = 80
            threshold_down_action        = "deny"
            preferred_group              = true
            rewrite_source_mac           = true
            config_static_mac            = true
            is_backup_redirect_ip        = true
            pod_aware_redirection        = true
            resilient_hashing            = true
            tag_based_sorting            = true
        }
    }`, testAccMSOServiceDeviceClusterDependencies())
}

// testAccMSOServiceDeviceClusterConfigRenameTwoInterfaces keeps the two-interface
// shape but renames the cluster, exercising the in-place name change path
// (cluster.name is not ForceNew) and is also reused by the drift-recovery step.
func testAccMSOServiceDeviceClusterConfigRenameTwoInterfaces() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster" "cluster" {
        template_id = mso_template.device_template.id
        name        = "%s"
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
        }
    }`, testAccMSOServiceDeviceClusterDependencies(), testServiceDeviceClusterNameRenamed)
}

// The ExpectError configurations below use literal placeholder values for
// template_id and any referenced UUIDs so the test does not need to provision
// real infrastructure. SDK validation and CustomizeDiff both fire before any
// API request is made, so these checks are effectively offline.

func testAccMSOServiceDeviceClusterConfigErrBothInterfaceTargets() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name              = "interface1"
        bd_uuid           = "00000000-0000-0000-0000-000000000001"
        external_epg_uuid = "00000000-0000-0000-0000-000000000002"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrNeitherInterfaceTarget() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name = "interface1"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrInvalidDeviceMode() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer42"
    device_type = "firewall"
    interface_properties {
        name    = "interface1"
        bd_uuid = "00000000-0000-0000-0000-000000000001"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrInvalidDeviceType() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "router"
    interface_properties {
        name    = "interface1"
        bd_uuid = "00000000-0000-0000-0000-000000000001"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrInvalidLoadBalanceHashing() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name                 = "interface1"
        bd_uuid              = "00000000-0000-0000-0000-000000000001"
        load_balance_hashing = "roundRobin"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrInvalidThresholdDownAction() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name                  = "interface1"
        bd_uuid               = "00000000-0000-0000-0000-000000000001"
        threshold_down_action = "drop"
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrMinThresholdTooHigh() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name          = "interface1"
        bd_uuid       = "00000000-0000-0000-0000-000000000001"
        min_threshold = 150
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrMaxThresholdTooHigh() string {
	return `
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "err_cluster"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name          = "interface1"
        bd_uuid       = "00000000-0000-0000-0000-000000000001"
        max_threshold = 200
    }
}
`
}

func testAccMSOServiceDeviceClusterConfigErrNameTooLong() string {
	return fmt.Sprintf(`
resource "mso_service_device_cluster" "cluster" {
    template_id = "00000000-0000-0000-0000-000000000000"
    name        = "%s"
    device_mode = "layer3"
    device_type = "firewall"
    interface_properties {
        name    = "interface1"
        bd_uuid = "00000000-0000-0000-0000-000000000001"
    }
}
`, strings.Repeat("a", 65))
}
