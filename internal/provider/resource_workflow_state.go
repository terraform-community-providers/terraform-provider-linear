package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/frankgreco/terraform-helpers/validators"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = workflowStateResourceType{}
var _ tfsdk.Resource = workflowStateResource{}
var _ tfsdk.ResourceWithImportState = workflowStateResource{}

type workflowStateResourceType struct{}

func (t workflowStateResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team workflow state.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the workflow state.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the workflow state.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"type": {
				MarkdownDescription: "Type of the workflow state.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSlice(true, "backlog", "unstarted", "started", "completed", "canceled"),
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"description": {
				MarkdownDescription: "Description of the workflow state.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"color": {
				MarkdownDescription: "Color of the workflow state.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.Match(colorRegex()),
				},
			},
			"position": {
				MarkdownDescription: "Position of the workflow state.",
				Type:                types.NumberType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"team_id": {
				MarkdownDescription: "Identifier of the team.",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
			"default": {
				MarkdownDescription: "Whether the workflow state is used for issues that are opened.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"draft": {
				MarkdownDescription: "Whether the workflow state is used for PRs that are opened as drafts.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"start": {
				MarkdownDescription: "Whether the workflow state is used for PRs that are opened.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"review": {
				MarkdownDescription: "Whether the workflow state is used for PRs that have reviews requested.",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"merge": {
				MarkdownDescription: "Whether the workflow state is used for PRs that are merged.",
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

func (t workflowStateResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return workflowStateResource{
		provider: provider,
	}, diags
}

type workflowStateResourceData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Position    types.Number `tfsdk:"position"`
	TeamId      types.String `tfsdk:"team_id"`
	Default     types.Bool   `tfsdk:"default"`
	Draft       types.Bool   `tfsdk:"draft"`
	Start       types.Bool   `tfsdk:"start"`
	Review      types.Bool   `tfsdk:"review"`
	Merge       types.Bool   `tfsdk:"merge"`
}

type workflowStateResource struct {
	provider provider
}

func (r workflowStateResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data workflowStateResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := WorkflowStateCreateInput{
		Name:        data.Name.Value,
		Type:        data.Type.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
		TeamId:      data.TeamId.Value,
	}

	if data.Position.Value != nil {
		input.Position, _ = data.Position.Value.Float64()
	}

	response, err := createWorkflowState(context.Background(), r.provider.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a workflow state")

	data.Id = types.String{Value: response.WorkflowStateCreate.WorkflowState.Id}
	data.Description = types.String{Value: response.WorkflowStateCreate.WorkflowState.Description}
	data.Position = types.Number{Value: big.NewFloat(response.WorkflowStateCreate.WorkflowState.Position)}

	teamInput := TeamUpdateInput{}

	if !data.Default.IsNull() && !data.Default.IsUnknown() && data.Default.Value {
		teamInput.DefaultIssueStateId = data.Id.Value
	}

	if !data.Draft.IsNull() && !data.Draft.IsUnknown() && data.Draft.Value {
		teamInput.DraftWorkflowStateId = data.Id.Value
	}

	if !data.Start.IsNull() && !data.Start.IsUnknown() && data.Start.Value {
		teamInput.StartWorkflowStateId = data.Id.Value
	}

	if !data.Review.IsNull() && !data.Review.IsUnknown() && data.Review.Value {
		teamInput.ReviewWorkflowStateId = data.Id.Value
	}

	if !data.Merge.IsNull() && !data.Merge.IsUnknown() && data.Merge.Value {
		teamInput.MergeWorkflowStateId = data.Id.Value
	}

	teamResponse, teamErr := updateTeamWorkflowAutomation(context.Background(), r.provider.client, teamInput, data.TeamId.Value)

	if teamErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", teamErr))
		return
	}

	tflog.Trace(ctx, "updated a team")

	if teamResponse.TeamUpdate.Team.DefaultIssueState.Id == data.Id.Value {
		data.Default = types.Bool{Value: true}
	} else {
		data.Default = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.DraftWorkflowState.Id == data.Id.Value {
		data.Draft = types.Bool{Value: true}
	} else {
		data.Draft = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.StartWorkflowState.Id == data.Id.Value {
		data.Start = types.Bool{Value: true}
	} else {
		data.Start = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.ReviewWorkflowState.Id == data.Id.Value {
		data.Review = types.Bool{Value: true}
	} else {
		data.Review = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.MergeWorkflowState.Id == data.Id.Value {
		data.Merge = types.Bool{Value: true}
	} else {
		data.Merge = types.Bool{Value: false}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workflowStateResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data workflowStateResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkflowState(context.Background(), r.provider.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow state, got error: %s", err))
		return
	}

	data.Name = types.String{Value: response.WorkflowState.Name}
	data.Type = types.String{Value: response.WorkflowState.Type}
	data.Description = types.String{Value: response.WorkflowState.Description}
	data.Color = types.String{Value: response.WorkflowState.Color}
	data.Position = types.Number{Value: big.NewFloat(response.WorkflowState.Position)}
	data.TeamId = types.String{Value: response.WorkflowState.Team.Id}

	teamResponse, teamErr := getTeamWorkflowAutomation(context.Background(), r.provider.client, data.TeamId.Value)

	if teamErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", teamErr))
		return
	}

	if teamResponse.Team.DefaultIssueState.Id == data.Id.Value {
		data.Default = types.Bool{Value: true}
	} else {
		data.Default = types.Bool{Value: false}
	}

	if teamResponse.Team.DraftWorkflowState.Id == data.Id.Value {
		data.Draft = types.Bool{Value: true}
	} else {
		data.Draft = types.Bool{Value: false}
	}

	if teamResponse.Team.StartWorkflowState.Id == data.Id.Value {
		data.Start = types.Bool{Value: true}
	} else {
		data.Start = types.Bool{Value: false}
	}

	if teamResponse.Team.ReviewWorkflowState.Id == data.Id.Value {
		data.Review = types.Bool{Value: true}
	} else {
		data.Review = types.Bool{Value: false}
	}

	if teamResponse.Team.MergeWorkflowState.Id == data.Id.Value {
		data.Merge = types.Bool{Value: true}
	} else {
		data.Merge = types.Bool{Value: false}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workflowStateResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data workflowStateResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := WorkflowStateUpdateInput{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
	}

	if data.Position.Value != nil {
		input.Position, _ = data.Position.Value.Float64()
	}

	response, err := updateWorkflowState(context.Background(), r.provider.client, input, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a workflow state")

	data.Name = types.String{Value: response.WorkflowStateUpdate.WorkflowState.Name}
	data.Description = types.String{Value: response.WorkflowStateUpdate.WorkflowState.Description}
	data.Color = types.String{Value: response.WorkflowStateUpdate.WorkflowState.Color}
	data.Position = types.Number{Value: big.NewFloat(response.WorkflowStateUpdate.WorkflowState.Position)}

	teamInput := TeamUpdateInput{}

	if !data.Default.IsNull() && !data.Default.IsUnknown() {
		teamInput.DefaultIssueStateId = data.Id.Value
	}

	if !data.Draft.IsNull() && !data.Draft.IsUnknown() {
		teamInput.DraftWorkflowStateId = data.Id.Value
	}

	if !data.Start.IsNull() && !data.Start.IsUnknown() {
		teamInput.StartWorkflowStateId = data.Id.Value
	}

	if !data.Review.IsNull() && !data.Review.IsUnknown() {
		teamInput.ReviewWorkflowStateId = data.Id.Value
	}

	if !data.Merge.IsNull() && !data.Merge.IsUnknown() {
		teamInput.MergeWorkflowStateId = data.Id.Value
	}

	teamResponse, teamErr := updateTeamWorkflowAutomation(context.Background(), r.provider.client, teamInput, data.TeamId.Value)

	if teamErr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", teamErr))
		return
	}

	tflog.Trace(ctx, "updated a team")

	if teamResponse.TeamUpdate.Team.DefaultIssueState.Id == data.Id.Value {
		data.Default = types.Bool{Value: true}
	} else {
		data.Default = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.DraftWorkflowState.Id == data.Id.Value {
		data.Draft = types.Bool{Value: true}
	} else {
		data.Draft = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.StartWorkflowState.Id == data.Id.Value {
		data.Start = types.Bool{Value: true}
	} else {
		data.Start = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.ReviewWorkflowState.Id == data.Id.Value {
		data.Review = types.Bool{Value: true}
	} else {
		data.Review = types.Bool{Value: false}
	}

	if teamResponse.TeamUpdate.Team.MergeWorkflowState.Id == data.Id.Value {
		data.Merge = types.Bool{Value: true}
	} else {
		data.Merge = types.Bool{Value: false}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r workflowStateResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data workflowStateResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// #7: Delete workflow state
	// _, err := deleteWorkflowState(context.Background(), r.provider.client, data.Id.Value)
	var err error

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a workflow state")
}

func (r workflowStateResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: workflow_state_name:team_key. Got: %q", req.ID),
		)

		return
	}

	response, err := findWorkflowState(context.Background(), r.provider.client, parts[0], parts[1])

	if err != nil || len(response.WorkflowStates.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workflow state, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), response.WorkflowStates.Nodes[0].Id)...)
}
