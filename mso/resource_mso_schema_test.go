package mso

// Note: Only the default aci_multi_site template type is tested.
// Other template types (aci_autonomous, ndfc, cloud_local, sr_mpls) are excluded.
//
// Note: Deprecated template_name/tenant_id attributes are not tested.
// Only the modern template block is used.

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaId is used to capture the schema ID from the first test step for use in the manual delete/recreate step.
var msoSchemaId string

func TestAccMSOSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Schema with one template") },
				Config:    testAccMSOSchemaConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName,
						"template_type": "aci_multi_site",
						"description":   "",
					}),
					// Capture the schema ID for the manual delete/recreate step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema."+msoSchemaName]
						if !ok {
							return fmt.Errorf("mso_schema.%s not found in state", msoSchemaName)
						}
						msoSchemaId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema basic fields") },
				Config:    testAccMSOSchemaConfigUpdateBasicFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add Schema description") },
				Config:    testAccMSOSchemaConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add template description") },
				Config:    testAccMSOSchemaConfigAddTemplateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "Terraform test template",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove template description") },
				Config:    testAccMSOSchemaConfigRemoveTemplateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove schema description") },
				Config:    testAccMSOSchemaConfigUpdateBasicFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Schema") },
				ResourceName:      "mso_schema." + msoSchemaName,
				ImportState:       true,
				ImportStateVerify: true,
				// Deprecated attributes are set to empty on import and are not used in the config.
				ImportStateVerifyIgnore: []string{"template_name", "tenant_id"},
			},
			{
				PreConfig: func() { fmt.Println("Test: Add template children (VRF, BD, ANP, EPG)") },
				Config:    testAccMSOSchemaConfigWithChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update schema description with children present") },
				Config:    testAccMSOSchemaConfigWithChildrenUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove template children") },
				Config:    testAccMSOSchemaConfigRemoveChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add extra template") },
				Config:    testAccMSOSchemaConfigAddExtraTemplate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName2,
						"display_name":  msoSchemaTemplateName2,
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove extra template") },
				Config:    testAccMSOSchemaConfigRemoveTemplate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName + " updated",
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove all templates") },
				Config:    testAccMSOSchemaConfigRemoveAllTemplates(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", "Terraform test schema"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove description") },
				Config:    testAccMSOSchemaConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName+" updated"),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Schema after template changes") },
				ResourceName:      "mso_schema." + msoSchemaName,
				ImportState:       true,
				ImportStateVerify: true,
				// Deprecated attributes are set to empty on import and are not used in the config.
				ImportStateVerifyIgnore: []string{"template_name", "tenant_id"},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Verify recreation of manually deleted Schema")
					client := testAccProvider.Meta().(*client.Client)
					err := client.DeletebyId("api/v1/schemas/" + msoSchemaId)
					if err != nil {
						panic(fmt.Sprintf("Failed to delete schema %s via API: %s", msoSchemaId, err))
					}
				},
				Config: testAccMSOSchemaConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "name", msoSchemaName),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema."+msoSchemaName, "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema."+msoSchemaName, "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName,
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
		},
	})
}

func testAccMSOSchemaPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s`, testSiteConfigAnsibleTest(), testTenantConfig())
}

func testAccMSOSchemaConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name = "%[2]s"
		template {
			name         = "%[3]s"
			display_name = "%[3]s"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigUpdateBasicFields() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name = "%[2]s updated"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigAddTemplateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			description  = "Terraform test template"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigRemoveTemplateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
			description  = ""
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigAddExtraTemplate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
		template {
			name         = "%[5]s"
			display_name = "%[5]s"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName, msoSchemaTemplateName2)
}

func testAccMSOSchemaConfigRemoveTemplate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testAccMSOSchemaConfigRemoveAllTemplates() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema"
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName)
}

func testAccMSOSchemaConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name = "%[2]s updated"
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName)
}

func testAccMSOSchemaConfigWithChildren() string {
	return testAccMSOSchemaConfigAddDescription() + testSchemaTemplateVrfConfig() + testSchemaTemplateBdConfig() + testSchemaTemplateAnpConfigWithBdDependency() + testSchemaTemplateAnpEpgConfig()
}

func testAccMSOSchemaConfigWithChildrenUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema updated"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName) + testSchemaTemplateVrfConfig() + testSchemaTemplateBdConfig() + testSchemaTemplateAnpConfigWithBdDependency() + testSchemaTemplateAnpEpgConfig()
}

// testAccMSOSchemaConfigRemoveChildren returns the same schema config as testAccMSOSchemaConfigWithChildrenUpdateDescription
// but without children. Keeping the schema unchanged avoids a concurrent PATCH to NDO during child deletion.
func testAccMSOSchemaConfigRemoveChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema" "%[2]s" {
		name        = "%[2]s updated"
		description = "Terraform test schema updated"
		template {
			name         = "%[3]s"
			display_name = "%[3]s updated"
			tenant_id    = mso_tenant.%[4]s.id
		}
	}`, testAccMSOSchemaPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

// testSchemaTemplateAnpConfigWithBdDependency creates ANP config with an explicit dependency on BD.
// This serializes all child PATCH operations to the schema (EPG → ANP → BD → VRF) to prevent
// concurrent PATCH requests to the same NDO schema endpoint causing "Err Missing Ref" errors.
func testSchemaTemplateAnpConfigWithBdDependency() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	schema_id    = mso_schema.%[2]s.id
	template     = "%[3]s"
	depends_on   = [mso_schema_template_bd.%[4]s]
}
`, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName)
}

// testAccCheckMSOSchemaDestroy verifies that the schema is deleted from MSO.
func testAccCheckMSOSchemaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema" {
			_, err := client.GetViaURL("api/v1/schemas/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Schema still exists")
			}
		} else {
			continue
		}

	}
	return nil
}
