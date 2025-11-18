package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesBGPPeerPrefixPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: BGP Peer Prefix Policy Data Source") },
				Config:    testAccMSOTenantPoliciesBGPPeerPrefixPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "name", "test_bgp_peer_prefix_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "description", "Test BGP Peer Prefix Policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "action", "restart"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "max_number_of_prefixes", "1000"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "threshold_percentage", "50"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy", "restart_time", "60"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesBGPPeerPrefixPolicyDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
        template_id = mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy.template_id
        name        = "test_bgp_peer_prefix_policy"
    }`, testAccMSOTenantPoliciesBGPPeerPrefixPolicyConfigCreate())
}
