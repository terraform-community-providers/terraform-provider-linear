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
					resource.TestCheckResourceAttr("linear_team.test", "id", "example-id"),
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance Tests"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "linear_team.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	ImportStateVerifyIgnore: []string{"configurable_attribute"},
			// },
			// Update and Read testing
			// {
			// 	Config: testAccTeamResourceConfigUpdation(),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("linear_team.test", "configurable_attribute", "two"),
			// 	),
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
  key = "ACC"
  name = "Acceptance"
}
`
}
