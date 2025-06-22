package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type TeamWorkflowResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
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

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := updateTeamWorkflow(ctx, r.client, data)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", err))
		return
	}

	tflog.Trace(ctx, "created a team workflow")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamWorkflowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeamWorkflow(ctx, *r.client, data.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team workflow, got error: %s", err))
		return
	}

	read(data, response.Team)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamWorkflowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := updateTeamWorkflow(ctx, r.client, data)

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

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Key = state.Key

	err := updateTeamWorkflow(ctx, r.client, &data)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team workflow, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team workflow")
}

func (r *TeamWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}

func updateTeamWorkflow(ctx context.Context, client *graphql.Client, data *TeamWorkflowResourceModel) error {
	teamKey := data.Key.ValueString()

	existing, err := getTeamWorkflow(ctx, *client, teamKey)

	if err != nil {
		return fmt.Errorf("unable to get team workflow: %w", err)
	}

	teamId := existing.Team.Id

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesDraft, data.Draft.ValueStringPointer())

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesStart, data.Start.ValueStringPointer())

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesReview, data.Review.ValueStringPointer())

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesMergeable, data.Mergeable.ValueStringPointer())

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	err = updateEvent(ctx, client, existing.Team, teamId, GitAutomationStatesMerge, data.Merge.ValueStringPointer())

	if err != nil {
		return fmt.Errorf("unable to update team workflow: %w", err)
	}

	existing, err = getTeamWorkflow(ctx, *client, teamKey)

	if err != nil {
		return fmt.Errorf("unable to get team workflow: %w", err)
	}

	read(data, existing.Team)

	return nil
}

func updateEvent(ctx context.Context, client *graphql.Client, existing getTeamWorkflowTeam, teamId string, event GitAutomationStates, stateId *string) error {
	var foundState *TeamWorkflowGitAutomationStatesGitAutomationStateConnectionNodesGitAutomationState

	for _, n := range existing.GitAutomationStates.Nodes {
		if n.Event == event && n.TargetBranch == nil {
			foundState = &n
			break
		}
	}

	var err error

	if foundState != nil {
		if stateId == nil {
			_, err = deleteGitAutomationState(ctx, *client, foundState.Id)
		} else {
			input := GitAutomationStateUpdateInput{
				Event:          event,
				StateId:        stateId,
				TargetBranchId: nil,
			}

			_, err = updateGitAutomationState(ctx, *client, foundState.Id, input)
		}
	} else if stateId != nil {
		input := GitAutomationStateCreateInput{
			Event:          event,
			TeamId:         teamId,
			StateId:        stateId,
			TargetBranchId: nil,
		}

		_, err = createGitAutomationState(ctx, *client, input)
	}

	return err
}

func read(data *TeamWorkflowResourceModel, team getTeamWorkflowTeam) {
	data.Id = types.StringValue(team.Id)
	data.Key = types.StringValue(team.Key)

	for _, n := range team.GitAutomationStates.Nodes {
		if n.State == nil {
			continue
		}

		if n.TargetBranch != nil {
			// We only care about the target branch, so we skip the rest of the fields
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
