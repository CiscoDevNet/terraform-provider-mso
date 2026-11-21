package models

import (
	"context"
	"fmt"
	"maps"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var l3OutRouteMapsAttributeTypes = map[string]attr.Type{
	"outbound_uuid":                           types.StringType,
	"interleak_uuid":                          types.StringType,
	"static_route_redistribution_uuid":        types.StringType,
	"connected_route_redistribution_uuid":     types.StringType,
	"attached_host_route_redistribution_uuid": types.StringType,
}

// L3OutRouteMapsModel contains the common L3Out route-map references. BGP
// dampening references are modeled with the BGP settings.
type L3OutRouteMapsModel struct {
	OutboundUUID                        types.String `tfsdk:"outbound_uuid"`
	InterleakUUID                       types.String `tfsdk:"interleak_uuid"`
	StaticRouteRedistributionUUID       types.String `tfsdk:"static_route_redistribution_uuid"`
	ConnectedRouteRedistributionUUID    types.String `tfsdk:"connected_route_redistribution_uuid"`
	AttachedHostRouteRedistributionUUID types.String `tfsdk:"attached_host_route_redistribution_uuid"`
}

type l3OutRouteMapReference struct {
	terraformName string
	apiName       string
	value         types.String
}

func (model L3OutRouteMapsModel) references() []l3OutRouteMapReference {
	return []l3OutRouteMapReference{
		{terraformName: "outbound_uuid", apiName: "exportRouteMapRef", value: model.OutboundUUID},
		{terraformName: "interleak_uuid", apiName: "interleakRef", value: model.InterleakUUID},
		{
			terraformName: "static_route_redistribution_uuid",
			apiName:       "staticRouteRedistRef",
			value:         model.StaticRouteRedistributionUUID,
		},
		{
			terraformName: "connected_route_redistribution_uuid",
			apiName:       "connectedRouteRedistRef",
			value:         model.ConnectedRouteRedistributionUUID,
		},
		{
			terraformName: "attached_host_route_redistribution_uuid",
			apiName:       "attachedHostRouteRedistRef",
			value:         model.AttachedHostRouteRedistributionUUID,
		},
	}
}

func l3OutRouteMapsModelFromTerraform(ctx context.Context, value types.Object) (L3OutRouteMapsModel, bool, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return L3OutRouteMapsModel{}, false, nil
	}

	var model L3OutRouteMapsModel
	diagnostics := value.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, true, diagnostics
}

// l3OutRouteMapObjectsFromNDO reads the common and BGP-specific route-map
// references together because they share the advancedRouteMapRefs API object.
func l3OutRouteMapObjectsFromNDO(ctx context.Context, object map[string]any, bgpEnabled bool, previousRouteMaps, previousBGP types.Object) (types.Object, types.Object, error) {
	outbound, outboundExists, err := ndoapi.StringField(object, "exportRouteMapRef", ndoapi.OptionalField)
	if err != nil {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			err
	}
	advanced, _, err := ndoapi.MapField(object, "advancedRouteMapRefs", ndoapi.OptionalField)
	if err != nil {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			err
	}
	previousModel, _, diagnostics := l3OutRouteMapsModelFromTerraform(ctx, previousRouteMaps)
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("unable to decode previous L3Out route maps: %s", diagnostics.Errors()[0].Detail())
	}
	values := map[string]attr.Value{
		"outbound_uuid": stringFieldOrPreviousEmpty(outbound, outboundExists, previousModel.OutboundUUID),
	}
	for _, reference := range (L3OutRouteMapsModel{}).references()[1:] {
		value, exists, fieldErr := ndoapi.StringField(advanced, reference.apiName, ndoapi.OptionalField)
		if fieldErr != nil {
			return types.ObjectNull(l3OutRouteMapsAttributeTypes),
				types.ObjectNull(l3OutBGPAttributeTypes),
				fmt.Errorf("NDO advancedRouteMapRefs: %w", fieldErr)
		}
		values[reference.terraformName] = stringFieldOrPreviousEmpty(
			value,
			exists,
			previousRouteMapValue(previousModel, reference.terraformName),
		)
	}
	result, diagnostics := types.ObjectValue(l3OutRouteMapsAttributeTypes, values)
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("unable to represent L3Out route maps: %s", diagnostics.Errors()[0].Detail())
	}
	routeDampeningIPv4UUID, routeDampeningIPv4Exists, err := ndoapi.StringField(
		advanced,
		"routeDampeningV4Ref",
		ndoapi.OptionalField,
	)
	if err != nil {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("NDO advancedRouteMapRefs: %w", err)
	}
	routeDampeningIPv6UUID, routeDampeningIPv6Exists, err := ndoapi.StringField(
		advanced,
		"routeDampeningV6Ref",
		ndoapi.OptionalField,
	)
	if err != nil {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("NDO advancedRouteMapRefs: %w", err)
	}
	previousBGPModel, _, diagnostics := l3OutBGPModelFromTerraform(ctx, previousBGP)
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutRouteMapsAttributeTypes),
			types.ObjectNull(l3OutBGPAttributeTypes),
			fmt.Errorf("unable to decode previous L3Out BGP settings: %s", diagnostics.Errors()[0].Detail())
	}
	bgp, err := l3OutBGPObjectFromNDO(
		bgpEnabled,
		stringFieldOrPreviousEmpty(
			routeDampeningIPv4UUID,
			routeDampeningIPv4Exists,
			previousBGPModel.RouteDampeningIPv4UUID,
		),
		stringFieldOrPreviousEmpty(
			routeDampeningIPv6UUID,
			routeDampeningIPv6Exists,
			previousBGPModel.RouteDampeningIPv6UUID,
		),
	)
	return result, bgp, err
}

