package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMSOL3OutAnnotationResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create one standalone L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttrPair(
						"mso_l3out_annotation.owner",
						"template_id",
						"mso_template.l3out_test",
						"id",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out_annotation.owner",
						"l3out_uuid",
						"mso_l3out.test",
						"uuid",
					),
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
				PreConfig:   func() { fmt.Println("Test: Reject an empty standalone L3Out annotation value") },
				Config:      testAccMSOL3OutAnnotationEmptyValueConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`(?s)Key\s+or\s+Value\s+cannot\s+be\s+empty\s+for\s+Annotations`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Preserve the prior value after rejection and add another annotation")
					if err := testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "operations"},
					)(nil); err != nil {
						t.Fatalf("prior annotation changed after rejected empty value: %v", err)
					}
				},
				Config: testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "operations"),
					resource.TestCheckResourceAttr("mso_l3out_annotation.purpose", "value", "operations"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "operations", "purpose": "operations"},
					),
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
					if attributes["template_id"] != templateID ||
						attributes["l3out_uuid"] != l3outUUID || attributes["key"] != "owner" {
						return fmt.Errorf(
							"unexpected imported L3Out annotation identity: %q/%q/%q",
							attributes["template_id"],
							attributes["l3out_uuid"],
							attributes["key"],
						)
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
					resource.TestCheckResourceAttrPair(
						"mso_l3out_annotation.owner",
						"l3out_uuid",
						"mso_l3out.test",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out_annotation.purpose",
						"l3out_uuid",
						"mso_l3out.test",
						"uuid",
					),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "operations", "purpose": "operations"},
					),
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
				PreConfig: func() { fmt.Println("Test: Replace an annotation when its key changes") },
				Config:    testAccMSOL3OutAnnotationRenamedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(
						"mso_l3out_annotation.purpose",
						plancheck.ResourceActionDestroyBeforeCreate,
					)},
				},
				Check: testAccCheckL3OutAnnotations(
					&templateID,
					&l3outUUID,
					map[string]string{"service": "operations"},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete the last standalone L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationParentConfig(siteName, l3outName),
				Check:     testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{}),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create an annotation with a computed value") },
				Config:    testAccMSOL3OutAnnotationComputedValueConfig(siteName, l3outName, "computed-network"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("mso_l3out_annotation.owner", tfjsonpath.New("value")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "computed-network"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "computed-network"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update an annotation with a computed value") },
				Config:    testAccMSOL3OutAnnotationComputedValueConfig(siteName, l3outName, "computed-operations"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("mso_l3out_annotation.owner", tfjsonpath.New("value")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "computed-operations"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "computed-operations"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete the annotation after its computed value changes") },
				Config:    testAccMSOL3OutAnnotationParentConfig(siteName, l3outName),
				Check:     testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{}),
			},
		},
	})
}

func testAccMSOL3OutAnnotationParentConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "network"
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func testAccMSOL3OutAnnotationComputedValueConfig(siteName, l3outName, value string) string {
	return fmt.Sprintf(`%s
resource "terraform_data" "annotation_value" {
  input = %q
}

resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = terraform_data.annotation_value.output
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName), value)
}

func testAccMSOL3OutAnnotationDuplicateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "duplicate" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "duplicate"

  depends_on = [mso_l3out_annotation.owner]
}
`, testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName))
}

func testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "operations"
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func testAccMSOL3OutAnnotationTwoConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "purpose"
  value       = "operations"

  depends_on = [mso_l3out_annotation.owner]
}
`, testAccMSOL3OutAnnotationUpdateConfig(siteName, l3outName))
}

func testAccMSOL3OutAnnotationRemainingConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "purpose"
  value       = "operations"
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func TestAccMSOL3OutAnnotationInvalidReferences(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation_invalid")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
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
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "invalid" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = "owner"
  value       = "network"
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func TestAccMSOL3OutAnnotationOwnership(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation_ownership")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Manage a standalone annotation with the parent map omitted") },
				Config:    testAccMSOL3OutAnnotationCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "network"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "network"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the L3Out without managing its annotations") },
				Config:    testAccMSOL3OutAnnotationOmittedParentUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "parent updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.owner", "network"),
					resource.TestCheckResourceAttr("mso_l3out_annotation.owner", "value", "network"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "network"}),
				),
			},
		},
	})
}

func testAccMSOL3OutAnnotationOmittedParentUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "parent updated"
}

resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = "network"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func TestAccMSOL3OutAnnotationOwnershipConflict(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotation_conflict")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Configure one annotation on the L3Out") },
				Config:    testAccMSOL3OutAnnotationOwnershipParentConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.owner", "network"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"owner": "network"}),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Expect parent drift from a child annotation") },
				Config:             testAccMSOL3OutAnnotationOwnershipConflictConfig(siteName, l3outName),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"mso_l3out.test",
							tfjsonpath.New("annotations"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"owner": knownvalue.StringExact("network"),
							}),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_annotation.purpose", "value", "routing"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"owner": "network", "purpose": "routing"},
					),
				),
			},
		},
	})
}

func testAccMSOL3OutAnnotationOwnershipParentConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  annotations = {
    owner = "network"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationOwnershipConflictConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "purpose"
  value       = "routing"
}
`, testAccMSOL3OutAnnotationOwnershipParentConfig(siteName, l3outName))
}
func testAccCheckL3OutAnnotationCompositeID(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in Terraform state", resourceName)
		}
		attributes := resourceState.Primary.Attributes
		expected := fmt.Sprintf("%s/%s/%s", attributes["template_id"], attributes["l3out_uuid"], attributes["key"])
		if attributes["id"] != expected {
			return fmt.Errorf("unexpected L3Out annotation ID: got %q, want %q", attributes["id"], expected)
		}
		return nil
	}
}

func testAccMSOL3OutAnnotationEmptyValueConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "owner"
  value       = ""
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func testAccMSOL3OutAnnotationRenamedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  key         = "service"
  value       = "operations"
}
`, testAccMSOL3OutAnnotationParentConfig(siteName, l3outName))
}

func TestAccMSOL3OutAnnotationInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Reject an empty standalone annotation key") },
				Config: `resource "mso_l3out_annotation" "invalid" {
  template_id = "00000000-0000-0000-0000-000000000000"
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  key         = ""
  value       = "network"
}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject an empty key in the inline annotation map") },
				Config: `resource "mso_l3out" "invalid" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "invalid_annotations"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  annotations = { "" = "network" }
}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject a null value in the inline annotation map") },
				Config: `resource "mso_l3out" "invalid" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "invalid_annotations"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  annotations = { owner = null }
}`,
				ExpectError: regexp.MustCompile(`Null Map Value`),
			},
		},
	})
}
