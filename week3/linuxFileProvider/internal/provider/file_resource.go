package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &fileResource{}

func NewFileResource() resource.Resource {
	return &fileResource{}
}

type fileResource struct{}

type fileResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Content  types.String `tfsdk:"content"`
	Format   types.String `tfsdk:"format"`
	Created  types.Bool   `tfsdk:"created"`
}

func (r *fileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *fileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "The name of the file to manage.",
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

func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := os.WriteFile(plan.Filename.ValueString(), []byte(plan.Content.ValueString()), 0644)
	if err != nil {
		resp.Diagnostics.AddError(
			"File Creation Error",
			fmt.Sprintf("Failed to create file %s: %s", plan.Filename.ValueString(), err.Error()),
		)
		return
	}

	_ = logAction(plan.Filename.ValueString(), "Create \n")

	plan.Created = types.BoolValue(true)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 파일 존재 여부 확인 (없는 경우 Created = false)
	if _, err := os.Stat(state.Filename.ValueString()); os.IsNotExist(err) {
		state.Created = types.BoolValue(false)
	} else {
		state.Created = types.BoolValue(true)
	}

	_ = logAction(state.Filename.ValueString(), "Read \n")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update (content나 format 변경시 파일 다시 씀)
func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := os.WriteFile(plan.Filename.ValueString(), []byte(plan.Content.ValueString()), 0644)
	if err != nil {
		resp.Diagnostics.AddError(
			"File Update Error",
			fmt.Sprintf("Failed to update file %s: %s", plan.Filename.ValueString(), err.Error()),
		)
		return
	}

	_ = logAction(plan.Filename.ValueString(), "Update \n")

	plan.Created = types.BoolValue(true)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete (실제로 파일 삭제)
func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_ = logAction(state.Filename.ValueString(), "Delete \n")

	err := os.Remove(state.Filename.ValueString())
	if err != nil && !os.IsNotExist(err) {
		resp.Diagnostics.AddError(
			"File Deletion Error",
			fmt.Sprintf("Failed to delete file %s: %s", state.Filename.ValueString(), err.Error()),
		)
	}
}

func logAction(filename, action string) error {
	logname := filename + ".log"
	f, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	logLine := fmt.Sprintf("%s called at %s\n", action, time.Now().Format(time.RFC3339))
	_, err = f.WriteString(logLine)
	return err
}
