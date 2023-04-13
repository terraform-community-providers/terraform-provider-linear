package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &WorkspaceSettingsResource{}
var _ resource.ResourceWithImportState = &WorkspaceSettingsResource{}

func NewWorkspaceSettingsResource() resource.Resource {
	return &WorkspaceSettingsResource{}
}

type WorkspaceSettingsResource struct {
	client *graphql.Client
}

type WorkspaceSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	EnableRoadmap                   types.Bool   `tfsdk:"enable_roadmap"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
}

func (r *WorkspaceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_settings"
}

func (r *WorkspaceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear workspace settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the workspace.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_roadmap": schema.BoolAttribute{
				MarkdownDescription: "Enable roadmap for the workspace. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_git_linkback_messages": schema.BoolAttribute{
				MarkdownDescription: "Enable git linkbacks for private repositories. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"enable_git_linkback_messages_public": schema.BoolAttribute{
				MarkdownDescription: "Enable git linkbacks for public repositories. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *WorkspaceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *WorkspaceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.ValueBool(),
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.ValueBool(),
	}

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace settings, got error: %s", err))
		return
	}

	organization := response.OrganizationUpdate.Organization

	data.Id = types.StringValue(organization.Id)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkspaceSettings(ctx, *r.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace settings, got error: %s", err))
		return
	}

	organization := response.Organization

	data.Id = types.StringValue(organization.Id)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   data.EnableRoadmap.ValueBool(),
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.ValueBool(),
	}

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated workspace settings")

	organization := response.OrganizationUpdate.Organization

	data.Id = types.StringValue(organization.Id)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := UpdateOrganizationInput{
		RoadmapEnabled:                   false,
		GitLinkbackMessagesEnabled:       false,
		GitPublicLinkbackMessagesEnabled: false,
	}

	_, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted workspace settings")
}

func (r *WorkspaceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
