resource "linear_team" "parent" {
  key  = "ENG"
  name = "Engineering"
}

resource "linear_team" "child" {
  key                       = "MOB"
  name                      = "Mobile"
  parent_id                 = linear_team.parent.id
  inherit_issue_estimation  = true
  inherit_workflow_statuses = true
}
