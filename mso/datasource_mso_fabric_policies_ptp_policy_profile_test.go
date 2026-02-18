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
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "profile_template", "aes67"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "delay_interval", "-2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "sync_interval", "-3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "announce_timeout", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "override_node_profile", "false"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_ptp_policy_profile.ptp_policy_profile", "uuid"),
				),
			},
		},
	})
}

func testAccMSOPtpPolicyProfileDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
	    template_id        = mso_fabric_policies_ptp_policy_profile.ptp_policy_profile.template_id
	    name               = "tf_ptp_profile"
    }`, testAccMSOPtpPolicyProfileConfigCreate())
}
