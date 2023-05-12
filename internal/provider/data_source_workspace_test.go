package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccWorkspaceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.linear_workspace.test", "id", "1e73fcad-aac6-4bbe-a5e1-e08cffe04eb5"),
					resource.TestCheckResourceAttr("data.linear_workspace.test", "name", "terraform"),
					resource.TestCheckResourceAttr("data.linear_workspace.test", "url_key", "terraform-test"),
				),
			},
		},
	})
}

const testAccWorkspaceDataSourceConfig = `
data "linear_workspace" "test" {}
`
