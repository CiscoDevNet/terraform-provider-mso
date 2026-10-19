package mso

import (
	"reflect"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/container"
)

func TestSplitCommaString(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty", "", []string{}},
		{"whitespace_only", "   \t  ", []string{}},
		{"single_token", "1/1", []string{"1/1"}},
		{"single_token_with_spaces", " 1/1 ", []string{"1/1"}},
		{"multiple_tokens", "1/1,1/2", []string{"1/1", "1/2"}},
		{"multiple_tokens_with_spaces", " 1/1 , 1/2 ", []string{"1/1", "1/2"}},
		{"drops_empty_tokens_leading_trailing_commas", ",1/1,1/2,", []string{"1/1", "1/2"}},
		{"drops_empty_tokens_repeated_commas", "1/1,,1/2", []string{"1/1", "1/2"}},
		{"all_empty_tokens", ",,,", []string{}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := splitCommaString(testCase.input)
			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("splitCommaString(%q) = %#v, expected %#v", testCase.input, got, testCase.expected)
			}
		})
	}
}

// TestBuildTaskErrorMessage verifies that buildTaskErrorMessage surfaces the
// site-specific failure text found in operDetails.execSiteStatus (e.g. the
// real APIC validation error) instead of only the generic errMessage from
// operDetails.detailedStatus. See issue #548.
func TestBuildTaskErrorMessage(t *testing.T) {
	cases := []struct {
		name              string
		taskJSON          string
		defaultErrMessage string
		expected          string
	}{
		{
			name:              "no_detailed_status_or_exec_site_status_returns_default",
			taskJSON:          `{"operDetails": {}}`,
			defaultErrMessage: "Could not determine specific deployment error message.",
			expected:          "Could not determine specific deployment error message.",
		},
		{
			name: "detailed_status_errmessage_without_failed_sites",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "Template deployment failed: all stretch fabrics failed to deploy"}],
				"execSiteStatus": [{"siteID": "site1", "status": {"siteStatus": "Succeeded"}}]
			}}`,
			defaultErrMessage: "default",
			expected:          "Template deployment failed: all stretch fabrics failed to deploy",
		},
		{
			name: "includes_single_failed_site_message",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "Template deployment failed: all stretch fabrics failed to deploy"}],
				"execSiteStatus": [{"siteID": "67ad8a45b867a59a418e6550", "status": {"siteStatus": "Failed", "msg": "code: \"100\", text: \"Validation failed: Port has encap: vlan-1001 which is also configured for being used as encapsulation 802.1q on another EPG.\""}}]
			}}`,
			defaultErrMessage: "default",
			expected:          `Template deployment failed: all stretch fabrics failed to deploy: Site 67ad8a45b867a59a418e6550: code: "100", text: "Validation failed: Port has encap: vlan-1001 which is also configured for being used as encapsulation 802.1q on another EPG."`,
		},
		{
			name: "includes_multiple_failed_sites_joined_with_semicolon",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "deploy failed"}],
				"execSiteStatus": [
					{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "error on site1"}},
					{"siteID": "site2", "status": {"siteStatus": "Failed", "msg": "error on site2"}}
				]
			}}`,
			defaultErrMessage: "default",
			expected:          "deploy failed: Site site1: error on site1; Site site2: error on site2",
		},
		{
			name: "ignores_succeeded_sites_and_sites_without_msg",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "deploy failed"}],
				"execSiteStatus": [
					{"siteID": "site1", "status": {"siteStatus": "Succeeded"}},
					{"siteID": "site2", "status": {"siteStatus": "Failed"}}
				]
			}}`,
			defaultErrMessage: "default",
			expected:          "deploy failed",
		},
		{
			name: "does_not_duplicate_site_message_already_present_in_errmessage",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "Template deployment failed: code: \"100\", text: \"Validation failed\""}],
				"execSiteStatus": [{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "code: \"100\", text: \"Validation failed\""}}]
			}}`,
			defaultErrMessage: "default",
			expected:          `Template deployment failed: code: "100", text: "Validation failed"`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cont, err := container.ParseJSON([]byte(testCase.taskJSON))
			if err != nil {
				t.Fatalf("failed to parse test task JSON: %s", err)
			}
			got := buildTaskErrorMessage(cont, testCase.defaultErrMessage)
			if got != testCase.expected {
				t.Fatalf("buildTaskErrorMessage() = %q, expected %q", got, testCase.expected)
			}
		})
	}
}
