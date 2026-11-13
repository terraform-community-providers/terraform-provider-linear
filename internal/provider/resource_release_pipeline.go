package provider

import (
	"strings"
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ReleasePipelineResource{}
var _ resource.ResourceWithImportState = &ReleasePipelineResource{}

func NewReleasePipelineResource() resource.Resource {
	return &ReleasePipelineResource{}
}

type ReleasePipelineResource struct {
	client *graphql.Client
}

type ReleasePipelineResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	IsProduction        types.Bool   `tfsdk:"is_production"`
	IncludePathPatterns types.List   `tfsdk:"include_path_patterns"`
	Teams               types.Set    `tfsdk:"teams"`
	DefaultStage        types.Object `tfsdk:"default_stage"`
}

func (r *ReleasePipelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_release_pipeline"
}

func (r *ReleasePipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linear release pipeline.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the release pipeline.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the release pipeline.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the pipeline: `continuous` or `scheduled`. **Default** `continuous`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("continuous"),
				Validators: []validator.String{
					stringvalidator.OneOf("continuous", "scheduled"),
				},
			},
			"is_production": schema.BoolAttribute{
				MarkdownDescription: "Whether this pipeline targets a production environment.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"include_path_patterns": schema.ListAttribute{
				MarkdownDescription: "Glob patterns to filter commits by file path.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"teams": schema.SetAttribute{
				MarkdownDescription: "Identifiers of the teams associated with this pipeline.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"default_stage": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for the default `started` stage created with the pipeline. *This can not be deleted.*",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						releaseStageAttrTypes,
						map[string]attr.Value{
							"id":       types.StringUnknown(),
							"name":     types.StringValue("Started"),
							"color":    types.StringValue("#f59e0b"),
							"position": types.Float64Unknown(),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the release stage.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the stage. **Default** `Started`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("Started"),
					},
					"color": schema.StringAttribute{
						MarkdownDescription: "Color of the stage. **Default** `#f59e0b`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("#f59e0b"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(colorRegex(), "must be a hex color"),
						},
					},
					"position": schema.Float64Attribute{
						MarkdownDescription: "Position of the stage.",
						Computed:            true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *ReleasePipelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReleasePipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ReleasePipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ReleasePipelineCreateInput{
		Name:         data.Name.ValueString(),
		IsProduction: data.IsProduction.ValueBool(),
		Type:         ReleasePipelineType(data.Type.ValueString()),
	}

	if !data.IncludePathPatterns.IsNull() {
		var patterns []string
		resp.Diagnostics.Append(data.IncludePathPatterns.ElementsAs(ctx, &patterns, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.IncludePathPatterns = patterns
	}

	if !data.Teams.IsNull() {
		var teamIds []string
		resp.Diagnostics.Append(data.Teams.ElementsAs(ctx, &teamIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.TeamIds = teamIds
	}

	response, err := createReleasePipeline(ctx, *r.client, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create release pipeline, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a release pipeline")

	pipeline := response.ReleasePipelineCreate.ReleasePipeline.ReleasePipelineFields
	r.readModel(ctx, data, &pipeline, &resp.Diagnostics)

	// Create the default stage on the pipeline.
	var stageData ReleasePipelineStageModel
	resp.Diagnostics.Append(data.DefaultStage.As(ctx, &stageData, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	stageInput := ReleaseStageCreateInput{
		Name:       stageData.Name.ValueString(),
		Type:       ReleaseStageType("started"),
		Color:      stageData.Color.ValueString(),
		Position:   0,
		PipelineId: pipeline.Id,
	}

	stageResp, err := createReleaseStage(ctx, *r.client, stageInput)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create default stage, got error: %s", err))
		return
	}

	stage := stageResp.ReleaseStageCreate.ReleaseStage.ReleaseStageFields
	data.DefaultStage = types.ObjectValueMust(releaseStageAttrTypes, map[string]attr.Value{
		"id":       types.StringValue(stage.Id),
		"name":     types.StringValue(stage.Name),
		"color":    types.StringValue(stage.Color),
		"position": types.Float64Value(stage.Position),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleasePipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ReleasePipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getReleasePipeline(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "Entity not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read release pipeline, got error: %s", err))
		return
	}

	pipeline := response.ReleasePipeline.ReleasePipelineFields
	r.readModel(ctx, data, &pipeline, &resp.Diagnostics)

	// Read the default stage (first started stage) from the pipeline.
	r.readDefaultStage(data, &pipeline)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleasePipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ReleasePipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := ReleasePipelineUpdateInput{
		Name:         data.Name.ValueString(),
		IsProduction: data.IsProduction.ValueBool(),
		Type:         ReleasePipelineType(data.Type.ValueString()),
	}

	if !data.IncludePathPatterns.IsNull() {
		var patterns []string
		resp.Diagnostics.Append(data.IncludePathPatterns.ElementsAs(ctx, &patterns, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.IncludePathPatterns = patterns
	}

	if !data.Teams.IsNull() {
		var teamIds []string
		resp.Diagnostics.Append(data.Teams.ElementsAs(ctx, &teamIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.TeamIds = teamIds
	}

	response, err := updateReleasePipeline(ctx, *r.client, input, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update release pipeline, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a release pipeline")

	pipeline := response.ReleasePipelineUpdate.ReleasePipeline.ReleasePipelineFields
	r.readModel(ctx, data, &pipeline, &resp.Diagnostics)

	// Update the default stage if changed.
	var stageData ReleasePipelineStageModel
	resp.Diagnostics.Append(data.DefaultStage.As(ctx, &stageData, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !stageData.Id.IsUnknown() && !stageData.Id.IsNull() {
		stageInput := ReleaseStageUpdateInput{
			Name:  stageData.Name.ValueString(),
			Color: stageData.Color.ValueString(),
		}

		stageResp, err := updateReleaseStage(ctx, *r.client, stageInput, stageData.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update default stage, got error: %s", err))
			return
		}

		stage := stageResp.ReleaseStageUpdate.ReleaseStage.ReleaseStageFields
		data.DefaultStage = types.ObjectValueMust(releaseStageAttrTypes, map[string]attr.Value{
			"id":       types.StringValue(stage.Id),
			"name":     types.StringValue(stage.Name),
			"color":    types.StringValue(stage.Color),
			"position": types.Float64Value(stage.Position),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReleasePipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ReleasePipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteReleasePipeline(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete release pipeline, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a release pipeline")
}

func (r *ReleasePipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ReleasePipelineResource) readModel(ctx context.Context, data *ReleasePipelineResourceModel, pipeline *ReleasePipelineFields, diags *diag.Diagnostics) {
	data.Id = types.StringValue(pipeline.Id)
	data.Name = types.StringValue(pipeline.Name)
	data.Type = types.StringValue(string(pipeline.Type))
	data.IsProduction = types.BoolValue(pipeline.IsProduction)

	if len(pipeline.IncludePathPatterns) > 0 {
		patterns, d := types.ListValueFrom(ctx, types.StringType, pipeline.IncludePathPatterns)
		diags.Append(d...)
		data.IncludePathPatterns = patterns
	} else if !data.IncludePathPatterns.IsNull() {
		data.IncludePathPatterns = types.ListValueMust(types.StringType, []attr.Value{})
	}

	if len(pipeline.Teams.Nodes) > 0 {
		ids := make([]string, len(pipeline.Teams.Nodes))
		for i, t := range pipeline.Teams.Nodes {
			ids[i] = t.Id
		}
		teams, d := types.SetValueFrom(ctx, types.StringType, ids)
		diags.Append(d...)
		data.Teams = teams
	} else if !data.Teams.IsNull() {
		data.Teams = types.SetValueMust(types.StringType, []attr.Value{})
	}
}

// ReleasePipelineStageModel represents the nested default stage attribute.
type ReleasePipelineStageModel struct {
	Id       types.String  `tfsdk:"id"`
	Name     types.String  `tfsdk:"name"`
	Color    types.String  `tfsdk:"color"`
	Position types.Float64 `tfsdk:"position"`
}

var releaseStageAttrTypes = map[string]attr.Type{
	"id":       types.StringType,
	"name":     types.StringType,
	"color":    types.StringType,
	"position": types.Float64Type,
}

func (r *ReleasePipelineResource) readDefaultStage(data *ReleasePipelineResourceModel, pipeline *ReleasePipelineFields) {
	// Find the default stage by matching the stored ID, or fall back to the
	// first started stage if no ID is stored yet (import case).
	var stageData ReleasePipelineStageModel
	if !data.DefaultStage.IsNull() && !data.DefaultStage.IsUnknown() {
		_ = data.DefaultStage.As(context.Background(), &stageData, basetypes.ObjectAsOptions{})
	}

	for _, s := range pipeline.Stages.Nodes {
		if (stageData.Id.ValueString() != "" && s.Id == stageData.Id.ValueString()) ||
			(stageData.Id.ValueString() == "" && string(s.Type) == "started") {
			data.DefaultStage = types.ObjectValueMust(releaseStageAttrTypes, map[string]attr.Value{
				"id":       types.StringValue(s.Id),
				"name":     types.StringValue(s.Name),
				"color":    types.StringValue(s.Color),
				"position": types.Float64Value(s.Position),
			})
			return
		}
	}
}
