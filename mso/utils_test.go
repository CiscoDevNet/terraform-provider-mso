package mso

import (
	"encoding/json"
	"reflect"
	"strings"
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
		siteNames         map[string]string
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
			expected: "Template deployment failed: all stretch fabrics failed to deploy\n" +
				"Site 67ad8a45b867a59a418e6550:\n" +
				"  code: \"100\", text: \"Validation failed:\n" +
				"  Port has encap: vlan-1001 which is also configured for being used as encapsulation 802.1q on another EPG.\"",
		},
		{
			name: "includes_multiple_failed_sites_on_separate_lines",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "deploy failed"}],
				"execSiteStatus": [
					{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "error on site1"}},
					{"siteID": "site2", "status": {"siteStatus": "Failed", "msg": "error on site2"}}
				]
			}}`,
			defaultErrMessage: "default",
			expected:          "deploy failed\nSite site1:\n  error on site1\nSite site2:\n  error on site2",
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
			name: "moves_site_message_from_errmessage_to_site_entry",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "Template deployment failed: code: \"100\", text: \"Validation failed\""}],
				"execSiteStatus": [
					{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "code: \"100\", text: \"Validation failed\""}},
					{"siteID": "site2", "status": {"siteStatus": "Failed", "msg": "error on site2"}}
				]
			}}`,
			defaultErrMessage: "default",
			expected:          "Template deployment failed\nSite site1:\n  code: \"100\", text: \"Validation failed\"\nSite site2:\n  error on site2",
		},
		{
			name: "omits_errmessage_identical_to_site_message",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "error on site1"}],
				"execSiteStatus": [{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "error on site1"}}]
			}}`,
			defaultErrMessage: "default",
			expected:          "Site site1:\n  error on site1",
		},
		{
			name: "labels_sites_with_name_and_falls_back_to_id",
			taskJSON: `{"operDetails": {
				"detailedStatus": [{"errMessage": "deploy failed"}],
				"execSiteStatus": [
					{"siteID": "site1", "status": {"siteStatus": "Failed", "msg": "error on site1"}},
					{"siteID": "site2", "status": {"siteStatus": "Failed", "msg": "error on site2"}}
				]
			}}`,
			defaultErrMessage: "default",
			siteNames:         map[string]string{"site1": "ansible_test"},
			expected:          "deploy failed\nSite ansible_test (site1):\n  error on site1\nSite site2:\n  error on site2",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cont, err := container.ParseJSON([]byte(testCase.taskJSON))
			if err != nil {
				t.Fatalf("failed to parse test task JSON: %s", err)
			}
			got := buildTaskErrorMessage(cont, testCase.defaultErrMessage, testCase.siteNames)
			if got != testCase.expected {
				t.Fatalf("buildTaskErrorMessage() = %q, expected %q", got, testCase.expected)
			}
		})
	}
}

func TestFormatTaskErrorMessage(t *testing.T) {
	cases := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:    "literal_newline_separators",
			message: `code: "100", text: "Validation failed: Port: eth1/1 has encap: vlan-2000 on another EPG\nPort: eth1/1 has encap: vlan-2001 on another EPG\n"`,
			expected: []string{
				`code: "100", text: "Validation failed:`,
				`Port: eth1/1 has encap: vlan-2000 on another EPG`,
				`Port: eth1/1 has encap: vlan-2001 on another EPG"`,
			},
		},
		{
			name:    "period_space_separators_without_newlines",
			message: `text: "Validation failed: Port: eth1/1 has encap: vlan-2000 on another EPG. Port: eth1/1 has encap: vlan-2001 on another EPG."`,
			expected: []string{
				`text: "Validation failed:`,
				`Port: eth1/1 has encap: vlan-2000 on another EPG.`,
				`Port: eth1/1 has encap: vlan-2001 on another EPG."`,
			},
		},
		{
			name:     "period_space_not_split_when_newlines_present",
			message:  `first fault. Second sentence\nthird fault`,
			expected: []string{`first fault. Second sentence`, `third fault`},
		},
		{
			name:    "period_space_followed_by_lowercase_split_and_dn_not_split",
			message: `on another EPG.  and uni/tn-t1/ap-a1/epg-e1/rspathAtt-[topology/pod-1/paths-101/pathep-[eth1/1]]"`,
			expected: []string{
				`on another EPG.`,
				`and uni/tn-t1/ap-a1/epg-e1/rspathAtt-[topology/pod-1/paths-101/pathep-[eth1/1]]"`,
			},
		},
		{
			name:    "run_together_faults_split",
			message: `Port: eth1/1 has more than 1 native encapPort: eth1/1 has encap: vlan-2000 on another EPG`,
			expected: []string{
				`Port: eth1/1 has more than 1 native encap`,
				`Port: eth1/1 has encap: vlan-2000 on another EPG`,
			},
		},
		{
			name:     "class_name_with_consecutive_capitals_not_split",
			message:  `fvAEPg is invalid`,
			expected: []string{`fvAEPg is invalid`},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := formatTaskErrorMessage(testCase.message)
			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("formatTaskErrorMessage(%q) = %#v, expected %#v", testCase.message, got, testCase.expected)
			}
		})
	}
}

