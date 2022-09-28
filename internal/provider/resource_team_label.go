package provider

import (
	"context"
	"fmt"
	"strings"

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
var _ provider.ResourceType = teamLabelResourceType{}
var _ resource.Resource = teamLabelResource{}
var _ resource.ResourceWithImportState = teamLabelResource{}

type teamLabelResourceType struct{}

func (t teamLabelResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team label.",
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
			"team_id": {
				MarkdownDescription: "Identifier of the team.",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
		},
	}, nil
}

func (t teamLabelResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return teamLabelResource{
		provider: provider,
	}, diags
}

type teamLabelResourceData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	TeamId      types.String `tfsdk:"team_id"`
}

type teamLabelResource struct {
	provider linearProvider
}

func (r teamLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data teamLabelResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name:   data.Name.Value,
		TeamId: &data.TeamId.Value,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	if !data.Color.IsUnknown() {
		input.Color = &data.Color.Value
	}

	response, err := createLabel(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team label")

	issueLabel := response.IssueLabelCreate.IssueLabel

	data.Id = types.String{Value: issueLabel.Id}
	data.Name = types.String{Value: issueLabel.Name}

	if issueLabel.Description != nil {
		data.Description = types.String{Value: *issueLabel.Description}
	}

	if issueLabel.Color != nil {
		data.Color = types.String{Value: *issueLabel.Color}
	}

	if issueLabel.Team != nil {
		data.TeamId = types.String{Value: issueLabel.Team.Id}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data teamLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team label, got error: %s", err))
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

	if issueLabel.Team != nil {
		data.TeamId = types.String{Value: issueLabel.Team.Id}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data teamLabelResourceData

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team label")

	issueLabel := response.IssueLabelUpdate.IssueLabel

	data.Id = types.String{Value: issueLabel.Id}
	data.Name = types.String{Value: issueLabel.Name}

	if issueLabel.Description != nil {
		data.Description = types.String{Value: *issueLabel.Description}
	}

	if issueLabel.Color != nil {
		data.Color = types.String{Value: *issueLabel.Color}
	}

	if issueLabel.Team != nil {
		data.TeamId = types.String{Value: issueLabel.Team.Id}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data teamLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team label")
}

func (r teamLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: label_name:team_key. Got: %q", req.ID),
		)

		return
	}

	response, err := findTeamLabel(context.Background(), r.provider.client, parts[0], parts[1])

	if err != nil || len(response.IssueLabels.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import team label, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.IssueLabels.Nodes[0].Id)...)
}
