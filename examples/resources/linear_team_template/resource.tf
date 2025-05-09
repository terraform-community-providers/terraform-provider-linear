variable "your_team_id" {
  description = "The ID of the Linear team"
  type        = string
  default     = "uuid-goes-here"
}

resource "linear_team_template" "test_template" {
  name = "Example template"
  # description = "Test Description" # optional
  template_data = jsonencode({
    "title" = "Test Title"
  })
  # can also do:
  # template_data = file("your-template.json")
  team_id = var.your_team_id
  type    = "issue"
}

