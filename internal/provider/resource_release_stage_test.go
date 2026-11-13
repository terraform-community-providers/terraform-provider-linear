package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReleaseStageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReleaseStageResourceConfig("Staging", "#00ff00", "1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_release_stage.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_release_stage.test", "name", "Staging"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "type", "started"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "position", "1"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "frozen", "false"),
				),
			},
			{
				ResourceName:      "linear_release_stage.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccReleaseStageResourceConfigUpdated("Production", "#ff0000", "2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_release_stage.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_release_stage.test", "name", "Production"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "color", "#ff0000"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "position", "2"),
					resource.TestCheckResourceAttr("linear_release_stage.test", "frozen", "true"),
				),
			},
			{
				ResourceName:      "linear_release_stage.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccReleaseStageResourceConfig(name string, color string, position string) string {
	return `
resource "linear_release_pipeline" "test" {
  name = "Stage Test Pipeline"
  type = "scheduled"
}

resource "linear_release_stage" "keep" {
  name        = "Keep Stage"
  type        = "started"
  color       = "#0000ff"
  position    = 0
  pipeline_id = linear_release_pipeline.test.id
}

resource "linear_release_stage" "test" {
  name        = "` + name + `"
  type        = "started"
  color       = "` + color + `"
  position    = ` + position + `
  pipeline_id = linear_release_pipeline.test.id
}
`
}

func testAccReleaseStageResourceConfigUpdated(name string, color string, position string) string {
	return `
resource "linear_release_pipeline" "test" {
  name = "Stage Test Pipeline"
  type = "scheduled"
}

resource "linear_release_stage" "keep" {
  name        = "Keep Stage"
  type        = "started"
  color       = "#0000ff"
  position    = 0
  pipeline_id = linear_release_pipeline.test.id
}

resource "linear_release_stage" "test" {
  name        = "` + name + `"
  type        = "started"
  color       = "` + color + `"
  position    = ` + position + `
  frozen      = true
  pipeline_id = linear_release_pipeline.test.id
}
`
}
