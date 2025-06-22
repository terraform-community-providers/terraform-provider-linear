## 0.3.0

### BREAKING
* Removed `no_priority_issues_first` from `linear_team`

### Enhancements
* Added `enable_thread_summaries`, `auto_close_parent_issues` & `auto_close_child_issues` to `linear_team`
* Added `triage.require_priority` to `linear_team`
* Added `mergeable` to `linear_team_workflow`
* Added `allow_members_to_create_teams` & `allow_members_to_manage_labels` to `linear_workspace_settings`

### Bug Fixes
* Fix issue with creating/updating team when cycles are enabled
* Fix issue with reading/creating/updating team workflow

## 0.2.6

### Bug Fixes
* Fix issue with creating/updating team caused by API change

## 0.2.5

### Enhancements
* Add `allow_members_to_invite` to workspace settings

### Bug Fixes
* Fix issues with API token
* Fix issues with updating workspace settings

## 0.2.3

### Enhancements
* Position of default workflow states is always `0`.

## 0.2.1

### Enhancements
* Improve recognition of existing workflows when importing a team

## 0.2.0 (First release)
