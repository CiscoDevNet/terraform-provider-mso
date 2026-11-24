package models

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var l3OutBGPAttributeTypes = map[string]attr.Type{
	"enabled":                   types.BoolType,
	"route_dampening_ipv4_uuid": types.StringType,
	"route_dampening_ipv6_uuid": types.StringType,
}

// L3OutBGPModel represents BGP enablement and BGP-specific route settings.
type L3OutBGPModel struct {
	Enabled                types.Bool   `tfsdk:"enabled"`
	RouteDampeningIPv4UUID types.String `tfsdk:"route_dampening_ipv4_uuid"`
	RouteDampeningIPv6UUID types.String `tfsdk:"route_dampening_ipv6_uuid"`
}

func l3OutBGPModelFromTerraform(ctx context.Context, value types.Object) (L3OutBGPModel, bool, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return L3OutBGPModel{}, false, nil
	}

	var model L3OutBGPModel
	diagnostics := value.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, true, diagnostics
}

func l3OutBGPObjectFromNDO(enabled bool, routeDampeningIPv4UUID, routeDampeningIPv6UUID types.String) (types.Object, error) {
	value, diagnostics := types.ObjectValue(l3OutBGPAttributeTypes, map[string]attr.Value{
		"enabled":                   types.BoolValue(enabled),
		"route_dampening_ipv4_uuid": routeDampeningIPv4UUID,
		"route_dampening_ipv6_uuid": routeDampeningIPv6UUID,
	})
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("unable to represent BGP state: %s", diagnostics.Errors()[0].Detail())
	}
	return value, nil
}

func l3OutBGPResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether BGP is enabled on the L3Out. Disabling BGP removes both route-dampening policy associations.",
		},
		"route_dampening_ipv4_uuid": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "UUID of the IPv4 route-dampening route-map policy. Route dampening is supported only when BGP is enabled on the L3Out. Omit while BGP remains enabled to preserve the association, or set to an empty string to clear it.",
		},
		"route_dampening_ipv6_uuid": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			MarkdownDescription: "UUID of the IPv6 route-dampening route-map policy. Route dampening is supported only when BGP is enabled on the L3Out. Omit while BGP remains enabled to preserve the association, or set to an empty string to clear it.",
		},
	}
}

func l3OutBGPDataSourceSchema() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether BGP is enabled on the L3Out.",
		},
		"route_dampening_ipv4_uuid": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "UUID of the IPv4 route-dampening route-map policy, if configured.",
		},
		"route_dampening_ipv6_uuid": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "UUID of the IPv6 route-dampening route-map policy, if configured.",
		},
	}
}
