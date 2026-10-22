package provider

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &L3OutDataSource{}
var _ datasource.DataSourceWithConfigure = &L3OutDataSource{}

// NewL3OutDataSource returns a data source for reading an existing L3Out from
// an NDO L3Out template.
func NewL3OutDataSource() datasource.DataSource {
	return &L3OutDataSource{}
}

func init() {
	registerDataSource("mso_l3out", NewL3OutDataSource)
}

// L3OutDataSource reads L3Out attributes from an NDO template.
type L3OutDataSource struct {
	client *client.Client
}

func (d *L3OutDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out"
}

func (d *L3OutDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = models.L3OutDataSourceSchema()
}

func (d *L3OutDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "MSO L3Out Data Source: Beginning Configure")
	if request.ProviderData == nil {
		return
	}

	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	d.client = configuredClient
	tflog.Debug(ctx, "MSO L3Out Data Source: Configure Completed")
}

func (d *L3OutDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "MSO L3Out Data Source: Beginning Read")
	var config models.L3OutModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	templateID := config.TemplateID.ValueString()
	l3OutName := config.Name.ValueString()
	l3Out, _, found, err := getL3OutResourceModel(
		ctx,
		d.client,
		templateID,
		config.Path(),
		models.L3OutModel{},
	)
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	if !found {
		response.Diagnostics.AddError(
			"L3Out Not Found",
			fmt.Sprintf("L3Out %q was not found in template %q.", l3OutName, templateID),
		)
		return
	}

	l3Out = l3Out.DataSourceValue(ctx, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &l3Out)...)
	if !response.Diagnostics.HasError() {
		tflog.Debug(ctx, "MSO L3Out Data Source: Read Completed", map[string]interface{}{"name": l3Out.Name.ValueString()})
	}
}
