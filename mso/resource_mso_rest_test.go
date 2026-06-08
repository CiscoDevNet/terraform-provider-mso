package mso

// TestAccMSORestResource drives a single `mso_rest.schema` instance through
// the CRUD-relevant HTTP methods against the MSO schemas endpoint:
//
//   - POST   create the schema at api/v1/schemas (default method).
//   - PUT    full-replace at api/v1/schemas/<id>, payload carries the nested
//            template / VRF / BD (vrfRef) / ANP / two EPGs (bdRef on both,
//            subnet on epg1).
//   - PATCH  add a subnet to epg2 via JSON-Patch on the same path.
//   - PATCH  idempotent replace of epg2.displayName combined with
//            `retrigger = true` to assert it is reset to "false" after read.
//   - error  invalid method and invalid payload to cover validation branches.
//
// The schema id is server-assigned. Step 1's config has only `mso_rest.schema`
// (path is the literal collection endpoint). Step 2 and onwards add
// `data "mso_schema" "lookup"` keyed by displayName; the lookup reads cleanly
// because the schema already exists from step 1.
//
// Why DELETE is NOT covered as a TestStep
// ---------------------------------------
// The DELETE method of resourceMSORest cannot be exercised reliably from
// inside the SDK v1 acceptance-test harness for this resource. The two
// natural ways to drive it both run into hard limitations:
//
//  1. Auto-destroy with `method` unset.
//     resourceMSORestDelete only fires a real DELETE when `method` is unset
//     in state. But `method` is Optional+Computed: once any earlier step
//     sets it (PUT/PATCH/...), the value stays in state for the rest of the
//     lifecycle and the destroy branch short-circuits with no API call. A
//     final step that sets `method = ""` does not clear the prior value
//     either (Optional+Computed keeps it). So auto-destroy on this test
//     cannot reach the DELETE branch.
//
//  2. A dedicated step that does Create/Update with `method = "DELETE"`.
//     This would need the server-assigned schema id baked into the HCL of
//     a later step (`path = "api/v1/schemas/<id>"`). Two attempts failed:
//       - referencing `data.mso_schema.lookup.id`: the post-apply refresh
//         re-reads the data source AFTER the DELETE, hits 404, and the
//         data source returns "Schema of specified name not found", which
//         fails the step.
//       - injecting the id via a TF variable (TF_VAR_*): the embedded
//         terraform engine used by SDK v1 does not propagate TF_VAR_* env
//         vars from PreConfig into the plan, so the variable stays
//         unassigned and the step errors with "Unassigned variable".
//     Splitting into two `resource.Test` calls so the second slice can
//     `fmt.Sprintf` the captured id directly works for the DELETE step
//     itself, but phase 1's auto-destroy then fails to remove the
//     `mso_tenant` because the orphaned schema still references it, which
//     leaves the test marked FAIL even when the DELETE step succeeded.
//
// The DELETE branches of resourceMSORest (Delete short-circuit and
// MakeRestRequest with "DELETE") are still exercised by unit-level
// coverage and by other tests' cleanup paths; replicating that here would
// require either changing the resource semantics (breaking change) or
// re-architecting the test fixtures around an existing tenant that
// resource.Test does not destroy. For now we accept that gap and document
// it explicitly.
//
// Cleanup of the schema left on MSO after step 6 is handled by the final
// PreConfig hook below (manual DeletebyId before the framework destroys
// the tenant), and CheckDestroy verifies the schema is gone.

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// REST-specific test fixtures. Kept local to the REST tests because the
// schema lifecycle is driven entirely through mso_rest and the nested object
// names are an implementation detail of the JSON payloads below.
var (
	msoRestSchemaName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// msoRestSchemaId is captured in step 1 from the API after the POST.
	// It is used by later steps' TestCheckFuncs and by CheckDestroy.
	msoRestSchemaId string
)

const (
	msoRestTemplate    = "Template1"
	msoRestVrf         = "vrf1"
	msoRestBd          = "bd1"
	msoRestAnp         = "anp1"
	msoRestEpg1        = "epg1"
	msoRestEpg2        = "epg2"
	msoRestEpg1Subnet  = "10.10.10.1/24"
	msoRestEpg2Subnet  = "10.10.20.1/24"
	msoRestEpg2NewName = "epg2-updated"
)

