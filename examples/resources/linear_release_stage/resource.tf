resource "linear_release_stage" "qa" {
  name        = "QA Verified"
  color       = "#4cb782"
  position    = 1000
  pipeline_id = linear_release_pipeline.production.id
}

resource "linear_release_stage" "staging" {
  name        = "Staging"
  color       = "#3b82f6"
  position    = 2000
  pipeline_id = linear_release_pipeline.production.id
}
