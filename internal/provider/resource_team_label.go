package provider

import (
	"context"
	"fmt"
	"strings"

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

var _ resource.Resource = &TeamLabelResource{}
var _ resource.ResourceWithImportState = &TeamLabelResource{}

func NewTeamLabelResource() resource.Resource {
	return &TeamLabelResource{}
}

type TeamLabelResource struct {
	client *graphql.Client
}

type TeamLabelResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	ParentId    types.String `tfsdk:"parent_id"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (r *TeamLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_label"
}

func (r *TeamLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear team label.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the label.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the label.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the label.",
				Optional:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the label.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
				},
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Parent (label group) of the label.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the team.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
		},
	}
}

func (r *TeamLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		ParentId:    data.ParentId.ValueStringPointer(),
		TeamId:      data.TeamId.ValueStringPointer(),
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	response, err := createLabel(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a team label")

	issueLabel := response.IssueLabelCreate.IssueLabel

	data.Id = types.StringValue(issueLabel.Id)
	data.Name = types.StringValue(issueLabel.Name)
	data.Description = types.StringPointerValue(issueLabel.Description)
	data.Color = types.StringPointerValue(issueLabel.Color)

	if issueLabel.Parent != nil {
		data.ParentId = types.StringValue(issueLabel.Parent.Id)
	}

	if issueLabel.Team != nil {
		data.TeamId = types.StringValue(issueLabel.Team.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getLabel(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team label, got error: %s", err))
		return
	}

	issueLabel := response.IssueLabel

	data.Id = types.StringValue(issueLabel.Id)
	data.Name = types.StringValue(issueLabel.Name)
	data.Description = types.StringPointerValue(issueLabel.Description)
	data.Color = types.StringPointerValue(issueLabel.Color)

	if issueLabel.Parent != nil {
		data.ParentId = types.StringValue(issueLabel.Parent.Id)
	}

	if issueLabel.Team != nil {
		data.TeamId = types.StringValue(issueLabel.Team.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelUpdateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		ParentId:    data.ParentId.ValueStringPointer(),
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	response, err := updateLabel(ctx, *r.client, input, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team label")

	issueLabel := response.IssueLabelUpdate.IssueLabel

	data.Id = types.StringValue(issueLabel.Id)
	data.Name = types.StringValue(issueLabel.Name)
	data.Description = types.StringPointerValue(issueLabel.Description)
	data.Color = types.StringPointerValue(issueLabel.Color)

	if issueLabel.Parent != nil {
		data.ParentId = types.StringValue(issueLabel.Parent.Id)
	}

	if issueLabel.Team != nil {
		data.TeamId = types.StringValue(issueLabel.Team.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteLabel(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team label")
}

func (r *TeamLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: label_name:team_key. Got: %q", req.ID),
		)

		return
	}

	response, err := findTeamLabel(ctx, *r.client, parts[0], parts[1])

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import team label, got error: %s", err))
		return
	}

	if len(response.IssueLabels.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", "Unable to import team label, got error: label not found")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.IssueLabels.Nodes[0].Id)...)
}
