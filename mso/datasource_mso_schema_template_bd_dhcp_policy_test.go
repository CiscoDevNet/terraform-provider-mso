package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateBDDHCPPolicy_DataSource(t *testing.T) {
	var pol1 models.TemplateBDDHCPPolicy
	schema := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resourceName := "mso_schema_template_bd_dhcp_policy.test"
	dataSourceName := "data.mso_schema_template_bd_dhcp_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDHCPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithoutRequired(tenantNames[0], schema, name, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithoutRequired(tenantNames[0], schema, name, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithoutRequired(tenantNames[0], schema, name, "bd_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithoutRequired(tenantNames[0], schema, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithInvalidParentResourceName(tenantNames[0], schema, name),
				ExpectError: regexp.MustCompile(`Object Not found`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithRandomAttr(tenantNames[0], schema, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithRequired(tenantNames[0], schema, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resourceName, &pol1),
					resource.TestCheckResourceAttrPair(resourceName, "bd_name", dataSourceName, "bd_name"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "template_name", dataSourceName, "template_name"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_option_name", dataSourceName, "dhcp_option_name"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_option_version", dataSourceName, "dhcp_option_version"),
					resource.TestCheckResourceAttrPair(resourceName, "version", dataSourceName, "version"),
					resource.TestCheckResourceAttrPair(resourceName, "schema_id", dataSourceName, "schema_id"),
				),
			},
		},
	})
}

func CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithRandomAttr(tenant, schema, name, key, value string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintf(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	data "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
		template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
		bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
		name                = mso_schema_template_bd_dhcp_policy.test.name
		%s                  = "%s"
	}
	`, key, value)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithInvalidParentResourceName(tenant, schema, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintln(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	data "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
		template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
		bd_name             = "${mso_schema_template_bd_dhcp_policy.test.bd_name}_invalid"
		name                = mso_schema_template_bd_dhcp_policy.test.name
	}
	`)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithRequired(tenant, schema, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintln(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	data "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
		template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
		bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
		name                = mso_schema_template_bd_dhcp_policy.test.name
	}
	`)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyDataSourceWithoutRequired(tenant, schema, name, attr string) string {
	rBlock := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	rBlock += `
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	`
	switch attr {
	case "schema_id":
		rBlock += `
		data "mso_schema_template_bd_dhcp_policy" "test" {
		#	schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
			template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
			bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
			name                = mso_schema_template_bd_dhcp_policy.test.name
		}
		`
	case "template_name":
		rBlock += `
		data "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
		#	template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
			bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
			name                = mso_schema_template_bd_dhcp_policy.test.name
		}
		`
	case "bd_name":
		rBlock += `
		data "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
			template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
		#	bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
			name                = mso_schema_template_bd_dhcp_policy.test.name
		}
		`
	case "name":
		rBlock += `
		data "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema_template_bd_dhcp_policy.test.schema_id
			template_name       = mso_schema_template_bd_dhcp_policy.test.template_name
			bd_name             = mso_schema_template_bd_dhcp_policy.test.bd_name
		#	name                = mso_schema_template_bd_dhcp_policy.test.name
		}
		`
	}
	return fmt.Sprintln(rBlock)
}
