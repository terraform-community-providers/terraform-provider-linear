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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TemplateResource{}
var _ resource.ResourceWithImportState = &TemplateResource{}

func NewTemplateResource() resource.Resource {
	return &TemplateResource{}
}

type TemplateResource struct {
	client *graphql.Client
}

type TemplateResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	TeamId      types.String `tfsdk:"team_id"`
	Data        types.String `tfsdk:"data"`
}

func (r *TemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear template resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the template.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the template.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the template.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the template. **Default** `issue`.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("issue"),
				Validators: []validator.String{
					stringvalidator.OneOf("issue", "project", "document"),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the team. If not provided, creates a workspace-level template.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "Template data of the template.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *TemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TemplateCreateInput{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueStringPointer(),
		Type:         data.Type.ValueString(),
		TeamId:       data.TeamId.ValueStringPointer(),
		TemplateData: data.Data.ValueString(),
	}

	response, err := templateCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create template, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a template")

	template := response.TemplateCreate.Template

	data.Id = types.StringValue(template.Id)
	data.Name = types.StringValue(template.Name)
	data.Description = types.StringPointerValue(template.Description)
	data.Type = types.StringValue(template.Type)
	data.Data = types.StringValue(template.TemplateData)

	if template.Team != nil {
		data.TeamId = types.StringValue(template.Team.Id)
	} else {
		data.TeamId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data *TemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTemplate(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read template, got error: %s", err))
		return
	}

	template := response.Template

	data.Id = types.StringValue(template.Id)
	data.Name = types.StringValue(template.Name)
	data.Description = types.StringPointerValue(template.Description)
	data.Type = types.StringValue(template.Type)
	data.Data = types.StringValue(template.TemplateData)

	if template.Team != nil {
		data.TeamId = types.StringValue(template.Team.Id)
	} else {
		data.TeamId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TemplateUpdateInput{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueStringPointer(),
		TeamId:       data.TeamId.ValueStringPointer(),
		TemplateData: data.Data.ValueString(),
	}

	response, err := templateUpdate(ctx, *r.client, input, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update template, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a template")

	template := response.TemplateUpdate.Template

	data.Id = types.StringValue(template.Id)
	data.Name = types.StringValue(template.Name)
	data.Description = types.StringPointerValue(template.Description)
	data.Type = types.StringValue(template.Type)
	data.Data = types.StringValue(template.TemplateData)

	if template.Team != nil {
		data.TeamId = types.StringValue(template.Team.Id)
	} else {
		data.TeamId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := templateDelete(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete template, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a template")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
