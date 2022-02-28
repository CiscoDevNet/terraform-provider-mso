package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateBDDHCPPolicy_Basic(t *testing.T) {
	var pol1 models.TemplateBDDHCPPolicy
	var pol2 models.TemplateBDDHCPPolicy
	schema := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	nameOther := makeTestVariable(acctest.RandString(5))
	option := makeTestVariable(acctest.RandString(5))
	resourceName := "mso_schema_template_bd_dhcp_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDHCPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenantNames[0], schema, name, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenantNames[0], schema, name, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenantNames[0], schema, name, "bd_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenantNames[0], schema, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenantNames[0], schema, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resourceName, &pol1),
					resource.TestCheckResourceAttr(resourceName, "bd_name", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "template_name", schema),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_name", ""),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_version", "0"),
					resource.TestCheckResourceAttr(resourceName, "version", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
				),
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithOptionalValues(tenantNames[0], schema, name, option),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resourceName, &pol2),
					resource.TestCheckResourceAttr(resourceName, "bd_name", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "template_name", schema),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_name", option),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_version", "1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					testAccCheckMSOSchemaTemplateBDDHCPPolicyIdEqual(&pol1, &pol2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyDestroy(tenantNames[0], schema, name),
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenantNames[0], schema, nameOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resourceName, &pol2),
					resource.TestCheckResourceAttr(resourceName, "bd_name", nameOther),
					resource.TestCheckResourceAttr(resourceName, "name", nameOther),
					resource.TestCheckResourceAttr(resourceName, "template_name", schema),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_name", ""),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_version", "0"),
					resource.TestCheckResourceAttr(resourceName, "version", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					testAccCheckMSOSchemaTemplateBDDHCPPolicyIdNotEqual(&pol1, &pol2),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateBDDHCPPolicy_Negtive(t *testing.T) {
	schema := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDHCPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithInvalidParentName(tenantNames[0], schema, name),
				ExpectError: regexp.MustCompile(`Resource Not Found`),
			},
			{
				Config:      CreateMSOSchemaTemplateBDDHCPPolicyWithRandomAttr(tenantNames[0], schema, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenantNames[0], schema, name),
			},
		},
	})
}

func testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resource string, m *models.TemplateBDDHCPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[resource]

		if !err1 {
			return fmt.Errorf("BD DHCP Policy %s not found", resource)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No BD DHCP Policy id was set")
		}
		BDDHCPPolicyModel := modelFromMSOTemplateBDDHCPPolicyId(rs1.Primary.ID)
		remoteModel, err := getMSOTemplateBDDHCPPolicy(client, BDDHCPPolicyModel)
		if err != nil {
			return err
		}
		*m = *remoteModel
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBDDHCPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_bd_dhcp_policy" {
			BDDHCPPolicyModel := modelFromMSOTemplateBDDHCPPolicyId(rs.Primary.ID)
			_, err := getMSOTemplateBDDHCPPolicy(client, BDDHCPPolicyModel)
			if err == nil {
				return fmt.Errorf("Schema Template BD DHCP Policy with id %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMSOSchemaTemplateBDDHCPPolicyIdEqual(m1, m2 *models.TemplateBDDHCPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := createMSOTemplateBDDHCPPolicyId(m1)
		id2 := createMSOTemplateBDDHCPPolicyId(m2)
		if id1 != id2 {
			return fmt.Errorf("Schema Template BD DHCP Policy ids are not equal")
		}
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBDDHCPPolicyIdNotEqual(m1, m2 *models.TemplateBDDHCPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := createMSOTemplateBDDHCPPolicyId(m1)
		id2 := createMSOTemplateBDDHCPPolicyId(m2)
		if id1 == id2 {
			return fmt.Errorf("Schema Template BD DHCP Policy ids are equal")
		}
		return nil
	}
}

func CreateMSOSchemaTemplateBDDHCPPolicyDestroy(tenant, schema, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithRandomAttr(tenant, schema, name, key, value string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintf(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
		%s                  = "%s"
	}
	`, key, value)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithInvalidParentName(tenant, schema, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintln(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = "${mso_schema_template_bd.test.name}_invalid"
		name                = mso_dhcp_relay_policy.test.name
	}
	`)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithOptionalValues(tenant, scheme, name, option string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, scheme, name)
	resource += fmt.Sprintf(`
	resource "mso_dhcp_option_policy" "test" {
		tenant_id   = data.mso_tenant.test.id
		name        = "%s"
	}
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
		dhcp_option_name    = mso_dhcp_option_policy.test.name
		version             = 1
		dhcp_option_version = 1
	}
	`, option)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenant, schema, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	resource += fmt.Sprintln(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	`)
	return resource
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenant, schema, name, attr string) string {
	rBlock := GetParentConfigBDDHCPPolicy(tenant, schema, name)
	switch attr {
	case "schema_id":
		rBlock += `
		resource "mso_schema_template_bd_dhcp_policy" "test" {
		#	schema_id           = mso_schema.test.id
			template_name       = mso_schema.test.template_name
			bd_name             = mso_schema_template_bd.test.name
			name                = mso_dhcp_relay_policy.test.name
		}
		`
	case "template_name":
		rBlock += `
		resource "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema.test.id
		#	template_name       = mso_schema.test.template_name
			bd_name             = mso_schema_template_bd.test.name
			name                = mso_dhcp_relay_policy.test.name
		}
		`
	case "bd_name":
		rBlock += `
		resource "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema.test.id
			template_name       = mso_schema.test.template_name
		#	bd_name             = mso_schema_template_bd.test.name
			name                = mso_dhcp_relay_policy.test.name
		}
		`
	case "name":
		rBlock += `
		resource "mso_schema_template_bd_dhcp_policy" "test" {
			schema_id           = mso_schema.test.id
			template_name       = mso_schema.test.template_name
			bd_name             = mso_schema_template_bd.test.name
		#	name                = mso_dhcp_relay_policy.test.name
		}
		`
	}
	return fmt.Sprintln(rBlock)
}
