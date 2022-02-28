package mso

import (
	"fmt"
	"regexp"

	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPOptionPolicyOptionDataSource_Basic(t *testing.T) {
	var optionPolicyOption models.DHCPOptionPolicyOption
	resourceName := "mso_dhcp_option_policy_option.test"
	dataSourceName := "data.mso_dhcp_option_policy_option.test"
	tenant := tenantNames[0]
	policyName := makeTestVariable(acctest.RandString(5))
	name := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyOptionDataSourceWithoutRequiredArguments(tenant, policyName, name, "option_policy_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionDataSourceWithoutRequiredArguments(tenant, policyName, name, "option_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionDataSourceAttr(tenant, policyName, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionDataSourceWithInvalidName(tenant, policyName, name),
				ExpectError: regexp.MustCompile(`No DHCP Option Policy found`),
			},
			{
				Config: MSODHCPOptionPolicyOptionDataSourceConfig(tenant, policyName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &optionPolicyOption),
					resource.TestCheckResourceAttrPair(resourceName, "option_policy_name", dataSourceName, "option_policy_name"),
					resource.TestCheckResourceAttrPair(resourceName, "option_name", dataSourceName, "option_name"),
					resource.TestCheckResourceAttrPair(resourceName, "option_id", dataSourceName, "option_id"),
					resource.TestCheckResourceAttrPair(resourceName, "option_data", dataSourceName, "option_data"),
				),
			},
			{
				Config: MSODHCPOptionPolicyOptionDataSourceWithUpdatedResource(tenant, policyName, name, "option_data", randomValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &optionPolicyOption),
					resource.TestCheckResourceAttrPair(resourceName, "option_data", dataSourceName, "option_data"),
				),
			},
		},
	})
}

func MSODHCPOptionPolicyOptionDataSourceConfig(tenant, policyName, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
	}
	data "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
		option_name = mso_dhcp_option_policy_option.test.option_name
	}
	`, tenant, tenant, policyName, name)
	return resource
}

func MSODHCPOptionPolicyOptionDataSourceWithoutRequiredArguments(tenant, policyName, name, attr string) string {
	resource := `
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
	}
	`
	switch attr {
	case "option_policy_name":
		resource += `
		data "mso_dhcp_option_policy_option" "test"{
		#	option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
			option_name = mso_dhcp_option_policy_option.test.option_name
		}
		`
	case "option_name":
		resource += `
		data "mso_dhcp_option_policy_option" "test"{
			option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
		#	option_name = mso_dhcp_option_policy_option.test.option_name
		}
	`
	}
	return fmt.Sprintf(resource, tenant, tenant, policyName, name)
}

func MSODHCPOptionPolicyOptionDataSourceWithUpdatedResource(tenant, policyName, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
		%s = "%s"
	}
	data "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
		option_name = mso_dhcp_option_policy_option.test.option_name
	}
	`, tenant, tenant, policyName, name, key, val)
	return resource
}

func MSODHCPOptionPolicyOptionDataSourceAttr(tenant, policyName, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
	}

	data "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
		option_name = mso_dhcp_option_policy_option.test.option_name
		%s = "%s"
	}
	`, tenant, tenant, policyName, name, key, val)
	return resource
}

func MSODHCPOptionPolicyOptionDataSourceWithInvalidName(tenant, policyName, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
	}
	data "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy_option.test.option_policy_name
		option_name = "${mso_dhcp_option_policy_option.test.option_name}_invalid"
	}
	`, tenant, tenant, policyName, name)
	return resource
}
