package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOL3DomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Missing fabric_policies_l3_domain data source remains an error") },
				Config:      testAccMSOL3DomainDataSourceMissing(),
				ExpectError: regexp.MustCompile(`Policy name missing_out_of_band_test_policy not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: L3 Domain Data Source") },
				Config:    testAccMSOL3DomainDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_l3_domain.l3_domain", "description", "Test L3 Domain"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_l3_domain.l3_domain", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_l3_domain.l3_domain", "template_id"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid"),
					resource.TestCheckResourceAttrPair(
						"data.mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid",
						"mso_fabric_policies_vlan_pool.vlan_pool", "uuid",
					),
				),
			},
		},
	})
}

func testAccMSOL3DomainDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id = mso_template.template_fabric_policy.id
        name        = mso_fabric_policies_l3_domain.l3_domain.name
    }`, testAccMSOL3DomainConfigCreate())
}

func testAccMSOL3DomainDataSourceMissing() string {
	return fmt.Sprintf(`%s
    data "mso_fabric_policies_l3_domain" "missing" {
        template_id = mso_fabric_policies_l3_domain.l3_domain.template_id
        name        = "missing_out_of_band_test_policy"
    }`, testAccMSOL3DomainConfigCreate())
}
