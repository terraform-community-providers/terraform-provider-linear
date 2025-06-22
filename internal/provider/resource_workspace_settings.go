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

type WorkspaceSettingsResourceProjectsModel struct {
	UpdateReminderDay       types.String `tfsdk:"update_reminder_day"`
	UpdateReminderHour      types.Int64  `tfsdk:"update_reminder_hour"`
	UpdateReminderFrequency types.Int64  `tfsdk:"update_reminder_frequency"`
}

var projectsAttrTypes = map[string]attr.Type{
	"update_reminder_day":       types.StringType,
	"update_reminder_hour":      types.Int64Type,
	"update_reminder_frequency": types.Int64Type,
}

type WorkspaceSettingsResourceInitiativesModel struct {
	Enabled                 types.Bool   `tfsdk:"enabled"`
	UpdateReminderDay       types.String `tfsdk:"update_reminder_day"`
	UpdateReminderHour      types.Int64  `tfsdk:"update_reminder_hour"`
	UpdateReminderFrequency types.Int64  `tfsdk:"update_reminder_frequency"`
}

var initiativesAttrTypes = map[string]attr.Type{
	"enabled":                   types.BoolType,
	"update_reminder_day":       types.StringType,
	"update_reminder_hour":      types.Int64Type,
	"update_reminder_frequency": types.Int64Type,
}

type WorkspaceSettingsResourceFeedModel struct {
	Enabled  types.Bool   `tfsdk:"enabled"`
	Schedule types.String `tfsdk:"schedule"`
}

var feedAttrTypes = map[string]attr.Type{
	"enabled":  types.BoolType,
	"schedule": types.StringType,
}

type WorkspaceSettingsResourceCustomersModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

var customersAttrTypes = map[string]attr.Type{
	"enabled": types.BoolType,
}

type WorkspaceSettingsResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	AllowMembersToInvite            types.Bool   `tfsdk:"allow_members_to_invite"`
	AllowMembersToCreateTeams       types.Bool   `tfsdk:"allow_members_to_create_teams"`
	AllowMembersToManageLabels      types.Bool   `tfsdk:"allow_members_to_manage_labels"`
	EnableGitLinkbackMessages       types.Bool   `tfsdk:"enable_git_linkback_messages"`
	EnableGitLinkbackMessagesPublic types.Bool   `tfsdk:"enable_git_linkback_messages_public"`
	FiscalYearStartMonth            types.Int64  `tfsdk:"fiscal_year_start_month"`
	Projects                        types.Object `tfsdk:"projects"`
	Initiatives                     types.Object `tfsdk:"initiatives"`
	Feed                            types.Object `tfsdk:"feed"`
	Customers                       types.Object `tfsdk:"customers"`
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
			"fiscal_year_start_month": schema.Int64Attribute{
				MarkdownDescription: "Month at which the fiscal year starts. **Default** `0` representing January.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 11),
				},
			},
			"projects": schema.SingleNestedAttribute{
				MarkdownDescription: "Project settings for the workspace.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						projectsAttrTypes,
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
			"initiatives": schema.SingleNestedAttribute{
				MarkdownDescription: "Initiative settings for the workspace.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						initiativesAttrTypes,
						map[string]attr.Value{
							"enabled":                   types.BoolValue(false),
							"update_reminder_day":       types.StringValue("Friday"),
							"update_reminder_hour":      types.Int64Value(14),
							"update_reminder_frequency": types.Int64Value(0),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable initiatives. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"update_reminder_day": schema.StringAttribute{
						MarkdownDescription: "Day on which to prompt for initiative updates. **Default** `Friday`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("Friday"),
						Validators: []validator.String{
							stringvalidator.OneOf("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"),
						},
					},
					"update_reminder_hour": schema.Int64Attribute{
						MarkdownDescription: "Hour of day (0-23) at which to prompt for initiative updates. **Default** `14`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(2),
						Validators: []validator.Int64{
							int64validator.Between(0, 23),
						},
					},
					"update_reminder_frequency": schema.Int64Attribute{
						MarkdownDescription: "Frequency in weeks to send initiative update reminders. **Default** `0`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(0),
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
			},
			"feed": schema.SingleNestedAttribute{
				MarkdownDescription: "Feed settings for the workspace.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						feedAttrTypes,
						map[string]attr.Value{
							"enabled":  types.BoolValue(false),
							"schedule": types.StringValue("daily"),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable feed summaries for the workspace. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"schedule": schema.StringAttribute{
						MarkdownDescription: "Schedule for feed summaries (daily, weekly). **Default** `daily`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("daily"),
						Validators: []validator.String{
							stringvalidator.OneOf("daily", "weekly"),
						},
					},
				},
			},
			"customers": schema.SingleNestedAttribute{
				MarkdownDescription: "Customer Requests settings for the workspace.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						customersAttrTypes,
						map[string]attr.Value{
							"enabled": types.BoolValue(false),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable customer requests. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
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
	var projectsData *WorkspaceSettingsResourceProjectsModel
	var initiativesData *WorkspaceSettingsResourceInitiativesModel
	var feedData *WorkspaceSettingsResourceFeedModel
	var customersData *WorkspaceSettingsResourceCustomersModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Projects.As(ctx, &projectsData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Initiatives.As(ctx, &initiativesData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Feed.As(ctx, &feedData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Customers.As(ctx, &customersData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := OrganizationUpdateInput{
		AllowMembersToInvite:                     data.AllowMembersToInvite.ValueBool(),
		RestrictTeamCreationToAdmins:             !data.AllowMembersToCreateTeams.ValueBool(),
		RestrictLabelManagementToAdmins:          !data.AllowMembersToManageLabels.ValueBool(),
		RoadmapEnabled:                           initiativesData.Enabled.ValueBool(),
		GitLinkbackMessagesEnabled:               data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled:         data.EnableGitLinkbackMessagesPublic.ValueBool(),
		FiscalYearStartMonth:                     float64(data.FiscalYearStartMonth.ValueInt64()),
		ProjectUpdateReminderFrequencyInWeeks:    float64(projectsData.UpdateReminderFrequency.ValueInt64()),
		ProjectUpdateRemindersDay:                Day(projectsData.UpdateReminderDay.ValueString()),
		ProjectUpdateRemindersHour:               float64(projectsData.UpdateReminderHour.ValueInt64()),
		InitiativeUpdateReminderFrequencyInWeeks: float64(initiativesData.UpdateReminderFrequency.ValueInt64()),
		InitiativeUpdateRemindersDay:             Day(initiativesData.UpdateReminderDay.ValueString()),
		InitiativeUpdateRemindersHour:            float64(initiativesData.UpdateReminderHour.ValueInt64()),
		FeedEnabled:                              feedData.Enabled.ValueBool(),
		DefaultFeedSummarySchedule:               FeedSummarySchedule(feedData.Schedule.ValueString()),
		CustomersEnabled:                         customersData.Enabled.ValueBool(),
	}

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
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)
	data.FiscalYearStartMonth = types.Int64Value(int64(organization.FiscalYearStartMonth))

	data.Projects = types.ObjectValueMust(
		projectsAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

	data.Initiatives = types.ObjectValueMust(
		initiativesAttrTypes,
		map[string]attr.Value{
			"enabled":                   types.BoolValue(organization.RoadmapEnabled),
			"update_reminder_frequency": types.Int64Value(int64(organization.InitiativeUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.InitiativeUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.InitiativeUpdateRemindersHour)),
		},
	)

	data.Feed = types.ObjectValueMust(
		feedAttrTypes,
		map[string]attr.Value{
			"enabled":  types.BoolValue(organization.FeedEnabled),
			"schedule": types.StringValue(string(organization.DefaultFeedSummarySchedule)),
		},
	)

	data.Customers = types.ObjectValueMust(
		customersAttrTypes,
		map[string]attr.Value{
			"enabled": types.BoolValue(customersData.Enabled.ValueBool()),
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
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)
	data.FiscalYearStartMonth = types.Int64Value(int64(organization.FiscalYearStartMonth))

	data.Projects = types.ObjectValueMust(
		projectsAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

	data.Initiatives = types.ObjectValueMust(
		initiativesAttrTypes,
		map[string]attr.Value{
			"enabled":                   types.BoolValue(organization.RoadmapEnabled),
			"update_reminder_frequency": types.Int64Value(int64(organization.InitiativeUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.InitiativeUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.InitiativeUpdateRemindersHour)),
		},
	)

	data.Feed = types.ObjectValueMust(
		feedAttrTypes,
		map[string]attr.Value{
			"enabled":  types.BoolValue(organization.FeedEnabled),
			"schedule": types.StringValue(string(organization.DefaultFeedSummarySchedule)),
		},
	)

	data.Customers = types.ObjectValueMust(
		customersAttrTypes,
		map[string]attr.Value{
			"enabled": types.BoolValue(organization.CustomersEnabled),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkspaceSettingsResourceModel
	var projectsData *WorkspaceSettingsResourceProjectsModel
	var initiativesData *WorkspaceSettingsResourceInitiativesModel
	var feedData *WorkspaceSettingsResourceFeedModel
	var customersData *WorkspaceSettingsResourceCustomersModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Projects.As(ctx, &projectsData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Initiatives.As(ctx, &initiativesData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Feed.As(ctx, &feedData, basetypes.ObjectAsOptions{})...)
	resp.Diagnostics.Append(data.Customers.As(ctx, &customersData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := OrganizationUpdateInput{
		AllowMembersToInvite:                     data.AllowMembersToInvite.ValueBool(),
		RestrictTeamCreationToAdmins:             !data.AllowMembersToCreateTeams.ValueBool(),
		RestrictLabelManagementToAdmins:          !data.AllowMembersToManageLabels.ValueBool(),
		RoadmapEnabled:                           initiativesData.Enabled.ValueBool(),
		GitLinkbackMessagesEnabled:               data.EnableGitLinkbackMessages.ValueBool(),
		GitPublicLinkbackMessagesEnabled:         data.EnableGitLinkbackMessagesPublic.ValueBool(),
		FiscalYearStartMonth:                     float64(data.FiscalYearStartMonth.ValueInt64()),
		ProjectUpdateReminderFrequencyInWeeks:    float64(projectsData.UpdateReminderFrequency.ValueInt64()),
		ProjectUpdateRemindersDay:                Day(projectsData.UpdateReminderDay.ValueString()),
		ProjectUpdateRemindersHour:               float64(projectsData.UpdateReminderHour.ValueInt64()),
		InitiativeUpdateReminderFrequencyInWeeks: float64(initiativesData.UpdateReminderFrequency.ValueInt64()),
		InitiativeUpdateRemindersDay:             Day(initiativesData.UpdateReminderDay.ValueString()),
		InitiativeUpdateRemindersHour:            float64(initiativesData.UpdateReminderHour.ValueInt64()),
		FeedEnabled:                              feedData.Enabled.ValueBool(),
		DefaultFeedSummarySchedule:               FeedSummarySchedule(feedData.Schedule.ValueString()),
		CustomersEnabled:                         customersData.Enabled.ValueBool(),
	}

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
	data.EnableGitLinkbackMessages = types.BoolValue(organization.GitLinkbackMessagesEnabled)
	data.EnableGitLinkbackMessagesPublic = types.BoolValue(organization.GitPublicLinkbackMessagesEnabled)
	data.FiscalYearStartMonth = types.Int64Value(int64(organization.FiscalYearStartMonth))

	data.Projects = types.ObjectValueMust(
		projectsAttrTypes,
		map[string]attr.Value{
			"update_reminder_frequency": types.Int64Value(int64(organization.ProjectUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.ProjectUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.ProjectUpdateRemindersHour)),
		},
	)

	data.Initiatives = types.ObjectValueMust(
		initiativesAttrTypes,
		map[string]attr.Value{
			"enabled":                   types.BoolValue(organization.RoadmapEnabled),
			"update_reminder_frequency": types.Int64Value(int64(organization.InitiativeUpdateReminderFrequencyInWeeks)),
			"update_reminder_day":       types.StringValue(string(organization.InitiativeUpdateRemindersDay)),
			"update_reminder_hour":      types.Int64Value(int64(organization.InitiativeUpdateRemindersHour)),
		},
	)

	data.Feed = types.ObjectValueMust(
		feedAttrTypes,
		map[string]attr.Value{
			"enabled":  types.BoolValue(organization.FeedEnabled),
			"schedule": types.StringValue(string(organization.DefaultFeedSummarySchedule)),
		},
	)

	data.Customers = types.ObjectValueMust(
		customersAttrTypes,
		map[string]attr.Value{
			"enabled": types.BoolValue(customersData.Enabled.ValueBool()),
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
		AllowMembersToInvite:                     true,
		RestrictTeamCreationToAdmins:             false,
		RestrictLabelManagementToAdmins:          false,
		GitLinkbackMessagesEnabled:               true,
		GitPublicLinkbackMessagesEnabled:         false,
		FiscalYearStartMonth:                     0,
		ProjectUpdateReminderFrequencyInWeeks:    0,
		ProjectUpdateRemindersDay:                Day("Friday"),
		ProjectUpdateRemindersHour:               14,
		RoadmapEnabled:                           false,
		InitiativeUpdateReminderFrequencyInWeeks: 0,
		InitiativeUpdateRemindersDay:             Day("Friday"),
		InitiativeUpdateRemindersHour:            14,
		FeedEnabled:                              false,
		DefaultFeedSummarySchedule:               FeedSummarySchedule("daily"),
		CustomersEnabled:                         false,
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
