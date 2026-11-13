package provider

import (
	"strings"
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ReleaseStageResource{}
var _ resource.ResourceWithImportState = &ReleaseStageResource{}

func NewReleaseStageResource() resource.Resource {
	return &ReleaseStageResource{}
}

type ReleaseStageResource struct {
	client *graphql.Client
}

type ReleaseStageResourceModel struct {
	Id         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Type       types.String  `tfsdk:"type"`
	Color      types.String  `tfsdk:"color"`
	Position   types.Float64 `tfsdk:"position"`
	Frozen     types.Bool    `tfsdk:"frozen"`
	PipelineId types.String  `tfsdk:"pipeline_id"`
}

func (r *ReleaseStageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_release_stage"
}

func (r *ReleaseStageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A custom stage in a Linear release pipeline. Custom stages are inserted between the system-managed Started and Released stages. Each pipeline has four fixed lifecycle stages (Planned, Started, Released, Canceled) that cannot be created or deleted via the API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the release stage.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the release stage. Can be any descriptive name (e.g. \"QA\", \"Staging\", \"Production\").",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the stage. Always `started` for user-created stages. The system-managed lifecycle stages (Planned, Started, Released, Canceled) are not accessible via the API.",
				Computed:            true,
				Default:             stringdefault.StaticString("started"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the stage as a HEX string (e.g. `#22c55e`). Each stage can have a unique color displayed in the Linear UI.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
				},
			},
			"position": schema.Float64Attribute{
				MarkdownDescription: "Position of the stage within the pipeline as a float value. Linear uses fractional positioning for ordering (e.g. 970.71, 2008.73). Stages are ordered by position between the system Started and Released stages.",
				Required:            true,
			},
			"frozen": schema.BoolAttribute{
				MarkdownDescription: "Whether the stage is frozen. Syncs will not automatically add issues to frozen stages.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"pipeline_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the release pipeline.",
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

func (r *ReleaseStageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ReleaseStageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ReleaseStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ReleaseStageCreateInput{
		Name:       data.Name.ValueString(),
		Type:       ReleaseStageType(data.Type.ValueString()),
		Color:      data.Color.ValueString(),
		Position:   data.Position.ValueFloat64(),
		Frozen:     data.Frozen.ValueBool(),
		PipelineId: data.PipelineId.ValueString(),
	}

	response, err := createReleaseStage(ctx, *r.client, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create release stage, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a release stage")

	stage := response.ReleaseStageCreate.ReleaseStage.ReleaseStageFields
	data.Id = types.StringValue(stage.Id)
	data.Name = types.StringValue(stage.Name)
	data.Type = types.StringValue(string(stage.Type))
	data.Color = types.StringValue(stage.Color)
	data.Position = types.Float64Value(stage.Position)
	data.Frozen = types.BoolValue(stage.Frozen)
	data.PipelineId = types.StringValue(stage.Pipeline.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleaseStageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ReleaseStageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getReleaseStage(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "Entity not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read release stage, got error: %s", err))
		return
	}

	stage := response.ReleaseStage.ReleaseStageFields
	data.Id = types.StringValue(stage.Id)
	data.Name = types.StringValue(stage.Name)
	data.Type = types.StringValue(string(stage.Type))
	data.Color = types.StringValue(stage.Color)
	data.Position = types.Float64Value(stage.Position)
	data.Frozen = types.BoolValue(stage.Frozen)
	data.PipelineId = types.StringValue(stage.Pipeline.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleaseStageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ReleaseStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ReleaseStageUpdateInput{
		Name:     data.Name.ValueString(),
		Color:    data.Color.ValueString(),
		Position: data.Position.ValueFloat64(),
		Frozen:   data.Frozen.ValueBool(),
	}

	response, err := updateReleaseStage(ctx, *r.client, input, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update release stage, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a release stage")

	stage := response.ReleaseStageUpdate.ReleaseStage.ReleaseStageFields
	data.Id = types.StringValue(stage.Id)
	data.Name = types.StringValue(stage.Name)
	data.Type = types.StringValue(string(stage.Type))
	data.Color = types.StringValue(stage.Color)
	data.Position = types.Float64Value(stage.Position)
	data.Frozen = types.BoolValue(stage.Frozen)
	data.PipelineId = types.StringValue(stage.Pipeline.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleaseStageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ReleaseStageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteReleaseStage(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		// Archive may fail if this is the last stage of its type in the pipeline.
		// The stage will be cleaned up when the pipeline is archived.
		tflog.Warn(ctx, fmt.Sprintf("unable to archive release stage (may be last of type): %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a release stage")
}

func (r *ReleaseStageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
