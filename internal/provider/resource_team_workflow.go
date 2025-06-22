package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TeamWorkflowResource{}
var _ resource.ResourceWithImportState = &TeamWorkflowResource{}

func NewTeamWorkflowResource() resource.Resource {
	return &TeamWorkflowResource{}
}

type TeamWorkflowResource struct {
	client *graphql.Client
}

type TargetBranch struct {
	Id            string `json:"id"`
	BranchPattern string `json:"branchPattern"`
	IsRegex       bool   `json:"isRegex"`
}

type TeamWorkflowResourceBranchModel struct {
	Id      types.String `tfsdk:"id"`
	Pattern types.String `tfsdk:"pattern"`
	IsRegex types.Bool   `tfsdk:"is_regex"`
}

var branchAttrTypes = map[string]attr.Type{
	"id":       types.StringType,
	"pattern":  types.StringType,
	"is_regex": types.BoolType,
}

type TeamWorkflowResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Branch    types.Object `tfsdk:"branch"`
	Draft     types.String `tfsdk:"draft"`
	Start     types.String `tfsdk:"start"`
	Review    types.String `tfsdk:"review"`
	Mergeable types.String `tfsdk:"mergeable"`
	Merge     types.String `tfsdk:"merge"`
}

func (r *TeamWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_workflow"
}

func (r *TeamWorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear team workflow.",
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
			"branch": schema.SingleNestedAttribute{
				MarkdownDescription: "Branch settings for this workflow state.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the branch pattern.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"pattern": schema.StringAttribute{
						MarkdownDescription: "Branch pattern to match.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
					},
					"is_regex": schema.BoolAttribute{
						MarkdownDescription: "Whether the branch pattern is a regex. **Default** `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"draft": schema.StringAttribute{
				MarkdownDescription: "Workflow state used when draft PRs are opened.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"start": schema.StringAttribute{
				MarkdownDescription: "Workflow state used when PRs are opened.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"review": schema.StringAttribute{
				MarkdownDescription: "Workflow state used when reviews are requested on PRs.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"mergeable": schema.StringAttribute{
				MarkdownDescription: "Workflow state used when PRs become mergeable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"merge": schema.StringAttribute{
				MarkdownDescription: "Workflow state used when PRs are merged.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
		},
	}
}

func (r *TeamWorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamWorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamWorkflowResourceModel
	var branchData *TeamWorkflowResourceBranchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := updateTeamWorkflow(ctx, r.client, data, branchData, nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", err))
		return
	}

	tflog.Trace(ctx, "created a team workflow")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamWorkflowResourceModel
	var branchData *TeamWorkflowResourceBranchModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeamWorkflow(ctx, *r.client, data.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team workflow, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	readTeamWorkflow(data, branchData, branchData, response.Team)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamWorkflowResourceModel
	var branchData *TeamWorkflowResourceBranchModel
	var state *TeamWorkflowResourceModel
	var branchState *TeamWorkflowResourceBranchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.Branch.As(ctx, &branchState, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := updateTeamWorkflow(ctx, r.client, data, branchData, branchState)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", err))
		return
	}

	tflog.Trace(ctx, "updated a team workflow")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamWorkflowResourceModel
	var state *TeamWorkflowResourceModel
	var branchState *TeamWorkflowResourceBranchModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.Branch.As(ctx, &branchState, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Key = state.Key

	err := updateTeamWorkflow(ctx, r.client, &data, nil, branchState)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team workflow, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team workflow")
}

func (r *TeamWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	// If only one part is provided, treat the resource for no target branch.
	if len(parts) == 1 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), req.ID)...)
		return
	}

	// If three parts are provided, treat the resource for a target branch.
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: team_key:branch_pattern:is_regex. Got: %q", req.ID),
		)

		return
	}

	// Set the branch object with the parsed values
	branchObj := types.ObjectValueMust(
		branchAttrTypes,
		map[string]attr.Value{
			"id":       types.StringNull(),
			"pattern":  types.StringValue(parts[1]),
			"is_regex": types.BoolValue(parts[2] == "true"),
		},
	)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branch"), branchObj)...)
}

