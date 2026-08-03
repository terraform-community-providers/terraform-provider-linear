package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDocumentResourceSchema(t *testing.T) {
	t.Parallel()

	var response resource.SchemaResponse
	(&DocumentResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)

	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", response.Diagnostics)
	}

	required := []string{"title", "content"}
	for _, name := range required {
		attribute, ok := response.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema is missing %q", name)
		}
		if !attribute.IsRequired() {
			t.Errorf("%q should be required", name)
		}
	}

	computed, ok := response.Schema.Attributes["id"]
	if !ok {
		t.Fatal("schema is missing \"id\"")
	}
	if !computed.IsComputed() {
		t.Error("\"id\" should be computed")
	}

	for _, name := range []string{"project_id", "issue_id"} {
		attribute, ok := response.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema is missing %q", name)
		}
		if !attribute.IsOptional() {
			t.Errorf("%q should be optional", name)
		}
	}

	for _, name := range []string{"initiative_id", "team_id", "release_id", "cycle_id"} {
		if _, ok := response.Schema.Attributes[name]; ok {
			t.Errorf("schema should not expose internal Linear parent attribute %q", name)
		}
	}
}

func TestDocumentResourceParentValidation(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parents     map[string]string
		expectError bool
	}{
		"no parent": {
			expectError: true,
		},
		"one parent": {
			parents: map[string]string{
				"project_id": "project-id",
			},
		},
		"issue identifier parent": {
			parents: map[string]string{
				"issue_id": "LIN-123",
			},
		},
		"multiple parents": {
			parents: map[string]string{
				"issue_id":   "LIN-123",
				"project_id": "project-id",
			},
			expectError: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var schemaResponse resource.SchemaResponse
			documentResource := &DocumentResource{}
			documentResource.Schema(context.Background(), resource.SchemaRequest{}, &schemaResponse)

			request := resource.ValidateConfigRequest{
				Config: documentResourceTestConfig(schemaResponse.Schema, testCase.parents),
			}
			var response resource.ValidateConfigResponse

			for _, configValidator := range documentResource.ConfigValidators(context.Background()) {
				configValidator.ValidateResource(context.Background(), request, &response)
			}

			if got := response.Diagnostics.HasError(); got != testCase.expectError {
				t.Fatalf("expected validation error %t, got %t: %v", testCase.expectError, got, response.Diagnostics)
			}
		})
	}
}

func TestDocumentResourceParentIdentifierValidation(t *testing.T) {
	t.Parallel()

	const validUUID = "53c7964a-5bd4-4679-8cca-a5b78498b2b3"

	testCases := map[string]struct {
		attribute   string
		value       string
		expectError bool
	}{
		"project UUID": {
			attribute: "project_id",
			value:     validUUID,
		},
		"invalid project": {
			attribute:   "project_id",
			value:       "project-id",
			expectError: true,
		},
		"issue UUID": {
			attribute: "issue_id",
			value:     validUUID,
		},
		"issue key": {
			attribute: "issue_id",
			value:     "LIN-123",
		},
		"lowercase issue key": {
			attribute:   "issue_id",
			value:       "lin-123",
			expectError: true,
		},
		"issue key without number": {
			attribute:   "issue_id",
			value:       "LIN",
			expectError: true,
		},
		"empty issue identifier": {
			attribute:   "issue_id",
			value:       "",
			expectError: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := documentStringAttributeValidationHasError(t, testCase.attribute, testCase.value)
			if got != testCase.expectError {
				t.Fatalf("expected validation error %t for %s=%q, got %t", testCase.expectError, testCase.attribute, testCase.value, got)
			}
		})
	}
}

func TestDocumentResourceManagedFieldValidation(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute   string
		value       string
		expectError bool
	}{
		"empty title": {
			attribute: "title",
			value:     "",
		},
		"empty content": {
			attribute:   "content",
			value:       "",
			expectError: true,
		},
		"Markdown content": {
			attribute: "content",
			value:     "# Managed by Terraform",
		},
		"newline-terminated content": {
			attribute: "content",
			value:     "Managed by Terraform\n",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := documentStringAttributeValidationHasError(t, testCase.attribute, testCase.value)
			if got != testCase.expectError {
				t.Fatalf("expected validation error %t for %s=%q, got %t", testCase.expectError, testCase.attribute, testCase.value, got)
			}
		})
	}
}

