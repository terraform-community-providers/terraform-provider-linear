package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTemplateResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTemplateResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_template.test", "description"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckNoResourceAttr("linear_template.test", "team_id"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"title\":\"\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_template.test",
				ImportState:       true,
				ImportStateIdFunc: templateImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTemplateResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_template.test", "description"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckNoResourceAttr("linear_template.test", "team_id"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"title\":\"\"}"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTemplateResourceConfigNonDefault("Easy Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Easy Tech Debt"),
					resource.TestCheckResourceAttr("linear_template.test", "description", "Tech debt that is easy to fix"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckResourceAttr("linear_template.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"labelIds\":[\"53c7964a-5bd4-4679-8cca-a5b78498b2b3\"],\"priority\":0,\"title\":\"\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_template.test",
				ImportState:       true,
				ImportStateIdFunc: templateImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTemplateResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTemplateResourceConfigNonDefault("Easy Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Easy Tech Debt"),
					resource.TestCheckResourceAttr("linear_template.test", "description", "Tech debt that is easy to fix"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckResourceAttr("linear_template.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"labelIds\":[\"53c7964a-5bd4-4679-8cca-a5b78498b2b3\"],\"priority\":0,\"title\":\"\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_template.test",
				ImportState:       true,
				ImportStateIdFunc: templateImportIdFunc,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTemplateResourceConfigNonDefault("Easy Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Easy Tech Debt"),
					resource.TestCheckResourceAttr("linear_template.test", "description", "Tech debt that is easy to fix"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckResourceAttr("linear_template.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"labelIds\":[\"53c7964a-5bd4-4679-8cca-a5b78498b2b3\"],\"priority\":0,\"title\":\"\"}"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTemplateResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_template.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_template.test", "description"),
					resource.TestCheckResourceAttr("linear_template.test", "type", "issue"),
					resource.TestCheckNoResourceAttr("linear_template.test", "team_id"),
					resource.TestCheckResourceAttr("linear_template.test", "data", "{\"title\":\"\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_template.test",
				ImportState:       true,
				ImportStateIdFunc: templateImportIdFunc,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTemplateResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_template" "test" {
  name = "%s"
  type = "issue"
  data = jsonencode({
    "title" = ""
  })
}
`, name)
}

func testAccTemplateResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_template" "test" {
  name = "%s"
  description = "Tech debt that is easy to fix"
  type = "issue"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
  data = jsonencode({
    "title" = ""
    "priority" = 0
    "labelIds" = ["53c7964a-5bd4-4679-8cca-a5b78498b2b3"]
  })
}
`, name)
}

func templateImportIdFunc(state *terraform.State) (string, error) {
	rawState, ok := state.RootModule().Resources["linear_template.test"]

	if !ok {
		return "", fmt.Errorf("Resource Not found")
	}

	return rawState.Primary.Attributes["id"], nil
}
