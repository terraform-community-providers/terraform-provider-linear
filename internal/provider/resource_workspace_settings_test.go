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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "daily"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "false"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "daily"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "false"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "7"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "weekly"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "true"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "7"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "weekly"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "true"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "7"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Wednesday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "9"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "2"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "weekly"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "true"),
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
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages", "true"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "enable_git_linkback_messages_public", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "fiscal_year_start_month", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "projects.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_day", "Friday"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_hour", "14"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "initiatives.update_reminder_frequency", "0"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.enabled", "false"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "feed.schedule", "daily"),
					resource.TestCheckResourceAttr("linear_workspace_settings.test", "customers.enabled", "false"),
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
	enable_git_linkback_messages = false
	enable_git_linkback_messages_public = true
	fiscal_year_start_month = 7

	projects = {
		update_reminder_day       = "Wednesday"
		update_reminder_hour      = 9
		update_reminder_frequency = 2
	}

	initiatives = {
		enabled                   = true
		update_reminder_day       = "Wednesday"
		update_reminder_hour      = 9
		update_reminder_frequency = 2
	}

	feed = {
		enabled  = true
		schedule = "weekly"
	}

	customers = {
		enabled = true
	}
}
`
}