func TestApplyDocumentValuesToModel(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentIssueID string
		values         documentValues
		expectedParent string
		expectedValue  string
	}{
		"issue UUID": {
			values: documentValues{
				Issue: &documentIssueValues{
					ID:         "issue-uuid",
					Identifier: "LIN-123",
				},
			},
			expectedParent: "issue_id",
			expectedValue:  "issue-uuid",
		},
		"issue identifier remains stable": {
			currentIssueID: "LIN-123",
			values: documentValues{
				Issue: &documentIssueValues{
					ID:         "issue-uuid",
					Identifier: "LIN-123",
				},
			},
			expectedParent: "issue_id",
			expectedValue:  "LIN-123",
		},
		"project": {
			values: documentValues{
				ProjectID: stringPointer("project-id"),
			},
			expectedParent: "project_id",
			expectedValue:  "project-id",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			model := &DocumentResourceModel{
				IssueId: types.StringValue(testCase.currentIssueID),
			}
			testCase.values.ID = "document-id"
			testCase.values.Title = ""
			testCase.values.Content = nil

			applyDocumentValuesToModel(model, testCase.values)

			if got := model.Id.ValueString(); got != "document-id" {
				t.Errorf("expected id document-id, got %q", got)
			}
			if got := model.Title.ValueString(); got != "" {
				t.Errorf("expected empty title, got %q", got)
			}
			if model.Content.IsNull() || model.Content.IsUnknown() {
				t.Error("null API content should normalize to a known empty string")
			} else if got := model.Content.ValueString(); got != "" {
				t.Errorf("expected empty content, got %q", got)
			}

			for _, parentName := range documentParentAttributeNames {
				parentValue := documentParentModelValue(model, parentName)
				if parentName == testCase.expectedParent {
					if parentValue.IsNull() || parentValue.IsUnknown() {
						t.Errorf("%s should be known", parentName)
					} else if got := parentValue.ValueString(); got != testCase.expectedValue {
						t.Errorf("expected %s %q, got %q", parentName, testCase.expectedValue, got)
					}
				} else if !parentValue.IsNull() {
					t.Errorf("%s should be null, got %q", parentName, parentValue.ValueString())
				}
			}
		})
	}
}

func TestApplyDocumentValuesToModelUsesContent(t *testing.T) {
	t.Parallel()

	content := "# Managed by Terraform"
	model := &DocumentResourceModel{}

	applyDocumentValuesToModel(model, documentValues{
		ID:        "document-id",
		Title:     "Runbook",
		Content:   &content,
		ProjectID: stringPointer("project-id"),
	})

	if got := model.Content.ValueString(); got != content {
		t.Fatalf("expected content %q, got %q", content, got)
	}
}

func TestApplyDocumentValuesToModelNormalizesOnlyTerminalLineEndings(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		current  types.String
		remote   *string
		expected string
	}{
		"single terminal LF": {
			current:  types.StringValue("body\n"),
			remote:   stringPointer("body"),
			expected: "body\n",
		},
		"multiple terminal LF": {
			current:  types.StringValue("body\n\n\n"),
			remote:   stringPointer("body\n"),
			expected: "body\n\n\n",
		},
		"terminal CRLF": {
			current:  types.StringValue("body\r\n"),
			remote:   stringPointer("body\n"),
			expected: "body\r\n",
		},
		"exact match": {
			current:  types.StringValue("body\n"),
			remote:   stringPointer("body\n"),
			expected: "body\n",
		},
		"null current adopts remote": {
			current:  types.StringNull(),
			remote:   stringPointer("remote\n"),
			expected: "remote\n",
		},
		"unknown current adopts remote": {
			current:  types.StringUnknown(),
			remote:   stringPointer("remote\n"),
			expected: "remote\n",
		},
		"nil API content with empty current": {
			current:  types.StringValue(""),
			remote:   nil,
			expected: "",
		},
		"nil API content with newline-only current": {
			current:  types.StringValue("\n\r\n"),
			remote:   nil,
			expected: "\n\r\n",
		},
		"substantive text change adopts remote": {
			current:  types.StringValue("old body\n"),
			remote:   stringPointer("new body"),
			expected: "new body",
		},
		"internal newline change adopts remote": {
			current:  types.StringValue("line one\nline two\n"),
			remote:   stringPointer("line one line two"),
			expected: "line one line two",
		},
		"trailing spaces remain non-equivalent": {
			current:  types.StringValue("body "),
			remote:   stringPointer("body"),
			expected: "body",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			model := &DocumentResourceModel{
				Content: testCase.current,
			}
			applyDocumentValuesToModel(model, documentValues{
				ID:      "document-id",
				Title:   "Document",
				Content: testCase.remote,
			})

			if model.Content.IsNull() || model.Content.IsUnknown() {
				t.Fatalf("expected known content %q, got %s", testCase.expected, model.Content)
			}
			if got := model.Content.ValueString(); got != testCase.expected {
				t.Fatalf("expected content %q, got %q", testCase.expected, got)
			}
		})
	}
}

