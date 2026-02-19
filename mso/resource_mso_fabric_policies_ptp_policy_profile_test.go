package mso

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccMSOPtpPolicyProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create PTP Policy Profile") },
				Config:    testAccMSOPtpPolicyProfileConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "name", "tf_ptp_profile"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "description", "Terraform test PTP Policy Profile"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "profile_template", "aes67"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "delay_interval", "-2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "sync_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_timeout", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "override_node_profile", "false"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "ptp_policy_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy Profile") },
				Config:    testAccMSOPtpPolicyProfileConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "name", "tf_ptp_profile"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "description", "Terraform test PTP Policy Profile updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "profile_template", "smpte"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "delay_interval", "-2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "sync_interval", "-2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_timeout", "10"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "override_node_profile", "true"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "ptp_policy_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Telecom PTP Policy Profile") },
				Config:    testAccMSOPtpPolicyProfileConfigCreateTelecom(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "name", "tf_ptp_profile_2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "description", "Terraform test PTP Policy Profile 2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "profile_template", "telecom"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "delay_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "sync_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "announce_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "announce_timeout", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "local_priority", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "destination_mac_type", "forwardable"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "mismatched_mac_handling", "reply_with_config_mac"),
					resource.TestCheckNoResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "override_node_profile"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "ptp_policy_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Telecom PTP Policy Profile") },
				Config:    testAccMSOPtpPolicyProfileConfigUpdateTelecom(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "name", "tf_ptp_profile_2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "description", "Terraform test PTP Policy Profile 2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "profile_template", "telecom"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "delay_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "sync_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "announce_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "announce_timeout", "2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "local_priority", "99"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "destination_mac_type", "non_forwardable"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "mismatched_mac_handling", "reply_with_received_mac"),
					resource.TestCheckNoResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "override_node_profile"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile_2", "ptp_policy_uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import PTP Policy Profile") },
				ResourceName:      "mso_fabric_policies_ptp_policy_profile.ptp_policy_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_ptp_policy_profile", "fabricPolicyTemplate", "template", "ptpPolicy", "profiles"),
	})
}

func testAccMSOPtpPolicyProfileConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
		template_id           = mso_template.template_fabric_policy.id
		ptp_policy_uuid       = mso_fabric_policies_ptp_policy.ptp_policy.uuid
		name                  = "tf_ptp_profile"
		description           = "Terraform test PTP Policy Profile"
		profile_template      = "aes67"
		delay_interval        = -2
		sync_interval         = -3
		announce_timeout      = 3
		announce_interval     = 1
		override_node_profile = false

		# Explicit dependency on PTP Policy
		depends_on = [
			mso_fabric_policies_ptp_policy.ptp_policy
		]
	}`, testAccMSOPtpPolicyConfigCreate())
}

func testAccMSOPtpPolicyProfileConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
		template_id           = mso_template.template_fabric_policy.id
		ptp_policy_uuid       = mso_fabric_policies_ptp_policy.ptp_policy.uuid
		name                  = "tf_ptp_profile"
		description           = "Terraform test PTP Policy Profile updated"
		profile_template      = "smpte"
		delay_interval        = -2
		sync_interval         = -2
		announce_timeout      = 10
		announce_interval     = -3
		override_node_profile = true

		# Explicit dependency on PTP Policy
		depends_on = [
			mso_fabric_policies_ptp_policy.ptp_policy
		]
	}`, testAccMSOPtpPolicyConfigCreate())
}

func testAccMSOPtpPolicyProfileConfigCreateTelecom() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile_2" {
		template_id                = mso_template.template_fabric_policy.id
		ptp_policy_uuid            = mso_fabric_policies_ptp_policy.ptp_policy.uuid
		name                       = "tf_ptp_profile_2"
		description                = "Terraform test PTP Policy Profile 2"
		profile_template           = "telecom"
		announce_interval          = -3
		delay_interval             = -4
		sync_interval              = -4
		announce_timeout           = 3
		local_priority             = 120
		destination_mac_type       = "forwardable"
		mismatched_mac_handling    = "reply_with_config_mac"
	}`, testAccMSOPtpPolicyProfileConfigUpdate())
}

func testAccMSOPtpPolicyProfileConfigUpdateTelecom() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile_2" {
		template_id                = mso_template.template_fabric_policy.id
		ptp_policy_uuid            = mso_fabric_policies_ptp_policy.ptp_policy.uuid
		name                       = "tf_ptp_profile_2"
		description                = "Terraform test PTP Policy Profile 2"
		profile_template           = "telecom"
		announce_interval          = -3
		delay_interval             = -4
		sync_interval              = -4
		announce_timeout           = 2
		local_priority             = 99
		destination_mac_type       = "non_forwardable"
		mismatched_mac_handling    = "reply_with_received_mac"

		# Explicit dependency on first PTP Profile
		# This is to ensure the test deletes the profiles in order.
		depends_on = [
			mso_fabric_policies_ptp_policy_profile.ptp_policy_profile
		]
	}`, testAccMSOPtpPolicyProfileConfigUpdate())
}
