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

func TestAccMSOSchemaSiteAnpEpgUsegAttr_Basic(t *testing.T) {
	var useg1 models.SiteUsegAttr
	var useg2 models.SiteUsegAttr
	resourceName := "mso_schema_site_anp_epg_useg_attr.test"
	template_name := "Template1"
	usegName := makeTestVariable(acctest.RandString(5))
	usegNameOther := makeTestVariable(acctest.RandString(5))
	anpName := makeTestVariable(acctest.RandString(5))
	epgName := makeTestVariable(acctest.RandString(5))
	epgNameOther := makeTestVariable(acctest.RandString(5))
	usegTypeMAC := "mac"
	usegTypeTag := "tag"
	usegTypeIP := "ip"
	value := "00:00:5e:00:53:af"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "site_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "anp_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "epg_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "useg_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "useg_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, "value"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg1),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", template_name),
					resource.TestCheckResourceAttr(resourceName, "anp_name", anpName),
					resource.TestCheckResourceAttr(resourceName, "epg_name", epgName),
					resource.TestCheckResourceAttr(resourceName, "useg_name", usegName),
					resource.TestCheckResourceAttr(resourceName, "useg_type", usegTypeMAC),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// checking forceNew for the primary attr(useg_name) of the resource
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegNameOther, usegTypeMAC, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg2),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", template_name),
					resource.TestCheckResourceAttr(resourceName, "anp_name", anpName),
					resource.TestCheckResourceAttr(resourceName, "epg_name", epgName),
					resource.TestCheckResourceAttr(resourceName, "useg_name", usegNameOther),
					resource.TestCheckResourceAttr(resourceName, "useg_type", usegTypeMAC),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrIdNotEqual(&useg1, &useg2),
				),
			},
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value),
			},
			{
				// checking forceNew for the immediate parent
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgNameOther, usegName, usegTypeMAC, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg2),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", template_name),
					resource.TestCheckResourceAttr(resourceName, "anp_name", anpName),
					resource.TestCheckResourceAttr(resourceName, "epg_name", epgNameOther),
					resource.TestCheckResourceAttr(resourceName, "useg_name", usegName),
					resource.TestCheckResourceAttr(resourceName, "useg_type", usegTypeMAC),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrIdNotEqual(&useg1, &useg2),
				),
			},
			{
				// checking with 'tag' type so that we can check all the type dependent required params
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "startsWith", "MAC address"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg1),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", template_name),
					resource.TestCheckResourceAttr(resourceName, "anp_name", anpName),
					resource.TestCheckResourceAttr(resourceName, "epg_name", epgName),
					resource.TestCheckResourceAttr(resourceName, "useg_name", usegName),
					resource.TestCheckResourceAttr(resourceName, "useg_type", usegTypeTag),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "operator", "startsWith"),
					resource.TestCheckResourceAttr(resourceName, "category", "MAC address"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
				),
			},
			{
				// checking with 'ip' type so that we can check all the type dependent required params
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeIP, "10.0.0.1", "", "", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(resourceName, &useg1),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", template_name),
					resource.TestCheckResourceAttr(resourceName, "anp_name", anpName),
					resource.TestCheckResourceAttr(resourceName, "epg_name", epgName),
					resource.TestCheckResourceAttr(resourceName, "useg_name", usegName),
					resource.TestCheckResourceAttr(resourceName, "useg_type", usegTypeIP),
					resource.TestCheckResourceAttr(resourceName, "value", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "fv_subnet", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgUsegAttr_Negative(t *testing.T) {
	template_name := "Template1"
	usegName := makeTestVariable(acctest.RandString(5))
	anpName := makeTestVariable(acctest.RandString(5))
	epgName := makeTestVariable(acctest.RandString(5))
	usegTypeMAC := "mac"
	usegTypeTag := "tag"
	usegTypeIP := "ip"
	value := "00:00:5e:00:53:af"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				// max length check
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, acctest.RandString(1001), usegTypeMAC, value),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				// invalid schema attribute
				Config:      MSOSchemaSiteAnpEpgUsegAttrsForRandomAttrName(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				// invalid mac
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, randomValue),
				ExpectError: regexp.MustCompile(`expected "value" to be a valid MAC address`),
			},
			{
				// invalid ip
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeIP, randomValue, "", "", "false"),
				ExpectError: regexp.MustCompile(`expected value to contain a valid IP`),
			},
			{
				// invalid operator
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeTag, value, "starts_with", "MAC address"),
				ExpectError: regexp.MustCompile(`expected operator to be one of`),
			},
			{
				// invalid useg type
				Config:      MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, randomValue, value, "startsWith", "MAC address"),
				ExpectError: regexp.MustCompile(`expected useg_type to be one of`),
			},
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrWithRequired(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgUsegAttr_MultipleCreateDelete(t *testing.T) {
	usegName := makeTestVariable(acctest.RandString(5))
	anpName := makeTestVariable(acctest.RandString(5))
	epgName := makeTestVariable(acctest.RandString(5))
	usegTypeMAC := "mac"
	value := "00:00:5e:00:53:af"
	template_name := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSOSchemaSiteAnpEpgUsegAttrMultiple(siteNames[1], tenantNames[1], template_name, anpName, epgName, usegName, usegTypeMAC, value),
			},
		},
	})
}

