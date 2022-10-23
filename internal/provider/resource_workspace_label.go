package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/validators"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = workspaceLabelResourceType{}
var _ resource.Resource = workspaceLabelResource{}
var _ resource.ResourceWithImportState = workspaceLabelResource{}

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
					resource.UseStateForUnknown(),
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
					modifiers.NullableString(),
				},
			},
			"color": {
				MarkdownDescription: "Color of the label.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(colorRegex()),
				},
			},
		},
	}, nil
}

func (t workspaceLabelResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
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
	provider linearProvider
}

func (r workspaceLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data workspaceLabelResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name: data.Name.Value,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	if !data.Color.IsUnknown() {
		input.Color = &data.Color.Value
	}

	response, err := createLabel(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a workspace label")

	issueLabel := response.IssueLabelCreate.IssueLabel

	data.Id = types.String{Value: issueLabel.Id}
	data.Name = types.String{Value: issueLabel.Name}

	if issueLabel.Description != nil {
		data.Description = types.String{Value: *issueLabel.Description}
	}

	if issueLabel.Color != nil {
		data.Color = types.String{Value: *issueLabel.Color}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data workspaceLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace label, got error: %s", err))
		return
	}

	issueLabel := response.IssueLabel

	data.Id = types.String{Value: issueLabel.Id}
	data.Name = types.String{Value: issueLabel.Name}

	if issueLabel.Description != nil {
		data.Description = types.String{Value: *issueLabel.Description}
	}

	if issueLabel.Color != nil {
		data.Color = types.String{Value: *issueLabel.Color}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data workspaceLabelResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelUpdateInput{
		Name: data.Name.Value,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	if !data.Color.IsUnknown() {
		input.Color = &data.Color.Value
	}

	response, err := updateLabel(context.Background(), r.provider.client, input, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a workspace label")

	issueLabel := response.IssueLabelUpdate.IssueLabel

	data.Id = types.String{Value: issueLabel.Id}
	data.Name = types.String{Value: issueLabel.Name}

	if issueLabel.Description != nil {
		data.Description = types.String{Value: *issueLabel.Description}
	}

	if issueLabel.Color != nil {
		data.Color = types.String{Value: *issueLabel.Color}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data workspaceLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a workspace label")
}

func (r workspaceLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := findWorkspaceLabel(context.Background(), r.provider.client, req.ID)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace label, got error: %s", err))
		return
	}

	if len(response.IssueLabels.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", "Unable to import team label, got error: label not found")
		return
	}

	if response.IssueLabels.Nodes[0].Team.Id != "" {
		resp.Diagnostics.AddError("Client Error", "Unable to import team label, got error: label is a team label")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.IssueLabels.Nodes[0].Id)...)
}
