package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceConfigDefault("ACC", "Acc Tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acc Tests"),
					resource.TestCheckResourceAttr("linear_team.test", "description", ""),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Bank"),
					resource.TestMatchResourceAttr("linear_team.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "2"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "notUsed"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "ACC",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTeamResourceConfigDefault("ACC", "Acc Tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acc Tests"),
					resource.TestCheckResourceAttr("linear_team.test", "description", ""),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Bank"),
					resource.TestMatchResourceAttr("linear_team.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "2"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "notUsed"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "1"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceConfigNonDefault("AC", "Acceptance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "AC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "true"),
					// #2
					// resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "2"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "4"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "linear"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "AC",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTeamResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceConfigNonDefault("DEV", "DevOps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "DEV"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "DevOps"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "true"),
					// #2
					// resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "2"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "4"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "linear"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "DEV",
				ImportStateVerify: true,
			},
			// Update with null values
			{
				Config: testAccTeamResourceConfigDefault("DEV", "DevOps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linear_team.test", "key", "DEV"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "DevOps"),
					resource.TestCheckResourceAttr("linear_team.test", "description", ""),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "no_priority_issues_first", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "2"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "notUsed"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamResourceConfigDefault(key string, name string) string {
	return fmt.Sprintf(`
resource "linear_team" "test" {
  key = "%s"
  name = "%s"
}
`, key, name)
}

func testAccTeamResourceConfigNonDefault(key string, name string) string {
	return fmt.Sprintf(`
resource "linear_team" "test" {
  key = "%s"
  name = "%s"
  private = true
  description = "nice team"
  icon = "Image"
  color = "#00ff00"
  timezone = "Europe/London"
  no_priority_issues_first = false
  enable_issue_history_grouping = false
  enable_issue_default_to_bottom = true
  # auto_archive_period = 6

  triage = {
    enabled = true
  }

  cycles = {
    enabled = true
    start_day = 6
    duration = 2
    cooldown = 1
    upcoming = 4
    auto_add_started = false
    auto_add_completed = false
    need_for_active = true
  }

  estimation = {
    type = "linear"
    extended = true
    allow_zero = true
    default = 0
  }
}
`, key, name)
}
