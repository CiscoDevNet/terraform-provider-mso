package mso

import (
	"testing"
)

func TestParseChildResourceId(t *testing.T) {
	cases := []struct {
		name             string
		importId         string
		separator        string
		expectedParentId string
		expectedChild    string
		expectError      bool
	}{
		{
			name:             "valid_context_separator",
			importId:         "templateId/abc123/RouteMapPolicy/test/context/ctx_1",
			separator:        "/context/",
			expectedParentId: "templateId/abc123/RouteMapPolicy/test",
			expectedChild:    "ctx_1",
			expectError:      false,
		},
		{
			name:             "valid_with_different_separator",
			importId:         "parentId/entry/myEntry",
			separator:        "/entry/",
			expectedParentId: "parentId",
			expectedChild:    "myEntry",
			expectError:      false,
		},
		{
			name:             "child_name_contains_separator_characters",
			importId:         "templateId/abc123/RouteMapPolicy/test/context/ctx/with/slashes",
			separator:        "/context/",
			expectedParentId: "templateId/abc123/RouteMapPolicy/test",
			expectedChild:    "ctx/with/slashes",
			expectError:      false,
		},
		{
			name:        "missing_separator",
			importId:    "templateId/abc123/RouteMapPolicy/test",
			separator:   "/context/",
			expectError: true,
		},
		{
			name:        "empty_import_id",
			importId:    "",
			separator:   "/context/",
			expectError: true,
		},
		{
			name:        "empty_parent_id",
			importId:    "/context/ctx_1",
			separator:   "/context/",
			expectError: true,
		},
		{
			name:        "empty_child_name",
			importId:    "templateId/abc123/RouteMapPolicy/test/context/",
			separator:   "/context/",
			expectError: true,
		},
		{
			name:        "separator_only",
			importId:    "/context/",
			separator:   "/context/",
			expectError: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parentId, childName, err := ParseChildResourceId(testCase.importId, testCase.separator)

			if testCase.expectError {
				if err == nil {
					t.Fatalf("ParseChildResourceId(%q, %q) expected error, got parentId=%q childName=%q", testCase.importId, testCase.separator, parentId, childName)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseChildResourceId(%q, %q) unexpected error: %s", testCase.importId, testCase.separator, err)
			}
			if parentId != testCase.expectedParentId {
				t.Fatalf("ParseChildResourceId(%q, %q) parentId = %q, expected %q", testCase.importId, testCase.separator, parentId, testCase.expectedParentId)
			}
			if childName != testCase.expectedChild {
				t.Fatalf("ParseChildResourceId(%q, %q) childName = %q, expected %q", testCase.importId, testCase.separator, childName, testCase.expectedChild)
			}
		})
	}
}
