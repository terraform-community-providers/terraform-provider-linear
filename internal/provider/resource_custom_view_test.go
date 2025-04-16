package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceViewResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceViewResource("name1", "name2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// name1
					resource.TestMatchResourceAttr("linear_custom_view.test1", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_custom_view.test1", "name", "name1"),

					// name2
					resource.TestMatchResourceAttr("linear_custom_view.test2", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_custom_view.test2", "name", "name2"),
				),
			},
			// Update with same values
			{
				Config: testAccWorkspaceViewResource("name2", "name2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// name1
					resource.TestMatchResourceAttr("linear_custom_view.test1", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_custom_view.test1", "name", "name2"),

					// name2
					resource.TestMatchResourceAttr("linear_custom_view.test2", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_custom_view.test2", "name", "name2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceViewResource(name1, name2 string) string {
	return fmt.Sprintf(`
resource "linear_custom_view" "test1" {
  name = "%s"
}

resource "linear_custom_view" "test2" {
  name = "%s"
}
`, name1, name2)
}