func updateTeamWorkflow(ctx context.Context, client *graphql.Client, data *TeamWorkflowResourceModel, branchPlan *TeamWorkflowResourceBranchModel, branchState *TeamWorkflowResourceBranchModel) error {
	teamKey := data.Key.ValueString()

	existing, err := getTeamWorkflow(ctx, *client, teamKey)

	if err != nil {
		return fmt.Errorf("unable to get team workflow: %w", err)
	}

	teamId := existing.Team.Id

	planBranch := findTeamWorkflowTargetBranch(branchPlan, existing.Team)

	// If the branch is specified in the state but not in the plan, delete it
	if branchState != nil && branchPlan == nil {
		_, err = deleteGitAutomationTargetBranch(ctx, *client, branchState.Id.ValueString())

		if err != nil {
			return fmt.Errorf("unable to delete team workflow branch: %w", err)
		}
	}

	// If the branch is not found, we need to create it.
	if branchPlan != nil && planBranch == nil {
		response, err := createGitAutomationTargetBranch(ctx, *client, GitAutomationTargetBranchCreateInput{
			TeamId:        teamId,
			BranchPattern: branchPlan.Pattern.ValueString(),
			IsRegex:       branchPlan.IsRegex.ValueBool(),
		})

		if err != nil {
			return fmt.Errorf("unable to create team workflow branch: %w", err)
		}

		planBranch = &TargetBranch{
			Id:            response.GitAutomationTargetBranchCreate.TargetBranch.Id,
			BranchPattern: response.GitAutomationTargetBranchCreate.TargetBranch.BranchPattern,
			IsRegex:       response.GitAutomationTargetBranchCreate.TargetBranch.IsRegex,
		}
	}

	var branchId *string

	if planBranch != nil {
		branchId = &planBranch.Id
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesDraft, data.Draft.ValueStringPointer(), branchId)

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesStart, data.Start.ValueStringPointer(), branchId)

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesReview, data.Review.ValueStringPointer(), branchId)

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesMergeable, data.Mergeable.ValueStringPointer(), branchId)

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesMerge, data.Merge.ValueStringPointer(), branchId)

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	existing, err = getTeamWorkflow(ctx, *client, teamKey)

	if err != nil {
		return fmt.Errorf("unable to get team workflow: %w", err)
	}

	readTeamWorkflow(data, branchPlan, branchState, existing.Team)

	return nil
}

func updateEvent(ctx context.Context, client *graphql.Client, existing getTeamWorkflowTeam, teamId string, event GitAutomationStates, stateId *string, branchId *string) error {
	var foundState *TeamWorkflowGitAutomationStatesGitAutomationStateConnectionNodesGitAutomationState

	for _, n := range existing.GitAutomationStates.Nodes {
		if n.Event != event {
			continue
		}

		if (branchId == nil && n.TargetBranch != nil) ||
			(branchId != nil && (n.TargetBranch == nil || n.TargetBranch.Id != *branchId)) {
			continue
		}

		foundState = &n
		break
	}

	var err error

	if foundState != nil {
		if stateId == nil {
			_, err = deleteGitAutomationState(ctx, *client, foundState.Id)
		} else {
			input := GitAutomationStateUpdateInput{
				Event:          event,
				StateId:        stateId,
				TargetBranchId: branchId,
			}

			_, err = updateGitAutomationState(ctx, *client, foundState.Id, input)
		}
	} else if stateId != nil {
		input := GitAutomationStateCreateInput{
			Event:          event,
			TeamId:         teamId,
			StateId:        stateId,
			TargetBranchId: branchId,
		}

		_, err = createGitAutomationState(ctx, *client, input)
	}

	return err
}

func readTeamWorkflow(data *TeamWorkflowResourceModel, branchPlan *TeamWorkflowResourceBranchModel, branchState *TeamWorkflowResourceBranchModel, existing getTeamWorkflowTeam) {
	data.Id = types.StringValue(existing.Id)
	data.Key = types.StringValue(existing.Key)

	branch := findTeamWorkflowTargetBranch(branchPlan, existing)

	if branch == nil {
		if branchPlan == nil {
			data.Branch = types.ObjectNull(branchAttrTypes)
		} else {
			data.Branch = types.ObjectValueMust(
				branchAttrTypes,
				map[string]attr.Value{
					"id":       branchState.Id,
					"pattern":  branchState.Pattern,
					"is_regex": branchState.IsRegex,
				},
			)
		}
	} else {
		data.Branch = types.ObjectValueMust(
			branchAttrTypes,
			map[string]attr.Value{
				"id":       types.StringValue(branch.Id),
				"pattern":  types.StringValue(branch.BranchPattern),
				"is_regex": types.BoolValue(branch.IsRegex),
			},
		)
	}

	for _, n := range existing.GitAutomationStates.Nodes {
		if n.State == nil {
			continue
		}

		if (branch == nil && n.TargetBranch != nil) ||
			(branch != nil && (n.TargetBranch == nil || n.TargetBranch.Id != branch.Id)) {
			continue
		}

		switch n.Event {
		case GitAutomationStatesDraft:
			data.Draft = types.StringValue(n.State.Id)
		case GitAutomationStatesStart:
			data.Start = types.StringValue(n.State.Id)
		case GitAutomationStatesReview:
			data.Review = types.StringValue(n.State.Id)
		case GitAutomationStatesMergeable:
			data.Mergeable = types.StringValue(n.State.Id)
		case GitAutomationStatesMerge:
			data.Merge = types.StringValue(n.State.Id)
		}
	}
}

func findTeamWorkflowTargetBranch(
	branchData *TeamWorkflowResourceBranchModel,
	existing getTeamWorkflowTeam,
) *TargetBranch {
	// Find the branch that matches the pattern and is_regex if specified.
	if branchData != nil {
		for _, n := range existing.GitAutomationStates.Nodes {
			if n.TargetBranch.BranchPattern == branchData.Pattern.ValueString() && n.TargetBranch.IsRegex == branchData.IsRegex.ValueBool() {
				return &TargetBranch{
					Id:            n.TargetBranch.Id,
					BranchPattern: n.TargetBranch.BranchPattern,
					IsRegex:       n.TargetBranch.IsRegex,
				}
			}
		}
	}

	return nil
}
