package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceLabelResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceLabelResourceConfigDefault("UX"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "UX"),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "description"),
					resource.TestMatchResourceAttr("linear_workspace_label.test", "color", colorRegex()),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "parent_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_label.test",
				ImportState:       true,
				ImportStateId:     "UX",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccWorkspaceLabelResourceConfigDefault("UX"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "UX"),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "description"),
					resource.TestMatchResourceAttr("linear_workspace_label.test", "color", colorRegex()),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "parent_id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceLabelResourceConfigNonDefault("Easy UX"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "Easy UX"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "parent_id", "09b38784-8d8b-453a-83b6-84c08d094803"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_label.test",
				ImportState:       true,
				ImportStateId:     "Easy UX",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccWorkspaceLabelResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceLabelResourceConfigNonDefault("Needs product"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "Needs product"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "parent_id", "09b38784-8d8b-453a-83b6-84c08d094803"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_label.test",
				ImportState:       true,
				ImportStateId:     "Needs product",
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccWorkspaceLabelResourceConfigNonDefault("Needs product"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "Needs product"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "parent_id", "09b38784-8d8b-453a-83b6-84c08d094803"),
				),
			},
			// Update with null values
			{
				Config: testAccWorkspaceLabelResourceConfigDefault("UX"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "UX"),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "description"),
					resource.TestMatchResourceAttr("linear_workspace_label.test", "color", colorRegex()),
					resource.TestCheckNoResourceAttr("linear_workspace_label.test", "parent_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_label.test",
				ImportState:       true,
				ImportStateId:     "UX",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceLabelResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_workspace_label" "test" {
  name = "%s"
}
`, name)
}

func testAccWorkspaceLabelResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_workspace_label" "test" {
  name = "%s"
  description = "lots of it"
  color = "#00ff00"
  parent_id = "09b38784-8d8b-453a-83b6-84c08d094803"
}
`, name)
}
