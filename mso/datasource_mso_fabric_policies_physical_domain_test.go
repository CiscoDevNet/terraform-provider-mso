package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOPhysicalDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject missing policy data source") },
				Config:      testAccMSOPhysicalDomainDataSourceMissing(),
				ExpectError: regexp.MustCompile("Policy name tf_missing_oob not found"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Physical Domain Data Source") },
				Config:    testAccMSOPhysicalDomainDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_physical_domain.physical_domain", "name", "tf_test_physical_domain"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_physical_domain.physical_domain", "description", "Terraform test Physical Domain"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_physical_domain.physical_domain", "vlan_pool_uuid"),
				),
			},
		},
	})
}

func testAccMSOPhysicalDomainDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_physical_domain" "physical_domain" {
	    template_id        = mso_fabric_policies_physical_domain.physical_domain.template_id
	    name               = "tf_test_physical_domain"
    }`, testAccMSOPhysicalDomainConfigCreate())
}

func testAccMSOPhysicalDomainDataSourceMissing() string {
	return fmt.Sprintf(`%s
    data "mso_fabric_policies_physical_domain" "missing" {
        template_id = mso_fabric_policies_physical_domain.physical_domain.template_id
        name        = "tf_missing_oob"
    }`, testAccMSOPhysicalDomainConfigCreate())
}
