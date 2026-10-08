package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesBGPPeerPrefixPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create BGP Peer Prefix Policy") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Test BGP Peer Prefix Policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "restart"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "1000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "60"),
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
					fmt.Println("Test: Recreate BGP Peer Prefix Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "bgpPeerPrefixPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Test BGP Peer Prefix Policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "restart"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "1000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "60"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BGP Peer Prefix Policy with Log Action") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithLogAction(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Updated with Log Action"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "log"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "5000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "80"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "100"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BGP Peer Prefix Policy with Reject Action") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithRejectAction(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Updated with Reject Action"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "reject"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "10000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "75"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BGP Peer Prefix Policy with Shutdown Action") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithShutdownAction(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy_shutdown"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Updated with Shutdown Action"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "shutdown"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "25000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "90"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "300"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BGP Peer Prefix Policy with Maximum Values") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy_max"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Maximum Values Test"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "restart"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "300000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "65535"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BGP Peer Prefix Policy with Minimum Values") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy_min"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Minimum Values Test"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "log"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "1"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import BGP Peer Prefix Policy") },
				ResourceName:      "mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing BGP Peer Prefix Policy") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_bgp_peer_prefix_policy", "bgpPeerPrefixPol"),
	})
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy"
        description             = "Test BGP Peer Prefix Policy"
        action                  = "restart"
        max_number_of_prefixes  = 1000
        threshold_percentage    = 50
        restart_time            = 60
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithLogAction() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy"
        description             = "Updated with Log Action"
        action                  = "log"
        max_number_of_prefixes  = 5000
        threshold_percentage    = 80
        restart_time            = 100
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithRejectAction() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy"
        description             = "Updated with Reject Action"
        action                  = "reject"
        max_number_of_prefixes  = 10000
        threshold_percentage    = 75
        restart_time            = 1
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithShutdownAction() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy_shutdown"
        description             = "Updated with Shutdown Action"
        action                  = "shutdown"
        max_number_of_prefixes  = 25000
        threshold_percentage    = 90
        restart_time            = 300
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy_max"
        description             = "Maximum Values Test"
        action                  = "restart"
        max_number_of_prefixes  = 300000
        threshold_percentage    = 100
        restart_time            = 65535
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigUpdateWithMinValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id             = mso_template.template_tenant.id
        name                    = "test_bgp_peer_prefix_policy_min"
        description             = "Minimum Values Test"
        action                  = "log"
        max_number_of_prefixes  = 1
        threshold_percentage    = 1
        restart_time            = 1
    }`, testAccMSOTemplateResourceTenantConfig())
}
