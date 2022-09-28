package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
var _ provider.ResourceType = teamResourceType{}
var _ resource.Resource = teamResource{}
var _ resource.ResourceWithImportState = teamResource{}

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
					resource.UseStateForUnknown(),
				},
			},
			"key": {
				MarkdownDescription: "Key of the team.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.MaxLength(5),
					validators.Match(regexp.MustCompile("^[A-Z0-9]+$")),
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
					modifiers.DefaultBool(false),
				},
			},
			"description": {
				MarkdownDescription: "Description of the team.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
			},
			"icon": {
				MarkdownDescription: "Icon of the team.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
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
					resource.UseStateForUnknown(),
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
					modifiers.DefaultString("Etc/GMT"),
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
					modifiers.DefaultBool(true),
				},
			},
			"enable_issue_history_grouping": {
				MarkdownDescription: "Enable issue history grouping for the team. **Default** `true`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(true),
				},
			},
			"enable_issue_default_to_bottom": {
				MarkdownDescription: "Enable moving issues to bottom of the column when changing state. **Default** `false`.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(false),
				},
			},
			"auto_archive_period": {
				MarkdownDescription: "Period after which closed and completed issues are automatically archived, in months. **Default** `3`.",
				Type:                types.Float64Type,
				// #2
				// Optional:            true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultFloat(3),
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
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enabled": {
						MarkdownDescription: "Enable triage mode for the team. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultBool(false),
						},
					},
				}),
			},
			"cycles": {
				MarkdownDescription: "Cycle settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enabled": {
						MarkdownDescription: "Enable cycles for the team. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultBool(false),
						},
					},
					"start_day": {
						MarkdownDescription: "Start day of the cycle. Sunday is 0, Saturday is 6. **Default** `0`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultFloat(0),
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
							modifiers.DefaultFloat(1),
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
							modifiers.DefaultFloat(0),
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
							modifiers.DefaultFloat(2),
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
							modifiers.DefaultBool(true),
						},
					},
					"auto_add_completed": {
						MarkdownDescription: "Auto add completed issues that don't belong to any cycle to the active cycle. **Default** `true`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultBool(true),
						},
					},
					"need_for_active": {
						MarkdownDescription: "Whether all active issues need to have a cycle. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultBool(false),
						},
					},
				}),
			},
			"estimation": {
				MarkdownDescription: "Issue estimation settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"type": {
						MarkdownDescription: "Issue estimation type for the team. **Default** `notUsed`.",
						Type:                types.StringType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultString("notUsed"),
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
							modifiers.DefaultBool(false),
						},
					},
					"allow_zero": {
						MarkdownDescription: "Whether zero is allowed as an estimation. **Default** `false`.",
						Type:                types.BoolType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultBool(false),
						},
					},
					"default": {
						MarkdownDescription: "Default estimation for issues that are unestimated. **Default** `1`.",
						Type:                types.Float64Type,
						Optional:            true,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							modifiers.DefaultFloat(1),
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

func (t teamResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
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
	Triage                     types.Object  `tfsdk:"triage"`
	Cycles                     types.Object  `tfsdk:"cycles"`
	Estimation                 types.Object  `tfsdk:"estimation"`
}

type teamResource struct {
	provider linearProvider
}

func (r teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data teamResourceData
	var triageData triage
	var cyclesData cycles
	var estimationData estimation

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TeamCreateInput{
		Key:                           data.Key.Value,
		Name:                          data.Name.Value,
		Private:                       data.Private.Value,
		Timezone:                      data.Timezone.Value,
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.Value,
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.Value,
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.Value,
		// #2
		// AutoArchivePeriod:             data.AutoArchivePeriod.Value,
		AutoArchivePeriod: 3,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	if !data.Icon.IsUnknown() {
		input.Icon = &data.Icon.Value
	}

	if !data.Color.IsUnknown() {
		input.Color = &data.Color.Value
	}

	diags = data.Triage.As(ctx, &triageData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.Value

	diags = data.Cycles.As(ctx, &cyclesData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.CyclesEnabled = cyclesData.Enabled.Value
	input.CycleStartDay = cyclesData.StartDay.Value
	input.CycleDuration = int(cyclesData.Duration.Value)
	input.CycleCooldownTime = int(cyclesData.Cooldown.Value)
	input.UpcomingCycleCount = cyclesData.Upcoming.Value
	input.CycleIssueAutoAssignStarted = cyclesData.AutoAddStarted.Value
	input.CycleIssueAutoAssignCompleted = cyclesData.AutoAddCompleted.Value
	input.CycleLockToActive = cyclesData.NeedForActive.Value

	diags = data.Estimation.As(ctx, &estimationData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.IssueEstimationType = estimationData.Type.Value
	input.IssueEstimationExtended = estimationData.Extended.Value
	input.IssueEstimationAllowZero = estimationData.AllowZero.Value
	input.DefaultIssueEstimate = estimationData.Default.Value

	response, err := createTeam(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team")

	team := response.TeamCreate.Team

	data.Id = types.String{Value: team.Id}
	data.Private = types.Bool{Value: team.Private}
	data.Timezone = types.String{Value: team.Timezone}
	data.NoPriorityIssuesFirst = types.Bool{Value: team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: team.AutoArchivePeriod}

	if team.Description != nil {
		data.Description = types.String{Value: *team.Description}
	}

	if team.Icon != nil {
		data.Icon = types.String{Value: *team.Icon}
	}

	if team.Color != nil {
		data.Color = types.String{Value: *team.Color}
	}

	data.Triage = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled": types.Bool{Value: team.TriageEnabled},
		},
	}

	data.Cycles = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled":            types.Bool{Value: team.CyclesEnabled},
			"start_day":          types.Float64{Value: team.CycleStartDay},
			"duration":           types.Float64{Value: team.CycleDuration},
			"cooldown":           types.Float64{Value: team.CycleCooldownTime},
			"upcoming":           types.Float64{Value: team.UpcomingCycleCount},
			"auto_add_started":   types.Bool{Value: team.CycleIssueAutoAssignStarted},
			"auto_add_completed": types.Bool{Value: team.CycleIssueAutoAssignCompleted},
			"need_for_active":    types.Bool{Value: team.CycleLockToActive},
		},
	}

	data.Estimation = types.Object{
		AttrTypes: map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		Attrs: map[string]attr.Value{
			"type":       types.String{Value: team.IssueEstimationType},
			"extended":   types.Bool{Value: team.IssueEstimationExtended},
			"allow_zero": types.Bool{Value: team.IssueEstimationAllowZero},
			"default":    types.Float64{Value: team.DefaultIssueEstimate},
		},
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	team := response.Team

	data.Id = types.String{Value: team.Id}
	data.Name = types.String{Value: team.Name}
	data.Private = types.Bool{Value: team.Private}
	data.Timezone = types.String{Value: team.Timezone}
	data.NoPriorityIssuesFirst = types.Bool{Value: team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: team.AutoArchivePeriod}

	if team.Description != nil {
		data.Description = types.String{Value: *team.Description}
	}

	if team.Icon != nil {
		data.Icon = types.String{Value: *team.Icon}
	}

	if team.Color != nil {
		data.Color = types.String{Value: *team.Color}
	}

	data.Triage = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled": types.Bool{Value: team.TriageEnabled},
		},
	}

	data.Cycles = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled":            types.Bool{Value: team.CyclesEnabled},
			"start_day":          types.Float64{Value: team.CycleStartDay},
			"duration":           types.Float64{Value: team.CycleDuration},
			"cooldown":           types.Float64{Value: team.CycleCooldownTime},
			"upcoming":           types.Float64{Value: team.UpcomingCycleCount},
			"auto_add_started":   types.Bool{Value: team.CycleIssueAutoAssignStarted},
			"auto_add_completed": types.Bool{Value: team.CycleIssueAutoAssignCompleted},
			"need_for_active":    types.Bool{Value: team.CycleLockToActive},
		},
	}

	data.Estimation = types.Object{
		AttrTypes: map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		Attrs: map[string]attr.Value{
			"type":       types.String{Value: team.IssueEstimationType},
			"extended":   types.Bool{Value: team.IssueEstimationExtended},
			"allow_zero": types.Bool{Value: team.IssueEstimationAllowZero},
			"default":    types.Float64{Value: team.DefaultIssueEstimate},
		},
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data teamResourceData
	var triageData triage
	var cyclesData cycles
	var estimationData estimation

	var state teamResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TeamUpdateInput{
		Private:                       data.Private.Value,
		Timezone:                      data.Timezone.Value,
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.Value,
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.Value,
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.Value,
		// #2
		// AutoArchivePeriod:             data.AutoArchivePeriod.Value,
		AutoArchivePeriod: 3,
	}

	if data.Name.Value != state.Name.Value {
		input.Name = data.Name.Value
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	if !data.Icon.IsUnknown() {
		input.Icon = &data.Icon.Value
	}

	if !data.Color.IsUnknown() {
		input.Color = &data.Color.Value
	}

	diags = data.Triage.As(ctx, &triageData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.Value

	diags = data.Cycles.As(ctx, &cyclesData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.CyclesEnabled = cyclesData.Enabled.Value
	input.CycleStartDay = cyclesData.StartDay.Value
	input.CycleDuration = int(cyclesData.Duration.Value)
	input.CycleCooldownTime = int(cyclesData.Cooldown.Value)
	input.UpcomingCycleCount = cyclesData.Upcoming.Value
	input.CycleIssueAutoAssignStarted = cyclesData.AutoAddStarted.Value
	input.CycleIssueAutoAssignCompleted = cyclesData.AutoAddCompleted.Value
	input.CycleLockToActive = cyclesData.NeedForActive.Value

	diags = data.Estimation.As(ctx, &estimationData, types.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.IssueEstimationType = estimationData.Type.Value
	input.IssueEstimationExtended = estimationData.Extended.Value
	input.IssueEstimationAllowZero = estimationData.AllowZero.Value
	input.DefaultIssueEstimate = estimationData.Default.Value

	if input.CyclesEnabled {
		input.CycleEnabledStartWeek = "nextWeek"
	}

	var key string

	diags = req.State.GetAttribute(ctx, path.Root("key"), &key)
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

	team := response.TeamUpdate.Team

	data.Id = types.String{Value: team.Id}
	data.Private = types.Bool{Value: team.Private}
	data.Timezone = types.String{Value: team.Timezone}
	data.NoPriorityIssuesFirst = types.Bool{Value: team.IssueOrderingNoPriorityFirst}
	data.EnableIssueHistoryGrouping = types.Bool{Value: team.GroupIssueHistory}
	data.EnableIssueDefaultToBottom = types.Bool{Value: team.IssueSortOrderDefaultToBottom}
	data.AutoArchivePeriod = types.Float64{Value: team.AutoArchivePeriod}

	if team.Description != nil {
		data.Description = types.String{Value: *team.Description}
	}

	if team.Icon != nil {
		data.Icon = types.String{Value: *team.Icon}
	}

	if team.Color != nil {
		data.Color = types.String{Value: *team.Color}
	}

	data.Triage = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled": types.Bool{Value: team.TriageEnabled},
		},
	}

	data.Cycles = types.Object{
		AttrTypes: map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		Attrs: map[string]attr.Value{
			"enabled":            types.Bool{Value: team.CyclesEnabled},
			"start_day":          types.Float64{Value: team.CycleStartDay},
			"duration":           types.Float64{Value: team.CycleDuration},
			"cooldown":           types.Float64{Value: team.CycleCooldownTime},
			"upcoming":           types.Float64{Value: team.UpcomingCycleCount},
			"auto_add_started":   types.Bool{Value: team.CycleIssueAutoAssignStarted},
			"auto_add_completed": types.Bool{Value: team.CycleIssueAutoAssignCompleted},
			"need_for_active":    types.Bool{Value: team.CycleLockToActive},
		},
	}

	data.Estimation = types.Object{
		AttrTypes: map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		Attrs: map[string]attr.Value{
			"type":       types.String{Value: team.IssueEstimationType},
			"extended":   types.Bool{Value: team.IssueEstimationExtended},
			"allow_zero": types.Bool{Value: team.IssueEstimationAllowZero},
			"default":    types.Float64{Value: team.DefaultIssueEstimate},
		},
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}
