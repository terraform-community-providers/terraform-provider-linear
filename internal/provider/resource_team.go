package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/frankgreco/terraform-helpers/validators"
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
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
					validators.MaxLength(5),
					validators.NoWhitespace(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the team.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(2),
				},
			},
			"private": {
				MarkdownDescription: "Privacy of the team. **Default** `false`.",
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
				Validators: []tfsdk.AttributeValidator{
					validators.Match(regexp.MustCompile("^[a-zA-Z]+$")),
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
				Validators: []tfsdk.AttributeValidator{
					validators.Match(colorRegex()),
				},
			},
			"timezone": {
				MarkdownDescription: "Timezone of the team. **Default** `Etc/GMT`.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"no_priority_issues_first": {
				MarkdownDescription: "Prefer issues without priority at the top during issue prioritization order. **Default** `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_issue_history_grouping": {
				MarkdownDescription: "Enable issue history grouping for the team. **Default** `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_issue_default_to_bottom": {
				MarkdownDescription: "Enable moving issues to bottom of the column when changing state. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"auto_archive_period": {
				MarkdownDescription: "Period after which closed and completed issues are automatically archived, in months. **Default** `3`.",
				Type:                types.Float64Type,
				// #2
				// Optional:            true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.FloatInSlice(1, 3, 6, 9, 12),
				},
			},
			"triage": {
				MarkdownDescription: "Triage settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enabled": {
						MarkdownDescription: "Enable triage mode for the team. **Default** `false`.",
						Type:                types.BoolType,
						Required:            true,
					},
				}),
			},
			"cycles": {
				MarkdownDescription: "Cycle settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enabled": {
						MarkdownDescription: "Enable cycles for the team. **Default** `false`.",
						Type:                types.BoolType,
						Required:            true,
					},
					"start_day": {
						MarkdownDescription: "Start day of the cycle. Sunday is 0, Saturday is 6. **Default** `0`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.FloatInSlice(0, 1, 2, 3, 4, 5, 6),
						},
					},
					"duration": {
						MarkdownDescription: "Duration of the cycle in weeks. **Default** `1`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.FloatInSlice(1, 2, 3, 4, 5, 6, 7, 8),
						},
					},
					"cooldown": {
						MarkdownDescription: "Cooldown time between cycles in weeks. **Default** `0`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.FloatInSlice(0, 1, 2, 3),
						},
					},
					"upcoming": {
						MarkdownDescription: "Number of upcoming cycles to automatically create. **Default** `2`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.FloatInSlice(1, 2, 3, 4, 6, 8, 10),
						},
					},
					"auto_add_started": {
						MarkdownDescription: "Auto add started issues that don't belong to any cycle to the active cycle. **Default** `true`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
					},
					"auto_add_completed": {
						MarkdownDescription: "Auto add completed issues that don't belong to any cycle to the active cycle. **Default** `true`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
					},
					"need_for_active": {
						MarkdownDescription: "Whether all active issues need to have a cycle. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
					},
				}),
			},
			"estimation": {
				MarkdownDescription: "Issue estimation settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"type": {
						MarkdownDescription: "Issue estimation type for the team. **Default** `notUsed`.",
						Type:                types.StringType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, "notUsed", "exponential", "fibonacci", "linear", "tShirt"),
						},
					},
					"extended": {
						MarkdownDescription: "Whether the team uses extended estimation. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
					},
					"allow_zero": {
						MarkdownDescription: "Whether zero is allowed as an estimation. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
					},
					"default": {
						MarkdownDescription: "Default estimation for issues that are unestimated. **Default** `1`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.UseStateForUnknown(),
						},
						Validators: []tfsdk.AttributeValidator{
							validators.FloatInSlice(0, 1),
						},
					},
				}),
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

type triage struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type cycles struct {
	Enabled          types.Bool    `tfsdk:"enabled"`
	StartDay         types.Float64 `tfsdk:"start_day"`
	Duration         types.Float64 `tfsdk:"duration"`
	Cooldown         types.Float64 `tfsdk:"cooldown"`
	Upcoming         types.Float64 `tfsdk:"upcoming"`
	AutoAddStarted   types.Bool    `tfsdk:"auto_add_started"`
	AutoAddCompleted types.Bool    `tfsdk:"auto_add_completed"`
	NeedForActive    types.Bool    `tfsdk:"need_for_active"`
}

type estimation struct {
	Type      types.String  `tfsdk:"type"`
	Extended  types.Bool    `tfsdk:"extended"`
	AllowZero types.Bool    `tfsdk:"allow_zero"`
	Default   types.Float64 `tfsdk:"default"`
}