func TestIsTaskErrorMessageTruncated(t *testing.T) {
	long := strings.Repeat("x", taskErrorTruncationThreshold)
	cases := []struct {
		name     string
		message  string
		expected bool
	}{
		{"short_without_full_stop", "Validation failed", false},
		{"long_cut_off_mid_word", long + ` topology/pod-1/paths-101/pathe"`, true},
		{"long_ending_with_full_stop_and_quote", long + ` on another EPG."`, false},
		{"long_ending_with_full_stop_literal_newline_and_quote", long + ` on another EPG.\n"`, false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := isTaskErrorMessageTruncated(testCase.message); got != testCase.expected {
				t.Fatalf("isTaskErrorMessageTruncated() = %t, expected %t", got, testCase.expected)
			}
		})
	}
}

func TestBuildTaskErrorMessageTruncatedSiteMessage(t *testing.T) {
	truncatedMsg := (`code: "100", text: "Validation failed: ` +
		strings.Repeat(`Port: topology/pod-1/paths-101/pathep-[eth1/1] has encap: vlan-2000 on another EPG\n`, 40))[:2450]
	task := map[string]interface{}{
		"operDetails": map[string]interface{}{
			"detailedStatus": []interface{}{map[string]interface{}{"errMessage": "deploy failed"}},
			"execSiteStatus": []interface{}{
				map[string]interface{}{"siteID": "site1", "status": map[string]interface{}{"siteStatus": "Failed", "msg": "short error on site1."}},
				map[string]interface{}{"siteID": "site2", "status": map[string]interface{}{"siteStatus": "Failed", "msg": truncatedMsg}},
			},
		},
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("failed to marshal test task JSON: %s", err)
	}
	cont, err := container.ParseJSON(taskJSON)
	if err != nil {
		t.Fatalf("failed to parse test task JSON: %s", err)
	}

	got := buildTaskErrorMessage(cont, "default", nil)
	if count := strings.Count(got, taskErrorTruncationNotice); count != 1 {
		t.Fatalf("expected exactly one truncation notice, got %d in %q", count, got)
	}
	if !strings.HasPrefix(got, "deploy failed\nSite site1:\n  short error on site1.\nSite site2:\n") {
		t.Fatalf("unexpected message structure: %q", got)
	}
	if !strings.HasSuffix(got, "\n  "+taskErrorTruncationNotice) {
		t.Fatalf("expected indented truncation notice after site2 message, got %q", got)
	}
}

func TestParseSiteNamesById(t *testing.T) {
	cont, err := container.ParseJSON([]byte(`{"sites": [
		{"id": "nd-site", "common": {"name": "ansible-test"}},
		{"id": "mso-site", "name": "mso-test"},
		{"id": "no-name"}
	]}`))
	if err != nil {
		t.Fatalf("failed to parse test sites JSON: %s", err)
	}
	expected := map[string]string{"nd-site": "ansible-test", "mso-site": "mso-test"}
	if got := parseSiteNamesById(cont); !reflect.DeepEqual(got, expected) {
		t.Fatalf("parseSiteNamesById() = %#v, expected %#v", got, expected)
	}
}
