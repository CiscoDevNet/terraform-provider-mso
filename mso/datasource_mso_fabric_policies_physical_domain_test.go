package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOPhysicalDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Physical Domain Data Source") },
				Config:    testAccMSOPhysicalDomainDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_physical_domain.physical_domain", "name", "tf_test_physical_domain"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_physical_domain.physical_domain", "description", "Terraform test Physical Domain"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_physical_domain.physical_domain", "vlan_pool_uuid"),
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