type teamResourceData struct {
	Id                         types.String  `tfsdk:"id"`
	Key                        types.String  `tfsdk:"key"`
	Name                       types.String  `tfsdk:"name"`
	Private                    types.Bool    `tfsdk:"private"`
	Description                types.String  `tfsdk:"description"`
	Icon                       types.String  `tfsdk:"icon"`
	Color                      types.String  `tfsdk:"color"`
	Timezone                   types.String  `tfsdk:"timezone"`
	NoPriorityIssuesFirst      types.Bool    `tfsdk:"no_priority_issues_first"`
	EnableIssueHistoryGrouping types.Bool    `tfsdk:"enable_issue_history_grouping"`
	EnableIssueDefaultToBottom types.Bool    `tfsdk:"enable_issue_default_to_bottom"`
	AutoArchivePeriod          types.Float64 `tfsdk:"auto_archive_period"`
	Triage                     *triage       `tfsdk:"triage"`
	Cycles                     *cycles       `tfsdk:"cycles"`
	Estimation                 *estimation   `tfsdk:"estimation"`
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
		Key:                           data.Key.Value,
		Name:                          data.Name.Value,
		Private:                       data.Private.Value,
		Description:                   data.Description.Value,
		Icon:                          data.Icon.Value,
		Color:                         data.Color.Value,
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.Value,
	}

	if data.Timezone.IsNull() {
		input.Timezone = "Etc/GMT"
	} else {
		input.Timezone = data.Timezone.Value
	}

	if data.NoPriorityIssuesFirst.IsNull() {
		input.IssueOrderingNoPriorityFirst = true
	} else {
		input.IssueOrderingNoPriorityFirst = data.NoPriorityIssuesFirst.Value
	}

	if data.EnableIssueHistoryGrouping.IsNull() {
		input.GroupIssueHistory = true
	} else {
		input.GroupIssueHistory = data.EnableIssueHistoryGrouping.Value
	}

	// #2
	// if data.AutoArchivePeriod.IsNull() {
	input.AutoArchivePeriod = 3
	// } else {
	// 	input.AutoArchivePeriod = data.AutoArchivePeriod.Value
	// }

	if data.Triage != nil {
		input.TriageEnabled = data.Triage.Enabled.Value
	}

	if data.Cycles != nil {
		input.CyclesEnabled = data.Cycles.Enabled.Value
		input.CycleStartDay = data.Cycles.StartDay.Value
		input.CycleCooldownTime = int(data.Cycles.Cooldown.Value)
		input.CycleLockToActive = data.Cycles.NeedForActive.Value
	}

	if data.Cycles == nil || data.Cycles.Duration.IsNull() {
		input.CycleDuration = 1
	} else {
		input.CycleDuration = int(data.Cycles.Duration.Value)
	}

	if data.Cycles == nil || data.Cycles.Upcoming.IsNull() {
		input.UpcomingCycleCount = 2
	} else {
		input.UpcomingCycleCount = data.Cycles.Upcoming.Value
	}

	if data.Cycles == nil || data.Cycles.AutoAddStarted.IsNull() {
		input.CycleIssueAutoAssignStarted = true
	} else {
		input.CycleIssueAutoAssignStarted = data.Cycles.AutoAddStarted.Value
	}

	if data.Cycles == nil || data.Cycles.AutoAddCompleted.IsNull() {
		input.CycleIssueAutoAssignCompleted = true
	} else {
		input.CycleIssueAutoAssignCompleted = data.Cycles.AutoAddCompleted.Value
	}

	if data.Estimation != nil {
		input.IssueEstimationExtended = data.Estimation.Extended.Value
		input.IssueEstimationAllowZero = data.Estimation.AllowZero.Value
	}

	if data.Estimation == nil || data.Estimation.Type.IsNull() {
		input.IssueEstimationType = "notUsed"
	} else {
		input.IssueEstimationType = data.Estimation.Type.Value
	}

	if data.Estimation == nil || data.Estimation.Default.IsNull() {
		input.DefaultIssueEstimate = 1
	} else {
		input.DefaultIssueEstimate = data.Estimation.Default.Value
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
	data.NoPriorityIssuesFirst = types.Bool{Value: response.TeamCreate.Team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.TeamCreate.Team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: response.TeamCreate.Team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: response.TeamCreate.Team.AutoArchivePeriod}

	data.Triage = &triage{
		Enabled: types.Bool{Value: response.TeamCreate.Team.TriageEnabled},
	}

	data.Cycles = &cycles{
		Enabled:          types.Bool{Value: response.TeamCreate.Team.CyclesEnabled},
		StartDay:         types.Float64{Value: response.TeamCreate.Team.CycleStartDay},
		Duration:         types.Float64{Value: response.TeamCreate.Team.CycleDuration},
		Cooldown:         types.Float64{Value: response.TeamCreate.Team.CycleCooldownTime},
		Upcoming:         types.Float64{Value: response.TeamCreate.Team.UpcomingCycleCount},
		AutoAddStarted:   types.Bool{Value: response.TeamCreate.Team.CycleIssueAutoAssignStarted},
		AutoAddCompleted: types.Bool{Value: response.TeamCreate.Team.CycleIssueAutoAssignCompleted},
		NeedForActive:    types.Bool{Value: response.TeamCreate.Team.CycleLockToActive},
	}

	data.Estimation = &estimation{
		Type:      types.String{Value: response.TeamCreate.Team.IssueEstimationType},
		Extended:  types.Bool{Value: response.TeamCreate.Team.IssueEstimationExtended},
		AllowZero: types.Bool{Value: response.TeamCreate.Team.IssueEstimationAllowZero},
		Default:   types.Float64{Value: response.TeamCreate.Team.DefaultIssueEstimate},
	}

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
	data.Name = types.String{Value: response.Team.Name}
	data.Private = types.Bool{Value: response.Team.Private}
	data.Description = types.String{Value: response.Team.Description}
	data.Icon = types.String{Value: response.Team.Icon}
	data.Color = types.String{Value: response.Team.Color}
	data.Timezone = types.String{Value: response.Team.Timezone}
	data.NoPriorityIssuesFirst = types.Bool{Value: response.Team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.Team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: response.Team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: response.Team.AutoArchivePeriod}

	data.Triage = &triage{
		Enabled: types.Bool{Value: response.Team.TriageEnabled},
	}

	data.Cycles = &cycles{
		Enabled:          types.Bool{Value: response.Team.CyclesEnabled},
		StartDay:         types.Float64{Value: response.Team.CycleStartDay},
		Duration:         types.Float64{Value: response.Team.CycleDuration},
		Cooldown:         types.Float64{Value: response.Team.CycleCooldownTime},
		Upcoming:         types.Float64{Value: response.Team.UpcomingCycleCount},
		AutoAddStarted:   types.Bool{Value: response.Team.CycleIssueAutoAssignStarted},
		AutoAddCompleted: types.Bool{Value: response.Team.CycleIssueAutoAssignCompleted},
		NeedForActive:    types.Bool{Value: response.Team.CycleLockToActive},
	}

	data.Estimation = &estimation{
		Type:      types.String{Value: response.Team.IssueEstimationType},
		Extended:  types.Bool{Value: response.Team.IssueEstimationExtended},
		AllowZero: types.Bool{Value: response.Team.IssueEstimationAllowZero},
		Default:   types.Float64{Value: response.Team.DefaultIssueEstimate},
	}

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
		Name:                          data.Name.Value,
		Private:                       data.Private.Value,
		Description:                   data.Description.Value,
		Icon:                          data.Icon.Value,
		Color:                         data.Color.Value,
		Timezone:                      data.Timezone.Value,
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.Value,
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.Value,
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.Value,
		AutoArchivePeriod:             data.AutoArchivePeriod.Value,
		TriageEnabled:                 data.Triage.Enabled.Value,
		CyclesEnabled:                 data.Cycles.Enabled.Value,
		CycleStartDay:                 data.Cycles.StartDay.Value,
		CycleDuration:                 int(data.Cycles.Duration.Value),
		CycleCooldownTime:             int(data.Cycles.Cooldown.Value),
		UpcomingCycleCount:            data.Cycles.Upcoming.Value,
		CycleIssueAutoAssignStarted:   data.Cycles.AutoAddStarted.Value,
		CycleIssueAutoAssignCompleted: data.Cycles.AutoAddCompleted.Value,
		CycleLockToActive:             data.Cycles.NeedForActive.Value,
		CycleEnabledStartWeek:         "nextWeek",
		IssueEstimationType:           data.Estimation.Type.Value,
		IssueEstimationExtended:       data.Estimation.Extended.Value,
		IssueEstimationAllowZero:      data.Estimation.AllowZero.Value,
		DefaultIssueEstimate:          data.Estimation.Default.Value,
	}

	var key string

	diags = req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("key"), &key)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Key.Value != key {
		input.Key = data.Key.Value
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
	data.NoPriorityIssuesFirst = types.Bool{Value: response.TeamUpdate.Team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: response.TeamUpdate.Team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: response.TeamUpdate.Team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: response.TeamUpdate.Team.AutoArchivePeriod}

	data.Triage = &triage{
		Enabled: types.Bool{Value: response.TeamUpdate.Team.TriageEnabled},
	}

	data.Cycles = &cycles{
		Enabled:          types.Bool{Value: response.TeamUpdate.Team.CyclesEnabled},
		StartDay:         types.Float64{Value: response.TeamUpdate.Team.CycleStartDay},
		Duration:         types.Float64{Value: response.TeamUpdate.Team.CycleDuration},
		Cooldown:         types.Float64{Value: response.TeamUpdate.Team.CycleCooldownTime},
		Upcoming:         types.Float64{Value: response.TeamUpdate.Team.UpcomingCycleCount},
		AutoAddStarted:   types.Bool{Value: response.TeamUpdate.Team.CycleIssueAutoAssignStarted},
		AutoAddCompleted: types.Bool{Value: response.TeamUpdate.Team.CycleIssueAutoAssignCompleted},
		NeedForActive:    types.Bool{Value: response.TeamUpdate.Team.CycleLockToActive},
	}

	data.Estimation = &estimation{
		Type:      types.String{Value: response.TeamUpdate.Team.IssueEstimationType},
		Extended:  types.Bool{Value: response.TeamUpdate.Team.IssueEstimationExtended},
		AllowZero: types.Bool{Value: response.TeamUpdate.Team.IssueEstimationAllowZero},
		Default:   types.Float64{Value: response.TeamUpdate.Team.DefaultIssueEstimate},
	}

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