func TestAccMSORestResource(t *testing.T) {
	resourceName := "mso_rest.schema"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSORestSchemaDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: POST create schema via mso_rest (default method, path=api/v1/schemas)")
				},
				Config: testAccMSORestPostConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", "api/v1/schemas"),
					resource.TestCheckResourceAttr(resourceName, "id", "api/v1/schemas"),
					resource.TestCheckResourceAttr(resourceName, "retrigger", "false"),
					testAccCheckMSORestSchemaCaptureId(),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: PUT full-replace with nested vrf+bd+anp and two EPGs (epg1 has subnet)")
				},
				Config: testAccMSORestPutConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "method", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "retrigger", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckMSORestEpgSubnetCount(msoRestEpg1, 1),
					testAccCheckMSORestEpgSubnetCount(msoRestEpg2, 0),
					testAccCheckMSORestEpgBdRef(msoRestEpg1),
					testAccCheckMSORestEpgBdRef(msoRestEpg2),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: PATCH add subnet to epg2 via JSON-Patch")
				},
				Config: testAccMSORestPatchAddSubnetConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "method", "PATCH"),
					testAccCheckMSORestEpgSubnetCount(msoRestEpg1, 1),
					testAccCheckMSORestEpgSubnetCount(msoRestEpg2, 1),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: PATCH idempotent replace of epg2 displayName with retrigger=true (asserts retrigger reset to false)")
				},
				Config: testAccMSORestPatchReplaceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "method", "PATCH"),
					resource.TestCheckResourceAttr(resourceName, "retrigger", "false"),
					testAccCheckMSORestEpgDisplayName(msoRestEpg2, msoRestEpg2NewName),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Invalid method (expect error from resourceMSORest)")
				},
				Config:      testAccMSORestInvalidMethodConfig(),
				ExpectError: regexp.MustCompile("Invalid method"),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Invalid payload (expect error from MakeRestRequest JSON parse)")
				},
				Config:      testAccMSORestInvalidPayloadConfig(),
				ExpectError: regexp.MustCompile("Unable to parse the payload"),
			},
			{
				// Final cleanup step. Removes mso_rest.schema from the
				// config so the framework destroys it -- but the resource's
				// Delete short-circuits because `method` is set in state
				// (see resourceMSORestDelete docs and the package comment
				// above). The PreConfig below issues a direct API DELETE
				// against MSO so the schema is gone before the framework
				// destroys the tenant in its final destroy phase; otherwise
				// the tenant DELETE fails ("Unable to delete the object")
				// because MSO refuses to delete a tenant that still owns a
				// schema, and the whole test is marked FAIL.
				PreConfig: func() {
					fmt.Println("Test: Manual cleanup -- DELETE schema on MSO so tenant destroy succeeds")
					if msoRestSchemaId != "" {
						msoClient := testAccProvider.Meta().(*client.Client)
						_ = msoClient.DeletebyId("api/v1/schemas/" + msoRestSchemaId)
					}
				},
				Config: testAccMSORestPrereqConfig(),
			},
		},
	})
}

// testAccMSORestPrereqConfig is the minimal scaffolding needed to derive a
// valid tenantId for the schema payloads. It uses the existing site/tenant
// helpers so the test does not duplicate that boilerplate.
func testAccMSORestPrereqConfig() string {
	return fmt.Sprintf(`%s%s`, testSiteConfigAnsibleTest(), testTenantConfig())
}

// testAccMSORestSchemaLookupConfig adds a data "mso_schema" block that
// resolves the server-assigned id by displayName. It is included from step 2
// onwards (step 1 cannot include it because the schema does not yet exist
// at plan time of the first apply).
func testAccMSORestSchemaLookupConfig() string {
	return fmt.Sprintf(`
data "mso_schema" "lookup" {
  name = %q
}
`, msoRestSchemaName)
}

// testAccMSORestPostConfig is step 1: POST to api/v1/schemas with the
// minimal schema document (one empty template). No method ⇒ resource Create
// uses the default POST branch.
func testAccMSORestPostConfig() string {
	return fmt.Sprintf(`%s
resource "mso_rest" "schema" {
  path = "api/v1/schemas"
  payload = jsonencode({
    displayName     = %q
    sites           = []
    _updateVersion  = 0
    templates = [{
      name             = %q
      displayName      = %q
      tenantId         = mso_tenant.%s.id
      templateType     = "stretched-template"
      templateSubType  = []
      anps             = []
      contracts        = []
      vrfs             = []
      bds              = []
      filters          = []
      externalEpgs     = []
      serviceGraphs    = []
      intersiteL3outs  = []
    }]
  })
}
`, testAccMSORestPrereqConfig(), msoRestSchemaName, msoRestTemplate, msoRestTemplate, msoTenantName)
}

