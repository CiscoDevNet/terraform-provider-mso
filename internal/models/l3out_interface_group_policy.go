package models

import (
	"context"
	"fmt"
	"maps"
	"reflect"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InterfaceGroupPolicyModel contains the fields shared by a standalone policy
// and a policy nested in the L3Out interface_groups map.
type InterfaceGroupPolicyModel struct {
	Description                types.String `tfsdk:"description"`
	InterfaceRoutingPolicyUUID types.String `tfsdk:"interface_routing_policy_uuid"`
	CustomQoSPolicyUUID        types.String `tfsdk:"custom_qos_policy_uuid"`
	QoSPriority                types.String `tfsdk:"qos_priority"`
	NetFlowMonitorUUIDs        types.Map    `tfsdk:"netflow_monitor_uuids"`
	BFD                        types.Object `tfsdk:"bfd"`
	BFDMultiHop                types.Object `tfsdk:"bfd_multi_hop"`
	OSPF                       types.Object `tfsdk:"ospf"`
}

// L3OutInterfaceGroupPolicyModel manages one named interfaceGroups entry.
type L3OutInterfaceGroupPolicyModel struct {
	ID         types.String `tfsdk:"id"`
	TemplateID types.String `tfsdk:"template_id"`
	L3OutUUID  types.String `tfsdk:"l3out_uuid"`
	Name       types.String `tfsdk:"name"`
	InterfaceGroupPolicyModel
}

type L3OutInterfaceGroupPolicyResourceIdentityModel struct {
	TemplateID types.String `tfsdk:"template_id"`
	L3OutUUID  types.String `tfsdk:"l3out_uuid"`
	Name       types.String `tfsdk:"name"`
}

func (identity L3OutInterfaceGroupPolicyResourceIdentityModel) ID() string {
	return fmt.Sprintf("%s/%s/%s", identity.TemplateID.ValueString(), identity.L3OutUUID.ValueString(), identity.Name.ValueString())
}

func L3OutInterfaceGroupPolicyResourceIdentitySchema() identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"template_id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "UUID of the L3Out template.",
			},
			"l3out_uuid": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "UUID of the parent L3Out.",
			},
			"name": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "Name of the interface group policy.",
			},
		},
	}
}

func (data L3OutInterfaceGroupPolicyModel) Path() ndoapi.Path {
	return NewL3OutPath(data.L3OutUUID.ValueString(), "").WithSteps(ndoapi.PathStep{
		Field: "interfaceGroups",
		Selector: &ndoapi.ObjectSelector{
			Keys: []ndoapi.ObjectIdentifier{
				{Field: "name", Value: data.Name.ValueString()},
			},
		},
	})
}

func (data *L3OutInterfaceGroupPolicyModel) SetFromNDOObject(ctx context.Context, templateID, l3outUUID string, object map[string]any, previous L3OutInterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) error {
	name, _, err := ndoapi.StringField(object, "name", ndoapi.RequiredNonEmptyField)
	if err != nil {
		return err
	}
	data.TemplateID = types.StringValue(templateID)
	data.L3OutUUID = types.StringValue(l3outUUID)
	data.Name = types.StringValue(name)
	data.ID = types.StringValue((L3OutInterfaceGroupPolicyResourceIdentityModel{
		TemplateID: data.TemplateID,
		L3OutUUID:  data.L3OutUUID,
		Name:       data.Name,
	}).ID())
	return data.InterfaceGroupPolicyModel.SetFromNDOObject(ctx, object, previous.InterfaceGroupPolicyModel, diagnostics)
}

