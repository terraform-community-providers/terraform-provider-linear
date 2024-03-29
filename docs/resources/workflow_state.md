---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "linear_workflow_state Resource - terraform-provider-linear"
subcategory: ""
description: |-
  Linear team workflow state.
---

# linear_workflow_state (Resource)

Linear team workflow state.

## Example Usage

```terraform
resource "linear_workflow_state" "example" {
  name    = "Deployed"
  type    = "completed"
  color   = "#ffff00"
  team_id = linear_team.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `color` (String) Color of the workflow state.
- `name` (String) Name of the workflow state.
- `position` (Number) Position of the workflow state.
- `team_id` (String) Identifier of the team.
- `type` (String) Type of the workflow state.

### Optional

- `description` (String) Description of the workflow state.

### Read-Only

- `id` (String) Identifier of the workflow state.

## Import

Import is supported using the following syntax:

```shell
terraform import linear_worflow_state.example Done:SOME
```
