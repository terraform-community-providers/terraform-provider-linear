package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/validators"
)

var _ resource.Resource = &WorkflowStateResource{}
var _ resource.ResourceWithImportState = &WorkflowStateResource{}

func NewWorkflowStateResource() resource.Resource {
	return &WorkflowStateResource{}
}

type WorkflowStateResource struct {
	client *graphql.Client
}

type WorkflowStateResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Position    types.Number `tfsdk:"position"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (r *WorkflowStateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_state"
}

func (r *WorkflowStateResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team workflow state.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the workflow state.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
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
					resource.RequiresReplace(),
				},
			},
			"position": {
				MarkdownDescription: "Position of the workflow state.",
				Type:                types.NumberType,
				Required:            true,
			},
			"color": {
				MarkdownDescription: "Color of the workflow state.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.Match(colorRegex()),
				},
			},
			"description": {
				MarkdownDescription: "Description of the workflow state.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
			},
			"team_id": {
				MarkdownDescription: "Identifier of the team.",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
		},
	}, nil
}

func (r *WorkflowStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WorkflowStateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	position, _ := data.Position.Value.Float64()

	input := WorkflowStateCreateInput{
		Name:     data.Name.Value,
		Type:     data.Type.Value,
		Position: position,
		Color:    data.Color.Value,
		TeamId:   data.TeamId.Value,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	response, err := createWorkflowState(context.Background(), *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a workflow state")

	workflowState := response.WorkflowStateCreate.WorkflowState

	data.Id = types.String{Value: workflowState.Id}
	data.Name = types.String{Value: workflowState.Name}
	data.Type = types.String{Value: workflowState.Type}
	data.Position = types.Number{Value: big.NewFloat(workflowState.Position)}
	data.Color = types.String{Value: workflowState.Color}

	if workflowState.Description != nil {
		data.Description = types.String{Value: *workflowState.Description}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WorkflowStateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getWorkflowState(context.Background(), *r.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a workflow state")

	workflowState := response.WorkflowState

	data.Name = types.String{Value: workflowState.Name}
	data.Type = types.String{Value: workflowState.Type}
	data.Position = types.Number{Value: big.NewFloat(workflowState.Position)}
	data.Color = types.String{Value: workflowState.Color}
	data.TeamId = types.String{Value: workflowState.Team.Id}

	if workflowState.Description != nil {
		data.Description = types.String{Value: *workflowState.Description}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkflowStateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	position, _ := data.Position.Value.Float64()

	input := WorkflowStateUpdateInput{
		Name:     data.Name.Value,
		Color:    data.Color.Value,
		Position: position,
	}

	if !data.Description.IsNull() {
		input.Description = &data.Description.Value
	}

	response, err := updateWorkflowState(context.Background(), *r.client, input, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a workflow state")

	workflowState := response.WorkflowStateUpdate.WorkflowState

	data.Name = types.String{Value: workflowState.Name}
	data.Position = types.Number{Value: big.NewFloat(workflowState.Position)}
	data.Color = types.String{Value: workflowState.Color}

	if workflowState.Description != nil {
		data.Description = types.String{Value: *workflowState.Description}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WorkflowStateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteWorkflowState(context.Background(), *r.client, data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workflow state, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a workflow state")
}

func (r *WorkflowStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: workflow_state_name:team_key. Got: %q", req.ID),
		)

		return
	}

	response, err := findWorkflowState(context.Background(), *r.client, parts[0], parts[1])

	if err != nil || len(response.WorkflowStates.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workflow state, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.WorkflowStates.Nodes[0].Id)...)
}
