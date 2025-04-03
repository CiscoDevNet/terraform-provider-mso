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
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "description", ""),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "template_id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "providers.#", "2"),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "providers",
						map[string]string{
							"application_epg_uuid":       fmt.Sprintf("mso_schema_template_anp_epg.%s.uuid", msoSchemaTemplateAnpEpgName),
							"dhcp_server_address":        "1.1.1.1",
							"dhcp_server_vrf_preference": "false",
							"external_epg_uuid":          "",
						},
					),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "providers",
						map[string]string{
							"application_epg_uuid":       "",
							"dhcp_server_address":        "2.2.2.2",
							"dhcp_server_vrf_preference": "true",
							"external_epg_uuid":          fmt.Sprintf("mso_schema_template_external_epg.%s.uuid", msoSchemaTemplateExtEpgName),
						},
					),
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
