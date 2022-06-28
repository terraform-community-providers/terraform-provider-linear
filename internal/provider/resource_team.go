package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = teamResourceType{}
var _ tfsdk.Resource = teamResource{}
var _ tfsdk.ResourceWithImportState = teamResource{}

type teamResourceType struct{}

func (t teamResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the team.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"key": {
				MarkdownDescription: "Key of the team.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "Name of the team.",
				Type:                types.StringType,
				Required:            true,
			},
			"private": {
				MarkdownDescription: "Privacy of the team. Default `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"description": {
				MarkdownDescription: "Description of the team.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"icon": {
				MarkdownDescription: "Icon of the team.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"color": {
				MarkdownDescription: "Color of the team.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"timezone": {
				MarkdownDescription: "Timezone of the team. Default `Etc/GMT`.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_issue_history_grouping": {
				MarkdownDescription: "Enable issue history grouping for the team. Default `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"no_priority_issues_first": {
				MarkdownDescription: "Prefer issues without priority during issue prioritization order. Default `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t teamResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return teamResource{
		provider: provider,
	}, diags
}

type teamResourceData struct {
	Id                         types.String `tfsdk:"id"`
	Key                        types.String `tfsdk:"key"`
	Name                       types.String `tfsdk:"name"`
	Private                    types.Bool   `tfsdk:"private"`
	Description                types.String `tfsdk:"description"`
	Icon                       types.String `tfsdk:"icon"`
	Color                      types.String `tfsdk:"color"`
	Timezone                   types.String `tfsdk:"timezone"`
	EnableIssueHistoryGrouping types.Bool   `tfsdk:"enable_issue_history_grouping"`
	NoPriorityIssuesFirst      types.Bool   `tfsdk:"no_priority_issues_first"`
}

type teamResource struct {
	provider provider
}

func (r teamResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data teamResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TeamCreateInput{
		Key:         data.Key.Value,
		Name:        data.Name.Value,
		Private:     data.Private.Value,
		Description: data.Description.Value,
		Icon:        data.Icon.Value,
		Color:       data.Color.Value,
	}

	if data.Timezone.IsNull() {
		input.Timezone = "Etc/GMT"
	} else {
		input.Timezone = data.Timezone.Value
	}

	if data.EnableIssueHistoryGrouping.IsNull() {
		input.GroupIssueHistory = true
	} else {
		input.GroupIssueHistory = data.EnableIssueHistoryGrouping.Value
	}

	if data.NoPriorityIssuesFirst.IsNull() {
		input.IssueOrderingNoPriorityFirst = true
	} else {
		input.IssueOrderingNoPriorityFirst = data.NoPriorityIssuesFirst.Value
	}

	response, err := createTeam(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team")

	data.Id = types.String{Value: response.TeamCreate.Team.Id}
	data.Private = types.Bool{Value: response.TeamCreate.Team.Private}
	data.Description = types.String{Value: response.TeamCreate.Team.Description}
	data.Icon = types.String{Value: response.TeamCreate.Team.Icon}
	data.Color = types.String{Value: response.TeamCreate.Team.Color}
	data.Timezone = types.String{Value: response.TeamCreate.Team.Timezone}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.TeamCreate.Team.GroupIssueHistory}
	data.NoPriorityIssuesFirst = types.Bool{Value: response.TeamCreate.Team.IssueOrderingNoPriorityFirst}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data teamResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeam(context.Background(), r.provider.client, data.Key.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Team.Id}
	data.Private = types.Bool{Value: response.Team.Private}
	data.Description = types.String{Value: response.Team.Description}
	data.Icon = types.String{Value: response.Team.Icon}
	data.Color = types.String{Value: response.Team.Color}
	data.Timezone = types.String{Value: response.Team.Timezone}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.Team.GroupIssueHistory}
	data.NoPriorityIssuesFirst = types.Bool{Value: response.Team.IssueOrderingNoPriorityFirst}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data teamResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TeamUpdateInput{
		Key:                          data.Key.Value,
		Name:                         data.Name.Value,
		Private:                      data.Private.Value,
		Description:                  data.Description.Value,
		Icon:                         data.Icon.Value,
		Color:                        data.Color.Value,
		Timezone:                     data.Timezone.Value,
		GroupIssueHistory:            data.EnableIssueHistoryGrouping.Value,
		IssueOrderingNoPriorityFirst: data.NoPriorityIssuesFirst.Value,
	}

	var key string

	diags = req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("key"), &key)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := updateTeam(context.Background(), r.provider.client, input, key)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team")

	data.Id = types.String{Value: response.TeamUpdate.Team.Id}
	data.Private = types.Bool{Value: response.TeamUpdate.Team.Private}
	data.Description = types.String{Value: response.TeamUpdate.Team.Description}
	data.Icon = types.String{Value: response.TeamUpdate.Team.Icon}
	data.Color = types.String{Value: response.TeamUpdate.Team.Color}
	data.Timezone = types.String{Value: response.TeamUpdate.Team.Timezone}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.TeamUpdate.Team.GroupIssueHistory}
	data.NoPriorityIssuesFirst = types.Bool{Value: response.TeamUpdate.Team.IssueOrderingNoPriorityFirst}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data teamResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteTeam(context.Background(), r.provider.client, data.Key.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team")
}

func (r teamResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("key"), req, resp)
}
