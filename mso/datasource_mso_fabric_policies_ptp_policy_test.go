package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOPtpPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: PTP Policy Data Source") },
				Config:    testAccMSOPtpPolicyDataSource(),
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
				),
			},
		},
	})
}

func testAccMSOPtpPolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_ptp_policy" "ptp_policy" {
	    template_id        = mso_fabric_policies_ptp_policy.ptp_policy.template_id
	    name               = "tf_test_ptp_policy"
    }`, testAccMSOPtpPolicyConfigCreate())
}
