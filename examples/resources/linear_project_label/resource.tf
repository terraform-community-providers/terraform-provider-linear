resource "linear_project_label" "release_tier" {
  name     = "Release Tier"
  is_group = true
}

resource "linear_project_label" "tier_1" {
  name      = "T1"
  parent_id = linear_project_label.release_tier.id
}