func MSOSchemaSiteAnpEpgUsegAttrMultiple(site, tenant, template, anp, epg, useg, useg_type, value string) string {
	resource := CreatSchemaSiteConfig(site, tenant, template)
	resource += MSOSchemaSiteAnpEpg(anp, epg)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s_${count.index}"
		useg_type     = "%s"
		value         = "%s"
		count         = 3
	}
	`, useg, useg_type, value)
	return resource
}

func testAccCheckMSOSchemaSiteAnpEpgUsegAttrIdNotEqual(m1, m2 *models.SiteUsegAttr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := UsegAttrModelToUsegId(m1)
		id2 := UsegAttrModelToUsegId(m2)
		if id1 == id2 {
			return fmt.Errorf("Schema Site Anp Epg Useg Attr Ids are equal")
		}
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpEpgUsegAttrDestroy(s *terraform.State) error {
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

func testAccCheckMSOSchemaSiteAnpEpgUsegAttrExists(usegName string, m *models.SiteUsegAttr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, ok := s.RootModule().Resources[usegName]
		if !ok {
			return fmt.Errorf("Schema Site Anp Epg Useg Attr %s not found", usegName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schema Site Anp Epg Useg Attr ID was set")
		}
		useg, err := UsegIdToUsegAttrModel(rs.Primary.ID)
		if err != nil {
			return err
		}
		var read *models.SiteUsegAttr
		read, _, err = client.ReadAnpEpgUsegAttr(useg)
		if err != nil {
			return err
		}
		*m = *read
		return nil
	}
}

func MSOSchemaSiteAnpEpgUsegAttrWithoutRequired(site, tenant, template, anp, epg, useg, useg_type, value, keyAttr string) string {
	rBlock := CreatSchemaSiteConfig(site, tenant, template)
	rBlock += MSOSchemaSiteAnpEpg(anp, epg)
	switch keyAttr {
	case "schema_id":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
		#	schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	case "site_id":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
		#	site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	case "anp_name":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
		#	anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	case "epg_name":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
		#	epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	case "useg_name":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
		#	useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	case "useg_type":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
		#	useg_type     = "%s"
			value         = "%s"
		}`
	case "value":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
			template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
		#	value         = "%s"
		}`
	case "template_name":
		rBlock += `
		resource "mso_schema_site_anp_epg_useg_attr" "test" {
			schema_id     = mso_schema.test.id
			site_id       = data.mso_site.test.id
			anp_name      = mso_schema_site_anp.test.anp_name
			epg_name      = mso_schema_site_anp_epg.test.epg_name
		#	template_name = mso_schema_site.test.template_name
			useg_name     = "%s"
			useg_type     = "%s"
			value         = "%s"
		}`
	}
	return fmt.Sprintf(rBlock, useg, useg_type, value)
}

func MSOSchemaSiteAnpEpgUsegAttrsForRandomAttrName(site, tenant, template, anp, epg, useg, useg_type, useg_value, key, value string) string {
	resource := CreatSchemaSiteConfig(site, tenant, template)
	resource += MSOSchemaSiteAnpEpg(anp, epg)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		useg_type     = "%s"
		value		  = "%s"
		%s            = "%s"
	}
	`, useg, useg_type, useg_value, key, value)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrWithRequired(site, tenant, template, anp, epg, useg, useg_type, value string, opAttr ...string) string {
	// here opAttr param will be used for the extra attributes of other types
	// format of opAttr: type opAttr []string{operator, category, fvSubnet}
	// for 'ip' you have to pass first two opAttr as empty string("")
	// ----------------------------------------------------------------------
	resource := CreatSchemaSiteConfig(site, tenant, template)
	resource += MSOSchemaSiteAnpEpg(anp, epg)

	// other types
	if StringInSlice(useg_type, []string{"domain", "guest-os", "hv", "rootContName", "vm", "vm-name", "vnic"}) {
		resource += MSOSchemaSiteAnpEpgUsegAttrOtherTypes(useg, useg_type, opAttr[0], "", value)
		return resource
	}
	// 'tag' useg_type
	if useg_type == "tag" {
		resource += MSOSchemaSiteAnpEpgUsegAttrOtherTypes(useg, useg_type, opAttr[0], opAttr[1], value)
		return resource
	}
	// 'ip' useg_type
	if useg_type == "ip" {
		if opAttr[2] == "true" {
			resource += MSOSchemaSiteAnpEpgUsegAttrIPType(useg, "0.0.0.0", "true")
		} else {
			resource += MSOSchemaSiteAnpEpgUsegAttrIPType(useg, value, "false")
		}
		return resource
	}
	// 'mac' useg_type and other invalid types
	resource += MSOSchemaSiteAnpEpgUsegAttr(useg, useg_type, value)
	return resource
}

