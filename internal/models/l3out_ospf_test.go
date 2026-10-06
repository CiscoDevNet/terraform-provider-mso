package models

import (
	"strings"
	"testing"
)

func TestL3OutOSPFModelSetFromNDOObjectRejectsMalformedFields(t *testing.T) {
	tests := []struct {
		name     string
		object   map[string]any
		expected string
	}{
		{
			name:     "invalid area configuration",
			object:   map[string]any{"ospfAreaConfig": "invalid"},
			expected: `NDO field "ospfAreaConfig" has unexpected type`,
		},
		{
			name:     "invalid area ID",
			object:   map[string]any{"ospfAreaConfig": map[string]any{"id": 42}},
			expected: `NDO field "id" has unexpected type`,
		},
		{
			name:     "invalid control",
			object:   map[string]any{"ospfAreaConfig": map[string]any{"control": "invalid"}},
			expected: `NDO field "control" has unexpected type`,
		},
		{
			name:     "invalid control value",
			object:   map[string]any{"ospfAreaConfig": map[string]any{"control": map[string]any{"redistribute": "true"}}},
			expected: `NDO field "redistribute" has unexpected type`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var model L3OutOSPFModel
			if err := model.setFromNDOObject(test.object); err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Fatalf("expected error containing %q, got %v", test.expected, err)
			}
		})
	}
}
