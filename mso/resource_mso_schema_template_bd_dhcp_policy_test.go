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
	schema := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	nameOther:=makeTestVariable(acctest.RandString(5))
	resourceName := "mso_schema_template_bd_dhcp_policy.test"
	resource.Test(t, resource.TestCase{
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
			},{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithOptionalValues(tenantNames[0], schema, name),
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenantNames[0], schema, nameOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDDHCPPolicyExists(resourceName, &pol1),
					resource.TestCheckResourceAttr(resourceName, "bd_name", nameOther),
					resource.TestCheckResourceAttr(resourceName, "name", nameOther),
					resource.TestCheckResourceAttr(resourceName, "template_name", schema),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_name", ""),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_version", "0"),
					resource.TestCheckResourceAttr(resourceName, "version", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
				),
			},
		},
	})
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithOptionalValues(tenant,scheme,name string) string{
	resource := GetParentConfigBDDHCPPolicy(tenant, scheme, name)
	resource += fmt.Sprintln(`
	resource "mso_schema_template_bd_dhcp_policy" "test" {
		schema_id           = mso_schema.test.id
		template_name       = mso_schema.test.template_name
		bd_name             = mso_schema_template_bd.test.name
		name                = mso_dhcp_relay_policy.test.name
	}
	`)
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithRequired(tenant, scheme, name string) string {
	resource := GetParentConfigBDDHCPPolicy(tenant, scheme, name)
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

// func TestAccMSOSchemaTemplateBD_Update(t *testing.T) {
// 	var ss TemplateBD

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckMSOSchemaTemplateBDDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckMSOTemplateBDConfig_basic("flood"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMSOSchemaTemplateBDExists("mso_schema_template_bd.bridge_domain", &ss),
// 					testAccCheckMSOSchemaTemplateBDAttributes("flood", &ss),
// 				),
// 			},
// 			{
// 				Config: testAccCheckMSOTemplateBDConfig_basic("proxy"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMSOSchemaTemplateBDExists("mso_schema_template_bd.bridge_domain", &ss),
// 					testAccCheckMSOSchemaTemplateBDAttributes("proxy", &ss),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccCheckMSOTemplateBDConfig_basic(unicast string) string {
// 	return fmt.Sprintf(`
// 	resource "mso_schema_template_bd" "bridge_domain" {
// 		schema_id = "5ea809672c00003bc40a2799"
// 		template_name = "Template1"
// 		name = "testAccBD"
// 		display_name = "testAcc"
// 		vrf_name = "demo"
// 		layer2_unknown_unicast = "%s"
// 	}
// `, unicast)
// }

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

// func testAccCheckMSOSchemaTemplateBDAttributes(layer2_unknown_unicast string, ss *TemplateBD) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if layer2_unknown_unicast != ss.layer2_unknown_unicast {
// 			return fmt.Errorf("Bad Template BD layer2_unknown_unicast %s", ss.layer2_unknown_unicast)
// 		}

// 		if "testAcc" != ss.display_name {
// 			return fmt.Errorf("Bad Template BD display name %s", ss.display_name)
// 		}

// 		if "demo" != ss.vrf_name {
// 			return fmt.Errorf("Bad Template BD VRF name %s", ss.vrf_name)
// 		}
// 		return nil
// 	}
// }

// type TemplateBD struct {
// 	display_name           string
// 	vrf_name               string
// 	layer2_unknown_unicast string
// }
