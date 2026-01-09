package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPOptionPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: DHCP Option Policy Data Source") },
				Config:    testAccMSODHCPOptionPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "description", "Test DHCP Option Policy"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "template_id"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.#", "2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.0.name", "option_1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.0.id", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.0.data", "data_1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.1.name", "option_2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.1.id", "2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options.1.data", "data_2"),
				),
			},
		},
	})
}

func testAccMSODHCPOptionPolicyDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_tenant_policies_dhcp_option_policy.dhcp_policy.template_id
        name        = "test_dhcp_option_policy"
    }`, testAccMSODHCPOptionPolicyConfigCreate())
}
