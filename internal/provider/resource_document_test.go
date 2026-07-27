package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	documentProjectIDEnvVar = "LINEAR_TEST_DOCUMENT_PROJECT_ID"
	documentIssueIDEnvVar   = "LINEAR_TEST_DOCUMENT_ISSUE_ID"

	documentInitialContent = "# Terraform acceptance test\n\nCreated through the provider.\n"
	documentUpdatedContent = "# Updated acceptance test\n\nThe complete Markdown body is managed by Terraform.\n"
)

func TestAccDocumentResource(t *testing.T) {
	projectID := os.Getenv(documentProjectIDEnvVar)
	issueID := os.Getenv(documentIssueIDEnvVar)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccDocumentParentPreCheck(t, documentProjectIDEnvVar)
			testAccDocumentParentPreCheck(t, documentIssueIDEnvVar)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with multiline Markdown content.
			{
				Config: testAccDocumentResourceConfig(
					"Terraform Provider Document Acceptance Test",
					documentInitialContent,
					"project_id",
					projectID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_document.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_document.test", "title", "Terraform Provider Document Acceptance Test"),
					resource.TestCheckResourceAttr("linear_document.test", "content", documentInitialContent),
					resource.TestCheckResourceAttrWith("linear_document.test", "project_id", func(value string) error {
						if value != projectID {
							return fmt.Errorf("expected project_id %q, got %q", projectID, value)
						}
						return nil
					}),
					resource.TestCheckNoResourceAttr("linear_document.test", "issue_id"),
				),
			},
			// ImportState testing.
			{
				ResourceName:            "linear_document.test",
				ImportState:             true,
				ImportStateIdFunc:       documentImportIDFunc,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content"},
				ImportStateCheck:        documentImportContentCheck,
			},
			// Update title and complete Markdown content.
			{
				Config: testAccDocumentResourceConfig(
					"Updated Terraform Document",
					documentUpdatedContent,
					"project_id",
					projectID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_document.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_document.test", "title", "Updated Terraform Document"),
					resource.TestCheckResourceAttr("linear_document.test", "content", documentUpdatedContent),
				),
			},
			// Move the document to a different parent type.
			{
				Config: testAccDocumentResourceConfig(
					"Updated Terraform Document",
					documentUpdatedContent,
					"issue_id",
					issueID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_document.test", "id", uuidRegex()),
					resource.TestCheckResourceAttrWith("linear_document.test", "issue_id", func(value string) error {
						if value != issueID {
							return fmt.Errorf("expected issue_id %q, got %q", issueID, value)
						}
						return nil
					}),
					resource.TestCheckNoResourceAttr("linear_document.test", "project_id"),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}

func testAccDocumentParentPreCheck(t *testing.T, envVar string) {
	t.Helper()

	value := os.Getenv(envVar)
	if value == "" {
		t.Skipf("%s must be set to run the document acceptance test", envVar)
	}

	if !uuidRegex().MatchString(value) {
		t.Fatalf("%s must contain a UUID, got %q", envVar, value)
	}
}

func testAccDocumentResourceConfig(
	title string,
	content string,
	parentAttribute string,
	parentID string,
) string {
	return fmt.Sprintf(`
resource "linear_document" "test" {
  title   = %q
  content = %q

  %s = %q
}
`, title, content, parentAttribute, parentID)
}

func documentImportIDFunc(state *terraform.State) (string, error) {
	rawState, ok := state.RootModule().Resources["linear_document.test"]
	if !ok {
		return "", fmt.Errorf("resource not found")
	}

	return rawState.Primary.Attributes["id"], nil
}

func documentImportContentCheck(states []*terraform.InstanceState) error {
	expectedContent := strings.TrimRight(documentInitialContent, "\r\n")

	for _, state := range states {
		if state.Ephemeral.Type != "linear_document" {
			continue
		}

		if content, ok := state.Attributes["content"]; ok {
			if content != expectedContent {
				return fmt.Errorf("expected imported content %q, got %q", expectedContent, content)
			}
			return nil
		}
	}

	return fmt.Errorf("content not found in imported document state")
}