func (data *InterfaceGroupPolicyModel) SetFromNDOObject(ctx context.Context, object map[string]any, previous InterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) error {
	description, descriptionExists, err := ndoapi.StringField(object, "description", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.Description = tfplugin.StringOrEmpty(description, descriptionExists, nil)
	interfaceRoutingPolicyRef, interfaceRoutingPolicyRefExists, err := ndoapi.StringField(object, "interfaceRoutingPolicyRef", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.InterfaceRoutingPolicyUUID = tfplugin.StringOrEmpty(interfaceRoutingPolicyRef, interfaceRoutingPolicyRefExists, nil)
	qosRef, qosRefExists, err := ndoapi.StringField(object, "qosRef", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.CustomQoSPolicyUUID = tfplugin.StringOrEmpty(qosRef, qosRefExists, nil)
	qosPriority, qosPriorityExists, err := ndoapi.StringField(object, "qosPriority", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.QoSPriority = tfplugin.StringOrNull(qosPriority, qosPriorityExists, nil)
	refs, _, err := ndoapi.MapField(object, "netFlowMonitorRefs", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	monitorUUIDs := make(map[string]string, len(refs))
	for _, kind := range []string{"ipv4", "ipv6", "ce", "unspecified"} {
		if value, exists, err := ndoapi.StringField(refs, kind, ndoapi.OptionalField); err != nil {
			return fmt.Errorf("NDO netFlowMonitorRefs: %w", err)
		} else if exists && value != "" {
			monitorUUIDs[kind] = value
		}
	}
	for kind := range refs {
		if kind != "ipv4" && kind != "ipv6" && kind != "ce" && kind != "unspecified" {
			return fmt.Errorf("NDO netFlowMonitorRefs has unsupported type %q", kind)
		}
	}
	data.NetFlowMonitorUUIDs = tfplugin.StringMap(ctx, monitorUUIDs, diagnostics)
	if diagnostics.HasError() {
		return nil
	}
	if data.BFD, err = interfaceGroupBFDFromNDO(object, "bfd", previous.BFD); err != nil {
		return err
	}
	if data.BFDMultiHop, err = interfaceGroupBFDFromNDO(object, "bfdMultiHop", previous.BFDMultiHop); err != nil {
		return err
	}
	data.OSPF, err = interfaceGroupOSPFFromNDO(object, previous.OSPF)
	return err
}

// DataSourceValue omits authentication keys, which Orchestration never returns.
func (data InterfaceGroupPolicyModel) DataSourceValue(ctx context.Context, diagnostics *diag.Diagnostics) InterfaceGroupPolicyModel {
	data.BFD = interfaceGroupProtocolDataSourceValue(ctx, data.BFD, diagnostics)
	data.BFDMultiHop = interfaceGroupProtocolDataSourceValue(ctx, data.BFDMultiHop, diagnostics)
	data.OSPF = interfaceGroupProtocolDataSourceValue(ctx, data.OSPF, diagnostics)
	return data
}

func interfaceGroupProtocolDataSourceValue(ctx context.Context, value types.Object, diagnostics *diag.Diagnostics) types.Object {
	attributeTypes := maps.Clone(value.AttributeTypes(ctx))
	delete(attributeTypes, "key")
	if value.IsNull() {
		return types.ObjectNull(attributeTypes)
	}
	if value.IsUnknown() {
		return types.ObjectUnknown(attributeTypes)
	}
	attributes := maps.Clone(value.Attributes())
	delete(attributes, "key")
	result, conversionDiagnostics := types.ObjectValue(attributeTypes, attributes)
	diagnostics.Append(conversionDiagnostics...)
	return result
}

func (data L3OutInterfaceGroupPolicyModel) DataSourceValue(ctx context.Context, diagnostics *diag.Diagnostics) L3OutInterfaceGroupPolicyModel {
	data.InterfaceGroupPolicyModel = data.InterfaceGroupPolicyModel.DataSourceValue(ctx, diagnostics)
	return data
}

func (data L3OutInterfaceGroupPolicyModel) ToPayload(ctx context.Context, configuration L3OutInterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) map[string]any {
	payload := data.InterfaceGroupPolicyModel.ToPayload(ctx, configuration.InterfaceGroupPolicyModel, diagnostics)
	payload["name"] = data.Name.ValueString()
	return payload
}

func (data InterfaceGroupPolicyModel) ToPayload(ctx context.Context, configuration InterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) map[string]any {
	payload := make(map[string]any)
	if description := tfplugin.KnownStringPointer(data.Description); description != nil && *description != "" {
		payload["description"] = *description
	}
	if interfaceRoutingPolicyRef := tfplugin.KnownStringPointer(data.InterfaceRoutingPolicyUUID); interfaceRoutingPolicyRef != nil && *interfaceRoutingPolicyRef != "" {
		payload["interfaceRoutingPolicyRef"] = *interfaceRoutingPolicyRef
	}
	if qosRef := tfplugin.KnownStringPointer(data.CustomQoSPolicyUUID); qosRef != nil && *qosRef != "" {
		payload["qosRef"] = *qosRef
	}
	if qosPriority := tfplugin.KnownStringPointer(data.QoSPriority); qosPriority != nil && *qosPriority != "" {
		payload["qosPriority"] = *qosPriority
	}
	if !configuration.NetFlowMonitorUUIDs.IsNull() {
		payload["netFlowMonitorRefs"] = interfaceGroupMonitorPayload(tfplugin.StringMapFromTerraform(ctx, data.NetFlowMonitorUUIDs, diagnostics))
	}
	if !configuration.BFD.IsNull() {
		if desired := interfaceGroupBFDPayload(data.BFD); desired != nil {
			payload["bfd"] = desired
		}
	}
	if !configuration.BFDMultiHop.IsNull() {
		if desired := interfaceGroupBFDPayload(data.BFDMultiHop); desired != nil {
			payload["bfdMultiHop"] = desired
		}
	}
	if !configuration.OSPF.IsNull() {
		if desired := interfaceGroupOSPFPayload(data.OSPF); desired != nil {
			payload["ospf"] = desired
		}
	}
	return payload
}

func interfaceGroupMonitorPayload(monitorUUIDs map[string]string) map[string]any {
	result := make(map[string]any, len(monitorUUIDs))
	for kind, uuid := range monitorUUIDs {
		result[kind] = uuid
	}
	return result
}

func (plan L3OutInterfaceGroupPolicyModel) AddPatchOperations(ctx context.Context, existing map[string]any, operations *ndoapi.PatchOperations, configuration, state L3OutInterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) error {
	return plan.InterfaceGroupPolicyModel.AddPatchOperations(ctx, existing, operations, configuration.InterfaceGroupPolicyModel, state.InterfaceGroupPolicyModel, diagnostics)
}

func (plan InterfaceGroupPolicyModel) AddPatchOperations(ctx context.Context, existing map[string]any, operations *ndoapi.PatchOperations, configuration, state InterfaceGroupPolicyModel, diagnostics *diag.Diagnostics) error {
	operations.SetString("description", tfplugin.KnownStringPointer(plan.Description), nil, true)
	operations.SetString("interfaceRoutingPolicyRef", tfplugin.KnownStringPointer(plan.InterfaceRoutingPolicyUUID), nil, true)
	operations.SetString("qosRef", tfplugin.KnownStringPointer(plan.CustomQoSPolicyUUID), nil, true)
	operations.SetString("qosPriority", tfplugin.KnownStringPointer(plan.QoSPriority), nil, false)
	if !configuration.NetFlowMonitorUUIDs.IsNull() {
		desired := interfaceGroupMonitorPayload(tfplugin.StringMapFromTerraform(ctx, plan.NetFlowMonitorUUIDs, diagnostics))
		current, _, err := ndoapi.MapField(existing, "netFlowMonitorRefs", ndoapi.OptionalField)
		if err != nil {
			return err
		}
		actual := make(map[string]any)
		for kind := range current {
			uuid, _, err := ndoapi.StringField(current, kind, ndoapi.RequiredField)
			if err != nil {
				return fmt.Errorf("NDO netFlowMonitorRefs: %w", err)
			}
			if uuid != "" {
				actual[kind] = uuid
			}
		}
		if !reflect.DeepEqual(actual, desired) {
			operations.Set("netFlowMonitorRefs", desired)
		}
	}
	if err := interfaceGroupBFDAddPatchOperation(operations, existing, "bfd", plan.BFD, configuration.BFD, state.BFD); err != nil {
		return err
	}
	if err := interfaceGroupBFDAddPatchOperation(operations, existing, "bfdMultiHop", plan.BFDMultiHop, configuration.BFDMultiHop, state.BFDMultiHop); err != nil {
		return err
	}
	if err := interfaceGroupOSPFAddPatchOperation(operations, existing, plan.OSPF, configuration.OSPF, state.OSPF); err != nil {
		return err
	}
	return nil
}

// setInterfaceGroupProtocolPatch preserves NDO's secret reference when the
// configured key is unchanged. BFD and OSPF decide authentication separately.
func setInterfaceGroupProtocolPatch(operations *ndoapi.PatchOperations, existing map[string]any, field string, desired map[string]any, configuration, state types.Object, isAuthenticated func(map[string]any) bool) error {
	current, _, err := ndoapi.MapField(existing, field, ndoapi.OptionalField)
	if err != nil {
		return err
	}
	configuredKey, _ := configuration.Attributes()["key"].(types.String)
	var previousKey types.String
	if !state.IsNull() && !state.IsUnknown() {
		previousKey, _ = state.Attributes()["key"].(types.String)
	}
	if configuredKey.IsNull() || configuredKey.IsUnknown() || configuredKey.Equal(previousKey) {
		delete(desired, "key")
	}
	merged := maps.Clone(current)
	if merged == nil {
		merged = make(map[string]any)
	}
	maps.Copy(merged, desired)
	if _, exists := desired["key"]; !exists {
		delete(merged, "key")
		if currentKey, exists := current["key"]; exists {
			merged["key"] = currentKey
		}
	}
	authenticated := isAuthenticated(merged)
	if !authenticated {
		delete(merged, "key")
		delete(merged, "keyID")
	}
	// NDO omits enabled=true and may retain a placeholder key ID without
	// authentication. JSON decoding also gives keyID a float64 type, while
	// Terraform supplies an int64. Those representations are equivalent.
	normalizedCurrent := maps.Clone(current)
	if normalizedCurrent != nil {
		if _, exists := normalizedCurrent["enabled"]; !exists {
			normalizedCurrent["enabled"] = true
		}
		if !authenticated {
			delete(normalizedCurrent, "key")
			delete(normalizedCurrent, "keyID")
		} else if keyID, exists, err := ndoapi.Int64Field(current, "keyID", ndoapi.OptionalField); err != nil {
			return fmt.Errorf("NDO %s: %w", field, err)
		} else if exists {
			normalizedCurrent["keyID"] = keyID
		}
		if reflect.DeepEqual(normalizedCurrent, merged) {
			return nil
		}
	}
	operations.Set(field, merged)
	return nil
}

func L3OutInterfaceGroupPolicyResourceSchema() schema.Schema {
	attributes := interfaceGroupPolicyResourceAttributes()
	attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Identifier of the interface group policy.",
	}
	attributes["template_id"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "UUID of the L3Out template.",
	}
	attributes["l3out_uuid"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "UUID of the parent L3Out.",
	}
	attributes["name"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "Name of the interface group policy.",
	}
	return schema.Schema{
		MarkdownDescription: "Manages one interface group policy on an L3Out.",
		Attributes:          attributes,
	}
}

func interfaceGroupPolicyResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"description": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "Description of the interface group policy.",
		},
		"interface_routing_policy_uuid": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "UUID of the referenced tenant interface routing policy.",
		},
		"custom_qos_policy_uuid": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseNonNullStateForUnknown(),
			},
			MarkdownDescription: "UUID of the referenced tenant custom QoS policy.",
		},
		"qos_priority": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseNonNullStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					"level1",
					"level2",
					"level3",
					"level4",
					"level5",
					"level6",
					"unspecified",
				),
			},
			MarkdownDescription: "QoS priority of the interface group.",
		},
		"netflow_monitor_uuids": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Map{
				mapplanmodifier.UseNonNullStateForUnknown(),
			},
			Validators: []validator.Map{
				mapvalidator.KeysAre(
					stringvalidator.OneOf(
						"ipv4",
						"ipv6",
						"ce",
						"unspecified",
					),
				),
				mapvalidator.NoNullValues(),
				mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
			MarkdownDescription: "NetFlow monitor UUIDs keyed by traffic type.",
		},
		"bfd": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseNonNullStateForUnknown(),
				tfplugin.UseConfiguredDisabledObject("enabled"),
			},
			Validators: []validator.Object{
				tfplugin.RequireAttributesWhenEnabled("enabled"),
				tfplugin.AllowAttributesOnlyWhenValue(
					"authentication_enabled",
					[]attr.Value{types.BoolValue(true)},
					"key",
					"key_id",
				),
			},
			Attributes:          interfaceGroupBFDResourceAttributes(),
			MarkdownDescription: "Single-hop BFD settings of the interface group.",
		},
		"bfd_multi_hop": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseNonNullStateForUnknown(),
				tfplugin.UseConfiguredDisabledObject("enabled"),
			},
			Validators: []validator.Object{
				tfplugin.RequireAttributesWhenEnabled("enabled"),
				tfplugin.AllowAttributesOnlyWhenValue(
					"authentication_enabled",
					[]attr.Value{types.BoolValue(true)},
					"key",
					"key_id",
				),
			},
			Attributes:          interfaceGroupBFDResourceAttributes(),
			MarkdownDescription: "Multi-hop BFD settings of the interface group.",
		},
		"ospf": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseNonNullStateForUnknown(),
				tfplugin.UseConfiguredDisabledObject("enabled"),
			},
			Validators: []validator.Object{
				tfplugin.RequireAttributesWhenEnabled("enabled"),
				tfplugin.AllowAttributesOnlyWhenValue(
					"authentication_type",
					[]attr.Value{
						types.StringValue("simple"),
						types.StringValue("md5"),
					},
					"key",
					"key_id",
				),
			},
			Attributes:          interfaceGroupOSPFResourceAttributes(),
			MarkdownDescription: "OSPF authentication settings of the interface group.",
		},
	}
}

