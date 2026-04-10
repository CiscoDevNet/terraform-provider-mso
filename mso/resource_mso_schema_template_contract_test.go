package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateContractSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateContractSchemaId string

func TestAccMSOSchemaTemplateContractResourceTwoWay(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create contract with name, display_name, filter_type, scope, and filter_relationship")
				},
				Config: testAccMSOSchemaTemplateContractConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "display_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_type", "bothWay"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "scope", "context"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.#", "1"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.0.filter_name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.0.filter_type", "bothWay"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_contract."+msoSchemaTemplateContractName]
						if !ok {
							return fmt.Errorf("Contract resource not found in state")
						}
						msoSchemaTemplateContractSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Expect error when changing filter_type from bothWay to oneWay") },
				Config:      testAccMSOSchemaTemplateContractConfigChangeFilterType(),
				ExpectError: regexp.MustCompile(`The filter_type cannot be changed`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update contract display_name") },
				Config:    testAccMSOSchemaTemplateContractConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "display_name", msoSchemaTemplateContractName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add contract description") },
				Config:    testAccMSOSchemaTemplateContractConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "display_name", msoSchemaTemplateContractName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "description", "Terraform test contract"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove contract description") },
				Config:    testAccMSOSchemaTemplateContractConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "display_name", msoSchemaTemplateContractName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update contract scope") },
				Config:    testAccMSOSchemaTemplateContractConfigUpdateScope(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "scope", "tenant"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update contract priority") },
				Config:    testAccMSOSchemaTemplateContractConfigUpdatePriority(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "scope", "tenant"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "priority", "level1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update contract target_dscp") },
				Config:    testAccMSOSchemaTemplateContractConfigUpdateDscpPriority(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "scope", "tenant"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "target_dscp", "af11"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "priority", "level1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add second filter_relationship with directives and action") },
				Config:    testAccMSOSchemaTemplateContractConfigTwoFilters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.#", "2"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import contract") },
				ResourceName: "mso_schema_template_contract." + msoSchemaTemplateContractName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_contract."+msoSchemaTemplateContractName]
					if !ok {
						return "", fmt.Errorf("Contract resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/contracts/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["contract_name"]), nil
				},
				ImportStateVerify: true,
				// filter_relationships and directives are deprecated computed attributes that get set during import but not in config
				ImportStateVerifyIgnore: []string{"filter_relationships", "directives"},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate contract after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					contractRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/contracts/%s", msoSchemaTemplateName, msoSchemaTemplateContractName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateContractSchemaId), contractRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete contract: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateContractConfigTwoFilters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove second filter_relationship") },
				Config:    testAccMSOSchemaTemplateContractConfigSingleFilter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.#", "1"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.0.filter_name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractName, "filter_relationship.0.filter_type", "bothWay"),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateContractResourceOneWay(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create oneWay contract with provider_to_consumer and consumer_to_provider filter_relationships")
				},
				Config: testAccMSOSchemaTemplateContractOneWayConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "contract_name", msoSchemaTemplateContractOneWayName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "display_name", msoSchemaTemplateContractOneWayName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "filter_type", "oneWay"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "scope", "context"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "filter_relationship.#", "2"),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Expect error when changing filter_type from oneWay to bothWay") },
				Config:      testAccMSOSchemaTemplateContractOneWayConfigChangeFilterType(),
				ExpectError: regexp.MustCompile(`The filter_type cannot be changed`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add bothWay filter alongside directional filters") },
				Config:    testAccMSOSchemaTemplateContractOneWayConfigAddBothWay(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "contract_name", msoSchemaTemplateContractOneWayName),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "filter_type", "oneWay"),
					resource.TestCheckResourceAttr("mso_schema_template_contract."+msoSchemaTemplateContractOneWayName, "filter_relationship.#", "3"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import oneWay contract") },
				ResourceName: "mso_schema_template_contract." + msoSchemaTemplateContractOneWayName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_contract."+msoSchemaTemplateContractOneWayName]
					if !ok {
						return "", fmt.Errorf("Contract resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/contracts/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["contract_name"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filter_relationships", "directives"},
			},
		},
	})
}

