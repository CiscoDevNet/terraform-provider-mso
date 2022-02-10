package mso

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPOptionPolicy_DataSource(t *testing.T) {
	var optPolicy1 models.DHCPOptionPolicy
	resourceName := "mso_dhcp_option_policy.test"
	dataSourceName := "data.mso_dhcp_option_policy.test"
	tenant := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	optionName := "acctest" + acctest.RandString(5)
	optionId := strconv.Itoa(acctest.RandIntRange(1, 1000))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyDataSourceWithoutName(tenant, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyDataSourceAttr(tenant, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPOptionPolicyDataSourceWithInvalidName(tenant, name),
				ExpectError: regexp.MustCompile(`DHCP Option Policy with name: (.)+ not found`),
			},
			{
				Config: MSODHCPOptionPolicyDataSourceConfig(tenant, name, optionName, optionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy1),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "tenant_id", dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "option.#", dataSourceName, "option.#"),
					resource.TestCheckResourceAttrPair(resourceName, "option.0.name", dataSourceName, "option.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "option.0.id", dataSourceName, "option.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "option.0.data", dataSourceName, "option.0.data"),
				),
			},
			{
				Config:             MSODHCPOptionPolicyDataSourceWithUpdatedResource(tenant, name, "description", randomValue),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy1),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
				),
			},
		},
	})
}

func MSODHCPOptionPolicyDataSourceConfig(tenant, name, optionname, optionid string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		option {
			name = "%s"
			id = "%s"
		}
	}
	data "mso_dhcp_option_policy" "test" {
		name = mso_dhcp_option_policy.test.name
	}
	`, tenant, tenant, name, optionname, optionid)
	return resource
}

func MSODHCPOptionPolicyDataSourceWithoutName(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
			
	}
	data "mso_dhcp_option_policy" "test" {}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPOptionPolicyDataSourceWithUpdatedResource(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "mso_dhcp_option_policy" "test" {
		name = mso_dhcp_option_policy.test.name
		depends_on = [mso_dhcp_option_policy.test]
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPOptionPolicyDataSourceAttr(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
			
	}
	data "mso_dhcp_option_policy" "test" {
		name = mso_dhcp_option_policy.test.name
		%s = "%s"
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPOptionPolicyDataSourceWithInvalidName(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
	}
	data "mso_dhcp_option_policy" "test" {
		name = "${mso_dhcp_option_policy.test.name}_invalid"
	}
	`, tenant, tenant, name)
	return resource
}