func L3OutInterfaceGroupPolicyDataSourceSchema() datasourceschema.Schema {
	attributes := interfaceGroupPolicyDataSourceAttributes()
	attributes["id"] = datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Identifier of the interface group policy.",
	}
	attributes["template_id"] = datasourceschema.StringAttribute{
		Required:            true,
		MarkdownDescription: "UUID of the L3Out template.",
	}
	attributes["l3out_uuid"] = datasourceschema.StringAttribute{
		Required:            true,
		MarkdownDescription: "UUID of the parent L3Out.",
	}
	attributes["name"] = datasourceschema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Name of the interface group policy.",
	}
	return datasourceschema.Schema{
		MarkdownDescription: "Reads one interface group policy on an L3Out.",
		Attributes:          attributes,
	}
}

func interfaceGroupPolicyDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"description": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Description of the interface group policy.",
		},
		"interface_routing_policy_uuid": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "UUID of the referenced tenant interface routing policy.",
		},
		"custom_qos_policy_uuid": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "UUID of the referenced tenant custom QoS policy.",
		},
		"qos_priority": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "QoS priority of the interface group.",
		},
		"netflow_monitor_uuids": datasourceschema.MapAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			MarkdownDescription: "NetFlow monitor UUIDs keyed by traffic type.",
		},
		"bfd": datasourceschema.SingleNestedAttribute{
			Computed:            true,
			Attributes:          interfaceGroupBFDDataSourceAttributes(),
			MarkdownDescription: "Single-hop BFD settings of the interface group.",
		},
		"bfd_multi_hop": datasourceschema.SingleNestedAttribute{
			Computed:            true,
			Attributes:          interfaceGroupBFDDataSourceAttributes(),
			MarkdownDescription: "Multi-hop BFD settings of the interface group.",
		},
		"ospf": datasourceschema.SingleNestedAttribute{
			Computed:            true,
			Attributes:          interfaceGroupOSPFDataSourceAttributes(),
			MarkdownDescription: "OSPF authentication settings of the interface group.",
		},
	}
}
