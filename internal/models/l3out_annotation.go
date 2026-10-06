package models

import (
	"fmt"

	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// L3OutAnnotationModel locates one annotation on an L3Out. AnnotationModel
// owns the NDO key/value translation and can be reused under other parents.
type L3OutAnnotationModel struct {
	ID         types.String `tfsdk:"id"`
	TemplateID types.String `tfsdk:"template_id"`
	L3OutUUID  types.String `tfsdk:"l3out_uuid"`
	AnnotationModel
}

// L3OutAnnotationResourceIdentityModel identifies one annotation independently
// of its mutable value.
type L3OutAnnotationResourceIdentityModel struct {
	TemplateID types.String `tfsdk:"template_id"`
	L3OutUUID  types.String `tfsdk:"l3out_uuid"`
	Key        types.String `tfsdk:"key"`
}

// ID returns the string import identifier for this annotation identity.
func (identity L3OutAnnotationResourceIdentityModel) ID() string {
	return fmt.Sprintf("%s/%s/%s", identity.TemplateID.ValueString(), identity.L3OutUUID.ValueString(), identity.Key.ValueString())
}

func (data L3OutAnnotationModel) Path() ndoapi.Path {
	return data.AnnotationModel.Path(NewL3OutPath(data.L3OutUUID.ValueString(), ""))
}

func (data *L3OutAnnotationModel) SetFromNDOObject(templateID, l3outUUID string, object map[string]any) error {
	if err := data.AnnotationModel.SetFromNDOObject(object); err != nil {
		return err
	}
	data.TemplateID = types.StringValue(templateID)
	data.L3OutUUID = types.StringValue(l3outUUID)
	identity := L3OutAnnotationResourceIdentityModel{
		TemplateID: data.TemplateID,
		L3OutUUID:  data.L3OutUUID,
		Key:        data.Key,
	}
	data.ID = types.StringValue(identity.ID())
	return nil
}

func L3OutAnnotationResourceIdentitySchema() identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"template_id": identityschema.StringAttribute{
				Description:       "UUID of the L3Out template.",
				RequiredForImport: true,
			},
			"l3out_uuid": identityschema.StringAttribute{
				Description:       "UUID of the L3Out containing the annotation.",
				RequiredForImport: true,
			},
			"key": identityschema.StringAttribute{
				Description:       "Key of the L3Out annotation.",
				RequiredForImport: true,
			},
		},
	}
}

func L3OutAnnotationResourceSchema() schema.Schema {
	attributes := annotationResourceAttributes()
	attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Stable identity composed of the template UUID, L3Out UUID, and annotation key.",
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
		MarkdownDescription: "UUID of the L3Out containing the annotation.",
	}
	return schema.Schema{
		MarkdownDescription: "Manages one annotation on an L3Out.",
		Attributes:          attributes,
	}
}

func L3OutAnnotationDataSourceSchema() datasourceschema.Schema {
	attributes := annotationDataSourceAttributes()
	attributes["id"] = datasourceschema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Stable identity composed of the template UUID, L3Out UUID, and annotation key.",
	}
	attributes["template_id"] = datasourceschema.StringAttribute{
		Required:            true,
		MarkdownDescription: "UUID of the L3Out template.",
	}
	attributes["l3out_uuid"] = datasourceschema.StringAttribute{
		Required:            true,
		MarkdownDescription: "UUID of the L3Out containing the annotation.",
	}
	return datasourceschema.Schema{
		MarkdownDescription: "Reads one annotation on an L3Out.",
		Attributes:          attributes,
	}
}
