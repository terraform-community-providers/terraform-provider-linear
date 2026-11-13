package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReleasePipelineResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReleasePipelineResourceConfigDefault("Test Pipeline"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_release_pipeline.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "name", "Test Pipeline"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "type", "continuous"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "is_production", "true"),
					resource.TestMatchResourceAttr("linear_release_pipeline.test", "default_stage.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "default_stage.name", "Started"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "default_stage.color", "#f59e0b"),
				),
			},
			{
				ResourceName:      "linear_release_pipeline.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccReleasePipelineResourceConfigNonDefault("Updated Pipeline"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_release_pipeline.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "name", "Updated Pipeline"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "type", "scheduled"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "is_production", "false"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "include_path_patterns.0", "src/**"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "default_stage.name", "Deploy"),
					resource.TestCheckResourceAttr("linear_release_pipeline.test", "default_stage.color", "#10b981"),
				),
			},
			{
				ResourceName:      "linear_release_pipeline.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccReleasePipelineResourceConfigDefault(name string) string {
	return `
resource "linear_release_pipeline" "test" {
  name = "` + name + `"
}
`
}

func testAccReleasePipelineResourceConfigNonDefault(name string) string {
	return `
resource "linear_release_pipeline" "test" {
  name                  = "` + name + `"
  type                  = "scheduled"
  is_production         = false
  include_path_patterns = ["src/**"]

  default_stage = {
    name  = "Deploy"
    color = "#10b981"
  }
}
`
}