func TestDocumentCreateInputIncludesManagedFieldsAndOnlySelectedParent(t *testing.T) {
	t.Parallel()

	input := documentCreateInputFromModel(&DocumentResourceModel{
		Title:     types.StringValue(""),
		Content:   types.StringValue("# Managed by Terraform"),
		ProjectId: types.StringValue("project-id"),
	})

	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("could not marshal create input: %v", err)
	}
	got := string(encoded)

	for _, expected := range []string{`"title":""`, `"content":"# Managed by Terraform"`, `"projectId":"project-id"`} {
		if !strings.Contains(got, expected) {
			t.Errorf("expected create input to contain %s, got %s", expected, got)
		}
	}
	for _, unexpected := range []string{`"initiativeId"`, `"teamId"`, `"issueId"`, `"releaseId"`, `"cycleId"`} {
		if strings.Contains(got, unexpected) {
			t.Errorf("expected create input to omit %s, got %s", unexpected, got)
		}
	}
}

func TestDocumentUpdateInputIncludesManagedFieldsAndOnlySelectedParent(t *testing.T) {
	t.Parallel()

	input := documentUpdateInputFromModel(&DocumentResourceModel{
		Title:   types.StringValue(""),
		Content: types.StringValue("# Managed by Terraform"),
		IssueId: types.StringValue("LIN-123"),
	})

	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("could not marshal update input: %v", err)
	}
	got := string(encoded)

	for _, expected := range []string{`"title":""`, `"content":"# Managed by Terraform"`, `"issueId":"LIN-123"`} {
		if !strings.Contains(got, expected) {
			t.Errorf("expected update input to contain %s, got %s", expected, got)
		}
	}
	for _, unexpected := range []string{`"initiativeId"`, `"teamId"`, `"releaseId"`, `"cycleId"`, `"projectId"`} {
		if strings.Contains(got, unexpected) {
			t.Errorf("expected update input to omit %s, got %s", unexpected, got)
		}
	}
}

func TestDocumentInputsPreserveTerminalLineEndings(t *testing.T) {
	t.Parallel()

	const content = "body\r\n"

	testCases := map[string]any{
		"create": documentCreateInputFromModel(&DocumentResourceModel{
			Title:     types.StringValue("Document"),
			Content:   types.StringValue(content),
			ProjectId: types.StringValue("project-id"),
		}),
		"update": documentUpdateInputFromModel(&DocumentResourceModel{
			Title:   types.StringValue("Document"),
			Content: types.StringValue(content),
			IssueId: types.StringValue("LIN-123"),
		}),
	}

	for name, input := range testCases {
		name, input := name, input

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encoded, err := json.Marshal(input)
			if err != nil {
				t.Fatalf("could not marshal %s input: %v", name, err)
			}

			var decoded map[string]any
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("could not decode %s input: %v", name, err)
			}
			if got := decoded["content"]; got != content {
				t.Fatalf("expected %s input content %q, got %#v", name, content, got)
			}
		})
	}
}

func TestDocumentResourceCreateReconcilesOnlyReturnedID(t *testing.T) {
	t.Parallel()

	testCases := map[string]DocumentResourceModel{
		"empty title with issue parent": documentResourceModelForMutationTest(
			"",
			"body",
			"issue_id",
			"LIN-123",
		),
		"trailing newline with project parent": documentResourceModelForMutationTest(
			"Planned title",
			"body\n",
			"project_id",
			"53c7964a-5bd4-4679-8cca-a5b78498b2b3",
		),
	}

	for name, planned := range testCases {
		name, planned := name, planned

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			documentResource := documentResourceWithMutationResponses(map[string]string{
				"createDocument": `{
					"documentCreate": {
						"success": true,
						"document": {
							"id": "returned-document-id"
						}
					}
				}`,
			})
			documentSchema, plan := documentResourcePlanForUnitTest(t, planned)
			response := resource.CreateResponse{
				State: tfsdk.State{Schema: documentSchema},
			}

			documentResource.Create(
				context.Background(),
				resource.CreateRequest{Plan: plan},
				&response,
			)

			if response.Diagnostics.HasError() {
				t.Fatalf("unexpected create diagnostics: %v", response.Diagnostics)
			}

			got := documentResourceModelFromState(t, response.State)
			expected := planned
			expected.Id = types.StringValue("returned-document-id")
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("create state should equal the plan plus returned ID\nexpected: %#v\nactual:   %#v", expected, got)
			}
		})
	}
}

