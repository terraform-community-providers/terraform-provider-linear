package provider

import (
	"context"
	"fmt"
	"regexp"

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

var _ resource.Resource = &CustomViewResource{}
var _ resource.ResourceWithImportState = &CustomViewResource{}

func NewCustomViewResource() resource.Resource {
	return &CustomViewResource{}
}

type CustomViewResource struct {
	client *graphql.Client
}

type CustomViewResourceModel struct {
	// Resource ID, provided by Linear.
	Id types.String `tfsdk:"id"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Color types.String `tfsdk:"color"`
	Icon  types.String `tfsdk:"icon"`
}

func (r *CustomViewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_view"
}

func (r *CustomViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear custom view.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the view.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the view.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the view.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the view.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
				},
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon of the view.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z]+$"), "must only contain letters"),
				},
			},
		},
	}
}

func (r *CustomViewResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CustomViewResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := CustomViewCreateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),

		Color: data.Color.ValueString(),
		Icon:  data.Icon.ValueString(),

		FilterData: map[string]interface{}{},
	}

	response, err := customViewCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom view, got error: %s", err))
		return
	}

	if !response.CustomViewCreate.Success {
		resp.Diagnostics.AddError("Client Error", "Linear reported custom view was not created")
		return
	}

	tflog.Trace(ctx, "created a custom view")
	resp.Diagnostics.Append(resp.State.Set(ctx, customViewToResourceModel(response.CustomViewCreate.CustomView.CustomView))...)
}

func (r *CustomViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CustomViewResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := customView(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom view, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a custom view")
	resp.Diagnostics.Append(resp.State.Set(ctx, customViewToResourceModel(response.CustomView.CustomView))...)
}

func (r *CustomViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CustomViewResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := CustomViewUpdateInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),

		Color: data.Color.ValueString(),
		Icon:  data.Icon.ValueString(),

		FilterData: map[string]interface{}{},
	}

	response, err := customViewUpdate(ctx, *r.client, input, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom view, got error: %s", err))
		return
	}

	if !response.CustomViewUpdate.Success {
		resp.Diagnostics.AddError("Client Error", "Linear reported custom view was not updated")
		return
	}

	tflog.Trace(ctx, "updated a custom view")
	resp.Diagnostics.Append(resp.State.Set(ctx, customViewToResourceModel(response.CustomViewUpdate.CustomView.CustomView))...)
}

func (r *CustomViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CustomViewResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := customViewDelete(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom view, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a custom view")
}

func (r *CustomViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := customView(ctx, *r.client, req.ID)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import custom view, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "imported a custom view")
	resp.Diagnostics.Append(resp.State.Set(ctx, customViewToResourceModel(response.CustomView.CustomView))...)
}

func customViewToResourceModel(view CustomView) *CustomViewResourceModel {
	data := &CustomViewResourceModel{}

	data.Id = types.StringValue(view.Id)

	data.Name = types.StringValue(view.Name)
	data.Description = types.StringPointerValue(view.Description)

	data.Color = types.StringPointerValue(view.Color)
	data.Icon = types.StringPointerValue(view.Icon)

	return data
}