// Prerequisite config: site + tenant + schema + filter
func testAccMSOSchemaTemplateContractPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateFilterEntryConfig())
}

// Prerequisite config with two filters
func testAccMSOSchemaTemplateContractPrerequisiteWithTwoFiltersConfig() string {
	return fmt.Sprintf(`%s%s`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), testSchemaTemplateFilterEntryConfig2())
}

func testAccMSOSchemaTemplateContractConfigChangeFilterType() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s"
		filter_type   = "oneWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s"
		filter_type   = "bothWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "context"
		description   = "Terraform test contract"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "context"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigUpdateScope() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "tenant"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigUpdatePriority() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "tenant"
		priority      = "level1"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigUpdateDscpPriority() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "tenant"
		target_dscp   = "af11"
		priority      = "level1"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateContractConfigTwoFilters() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "tenant"
		target_dscp   = "af11"
		priority      = "level1"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[6]s.name
			filter_type = "bothWay"
			directives  = ["log"]
			action      = "deny"
			priority    = "level3"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteWithTwoFiltersConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName, msoSchemaTemplateFilterName2)
}

func testAccMSOSchemaTemplateContractConfigSingleFilter() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s updated"
		filter_type   = "bothWay"
		scope         = "tenant"
		target_dscp   = "af11"
		priority      = "level1"
		description   = ""
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
	}`, testAccMSOSchemaTemplateContractPrerequisiteWithTwoFiltersConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

// oneWay contract configs

func testAccMSOSchemaTemplateContractOneWayPrerequisiteConfig() string {
	return testAccMSOSchemaTemplateContractPrerequisiteWithTwoFiltersConfig()
}

func testAccMSOSchemaTemplateContractOneWayConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s"
		filter_type   = "oneWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "provider_to_consumer"
		}
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[6]s.name
			filter_type = "consumer_to_provider"
		}
	}`, testAccMSOSchemaTemplateContractOneWayPrerequisiteConfig(), msoSchemaTemplateContractOneWayName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName, msoSchemaTemplateFilterName2)
}

func testAccMSOSchemaTemplateContractOneWayConfigChangeFilterType() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s"
		filter_type   = "bothWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "provider_to_consumer"
		}
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[6]s.name
			filter_type = "consumer_to_provider"
		}
	}`, testAccMSOSchemaTemplateContractOneWayPrerequisiteConfig(), msoSchemaTemplateContractOneWayName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName, msoSchemaTemplateFilterName2)
}

func testAccMSOSchemaTemplateContractOneWayConfigAddBothWay() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_contract" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		contract_name = "%[2]s"
		display_name  = "%[2]s"
		filter_type   = "oneWay"
		scope         = "context"
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "bothWay"
		}
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[5]s.name
			filter_type = "provider_to_consumer"
		}
		filter_relationship {
			filter_name = mso_schema_template_filter_entry.%[6]s.name
			filter_type = "consumer_to_provider"
		}
	}`, testAccMSOSchemaTemplateContractOneWayPrerequisiteConfig(), msoSchemaTemplateContractOneWayName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName, msoSchemaTemplateFilterName2)
}

func testAccCheckMSOSchemaTemplateContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_contract" {
			schemaID := rs.Primary.Attributes["schema_id"]
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := con.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := con.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				contractCount, err := tempCont.ArrayCount("contracts")
				if err != nil {
					return fmt.Errorf("Unable to get Contract list")
				}
				for j := 0; j < contractCount; j++ {
					contractCont, err := tempCont.ArrayElement(j, "contracts")
					if err != nil {
						return err
					}
					name := models.StripQuotes(contractCont.S("name").String())
					if rs.Primary.Attributes["contract_name"] == name {
						return fmt.Errorf("Schema Template Contract record still exists")
					}
				}
			}
		}
	}
	return nil
}
