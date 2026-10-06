package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMSOL3OutAnnotationResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create one standalone L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttrPair("mso_l3out_annotation.owner", "template_id", "mso_template.l3out_test", "id"),
					resource.TestCheckResourceAttrPair("mso_l3out_annotation.owner", "l3out_uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "key", "owner"),
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "network"),
					testAccCheckL3OutAnnotationCompositeID("mso_l3out_annotation.owner"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "network"}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate an annotation deleted out of band")
					testAccReplaceL3OutAnnotations(t, templateID, l3outUUID, models.AnnotationsModel{})
				},
				Config: testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "network"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "network"}),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject a second resource for the same L3Out annotation key") },
				Config:      testAccMSOL3OutAnnotationDuplicateConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Annotation Already Exists`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update a standalone L3Out annotation value") },
				Config:    testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "operations"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "operations"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create another L3Out annotation with the same value") },
				Config:    testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "operations"),
					resource.TestCheckResourceAttr("mso_l3out_annotation.purpose", "value", "operations"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "operations", "purpose": "operations"}),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import one standalone L3Out annotation") },
				Config:       testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName),
				ResourceName: "mso_l3out_annotation.owner",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					annotation, ok := state.RootModule().Resources["mso_l3out_annotation.owner"]
					if !ok {
						return "", fmt.Errorf("mso_l3out_annotation.owner not found in Terraform state")
					}
					return annotation.Primary.Attributes["id"], nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig:     func() { fmt.Println("Test: Reject invalid L3Out annotation import ID") },
				ResourceName:  "mso_l3out_annotation.owner",
				ImportState:   true,
				ImportStateId: "invalid-l3out-annotation-import-id",
				ExpectError: regexp.MustCompile(
					`The import ID must be in the form <template_id>/<l3out_uuid>/<key>`,
				),
			},
			{
				PreConfig:       func() { fmt.Println("Test: Import L3Out annotation by resource identity") },
				ResourceName:    "mso_l3out_annotation.owner",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported L3Out annotation, got %d", len(states))
					}
					attributes := states[0].Attributes
					if attributes["template_id"] != templateID || attributes["l3out_uuid"] != l3outUUID || attributes["key"] != "owner" {
						return fmt.Errorf("unexpected imported L3Out annotation identity: %q/%q/%q", attributes["template_id"], attributes["l3out_uuid"], attributes["key"])
					}
					return nil
				},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate annotations after their L3Out is deleted out of band")
					testAccDeleteL3Out(t, templateID, l3outUUID)
				},
				Config: testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttrPair("mso_l3out_annotation.owner", "l3out_uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_annotation.purpose", "l3out_uuid", "mso_l3out.test", "uuid"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "operations", "purpose": "operations"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete one L3Out annotation and preserve the other") },
				Config:    testAccMSOL3OutAnnotationRemainingConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.purpose", "value", "operations"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"purpose": "operations"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete the last standalone L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationParentConfig(siteName, l3outName),
				Check:     testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{}),
			},
		},
	})
}

func testAccMSOL3OutAnnotationParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`
}

func testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationParentConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "network"
}
`
}

func testAccMSOL3OutAnnotationDuplicateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "duplicate" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "duplicate"

  depends_on = [mso_l3out_annotation.owner]
}
`
}

func testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationParentConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "operations"
}
`
}

func testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "purpose"
  value       = "operations"

  depends_on = [mso_l3out_annotation.owner]
}
`
}

func testAccMSOL3OutAnnotationRemainingConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationParentConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "purpose"
  value       = "operations"
}
`
}

func TestAccMSOL3OutAnnotationInvalidReferences(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation_invalid")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject annotation creation in a nonexistent template") },
				Config:      testAccMSOL3OutAnnotationInvalidTemplateConfig(),
				ExpectError: regexp.MustCompile(`L3Out Template Not Found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create annotation parent L3Out") },
				Config:    testAccMSOL3OutAnnotationParentConfig(siteName, l3outName),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject annotation creation on a nonexistent L3Out") },
				Config:      testAccMSOL3OutAnnotationMissingParentConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Not Found`),
			},
		},
	})
}

func testAccMSOL3OutAnnotationInvalidTemplateConfig() string {
	return `
resource "mso_l3out_annotation" "invalid" {
  template_id = "00000000-0000-0000-0000-000000000000"
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = "owner"
  value       = "network"
}
`
}

func testAccMSOL3OutAnnotationMissingParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutAnnotationParentConfig(siteName, l3outName) + `
resource "mso_l3out_annotation" "invalid" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = "owner"
  value       = "network"
}
`
}

func testAccCheckL3OutAnnotationCompositeID(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in Terraform state", resourceName)
		}
		attributes := resourceState.Primary.Attributes
		expected := attributes["template_id"] + "/" + attributes["l3out_uuid"] + "/" + attributes["key"]
		if attributes["id"] != expected {
			return fmt.Errorf("unexpected L3Out annotation ID: got %q, want %q", attributes["id"], expected)
		}
		return nil
	}
}
