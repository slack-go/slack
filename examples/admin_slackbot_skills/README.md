# Admin Slackbot Skills examples

Run from the repository root. Set `SLACK_ADMIN_TOKEN` to an org admin/owner's
**user token from an org-level app installation** with `admin.skills:read` and
`admin.skills:write`. An Enterprise plan and the Skills API feature are required.
Workspace-level installations cannot obtain these scopes.

| Action | Client method | Scope |
| --- | --- | --- |
| `list` | `AdminSlackbotSkillsList` | `admin.skills:read` |
| `publish` | `AdminSlackbotSkillsPublish` | `admin.skills:write` |
| `set-auto-enabled` | `AdminSlackbotSkillsSetAutoEnabled` | `admin.skills:write` |
| `unpublish` | `AdminSlackbotSkillsUnpublish` | `admin.skills:write` |

All have `Context` variants, matching the paired methods used by the library's
admin app-list APIs and most other API groups. Inventory and deployment optional
fields use variadic `MethodNameOption` functions; omitted values use Slack's defaults.
Publish/unpublish return a response even on Slack errors so callers can inspect
`FailedSkills`. Deployment errors retain `UsergroupIDs` and `InvalidUsergroupIDs`.
Requests use the existing form transport; array arguments are JSON-encoded form
fields. This is confined to the new skills methods.

References: [skills guide](https://docs.slack.dev/ai/slackbot-skills) and
[`admin.slackbot.skills.list`](https://docs.slack.dev/reference/methods/admin.slackbot.skills.list).

`token_revoked` requires replacing the exported `SLACK_ADMIN_TOKEN` after
reauthorizing the org-level installation. Use a current org admin/owner user
access token, not a bot token or refresh token. The example prints a recovery
hint without displaying the credential.

## Inventory — `admin.slackbot.skills.list`

```sh
go run ./examples/admin_slackbot_skills -action list -all -limit 100 \
  -category engineering -creators U_CREATOR,W_CREATOR \
  -publish-status unpublished -auto-enabled-status all \
  -sort date_updated -sort-dir desc
```

Every optional argument has a flag; `-cursor NEXT_CURSOR` resumes a listing.
With `-all`, filters and sorting remain fixed, and empty pages with cursors do
not terminate pagination. Inventory includes invite-only skill metadata without
granting canvas access. Index updates are asynchronous and counts approximate.

## Publish — `admin.slackbot.skills.publish`

```sh
go run ./examples/admin_slackbot_skills -action publish -skills F_SKILL,F_OTHER
go run ./examples/admin_slackbot_skills -action list -publish-status published -all
```

Publishing adds skills to the org catalog and shares their canvases (including
instructions) org-wide if needed. The caller must be able to edit and share each
canvas requiring that change. It does not enable the skills for users. Up to 50
IDs are accepted; duplicates count toward the limit. On `failed_for_some_skills`,
the CLI prints every `failed_skills` item before exiting nonzero. Successful
items and completed sharing changes are retained.

## Auto-enable — `admin.slackbot.skills.setAutoEnabled`

```sh
go run ./examples/admin_slackbot_skills -action set-auto-enabled -skill F_SKILL \
  -permission named_entities -usergroups S_GROUP
go run ./examples/admin_slackbot_skills -action set-auto-enabled -skill F_SKILL \
  -permission no_one
```

Use `-permission everyone` to deploy organization-wide. The canvas must already
be accessible to the selected audience; this method never grants access. Omit
`-usergroups` for `everyone` and `no_one`. Named groups accept 1–100 IDs. Without
org-wide sharing, validation allows at most 1,000 member checks (overlap counts
per group). Responses show the resulting deployment or inaccessible/invalid
group IDs on errors. Publication and deployment are independent.

## Unpublish — `admin.slackbot.skills.unpublish`

```sh
go run ./examples/admin_slackbot_skills -action unpublish -skills F_SKILL,F_OTHER
go run ./examples/admin_slackbot_skills -action list -publish-status unpublished -all
```

Unpublishing preserves canvas sharing, personal libraries, and auto-enablement.
Already-unpublished skills succeed. Batch errors include individual failed IDs,
and successful changes are not rolled back.

## Complete admin group verification

For this workflow, the same `SLACK_ADMIN_TOKEN` also needs `skills:write` for
fixture creation/deletion. Supply a test group in your organization:

```sh
go run ./examples/admin_slackbot_skills -action workflow -usergroups S_TEST_GROUP
```

The workflow creates a uniquely named temporary skill, finds it in inventory,
publishes it (sharing the fixture's instructions org-wide), deploys it to the
selected groups, verifies the deployment in inventory, unpublishes it and checks
that deployment remains, disables auto-enablement, and deletes the skill. It
does not invoke the skill. Inventory verification follows every page, retries at
three-second intervals, and allows up to one minute per state for index updates.

If a step fails or you press Ctrl-C, cleanup attempts to disable, unpublish, and
delete the fixture. Its ID is printed for recovery if cleanup fails. Success ends
with `PASS`. Individual commands accept existing IDs; the workflow only changes
the fixture it creates.
