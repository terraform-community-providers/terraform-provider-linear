package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowStateResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowStateResourceConfigDefault("Draft"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Draft"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", ""),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#ffff00"),
					resource.TestCheckResourceAttrSet("linear_workflow_state.test", "position"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workflow_state.test",
				ImportState:       true,
				ImportStateId:     "Draft:DEF",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccWorkflowStateResourceConfigDefault("Draft"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Draft"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", ""),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#ffff00"),
					resource.TestCheckResourceAttrSet("linear_workflow_state.test", "position"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// Update and Read testing
			{
				Config: testAccWorkflowStateResourceConfigNonDefault("In review", "started", "For review"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "In review"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", "For review"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#00ffff"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "position", "20"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workflow_state.test",
				ImportState:       true,
				ImportStateId:     "In review:DEF",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccWorkflowStateResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowStateResourceConfigNonDefault("Deployed", "completed", "Deployed to prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Deployed"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "completed"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", "Deployed to prod"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#00ffff"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "position", "20"),
					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workflow_state.test",
				ImportState:       true,
				ImportStateId:     "Deployed:DEF",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkflowStateResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_workflow_state" "test" {
  name = "%s"
  type = "started"
  color = "#ffff00"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
}
`, name)
}

func testAccWorkflowStateResourceConfigNonDefault(name string, ty string, description string) string {
	return fmt.Sprintf(`
resource "linear_workflow_state" "test" {
  name = "%s"
  type = "%s"
  description = "%s"
  color = "#00ffff"
  position = 20
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
}
`, name, ty, description)
}
