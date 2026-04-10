package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateContractDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read contract datasource not found error") },
				Config:      testAccMSOSchemaTemplateContractDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the Contract"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read contract datasource") },
				Config:    testAccMSOSchemaTemplateContractDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_contract.contract", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "display_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "filter_type", "bothWay"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "scope", "context"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "filter_relationship.#", "1"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "filter_relationship.0.filter_name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract.contract", "filter_relationship.0.filter_type", "bothWay"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateContractDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_contract" "contract" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		contract_name = mso_schema_template_contract.%[4]s.contract_name
	}`, testAccMSOSchemaTemplateContractConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateContractName)
}

func testAccMSOSchemaTemplateContractDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_contract" "contract" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		contract_name = "non_existing_contract_name"
	}`, testAccMSOSchemaTemplateContractConfigCreate(), msoSchemaName, msoSchemaTemplateName)
}
