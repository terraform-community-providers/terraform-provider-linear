package provider

import (
	"context"
	"fmt"

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

func (r *WorkspaceLabelResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Linear workspace label.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Identifier of the label.",
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the label.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"description": {
				MarkdownDescription: "Description of the label.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.NullableString(),
				},
			},
			"color": {
				MarkdownDescription: "Color of the label.",
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
			"parent_id": {
				MarkdownDescription: "Parent (label group) of the label.",
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
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		value := data.Description.ValueString()
		input.Description = &value
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	if !data.ParentId.IsNull() {
		value := data.ParentId.ValueString()
		input.ParentId = &value
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

	if issueLabel.Description != nil {
		data.Description = types.StringValue(*issueLabel.Description)
	}

	if issueLabel.Color != nil {
		data.Color = types.StringValue(*issueLabel.Color)
	}

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

	if issueLabel.Description != nil {
		data.Description = types.StringValue(*issueLabel.Description)
	}

	if issueLabel.Color != nil {
		data.Color = types.StringValue(*issueLabel.Color)
	}

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
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		value := data.Description.ValueString()
		input.Description = &value
	}

	if !data.Color.IsUnknown() {
		value := data.Color.ValueString()
		input.Color = &value
	}

	if !data.ParentId.IsNull() {
		value := data.ParentId.ValueString()
		input.ParentId = &value
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

	if issueLabel.Description != nil {
		data.Description = types.StringValue(*issueLabel.Description)
	}

	if issueLabel.Color != nil {
		data.Color = types.StringValue(*issueLabel.Color)
	}

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
