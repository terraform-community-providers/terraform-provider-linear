package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/frankgreco/terraform-helpers/validators"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = teamLabelResourceType{}
var _ tfsdk.Resource = teamLabelResource{}
var _ tfsdk.ResourceWithImportState = teamLabelResource{}

type teamLabelResourceType struct{}

func (t teamLabelResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team label.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the team.",
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
					// TODO: Color value validation
				},
			},
			"team_id": {
				MarkdownDescription: "Identifier of the team.",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					// TODO: UUID validation
					validators.MinLength(1),
				},
			},
		},
	}, nil
}

func (t teamLabelResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
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
	provider provider
}

func (r teamLabelResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data teamLabelResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
		TeamId:      data.TeamId.Value,
	}

	response, err := createTeamLabel(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team label")

	data.Id = types.String{Value: response.IssueLabelCreate.IssueLabel.Id}
	data.Description = types.String{Value: response.IssueLabelCreate.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabelCreate.IssueLabel.Color}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data teamLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeamLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team label, got error: %s", err))
		return
	}

	data.Name = types.String{Value: response.IssueLabel.Name}
	data.Description = types.String{Value: response.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabel.Color}
	data.TeamId = types.String{Value: response.IssueLabel.Team.Id}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data teamLabelResourceData

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

	response, err := updateTeamLabel(context.Background(), r.provider.client, input, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team label")

	data.Name = types.String{Value: response.IssueLabelUpdate.IssueLabel.Name}
	data.Description = types.String{Value: response.IssueLabelUpdate.IssueLabel.Description}
	data.Color = types.String{Value: response.IssueLabelUpdate.IssueLabel.Color}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamLabelResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data teamLabelResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteTeamLabel(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team label")
}

func (r teamLabelResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)

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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), response.IssueLabels.Nodes[0].Id)...)
}
