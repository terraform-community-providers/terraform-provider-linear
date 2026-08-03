package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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

var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithConfigValidators = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

var documentParentAttributeNames = []string{
	"issue_id",
	"project_id",
}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

type DocumentResource struct {
	client *graphql.Client
}

type DocumentResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Title     types.String `tfsdk:"title"`
	Content   types.String `tfsdk:"content"`
	IssueId   types.String `tfsdk:"issue_id"`
	ProjectId types.String `tfsdk:"project_id"`
}

type documentIssueValues struct {
	ID         string
	Identifier string
}

type documentValues struct {
	ID        string
	Title     string
	Content   *string
	Trashed   bool
	Issue     *documentIssueValues
	ProjectID *string
}

func (r *DocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document"
}

func (r *DocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	uuidParentValidators := []validator.String{
		stringvalidator.RegexMatches(uuidRegex(), "must be a UUID"),
	}
	issueParentValidators := []validator.String{
		stringvalidator.Any(
			stringvalidator.RegexMatches(uuidRegex(), "must be a UUID"),
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[A-Z0-9]{1,5}-[1-9][0-9]*$`),
				"must be a UUID or an uppercase issue key such as LIN-123",
			),
		),
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "A Linear document whose complete Markdown content is managed by Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the document.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the document. An empty string creates an untitled document.",
				Required:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Complete document content in Markdown. Content must be non-empty.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"issue_id": schema.StringAttribute{
				MarkdownDescription: "Identifier or issue key, such as `LIN-123`, of the issue that owns the document. Exactly one parent identifier must be configured.",
				Optional:            true,
				Validators:          issueParentValidators,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the project that owns the document. Exactly one parent identifier must be configured.",
				Optional:            true,
				Validators:          uuidParentValidators,
			},
		},
	}
}

func (r *DocumentResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("issue_id"),
			path.MatchRoot("project_id"),
		),
	}
}

func (r *DocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := createDocument(ctx, *r.client, documentCreateInputFromModel(&data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create document, got error: %s", err))
		return
	}
	if err := documentMutationError("create", response.DocumentCreate.Success); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	data.Id = types.StringValue(response.DocumentCreate.Document.Id)
	tflog.Trace(ctx, "created a document")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getDocument(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		if isDocumentNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read document, got error: %s", err))
		return
	}

	values := documentValuesFromGraphQL(response.Document.Document)
	if values.Trashed {
		resp.State.RemoveResource(ctx)
		return
	}

	applyDocumentValuesToModel(&data, values)
	tflog.Trace(ctx, "read a document")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := updateDocument(ctx, *r.client, data.Id.ValueString(), documentUpdateInputFromModel(&data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update document, got error: %s", err))
		return
	}
	if err := documentMutationError("update", response.DocumentUpdate.Success); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	applyDocumentValuesToModel(&data, documentValuesFromGraphQL(response.DocumentUpdate.Document.Document))
	tflog.Trace(ctx, "updated a document")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := deleteDocument(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete document, got error: %s", err))
		return
	}
	if err := documentMutationError("delete", response.DocumentDelete.Success); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted a document")
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func documentCreateInputFromModel(data *DocumentResourceModel) DocumentCreateInput {
	return DocumentCreateInput{
		Title:     data.Title.ValueString(),
		Content:   data.Content.ValueString(),
		IssueId:   data.IssueId.ValueStringPointer(),
		ProjectId: data.ProjectId.ValueStringPointer(),
	}
}

func documentUpdateInputFromModel(data *DocumentResourceModel) DocumentUpdateInput {
	return DocumentUpdateInput{
		Title:     data.Title.ValueString(),
		Content:   data.Content.ValueString(),
		IssueId:   data.IssueId.ValueStringPointer(),
		ProjectId: data.ProjectId.ValueStringPointer(),
	}
}

// documentValuesFromGraphQL isolates generated GraphQL types from Terraform state mapping.
func documentValuesFromGraphQL(document Document) documentValues {
	values := documentValues{
		ID:      document.Id,
		Title:   document.Title,
		Content: document.Content,
		Trashed: document.Trashed,
	}

	if document.Issue != nil {
		values.Issue = &documentIssueValues{
			ID:         document.Issue.Id,
			Identifier: document.Issue.Identifier,
		}
	}
	if document.Project != nil {
		values.ProjectID = &document.Project.Id
	}

	return values
}

// applyDocumentValuesToModel replaces document state, including clearing stale parent identifiers.
func applyDocumentValuesToModel(data *DocumentResourceModel, values documentValues) {
	currentIssueID := data.IssueId.ValueString()

	data.Id = types.StringValue(values.ID)
	data.Title = types.StringValue(values.Title)
	data.Content = types.StringValue(documentContentForState(data.Content, values.Content))

	data.IssueId = types.StringNull()
	data.ProjectId = types.StringNull()

	switch {
	case values.Issue != nil:
		if currentIssueID == values.Issue.ID || currentIssueID == values.Issue.Identifier {
			data.IssueId = types.StringValue(currentIssueID)
		} else {
			data.IssueId = types.StringValue(values.Issue.ID)
		}
	case values.ProjectID != nil:
		data.ProjectId = types.StringValue(*values.ProjectID)
	}
}

// documentContentForState preserves known Terraform content when Linear differs only by terminal CR/LF characters.
func documentContentForState(current types.String, remote *string) string {
	remoteContent := ""
	if remote != nil {
		remoteContent = *remote
	}

	if current.IsNull() || current.IsUnknown() {
		return remoteContent
	}

	currentContent := current.ValueString()
	if strings.TrimRight(currentContent, "\r\n") == strings.TrimRight(remoteContent, "\r\n") {
		return currentContent
	}

	return remoteContent
}

func documentMutationError(operation string, success bool) error {
	if success {
		return nil
	}
	return fmt.Errorf("Linear reported that it was unable to %s document", operation)
}

func isDocumentNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Entity not found: Document")
}
