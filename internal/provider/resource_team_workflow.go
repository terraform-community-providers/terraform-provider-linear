package provider

import (
	"context"
	"fmt"
	"regexp"

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

var _ resource.Resource = &TeamWorkflowResource{}
var _ resource.ResourceWithImportState = &TeamWorkflowResource{}

func NewTeamWorkflowResource() resource.Resource {
	return &TeamWorkflowResource{}
}

type TeamWorkflowResource struct {
	client *graphql.Client
}

type TeamWorkflowResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Draft  types.String `tfsdk:"draft"`
	Start  types.String `tfsdk:"start"`
	Review types.String `tfsdk:"review"`
	Merge  types.String `tfsdk:"merge"`
}

func (r *TeamWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_workflow"
}

func (r *TeamWorkflowResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear team workflow.",
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
			"draft": {
				MarkdownDescription: "Workflow state used when draft PRs are opened.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
			"start": {
				MarkdownDescription: "Workflow state used when PRs are opened.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
			"review": {
				MarkdownDescription: "Workflow state used when reviews are requested on PRs.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
			"merge": {
				MarkdownDescription: "Workflow state used when PRs are merged.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.Match(uuidRegex()),
				},
			},
		},
	}, nil
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

	response, err := update(ctx, data, r.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team workflow, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team workflow")

	read(data, response)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamWorkflowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTeamWorkflow(ctx, *r.client, data.Key.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team workflow, got error: %s", err))
		return
	}

	team := response.Team

	data.Id = types.String{Value: team.Id}
	data.Key = types.String{Value: team.Key}

	if team.DraftWorkflowState != nil {
		data.Draft = types.String{Value: team.DraftWorkflowState.Id}
	}

	if team.StartWorkflowState != nil {
		data.Start = types.String{Value: team.StartWorkflowState.Id}
	}

	if team.ReviewWorkflowState != nil {
		data.Review = types.String{Value: team.ReviewWorkflowState.Id}
	}

	if team.MergeWorkflowState != nil {
		data.Merge = types.String{Value: team.MergeWorkflowState.Id}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamWorkflowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := update(ctx, data, r.client)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team workflow, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team workflow")

	read(data, response)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamWorkflowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := updateTeamWorkflow(ctx, *r.client, data.Key.Value, nil, nil, nil, nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team workflow, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team workflow")
}

func (r *TeamWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}

func update(ctx context.Context, data *TeamWorkflowResourceModel, client *graphql.Client) (*updateTeamWorkflowResponse, error) {
	var draft *string
	var start *string
	var review *string
	var merge *string

	if !data.Draft.IsNull() {
		draft = &data.Draft.Value
	}

	if !data.Start.IsNull() {
		start = &data.Start.Value
	}

	if !data.Review.IsNull() {
		review = &data.Review.Value
	}

	if !data.Merge.IsNull() {
		merge = &data.Merge.Value
	}

	return updateTeamWorkflow(ctx, *client, data.Key.Value, draft, start, review, merge)
}

func read(data *TeamWorkflowResourceModel, response *updateTeamWorkflowResponse) {
	team := response.TeamUpdate.Team

	data.Id = types.String{Value: team.Id}
	data.Key = types.String{Value: team.Key}

	if team.DraftWorkflowState != nil {
		data.Draft = types.String{Value: team.DraftWorkflowState.Id}
	}

	if team.StartWorkflowState != nil {
		data.Start = types.String{Value: team.StartWorkflowState.Id}
	}

	if team.ReviewWorkflowState != nil {
		data.Review = types.String{Value: team.ReviewWorkflowState.Id}
	}

	if team.MergeWorkflowState != nil {
		data.Merge = types.String{Value: team.MergeWorkflowState.Id}
	}
}