func MSOSchemaSiteAnpEpg(anp, epg string) string {
	resource := fmt.Sprintf(`
	resource "mso_schema_site_anp" "test" {
		schema_id     = mso_schema.test.id
		anp_name      = "%s"
		template_name = mso_schema_site.test.template_name
		site_id       = data.mso_site.test.id
	}

	resource "mso_schema_site_anp_epg" "test" {
		schema_id     = mso_schema.test.id
		template_name = mso_schema_site.test.template_name
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = "%s"
	}
	`, anp, epg)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttr(useg, useg_type, value string) string {
	// will be used for the mac type
	resource := fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		useg_type     = "%s"
		value         = "%s"
	}
	`, useg, useg_type, value)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrIPType(useg, value, fvSubnet string) string {
	resource := fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		useg_type     = "ip"
		value         = "%s"
		fv_subnet	  = "%s"
	}
	`, useg, value, fvSubnet)
	return resource
}

func MSOSchemaSiteAnpEpgUsegAttrOtherTypes(useg, useg_type, operator, category, value string) string {
	// pass category as empty string("") for other types except 'tag' type
	resource := ""
	if category == "" {
		resource += fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		useg_type     = "%s"
		operator      = "%s"
		value         = "%s"
	}
	`, useg, useg_type, operator, value)
	} else {
		resource += fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_useg_attr" "test" {
		schema_id     = mso_schema.test.id
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.test.anp_name
		epg_name      = mso_schema_site_anp_epg.test.epg_name
		template_name = mso_schema_site.test.template_name
		useg_name     = "%s"
		useg_type     = "%s"
		operator      = "%s"
		category      = "%s"
		value         = "%s"
	}
	`, useg, useg_type, operator, category, value)
	}
	return resource
}
