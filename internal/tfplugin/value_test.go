package tfplugin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestKnownBoolPointer(t *testing.T) {
	trueValue := types.BoolValue(true)
	falseValue := types.BoolValue(false)

	if got := KnownBoolPointer(types.BoolNull()); got != nil {
		t.Fatal("expected null bool to convert to nil")
	}
	if got := KnownBoolPointer(types.BoolUnknown()); got != nil {
		t.Fatal("expected unknown bool to convert to nil")
	}
	if got := KnownBoolPointer(trueValue); got == nil || !*got {
		t.Fatalf("expected true bool pointer, got %#v", got)
	}
	if got := KnownBoolPointer(falseValue); got == nil || *got {
		t.Fatalf("expected false bool pointer, got %#v", got)
	}
}

func TestBoolOrNull(t *testing.T) {
	if got := BoolOrNull(false, false); !got.IsNull() {
		t.Fatal("expected absent bool to convert to null")
	}
	if got := BoolOrNull(true, true); got.IsNull() || got.IsUnknown() || !got.ValueBool() {
		t.Fatalf("expected known true bool, got %#v", got)
	}
	if got := BoolOrNull(false, true); got.IsNull() || got.IsUnknown() || got.ValueBool() {
		t.Fatalf("expected known false bool, got %#v", got)
	}
}

func TestKnownStringPointer(t *testing.T) {
	value := "value"
	known := types.StringValue(value)

	if got := KnownStringPointer(types.StringNull()); got != nil {
		t.Fatal("expected null string to convert to nil")
	}
	if got := KnownStringPointer(types.StringUnknown()); got != nil {
		t.Fatal("expected unknown string to convert to nil")
	}
	if got := KnownStringPointer(known); got == nil || *got != value {
		t.Fatalf("unexpected string pointer: %#v", got)
	}
}

func TestKnownInt64Pointer(t *testing.T) {
	if got := KnownInt64Pointer(types.Int64Null()); got != nil {
		t.Fatal("expected null int64 to convert to nil")
	}
	if got := KnownInt64Pointer(types.Int64Unknown()); got != nil {
		t.Fatal("expected unknown int64 to convert to nil")
	}
	if got := KnownInt64Pointer(types.Int64Value(42)); got == nil || *got != 42 {
		t.Fatalf("unexpected int64 pointer: %#v", got)
	}
}

func TestStringOrNull(t *testing.T) {
	translate := func(value string) string {
		return value + "-translated"
	}

	if got := StringOrNull("ignored", false, translate); !got.IsNull() {
		t.Fatal("expected absent value to convert to null")
	}
	if got := StringOrNull("value", true, nil); got.ValueString() != "value" {
		t.Fatalf("unexpected untransformed value: got %q", got.ValueString())
	}
	if got := StringOrNull("value", true, translate); got.ValueString() != "value-translated" {
		t.Fatalf("unexpected transformed value: got %q", got.ValueString())
	}
	if got := StringOrNull("", true, nil); got.IsNull() || got.IsUnknown() || got.ValueString() != "" {
		t.Fatalf("expected known empty value, got %#v", got)
	}
}

func TestStringOrEmpty(t *testing.T) {
	translate := func(value string) string {
		return value + "-translated"
	}

	if got := StringOrEmpty("ignored", false, translate); got.IsNull() || got.ValueString() != "" {
		t.Fatalf("expected absent value to convert to known empty string, got %#v", got)
	}
	if got := StringOrEmpty("value", true, translate); got.ValueString() != "value-translated" {
		t.Fatalf("unexpected transformed value: got %q", got.ValueString())
	}
}

func TestStringMapConversions(t *testing.T) {
	ctx := context.Background()
	expected := map[string]string{"blue": "10", "red": "20"}
	terraformValue, diagnostics := types.MapValueFrom(ctx, types.StringType, expected)
	if diagnostics.HasError() {
		t.Fatalf("creating Terraform map: %v", diagnostics)
	}

	fromTerraformDiagnostics := makeDiagnostics()
	if got := StringMapFromTerraform(ctx, terraformValue, &fromTerraformDiagnostics); !mapsEqual(got, expected) {
		t.Fatalf("unexpected map from Terraform: got %#v, want %#v", got, expected)
	}
	if fromTerraformDiagnostics.HasError() {
		t.Fatalf("converting Terraform map: %v", fromTerraformDiagnostics)
	}
	if got := StringMapFromTerraform(ctx, types.MapNull(types.StringType), &fromTerraformDiagnostics); got != nil {
		t.Fatalf("expected null map to convert to nil, got %#v", got)
	}
	if got := StringMapFromTerraform(ctx, types.MapUnknown(types.StringType), &fromTerraformDiagnostics); got != nil {
		t.Fatalf("expected unknown map to convert to nil, got %#v", got)
	}

	toTerraformDiagnostics := makeDiagnostics()
	if got := StringMap(ctx, nil, &toTerraformDiagnostics); !got.IsNull() {
		t.Fatal("expected nil map to convert to null")
	}
	got := StringMap(ctx, expected, &toTerraformDiagnostics)
	if toTerraformDiagnostics.HasError() {
		t.Fatalf("converting map to Terraform: %v", toTerraformDiagnostics)
	}
	actual := map[string]string{}
	got.ElementsAs(ctx, &actual, false)
	if !mapsEqual(actual, expected) {
		t.Fatalf("unexpected Terraform map: got %#v, want %#v", actual, expected)
	}
}

func makeDiagnostics() diag.Diagnostics {
	return nil
}

func mapsEqual(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range right {
		if left[key] != value {
			return false
		}
	}
	return true
}
