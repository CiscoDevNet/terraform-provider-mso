package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTemplateDatasourceTenantErrors(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: No template_id or template_name provided in Template configuration") },
				Config:      testAccMSOTemplateDatasourceErrorNoIdOrNameConfig(),
				ExpectError: regexp.MustCompile("either `template_id` or `template_name` must be provided"),
			},
			{
				PreConfig:   func() { fmt.Println("Test: No template_type with name provided in Template configuration") },
				Config:      testAccMSOTemplateDatasourceErrorNoTypeConfig(),
				ExpectError: regexp.MustCompile("`template_type` must be provided when `template_name` is provided"),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Both template_id and template_name provided in Template configuration") },
				Config:      testAccMSOTemplateDatasourceErrorIdAndNameConfig(),
				ExpectError: regexp.MustCompile("only one of `template_id` or `template_name` must be provided"),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Non existing template name provided in Template configuration") },
				Config:      testAccMSOTemplateDatasourceErrorNonExistingConfig(),
				ExpectError: regexp.MustCompile("Template with name 'non_existing_template_name' not found."),
			},
		},
	})
}

func TestAccMSOTemplateDatasourceTenantName(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Tenant template with name and type provided in Template configuration") },
				Config:    testAccMSOTemplateDatasourceNameAndTypeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"data.mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1, msoTemplateSiteName2},
						},
						false,
					),
				),
			},
		},
	})
}

func TestAccMSOTemplateDatasourceTenantId(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Tenant template with id provided in Template configuration") },
				Config:    testAccMSOTemplateDatasourceIdConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"data.mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1, msoTemplateSiteName2},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateDatasourceNameAndTypeConfig() string {
	return fmt.Sprintf(`%s
	data "mso_template" "template_tenant" {
		template_name = mso_template.template_tenant.template_name
		template_type = "tenant"
	}
	`, testAccMSOTemplateResourceTenanTwoSitesConfig())
}

func testAccMSOTemplateDatasourceIdConfig() string {
	return fmt.Sprintf(`%s
	data "mso_template" "template_tenant" {
		template_id = mso_template.template_tenant.id
	}
	`, testAccMSOTemplateResourceTenanTwoSitesConfig())
}

func testAccMSOTemplateDatasourceErrorNoIdOrNameConfig() string {
	return fmt.Sprintf(`
	data "mso_template" "template_tenant" {
		template_type = "tenant"
	}
	`)
}

func testAccMSOTemplateDatasourceErrorNoTypeConfig() string {
	return fmt.Sprintf(`
	data "mso_template" "template_tenant" {
		template_name = "non_existing_template_name"
	}
	`)
}

func testAccMSOTemplateDatasourceErrorIdAndNameConfig() string {
	return fmt.Sprintf(`
	data "mso_template" "template_tenant" {
		template_id = "non_existing_template_id"
		template_name = "non_existing_template_name"
		template_type = "tenant"
	}
	`)
}

func testAccMSOTemplateDatasourceErrorNonExistingConfig() string {
	return fmt.Sprintf(`
	data "mso_template" "template_tenant" {
		template_name = "non_existing_template_name"
		template_type = "tenant"
	}
	`)
}
