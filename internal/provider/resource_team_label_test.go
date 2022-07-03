package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamLabelResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamLabelResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Tech Debt"),
					resource.TestCheckResourceAttr("linear_team_label.test", "description", ""),
					resource.TestMatchResourceAttr("linear_team_label.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Tech Debt:TEST",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTeamLabelResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Tech Debt"),
					resource.TestCheckResourceAttr("linear_team_label.test", "description", ""),
					resource.TestMatchResourceAttr("linear_team_label.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamLabelResourceConfigNonDefault("Easy Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Easy Tech Debt"),
					resource.TestCheckResourceAttr("linear_team_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_team_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Easy Tech Debt:TEST",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTeamLabelResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamLabelResourceConfigNonDefault("Needs design"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Needs design"),
					resource.TestCheckResourceAttr("linear_team_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_team_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Needs design:TEST",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamLabelResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_label" "test" {
  name = "%s"
  team_id = "4486be5a-706b-47be-81ab-1937d6ecf193"
}
`, name)
}

func testAccTeamLabelResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_label" "test" {
  name = "%s"
  description = "lots of it"
  color = "#00ff00"
  team_id = "4486be5a-706b-47be-81ab-1937d6ecf193"
}
`, name)
}
