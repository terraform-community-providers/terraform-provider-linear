resource "linear_team_label" "example" {
  name    = "Tech Debt"
  team_id = linear_team.example.id
}
