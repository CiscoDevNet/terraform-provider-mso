package models

import (
	"context"
	"fmt"
	"maps"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var targetDSCPTranslation = ndoapi.NewTranslationMap(map[string]string{
	"expedited_forwarding": "expeditedForwarding",
	"voice_admit":          "voiceAdmit",
})

// NewL3OutPath declares the L3Out location and lookup identity. A
// known UUID is authoritative; name is used only when UUID is absent.
func NewL3OutPath(uuid, name string) ndoapi.Path {
	return ndoapi.NewPath(
		ndoapi.PathStep{Field: "l3outTemplate"},
		ndoapi.PathStep{
			Field: "l3outs",
			Selector: &ndoapi.ObjectSelector{
				UUID: ndoapi.ObjectIdentifier{Field: "uuid", Value: uuid},
				Keys: []ndoapi.ObjectIdentifier{{Field: "name", Value: name}},
			},
		},
	)
}

// L3OutModel is the aggregate L3Out Terraform model. The routing protocol is
// derived from the OSPF and BGP objects.
type L3OutModel struct {
	ID                        types.String `tfsdk:"id"`
	TemplateID                types.String `tfsdk:"template_id"`
	UUID                      types.String `tfsdk:"uuid"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	VRFUUID                   types.String `tfsdk:"vrf_uuid"`
	L3Domain                  types.String `tfsdk:"l3_domain"`
	TargetDSCP                types.String `tfsdk:"target_dscp"`
	PIMEnabled                types.Bool   `tfsdk:"pim_enabled"`
	ImportRouteControlEnabled types.Bool   `tfsdk:"import_route_control_enabled"`
	BGP                       types.Object `tfsdk:"bgp"`
	OriginateDefaultRoute     types.String `tfsdk:"originate_default_route"`
	OSPF                      types.Object `tfsdk:"ospf"`
	Annotations               types.Map    `tfsdk:"annotations"`
	InterfaceGroups           types.Map    `tfsdk:"interface_groups"`
}

// L3OutResourceIdentityModel identifies an L3Out within its template.
type L3OutResourceIdentityModel struct {
	TemplateID types.String `tfsdk:"template_id"`
	UUID       types.String `tfsdk:"uuid"`
}

// ID returns the string import identifier for this L3Out identity.
func (identity L3OutResourceIdentityModel) ID() string {
	return fmt.Sprintf("%s/%s", identity.TemplateID.ValueString(), identity.UUID.ValueString())
}

func L3OutResourceIdentitySchema() identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"template_id": identityschema.StringAttribute{
				Description:       "The UUID of the L3Out template containing the L3Out.",
				RequiredForImport: true,
			},
			"uuid": identityschema.StringAttribute{
				Description:       "The UUID of the L3Out within the template.",
				RequiredForImport: true,
			},
		},
	}
}

func (data L3OutModel) Path() ndoapi.Path {
	return NewL3OutPath(data.UUID.ValueString(), data.Name.ValueString())
}

// SetFromNDOObject translates an NDO L3Out object into Terraform state.
// previousL3OutModel is the plan after Create/Update or the prior state during
// Read. NDO returns only references for interface group authentication keys,
// so their configured values must be carried forward from that model.
func (data *L3OutModel) SetFromNDOObject(ctx context.Context, templateID string, object map[string]any, previousL3OutModel L3OutModel) error {
	// Terraform validates configuration, but these fields come from NDO and
	// must be valid before they are written to state.
	uuid, _, err := ndoapi.StringField(object, "uuid", ndoapi.RequiredNonEmptyField)
	if err != nil {
		return err
	}
	name, _, err := ndoapi.StringField(object, "name", ndoapi.RequiredNonEmptyField)
	if err != nil {
		return err
	}
	vrfRef, _, err := ndoapi.StringField(object, "vrfRef", ndoapi.RequiredNonEmptyField)
	if err != nil {
		return err
	}
	routingProtocolValue, _, err := ndoapi.StringField(object, "routingProtocol", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	routingProtocol, err := newL3OutRoutingProtocolState(routingProtocolValue)
	if err != nil {
		return err
	}
	data.TemplateID = types.StringValue(templateID)
	data.UUID = types.StringValue(uuid)
	identity := L3OutResourceIdentityModel{TemplateID: data.TemplateID, UUID: data.UUID}
	data.ID = types.StringValue(identity.ID())
	data.Name = types.StringValue(name)
	data.VRFUUID = types.StringValue(vrfRef)
	description, descriptionExists, err := ndoapi.StringField(object, "description", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.Description = tfplugin.StringOrEmpty(description, descriptionExists, nil)
	l3Domain, l3DomainExists, err := ndoapi.StringField(object, "l3domain", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.L3Domain = tfplugin.StringOrEmpty(l3Domain, l3DomainExists, nil)
	targetDSCP, targetDSCPExists, err := ndoapi.StringField(object, "targetDscp", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.TargetDSCP = tfplugin.StringOrNull(targetDSCP, targetDSCPExists, targetDSCPTranslation.ToSchema)
	pim, pimExists, err := ndoapi.BoolField(object, "pim", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.PIMEnabled = tfplugin.BoolOrNull(pim, pimExists)
	importControl, importControlExists, err := ndoapi.BoolField(object, "importRouteControl", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	data.ImportRouteControlEnabled = tfplugin.BoolOrNull(importControl, importControlExists)
	data.BGP, err = l3OutBGPObjectFromNDO(routingProtocol.bgp)
	if err != nil {
		return err
	}
	var defaultRoute L3OutDefaultRouteModel
	err = defaultRoute.SetFromNDOObject(object)
	if err != nil {
		return err
	}
	data.OriginateDefaultRoute = defaultRoute.Mode
	data.OSPF, err = l3OutOSPFObjectFromNDO(object, defaultRoute.Always, routingProtocol.ospf)
	if err != nil {
		return err
	}
	var annotations AnnotationsModel
	if err := annotations.SetFromNDOObject(object); err != nil {
		return err
	}
	var diagnostics diag.Diagnostics
	data.Annotations = annotations.TerraformValue(ctx, &diagnostics)
	if diagnostics.HasError() {
		return fmt.Errorf("failed to convert L3Out annotations to Terraform state: %s", diagnostics.Errors()[0].Detail())
	}
	previousGroups := InterfaceGroupPoliciesFromTerraform(ctx, previousL3OutModel.InterfaceGroups, &diagnostics)
	if diagnostics.HasError() {
		return fmt.Errorf("failed to read previous L3Out interface groups: %s", diagnostics.Errors()[0].Detail())
	}
	var groups InterfaceGroupPoliciesModel
	if err := groups.SetFromNDOObject(ctx, object, previousGroups, &diagnostics); err != nil {
		return err
	}
	if diagnostics.HasError() {
		return fmt.Errorf("failed to read L3Out interface groups: %s", diagnostics.Errors()[0].Detail())
	}
	data.InterfaceGroups = groups.TerraformValue(ctx, &diagnostics)
	if diagnostics.HasError() {
		return fmt.Errorf("failed to convert L3Out interface groups to Terraform state: %s", diagnostics.Errors()[0].Detail())
	}
	return nil
}

// DataSourceValue removes fields that cannot be read from Orchestration.
func (data L3OutModel) DataSourceValue(ctx context.Context, diagnostics *diag.Diagnostics) L3OutModel {
	groups := InterfaceGroupPoliciesFromTerraform(ctx, data.InterfaceGroups, diagnostics)
	if diagnostics.HasError() {
		return data
	}
	data.InterfaceGroups = groups.DataSourceValue(ctx, diagnostics)
	return data
}

// ToPayload builds the complete create payload from planned values. The
// configuration identifies explicitly managed blocks; computed child values
// in the plan must not turn an explicit empty block into an enabled one.
func (data L3OutModel) ToPayload(ctx context.Context, configuration L3OutModel, diagnostics *diag.Diagnostics) map[string]any {
	payload := map[string]any{
		"name":            data.Name.ValueString(),
		"routingProtocol": "none",
		"vrfRef":          data.VRFUUID.ValueString(),
	}
	if description := tfplugin.KnownStringPointer(data.Description); description != nil {
		payload["description"] = *description
	}
	if l3Domain := tfplugin.KnownStringPointer(data.L3Domain); l3Domain != nil {
		payload["l3domain"] = *l3Domain
	}
	if targetDSCP := tfplugin.KnownStringPointer(data.TargetDSCP); targetDSCP != nil {
		payload["targetDscp"] = targetDSCPTranslation.ToAPI(*targetDSCP)
	}
	if pim := tfplugin.KnownBoolPointer(data.PIMEnabled); pim != nil {
		payload["pim"] = *pim
	}
	if importControl := tfplugin.KnownBoolPointer(data.ImportRouteControlEnabled); importControl != nil {
		payload["importRouteControl"] = *importControl
	}
	routingProtocolChange := l3OutRoutingProtocolChange{}
	ospf, ospfConfigured := l3OutOSPFModelFromTerraform(configuration.OSPF)
	if ospfConfigured && ospf.Enabled.ValueBool() {
		ospfEnabled := true
		routingProtocolChange.ospf = &ospfEnabled
		maps.Copy(payload, ospf.toPayload())
	}
	if _, configured := l3OutBGPModelFromTerraform(configuration.BGP); configured {
		bgp, _ := l3OutBGPModelFromTerraform(data.BGP)
		routingProtocolChange.bgp = tfplugin.KnownBoolPointer(bgp.Enabled)
	}
	payload["routingProtocol"] = l3OutRoutingProtocolValue(routingProtocolChange)
	defaultRoute := L3OutDefaultRouteModel{Mode: configuration.OriginateDefaultRoute}
	if ospfConfigured && ospf.Enabled.ValueBool() {
		defaultRoute.Always = ospf.OriginateDefaultRouteAlways
	}
	if defaultRouteLeak := defaultRoute.ToPayload(); len(defaultRouteLeak) > 0 {
		payload["defaultRouteLeak"] = defaultRouteLeak
	}
	if !configuration.Annotations.IsNull() {
		annotations := AnnotationsModelFromTerraform(ctx, data.Annotations, diagnostics)
		if annotations != nil {
			payload["tagAnnotations"] = annotations.ToPayload()
		}
	}
	if !configuration.InterfaceGroups.IsNull() {
		groups := InterfaceGroupPoliciesFromTerraform(ctx, data.InterfaceGroups, diagnostics)
		configured := InterfaceGroupPoliciesFromTerraform(ctx, configuration.InterfaceGroups, diagnostics)
		if !diagnostics.HasError() && groups != nil {
			payload["interfaceGroups"] = groups.ToPayload(ctx, configured, diagnostics)
		}
	}
	return payload
}

// L3OutResourceSchema returns the aggregate L3Out resource schema.
func L3OutResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an L3Out in a Nexus Dashboard Orchestration template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Stable identifier composed of the template UUID and L3Out UUID.",
			},
			"template_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "UUID of the existing L3Out template.",
			},
			"uuid": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "UUID of the remote L3Out.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the L3Out.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Description of the L3Out.",
			},
			"vrf_uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the existing VRF associated with the L3Out.",
			},
			"l3_domain": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Name of the L3 domain associated with the L3Out.",
			},
			"target_dscp": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
						"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6", "cs7",
						"expedited_forwarding", "unspecified", "voice_admit",
					),
				},
				MarkdownDescription: "DSCP level of the L3Out.",
			},
			"pim_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
				MarkdownDescription: "Whether protocol-independent multicast is enabled on the L3Out.",
			},
			"bgp": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
					tfplugin.UseConfiguredDisabledObject("enabled"),
				},
				Attributes:          l3OutBGPResourceSchema(),
				MarkdownDescription: "BGP settings of the L3Out.",
			},
			"import_route_control_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
				MarkdownDescription: "Enables import route control on this L3Out.",
			},
			"originate_default_route": l3OutDefaultRouteResourceAttribute(),
			"annotations":             annotationsResourceAttribute(),
			"interface_groups":        interfaceGroupPoliciesResourceAttribute(),
			"ospf": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
					tfplugin.UseConfiguredDisabledObject("enabled"),
				},
				Validators: []validator.Object{
					tfplugin.RequireAttributesWhenEnabled("enabled", "area_id", "area_type"),
				},
				Attributes:          l3OutOSPFResourceSchema(),
				MarkdownDescription: "OSPF settings of the L3Out.",
			},
		},
	}
}

// L3OutDataSourceSchema exposes the same L3Out values as the resource, with
// name and template ID used only to select the object to read.
func L3OutDataSourceSchema() datasourceschema.Schema {
	return datasourceschema.Schema{
		MarkdownDescription: "Reads an L3Out in a Nexus Dashboard Orchestration template.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Stable data source identifier composed of the template UUID and L3Out UUID.",
			},
			"template_id": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the L3Out template containing the L3Out.",
			},
			"uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the L3Out.",
			},
			"name": datasourceschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the L3Out to find in the template.",
			},
			"description": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the L3Out.",
			},
			"vrf_uuid": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the VRF associated with the L3Out.",
			},
			"l3_domain": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the L3 domain associated with the L3Out.",
			},
			"target_dscp": datasourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DSCP level of the L3Out in Terraform schema form.",
			},
			"pim_enabled": datasourceschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether protocol-independent multicast is enabled.",
			},
			"bgp": datasourceschema.SingleNestedAttribute{
				Computed:            true,
				Attributes:          l3OutBGPDataSourceSchema(),
				MarkdownDescription: "BGP configuration of the L3Out, including its enabled state.",
			},
			"import_route_control_enabled": datasourceschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether import route control is enabled.",
			},
			"originate_default_route": l3OutDefaultRouteDataSourceAttribute(),
			"annotations":             annotationsDataSourceAttribute(),
			"interface_groups":        interfaceGroupPoliciesDataSourceAttribute(),
			"ospf": datasourceschema.SingleNestedAttribute{
				Computed:            true,
				Attributes:          l3OutOSPFDataSourceSchema(),
				MarkdownDescription: "OSPF configuration of the L3Out, including its enabled state.",
			},
		},
	}
}

// BuildPatchOperations builds changes for the aggregate L3Out resource. Known
// planned scalar values are compared with the latest NDO object.
func BuildPatchOperations(ctx context.Context, resolved ndoapi.ResolvedObject, plan, configuration, state L3OutModel, diagnostics *diag.Diagnostics) *ndoapi.PatchOperations {
	operations := ndoapi.NewPatchOperations(resolved.Object, resolved.PatchPath())
	operations.SetString("name", tfplugin.KnownStringPointer(plan.Name), nil, false)
	operations.SetString("vrfRef", tfplugin.KnownStringPointer(plan.VRFUUID), nil, false)
	operations.SetString("description", tfplugin.KnownStringPointer(plan.Description), nil, true)
	operations.SetString("l3domain", tfplugin.KnownStringPointer(plan.L3Domain), nil, true)
	operations.SetString("targetDscp", tfplugin.KnownStringPointer(plan.TargetDSCP), targetDSCPTranslation.ToAPI, true)

	if pim := tfplugin.KnownBoolPointer(plan.PIMEnabled); pim != nil {
		operations.Set("pim", *pim)
	}
	if importControl := tfplugin.KnownBoolPointer(plan.ImportRouteControlEnabled); importControl != nil {
		operations.Set("importRouteControl", *importControl)
	}
	routingProtocolChange := l3OutRoutingProtocolChange{}
	ospf, ospfConfigured := l3OutOSPFModelFromTerraform(configuration.OSPF)
	if ospfConfigured {
		if err := setL3OutOSPFOperations(operations, resolved.Object, ospf); err != nil {
			diagnostics.AddError("Failed to Read L3Out OSPF", err.Error())
		}
		ospfEnabled := ospf.Enabled.ValueBool()
		routingProtocolChange.ospf = &ospfEnabled
	}
	defaultRoute := L3OutDefaultRouteModel{Mode: plan.OriginateDefaultRoute}
	if ospfConfigured {
		if !ospf.Enabled.ValueBool() {
			defaultRoute.Always = types.BoolValue(false)
		} else {
			defaultRoute.Always = ospf.OriginateDefaultRouteAlways
		}
	}
	if err := defaultRoute.AddPatchOperations(operations, resolved.Object); err != nil {
		diagnostics.AddError("Failed to Read L3Out Default Route", err.Error())
	}
	if _, configured := l3OutBGPModelFromTerraform(configuration.BGP); configured {
		bgp, _ := l3OutBGPModelFromTerraform(plan.BGP)
		routingProtocolChange.bgp = tfplugin.KnownBoolPointer(bgp.Enabled)
	}
	if routingProtocolChange.bgp != nil || routingProtocolChange.ospf != nil {
		if err := setL3OutRoutingProtocolOperations(operations, resolved.Object, routingProtocolChange); err != nil {
			diagnostics.AddError("Invalid L3Out Routing Protocol", err.Error())
		}
	}
	if !configuration.Annotations.IsNull() {
		annotations := AnnotationsModelFromTerraform(ctx, plan.Annotations, diagnostics)
		if !diagnostics.HasError() {
			if err := annotations.AddPatchOperations(operations, resolved.Object); err != nil {
				diagnostics.AddError("Failed to Read L3Out Annotations", err.Error())
			}
		}
	}
	if !configuration.InterfaceGroups.IsNull() {
		groups := InterfaceGroupPoliciesFromTerraform(ctx, plan.InterfaceGroups, diagnostics)
		configured := InterfaceGroupPoliciesFromTerraform(ctx, configuration.InterfaceGroups, diagnostics)
		previous := InterfaceGroupPoliciesFromTerraform(ctx, state.InterfaceGroups, diagnostics)
		if !diagnostics.HasError() {
			if err := groups.AddPatchOperations(ctx, resolved.Object, operations, resolved.PatchPath(), configured, previous, diagnostics); err != nil {
				diagnostics.AddError("Failed to Read L3Out Interface Groups", err.Error())
			}
		}
	}

	return operations
}

type l3OutRoutingProtocolState struct {
	bgp  bool
	ospf bool
}

type l3OutRoutingProtocolChange struct {
	bgp  *bool
	ospf *bool
}

func newL3OutRoutingProtocolState(value string) (l3OutRoutingProtocolState, error) {
	switch value {
	case "", "none":
		return l3OutRoutingProtocolState{}, nil
	case "bgp":
		return l3OutRoutingProtocolState{bgp: true}, nil
	case "ospf":
		return l3OutRoutingProtocolState{ospf: true}, nil
	case "bgpOspf":
		return l3OutRoutingProtocolState{bgp: true, ospf: true}, nil
	default:
		return l3OutRoutingProtocolState{}, fmt.Errorf("unexpected NDO routingProtocol %q", value)
	}
}

func (state *l3OutRoutingProtocolState) apply(change l3OutRoutingProtocolChange) {
	if change.bgp != nil {
		state.bgp = *change.bgp
	}
	if change.ospf != nil {
		state.ospf = *change.ospf
	}
}

func (state l3OutRoutingProtocolState) value() string {
	switch {
	case state.bgp && state.ospf:
		return "bgpOspf"
	case state.bgp:
		return "bgp"
	case state.ospf:
		return "ospf"
	default:
		return "none"
	}
}

func setL3OutRoutingProtocolOperations(operations *ndoapi.PatchOperations, existing map[string]any, change l3OutRoutingProtocolChange) error {
	value, _, err := ndoapi.StringField(existing, "routingProtocol", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	state, err := newL3OutRoutingProtocolState(value)
	if err != nil {
		return err
	}
	state.apply(change)
	operations.Set("routingProtocol", state.value())
	return nil
}

func l3OutRoutingProtocolValue(change l3OutRoutingProtocolChange) string {
	state := l3OutRoutingProtocolState{}
	state.apply(change)
	return state.value()
}