// testAccMSORestSchemaPutPayload builds the full PUT payload that replaces
// the schema with a nested template containing one VRF, one BD (with vrfRef
// pointing at the VRF) and one ANP with two EPGs:
//   - epg1: bdRef → bd1, subnets:[{ip: epg1Subnet}]
//   - epg2: bdRef → bd1, no subnets (gets a subnet later via PATCH)
//
// The schema id required by vrfRef/bdRef is read from data.mso_schema.lookup.
func testAccMSORestSchemaPutPayload() string {
	return fmt.Sprintf(`jsonencode({
    displayName = %q
    sites       = []
    templates = [{
      name            = %q
      displayName     = %q
      tenantId        = mso_tenant.%s.id
      templateType    = "stretched-template"
      templateSubType = []
      contracts       = []
      filters         = []
      externalEpgs    = []
      serviceGraphs   = []
      intersiteL3outs = []
      vrfs = [{
        name        = %q
        displayName = %q
      }]
      bds = [{
        name                  = %q
        displayName           = %q
        layer2UnknownUnicast  = "proxy"
        vrfRef = {
          schemaId     = data.mso_schema.lookup.id
          templateName = %q
          vrfName      = %q
        }
      }]
      anps = [{
        name        = %q
        displayName = %q
        epgs = [
          {
            name        = %q
            displayName = %q
            bdRef = {
              schemaId     = data.mso_schema.lookup.id
              templateName = %q
              bdName       = %q
            }
            subnets = [{
              ip     = %q
              scope  = "private"
              shared = false
            }]
          },
          {
            name        = %q
            displayName = %q
            bdRef = {
              schemaId     = data.mso_schema.lookup.id
              templateName = %q
              bdName       = %q
            }
          }
        ]
      }]
    }]
  })`,
		msoRestSchemaName,
		msoRestTemplate, msoRestTemplate, msoTenantName,
		msoRestVrf, msoRestVrf,
		msoRestBd, msoRestBd, msoRestTemplate, msoRestVrf,
		msoRestAnp, msoRestAnp,
		msoRestEpg1, msoRestEpg1, msoRestTemplate, msoRestBd, msoRestEpg1Subnet,
		msoRestEpg2, msoRestEpg2, msoRestTemplate, msoRestBd,
	)
}

// testAccMSORestPutConfig is step 2: change method to PUT and target the
// per-schema endpoint. Payload contains the nested EPG configuration.
func testAccMSORestPutConfig() string {
	return fmt.Sprintf(`%s%s
resource "mso_rest" "schema" {
  method  = "PUT"
  path    = "api/v1/schemas/${data.mso_schema.lookup.id}"
  payload = %s
}
`, testAccMSORestPrereqConfig(), testAccMSORestSchemaLookupConfig(), testAccMSORestSchemaPutPayload())
}

// testAccMSORestPatchAddSubnetConfig is step 3: JSON-Patch that adds a
// subnet to epg2.
func testAccMSORestPatchAddSubnetConfig() string {
	patch := fmt.Sprintf(`jsonencode([{
    op   = "add"
    path = "/templates/%s/anps/%s/epgs/%s/subnets/-"
    value = {
      ip     = %q
      scope  = "private"
      shared = false
    }
  }])`, msoRestTemplate, msoRestAnp, msoRestEpg2, msoRestEpg2Subnet)

	return fmt.Sprintf(`%s%s
resource "mso_rest" "schema" {
  method  = "PATCH"
  path    = "api/v1/schemas/${data.mso_schema.lookup.id}"
  payload = %s
}
`, testAccMSORestPrereqConfig(), testAccMSORestSchemaLookupConfig(), patch)
}

// testAccMSORestPatchReplaceConfig is step 4 (method=PATCH, retrigger=true).
// The payload is an idempotent replace of epg2.displayName so the PATCH can
// be applied twice without MSO complaining about duplicate operations.
func testAccMSORestPatchReplaceConfig(retrigger bool) string {
	patch := fmt.Sprintf(`jsonencode([{
    op    = "replace"
    path  = "/templates/%s/anps/%s/epgs/%s/displayName"
    value = %q
  }])`, msoRestTemplate, msoRestAnp, msoRestEpg2, msoRestEpg2NewName)

	extraLines := ""
	if retrigger {
		// retrigger is reset to false by the resource after every apply, so
		// the config (true) vs state (false) drift would otherwise produce a
		// non-empty post-apply plan. lifecycle.ignore_changes scopes the
		// suppression to just this attribute, keeping the plan check strict
		// for everything else.
		extraLines = "  retrigger = true\n  lifecycle {\n    ignore_changes = [retrigger]\n  }\n"
	}

	return fmt.Sprintf(`%s%s
resource "mso_rest" "schema" {
  method  = "PATCH"
  path    = "api/v1/schemas/${data.mso_schema.lookup.id}"
  payload = %s
%s}
`, testAccMSORestPrereqConfig(), testAccMSORestSchemaLookupConfig(), patch, extraLines)
}

