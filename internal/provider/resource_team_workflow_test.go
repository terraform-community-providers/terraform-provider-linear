package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamWorkflowResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamWorkflowResourceConfigDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "draft"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "start"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "review"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "mergeable"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "merge"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_workflow.test",
				ImportState:       true,
				ImportStateId:     "DEF",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTeamWorkflowResourceConfigDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "draft"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "start"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "review"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "mergeable"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "merge"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamWorkflowResourceConfigNonDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "5dbca6c1-9ee2-4bf7-a275-8b69ae27ad14"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "mergeable", "53099a59-c811-4b9c-8016-5443ce513de4"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "66df5c88-cae8-416b-b4e9-85a42b159e18"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_workflow.test",
				ImportState:       true,
				ImportStateId:     "DEF",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTeamWorkflowResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamWorkflowResourceConfigNonDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "5dbca6c1-9ee2-4bf7-a275-8b69ae27ad14"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "mergeable", "53099a59-c811-4b9c-8016-5443ce513de4"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "66df5c88-cae8-416b-b4e9-85a42b159e18"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_workflow.test",
				ImportState:       true,
				ImportStateId:     "DEF",
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccTeamWorkflowResourceConfigNonDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "5dbca6c1-9ee2-4bf7-a275-8b69ae27ad14"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "mergeable", "53099a59-c811-4b9c-8016-5443ce513de4"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "66df5c88-cae8-416b-b4e9-85a42b159e18"),
				),
			},
			// Update with null values
			{
				Config: testAccTeamWorkflowResourceConfigDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "draft"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "start"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "review"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "mergeable"),
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "merge"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_workflow.test",
				ImportState:       true,
				ImportStateId:     "DEF",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamWorkflowResourceConfigDefault(key string) string {
	return fmt.Sprintf(`
resource "linear_team_workflow" "test" {
  key = "%s"
}
`, key)
}

func testAccTeamWorkflowResourceConfigNonDefault(key string) string {
	return fmt.Sprintf(`
resource "linear_team_workflow" "test" {
  key = "%s"
  draft = "5dbca6c1-9ee2-4bf7-a275-8b69ae27ad14"
  start = "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"
  review = "9b6fdbd0-fd66-4ea2-a01d-a24ecf0c1191"
  mergeable = "53099a59-c811-4b9c-8016-5443ce513de4"
  merge = "66df5c88-cae8-416b-b4e9-85a42b159e18"
}
`, key)
}
