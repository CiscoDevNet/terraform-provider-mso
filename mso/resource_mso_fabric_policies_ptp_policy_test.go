package mso

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOPtpPolicyResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile",
						map[string]string{
							"name"                  : "profile1",
							"delay_interval"        : "-2",
							"profile_template"      : "aes67",
							"sync_interval"         : "-3",
							"announce_timeout"      : "3",
							"announce_interval"     : "1",
							"override_node_profile" : "false",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy disabled state") },
				Config:    testAccMSOPtpPolicyConfigUpdateDisable(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", "Terraform test PTP Policy disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_profile_template", "default"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority1", "255"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_priority2", "254"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "global_domain", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_sync_interval", "-1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_delay_interval", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "fabric_announce_timeout", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile",
						map[string]string{
							"name"                  : "profile1",
							"delay_interval"        : "-2",
							"profile_template"      : "aes67",
							"sync_interval"         : "-3",
							"announce_timeout"      : "3",
							"announce_interval"     : "1",
							"override_node_profile" : "false",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile",
						map[string]string{
							"name"                  : "profile2",
							"delay_interval"        : "-2",
							"profile_template"      : "aes67",
							"sync_interval"         : "-3",
							"announce_timeout"      : "3",
							"announce_interval"     : "1",
							"override_node_profile" : "false",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy remove one profile") },
				Config:    testAccMSOPtpPolicyConfigUpdateRemoveProfile(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "description", "Terraform test PTP Policy disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "admin_state", "disabled"),
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
			{
				PreConfig: func() { fmt.Println("Test: Update PTP Policy changing the profile template") },
				Config:    testAccMSOPtpPolicyConfigUpdateChangingTemplate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "name", "tf_test_ptp_policy"),
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
					resource.TestCheckResourceAttr("mso_fabric_policies_ptp_policy.ptp_policy", "ptp_profile.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import PTP Policy") },
				ResourceName:      "mso_fabric_policies_ptp_policy.ptp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_ptp_policy", "fabricPolicyTemplate", "template", "ptpProfile"),
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
		ptp_profile {
			name                  = "profile1"
			profile_template      = "aes67"
			delay_interval        = -2
			sync_interval         = -3
			announce_timeout      = 3
			announce_interval     = 1
			override_node_profile = false
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOPtpPolicyConfigUpdateDisable() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy"
		description                     = "Terraform test PTP Policy disabled"
		admin_state                     = "disabled"
		fabric_profile_template         = "default"
		global_priority1                = 255
		global_priority2                = 254
		global_domain                   = 100
		fabric_announce_interval        = 1
		fabric_sync_interval            = -1
		fabric_delay_interval           = 1
		fabric_announce_timeout         = 3
		ptp_profile {
			name                  = "profile1"
			profile_template      = "aes67"
			delay_interval        = -2
			sync_interval         = -3
			announce_timeout      = 3
			announce_interval     = 1
			override_node_profile = false
		}
		ptp_profile {
			name                  = "profile2"
			profile_template      = "telecom"
			announce_interval     = -3
			delay_interval        = -4
			sync_interval         = -4
			announce_timeout      = 3
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOPtpPolicyConfigUpdateRemoveProfile() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy"
		description                     = "Terraform test PTP Policy disabled"
		admin_state                     = "disabled"
		fabric_profile_template         = "default"
		global_priority1                = 255
		global_priority2                = 254
		global_domain                   = 100
		fabric_announce_interval        = 1
		fabric_sync_interval            = -1
		fabric_delay_interval           = 1
		fabric_announce_timeout         = 3
		ptp_profile {
			name                  = "profile2"
			profile_template      = "telecom"
			announce_interval     = -3
			delay_interval        = -4
			sync_interval         = -4
			announce_timeout      = 3
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOPtpPolicyConfigUpdateChangingTemplate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
		template_id                     = mso_template.template_fabric_policy.id
		name                            = "tf_test_ptp_policy"
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
