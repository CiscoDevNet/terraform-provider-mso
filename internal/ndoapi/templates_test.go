package ndoapi

import (
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/container"
)

func TestTemplateEndpoint(t *testing.T) {
	tests := []struct {
		templateID string
		expected   string
	}{
		{templateID: "template-123", expected: "api/v1/templates/template-123"},
		{templateID: "template/123", expected: "api/v1/templates/template%2F123"},
	}

	for _, test := range tests {
		t.Run(test.templateID, func(t *testing.T) {
			if actual := TemplateEndpoint(test.templateID); actual != test.expected {
				t.Fatalf("unexpected template endpoint: got %q, want %q", actual, test.expected)
			}
		})
	}
}

func TestTemplateFromResponse(t *testing.T) {
	response, err := container.ParseJSON([]byte(`{
		"l3outTemplate": {
			"l3outs": [{"uuid": "l3out-123", "name": "example"}]
		}
	}`))
	if err != nil {
		t.Fatalf("parsing template response: %v", err)
	}
	template, err := TemplateFromResponse("template-123", response)
	if err != nil {
		t.Fatalf("decoding template response: %v", err)
	}
	resolved, found, err := template.Find(NewPath(
		PathStep{Field: "l3outTemplate"},
		PathStep{
			Field:    "l3outs",
			Selector: &ObjectSelector{UUID: ObjectIdentifier{Field: "uuid", Value: "l3out-123"}},
		},
	))
	if err != nil {
		t.Fatalf("finding L3Out in the decoded template: %v", err)
	}
	if !found || resolved.Object["uuid"] != "l3out-123" {
		t.Fatalf("expected the L3Out in the decoded template, got found=%t object=%v", found, resolved.Object)
	}
}

func TestTemplateFromResponseRejectsMalformedResponses(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected string
	}{
		{
			name:     "nil response",
			expected: `template "template-123": empty response body`,
		},
		{
			name:     "array response",
			payload:  `[]`,
			expected: "unexpected template response type []interface {}",
		},
		{
			name:     "null response",
			payload:  `null`,
			expected: "unexpected template response type <nil>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response *container.Container
			if test.payload != "" {
				var err error
				response, err = container.ParseJSON([]byte(test.payload))
				if err != nil {
					t.Fatalf("parsing template response: %v", err)
				}
			}
			template, err := TemplateFromResponse("template-123", response)
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Fatalf("expected error containing %q, got %v", test.expected, err)
			}
			if template != nil {
				t.Fatal("expected no template for a malformed response")
			}
		})
	}
}
