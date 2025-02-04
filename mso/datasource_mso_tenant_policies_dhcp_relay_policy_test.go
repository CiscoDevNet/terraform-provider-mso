package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesDHCPRelayPolicyDataSource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: DHCP Relay Policy Data Source") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					customTestCheckResourceAttr(fmt.Sprintf("data.mso_tenant_policies_dhcp_relay_policy.%s", name), map[string]string{
						"name":                                   name,
						"description":                            "",
						"template_id":                            "reference",
						"providers.#":                            "2",
						"providers.0.application_epg":            "reference",
						"providers.0.dhcp_server_address":        "1.1.1.1",
						"providers.0.dhcp_server_vrf_preference": "false",
						"providers.0.external_epg":               "",
						"providers.1.application_epg":            "",
						"providers.1.dhcp_server_address":        "2.2.2.2",
						"providers.1.dhcp_server_vrf_preference": "true",
						"providers.1.external_epg":               "reference",
					}),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesDHCPRelayPolicyDataSource(name string) string {
	return fmt.Sprintf(`%[1]s
    data "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
        template_id = mso_tenant_policies_dhcp_relay_policy.%[2]s.template_id
        name        = "%[2]s"
    }`, testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreate(name), name)
}
