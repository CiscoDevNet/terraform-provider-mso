package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMSOL3OutAnnotationDataSource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation_data")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject annotation datasource with a nonexistent template") },
				Config:      testAccMSOL3OutAnnotationDataSourceInvalidTemplateConfig(),
				ExpectError: regexp.MustCompile(`L3Out Template Not Found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read a standalone L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationDataSourceCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out_annotation.test", "id", "mso_l3out_annotation.owner", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_annotation.test", "template_id", "mso_template.l3out_test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_annotation.test", "l3out_uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttr("data.mso_l3out_annotation.test", "key", "owner"),
					resource.TestCheckResourceAttr("data.mso_l3out_annotation.test", "value", "network"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Refresh the standalone L3Out annotation datasource after an update") },
				Config:    testAccMSOL3OutAnnotationDataSourceUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out_annotation.test", "value", "operations"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_annotation.test", "id", "mso_l3out_annotation.owner", "id"),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject annotation datasource with a nonexistent L3Out") },
				Config:      testAccMSOL3OutAnnotationDataSourceMissingParentConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Not Found`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject a missing standalone L3Out annotation") },
				Config:      testAccMSOL3OutAnnotationDataSourceMissingConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Annotation Not Found`),
			},
		},
	})
}

func testAccMSOL3OutAnnotationDataSourceInvalidTemplateConfig() string {
	return `
data "mso_l3out_annotation" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = "owner"
}
`
}

func testAccMSOL3OutAnnotationDataSourceMissingParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName) + `
data "mso_l3out_annotation" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = "owner"

  depends_on = [mso_l3out_annotation.owner]
}
`
}

func testAccMSOL3OutAnnotationDataSourceCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName) + `
data "mso_l3out_annotation" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"

  depends_on = [mso_l3out_annotation.owner]
}
`
}

func testAccMSOL3OutAnnotationDataSourceUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName) + `
data "mso_l3out_annotation" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"

  depends_on = [mso_l3out_annotation.owner]
}
`
}

func testAccMSOL3OutAnnotationDataSourceMissingConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName) + `
data "mso_l3out_annotation" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "missing"

  depends_on = [mso_l3out_annotation.owner]
}
`
}
