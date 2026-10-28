# Slackbot Skills examples

Run from the repository root. Set `SLACK_USER_TOKEN` to a user token with
`skills:read` and `skills:write`; the Skills API must be enabled for your org.
Tokens are read only from the environment. Commands print decoded responses,
including structured Slack errors, and exit nonzero on failure.

The examples contain their sample skill instructions in Go. No Markdown file or
other content file needs to be created before running them.

| Action | Client method | Scope |
| --- | --- | --- |
| `create` | `CreateSlackbotSkill` | `skills:write` |
| `info` | `GetSlackbotSkillInfo` | `skills:read` |
| `update` | `UpdateSlackbotSkill` | `skills:write` |
| `list` | `ListSlackbotSkills` | `skills:read` |
| `library-add` | `AddSlackbotSkillToLibrary` | `skills:write` |
| `library-remove` | `RemoveSlackbotSkillFromLibrary` | `skills:write` |
| `delete` | `DeleteSlackbotSkill` | `skills:write` |

Every method also has a `Context` variant. Methods with optional arguments take
variadic `MethodNameOption` functions. Required arguments remain positional.
For example:

```go
updated, err := api.UpdateSlackbotSkillContext(ctx, skillID,
    slack.UpdateSlackbotSkillOptionDescription(""),
    slack.UpdateSlackbotSkillOptionMCPServers(),
)
// Inspect updated even when err != nil: Slack can return structured details.
```

Categories: `communication`, `customer_support`, `data_and_analytics`,
`engineering`, `product_management`, `productivity`, `professional_development`,
`project_management`, `sales`, `system`, `team_management`, and `other`.

References: [skills guide](https://docs.slack.dev/ai/slackbot-skills) and
[`slackbot.skills.create`](https://docs.slack.dev/reference/methods/slackbot.skills.create).

## Authentication troubleshooting

`token_revoked` means Slack rejected the value of `SLACK_USER_TOKEN`; retrying
the same token will not fix it. Reauthorize/reinstall your app, obtain its current
**User OAuth Token** with `skills:read` and `skills:write`, and replace the value
exported in the terminal running this example. For rotating tokens, use the
current access token from your OAuth flow, not the refresh token. Do not use a
bot token, app-level token, or app client secret.

A failure at `create` with no returned skill ID creates no fixture for this
workflow to clean up. The CLI omits the empty skill object on such errors.
These new methods use the library's existing form transport. Array fields are
JSON-encoded within the form; omitted fields, empty strings, and empty arrays
remain distinct. The encoding does not repair revoked credentials. Any
`missing_charset` warning should be assessed separately from the API error.

## Create — `slackbot.skills.create`

```sh
go run ./examples/slackbot_skills -action create -name 'Release review' \
  -description 'Review a proposed release' -category engineering \
  -example-prompt 'Run {skill} on this release'
```

This uses the sample instructions built into the script. Override them inline
with `-instructions 'Your skill instructions'`, or convert a canvas you own with
`-canvas F_CANVAS`. These two flags are mutually exclusive. Save `skill.id` from
the output for later commands. Conversion preserves canvas sharing;
Markdown-created skills start invite-only. Creation is not an upsert.

`-mcp-servers-file path.json` accepts an array of MCP app declarations. Slack's
reference does not define the request item schema; provide a known valid payload
for your organization. The example does not assume the response schema is valid
for requests. A file containing `[]` sends an explicit empty array.

## Info — `slackbot.skills.info`

```sh
go run ./examples/slackbot_skills -action info -skill F_SKILL
```

Check `skill.instructions`, effective `access`, nullable fields, and MCP app
details. Built-in skills accept their `builtin:` reference. Even an org admin
needs read access to the underlying canvas.

## Update — `slackbot.skills.update`

```sh
go run ./examples/slackbot_skills -action update -skill F_SKILL \
  -name 'Release review v2' -instructions 'Review the release plan and list three risks.' \
  -description '' -category engineering -example-prompt 'Run {skill}'
```

Only supplied flags change fields. Instructions replace the whole document.
Use `-description ''`, `-category ''`, or `-example-prompt ''` to clear those
values. Use `-mcp-servers-file` to replace the whole MCP list (`[]` clears it).
Run `info` afterward to inspect the result. Failed updates may partially apply;
retrying the same update is safe.

## List — `slackbot.skills.list`

```sh
go run ./examples/slackbot_skills -action list -filter owned -limit 100 -all
go run ./examples/slackbot_skills -action list -filter library -category engineering
go run ./examples/slackbot_skills -action list -filter catalog -cursor 'NEXT_CURSOR'
```

`-all` follows `response_metadata.next_cursor`, even after empty pages. Filters
stay unchanged across requests. The default filter is `library`; instructions
are not returned. Call `info` to read them.

## Library add — `slackbot.skills.library.add`

```sh
go run ./examples/slackbot_skills -action library-add -skill F_SKILL
go run ./examples/slackbot_skills -action list -filter library -all
```

The first command exits successfully when the skill is enabled for the caller,
including when it was already on. The second lets you inspect the effective
library. Sharing, publication, and other people's libraries are unchanged.

## Library remove — `slackbot.skills.library.remove`

```sh
go run ./examples/slackbot_skills -action library-remove -skill F_SKILL
go run ./examples/slackbot_skills -action list -filter library -all
```

Verify the skill is absent from your effective library. This also disables skills
supplied through deployments or skill sets. Repeating removal succeeds.

## Delete — `slackbot.skills.delete`

```sh
go run ./examples/slackbot_skills -action delete -skill F_SKILL
go run ./examples/slackbot_skills -action info -skill F_SKILL
```

Deletion removes the skill **and its canvas**, including catalog, library, and
skill-set references. The subsequent `info` should fail with `skill_not_found`.
Use this command only for a skill you intend to delete.

## Complete group verification

```sh
go run ./examples/slackbot_skills -action workflow
```

This creates a uniquely named temporary skill, reads its instructions, updates
and reads them back, finds it in `owned`, removes it from the personal library,
checks its absence, adds it back, checks its presence, deletes it, and verifies
`skill_not_found`. Lists follow every page and allow up to 15 seconds for index
updates. Cleanup attempts deletion if an intermediate step fails or you press
Ctrl-C. The workflow only deletes the skill it created and prints its ID so you
can recover it if cleanup fails. Success ends with `PASS`.
