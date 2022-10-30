package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/validators"
)

var _ resource.Resource = &TeamResource{}
var _ resource.ResourceWithImportState = &TeamResource{}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

type TeamResource struct {
	client *graphql.Client
}

type TeamResourceTriageModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type TeamResourceCyclesModel struct {
	Enabled          types.Bool    `tfsdk:"enabled"`
	StartDay         types.Float64 `tfsdk:"start_day"`
	Duration         types.Float64 `tfsdk:"duration"`
	Cooldown         types.Float64 `tfsdk:"cooldown"`
	Upcoming         types.Float64 `tfsdk:"upcoming"`
	AutoAddStarted   types.Bool    `tfsdk:"auto_add_started"`
	AutoAddCompleted types.Bool    `tfsdk:"auto_add_completed"`
	NeedForActive    types.Bool    `tfsdk:"need_for_active"`
}

type TeamResourceEstimationModel struct {
	Type      types.String  `tfsdk:"type"`
	Extended  types.Bool    `tfsdk:"extended"`
	AllowZero types.Bool    `tfsdk:"allow_zero"`
	Default   types.Float64 `tfsdk:"default"`
}

type TeamResourceModel struct {
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
	AutoClosePeriod            types.Float64 `tfsdk:"auto_close_period"`
	Triage                     types.Object  `tfsdk:"triage"`
	Cycles                     types.Object  `tfsdk:"cycles"`
	Estimation                 types.Object  `tfsdk:"estimation"`
}

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				MarkdownDescription: "Period after which closed and completed issues are automatically archived, in months. **Default** `6`.",
				Type:                types.Float64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultFloat(6),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.FloatInSlice(1, 3, 6, 9, 12),
				},
			},
			"auto_close_period": {
				MarkdownDescription: "Period after which non-completed or non-canceled issues are automatically closed, in months. **Default** `6`. *Use `0` for turning this off.*",
				Type:                types.Float64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultFloat(6),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.FloatInSlice(0, 1, 3, 6, 9, 12),
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

func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamResourceModel
	var triageData *TeamResourceTriageModel
	var cyclesData *TeamResourceCyclesModel
	var estimationData *TeamResourceEstimationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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
		AutoArchivePeriod:             data.AutoArchivePeriod.Value,
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

	if data.AutoClosePeriod.Value != 0 {
		input.AutoClosePeriod = &data.AutoClosePeriod.Value
	}

	resp.Diagnostics.Append(data.Triage.As(ctx, &triageData, types.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.Value

	resp.Diagnostics.Append(data.Cycles.As(ctx, &cyclesData, types.ObjectAsOptions{})...)

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

	resp.Diagnostics.Append(data.Estimation.As(ctx, &estimationData, types.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.IssueEstimationType = estimationData.Type.Value
	input.IssueEstimationExtended = estimationData.Extended.Value
	input.IssueEstimationAllowZero = estimationData.AllowZero.Value
	input.DefaultIssueEstimate = estimationData.Default.Value

	response, err := createTeam(context.Background(), *r.client, input)

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

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64{Value: *team.AutoClosePeriod}
	} else {
		data.AutoClosePeriod = types.Float64{Value: 0}
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeam(context.Background(), *r.client, data.Key.Value)

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

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64{Value: *team.AutoClosePeriod}
	} else {
		data.AutoClosePeriod = types.Float64{Value: 0}
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamResourceModel
	var triageData *TeamResourceTriageModel
	var cyclesData *TeamResourceCyclesModel
	var estimationData *TeamResourceEstimationModel

	var state *TeamResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TeamUpdateInput{
		Private:                       data.Private.Value,
		Timezone:                      data.Timezone.Value,
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.Value,
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.Value,
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.Value,
		AutoArchivePeriod:             data.AutoArchivePeriod.Value,
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

	if data.AutoClosePeriod.Value != 0 {
		input.AutoClosePeriod = &data.AutoClosePeriod.Value
	}

	resp.Diagnostics.Append(data.Triage.As(ctx, &triageData, types.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.Value

	resp.Diagnostics.Append(data.Cycles.As(ctx, &cyclesData, types.ObjectAsOptions{})...)

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

	resp.Diagnostics.Append(data.Estimation.As(ctx, &estimationData, types.ObjectAsOptions{})...)

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

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("key"), &key)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Key.Value != key {
		input.Key = data.Key.Value
	}

	response, err := updateTeam(context.Background(), *r.client, input, key)

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

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64{Value: *team.AutoClosePeriod}
	} else {
		data.AutoClosePeriod = types.Float64{Value: 0}
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteTeam(context.Background(), *r.client, data.Key.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team")
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}
