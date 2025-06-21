package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acc Tests"),
					resource.TestCheckNoResourceAttr("linear_team.test", "description"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Bank"),
					resource.TestMatchResourceAttr("linear_team.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_close_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "false"),
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
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Backlog"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bec2c8"),
					resource.TestCheckNoResourceAttr("linear_team.test", "backlog_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Todo"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#e2e2e2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "unstarted_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In Progress"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#f2c94c"),
					resource.TestCheckNoResourceAttr("linear_team.test", "started_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Done"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5e6ad2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "completed_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Canceled"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#95a2b3"),
					resource.TestCheckNoResourceAttr("linear_team.test", "canceled_workflow_state.description"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "ACC",
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccTeamResourceConfigDefault("ACC", "Acc Tests"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "ACC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acc Tests"),
					resource.TestCheckNoResourceAttr("linear_team.test", "description"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Bank"),
					resource.TestMatchResourceAttr("linear_team.test", "color", colorRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_close_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "false"),
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
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Backlog"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bec2c8"),
					resource.TestCheckNoResourceAttr("linear_team.test", "backlog_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Todo"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#e2e2e2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "unstarted_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In Progress"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#f2c94c"),
					resource.TestCheckNoResourceAttr("linear_team.test", "started_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Done"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5e6ad2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "completed_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Canceled"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#95a2b3"),
					resource.TestCheckNoResourceAttr("linear_team.test", "canceled_workflow_state.description"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceConfigNonDefault("AC", "Acceptance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "AC"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "Acceptance"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_close_period", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "4"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "linear"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "0"),
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Icebox"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bbbbbb"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.description", "Not planned"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Ready to start"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#eeeeee"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.description", "Planned"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In flight"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#ffcccc"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.description", "Working on it"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Merged"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5566dd"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.description", "Merged to main"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Invalid"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#99aabb"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.description", "Not valid or not needed"),
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
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "DEV"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "DevOps"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_close_period", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "4"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "linear"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "0"),
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Icebox"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bbbbbb"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.description", "Not planned"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Ready to start"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#eeeeee"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.description", "Planned"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In flight"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#ffcccc"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.description", "Working on it"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Merged"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5566dd"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.description", "Merged to main"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Invalid"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#99aabb"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.description", "Not valid or not needed"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "DEV",
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccTeamResourceConfigNonDefault("DEV", "DevOps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "DEV"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "DevOps"),
					resource.TestCheckResourceAttr("linear_team.test", "description", "nice team"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_close_period", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.enabled", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.start_day", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.duration", "3"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.cooldown", "1"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.upcoming", "4"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_started", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.auto_add_completed", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "cycles.need_for_active", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.type", "linear"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.extended", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.allow_zero", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "estimation.default", "0"),
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Icebox"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bbbbbb"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.description", "Not planned"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Ready to start"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#eeeeee"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.description", "Planned"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In flight"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#ffcccc"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.description", "Working on it"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Merged"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5566dd"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.description", "Merged to main"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Invalid"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#99aabb"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.description", "Not valid or not needed"),
				),
			},
			// Update with null values
			{
				Config: testAccTeamResourceConfigDefault("DEV", "DevOps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "key", "DEV"),
					resource.TestCheckResourceAttr("linear_team.test", "name", "DevOps"),
					resource.TestCheckNoResourceAttr("linear_team.test", "description"),
					resource.TestCheckResourceAttr("linear_team.test", "icon", "Image"),
					resource.TestCheckResourceAttr("linear_team.test", "color", "#00ff00"),
					resource.TestCheckResourceAttr("linear_team.test", "private", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "timezone", "Etc/GMT"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_history_grouping", "true"),
					resource.TestCheckResourceAttr("linear_team.test", "enable_issue_default_to_bottom", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "auto_archive_period", "6"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.enabled", "false"),
					resource.TestCheckResourceAttr("linear_team.test", "triage.require_priority", "false"),
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
					resource.TestMatchResourceAttr("linear_team.test", "backlog_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.name", "Backlog"),
					resource.TestCheckResourceAttr("linear_team.test", "backlog_workflow_state.color", "#bec2c8"),
					resource.TestCheckNoResourceAttr("linear_team.test", "backlog_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "unstarted_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.name", "Todo"),
					resource.TestCheckResourceAttr("linear_team.test", "unstarted_workflow_state.color", "#e2e2e2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "unstarted_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "started_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.name", "In Progress"),
					resource.TestCheckResourceAttr("linear_team.test", "started_workflow_state.color", "#f2c94c"),
					resource.TestCheckNoResourceAttr("linear_team.test", "started_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "completed_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.name", "Done"),
					resource.TestCheckResourceAttr("linear_team.test", "completed_workflow_state.color", "#5e6ad2"),
					resource.TestCheckNoResourceAttr("linear_team.test", "completed_workflow_state.description"),
					resource.TestMatchResourceAttr("linear_team.test", "canceled_workflow_state.id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.position", "0"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.name", "Canceled"),
					resource.TestCheckResourceAttr("linear_team.test", "canceled_workflow_state.color", "#95a2b3"),
					resource.TestCheckNoResourceAttr("linear_team.test", "canceled_workflow_state.description"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linear_team.test",
				ImportState:       true,
				ImportStateId:     "DEV",
				ImportStateVerify: true,
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
  enable_issue_history_grouping = false
  enable_issue_default_to_bottom = true
  auto_archive_period = 3
  auto_close_period = 0

  triage = {
    enabled          = true
    require_priority = true
  }

  cycles = {
    enabled = true
    start_day = 6
    duration = 3
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

  backlog_workflow_state = {
    name = "Icebox"
    color = "#bbbbbb"
    description = "Not planned"
  }

  unstarted_workflow_state = {
    name = "Ready to start"
    color = "#eeeeee"
    description = "Planned"
  }

  started_workflow_state = {
    name = "In flight"
    color = "#ffcccc"
    description = "Working on it"
  }

  completed_workflow_state = {
    name = "Merged"
    color = "#5566dd"
    description = "Merged to main"
  }

  canceled_workflow_state = {
    name = "Invalid"
    color = "#99aabb"
    description = "Not valid or not needed"
  }
}
`, key, name)
}
