// Package tfplugin provides Terraform Plugin Framework conversion, planning, and validation helpers.
package tfplugin

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KnownBoolPointer returns nil for null or unknown Terraform values.
func KnownBoolPointer(value types.Bool) *bool {
	if value.IsUnknown() {
		return nil
	}
	return value.ValueBoolPointer()
}

// BoolOrNull converts an optional boolean value to Terraform state while
// preserving an explicitly configured false value.
func BoolOrNull(value bool, exists bool) types.Bool {
	if !exists {
		return types.BoolNull()
	}
	return types.BoolValue(value)
}

// KnownStringPointer returns nil for null or unknown Terraform values.
func KnownStringPointer(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}
	return value.ValueStringPointer()
}

// KnownInt64Pointer returns nil for null or unknown Terraform values.
func KnownInt64Pointer(value types.Int64) *int64 {
	if value.IsUnknown() {
		return nil
	}
	return value.ValueInt64Pointer()
}

// StringOrNull converts an optional string value to Terraform state. The
// translation function is applied only when the value is present and may be
// nil when no translation is needed.
func StringOrNull(value string, exists bool, translate func(string) string) types.String {
	if !exists {
		return types.StringNull()
	}
	if translate != nil {
		value = translate(value)
	}
	return types.StringValue(value)
}

// StringOrEmpty converts an optional string value to Terraform state while
// representing an absent API field as an explicitly known empty string.
// The translation function is applied only when the value is present.
func StringOrEmpty(value string, exists bool, translate func(string) string) types.String {
	if !exists {
		return types.StringValue("")
	}
	if translate != nil {
		value = translate(value)
	}
	return types.StringValue(value)
}

// StringMapFromTerraform converts a known Terraform string map to a Go map.
func StringMapFromTerraform(ctx context.Context, value types.Map, diagnostics *diag.Diagnostics) map[string]string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := make(map[string]string, len(value.Elements()))
	diagnostics.Append(value.ElementsAs(ctx, &result, false)...)
	return result
}

// StringMap converts a Go string map to a Terraform map, preserving nil as null.
func StringMap(ctx context.Context, value map[string]string, diagnostics *diag.Diagnostics) types.Map {
	if value == nil {
		return types.MapNull(types.StringType)
	}
	result, resultDiagnostics := types.MapValueFrom(ctx, types.StringType, value)
	diagnostics.Append(resultDiagnostics...)
	return result
}
