package provider

import (
	"context"
	"fmt"
	// "math/big"
	// "strings"
	// "encoding/json"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TeamTemplateResource{}
var _ resource.ResourceWithImportState = &TeamTemplateResource{}

func NewTeamTemplateResource() resource.Resource {
	return &TeamTemplateResource{}
}

type TeamTemplateResource struct {
	client *graphql.Client
}

type TeamTemplateResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	TemplateData types.String `tfsdk:"template_data"`
	TeamId      types.String `tfsdk:"team_id"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
}

func (r *TeamTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_template"
}

func (r *TeamTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear team template.",
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
			"template_data": schema.StringAttribute{
				MarkdownDescription: "Template data of the template.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the team.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex(), "must be an uuid"),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the template.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the template.",
				Optional:            true,
			},
		},
	}
}

func (r *TeamTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TemplateCreateInput{
		Name:         data.Name.ValueStringPointer(),
		TemplateData: data.TemplateData.ValueStringPointer(),
		TeamId:       data.TeamId.ValueStringPointer(),
		Type:         data.Type.ValueStringPointer(),
		Description:  data.Description.ValueStringPointer(),
	}

	response, err := templateCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team template, got error: %s", err))
		return
	}

	tflog.Info(ctx, "created a team template")

	teamTemplate := response.TemplateCreate.Template

	data.Id = types.StringValue(teamTemplate.Id)
	data.Name = types.StringPointerValue(teamTemplate.Name)
	data.TemplateData = types.StringPointerValue(teamTemplate.TemplateData)
	data.TeamId = types.StringValue(teamTemplate.Team.Id)
	data.Type = types.StringPointerValue(teamTemplate.Type)
	data.Description = types.StringValue(teamTemplate.Description)

	if teamTemplate.Description == "" {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(teamTemplate.Description)
	}
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data *TeamTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getTemplate(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team template, got error: %s", err))
		return
	}

	teamTemplate := response.Template

	data.Id = types.StringValue(teamTemplate.Id)
	data.Name = types.StringPointerValue(teamTemplate.Name)
	data.TemplateData = types.StringPointerValue(teamTemplate.TemplateData)
	data.TeamId = types.StringValue(teamTemplate.Team.Id)
	data.Type = types.StringPointerValue(teamTemplate.Type)
	data.Description = types.StringValue(teamTemplate.Description)

	if teamTemplate.Description == "" {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(teamTemplate.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := TemplateUpdateInput{
		Name:        data.Name.ValueString(),
		TemplateData: data.TemplateData.ValueString(),
		TeamId:      data.TeamId.ValueString(),
		Description: data.Description.ValueString(),
	}

	response, err := templateUpdate(ctx, *r.client, input, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team template, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a team template")

	teamTemplate := response.TemplateUpdate.Template

	data.Id = types.StringValue(teamTemplate.Id)
	data.Name = types.StringPointerValue(teamTemplate.Name)
	data.TemplateData = types.StringPointerValue(teamTemplate.TemplateData)
	data.TeamId = types.StringValue(teamTemplate.Team.Id)
	data.Description = types.StringValue(teamTemplate.Description)

	if teamTemplate.Description == "" {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(teamTemplate.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := templateDelete(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team template, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a team template")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError("Client Error", "ImportState not implemented")
	return
	/*
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: template_name:team_key. Got: %q", req.ID),
		)

		return
	}

	response, err := findTemplate(ctx, *r.client, parts[0], parts[1])

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import team template, got error: %s", err))
		return
	}

	if len(response.Templates.Nodes) != 1 {
		resp.Diagnostics.AddError("Client Error", "Unable to import team template, got error: template not found")
		return
	}

	data.Id = types.StringValue(response.Templates.Nodes[0].Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	*/	
}
