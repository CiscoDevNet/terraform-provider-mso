package ndoapi

import (
	"testing"
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
