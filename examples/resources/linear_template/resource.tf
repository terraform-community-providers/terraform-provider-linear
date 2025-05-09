resource "linear_template" "example" {
  name = "Team Example template"
  data = jsonencode({
    "title" = "Product Bug"
  })
  # can also do:
  # data = file("your-template.json")
}