func TestDocumentResourceUpdateReconcilesMutationDocument(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		planned          DocumentResourceModel
		response         string
		expectedTitle    string
		expectedBody     string
		expectedParent   string
		expectedParentID string
	}{
		"returned title and parent are authoritative while terminal line endings remain stable": {
			planned: documentResourceModelForMutationTest(
				"Planned title",
				"body\r\n",
				"project_id",
				"53c7964a-5bd4-4679-8cca-a5b78498b2b3",
			),
			response: `{
				"documentUpdate": {
					"success": true,
					"document": {
						"id": "document-id",
						"title": "Returned title",
						"content": "body",
						"trashed": false,
						"issue": {
							"id": "returned-issue-id",
							"identifier": "NEW-1"
						}
					}
				}
			}`,
			expectedTitle:    "Returned title",
			expectedBody:     "body\r\n",
			expectedParent:   "issue_id",
			expectedParentID: "returned-issue-id",
		},
		"substantive returned content is surfaced in state": {
			planned: documentResourceModelForMutationTest(
				"Planned title",
				"planned body\n",
				"project_id",
				"53c7964a-5bd4-4679-8cca-a5b78498b2b3",
			),
			response: `{
				"documentUpdate": {
					"success": true,
					"document": {
						"id": "document-id",
						"title": "Planned title",
						"content": "returned body",
						"trashed": false,
						"project": {
							"id": "53c7964a-5bd4-4679-8cca-a5b78498b2b3"
						}
					}
				}
			}`,
			expectedTitle:    "Planned title",
			expectedBody:     "returned body",
			expectedParent:   "project_id",
			expectedParentID: "53c7964a-5bd4-4679-8cca-a5b78498b2b3",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.planned.Id = types.StringValue("document-id")
			documentResource := documentResourceWithMutationResponses(map[string]string{
				"updateDocument": testCase.response,
			})
			documentSchema, plan := documentResourcePlanForUnitTest(t, testCase.planned)
			response := resource.UpdateResponse{
				State: tfsdk.State{Schema: documentSchema},
			}

			documentResource.Update(
				context.Background(),
				resource.UpdateRequest{Plan: plan},
				&response,
			)

			if response.Diagnostics.HasError() {
				t.Fatalf("unexpected update diagnostics: %v", response.Diagnostics)
			}

			got := documentResourceModelFromState(t, response.State)
			if got.Title.ValueString() != testCase.expectedTitle {
				t.Errorf("expected title %q, got %q", testCase.expectedTitle, got.Title.ValueString())
			}
			if got.Content.ValueString() != testCase.expectedBody {
				t.Errorf("expected content %q, got %q", testCase.expectedBody, got.Content.ValueString())
			}
			for _, parentName := range documentParentAttributeNames {
				parentValue := documentParentModelValue(&got, parentName)
				if parentName == testCase.expectedParent {
					if parentValue.ValueString() != testCase.expectedParentID {
						t.Errorf("expected %s %q, got %q", parentName, testCase.expectedParentID, parentValue.ValueString())
					}
				} else if !parentValue.IsNull() {
					t.Errorf("%s should be null, got %q", parentName, parentValue.ValueString())
				}
			}
		})
	}
}

func TestDocumentMutationError(t *testing.T) {
	t.Parallel()

	if err := documentMutationError("create", true); err != nil {
		t.Fatalf("successful mutation returned an error: %v", err)
	}

	err := documentMutationError("update", false)
	if err == nil {
		t.Fatal("unsuccessful mutation should return an error")
	}
	if !strings.Contains(err.Error(), "update") {
		t.Fatalf("expected operation in error, got %q", err)
	}
}

func TestIsDocumentNotFoundError(t *testing.T) {
	t.Parallel()

	if !isDocumentNotFoundError(errors.New("graphql: Entity not found: Document")) {
		t.Error("expected document not-found error to be recognized")
	}
	if isDocumentNotFoundError(errors.New("connection reset")) {
		t.Error("connection error should not be treated as document not found")
	}
	if isDocumentNotFoundError(nil) {
		t.Error("nil should not be treated as document not found")
	}
}

