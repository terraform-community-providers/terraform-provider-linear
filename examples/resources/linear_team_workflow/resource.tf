resource "linear_team_workflow" "example" {
  key   = "SOME"
  draft = linear_team.example.started_workflow_state.id
  merge = linear_team.example.completed_workflow_state.id
}
