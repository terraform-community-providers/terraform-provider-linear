package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	documentInitialContent = "# Terraform acceptance test\n\nCreated through the provider.\n"
	documentUpdatedContent = "# Updated acceptance test\n\nThe complete Markdown body is managed by Terraform.\n"
)

func TestAccDocumentResource(t *testing.T) {
	var fixture documentAcceptanceFixture

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			client := newDocumentAcceptanceClient()
			t.Cleanup(func() {
				if err := fixture.cleanup(context.Background(), client); err != nil {
					t.Errorf("could not clean up document acceptance fixtures: %v", err)
				}
			})

			fixtureName := fmt.Sprintf("Terraform document acceptance %d", time.Now().UnixNano())
			if err := fixture.setup(context.Background(), client, fixtureName); err != nil {
				t.Fatalf("could not set up document acceptance fixtures: %v", err)
			}

			t.Setenv("TF_VAR_document_project_id", fixture.projectID)
			t.Setenv("TF_VAR_document_issue_id", fixture.issueID)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with multiline Markdown content.
			{
				Config: testAccDocumentResourceConfig(
					"Terraform Provider Document Acceptance Test",
					documentInitialContent,
					"project_id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_document.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_document.test", "title", "Terraform Provider Document Acceptance Test"),
					resource.TestCheckResourceAttr("linear_document.test", "content", documentInitialContent),
					resource.TestCheckResourceAttrWith("linear_document.test", "project_id", func(value string) error {
						if value != fixture.projectID {
							return fmt.Errorf("expected project_id %q, got %q", fixture.projectID, value)
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
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_document.test", "id", uuidRegex()),
					resource.TestCheckResourceAttrWith("linear_document.test", "issue_id", func(value string) error {
						if value != fixture.issueID {
							return fmt.Errorf("expected issue_id %q, got %q", fixture.issueID, value)
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

func testAccDocumentResourceConfig(
	title string,
	content string,
	parentAttribute string,
) string {
	return fmt.Sprintf(`
variable "document_project_id" {
  type = string
}

variable "document_issue_id" {
  type = string
}

resource "linear_document" "test" {
  title   = %q
  content = %q

  %s = var.document_%s
}
`, title, content, parentAttribute, parentAttribute)
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
