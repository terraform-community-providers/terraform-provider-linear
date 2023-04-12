package provider

import (
	"context"
	"fmt"

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
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
)

var _ resource.Resource = &WorkspaceLabelResource{}
var _ resource.ResourceWithImportState = &WorkspaceLabelResource{}

func NewWorkspaceLabelResource() resource.Resource {
	return &WorkspaceLabelResource{}
}

type WorkspaceLabelResource struct {
	client *graphql.Client
}

type WorkspaceLabelResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	ParentId    types.String `tfsdk:"parent_id"`
}

func (r *WorkspaceLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_label"
}

func (r *WorkspaceLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear workspace label.",
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.NullableString(),
				},
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.NullableString(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
		},
	}
}

func (r *WorkspaceLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkspaceLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WorkspaceLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := IssueLabelCreateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		ParentId:    data.ParentId.ValueStringPointer(),
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	response, err := createLabel(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a workspace label")

	issueLabel := response.IssueLabelCreate.IssueLabel

	data.Id = types.StringValue(issueLabel.Id)
	data.Name = types.StringValue(issueLabel.Name)
	data.Description = types.StringPointerValue(issueLabel.Description)
	data.Color = types.StringPointerValue(issueLabel.Color)

	if issueLabel.Parent != nil {
		data.ParentId = types.StringValue(issueLabel.Parent.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WorkspaceLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getLabel(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace label, got error: %s", err))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WorkspaceLabelResourceModel

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a workspace label")

	issueLabel := response.IssueLabelUpdate.IssueLabel

	data.Id = types.StringValue(issueLabel.Id)
	data.Name = types.StringValue(issueLabel.Name)
	data.Description = types.StringPointerValue(issueLabel.Description)
	data.Color = types.StringPointerValue(issueLabel.Color)

	if issueLabel.Parent != nil {
		data.ParentId = types.StringValue(issueLabel.Parent.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WorkspaceLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteLabel(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a workspace label")
}

func (r *WorkspaceLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := findWorkspaceLabel(ctx, *r.client, req.ID)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import workspace label, got error: %s", err))
		return
	}

	if len(response.IssueLabels.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", "Unable to import workspace label, got error: label not found")
		return
	}

	if response.IssueLabels.Nodes[0].Team.Id != "" {
		resp.Diagnostics.AddError("Client Error", "Unable to import workspace label, got error: label is a team label")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.IssueLabels.Nodes[0].Id)...)
}
