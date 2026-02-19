package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOPtpPolicyProfileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: PTP Policy Profile Data Source") },
				Config:    testAccMSOPtpPolicyProfileDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "name", "tf_ptp_profile"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "description", "Terraform test PTP Policy Profile"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "profile_template", "telecom"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "delay_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "sync_interval", "-4"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_timeout", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "local_priority", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "destination_mac_type", "forwardable"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "mismatched_mac_handling", "reply_with_config_mac"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "override_node_profile", "false"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "ptp_policy_uuid"),
				),
			},
		},
	})
}

func testAccMSOPtpPolicyProfileConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
		template_id                = mso_template.template_fabric_policy.id
		ptp_policy_uuid            = mso_fabric_policies_ptp_policy.ptp_policy.uuid
		name                       = "tf_ptp_profile"
		description                = "Terraform test PTP Policy Profile"
		profile_template           = "telecom"
		announce_interval          = -3
		delay_interval             = -4
		sync_interval              = -4
		announce_timeout           = 3
		local_priority             = 120
		destination_mac_type       = "forwardable"
		mismatched_mac_handling    = "reply_with_config_mac"
		override_node_profile      = "false"
	}`, testAccMSOPtpPolicyConfigCreate())
}

func testAccMSOPtpPolicyProfileDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
	    template_id        = mso_fabric_policies_ptp_policy_profile.ptp_policy_profile.template_id
	    name               = "tf_ptp_profile"
    }`, testAccMSOPtpPolicyProfileConfig())
}
