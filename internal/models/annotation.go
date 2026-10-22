package models

import (
	"context"
	"fmt"
	"maps"
	"slices"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AnnotationModel represents one NDO tag annotation. Its parent can be any
// model whose NDO object contains a tagAnnotations collection.
type AnnotationModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// Path selects this annotation by tagKey beneath an arbitrary parent path.
func (model AnnotationModel) Path(parent ndoapi.Path) ndoapi.Path {
	return parent.WithSteps(ndoapi.PathStep{
		Field: "tagAnnotations",
		Selector: &ndoapi.ObjectSelector{
			Keys: []ndoapi.ObjectIdentifier{{Field: "tagKey", Value: model.Key.ValueString()}},
		},
	})
}

// SetFromNDOObject translates one decoded NDO annotation into the model.
func (model *AnnotationModel) SetFromNDOObject(object map[string]any) error {
	key, _, err := ndoapi.StringField(object, "tagKey", ndoapi.RequiredNonEmptyField)
	if err != nil {
		return fmt.Errorf("NDO annotation: %w", err)
	}
	value, _, err := ndoapi.StringField(object, "tagValue", ndoapi.RequiredField)
	if err != nil {
		return fmt.Errorf("NDO annotation %q: %w", key, err)
	}
	model.Key = types.StringValue(key)
	model.Value = types.StringValue(value)
	return nil
}

// ToPayload translates one annotation to the NDO representation.
func (model AnnotationModel) ToPayload() map[string]any {
	return map[string]any{
		"tagKey":   model.Key.ValueString(),
		"tagValue": model.Value.ValueString(),
	}
}

// AddPatchOperations updates only this annotation's value.
func (model AnnotationModel) AddPatchOperations(operations *ndoapi.PatchOperations) {
	operations.Set("tagValue", model.Value.ValueString())
}

// AnnotationsModel represents an entire tagAnnotations collection. A nil map
// is unmanaged; a non-nil empty map represents an explicitly empty collection.
type AnnotationsModel map[string]string

// AnnotationsModelFromTerraform reads a known map using the shared tfplugin
// conversion. Null and unknown Terraform maps remain unmanaged.
func AnnotationsModelFromTerraform(ctx context.Context, value types.Map, diagnostics *diag.Diagnostics) AnnotationsModel {
	return AnnotationsModel(tfplugin.StringMapFromTerraform(ctx, value, diagnostics))
}

// TerraformValue returns the annotations as a Terraform string map.
func (model AnnotationsModel) TerraformValue(ctx context.Context, diagnostics *diag.Diagnostics) types.Map {
	return tfplugin.StringMap(ctx, map[string]string(model), diagnostics)
}

// SetFromNDOObject reads the annotations from a parent NDO object. An absent
// collection becomes an empty map; duplicate keys and malformed entries fail.
func (model *AnnotationsModel) SetFromNDOObject(object map[string]any) error {
	annotations := make(AnnotationsModel)
	items, exists, err := ndoapi.ListField(object, "tagAnnotations", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	if !exists {
		*model = annotations
		return nil
	}
	for index, item := range items {
		annotationObject, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("NDO tagAnnotations[%d] has unexpected type %T", index, item)
		}
		var annotation AnnotationModel
		if err := annotation.SetFromNDOObject(annotationObject); err != nil {
			return fmt.Errorf("NDO tagAnnotations[%d]: %w", index, err)
		}
		key := annotation.Key.ValueString()
		if _, exists := annotations[key]; exists {
			return fmt.Errorf("NDO tagAnnotations contains duplicate tagKey %q", key)
		}
		annotations[key] = annotation.Value.ValueString()
	}
	*model = annotations
	return nil
}

// ToPayload returns the entire collection in stable key order.
func (model AnnotationsModel) ToPayload() []map[string]any {
	if model == nil {
		return nil
	}
	keys := slices.Sorted(maps.Keys(model))
	payload := make([]map[string]any, 0, len(keys))
	for _, key := range keys {
		payload = append(payload, AnnotationModel{Key: types.StringValue(key), Value: types.StringValue(model[key])}.ToPayload())
	}
	return payload
}

// AddPatchOperations replaces the collection only when its key/value pairs
// differ from the latest NDO object. API array ordering is immaterial.
func (model AnnotationsModel) AddPatchOperations(operations *ndoapi.PatchOperations, existing map[string]any) error {
	if model == nil {
		return nil
	}
	var current AnnotationsModel
	if err := current.SetFromNDOObject(existing); err != nil {
		return err
	}
	if maps.Equal(model, current) {
		return nil
	}
	operations.Set("tagAnnotations", model.ToPayload())
	return nil
}

func annotationsResourceAttribute() schema.MapAttribute {
	return schema.MapAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Map{
			mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
			mapvalidator.NoNullValues(),
		},
		MarkdownDescription: "Annotations on the parent object, keyed by annotation key.",
	}
}

func annotationResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			MarkdownDescription: "Unique annotation key on the parent object.",
		},
		"value": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Value of the annotation.",
		},
	}
}

func annotationsDataSourceAttribute() datasourceschema.MapAttribute {
	return datasourceschema.MapAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		MarkdownDescription: "Annotations on the parent object, keyed by annotation key.",
	}
}

func annotationDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"key": datasourceschema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			MarkdownDescription: "Annotation key to find on the parent object.",
		},
		"value": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Value of the annotation.",
		},
	}
}
