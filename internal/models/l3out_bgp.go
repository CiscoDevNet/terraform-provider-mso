package models

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var l3OutBGPAttributeTypes = map[string]attr.Type{
	"enabled": types.BoolType,
}

// L3OutBGPModel represents BGP enablement and future BGP-specific settings.
type L3OutBGPModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

func l3OutBGPModelFromTerraform(value types.Object) (L3OutBGPModel, bool) {
	if value.IsNull() || value.IsUnknown() {
		return L3OutBGPModel{}, false
	}
	enabled, _ := value.Attributes()["enabled"].(types.Bool)
	return L3OutBGPModel{Enabled: enabled}, true
}

func l3OutBGPObjectFromNDO(enabled bool) (types.Object, error) {
	value, diagnostics := types.ObjectValue(l3OutBGPAttributeTypes, map[string]attr.Value{
		"enabled": types.BoolValue(enabled),
	})
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutBGPAttributeTypes), fmt.Errorf("unable to represent BGP state: %s", diagnostics.Errors()[0].Detail())
	}
	return value, nil
}

func l3OutBGPResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether BGP is enabled on the L3Out.",
		},
	}
}

func l3OutBGPDataSourceSchema() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether BGP is enabled on the L3Out.",
		},
	}
}
