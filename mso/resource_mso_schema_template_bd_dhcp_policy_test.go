package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateBDDHCPPolicy_Basic(t *testing.T) {
	var pol1 models.TemplateBDDHCPPolicy
	schema:=makeTestVariable(acctest.RandString(5))
	name:=makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDHCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenantNames[0],schema,name),
			},
		},
	})
}

func CreateMSOSchemaTemplateBDDHCPPolicyWithoutRequired(tenant,schema,name string) string{
	resource:=GetParentConfigBDDHCPPolicy(tenant,schema,name)
	return resource
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
			if err != nil {
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
