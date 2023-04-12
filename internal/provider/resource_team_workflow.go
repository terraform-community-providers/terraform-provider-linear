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

	response, err := getTeamWorkflow(ctx, *r.client, data.Key.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team workflow, got error: %s", err))
		return
	}

	team := response.Team

	data.Id = types.StringValue(team.Id)
	data.Key = types.StringValue(team.Key)

	if team.DraftWorkflowState != nil {
		data.Draft = types.StringValue(team.DraftWorkflowState.Id)
	}

	if team.StartWorkflowState != nil {
		data.Start = types.StringValue(team.StartWorkflowState.Id)
	}

	if team.ReviewWorkflowState != nil {
		data.Review = types.StringValue(team.ReviewWorkflowState.Id)
	}

	if team.MergeWorkflowState != nil {
		data.Merge = types.StringValue(team.MergeWorkflowState.Id)
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

	_, err := updateTeamWorkflow(ctx, *r.client, data.Key.ValueString(), nil, nil, nil, nil)

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
	draft := data.Draft.ValueStringPointer()
	start := data.Start.ValueStringPointer()
	review := data.Review.ValueStringPointer()
	merge := data.Merge.ValueStringPointer()

	return updateTeamWorkflow(ctx, *client, data.Key.ValueString(), draft, start, review, merge)
}

func read(data *TeamWorkflowResourceModel, response *updateTeamWorkflowResponse) {
	team := response.TeamUpdate.Team

	data.Id = types.StringValue(team.Id)
	data.Key = types.StringValue(team.Key)

	if team.DraftWorkflowState != nil {
		data.Draft = types.StringValue(team.DraftWorkflowState.Id)
	}

	if team.StartWorkflowState != nil {
		data.Start = types.StringValue(team.StartWorkflowState.Id)
	}

	if team.ReviewWorkflowState != nil {
		data.Review = types.StringValue(team.ReviewWorkflowState.Id)
	}

	if team.MergeWorkflowState != nil {
		data.Merge = types.StringValue(team.MergeWorkflowState.Id)
	}
}
