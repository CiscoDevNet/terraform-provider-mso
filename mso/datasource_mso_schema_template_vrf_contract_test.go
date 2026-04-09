package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateVrfContractDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateVrfContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read VRF contract datasource not found error") },
				Config:      testAccMSOSchemaTemplateVrfContractDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the VRF Contract"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read VRF contract datasource") },
				Config:    testAccMSOSchemaTemplateVrfContractDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_vrf_contract.vrf_contract", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf_contract.vrf_contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf_contract.vrf_contract", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf_contract.vrf_contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf_contract.vrf_contract", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("data.mso_schema_template_vrf_contract.vrf_contract", "contract_schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf_contract.vrf_contract", "contract_template_name", msoSchemaTemplateName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateVrfContractDatasourceConfig() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf_contract" "vrf_contract" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = mso_schema_template_contract.%[5]s.contract_name
		relationship_type = "provider"
	}
	resource "mso_schema_template_vrf_contract" "vrf_contract_consumer" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = mso_schema_template_contract.%[5]s.contract_name
		relationship_type = "consumer"
	}`, testAccMSOSchemaTemplateVrfContractPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName, msoSchemaTemplateContractName)
}

func testAccMSOSchemaTemplateVrfContractDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_vrf_contract" "vrf_contract" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = mso_schema_template_vrf_contract.vrf_contract.contract_name
		relationship_type = "provider"
	}`, testAccMSOSchemaTemplateVrfContractDatasourceConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateVrfContractDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_vrf_contract" "vrf_contract" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = "non_existing_contract"
		relationship_type = "provider"
	}`, testAccMSOSchemaTemplateVrfContractDatasourceConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}