func previousRouteMapValue(model L3OutRouteMapsModel, name string) types.String {
	for _, reference := range model.references() {
		if reference.terraformName == name {
			return reference.value
		}
	}
	return types.StringNull()
}

func stringFieldOrPreviousEmpty(value string, exists bool, previous types.String) types.String {
	if exists {
		return types.StringValue(value)
	}
	if previous.IsNull() || previous.IsUnknown() || previous.ValueString() != "" {
		return types.StringNull()
	}
	return types.StringValue("")
}

// l3OutRouteMapPayload builds a complete payload fragment without changing
// caller-owned maps. Unconfigured models contain only null reference values.
func l3OutRouteMapPayload(routeMaps L3OutRouteMapsModel, bgp L3OutBGPModel) map[string]any {
	payload := make(map[string]any)
	advanced := make(map[string]any)
	for _, reference := range routeMaps.references() {
		value := tfplugin.KnownStringPointer(reference.value)
		if value == nil || *value == "" {
			continue
		}
		if reference.apiName == "exportRouteMapRef" {
			payload[reference.apiName] = *value
			continue
		}
		advanced[reference.apiName] = *value
	}
	if value := tfplugin.KnownStringPointer(bgp.RouteDampeningIPv4UUID); value != nil && *value != "" {
		advanced["routeDampeningV4Ref"] = *value
	}
	if value := tfplugin.KnownStringPointer(bgp.RouteDampeningIPv6UUID); value != nil && *value != "" {
		advanced["routeDampeningV6Ref"] = *value
	}
	if len(advanced) > 0 {
		payload["advancedRouteMapRefs"] = advanced
	}
	return payload
}

func l3OutRouteMapsResourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"outbound_uuid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the optional outbound route-map policy.",
			},
			"interleak_uuid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the interleak route-map policy.",
			},
			"static_route_redistribution_uuid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the static-route redistribution route-map policy.",
			},
			"connected_route_redistribution_uuid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the connected-route redistribution route-map policy.",
			},
			"attached_host_route_redistribution_uuid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the attached-host-route redistribution route-map policy.",
			},
		},
		MarkdownDescription: "Route-map policies associated with the L3Out. Omitted fields preserve their current values; set an individual UUID to an empty string to clear that reference.",
	}
}

func l3OutRouteMapsDataSourceSchema() datasourceschema.SingleNestedAttribute {
	return datasourceschema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]datasourceschema.Attribute{
			"outbound_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the outbound route-map policy, if configured.",
			},
			"interleak_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the interleak route-map policy, if configured.",
			},
			"static_route_redistribution_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the static-route redistribution route-map policy, if configured.",
			},
			"connected_route_redistribution_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the connected-route redistribution route-map policy, if configured.",
			},
			"attached_host_route_redistribution_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the attached-host-route redistribution route-map policy, if configured.",
			},
		},
		MarkdownDescription: "Route-map policy references configured on the L3Out.",
	}
}

func setL3OutRouteMapOperations(operations *ndoapi.PatchOperations, existing map[string]any, routeMaps L3OutRouteMapsModel, bgp L3OutBGPModel) error {
	for _, reference := range routeMaps.references() {
		if reference.apiName == "exportRouteMapRef" {
			operations.SetString(reference.apiName, tfplugin.KnownStringPointer(reference.value), nil, true)
		}
	}

	advanced, _, err := ndoapi.MapField(existing, "advancedRouteMapRefs", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO L3Out: %w", err)
	}
	updated := maps.Clone(advanced)
	if updated == nil {
		updated = make(map[string]any)
	}
	changed := false
	for _, reference := range routeMaps.references()[1:] {
		changed = applyRouteMapReference(
			updated,
			reference.apiName,
			tfplugin.KnownStringPointer(reference.value),
		) || changed
	}
	changed = applyRouteMapReference(
		updated,
		"routeDampeningV4Ref",
		tfplugin.KnownStringPointer(bgp.RouteDampeningIPv4UUID),
	) || changed
	changed = applyRouteMapReference(
		updated,
		"routeDampeningV6Ref",
		tfplugin.KnownStringPointer(bgp.RouteDampeningIPv6UUID),
	) || changed
	if changed {
		if len(updated) == 0 {
			operations.Remove("advancedRouteMapRefs")
		} else {
			operations.Set("advancedRouteMapRefs", updated)
		}
	}
	return nil
}

func applyRouteMapReference(advanced map[string]any, name string, value *string) bool {
	if value == nil {
		return false
	}
	if *value == "" {
		if _, exists := advanced[name]; !exists {
			return false
		}
		delete(advanced, name)
		return true
	}
	if current, exists := advanced[name]; exists && current == *value {
		return false
	}
	advanced[name] = *value
	return true
}