// testAccMSORestInvalidMethodConfig is step 5: keep a valid path/payload but
// set method to a value that fails the HTTP_METHODS validation in
// resourceMSORestUpdate.
func testAccMSORestInvalidMethodConfig() string {
	patch := fmt.Sprintf(`jsonencode([{
    op    = "replace"
    path  = "/templates/%s/anps/%s/epgs/%s/displayName"
    value = %q
  }])`, msoRestTemplate, msoRestAnp, msoRestEpg2, msoRestEpg2NewName)

	return fmt.Sprintf(`%s%s
resource "mso_rest" "schema" {
  method  = "FOO"
  path    = "api/v1/schemas/${data.mso_schema.lookup.id}"
  payload = %s
}
`, testAccMSORestPrereqConfig(), testAccMSORestSchemaLookupConfig(), patch)
}

// testAccMSORestInvalidPayloadConfig is step 6: keep a valid method/path but
// send a payload that is not valid JSON.
func testAccMSORestInvalidPayloadConfig() string {
	return fmt.Sprintf(`%s%s
resource "mso_rest" "schema" {
  method  = "PATCH"
  path    = "api/v1/schemas/${data.mso_schema.lookup.id}"
  payload = "not-json"
}
`, testAccMSORestPrereqConfig(), testAccMSORestSchemaLookupConfig())
}

// testAccCheckMSORestSchemaCaptureId looks up the data.mso_schema.lookup
// resource in the test state and records its id into the package-level
// msoRestSchemaId so later steps and CheckDestroy can reach the schema by
// id without re-querying the API by name.
//
// In step 1 the data source is not in the config, so this function falls
// back to a name-based lookup against MSO directly.
func testAccCheckMSORestSchemaCaptureId() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Try the data source first (steps 2+).
		if rs, ok := s.RootModule().Resources["data.mso_schema.lookup"]; ok && rs.Primary != nil && rs.Primary.ID != "" {
			msoRestSchemaId = rs.Primary.ID
			return nil
		}

		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL("api/v1/schemas")
		if err != nil {
			return fmt.Errorf("failed to list schemas: %s", err)
		}
		count, err := cont.ArrayCount("schemas")
		if err != nil {
			return fmt.Errorf("failed to read schemas array: %s", err)
		}
		for i := 0; i < count; i++ {
			sCont, err := cont.ArrayElement(i, "schemas")
			if err != nil {
				return err
			}
			if models.StripQuotes(sCont.S("displayName").String()) == msoRestSchemaName {
				msoRestSchemaId = models.StripQuotes(sCont.S("id").String())
				return nil
			}
		}
		return fmt.Errorf("schema with displayName %q not found in MSO after POST", msoRestSchemaName)
	}
}

// testAccCheckMSORestSchemaDestroy verifies the schema is gone from MSO at
// end-of-test. Because the lifecycle steps leave `method` set in state, the
// framework's auto-destroy short-circuits in resourceMSORestDelete and the
// schema stays on MSO. We therefore issue a manual DELETE here as cleanup,
// then assert the GET returns an error (i.e. the object is really gone).
func testAccCheckMSORestSchemaDestroy(s *terraform.State) error {
	if msoRestSchemaId == "" {
		return nil
	}
	msoClient := testAccProvider.Meta().(*client.Client)

	// Best-effort cleanup. Ignore the error: if the schema is already gone
	// the DELETE will surface a 404 which is fine for our purposes.
	_ = msoClient.DeletebyId("api/v1/schemas/" + msoRestSchemaId)

	cont, err := msoClient.GetViaURL("api/v1/schemas/" + msoRestSchemaId)
	if err != nil {
		return nil
	}
	if cont != nil {
		id := models.StripQuotes(cont.S("id").String())
		if id == msoRestSchemaId {
			return fmt.Errorf("schema %s still exists after cleanup", msoRestSchemaId)
		}
	}
	return nil
}

