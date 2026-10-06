package provider

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &L3OutAnnotationDataSource{}
var _ datasource.DataSourceWithConfigure = &L3OutAnnotationDataSource{}

func NewL3OutAnnotationDataSource() datasource.DataSource {
	return &L3OutAnnotationDataSource{}
}

func init() {
	registerDataSource("mso_l3out_annotation", NewL3OutAnnotationDataSource)
}

// L3OutAnnotationDataSource reads one annotation by its key on an L3Out.
type L3OutAnnotationDataSource struct {
	client *client.Client
}

func (d *L3OutAnnotationDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out_annotation"
}

func (d *L3OutAnnotationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = models.L3OutAnnotationDataSourceSchema()
}

func (d *L3OutAnnotationDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData))
		return
	}
	d.client = configuredClient
	tflog.Debug(ctx, "MSO L3Out Annotation Data Source: Configure Completed")
}

func (d *L3OutAnnotationDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var config models.L3OutAnnotationModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	template, err := ndoapi.GetTemplate(d.client, config.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	resolved, err := template.FindRequired(config.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "L3Out Annotation Not Found"), err.Error())
		return
	}
	var annotation models.L3OutAnnotationModel
	if err := annotation.SetFromNDOObject(config.TemplateID.ValueString(), config.L3OutUUID.ValueString(), resolved.Object); err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &annotation)...)
}
