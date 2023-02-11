package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaSiteL3out_Basic(t *testing.T) {
	var l3out1 models.IntersiteL3outs
	var l3out2 models.IntersiteL3outs
	resourceName := "mso_schema_site_l3out.test"
	vrf := makeTestVariable(acctest.RandString(5))
	vrfOther := makeTestVariable(acctest.RandString(5))
	l3out := makeTestVariable(acctest.RandString(5))
	l3outOther := makeTestVariable(acctest.RandString(5))
	prnames := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteL3outWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "l3out_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "vrf_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteL3outWithoutRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out, "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteL3outExists(resourceName, &l3out1),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "l3out_name", l3out),
					resource.TestCheckResourceAttr(resourceName, "template_name", prnames),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrf),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrf, l3outOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteL3outExists(resourceName, &l3out2),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "l3out_name", l3outOther),
					resource.TestCheckResourceAttr(resourceName, "template_name", prnames),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrf),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					testAccCheckMSOSchemaSiteL3outIdNotEqual(&l3out1, &l3out2),
				),
			},
			{
				Config: MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out),
			},
			{
				Config: MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrfOther, l3out),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteL3outExists(resourceName, &l3out2),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "l3out_name", l3out),
					resource.TestCheckResourceAttr(resourceName, "template_name", prnames),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", vrfOther),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					testAccCheckMSOSchemaSiteL3outIdNotEqual(&l3out1, &l3out2),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteL3out_Negative(t *testing.T) {
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
				Config:      MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrf, acctest.RandString(1001)),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSOSchemaSiteL3outAttr(siteNames[0], tenantNames[0], prnames, vrf, l3out, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: MSOSchemaSiteL3outWithRequired(siteNames[0], tenantNames[0], prnames, vrf, l3out),
			},
		},
	})
}

func TestAccMSOSchemaSiteL3out_MultipleCreateDelete(t *testing.T) {
	vrf := makeTestVariable(acctest.RandString(5))
	l3out := makeTestVariable(acctest.RandString(5))
	prnames := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSOSchemaSiteL3outMultiple(siteNames[0], tenantNames[0], prnames, vrf, l3out),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteL3outIdNotEqual(m1, m2 *models.IntersiteL3outs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := L3outModelToL3outId(m1)
		id2 := L3outModelToL3outId(m2)
		if id1 == id2 {
			return fmt.Errorf("Schema Site L3out Ids are equal")
		}
		return nil
	}
}

func testAccCheckMSOSchemaSiteL3outDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_l3out" {
			id := rs.Primary.ID
			l3out, _ := L3outIdToL3outModel(id)
			_, err := client.ReadIntersiteL3outs(l3out)
			if err == nil {
				return fmt.Errorf("Schema Site L3out still exist")
			}
		}
	}
	return nil
}

func testAccCheckMSOSchemaSiteL3outExists(l3outName string, m *models.IntersiteL3outs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, ok := s.RootModule().Resources[l3outName]
		if !ok {
			return fmt.Errorf("Schema Site L3out %s not found", l3outName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schema Site L3out Id was set")
		}
		l3out, err := L3outIdToL3outModel(rs.Primary.ID)
		if err != nil {
			return err
		}
		var read *models.IntersiteL3outs
		read, err = client.ReadIntersiteL3outs(l3out)
		if err != nil {
			return err
		}
		*m = *read
		return nil
	}
}

func MSOSchemaSiteL3outWithoutRequired(site, tenant, name, vrf, l3out, attr string) string {
	rBlock := CreatSchemaSiteConfig(site, tenant, name)
	rBlock += `
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	`
	switch attr {
	case "schema_id":
		rBlock += `
		resource "mso_schema_site_l3out" "test" {
		#	schema_id = mso_schema_site.test.schema_id
			l3out_name = "%s"
			template_name = mso_schema_site.test.template_name
			vrf_name = mso_schema_site_vrf.test.vrf_name
			site_id = mso_schema_site.test.site_id
		}`
	case "l3out_name":
		rBlock += `
		resource "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site.test.schema_id
		#	l3out_name = "%s"
			template_name = mso_schema_site.test.template_name
			vrf_name = mso_schema_site_vrf.test.vrf_name
			site_id = mso_schema_site.test.site_id
		}`
	case "template_name":
		rBlock += `
		resource "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site.test.schema_id
			l3out_name = "%s"
		#	template_name = mso_schema_site.test.template_name
			vrf_name = mso_schema_site_vrf.test.vrf_name
			site_id = mso_schema_site.test.site_id
		}
		`
	case "vrf_name":
		rBlock += `
		resource "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site.test.schema_id
			l3out_name = "%s"
			template_name = mso_schema_site.test.template_name
		#	vrf_name = mso_schema_site_vrf.test.vrf_name
			site_id = mso_schema_site.test.site_id
		}
		`
	case "site_id":
		rBlock += `
		resource "mso_schema_site_l3out" "test" {
			schema_id = mso_schema_site.test.schema_id
			l3out_name = "%s"
			template_name = mso_schema_site.test.template_name
			vrf_name = mso_schema_site_vrf.test.vrf_name
		#	site_id = mso_schema_site.test.site_id
		}
		`
	}
	return fmt.Sprintf(rBlock, vrf, l3out)
}

func MSOSchemaSiteL3outAttr(site, name, user, vrf, l3out, key, val string) string {
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
		%s = "%s"
	}
	`, vrf, l3out, key, val)
	return resource
}

func MSOSchemaSiteL3outWithRequired(site, name, user, vrf, l3out string) string {
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
	`, vrf, l3out)
	return resource
}

func MSOSchemaSiteL3outMultiple(site, tenant, name, vrf, l3out string) string {
	resource := CreatSchemaSiteConfig(site, tenant, name)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_vrf" "test" {
		template_name = mso_schema_site.test.template_name
		site_id       = mso_schema_site.test.site_id
		schema_id     = mso_schema_site.test.schema_id
		vrf_name      = "%s"
	}
	resource "mso_schema_site_l3out" "test" {
		schema_id = mso_schema_site.test.schema_id
		l3out_name = "%s${count.index}"
		template_name = mso_schema_site.test.template_name
		vrf_name = mso_schema_site_vrf.test.vrf_name
		site_id = mso_schema_site.test.site_id
		count = 5
	}
	`, vrf, l3out)
	return resource
}