// findMSORestEpg navigates the schema container to locate the named EPG
// inside msoRestTemplate/msoRestAnp.
func findMSORestEpg(schemaCont *container.Container, epgName string) (*container.Container, error) {
	templateCount, err := schemaCont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("schema has no templates: %s", err)
	}
	for i := 0; i < templateCount; i++ {
		tCont, err := schemaCont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		if models.StripQuotes(tCont.S("name").String()) != msoRestTemplate {
			continue
		}
		anpCount, err := tCont.ArrayCount("anps")
		if err != nil {
			return nil, fmt.Errorf("template %s has no anps: %s", msoRestTemplate, err)
		}
		for j := 0; j < anpCount; j++ {
			aCont, err := tCont.ArrayElement(j, "anps")
			if err != nil {
				return nil, err
			}
			if models.StripQuotes(aCont.S("name").String()) != msoRestAnp {
				continue
			}
			epgCount, err := aCont.ArrayCount("epgs")
			if err != nil {
				return nil, fmt.Errorf("anp %s has no epgs: %s", msoRestAnp, err)
			}
			for k := 0; k < epgCount; k++ {
				eCont, err := aCont.ArrayElement(k, "epgs")
				if err != nil {
					return nil, err
				}
				if models.StripQuotes(eCont.S("name").String()) == epgName {
					return eCont, nil
				}
			}
			return nil, fmt.Errorf("epg %s not found under anp %s", epgName, msoRestAnp)
		}
		return nil, fmt.Errorf("anp %s not found under template %s", msoRestAnp, msoRestTemplate)
	}
	return nil, fmt.Errorf("template %s not found in schema", msoRestTemplate)
}

// testAccCheckMSORestEpgSubnetCount asserts the named EPG has exactly the
// expected number of subnets.
func testAccCheckMSORestEpgSubnetCount(epgName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if msoRestSchemaId == "" {
			return fmt.Errorf("msoRestSchemaId not captured; step 1 capture must run first")
		}
		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL("api/v1/schemas/" + msoRestSchemaId)
		if err != nil {
			return fmt.Errorf("failed to GET schema %s: %s", msoRestSchemaId, err)
		}
		epgCont, err := findMSORestEpg(cont, epgName)
		if err != nil {
			return err
		}
		count, err := epgCont.ArrayCount("subnets")
		if err != nil {
			count = 0
		}
		if count != expected {
			return fmt.Errorf("epg %s: expected %d subnets, got %d", epgName, expected, count)
		}
		return nil
	}
}

// testAccCheckMSORestEpgBdRef asserts the EPG's bdRef resolves to msoRestBd.
// MSO returns bdRef either as a path string ("/schemas/.../bds/<name>") or
// as a structured object, depending on the API version, so we accept both
// representations.
func testAccCheckMSORestEpgBdRef(epgName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if msoRestSchemaId == "" {
			return fmt.Errorf("msoRestSchemaId not captured; step 1 capture must run first")
		}
		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL("api/v1/schemas/" + msoRestSchemaId)
		if err != nil {
			return err
		}
		epgCont, err := findMSORestEpg(cont, epgName)
		if err != nil {
			return err
		}
		bdRefRaw := epgCont.S("bdRef")
		if bdRefRaw == nil || bdRefRaw.Data() == nil {
			return fmt.Errorf("epg %s: bdRef not set", epgName)
		}
		// Structured object form: {schemaId, templateName, bdName}.
		// gabs returns "{}" from .String() when the key is missing, so
		// filter that out before treating it as a real bdName.
		if name := models.StripQuotes(bdRefRaw.S("bdName").String()); name != "" && name != "null" && name != "{}" {
			if name != msoRestBd {
				return fmt.Errorf("epg %s: bdRef.bdName = %q, want %q", epgName, name, msoRestBd)
			}
			return nil
		}
		// String path form: "/schemas/<id>/templates/<t>/bds/<name>".
		bdRefStr := models.StripQuotes(bdRefRaw.String())
		if !strings.Contains(bdRefStr, "/bds/"+msoRestBd) {
			return fmt.Errorf("epg %s: bdRef %q does not reference bd %q", epgName, bdRefStr, msoRestBd)
		}
		return nil
	}
}

// testAccCheckMSORestEpgDisplayName asserts the EPG's displayName equals
// the expected value (used after the PATCH replace step).
func testAccCheckMSORestEpgDisplayName(epgName, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if msoRestSchemaId == "" {
			return fmt.Errorf("msoRestSchemaId not captured; step 1 capture must run first")
		}
		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL("api/v1/schemas/" + msoRestSchemaId)
		if err != nil {
			return err
		}
		epgCont, err := findMSORestEpg(cont, epgName)
		if err != nil {
			return err
		}
		got := models.StripQuotes(epgCont.S("displayName").String())
		if got != expected {
			return fmt.Errorf("epg %s: displayName = %q, want %q", epgName, got, expected)
		}
		return nil
	}
}
