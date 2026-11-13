package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestProjectLabelUpdateInputOmitsRetiredAt(t *testing.T) {
	payload, err := json.Marshal(ProjectLabelUpdateInput{
		Name:    "Release Tier",
		IsGroup: true,
	})
	if err != nil {
		t.Fatalf("marshaling update input: %v", err)
	}
	if strings.Contains(string(payload), `"retiredAt"`) {
		t.Fatalf("ordinary update includes retiredAt: %s", payload)
	}
}

func TestDeleteProjectLabel(t *testing.T) {
	tests := map[string]struct {
		responseBody string
		wantError    string
	}{
		"already missing": {
			responseBody: `{"errors":[{"message":"Entity not found: ProjectLabel"}]}`,
		},
		"unsuccessful": {
			responseBody: `{"data":{"projectLabelDelete":{"success":false}}}`,
			wantError:    "Linear reported an unsuccessful deletion",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := graphql.NewClient("https://api.linear.app/graphql", &http.Client{
				Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader(test.responseBody)),
					}, nil
				}),
			})

			err := deleteProjectLabel(context.Background(), client, "7fbfc4ca-0ba4-4ad4-9729-ecdd0ab33d75")
			if test.wantError == "" && err != nil {
				t.Fatalf("deleteProjectLabel() error = %v", err)
			}
			if test.wantError != "" && (err == nil || !strings.Contains(err.Error(), test.wantError)) {
				t.Fatalf("deleteProjectLabel() error = %v, want error containing %q", err, test.wantError)
			}
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestAccProjectLabelResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectLabelResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_project_label.group", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_project_label.group", "name", "Release Tier"),
					resource.TestCheckResourceAttr("linear_project_label.group", "description", "Release tier group"),
					resource.TestMatchResourceAttr("linear_project_label.group", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_project_label.group", "is_group", "true"),
					resource.TestCheckNoResourceAttr("linear_project_label.group", "parent_id"),
					resource.TestMatchResourceAttr("linear_project_label.child", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_project_label.child", "name", "T1"),
					resource.TestCheckResourceAttr("linear_project_label.child", "is_group", "false"),
					resource.TestCheckResourceAttrPair("linear_project_label.child", "parent_id", "linear_project_label.group", "id"),
				),
			},
			{
				Config: testAccProjectLabelResourceConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_project_label.group", "name", "Release Tier Updated"),
					resource.TestCheckResourceAttr("linear_project_label.group", "color", "#5e6ad2"),
					resource.TestCheckResourceAttr("linear_project_label.child", "name", "Tier 1"),
					resource.TestCheckResourceAttrPair("linear_project_label.child", "parent_id", "linear_project_label.group", "id"),
				),
			},
			{
				ResourceName:      "linear_project_label.group",
				ImportState:       true,
				ImportStateIdFunc: projectLabelImportIdFunc("linear_project_label.group"),
				ImportStateVerify: true,
			},
			{
				ResourceName:      "linear_project_label.child",
				ImportState:       true,
				ImportStateIdFunc: projectLabelImportIdFunc("linear_project_label.child"),
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectLabelResourceConfigCleared(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("linear_project_label.child", "description"),
					resource.TestCheckNoResourceAttr("linear_project_label.child", "parent_id"),
				),
			},
			{
				Config: testAccProjectLabelResourceConfigCleared(),
				Check:  testAccDeleteProjectLabel("linear_project_label.child"),
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check:              testAccCheckProjectLabelRemoved("linear_project_label.child"),
			},
		},
	})
}

func testAccProjectLabelResourceConfig() string {
	return `
resource "linear_project_label" "group" {
  name        = "Release Tier"
  description = "Release tier group"
  is_group    = true
}

resource "linear_project_label" "child" {
  name        = "T1"
  description = "Highest release tier"
  color       = "#eb5757"
  parent_id   = linear_project_label.group.id
}
`
}

func testAccProjectLabelResourceConfigUpdated() string {
	return `
resource "linear_project_label" "group" {
  name        = "Release Tier Updated"
  description = "Release tier group"
  color       = "#5e6ad2"
  is_group    = true
}

resource "linear_project_label" "child" {
  name        = "Tier 1"
  description = "Highest release tier"
  color       = "#eb5757"
  parent_id   = linear_project_label.group.id
}
`
}

func testAccProjectLabelResourceConfigCleared() string {
	return `
resource "linear_project_label" "group" {
  name        = "Release Tier Updated"
  description = "Release tier group"
  color       = "#5e6ad2"
  is_group    = true
}

resource "linear_project_label" "child" {
  name  = "Tier 1"
  color = "#eb5757"
}
`
}

func projectLabelImportIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}
		return resourceState.Primary.Attributes["id"], nil
	}
}

func testAccDeleteProjectLabel(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		client := graphql.NewClient("https://api.linear.app/graphql", &http.Client{
			Transport: &authedTransport{
				token:   os.Getenv("LINEAR_TOKEN"),
				wrapped: http.DefaultTransport,
			},
		})
		response, err := projectLabelDelete(context.Background(), client, resourceState.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		if !response.ProjectLabelDelete.Success {
			return fmt.Errorf("deleting %s: Linear reported an unsuccessful deletion", resourceName)
		}
		return nil
	}
}

func testAccCheckProjectLabelRemoved(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if _, ok := state.RootModule().Resources[resourceName]; ok {
			return fmt.Errorf("resource %s still present after refresh", resourceName)
		}
		return nil
	}
}
