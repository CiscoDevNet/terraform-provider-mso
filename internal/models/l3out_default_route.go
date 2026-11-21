package models

import (
	"fmt"
	"maps"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var l3OutDefaultRouteTranslation = ndoapi.NewTranslationMap(map[string]string{
	"in_addition": "inAddition",
})

// L3OutDefaultRouteModel maps the public route mode and OSPF-specific always
// setting to the single defaultRouteLeak object used by NDO.
type L3OutDefaultRouteModel struct {
	Mode   types.String
	Always types.Bool
}

func (model *L3OutDefaultRouteModel) SetFromNDOObject(object map[string]any) error {
	defaultRouteLeak, _, err := ndoapi.MapField(object, "defaultRouteLeak", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	mode, exists, err := ndoapi.StringField(defaultRouteLeak, "originateDefaultRoute", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO defaultRouteLeak: %w", err)
	}
	model.Mode = tfplugin.StringOrEmpty(mode, exists, l3OutDefaultRouteTranslation.ToSchema)
	always, exists, err := ndoapi.BoolField(defaultRouteLeak, "always", ndoapi.OptionalField)
	if err != nil {
		return fmt.Errorf("NDO defaultRouteLeak: %w", err)
	}
	model.Always = tfplugin.BoolOrNull(always, exists)
	return nil
}

func (model L3OutDefaultRouteModel) ToPayload() map[string]any {
	mode := tfplugin.KnownStringPointer(model.Mode)
	if mode == nil || *mode == "" {
		return nil
	}
	always := false
	if configured := tfplugin.KnownBoolPointer(model.Always); configured != nil {
		always = *configured
	}
	return map[string]any{
		"originateDefaultRoute": l3OutDefaultRouteTranslation.ToAPI(*mode),
		"always":                always,
	}
}

func (model L3OutDefaultRouteModel) AddPatchOperations(operations *ndoapi.PatchOperations, existing map[string]any) error {
	mode := tfplugin.KnownStringPointer(model.Mode)
	always := tfplugin.KnownBoolPointer(model.Always)
	if mode == nil && always == nil {
		return nil
	}
	if mode != nil && *mode == "" {
		operations.Remove("defaultRouteLeak")
		return nil
	}
	current, _, err := ndoapi.MapField(existing, "defaultRouteLeak", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	defaultRouteLeak := maps.Clone(current)
	if defaultRouteLeak == nil {
		defaultRouteLeak = make(map[string]any)
	}
	if mode != nil {
		defaultRouteLeak["originateDefaultRoute"] = l3OutDefaultRouteTranslation.ToAPI(*mode)
	}
	if _, exists := defaultRouteLeak["originateDefaultRoute"]; !exists {
		return nil
	}
	if always != nil {
		defaultRouteLeak["always"] = *always
	} else if _, exists := defaultRouteLeak["always"]; !exists {
		defaultRouteLeak["always"] = false
	}
	operations.Set("defaultRouteLeak", defaultRouteLeak)
	return nil
}

func l3OutDefaultRouteResourceAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("", "only", "in_addition"),
		},
		MarkdownDescription: "Default-route origination mode for BGP or OSPF.",
	}
}

func l3OutDefaultRouteDataSourceAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Default-route origination mode for BGP or OSPF.",
	}
}
