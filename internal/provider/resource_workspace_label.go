package provider

import (
	"context"
	"fmt"

	"github.com/frankgreco/terraform-helpers/validators"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = workspaceLabelResourceType{}
var _ tfsdk.Resource = workspaceLabelResource{}
var _ tfsdk.ResourceWithImportState = workspaceLabelResource{}

type workspaceLabelResourceType struct{}

func (t workspaceLabelResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear workspace label.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the label.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the label.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"description": {
				MarkdownDescription: "Description of the label.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"color": {
				MarkdownDescription: "Color of the label.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(colorRegex()),
				},
			},
		},
	}, nil
}

func (t workspaceLabelResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return workspaceLabelResource{
		provider: provider,
	}, diags
}

type workspaceLabelResourceData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
}

type workspaceLabelResource struct {
	provider provider
}

func (r workspaceLabelResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data workspaceLabelResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
	}

	response, err := createWorkspaceLabel(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a workspace label")

	data.Id = types.String{Value: response.IssueLabelCreate.IssueLabel.Id}
	data.Description = types.String{Value: response.IssueLabelCreate.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabelCreate.IssueLabel.Color}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data workspaceLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspaceLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace label, got error: %s", err))
		return
	}

	data.Name = types.String{Value: response.IssueLabel.Name}
	data.Description = types.String{Value: response.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabel.Color}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data workspaceLabelResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelUpdateInput{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
	}

	response, err := updateWorkspaceLabel(context.Background(), r.provider.client, input, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a workspace label")

	data.Name = types.String{Value: response.IssueLabelUpdate.IssueLabel.Name}
	data.Description = types.String{Value: response.IssueLabelUpdate.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabelUpdate.IssueLabel.Color}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data workspaceLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteWorkspaceLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a workspace label")
}

func (r workspaceLabelResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	response, err := findWorkspaceLabel(context.Background(), r.provider.client, req.ID)

	if err != nil || len(response.IssueLabels.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace label, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), response.IssueLabels.Nodes[0].Id)...)
}
