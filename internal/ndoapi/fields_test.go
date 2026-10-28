package ndoapi

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestStringField(t *testing.T) {
	object := map[string]any{"value": "configured", "empty": "", "null": nil, "invalid": 42}
	for _, test := range []struct {
		name, key   string
		requirement FieldRequirement
		want        string
		exists      bool
		errorText   string
	}{
		{name: "optional value", key: "value", want: "configured", exists: true},
		{name: "optional empty", key: "empty", exists: true},
		{name: "optional null", key: "null"},
		{name: "optional missing", key: "missing"},
		{name: "optional wrong type", key: "invalid", exists: true, errorText: `unexpected type int`},
		{name: "required value", key: "value", requirement: RequiredField, want: "configured", exists: true},
		{name: "required empty", key: "empty", requirement: RequiredField, exists: true},
		{name: "required null", key: "null", requirement: RequiredField, errorText: `is required`},
		{name: "required missing", key: "missing", requirement: RequiredField, errorText: `is required`},
		{name: "required nonempty value", key: "value", requirement: RequiredNonEmptyField, want: "configured", exists: true},
		{name: "required nonempty empty", key: "empty", requirement: RequiredNonEmptyField, exists: true, errorText: `must be nonempty`},
		{name: "required nonempty wrong type", key: "invalid", requirement: RequiredNonEmptyField, exists: true, errorText: `unexpected type int`},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, exists, err := StringField(object, test.key, test.requirement)
			checkFieldResult(t, exists, err, test.exists, test.errorText)
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestBoolField(t *testing.T) {
	object := map[string]any{"true": true, "false": false, "null": nil, "invalid": "true"}
	for _, test := range []struct {
		key         string
		requirement FieldRequirement
		want        bool
		exists      bool
		errorText   string
	}{
		{key: "true", want: true, exists: true},
		{key: "false", exists: true},
		{key: "null"},
		{key: "missing"},
		{key: "invalid", exists: true, errorText: `unexpected type string`},
		{key: "missing", requirement: RequiredField, errorText: `is required`},
	} {
		t.Run(fmt.Sprintf("%s/%d", test.key, test.requirement), func(t *testing.T) {
			got, exists, err := BoolField(object, test.key, test.requirement)
			checkFieldResult(t, exists, err, test.exists, test.errorText)
			if got != test.want {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestMapField(t *testing.T) {
	object := map[string]any{"value": map[string]any{"name": "ospf"}, "empty": map[string]any{}, "null": nil, "invalid": "map"}
	for _, test := range []struct {
		key         string
		requirement FieldRequirement
		exists      bool
		errorText   string
	}{
		{key: "value", exists: true},
		{key: "empty", exists: true},
		{key: "null"},
		{key: "missing"},
		{key: "invalid", exists: true, errorText: `unexpected type string`},
		{key: "missing", requirement: RequiredField, errorText: `is required`},
	} {
		t.Run(fmt.Sprintf("%s/%d", test.key, test.requirement), func(t *testing.T) {
			got, exists, err := MapField(object, test.key, test.requirement)
			checkFieldResult(t, exists, err, test.exists, test.errorText)
			if test.key == "value" && got["name"] != "ospf" {
				t.Fatalf("unexpected map: %#v", got)
			}
		})
	}
}

func TestListField(t *testing.T) {
	object := map[string]any{"value": []any{"one"}, "empty": []any{}, "null": nil, "invalid": "list"}
	for _, test := range []struct {
		key         string
		requirement FieldRequirement
		exists      bool
		errorText   string
	}{
		{key: "value", exists: true},
		{key: "empty", exists: true},
		{key: "null"},
		{key: "missing"},
		{key: "invalid", exists: true, errorText: `unexpected type string`},
		{key: "missing", requirement: RequiredField, errorText: `is required`},
	} {
		t.Run(fmt.Sprintf("%s/%d", test.key, test.requirement), func(t *testing.T) {
			got, exists, err := ListField(object, test.key, test.requirement)
			checkFieldResult(t, exists, err, test.exists, test.errorText)
			if test.key == "value" && (len(got) != 1 || got[0] != "one") {
				t.Fatalf("unexpected list: %#v", got)
			}
		})
	}
}

func TestInt64Field(t *testing.T) {
	object := map[string]any{
		"int": int(10), "int32": int32(-10), "int64": int64(42), "json": float64(17),
		"fraction": 1.5, "overflow": math.Exp2(63), "nan": math.NaN(), "invalid": "17", "null": nil,
	}
	for _, test := range []struct {
		key         string
		requirement FieldRequirement
		want        int64
		exists      bool
		errorText   string
	}{
		{key: "int", want: 10, exists: true},
		{key: "int32", want: -10, exists: true},
		{key: "int64", want: 42, exists: true},
		{key: "json", want: 17, exists: true},
		{key: "fraction", exists: true, errorText: `whole number within int64 range`},
		{key: "overflow", exists: true, errorText: `whole number within int64 range`},
		{key: "nan", exists: true, errorText: `whole number within int64 range`},
		{key: "invalid", exists: true, errorText: `unexpected type string`},
		{key: "null"},
		{key: "missing"},
		{key: "missing", requirement: RequiredField, errorText: `is required`},
	} {
		t.Run(fmt.Sprintf("%s/%d", test.key, test.requirement), func(t *testing.T) {
			got, exists, err := Int64Field(object, test.key, test.requirement)
			checkFieldResult(t, exists, err, test.exists, test.errorText)
			if got != test.want {
				t.Fatalf("got %d, want %d", got, test.want)
			}
		})
	}
}

func checkFieldResult(t *testing.T, exists bool, err error, wantExists bool, errorText string) {
	t.Helper()
	if exists != wantExists {
		t.Fatalf("exists = %v, want %v", exists, wantExists)
	}
	if errorText == "" && err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if errorText != "" && (err == nil || !strings.Contains(err.Error(), errorText)) {
		t.Fatalf("expected error containing %q, got %v", errorText, err)
	}
}
