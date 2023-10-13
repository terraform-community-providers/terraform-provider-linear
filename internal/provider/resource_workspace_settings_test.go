package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceSettingsResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceSettingsResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccWorkspaceSettingsResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceSettingsResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccWorkspaceSettingsResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceSettingsResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccWorkspaceSettingsResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
				),
			},
			// Update with null values
			{
				Config: testAccWorkspaceSettingsResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_workspace_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceSettingsResourceConfigDefault() string {
	return `
resource "linear_workspace_settings" "test" {
}
`
}

func testAccWorkspaceSettingsResourceConfigNonDefault() string {
	return `
resource "linear_workspace_settings" "test" {
	allow_members_to_invite = false
	enable_roadmap = true
	enable_git_linkback_messages = false
	enable_git_linkback_messages_public = true
}
`
}
