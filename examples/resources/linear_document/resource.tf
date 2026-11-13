resource "linear_document" "example" {
  title      = "Engineering handbook"
  project_id = "ecb63be7-63fb-47fe-95fc-1585d24e172d"

  # Terraform owns the complete Markdown body. For larger documents, use file("${path.module}/document.md").
  content = <<-EOT
    # Engineering handbook

    This document is managed by Terraform.
  EOT

  # Destroying this resource moves the document to Linear's trash.
}
