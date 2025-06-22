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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
				),
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceSettingsResourceConfigNonDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
				),
			},
			// Update with null values
			{
				Config: testAccWorkspaceSettingsResourceConfigDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_workspace_settings.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_invite", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_create_teams", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "allow_members_to_manage_labels", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_roadmap", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
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
	allow_members_to_create_teams = false
	allow_members_to_manage_labels = false
	enable_roadmap = true
	enable_git_linkback_messages = false
	enable_git_linkback_messages_public = true

	projects = {
		update_reminder_day       = "Wednesday"
		update_reminder_hour      = 9
		update_reminder_frequency = 2
	}
}
`
}
