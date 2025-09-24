package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &fileDataSource{}
)

func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

type fileDataSource struct{}

func (d *fileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *fileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "The name of the file to create.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "The content to write inside the file.",
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "The file format type (e.g., txt, json).",
			},
			"created": schema.BoolAttribute{
				Computed:    true,
				Description: "True if the file was successfully created.",
			},
		},
	}
}

type fileDataSourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Content  types.String `tfsdk:"content"`
	Format   types.String `tfsdk:"format"`
	Created  types.Bool   `tfsdk:"created"`
}

func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state fileDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := os.WriteFile(state.Filename.ValueString(), []byte(state.Content.ValueString()), 0644)
	if err != nil {
		resp.Diagnostics.AddError(
			"File Creation Error",
			fmt.Sprintf("Failed to create file %s: %s", state.Filename.ValueString(), err.Error()),
		)
		return
	}

	state.Created = types.BoolValue(true)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
