package mso

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const msoTemplateSiteName1 = "ansible_test"
const msoTemplateSiteName2 = "ansible_test_2"

var msoTenantName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateVrfName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateExtEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

func testSiteConfigAnsibleTest() string {
	return fmt.Sprintf(`
data "mso_site" "%[1]s" {
	name = "%[1]s"
}
`, msoTemplateSiteName1)
}

func testSiteConfigAnsibleTest2() string {
	return fmt.Sprintf(`
data "mso_site" "%[1]s" {
	name = "%[1]s"
}
`, msoTemplateSiteName2)
}

func testTenantConfig() string {
	return fmt.Sprintf(`
resource "mso_tenant" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	site_associations { 
		site_id = data.mso_site.%[2]s.id
	}
}
`, msoTenantName, msoTemplateSiteName1)
}

func testTenantPolicyTemplateConfig() string {
	return fmt.Sprintf(`
resource "mso_template" "%[1]s" {
	template_name = "%[1]s"
	template_type = "tenant"
	tenant_id     = mso_tenant.%[2]s.id
}
`, msoTenantPolicyTemplateName, msoTenantName)
}

func testSchemaConfig() string {
	return fmt.Sprintf(`
resource "mso_schema" "%[1]s" {
	name = "%[1]s"
	template {
		name         = "%[2]s"
		display_name = "%[2]s"
		tenant_id    = mso_tenant.%[3]s.id
	}
}
`, msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testSchemaTemplateAnpConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	schema_id    = mso_schema.%[2]s.id
	template     = "%[3]s"
}
`, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateAnpEpgConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp_epg" "%[1]s" {
	name          = "%[1]s"
	display_name  = "%[1]s"
	anp_name      = "%[2]s"
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	depends_on = [
		mso_schema_template_anp.%[2]s,
	]
}
`, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateVrfConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_vrf" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	schema_id    = mso_schema.%[2]s.id
	template     = "%[3]s"
}
`, msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateExtEpgConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_external_epg" "%[1]s" {
	external_epg_name = "%[1]s"
	display_name      = "%[1]s"
	vrf_name          = mso_schema_template_vrf.%[2]s.name
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
}
`, msoSchemaTemplateExtEpgName, msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}
