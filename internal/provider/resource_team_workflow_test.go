package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "merge"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamWorkflowResourceConfigNonDefault("DEF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team_workflow.test", "id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "key", "DEF"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "fbb47815-9f97-4f7b-885b-1417a83b57c0"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "b3a08038-c253-4c3b-8019-a985a0ddb6d0"),
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
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "fbb47815-9f97-4f7b-885b-1417a83b57c0"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "b3a08038-c253-4c3b-8019-a985a0ddb6d0"),
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
					resource.TestCheckResourceAttr("linear_team_workflow.test", "draft", "fbb47815-9f97-4f7b-885b-1417a83b57c0"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "start", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "review", "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"),
					resource.TestCheckResourceAttr("linear_team_workflow.test", "merge", "b3a08038-c253-4c3b-8019-a985a0ddb6d0"),
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
					resource.TestCheckNoResourceAttr("linear_team_workflow.test", "merge"),
				),
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
  draft = "fbb47815-9f97-4f7b-885b-1417a83b57c0"
  start = "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"
  review = "4e868c3b-30d2-4d9e-9f1d-a6ed42c7926a"
  merge = "b3a08038-c253-4c3b-8019-a985a0ddb6d0"
}
`, key)
}
