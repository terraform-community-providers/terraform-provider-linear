package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceConfigCreation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance Tests"),
					resource.TestCheckResourceAttr("linear_team.test", "description", ""),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Shop"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "true"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "linear_team.test",
			// 	ImportState:       true,
			// 	ImportStateId:     "ACC",
			// 	ImportStateVerify: true,
			// },
			// Update with null values
			{
				Config: testAccTeamResourceConfigCreation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance Tests"),
					resource.TestCheckResourceAttr("linear_team.test", "description", ""),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Shop"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceConfigUpdation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "AC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "false"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "linear_team.test",
			// 	ImportState:       true,
			// 	ImportStateId:     "AC",
			// 	ImportStateVerify: true,
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamResourceConfigCreation() string {
	return `
resource "linear_team" "test" {
  key = "ACC"
  name = "Acceptance Tests"
}
`
}

func testAccTeamResourceConfigUpdation() string {
	return `
resource "linear_team" "test" {
  key = "AC"
  name = "Acceptance"
  private = true
  description = "nice team"
  icon = "Image"
  color = "#00ff00"
  timezone = "Europe/London"
  enable_issue_history_grouping = false
  no_priority_issues_first = false
}
`
}
