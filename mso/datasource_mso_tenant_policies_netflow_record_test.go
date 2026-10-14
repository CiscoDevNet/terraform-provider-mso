package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Missing tenant policy data source returns an error") },
				Config:      testAccMSOTenantPoliciesNetflowRecordDataSourceMissing(),
				ExpectError: regexp.MustCompile(`Policy name missing_out_of_band_deletion not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Record Data Source") },
				Config:    testAccMSOTenantPoliciesNetflowRecordDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_record.netflow_record", "description", "Updated NetFlow Record"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_netflow_record.netflow_record", "uuid"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_record.netflow_record", "match_parameters.0", "ethertype"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_record.netflow_record", "match_parameters.1", "destination_mac"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesNetflowRecordDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_tenant_policies_netflow_record.netflow_record.template_id
        name        = "test_netflow_record"
    }`, testAccMSOTenantPoliciesNetflowRecordConfigUpdateMatchParams())
}

func testAccMSOTenantPoliciesNetflowRecordDataSourceMissing() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_netflow_record" "missing" {
        template_id = mso_tenant_policies_netflow_record.netflow_record.template_id
        name        = "missing_out_of_band_deletion"
    }`, testAccMSOTenantPoliciesNetflowRecordConfigUpdateMatchParams())
}
