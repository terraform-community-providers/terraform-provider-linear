package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

type WorkspaceSettingsResourceProjectModel struct {
	UpdateReminderDay       types.String `tfsdk:"update_reminder_day"`
	UpdateReminderHour      types.Int64  `tfsdk:"update_reminder_hour"`
	UpdateReminderFrequency types.Int64  `tfsdk:"update_reminder_frequency"`
}

var projectAttrTypes = map[string]attr.Type{
	"update_reminder_day":       types.StringType,
	"update_reminder_hour":      types.Int64Type,
	"update_reminder_frequency": types.Int64Type,
}

type WorkspaceSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	AllowMembersToInvite            types.Bool   `tfsdk:"allow_members_to_invite"`
	AllowMembersToCreateTeams       types.Bool   `tfsdk:"allow_members_to_create_teams"`
	AllowMembersToManageLabels      types.Bool   `tfsdk:"allow_members_to_manage_labels"`
	EnableRoadmap                   types.Bool   `tfsdk:"enable_roadmap"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
	Projects                        types.Object `tfsdk:"projects"`
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
			"allow_members_to_invite": schema.BoolAttribute{
				MarkdownDescription: "Allow members to invite new members to the workspace. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"allow_members_to_create_teams": schema.BoolAttribute{
				MarkdownDescription: "Allow members to create new teams in the workspace. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"allow_members_to_manage_labels": schema.BoolAttribute{
				MarkdownDescription: "Allow members to manage labels in the workspace. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
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
			"projects": schema.SingleNestedAttribute{
				MarkdownDescription: "Project settings for the workspace.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						projectAttrTypes,
						map[string]attr.Value{
							"update_reminder_day":       types.StringValue("Friday"),
							"update_reminder_hour":      types.Int64Value(14),
							"update_reminder_frequency": types.Int64Value(0),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"update_reminder_day": schema.StringAttribute{
						MarkdownDescription: "Day on which to prompt for project updates. **Default** `Friday`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("Friday"),
						Validators: []validator.String{
							stringvalidator.OneOf("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"),
						},
					},
					"update_reminder_hour": schema.Int64Attribute{
						MarkdownDescription: "Hour of day (0-23) at which to prompt for project updates. **Default** `14`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(2),
						Validators: []validator.Int64{
							int64validator.Between(0, 23),
						},
					},
					"update_reminder_frequency": schema.Int64Attribute{
						MarkdownDescription: "Frequency in weeks to send project update reminders. **Default** `0`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(0),
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
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
	var projectsData *WorkspaceSettingsResourceProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := OrganizationUpdateInput{
		AllowMembersToInvite:             data.AllowMembersToInvite.ValueBool(),
		RestrictTeamCreationToAdmins:     !data.AllowMembersToCreateTeams.ValueBool(),
		RestrictLabelManagementToAdmins:  !data.AllowMembersToManageLabels.ValueBool(),
		RoadmapEnabled:                   data.EnableRoadmap.ValueBool(),
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.ValueBool(),
	}

	resp.Diagnostics.Append(data.Projects.As(ctx, &projectsData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.ProjectUpdateReminderFrequencyInWeeks = float64(projectsData.UpdateReminderFrequency.ValueInt64())
	input.ProjectUpdateRemindersDay = Day(projectsData.UpdateReminderDay.ValueString())
	input.ProjectUpdateRemindersHour = float64(projectsData.UpdateReminderHour.ValueInt64())

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace settings, got error: %s", err))
		return
	}

	organization := response.OrganizationUpdate.Organization

	data.Id = types.StringValue(organization.Id)
	data.AllowMembersToInvite = types.BoolValue(organization.AllowMembersToInvite)
	data.AllowMembersToCreateTeams = types.BoolValue(!organization.RestrictTeamCreationToAdmins)
	data.AllowMembersToManageLabels = types.BoolValue(!organization.RestrictLabelManagementToAdmins)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	data.Projects = types.ObjectValueMust(
		projectAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

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
	data.AllowMembersToInvite = types.BoolValue(organization.AllowMembersToInvite)
	data.AllowMembersToCreateTeams = types.BoolValue(!organization.RestrictTeamCreationToAdmins)
	data.AllowMembersToManageLabels = types.BoolValue(!organization.RestrictLabelManagementToAdmins)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	data.Projects = types.ObjectValueMust(
		projectAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkspaceSettingsResourceModel
	var projectsData *WorkspaceSettingsResourceProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := OrganizationUpdateInput{
		AllowMembersToInvite:             data.AllowMembersToInvite.ValueBool(),
		RestrictTeamCreationToAdmins:     !data.AllowMembersToCreateTeams.ValueBool(),
		RestrictLabelManagementToAdmins:  !data.AllowMembersToManageLabels.ValueBool(),
		RoadmapEnabled:                   data.EnableRoadmap.ValueBool(),
		GitLinkbackMessagesEnabled:       data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled: data.EnableGitLinkbackMessagesPublic.ValueBool(),
	}

	resp.Diagnostics.Append(data.Projects.As(ctx, &projectsData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.ProjectUpdateReminderFrequencyInWeeks = float64(projectsData.UpdateReminderFrequency.ValueInt64())
	input.ProjectUpdateRemindersDay = Day(projectsData.UpdateReminderDay.ValueString())
	input.ProjectUpdateRemindersHour = float64(projectsData.UpdateReminderHour.ValueInt64())

	response, err := updateWorkspaceSettings(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace settings, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated workspace settings")

	organization := response.OrganizationUpdate.Organization

	data.Id = types.StringValue(organization.Id)
	data.AllowMembersToInvite = types.BoolValue(organization.AllowMembersToInvite)
	data.AllowMembersToCreateTeams = types.BoolValue(!organization.RestrictTeamCreationToAdmins)
	data.AllowMembersToManageLabels = types.BoolValue(!organization.RestrictLabelManagementToAdmins)
	data.EnableRoadmap = types.BoolValue(organization.RoadmapEnabled)
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)

	data.Projects = types.ObjectValueMust(
		projectAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WorkspaceSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := OrganizationUpdateInput{
		AllowMembersToInvite:                  true,
		RestrictTeamCreationToAdmins:          false,
		RestrictLabelManagementToAdmins:       false,
		RoadmapEnabled:                        false,
		GitLinkbackMessagesEnabled:            true,
		GitPublicLinkbackMessagesEnabled:      false,
		ProjectUpdateReminderFrequencyInWeeks: 0,
		ProjectUpdateRemindersDay:             Day("Friday"),
		ProjectUpdateRemindersHour:            14,
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
