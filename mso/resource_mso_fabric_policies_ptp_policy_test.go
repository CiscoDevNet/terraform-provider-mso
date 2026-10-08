package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOPtpPolicyResource(t *testing.T) {
	resourceName := "mso_fabric_policies_ptp_policy.ptp_policy"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create PTP Policy") },
				Config:    testAccMSOPtpPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", "Terraform test PTP Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_profile_template", "default"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority1", "255"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority2", "254"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_domain", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_sync_interval", "-1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_delay_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_timeout", "3"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy.ptp_policy", "uuid"),
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
					fmt.Println("Test: Recreate PTP Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "ptpPolicy"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOPtpPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", "Terraform test PTP Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_profile_template", "default"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority1", "255"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority2", "254"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_domain", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_sync_interval", "-1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_delay_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_timeout", "3"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy.ptp_policy", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy disabled state") },
				Config:    testAccMSOPtpPolicyConfigUpdateDisable(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy_new"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_profile_template", "default"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority1", "255"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority2", "254"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_domain", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_sync_interval", "-1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_delay_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_timeout", "3"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy.ptp_policy", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy changing the profile template") },
				Config:    testAccMSOPtpPolicyConfigUpdateChangingTemplate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy_new"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", "Terraform test PTP Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_profile_template", "smpte"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority1", "200"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority2", "250"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_domain", "99"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_sync_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_delay_interval", "-2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_timeout", "10"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy.ptp_policy", "uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import PTP Policy") },
				ResourceName:      "mso_fabric_policies_ptp_policy.ptp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing PTP Policy") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_ptp_policy", "fabricPolicyTemplate", "template", "ptpPolicy"),
	})
}

func testAccMSOPtpPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy"
		description                     = "Terraform test PTP Policy"
		admin_state                     = "enabled"
		fabric_profile_template         = "default"
		global_priority1                = 255
		global_priority2                = 254
		global_domain                   = 100
		fabric_announce_interval        = 1
		fabric_sync_interval            = -1
		fabric_delay_interval           = 1
		fabric_announce_timeout         = 3
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOPtpPolicyConfigUpdateDisable() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy_new"
		description                     = ""
		admin_state                     = "disabled"
		fabric_profile_template         = "default"
		global_priority1                = 255
		global_priority2                = 254
		global_domain                   = 100
		fabric_announce_interval        = 1
		fabric_sync_interval            = -1
		fabric_delay_interval           = 1
		fabric_announce_timeout         = 3
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOPtpPolicyConfigUpdateChangingTemplate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy_new"
		description                     = "Terraform test PTP Policy"
		admin_state                     = "enabled"
		fabric_profile_template         = "smpte"
		global_priority1                = 200
		global_priority2                = 250
		global_domain                   = 99
		fabric_announce_interval        = -3
		fabric_sync_interval            = -4
		fabric_delay_interval           = -2
		fabric_announce_timeout         = 10
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
