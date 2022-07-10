resource "linear_workflow_state" "example" {
  name    = "Deployed"
  type    = "completed"
  color   = "#ffff00"
  team_id = linear_team.example.id
}
