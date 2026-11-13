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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ProjectLabelResource{}
var _ resource.ResourceWithImportState = &ProjectLabelResource{}

func NewProjectLabelResource() resource.Resource {
	return &ProjectLabelResource{}
}

type ProjectLabelResource struct {
	client *graphql.Client
}

type ProjectLabelResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	IsGroup     types.Bool   `tfsdk:"is_group"`
	ParentId    types.String `tfsdk:"parent_id"`
}

func (r *ProjectLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_label"
}

func (r *ProjectLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear project label.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the project label.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project label.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the project label.",
				Optional:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the project label.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
				},
			},
			"is_group": schema.BoolAttribute{
				MarkdownDescription: "Whether the project label is a group (parent) label. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the parent project label group.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
		},
	}
}

func (r *ProjectLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ProjectLabelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ProjectLabelCreateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		IsGroup:     data.IsGroup.ValueBool(),
		ParentId:    data.ParentId.ValueStringPointer(),
	}
	if !data.Color.IsUnknown() {
		input.Color = data.Color.ValueStringPointer()
	}

	response, err := projectLabelCreate(ctx, *r.client, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project label, got error: %s", err))
		return
	}
	if !response.ProjectLabelCreate.Success {
		resp.Diagnostics.AddError("Client Error", "Unable to create project label: Linear reported an unsuccessful creation")
		return
	}

	tflog.Trace(ctx, "created a project label")
	setProjectLabelState(data, response.ProjectLabelCreate.ProjectLabel.ProjectLabel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ProjectLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getProjectLabel(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "Entity not found: ProjectLabel") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project label, got error: %s", err))
		return
	}

	setProjectLabelState(data, response.ProjectLabel.ProjectLabel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ProjectLabelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ProjectLabelUpdateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		IsGroup:     data.IsGroup.ValueBool(),
		ParentId:    data.ParentId.ValueStringPointer(),
	}
	if !data.Color.IsUnknown() {
		input.Color = data.Color.ValueStringPointer()
	}

	response, err := projectLabelUpdate(ctx, *r.client, input, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project label, got error: %s", err))
		return
	}
	if !response.ProjectLabelUpdate.Success {
		resp.Diagnostics.AddError("Client Error", "Unable to update project label: Linear reported an unsuccessful update")
		return
	}

	tflog.Trace(ctx, "updated a project label")
	setProjectLabelState(data, response.ProjectLabelUpdate.ProjectLabel.ProjectLabel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ProjectLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := deleteProjectLabel(ctx, *r.client, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted a project label")
}

func deleteProjectLabel(ctx context.Context, client graphql.Client, id string) error {
	response, err := projectLabelDelete(ctx, client, id)
	if err != nil {
		if strings.Contains(err.Error(), "Entity not found: ProjectLabel") {
			return nil
		}
		return fmt.Errorf("unable to delete project label: %w", err)
	}
	if !response.ProjectLabelDelete.Success {
		return fmt.Errorf("unable to delete project label: Linear reported an unsuccessful deletion")
	}
	return nil
}

func (r *ProjectLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setProjectLabelState(data *ProjectLabelResourceModel, projectLabel ProjectLabel) {
	data.Id = types.StringValue(projectLabel.Id)
	data.Name = types.StringValue(projectLabel.Name)
	data.Description = types.StringPointerValue(projectLabel.Description)
	data.Color = types.StringValue(projectLabel.Color)
	data.IsGroup = types.BoolValue(projectLabel.IsGroup)
	if projectLabel.Parent != nil {
		data.ParentId = types.StringValue(projectLabel.Parent.Id)
	} else {
		data.ParentId = types.StringNull()
	}
}
