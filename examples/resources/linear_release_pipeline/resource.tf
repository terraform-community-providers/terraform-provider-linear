resource "linear_release_pipeline" "production" {
  name          = "Production"
  type          = "continuous"
  is_production = true
  teams         = [linear_team.engineering.id]

  default_stage = {
    name  = "Started"
    color = "#f59e0b"
  }
}

resource "linear_release_pipeline" "staging" {
  name = "Staging"
  type = "scheduled"

  include_path_patterns = ["src/**", "lib/**"]
}
