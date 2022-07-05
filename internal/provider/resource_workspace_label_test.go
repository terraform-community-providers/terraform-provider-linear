package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("linear_workspace_label.test", "description", ""),
					resource.TestMatchResourceAttr("linear_workspace_label.test", "color", colorRegex()),
				),
			},
			// ImportState testing
			// #4
			// {
			// 	ResourceName:      "linear_workspace_label.test",
			// 	ImportState:       true,
			// 	ImportStateId:     "UX",
			// 	ImportStateVerify: true,
			// },
			// Update with null values
			{
				Config: testAccWorkspaceLabelResourceConfigDefault("UX"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_label.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "name", "UX"),
					resource.TestCheckResourceAttr("linear_workspace_label.test", "description", ""),
					resource.TestMatchResourceAttr("linear_workspace_label.test", "color", colorRegex()),
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
				),
			},
			// ImportState testing
			// #4
			// {
			// 	ResourceName:      "linear_workspace_label.test",
			// 	ImportState:       true,
			// 	ImportStateId:     "Easy UX",
			// 	ImportStateVerify: true,
			// },
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
				),
			},
			// ImportState testing
			// #4
			// {
			// 	ResourceName:      "linear_workspace_label.test",
			// 	ImportState:       true,
			// 	ImportStateId:     "Needs product",
			// 	ImportStateVerify: true,
			// },
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
}
`, name)
}
