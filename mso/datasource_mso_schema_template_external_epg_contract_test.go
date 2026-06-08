package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateExternalEpgContractDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read External EPG Contract datasource not found error") },
				Config:      testAccMSOSchemaTemplateExtEpgContractDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile("Unable to find the External EPG Contract"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read External EPG Contract datasource") },
				Config:    testAccMSOSchemaTemplateExtEpgContractDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_external_epg_contract.contract", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg_contract.contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg_contract.contract", "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg_contract.contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg_contract.contract", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("data.mso_schema_template_external_epg_contract.contract", "contract_schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg_contract.contract", "contract_template_name", msoSchemaTemplateName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgContractDatasourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_schema_template_external_epg_contract" "contract" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	external_epg_name = "%[4]s"
	contract_name     = mso_schema_template_external_epg_contract.%[5]s_provider.contract_name
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateExtEpgContractConfigProvider(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName, msoSchemaTemplateContractName)
}

func testAccMSOSchemaTemplateExtEpgContractDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%s
data "mso_schema_template_external_epg_contract" "contract" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	external_epg_name = "%[4]s"
	contract_name     = "non_existent_contract"
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateExtEpgContractConfigProvider(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}
