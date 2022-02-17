package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteL3out_DataSource(t *testing.T) {
	var l3outModel models.IntersiteL3outs
	resourceName := "mso_schema_site_l3out.test"
	dataSourceName := "mso_schema_site_l3out.test"
	vrf := makeTestVariable(acctest.RandString(5))
	l3out := makeTestVariable(acctest.RandString(5))
	prnames := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteL3outDataSourceWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "l3out_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "vrf_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceAttr(siteNames[0], tenantNames[0], prnames, vrf, l3out, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSOSchemaSiteL3outDataSourceInvalidName(siteNames[0], tenantNames[0], prnames, vrf, l3out),
				ExpectError: regexp.MustCompile(`unable to find siteL3out`),
			},
			{
				Config: MSOSchemaSiteL3outDataSourceWithRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &l3outModel),
					resource.TestCheckResourceAttrPair(resourceName, "schema_id", dataSourceName, "schema_id"),
					resource.TestCheckResourceAttrPair(resourceName, "l3out_name", dataSourceName, "l3out_name"),
					resource.TestCheckResourceAttrPair(resourceName, "template_name", dataSourceName, "template_name"),
					resource.TestCheckResourceAttrPair(resourceName, "vrf_name", dataSourceName, "vrf_name"),
					resource.TestCheckResourceAttrPair(resourceName, "site_id", dataSourceName, "site_id"),
				),
			},
		},
	})
}

func MSOSchemaSiteL3outDataSourceWithoutRequired(site, tenant, name, vrf, l3out, attr string) string {
	rBlock := CreatSchemaSiteConfig(site, tenant, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	resource "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site.test.schema_id
		l3out_name = "%s"
		template_name = mso_schema_site.test.template_name
		vrf_name = mso_schema_site_vrf.test.vrf_name
		site_id = mso_schema_site.test.site_id
	}
	`
	switch attr {
	case "schema_id":
		rBlock += `
		data "mso_schema_site_l3out" "test" {
		#	schema_id = mso_schema_site_l3out.test.schema_id
			l3out_name = mso_schema_site_l3out.test.l3out_name
			template_name = mso_schema_site_l3out.test.template_name
			vrf_name = mso_schema_site_l3out.test.vrf_name
			site_id = mso_schema_site_l3out.test.site_id
		}`
	case "l3out_name":
		rBlock += `
		data "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site_l3out.test.schema_id
		#	l3out_name = mso_schema_site_l3out.test.l3out_name
			template_name = mso_schema_site_l3out.test.template_name
			vrf_name = mso_schema_site_l3out.test.vrf_name
			site_id = mso_schema_site_l3out.test.site_id
		}`
	case "template_name":
		rBlock += `
		data "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site_l3out.test.schema_id
			l3out_name = mso_schema_site_l3out.test.l3out_name
		#	template_name = mso_schema_site_l3out.test.template_name
			vrf_name = mso_schema_site_l3out.test.vrf_name
			site_id = mso_schema_site_l3out.test.site_id
		}
		`
	case "vrf_name":
		rBlock += `
		data "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site_l3out.test.schema_id
			l3out_name = mso_schema_site_l3out.test.l3out_name
			template_name = mso_schema_site_l3out.test.template_name
		#	vrf_name = mso_schema_site_l3out.test.vrf_name
			site_id = mso_schema_site_l3out.test.site_id
		}
		`
	case "site_id":
		rBlock += `
		data "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site_l3out.test.schema_id
			l3out_name = mso_schema_site_l3out.test.l3out_name
			template_name = mso_schema_site_l3out.test.template_name
			vrf_name = mso_schema_site_l3out.test.vrf_name
		#	site_id = mso_schema_site_l3out.test.site_id
		}
		`
	}
	return fmt.Sprintf(rBlock, vrf, l3out)
}

func MSOSchemaSiteL3outDataSourceInvalidName(site, name, user, vrf, l3out string) string {
	resource := CreatSchemaSiteConfig(site, name, user)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	resource "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site.test.schema_id
		l3out_name = "%s"
		template_name = mso_schema_site.test.template_name
		vrf_name = mso_schema_site_vrf.test.vrf_name
		site_id = mso_schema_site.test.site_id
	}
	data "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site_l3out.test.schema_id
		l3out_name = "${mso_schema_site_l3out.test.l3out_name}_invalid"
		template_name = mso_schema_site_l3out.test.template_name
		vrf_name = mso_schema_site_l3out.test.vrf_name
		site_id = mso_schema_site_l3out.test.site_id
	}
	`, vrf, l3out)
	return resource
}

func MSOSchemaSiteL3outDataSourceAttr(site, name, user, vrf, l3out, key, val string) string {
	resource := CreatSchemaSiteConfig(site, name, user)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	resource "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site.test.schema_id
		l3out_name = "%s"
		template_name = mso_schema_site.test.template_name
		vrf_name = mso_schema_site_vrf.test.vrf_name
		site_id = mso_schema_site.test.site_id
	}
	data "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site_l3out.test.schema_id
		l3out_name = mso_schema_site_l3out.test.l3out_name
		template_name = mso_schema_site_l3out.test.template_name
		vrf_name = mso_schema_site_l3out.test.vrf_name
		site_id = mso_schema_site_l3out.test.site_id
		%s = "%s"
	}
	`, vrf, l3out, key, val)
	return resource
}

func MSOSchemaSiteL3outDataSourceWithRequired(site, name, user, vrf, l3out string) string {
	resource := CreatSchemaSiteConfig(site, name, user)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	resource "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site.test.schema_id
		l3out_name = "%s"
		template_name = mso_schema_site.test.template_name
		vrf_name = mso_schema_site_vrf.test.vrf_name
		site_id = mso_schema_site.test.site_id
	}
	data "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site_l3out.test.schema_id
		l3out_name = mso_schema_site_l3out.test.l3out_name
		template_name = mso_schema_site_l3out.test.template_name
		vrf_name = mso_schema_site_l3out.test.vrf_name
		site_id = mso_schema_site_l3out.test.site_id
	}
	`, vrf, l3out)
	return resource
}
