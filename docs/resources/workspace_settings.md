---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "linear_workspace_settings Resource - terraform-provider-linear"
subcategory: ""
description: |-
  Linear workspace settings.
---

# linear_workspace_settings (Resource)

Linear workspace settings.

## Example Usage

```terraform
resource "linear_workspace_settings" "example" {
  enable_roadmap = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `enable_git_linkback_messages` (Boolean) Enable git linkbacks for private repositories. **Default** `false`.
- `enable_git_linkback_messages_public` (Boolean) Enable git linkbacks for public repositories. **Default** `false`.
- `enable_roadmap` (Boolean) Enable roadmap for the workspace. **Default** `false`.

### Read-Only

- `id` (String) Identifier of the workspace.

## Import

Import is supported using the following syntax:

```shell
terraform import linear_workspace_settings.example
```