package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
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

type TeamResourceWorkflowStateModel struct {
	Id          types.String  `tfsdk:"id"`
	Position    types.Float64 `tfsdk:"position"`
	Name        types.String  `tfsdk:"name"`
	Color       types.String  `tfsdk:"color"`
	Description types.String  `tfsdk:"description"`
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
	BacklogWorkflowState       types.Object  `tfsdk:"backlog_workflow_state"`
	UnstartedWorkflowState     types.Object  `tfsdk:"unstarted_workflow_state"`
	StartedWorkflowState       types.Object  `tfsdk:"started_workflow_state"`
	CompletedWorkflowState     types.Object  `tfsdk:"completed_workflow_state"`
	CanceledWorkflowState      types.Object  `tfsdk:"canceled_workflow_state"`
}

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear team.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the team.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key of the team.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(5),
					stringvalidator.RegexMatches(regexp.MustCompile("^[A-Z0-9]+$"), "must only contain uppercase letters and numbers"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(2),
				},
			},
			"private": schema.BoolAttribute{
				MarkdownDescription: "Privacy of the team. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.NullableString(),
				},
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z]+$"), "must only contain letters"),
				},
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
				},
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "Timezone of the team. **Default** `Etc/GMT`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Etc/GMT"),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"no_priority_issues_first": schema.BoolAttribute{
				MarkdownDescription: "Prefer issues without priority at the top during issue prioritization order. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"enable_issue_history_grouping": schema.BoolAttribute{
				MarkdownDescription: "Enable issue history grouping for the team. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"enable_issue_default_to_bottom": schema.BoolAttribute{
				MarkdownDescription: "Enable moving issues to bottom of the column when changing state. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auto_archive_period": schema.Float64Attribute{
				MarkdownDescription: "Period after which closed and completed issues are automatically archived, in months. **Default** `6`.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(6),
				Validators: []validator.Float64{
					float64validator.OneOf([]float64{1, 3, 6, 9, 12}...),
				},
			},
			"auto_close_period": schema.Float64Attribute{
				MarkdownDescription: "Period after which non-completed or non-canceled issues are automatically closed, in months. **Default** `6`. *Use `0` for turning this off.*",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(6),
				Validators: []validator.Float64{
					float64validator.OneOf([]float64{0, 1, 3, 6, 9, 12}...),
				},
			},
			"triage": schema.SingleNestedAttribute{
				MarkdownDescription: "Triage settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable triage mode for the team. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(false),
						},
					},
				},
			},
			"cycles": schema.SingleNestedAttribute{
				MarkdownDescription: "Cycle settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable cycles for the team. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(false),
						},
					},
					"start_day": schema.Float64Attribute{
						MarkdownDescription: "Start day of the cycle. Sunday is 0, Saturday is 6. **Default** `0`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							modifiers.DefaultFloat(0),
						},
						Validators: []validator.Float64{
							float64validator.OneOf([]float64{0, 1, 2, 3, 4, 5, 6}...),
						},
					},
					"duration": schema.Float64Attribute{
						MarkdownDescription: "Duration of the cycle in weeks. **Default** `1`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							modifiers.DefaultFloat(1),
						},
						Validators: []validator.Float64{
							float64validator.OneOf([]float64{1, 2, 3, 4, 5, 6, 7, 8}...),
						},
					},
					"cooldown": schema.Float64Attribute{
						MarkdownDescription: "Cooldown time between cycles in weeks. **Default** `0`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							modifiers.DefaultFloat(0),
						},
						Validators: []validator.Float64{
							float64validator.OneOf([]float64{0, 1, 2, 3}...),
						},
					},
					"upcoming": schema.Float64Attribute{
						MarkdownDescription: "Number of upcoming cycles to automatically create. **Default** `2`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							modifiers.DefaultFloat(2),
						},
						Validators: []validator.Float64{
							float64validator.OneOf([]float64{1, 2, 3, 4, 6, 8, 10}...),
						},
					},
					"auto_add_started": schema.BoolAttribute{
						MarkdownDescription: "Auto add started issues that don't belong to any cycle to the active cycle. **Default** `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(true),
						},
					},
					"auto_add_completed": schema.BoolAttribute{
						MarkdownDescription: "Auto add completed issues that don't belong to any cycle to the active cycle. **Default** `true`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(true),
						},
					},
					"need_for_active": schema.BoolAttribute{
						MarkdownDescription: "Whether all active issues need to have a cycle. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(false),
						},
					},
				},
			},
			"estimation": schema.SingleNestedAttribute{
				MarkdownDescription: "Issue estimation settings of the team.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Issue estimation type for the team. **Default** `notUsed`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("notUsed"),
						},
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"notUsed", "exponential", "fibonacci", "linear", "tShirt"}...),
						},
					},
					"extended": schema.BoolAttribute{
						MarkdownDescription: "Whether the team uses extended estimation. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(false),
						},
					},
					"allow_zero": schema.BoolAttribute{
						MarkdownDescription: "Whether zero is allowed as an estimation. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							modifiers.DefaultBool(false),
						},
					},
					"default": schema.Float64Attribute{
						MarkdownDescription: "Default estimation for issues that are unestimated. **Default** `1`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							modifiers.DefaultFloat(1),
						},
						Validators: []validator.Float64{
							float64validator.OneOf([]float64{0, 1}...),
						},
					},
				},
			},
			"backlog_workflow_state": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the `backlog` workflow state that is created by default for the team. *Position is always `0`. This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the workflow state. **Default** `Backlog`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("Backlog"),
						},
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the workflow state. **Default** `#bec2c8`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("#bec2c8"),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the workflow state.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.NullableString(),
						},
					},
				},
			},
			"unstarted_workflow_state": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the `unstarted` workflow state that is created by default for the team. *Position is always `0`. This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the workflow state. **Default** `Todo`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("Todo"),
						},
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the workflow state. **Default** `#e2e2e2`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("#e2e2e2"),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the workflow state.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.NullableString(),
						},
					},
				},
			},
			"started_workflow_state": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the `started` workflow state that is created by default for the team. *Position is always `0`. This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the workflow state. **Default** `In Progress`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("In Progress"),
						},
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the workflow state. **Default** `#f2c94c`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("#f2c94c"),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the workflow state.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.NullableString(),
						},
					},
				},
			},
			"completed_workflow_state": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the `completed` workflow state that is created by default for the team. *Position is always `0`. This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the workflow state. **Default** `Done`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("Done"),
						},
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the workflow state. **Default** `#5e6ad2`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("#5e6ad2"),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the workflow state.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.NullableString(),
						},
					},
				},
			},
			"canceled_workflow_state": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the `canceled` workflow state that is created by default for the team. *Position is always `0`. This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					modifiers.UnknownAttributesOnUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the workflow state.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the workflow state. **Default** `Canceled`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("Canceled"),
						},
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the workflow state. **Default** `#95a2b3`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("#95a2b3"),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the workflow state.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.NullableString(),
						},
					},
				},
			},
		},
	}
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
		Key:                           data.Key.ValueString(),
		Name:                          data.Name.ValueString(),
		Private:                       data.Private.ValueBool(),
		Timezone:                      data.Timezone.ValueString(),
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.ValueBool(),
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.ValueBool(),
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.ValueBool(),
		AutoArchivePeriod:             data.AutoArchivePeriod.ValueFloat64(),
	}

	if !data.Description.IsNull() {
		value := data.Description.ValueString()
		input.Description = &value
	}

	if !data.Icon.IsUnknown() {
		value := data.Icon.ValueString()
		input.Icon = &value
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	if data.AutoClosePeriod.ValueFloat64() != 0 {
		value := data.AutoClosePeriod.ValueFloat64()
		input.AutoClosePeriod = &value
	}

	resp.Diagnostics.Append(data.Triage.As(ctx, &triageData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.ValueBool()

	resp.Diagnostics.Append(data.Cycles.As(ctx, &cyclesData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.CyclesEnabled = cyclesData.Enabled.ValueBool()
	input.CycleStartDay = cyclesData.StartDay.ValueFloat64()
	input.CycleDuration = int(cyclesData.Duration.ValueFloat64())
	input.CycleCooldownTime = int(cyclesData.Cooldown.ValueFloat64())
	input.UpcomingCycleCount = cyclesData.Upcoming.ValueFloat64()
	input.CycleIssueAutoAssignStarted = cyclesData.AutoAddStarted.ValueBool()
	input.CycleIssueAutoAssignCompleted = cyclesData.AutoAddCompleted.ValueBool()
	input.CycleLockToActive = cyclesData.NeedForActive.ValueBool()

	resp.Diagnostics.Append(data.Estimation.As(ctx, &estimationData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.IssueEstimationType = estimationData.Type.ValueString()
	input.IssueEstimationExtended = estimationData.Extended.ValueBool()
	input.IssueEstimationAllowZero = estimationData.AllowZero.ValueBool()
	input.DefaultIssueEstimate = estimationData.Default.ValueFloat64()

	response, err := createTeam(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team")

	team := response.TeamCreate.Team

	data.Id = types.StringValue(team.Id)
	data.Private = types.BoolValue(team.Private)
	data.Timezone = types.StringValue(team.Timezone)
	data.NoPriorityIssuesFirst = types.BoolValue(team.IssueOrderingNoPriorityFirst)
	data.EnableIssueHistoryGrouping = types.BoolValue(team.GroupIssueHistory)
	data.EnableIssueDefaultToBottom = types.BoolValue(team.IssueSortOrderDefaultToBottom)
	data.AutoArchivePeriod = types.Float64Value(team.AutoArchivePeriod)

	if team.Description != nil {
		data.Description = types.StringValue(*team.Description)
	}

	if team.Icon != nil {
		data.Icon = types.StringValue(*team.Icon)
	}

	if team.Color != nil {
		data.Color = types.StringValue(*team.Color)
	}

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64Value(*team.AutoClosePeriod)
	} else {
		data.AutoClosePeriod = types.Float64Value(0)
	}

	data.Triage = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled": types.BoolType,
		},
		map[string]attr.Value{
			"enabled": types.BoolValue(team.TriageEnabled),
		},
	)

	data.Cycles = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		map[string]attr.Value{
			"enabled":            types.BoolValue(team.CyclesEnabled),
			"start_day":          types.Float64Value(team.CycleStartDay),
			"duration":           types.Float64Value(team.CycleDuration),
			"cooldown":           types.Float64Value(team.CycleCooldownTime),
			"upcoming":           types.Float64Value(team.UpcomingCycleCount),
			"auto_add_started":   types.BoolValue(team.CycleIssueAutoAssignStarted),
			"auto_add_completed": types.BoolValue(team.CycleIssueAutoAssignCompleted),
			"need_for_active":    types.BoolValue(team.CycleLockToActive),
		},
	)

	data.Estimation = types.ObjectValueMust(
		map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		map[string]attr.Value{
			"type":       types.StringValue(team.IssueEstimationType),
			"extended":   types.BoolValue(team.IssueEstimationExtended),
			"allow_zero": types.BoolValue(team.IssueEstimationAllowZero),
			"default":    types.Float64Value(team.DefaultIssueEstimate),
		},
	)

	// Read the workflow states so that we can update them

	workflowStatesResponse, workflowStatesErr := getTeamWorkflowStates(ctx, *r.client, team.Key)

	if workflowStatesErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team workflow states, got error: %s", workflowStatesErr))
		return
	}

	tflog.Trace(ctx, "read team workflow states")

	backlogWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "backlog")
	unstartedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "unstarted")
	startedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "started")
	completedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "completed")
	canceledWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "canceled")

	if backlogWorkflowState == nil || unstartedWorkflowState == nil || startedWorkflowState == nil || completedWorkflowState == nil || canceledWorkflowState == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to find all workflow states in a new team")
		return
	}

	// Update the workflow states

	backlog := updateTeamWorkflowStateInCreate(ctx, r, data.BacklogWorkflowState, resp, backlogWorkflowState.Id)
	unstarted := updateTeamWorkflowStateInCreate(ctx, r, data.UnstartedWorkflowState, resp, unstartedWorkflowState.Id)
	started := updateTeamWorkflowStateInCreate(ctx, r, data.StartedWorkflowState, resp, startedWorkflowState.Id)
	completed := updateTeamWorkflowStateInCreate(ctx, r, data.CompletedWorkflowState, resp, completedWorkflowState.Id)
	canceled := updateTeamWorkflowStateInCreate(ctx, r, data.CanceledWorkflowState, resp, canceledWorkflowState.Id)

	if backlog == nil || unstarted == nil || started == nil || completed == nil || canceled == nil {
		return
	}

	data.BacklogWorkflowState = *backlog
	data.UnstartedWorkflowState = *unstarted
	data.StartedWorkflowState = *started
	data.CompletedWorkflowState = *completed
	data.CanceledWorkflowState = *canceled

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeam(ctx, *r.client, data.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	team := response.Team

	data.Id = types.StringValue(team.Id)
	data.Name = types.StringValue(team.Name)
	data.Private = types.BoolValue(team.Private)
	data.Timezone = types.StringValue(team.Timezone)
	data.NoPriorityIssuesFirst = types.BoolValue(team.IssueOrderingNoPriorityFirst)
	data.EnableIssueHistoryGrouping = types.BoolValue(team.GroupIssueHistory)
	data.EnableIssueDefaultToBottom = types.BoolValue(team.IssueSortOrderDefaultToBottom)
	data.AutoArchivePeriod = types.Float64Value(team.AutoArchivePeriod)

	if team.Description != nil {
		data.Description = types.StringValue(*team.Description)
	}

	if team.Icon != nil {
		data.Icon = types.StringValue(*team.Icon)
	}

	if team.Color != nil {
		data.Color = types.StringValue(*team.Color)
	}

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64Value(*team.AutoClosePeriod)
	} else {
		data.AutoClosePeriod = types.Float64Value(0)
	}

	data.Triage = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled": types.BoolType,
		},
		map[string]attr.Value{
			"enabled": types.BoolValue(team.TriageEnabled),
		},
	)

	data.Cycles = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		map[string]attr.Value{
			"enabled":            types.BoolValue(team.CyclesEnabled),
			"start_day":          types.Float64Value(team.CycleStartDay),
			"duration":           types.Float64Value(team.CycleDuration),
			"cooldown":           types.Float64Value(team.CycleCooldownTime),
			"upcoming":           types.Float64Value(team.UpcomingCycleCount),
			"auto_add_started":   types.BoolValue(team.CycleIssueAutoAssignStarted),
			"auto_add_completed": types.BoolValue(team.CycleIssueAutoAssignCompleted),
			"need_for_active":    types.BoolValue(team.CycleLockToActive),
		},
	)

	data.Estimation = types.ObjectValueMust(
		map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		map[string]attr.Value{
			"type":       types.StringValue(team.IssueEstimationType),
			"extended":   types.BoolValue(team.IssueEstimationExtended),
			"allow_zero": types.BoolValue(team.IssueEstimationAllowZero),
			"default":    types.Float64Value(team.DefaultIssueEstimate),
		},
	)

	workflowStatesResponse, workflowStatesErr := getTeamWorkflowStates(ctx, *r.client, team.Key)

	if workflowStatesErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team workflow states, got error: %s", workflowStatesErr))
		return
	}

	tflog.Trace(ctx, "read team workflow states")

	backlogWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "backlog")
	unstartedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "unstarted")
	startedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "started")
	completedWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "completed")
	canceledWorkflowState := findWorkflowStateType(workflowStatesResponse.WorkflowStates.Nodes, "canceled")

	if backlogWorkflowState == nil || unstartedWorkflowState == nil || startedWorkflowState == nil || completedWorkflowState == nil || canceledWorkflowState == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to find all workflow states when reading team")
		return
	}

	data.BacklogWorkflowState = readWorkflowStateToObject(*backlogWorkflowState)
	data.UnstartedWorkflowState = readWorkflowStateToObject(*unstartedWorkflowState)
	data.StartedWorkflowState = readWorkflowStateToObject(*startedWorkflowState)
	data.CompletedWorkflowState = readWorkflowStateToObject(*completedWorkflowState)
	data.CanceledWorkflowState = readWorkflowStateToObject(*canceledWorkflowState)

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
		Private:                       data.Private.ValueBool(),
		Timezone:                      data.Timezone.ValueString(),
		IssueOrderingNoPriorityFirst:  data.NoPriorityIssuesFirst.ValueBool(),
		GroupIssueHistory:             data.EnableIssueHistoryGrouping.ValueBool(),
		IssueSortOrderDefaultToBottom: data.EnableIssueDefaultToBottom.ValueBool(),
		AutoArchivePeriod:             data.AutoArchivePeriod.ValueFloat64(),
	}

	if data.Key.ValueString() != state.Key.ValueString() {
		input.Key = data.Key.ValueString()
	}

	if data.Name.ValueString() != state.Name.ValueString() {
		input.Name = data.Name.ValueString()
	}

	if !data.Description.IsNull() {
		value := data.Description.ValueString()
		input.Description = &value
	}

	if !data.Icon.IsUnknown() {
		value := data.Icon.ValueString()
		input.Icon = &value
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	if data.AutoClosePeriod.ValueFloat64() != 0 {
		value := data.AutoClosePeriod.ValueFloat64()
		input.AutoClosePeriod = &value
	}

	resp.Diagnostics.Append(data.Triage.As(ctx, &triageData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.TriageEnabled = triageData.Enabled.ValueBool()

	resp.Diagnostics.Append(data.Cycles.As(ctx, &cyclesData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.CyclesEnabled = cyclesData.Enabled.ValueBool()
	input.CycleStartDay = cyclesData.StartDay.ValueFloat64()
	input.CycleDuration = int(cyclesData.Duration.ValueFloat64())
	input.CycleCooldownTime = int(cyclesData.Cooldown.ValueFloat64())
	input.UpcomingCycleCount = cyclesData.Upcoming.ValueFloat64()
	input.CycleIssueAutoAssignStarted = cyclesData.AutoAddStarted.ValueBool()
	input.CycleIssueAutoAssignCompleted = cyclesData.AutoAddCompleted.ValueBool()
	input.CycleLockToActive = cyclesData.NeedForActive.ValueBool()

	resp.Diagnostics.Append(data.Estimation.As(ctx, &estimationData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.IssueEstimationType = estimationData.Type.ValueString()
	input.IssueEstimationExtended = estimationData.Extended.ValueBool()
	input.IssueEstimationAllowZero = estimationData.AllowZero.ValueBool()
	input.DefaultIssueEstimate = estimationData.Default.ValueFloat64()

	if input.CyclesEnabled {
		input.CycleEnabledStartWeek = "nextWeek"
	}

	response, err := updateTeam(ctx, *r.client, input, state.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team")

	team := response.TeamUpdate.Team

	data.Id = types.StringValue(team.Id)
	data.Private = types.BoolValue(team.Private)
	data.Timezone = types.StringValue(team.Timezone)
	data.NoPriorityIssuesFirst = types.BoolValue(team.IssueOrderingNoPriorityFirst)
	data.EnableIssueHistoryGrouping = types.BoolValue(team.GroupIssueHistory)
	data.EnableIssueDefaultToBottom = types.BoolValue(team.IssueSortOrderDefaultToBottom)
	data.AutoArchivePeriod = types.Float64Value(team.AutoArchivePeriod)

	if team.Description != nil {
		data.Description = types.StringValue(*team.Description)
	}

	if team.Icon != nil {
		data.Icon = types.StringValue(*team.Icon)
	}

	if team.Color != nil {
		data.Color = types.StringValue(*team.Color)
	}

	if team.AutoClosePeriod != nil {
		data.AutoClosePeriod = types.Float64Value(*team.AutoClosePeriod)
	} else {
		data.AutoClosePeriod = types.Float64Value(0)
	}

	data.Triage = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled": types.BoolType,
		},
		map[string]attr.Value{
			"enabled": types.BoolValue(team.TriageEnabled),
		},
	)

	data.Cycles = types.ObjectValueMust(
		map[string]attr.Type{
			"enabled":            types.BoolType,
			"start_day":          types.Float64Type,
			"duration":           types.Float64Type,
			"cooldown":           types.Float64Type,
			"upcoming":           types.Float64Type,
			"auto_add_started":   types.BoolType,
			"auto_add_completed": types.BoolType,
			"need_for_active":    types.BoolType,
		},
		map[string]attr.Value{
			"enabled":            types.BoolValue(team.CyclesEnabled),
			"start_day":          types.Float64Value(team.CycleStartDay),
			"duration":           types.Float64Value(team.CycleDuration),
			"cooldown":           types.Float64Value(team.CycleCooldownTime),
			"upcoming":           types.Float64Value(team.UpcomingCycleCount),
			"auto_add_started":   types.BoolValue(team.CycleIssueAutoAssignStarted),
			"auto_add_completed": types.BoolValue(team.CycleIssueAutoAssignCompleted),
			"need_for_active":    types.BoolValue(team.CycleLockToActive),
		},
	)

	data.Estimation = types.ObjectValueMust(
		map[string]attr.Type{
			"type":       types.StringType,
			"extended":   types.BoolType,
			"allow_zero": types.BoolType,
			"default":    types.Float64Type,
		},
		map[string]attr.Value{
			"type":       types.StringValue(team.IssueEstimationType),
			"extended":   types.BoolValue(team.IssueEstimationExtended),
			"allow_zero": types.BoolValue(team.IssueEstimationAllowZero),
			"default":    types.Float64Value(team.DefaultIssueEstimate),
		},
	)

	// Update the workflow states

	backlog := updateTeamWorkflowStateInUpdate(ctx, r, data.BacklogWorkflowState, resp, state.BacklogWorkflowState.Attributes()["id"].(types.String).ValueString())
	unstarted := updateTeamWorkflowStateInUpdate(ctx, r, data.UnstartedWorkflowState, resp, state.UnstartedWorkflowState.Attributes()["id"].(types.String).ValueString())
	started := updateTeamWorkflowStateInUpdate(ctx, r, data.StartedWorkflowState, resp, state.StartedWorkflowState.Attributes()["id"].(types.String).ValueString())
	completed := updateTeamWorkflowStateInUpdate(ctx, r, data.CompletedWorkflowState, resp, state.CompletedWorkflowState.Attributes()["id"].(types.String).ValueString())
	canceled := updateTeamWorkflowStateInUpdate(ctx, r, data.CanceledWorkflowState, resp, state.CanceledWorkflowState.Attributes()["id"].(types.String).ValueString())

	if backlog == nil || unstarted == nil || started == nil || completed == nil || canceled == nil {
		return
	}

	data.BacklogWorkflowState = *backlog
	data.UnstartedWorkflowState = *unstarted
	data.StartedWorkflowState = *started
	data.CompletedWorkflowState = *completed
	data.CanceledWorkflowState = *canceled

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteTeam(ctx, *r.client, data.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team")
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}

func findWorkflowStateType(workflowStates []getTeamWorkflowStatesWorkflowStatesWorkflowStateConnectionNodesWorkflowState, ty string) *getTeamWorkflowStatesWorkflowStatesWorkflowStateConnectionNodesWorkflowState {
	for _, workflowState := range workflowStates {
		if workflowState.Type == ty && workflowState.Position == 0 {
			return &workflowState
		}
	}

	// If unable to find the exact workflow with the exact position, return the first one with the correct type
	sort.Slice(workflowStates, func(i, j int) bool {
		return workflowStates[i].Position < workflowStates[j].Position
	})

	for _, workflowState := range workflowStates {
		if workflowState.Type == ty {
			return &workflowState
		}
	}

	return nil
}

func readWorkflowStateToObject(workflowState getTeamWorkflowStatesWorkflowStatesWorkflowStateConnectionNodesWorkflowState) types.Object {
	attrs := map[string]attr.Value{
		"id":          types.StringValue(workflowState.Id),
		"position":    types.Float64Value(workflowState.Position),
		"name":        types.StringValue(workflowState.Name),
		"color":       types.StringValue(workflowState.Color),
		"description": types.StringNull(),
	}

	if workflowState.Description != nil {
		attrs["description"] = types.StringValue(*workflowState.Description)
	}

	ret := types.ObjectValueMust(
		map[string]attr.Type{
			"id":          types.StringType,
			"position":    types.Float64Type,
			"name":        types.StringType,
			"color":       types.StringType,
			"description": types.StringType,
		},
		attrs,
	)

	return ret
}

func updateWorkflowStateToObject(workflowState updateWorkflowStateWorkflowStateUpdateWorkflowStatePayloadWorkflowState) types.Object {
	attrs := map[string]attr.Value{
		"id":          types.StringValue(workflowState.Id),
		"position":    types.Float64Value(workflowState.Position),
		"name":        types.StringValue(workflowState.Name),
		"color":       types.StringValue(workflowState.Color),
		"description": types.StringNull(),
	}

	if workflowState.Description != nil {
		attrs["description"] = types.StringValue(*workflowState.Description)
	}

	ret := types.ObjectValueMust(
		map[string]attr.Type{
			"id":          types.StringType,
			"position":    types.Float64Type,
			"name":        types.StringType,
			"color":       types.StringType,
			"description": types.StringType,
		},
		attrs,
	)

	return ret
}

func updateTeamWorkflowStateInCreate(ctx context.Context, r *TeamResource, data types.Object, resp *resource.CreateResponse, id string) *types.Object {
	var workflowStateData *TeamResourceWorkflowStateModel

	resp.Diagnostics.Append(data.As(ctx, &workflowStateData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return nil
	}

	workflowStateInput := WorkflowStateUpdateInput{
		Name:  workflowStateData.Name.ValueString(),
		Color: workflowStateData.Color.ValueString(),
	}

	if !workflowStateData.Description.IsNull() {
		value := workflowStateData.Description.ValueString()
		workflowStateInput.Description = &value
	}

	workflowStateResponse, workflowStateErr := updateWorkflowState(ctx, *r.client, workflowStateInput, id)

	if workflowStateErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow state, got error: %s", workflowStateErr))
		return nil
	}

	ret := updateWorkflowStateToObject(workflowStateResponse.WorkflowStateUpdate.WorkflowState)

	return &ret
}

func updateTeamWorkflowStateInUpdate(ctx context.Context, r *TeamResource, data types.Object, resp *resource.UpdateResponse, id string) *types.Object {
	var workflowStateData *TeamResourceWorkflowStateModel

	resp.Diagnostics.Append(data.As(ctx, &workflowStateData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return nil
	}

	workflowStateInput := WorkflowStateUpdateInput{
		Name:  workflowStateData.Name.ValueString(),
		Color: workflowStateData.Color.ValueString(),
	}

	if !workflowStateData.Description.IsNull() {
		value := workflowStateData.Description.ValueString()
		workflowStateInput.Description = &value
	}

	workflowStateResponse, workflowStateErr := updateWorkflowState(ctx, *r.client, workflowStateInput, id)

	if workflowStateErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow state, got error: %s", workflowStateErr))
		return nil
	}

	ret := updateWorkflowStateToObject(workflowStateResponse.WorkflowStateUpdate.WorkflowState)

	return &ret
}
