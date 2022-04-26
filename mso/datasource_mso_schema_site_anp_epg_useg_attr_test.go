package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteAnpEpgUsegAttr_DataSource(t *testing.T) {
	var useg models.SiteUsegAttr
	resourceName := "mso_schema_site_anp_epg_useg_attr.test"
	datasourceName := "data.mso_schema_site_anp_epg_useg_attr.test"
	template_name := "Template1"
	usegName := makeTestVariable(acctest.RandString(5))
	anpName := makeTestVariable(acctest.RandString(5))
	epgName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	usegTypeTag := "tag"
	value := "test_tag_value"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "anp_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "epg_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "useg_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrDataSourceWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "startsWith", "Sample Test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg),
					resource.TestCheckResourceAttrPair(resourceName, "schema_id", datasourceName, "schema_id"),
					resource.TestCheckResourceAttrPair(resourceName, "site_id", datasourceName, "site_id"),
					resource.TestCheckResourceAttrPair(resourceName, "template_name", datasourceName, "template_name"),
					resource.TestCheckResourceAttrPair(resourceName, "anp_name", datasourceName, "anp_name"),
					resource.TestCheckResourceAttrPair(resourceName, "epg_name", datasourceName, "epg_name"),
					resource.TestCheckResourceAttrPair(resourceName, "useg_name", datasourceName, "useg_name"),
					resource.TestCheckResourceAttrPair(resourceName, "useg_type", datasourceName, "useg_type"),
					resource.TestCheckResourceAttrPair(resourceName, "value", datasourceName, "value"),
					resource.TestCheckResourceAttrPair(resourceName, "description", datasourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "category", datasourceName, "category"),
					resource.TestCheckResourceAttrPair(resourceName, "operator", datasourceName, "operator"),
					resource.TestCheckResourceAttrPair(resourceName, "fv_subnet", datasourceName, "fv_subnet"),
				),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrsDataSourceForRandomAttrName(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrsDataSourceUpdatedAttrName(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "description", datasourceName, "description"),
				),
			},
		},
	})
}

func MSOSchemaSiteAnpEpgUsegAttrDataSourceWithoutRequired(site, tenant, template, anp, epg, useg, useg_type, value, keyAttr string) string {
	rBlock := MSOSchemaSiteAnpEpgUsegAttrWithRequired(site, tenant, template, anp, epg, useg, useg_type, value, "startsWith", "Sample Test")

	switch keyAttr {
	case "schema_id":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
		#	schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	case "site_id":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
		#	site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	case "anp_name":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
		#	anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	case "epg_name":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
		#	epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	case "useg_name":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site_anp_epg_useg_attr.test.template_name
		#	useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	case "template_name":
		rBlock += `
		data "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
		#	template_name = mso_schema_site.test.template_name
			useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
		}`
	}
	resource := fmt.Sprintf(rBlock)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrsDataSourceForRandomAttrName(site, tenant, template, anp, epg, useg, useg_type, key, value string) string {
	resource := MSOSchemaSiteAnpEpgUsegAttrWithRequired(site, tenant, template, anp, epg, useg, useg_type, value, "startsWith", "Sample Test")

	resource += fmt.Sprintf(`
	data "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		%s            = "%s"
	}
	`, useg, key, value)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrsDataSourceUpdatedAttrName(site, tenant, template, anp, epg, useg, useg_type, useg_value, key, value string) string {
	resource := MSOSchemaSiteAnpEpgUsegAttrsForRandomAttrName(site, tenant, template, anp, epg, useg, useg_type, useg_value, key, value)

	resource += fmt.Sprintf(`
	data "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
	}
	`)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrDataSourceWithRequired(site, tenant, template, anp, epg, useg, useg_type, value string, opAttr ...string) string {
	// here opAttr param will be used for the extra attributes of other types
	// format of opAttr: type opAttr []string{operator, category, fvSubnet}
	// for 'ip' you have to pass first two opAttr as empty string("")
	// ----------------------------------------------------------------------
	resource := MSOSchemaSiteAnpEpgUsegAttrWithRequired(site, tenant, template, anp, epg, useg, useg_type, value, opAttr[0], opAttr[1])

	resource += fmt.Sprintf(`
	data "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = mso_schema_site_anp_epg_useg_attr.test.useg_name
	}`)
	return resource

}
