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
var _ tfsdk.ResourceType = workspaceSettingsResourceType{}
var _ tfsdk.Resource = workspaceSettingsResource{}
var _ tfsdk.ResourceWithImportState = workspaceSettingsResource{}

type workspaceSettingsResourceType struct{}

func (t workspaceSettingsResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear workspace settings.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the workspace.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_roadmap": {
				MarkdownDescription: "Enable roadmap for the workspace. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_git_linkback_messages": {
				MarkdownDescription: "Enable git linkbacks for private repositories. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_git_linkback_messages_public": {
				MarkdownDescription: "Enable git linkbacks for public repositories. **Default** `false`.",
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

func (t workspaceSettingsResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return workspaceSettingsResource{
		provider: provider,
	}, diags
}

type workspaceSettingsResourceData struct {
	Id                              types.String `tfsdk:"id"`
	EnableRoadmap                   types.Bool   `tfsdk:"enable_roadmap"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
}

type workspaceSettingsResource struct {
	provider provider
}

func (r workspaceSettingsResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data workspaceSettingsResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{}

	if data.EnableRoadmap.IsNull() || data.EnableRoadmap.IsUnknown() {
		input.RoadmapEnabled = false
	} else {
		input.RoadmapEnabled = data.EnableRoadmap.Value
	}

	if data.EnableGitLinkbackMessages.IsNull() || data.EnableGitLinkbackMessages.IsUnknown() {
		input.GitLinkbackMessagesEnabled = false
	} else {
		input.GitLinkbackMessagesEnabled = data.EnableGitLinkbackMessages.Value
	}

	if data.EnableGitLinkbackMessagesPublic.IsNull() || data.EnableGitLinkbackMessagesPublic.IsUnknown() {
		input.GitPublicLinkbackMessagesEnabled = false
	} else {
		input.GitPublicLinkbackMessagesEnabled = data.EnableGitLinkbackMessagesPublic.Value
	}

	response, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace settings, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.OrganizationUpdate.Organization.Id}
	data.EnableRoadmap = types.Bool{Value: response.OrganizationUpdate.Organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: response.OrganizationUpdate.Organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: response.OrganizationUpdate.Organization.GitPublicLinkbackMessagesEnabled}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data workspaceSettingsResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspaceSettings(context.Background(), r.provider.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace settings, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Organization.Id}
	data.EnableRoadmap = types.Bool{Value: response.Organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: response.Organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: response.Organization.GitPublicLinkbackMessagesEnabled}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data workspaceSettingsResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{}

	if data.EnableRoadmap.IsNull() {
		input.RoadmapEnabled = false
	} else {
		input.RoadmapEnabled = data.EnableRoadmap.Value
	}

	if data.EnableGitLinkbackMessages.IsNull() {
		input.GitLinkbackMessagesEnabled = false
	} else {
		input.GitLinkbackMessagesEnabled = data.EnableGitLinkbackMessages.Value
	}

	if data.EnableGitLinkbackMessagesPublic.IsNull() {
		input.GitPublicLinkbackMessagesEnabled = false
	} else {
		input.GitPublicLinkbackMessagesEnabled = data.EnableGitLinkbackMessagesPublic.Value
	}

	response, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated workspace settings")

	data.Id = types.String{Value: response.OrganizationUpdate.Organization.Id}
	data.EnableRoadmap = types.Bool{Value: response.OrganizationUpdate.Organization.RoadmapEnabled}
	data.EnableGitLinkbackMessages = types.Bool{Value: response.OrganizationUpdate.Organization.GitLinkbackMessagesEnabled}
	data.EnableGitLinkbackMessagesPublic = types.Bool{Value: response.OrganizationUpdate.Organization.GitPublicLinkbackMessagesEnabled}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workspaceSettingsResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data workspaceSettingsResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   false,
		GitLinkbackMessagesEnabled:       false,
		GitPublicLinkbackMessagesEnabled: false,
	}

	_, err := updateWorkspaceSettings(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted workspace settings")
}

func (r workspaceSettingsResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	response, err := getWorkspaceSettings(context.Background(), r.provider.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace settings, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), response.Organization.Id)...)
}
