package mso

// TestAccMSORestDataSource exercises the `data "mso_rest"` block. Two steps:
//
//  1. Negative: a path that does not exist on MSO must surface an error from
//     MakeRestRequest. We run this step first so it has minimal state to
//     refresh.
//  2. Positive: build a real schema via the existing resource helpers
//     (mso_schema -> template -> vrf -> bd -> anp -> epg) and read it back
//     through data.mso_rest at api/v1/schemas/<id>. Asserts `content` is set
//     and that the JSON body carries the expected displayName, ANP and EPG.

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSORestDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: data.mso_rest with non-existent schema id expects error")
				},
				// Use a syntactically valid but non-existent schema id so MSO
				// returns a structured 404 JSON body (caught by CheckForErrors).
				// A bare unknown collection path tends to return a 200 with an
				// empty body and silently passes.
				Config: `
data "mso_rest" "bad" {
  path = "api/v1/schemas/000000000000000000000000"
}
`,
				ExpectError: regexp.MustCompile(`(?i)(not\s*found|404|invalid|error)`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: data.mso_rest reads the schema created via the prereq resources")
				},
				Config: testAccMSORestDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_rest.schema", "content"),
					testAccCheckMSORestDataSourceContent(),
				),
			},
		},
	})
}

// testAccMSORestDataSourceConfig builds a real schema using the existing
// helpers and then reads it back via data.mso_rest.
func testAccMSORestDataSourceConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s%s
data "mso_rest" "schema" {
  # Reference the EPG's schema_id attribute so the path is "known after apply"
  # of the EPG resource. This naturally defers the data source read until the
  # ANP/EPG children have been PATCHed into the schema -- no depends_on
  # required (depends_on on data sources in SDK v1 leaves the post-apply plan
  # non-empty).
  path = "api/v1/schemas/${mso_schema_template_anp_epg.%s.schema_id}"
}
`,
		testSiteConfigAnsibleTest(),
		testTenantConfig(),
		testSchemaConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateAnpConfig(),
		testSchemaTemplateAnpEpgConfig(),
		msoSchemaTemplateAnpEpgName,
	)
}

// testAccCheckMSORestDataSourceContent decodes the data source's `content`
// attribute and asserts the schema document carries the expected name and
// the prereq ANP/EPG.
func testAccCheckMSORestDataSourceContent() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.mso_rest.schema"]
		if !ok {
			return fmt.Errorf("data.mso_rest.schema not found in state")
		}
		raw := rs.Primary.Attributes["content"]
		if raw == "" {
			return fmt.Errorf("data.mso_rest.schema content is empty")
		}

		var body struct {
			DisplayName string `json:"displayName"`
			Templates   []struct {
				Name string `json:"name"`
				Anps []struct {
					Name string `json:"name"`
					Epgs []struct {
						Name string `json:"name"`
					} `json:"epgs"`
				} `json:"anps"`
			} `json:"templates"`
		}
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			return fmt.Errorf("data.mso_rest.schema content is not valid JSON: %s", err)
		}
		if body.DisplayName != msoSchemaName {
			return fmt.Errorf("displayName = %q, want %q", body.DisplayName, msoSchemaName)
		}
		foundAnp, foundEpg := false, false
		for _, tmpl := range body.Templates {
			for _, anp := range tmpl.Anps {
				if anp.Name == msoSchemaTemplateAnpName {
					foundAnp = true
				}
				for _, epg := range anp.Epgs {
					if epg.Name == msoSchemaTemplateAnpEpgName {
						foundEpg = true
					}
				}
			}
		}
		if !foundAnp {
			return fmt.Errorf("anp %q not found in schema content", msoSchemaTemplateAnpName)
		}
		if !foundEpg {
			return fmt.Errorf("epg %q not found in schema content", msoSchemaTemplateAnpEpgName)
		}
		return nil
	}
}
