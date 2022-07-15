package provider

import (
	"fmt"
)

// func TestAccWorkflowStateResourceDefault(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: testAccWorkflowStateResourceConfigDefault("Draft"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Draft"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", ""),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#ffff00"),
// 					resource.TestCheckResourceAttrSet("linear_workflow_state.test", "position"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "default", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "draft", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "start", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "review", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "merge", "false"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:      "linear_workflow_state.test",
// 				ImportState:       true,
// 				ImportStateId:     "Draft:TEST",
// 				ImportStateVerify: true,
// 			},
// 			// Update with null values
// 			{
// 				Config: testAccWorkflowStateResourceConfigDefault("Draft"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Draft"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", ""),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#ffff00"),
// 					resource.TestCheckResourceAttrSet("linear_workflow_state.test", "position"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "default", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "draft", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "start", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "review", "false"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "merge", "false"),
// 				),
// 			},
// 			// Update and Read testing
// 			{
// 				Config: testAccWorkflowStateResourceConfigNonDefault("In review", "started", "For review"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "In review"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "started"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", "For review"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#00ffff"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "position", "20"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "default", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "draft", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "start", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "review", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "merge", "true"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:      "linear_workflow_state.test",
// 				ImportState:       true,
// 				ImportStateId:     "In review:TEST",
// 				ImportStateVerify: true,
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }

// func TestAccWorkflowStateResourceNonDefault(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: testAccWorkflowStateResourceConfigNonDefault("Deployed", "completed", "Deployed to prod"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("linear_workflow_state.test", "id", uuidRegex()),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "name", "Deployed"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "type", "completed"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "description", "Deployed to prod"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "color", "#00ffff"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "position", "20"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "team_id", "4486be5a-706b-47be-81ab-1937d6ecf193"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "default", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "draft", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "start", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "review", "true"),
// 					resource.TestCheckResourceAttr("linear_workflow_state.test", "merge", "true"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:      "linear_workflow_state.test",
// 				ImportState:       true,
// 				ImportStateId:     "Deployed:TEST",
// 				ImportStateVerify: true,
// 			},
// 			// TODO:(PR) Make true -> false
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }

func testAccWorkflowStateResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_workflow_state" "test" {
  name = "%s"
  type = "started"
  color = "#ffff00"
  team_id = "4486be5a-706b-47be-81ab-1937d6ecf193"
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
  team_id = "4486be5a-706b-47be-81ab-1937d6ecf193"
  default = true
  draft = true
  start = true
  review = true
  merge = true
}
`, name, ty, description)
}
