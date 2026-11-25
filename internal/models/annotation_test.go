package models

import (
	"strings"
	"testing"
)

func TestAnnotationsModelSetFromNDOObjectRejectsMalformedFields(t *testing.T) {
	tests := []struct {
		name     string
		object   map[string]any
		expected string
	}{
		{
			name:     "invalid entry",
			object:   map[string]any{"tagAnnotations": []any{42}},
			expected: "tagAnnotations[0] has unexpected type int",
		},
		{
			name: "duplicate key",
			object: map[string]any{
				"tagAnnotations": []any{
					map[string]any{"tagKey": "owner", "tagValue": "network"},
					map[string]any{"tagKey": "owner", "tagValue": "operations"},
				},
			},
			expected: `duplicate tagKey "owner"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var model AnnotationsModel
			if err := model.SetFromNDOObject(test.object); err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Fatalf("expected error containing %q, got %v", test.expected, err)
			}
		})
	}
}
