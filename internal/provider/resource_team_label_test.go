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
					resource.TestCheckNoResourceAttr("linear_team_label.test", "description"),
					resource.TestMatchResourceAttr("linear_team_label.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Tech Debt:DEF",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTeamLabelResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_team_label.test", "description"),
					resource.TestMatchResourceAttr("linear_team_label.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
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
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Easy Tech Debt:DEF",
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
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team_label.test",
				ImportState:       true,
				ImportStateId:     "Needs design:DEF",
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccTeamLabelResourceConfigNonDefault("Needs design"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Needs design"),
					resource.TestCheckResourceAttr("linear_team_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_team_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// Update with null values
			{
				Config: testAccTeamLabelResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_team_label.test", "description"),
					resource.TestMatchResourceAttr("linear_team_label.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team_label.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamLabelResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_label" "test" {
  name = "%s"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
}
`, name)
}

func testAccTeamLabelResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_label" "test" {
  name = "%s"
  description = "lots of it"
  color = "#00ff00"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
}
`, name)
}
