package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateAnpEpgContractDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgContractResourceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Read EPG Contract datasource") },
				Config:    testAccMSOSchemaTemplateAnpEpgContractDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp_epg_contract.contract", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp_epg_contract.contract", "contract_schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_contract.contract", "contract_template_name", msoSchemaTemplateName),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Read EPG Contract datasource not found error") },
				Config:      testAccMSOSchemaTemplateAnpEpgContractDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile("Unable to find the ANP EPG Contract"),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgContractDatasourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_schema_template_anp_epg_contract" "contract" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	anp_name      = "%[4]s"
	epg_name      = "%[5]s"
	contract_name = mso_schema_template_anp_epg_contract.%[6]s_provider.contract_name
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateAnpEpgContractConfigProvider(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateContractName)
}

func testAccMSOSchemaTemplateAnpEpgContractDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%s
data "mso_schema_template_anp_epg_contract" "contract" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	anp_name      = "%[4]s"
	epg_name      = "%[5]s"
	contract_name = "non_existent_contract"
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateAnpEpgContractConfigProvider(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}
