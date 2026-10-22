package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &L3OutAnnotationResource{}
var _ resource.ResourceWithIdentity = &L3OutAnnotationResource{}
var _ resource.ResourceWithImportState = &L3OutAnnotationResource{}

// Array indexes change when another annotation is removed. Serialize the
// lookup and patch of standalone annotations within this provider process.
var l3OutAnnotationMutationMu sync.Mutex

func NewL3OutAnnotationResource() resource.Resource {
	return &L3OutAnnotationResource{}
}

func init() {
	registerResource("mso_l3out_annotation", NewL3OutAnnotationResource)
}

// L3OutAnnotationResource manages one tagAnnotations entry by tagKey.
type L3OutAnnotationResource struct {
	client *client.Client
}

func (r *L3OutAnnotationResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out_annotation"
}

func (r *L3OutAnnotationResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = models.L3OutAnnotationResourceIdentitySchema()
}

func (r *L3OutAnnotationResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = models.L3OutAnnotationResourceSchema()
}

func (r *L3OutAnnotationResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData))
		return
	}
	r.client = configuredClient
	tflog.Debug(ctx, "MSO L3Out Annotation Resource: Configure Completed")
}

func (r *L3OutAnnotationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan models.L3OutAnnotationModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutAnnotationMutationMu.Lock()
	defer l3OutAnnotationMutationMu.Unlock()
	templateID := plan.TemplateID.ValueString()
	template, err := ndoapi.GetTemplate(r.client, templateID)
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	appendPath, err := template.AppendPath(plan.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "Failed to Resolve L3Out Annotation Path"), err.Error())
		return
	}
	_, found, err := template.Find(plan.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		response.Diagnostics.AddError("L3Out Annotation Already Exists", fmt.Sprintf("Annotation key %q already exists on L3Out %q. Import it to manage the existing annotation.", plan.Key.ValueString(), plan.L3OutUUID.ValueString()))
		return
	}

	patch := ndoapi.NewPatchPayload("add", appendPath, plan.AnnotationModel.ToPayload())
	if _, err := r.client.PatchbyID(ndoapi.TemplateEndpoint(templateID), patch); err != nil {
		response.Diagnostics.AddError("Failed to Create L3Out Annotation", err.Error())
		return
	}
	r.refreshState(ctx, plan, false, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	r.refreshState(ctx, state, true, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan models.L3OutAnnotationModel
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutAnnotationMutationMu.Lock()
	defer l3OutAnnotationMutationMu.Unlock()
	template, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	resolved, err := template.FindRequired(state.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "L3Out Annotation Not Found"), err.Error())
		return
	}
	operations := ndoapi.NewPatchOperations(resolved.Object, resolved.PatchPath())
	plan.AnnotationModel.AddPatchOperations(operations)
	if err := operations.Apply(r.client, ndoapi.TemplateEndpoint(plan.TemplateID.ValueString())); err != nil {
		response.Diagnostics.AddError("Failed to Update L3Out Annotation", err.Error())
		return
	}
	r.refreshState(ctx, plan, false, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutAnnotationMutationMu.Lock()
	defer l3OutAnnotationMutationMu.Unlock()
	template, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		if errors.Is(err, ndoapi.ErrNotFound) {
			response.State.RemoveResource(ctx)
		} else {
			response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		}
		return
	}
	resolved, found, err := template.Find(state.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		if _, err := r.client.PatchbyID(ndoapi.TemplateEndpoint(state.TemplateID.ValueString()), ndoapi.NewRemovePatchPayload(resolved.PatchPath())); err != nil {
			response.Diagnostics.AddError("Failed to Delete L3Out Annotation", err.Error())
			return
		}
	}
	response.State.RemoveResource(ctx)
}

func (r *L3OutAnnotationResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var identity models.L3OutAnnotationResourceIdentityModel
	if request.ID != "" {
		parts := strings.SplitN(request.ID, "/", 3)
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			response.Diagnostics.AddError("Invalid L3Out Annotation Import ID", "The import ID must be in the form <template_id>/<l3out_uuid>/<key>.")
			return
		}
		identity = models.L3OutAnnotationResourceIdentityModel{
			TemplateID: types.StringValue(parts[0]),
			L3OutUUID:  types.StringValue(parts[1]),
			Key:        types.StringValue(parts[2]),
		}
	} else {
		response.Diagnostics.Append(request.Identity.Get(ctx, &identity)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), identity.TemplateID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("l3out_uuid"), identity.L3OutUUID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("key"), identity.Key)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identity.ID())...)
	response.Diagnostics.Append(response.Identity.Set(ctx, identity)...)
}

func (r *L3OutAnnotationResource) refreshState(ctx context.Context, selector models.L3OutAnnotationModel, removeIfMissing bool, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) {
	template, err := ndoapi.GetTemplate(r.client, selector.TemplateID.ValueString())
	if err != nil {
		if removeIfMissing && errors.Is(err, ndoapi.ErrNotFound) {
			state.RemoveResource(ctx)
			return
		}
		diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	resolved, found, err := template.Find(selector.Path())
	if err != nil {
		diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if !found {
		if removeIfMissing {
			state.RemoveResource(ctx)
		} else {
			diagnostics.AddError("L3Out Annotation Not Found After Apply", fmt.Sprintf("Annotation key %q was not found on L3Out %q after the change was applied.", selector.Key.ValueString(), selector.L3OutUUID.ValueString()))
		}
		return
	}
	var data models.L3OutAnnotationModel
	if err := data.SetFromNDOObject(selector.TemplateID.ValueString(), selector.L3OutUUID.ValueString(), resolved.Object); err != nil {
		diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	diagnostics.Append(state.Set(ctx, &data)...)
	if identity != nil {
		diagnostics.Append(identity.Set(ctx, models.L3OutAnnotationResourceIdentityModel{
			TemplateID: data.TemplateID,
			L3OutUUID:  data.L3OutUUID,
			Key:        data.Key,
		})...)
	}
}