func documentStringAttributeValidationHasError(t *testing.T, attributeName, value string) bool {
	t.Helper()

	var schemaResponse resource.SchemaResponse
	(&DocumentResource{}).Schema(context.Background(), resource.SchemaRequest{}, &schemaResponse)

	attribute, ok := schemaResponse.Schema.Attributes[attributeName].(schema.StringAttribute)
	if !ok {
		t.Fatalf("%s is not a string attribute", attributeName)
	}

	var response schemavalidator.StringResponse
	request := schemavalidator.StringRequest{
		Path:        path.Root(attributeName),
		ConfigValue: types.StringValue(value),
	}
	for _, attributeValidator := range attribute.Validators {
		attributeValidator.ValidateString(context.Background(), request, &response)
	}

	return response.Diagnostics.HasError()
}

type documentMutationClient struct {
	responses map[string]string
}

func (c *documentMutationClient) MakeRequest(_ context.Context, request *graphql.Request, response *graphql.Response) error {
	payload, ok := c.responses[request.OpName]
	if !ok {
		return fmt.Errorf("unexpected GraphQL operation %q", request.OpName)
	}
	return json.Unmarshal([]byte(payload), response.Data)
}

func documentResourceWithMutationResponses(responses map[string]string) *DocumentResource {
	var client graphql.Client = &documentMutationClient{responses: responses}
	return &DocumentResource{client: &client}
}

func documentResourceModelForMutationTest(title, content, parentName, parentID string) DocumentResourceModel {
	model := DocumentResourceModel{
		Id:        types.StringUnknown(),
		Title:     types.StringValue(title),
		Content:   types.StringValue(content),
		IssueId:   types.StringNull(),
		ProjectId: types.StringNull(),
	}

	switch parentName {
	case "issue_id":
		model.IssueId = types.StringValue(parentID)
	case "project_id":
		model.ProjectId = types.StringValue(parentID)
	default:
		panic("unsupported mutation test parent: " + parentName)
	}

	return model
}

func documentResourcePlanForUnitTest(t *testing.T, model DocumentResourceModel) (schema.Schema, tfsdk.Plan) {
	t.Helper()

	var schemaResponse resource.SchemaResponse
	(&DocumentResource{}).Schema(context.Background(), resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", schemaResponse.Diagnostics)
	}

	plan := tfsdk.Plan{Schema: schemaResponse.Schema}
	if diagnostics := plan.Set(context.Background(), &model); diagnostics.HasError() {
		t.Fatalf("could not set document plan: %v", diagnostics)
	}

	return schemaResponse.Schema, plan
}

func documentResourceModelFromState(t *testing.T, state tfsdk.State) DocumentResourceModel {
	t.Helper()

	var model DocumentResourceModel
	if diagnostics := state.Get(context.Background(), &model); diagnostics.HasError() {
		t.Fatalf("could not get document state: %v", diagnostics)
	}
	return model
}

func documentResourceTestConfig(documentSchema schema.Schema, parents map[string]string) tfsdk.Config {
	attributeTypes := map[string]tftypes.Type{
		"id":         tftypes.String,
		"title":      tftypes.String,
		"content":    tftypes.String,
		"issue_id":   tftypes.String,
		"project_id": tftypes.String,
	}
	values := map[string]tftypes.Value{
		"id":      tftypes.NewValue(tftypes.String, nil),
		"title":   tftypes.NewValue(tftypes.String, ""),
		"content": tftypes.NewValue(tftypes.String, ""),
	}

	for _, name := range documentParentAttributeNames {
		value, ok := parents[name]
		if !ok {
			values[name] = tftypes.NewValue(tftypes.String, nil)
			continue
		}
		values[name] = tftypes.NewValue(tftypes.String, value)
	}

	return tfsdk.Config{
		Schema: documentSchema,
		Raw: tftypes.NewValue(
			tftypes.Object{AttributeTypes: attributeTypes},
			values,
		),
	}
}

func documentParentModelValue(model *DocumentResourceModel, name string) types.String {
	switch name {
	case "issue_id":
		return model.IssueId
	case "project_id":
		return model.ProjectId
	default:
		panic("unexpected parent attribute: " + name)
	}
}

func stringPointer(value string) *string {
	return &value
}
