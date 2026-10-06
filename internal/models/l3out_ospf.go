package models

import (
	"fmt"
	"maps"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var l3OutOSPFAttributeTypes = map[string]attr.Type{
	"enabled":                        types.BoolType,
	"area_id":                        types.StringType,
	"area_type":                      types.StringType,
	"cost":                           types.Int64Type,
	"originate_default_route_always": types.BoolType,
	"send_redistributed_lsa":         types.BoolType,
	"originate_summary_lsa":          types.BoolType,
	"suppress_forwarding_address_in_translated_lsa": types.BoolType,
}

// L3OutOSPFModel represents the OSPF block in the Terraform model and its
// corresponding fields in an NDO L3Out payload.
type L3OutOSPFModel struct {
	Enabled                                  types.Bool   `tfsdk:"enabled"`
	AreaID                                   types.String `tfsdk:"area_id"`
	AreaType                                 types.String `tfsdk:"area_type"`
	Cost                                     types.Int64  `tfsdk:"cost"`
	OriginateDefaultRouteAlways              types.Bool   `tfsdk:"originate_default_route_always"`
	SendRedistributedLSA                     types.Bool   `tfsdk:"send_redistributed_lsa"`
	OriginateSummaryLSA                      types.Bool   `tfsdk:"originate_summary_lsa"`
	SuppressForwardingAddressInTranslatedLSA types.Bool   `tfsdk:"suppress_forwarding_address_in_translated_lsa"`
}

func (model L3OutOSPFModel) toPayload() map[string]any {
	payload := make(map[string]any)
	areaConfig := make(map[string]any)
	if cost := tfplugin.KnownInt64Pointer(model.Cost); cost != nil {
		areaConfig["cost"] = *cost
	}
	if areaID := tfplugin.KnownStringPointer(model.AreaID); areaID != nil {
		areaConfig["id"] = *areaID
	}
	if areaType := tfplugin.KnownStringPointer(model.AreaType); areaType != nil {
		areaConfig["areaType"] = *areaType
	}
	control := make(map[string]any)
	if sendRedistributedLSA := tfplugin.KnownBoolPointer(model.SendRedistributedLSA); sendRedistributedLSA != nil {
		control["redistribute"] = *sendRedistributedLSA
	}
	if originateSummaryLSA := tfplugin.KnownBoolPointer(model.OriginateSummaryLSA); originateSummaryLSA != nil {
		control["originate"] = *originateSummaryLSA
	}
	if suppressForwardingAddressInTranslatedLSA := tfplugin.KnownBoolPointer(model.SuppressForwardingAddressInTranslatedLSA); suppressForwardingAddressInTranslatedLSA != nil {
		control["suppressFA"] = *suppressForwardingAddressInTranslatedLSA
	}
	if len(control) > 0 {
		areaConfig["control"] = control
	}
	payload["ospfAreaConfig"] = areaConfig
	return payload
}

func setL3OutOSPFOperations(operations *ndoapi.PatchOperations, existing map[string]any, model L3OutOSPFModel) error {
	if !model.Enabled.ValueBool() {
		operations.Remove("ospfAreaConfig")
		return nil
	}

	current, _, err := ndoapi.MapField(existing, "ospfAreaConfig", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	areaConfig := maps.Clone(current)
	if areaConfig == nil {
		areaConfig = make(map[string]any)
	}
	if cost, exists, err := ndoapi.Int64Field(areaConfig, "cost", ndoapi.OptionalField); err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	} else if exists {
		areaConfig["cost"] = cost
	}
	if cost := tfplugin.KnownInt64Pointer(model.Cost); cost != nil {
		areaConfig["cost"] = *cost
	}
	if areaID := tfplugin.KnownStringPointer(model.AreaID); areaID != nil {
		areaConfig["id"] = *areaID
	}
	if areaType := tfplugin.KnownStringPointer(model.AreaType); areaType != nil {
		areaConfig["areaType"] = *areaType
	}
	currentControl, _, err := ndoapi.MapField(areaConfig, "control", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	}
	control := maps.Clone(currentControl)
	if control == nil {
		control = make(map[string]any)
	}
	if sendRedistributedLSA := tfplugin.KnownBoolPointer(model.SendRedistributedLSA); sendRedistributedLSA != nil {
		control["redistribute"] = *sendRedistributedLSA
	}
	if originateSummaryLSA := tfplugin.KnownBoolPointer(model.OriginateSummaryLSA); originateSummaryLSA != nil {
		control["originate"] = *originateSummaryLSA
	}
	if suppressForwardingAddressInTranslatedLSA := tfplugin.KnownBoolPointer(model.SuppressForwardingAddressInTranslatedLSA); suppressForwardingAddressInTranslatedLSA != nil {
		control["suppressFA"] = *suppressForwardingAddressInTranslatedLSA
	}
	if len(control) > 0 {
		areaConfig["control"] = control
	}
	operations.Set("ospfAreaConfig", areaConfig)
	return nil
}

func (model *L3OutOSPFModel) setFromTerraformObject(value types.Object) {
	attributes := value.Attributes()
	model.Enabled, _ = attributes["enabled"].(types.Bool)
	model.AreaID, _ = attributes["area_id"].(types.String)
	model.AreaType, _ = attributes["area_type"].(types.String)
	model.Cost, _ = attributes["cost"].(types.Int64)
	model.OriginateDefaultRouteAlways, _ = attributes["originate_default_route_always"].(types.Bool)
	model.SendRedistributedLSA, _ = attributes["send_redistributed_lsa"].(types.Bool)
	model.OriginateSummaryLSA, _ = attributes["originate_summary_lsa"].(types.Bool)
	model.SuppressForwardingAddressInTranslatedLSA, _ = attributes["suppress_forwarding_address_in_translated_lsa"].(types.Bool)
}

func l3OutOSPFModelFromTerraform(value types.Object) (L3OutOSPFModel, bool) {
	if value.IsNull() || value.IsUnknown() {
		return L3OutOSPFModel{}, false
	}

	var model L3OutOSPFModel
	model.setFromTerraformObject(value)
	return model, true
}

func (model *L3OutOSPFModel) setFromNDOObject(object map[string]any) error {
	*model = L3OutOSPFModel{
		Enabled:                                  types.BoolValue(true),
		AreaID:                                   types.StringNull(),
		AreaType:                                 types.StringNull(),
		Cost:                                     types.Int64Null(),
		OriginateDefaultRouteAlways:              types.BoolNull(),
		SendRedistributedLSA:                     types.BoolNull(),
		OriginateSummaryLSA:                      types.BoolNull(),
		SuppressForwardingAddressInTranslatedLSA: types.BoolNull(),
	}

	areaConfig, exists, err := ndoapi.MapField(object, "ospfAreaConfig", ndoapi.OptionalField)
	if err != nil || !exists {
		return err
	}
	areaID, exists, err := ndoapi.StringField(areaConfig, "id", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	}
	if exists {
		model.AreaID = types.StringValue(areaID)
	}
	areaType, exists, err := ndoapi.StringField(areaConfig, "areaType", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	}
	if exists {
		model.AreaType = types.StringValue(areaType)
	}
	cost, exists, err := ndoapi.Int64Field(areaConfig, "cost", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	}
	if exists {
		model.Cost = types.Int64Value(cost)
	}
	control, exists, err := ndoapi.MapField(areaConfig, "control", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig: %w", err)
	}
	if !exists {
		return nil
	}
	redistribute, exists, err := ndoapi.BoolField(control, "redistribute", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig.control: %w", err)
	}
	if exists {
		model.SendRedistributedLSA = types.BoolValue(redistribute)
	}
	originate, exists, err := ndoapi.BoolField(control, "originate", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig.control: %w", err)
	}
	if exists {
		model.OriginateSummaryLSA = types.BoolValue(originate)
	}
	suppressFA, exists, err := ndoapi.BoolField(control, "suppressFA", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO ospfAreaConfig.control: %w", err)
	}
	if exists {
		model.SuppressForwardingAddressInTranslatedLSA = types.BoolValue(suppressFA)
	}

	return nil
}

func (model L3OutOSPFModel) toTerraformObject() (types.Object, error) {
	value, diagnostics := types.ObjectValue(l3OutOSPFAttributeTypes, map[string]attr.Value{
		"enabled":                        model.Enabled,
		"area_id":                        model.AreaID,
		"area_type":                      model.AreaType,
		"cost":                           model.Cost,
		"originate_default_route_always": model.OriginateDefaultRouteAlways,
		"send_redistributed_lsa":         model.SendRedistributedLSA,
		"originate_summary_lsa":          model.OriginateSummaryLSA,
		"suppress_forwarding_address_in_translated_lsa": model.SuppressForwardingAddressInTranslatedLSA,
	})
	if diagnostics.HasError() {
		return types.ObjectNull(l3OutOSPFAttributeTypes), fmt.Errorf("unable to represent OSPF state: %s", diagnostics.Errors()[0].Detail())
	}
	return value, nil
}

func l3OutOSPFObjectFromNDO(object map[string]any, always types.Bool, ospfEnabled bool) (types.Object, error) {
	if !ospfEnabled {
		return (L3OutOSPFModel{
			Enabled:                                  types.BoolValue(false),
			AreaID:                                   types.StringNull(),
			AreaType:                                 types.StringNull(),
			Cost:                                     types.Int64Null(),
			OriginateDefaultRouteAlways:              types.BoolNull(),
			SendRedistributedLSA:                     types.BoolNull(),
			OriginateSummaryLSA:                      types.BoolNull(),
			SuppressForwardingAddressInTranslatedLSA: types.BoolNull(),
		}).toTerraformObject()
	}

	var model L3OutOSPFModel
	if err := model.setFromNDOObject(object); err != nil {
		return types.ObjectNull(l3OutOSPFAttributeTypes), err
	}
	if always.IsNull() {
		always = types.BoolValue(false)
	}
	model.OriginateDefaultRouteAlways = always
	return model.toTerraformObject()
}

func l3OutOSPFResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether OSPF is enabled on the L3Out.",
		},
		"area_id": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "OSPF area ID.",
		},
		"area_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"regular",
					"stub",
					"nssa",
				),
			},
			MarkdownDescription: "OSPF area type.",
		},
		"cost": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				tfplugin.UseNonNullStateForUnknownInt64(),
			},
			MarkdownDescription: "OSPF area cost.",
		},
		"originate_default_route_always": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				tfplugin.SetBoolWhenStringEquals(
					path.MatchRoot("originate_default_route"),
					"",
					false,
				),
				tfplugin.NullBoolWhenParentDisabledOrAbsent(
					path.MatchRoot("ospf"),
					"enabled",
				),
				boolplanmodifier.UseNonNullStateForUnknown(),
			},
			Validators: []validator.Bool{
				tfplugin.RequireValueWhenBool(
					true,
					path.MatchRoot("originate_default_route"),
					types.StringValue("only"),
					types.StringValue("in_addition"),
				),
			},
			MarkdownDescription: "Whether OSPF originates the default route when it is absent from the routing table.",
		},
		"send_redistributed_lsa": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "Controls sending redistributed LSAs into the OSPF area.",
		},
		"originate_summary_lsa": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "Controls OSPF summary-LSA origination.",
		},
		"suppress_forwarding_address_in_translated_lsa": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "Controls suppression of forwarding-address-translated LSAs.",
		},
	}
}

func l3OutOSPFDataSourceSchema() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether OSPF is enabled on the L3Out.",
		},
		"area_id": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "OSPF area ID.",
		},
		"area_type": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "OSPF area type.",
		},
		"cost": datasourceschema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "OSPF area cost.",
		},
		"originate_default_route_always": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether OSPF always originates the default route.",
		},
		"send_redistributed_lsa": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether redistributed LSAs are sent into the OSPF area.",
		},
		"originate_summary_lsa": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether OSPF summary LSAs are originated.",
		},
		"suppress_forwarding_address_in_translated_lsa": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether forwarding addresses in translated LSAs are suppressed.",
		},
	}
}
