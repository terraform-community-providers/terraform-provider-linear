# Changelog

## 0.3.3

### Enhancements
* Allow `linear_team` to be updated by non-admins

## 0.3.1

### Enhancements
* Added `linear_template` resource

## 0.3.0

### BREAKING
* Removed `no_priority_issues_first` from `linear_team`
* Move `enable_roadmap` to `initiatives.enabled` in `linear_workspace_settings`

### Enhancements
* Added `parent_id` to `linear_team` inorder to create sub-teams
* Added `enable_thread_summaries`, `auto_close_parent_issues`, `auto_close_child_issues` & `triage.require_priority` to `linear_team`
* Added `mergeable` to `linear_team_workflow`
* Added `branch` nested block to `linear_team_workflow`
* Added `allow_members_to_create_teams`, `allow_members_to_manage_labels` & `fiscal_year_start_month` to `linear_workspace_settings`
* Added `projects`, `initiatives`, `feed` & `customers` nested blocks to `linear_workspace_settings`
* Allow `triage` in `type` for `linear_workflow_state`

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
