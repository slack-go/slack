# Slack Code example

This CLI demonstrates all nine `agents.conversations.*` methods. Each invocation
performs one action in a real workspace and prints the decoded response as JSON.
It also provides `create-canvas` to prepare a canvas for the canvas-view examples
and a read-only `info` action to inspect the channel's decoded metadata.

## Setup

1. Use a workspace where [Slack Code](https://docs.slack.dev/ai/slack-code/)
   is available to your app.
2. Enable code channels in the app manifest:

   ```yaml
   features:
     code_channels:
       enabled: true
   ```

3. Add `code_channels:manage` to the bot token scopes and install or reinstall
   the app. The optional `create-canvas` action also needs `canvases:write`.
   The `info` action needs `channels:read` for public channels or `groups:read`
   for private channels.
4. Set `SLACK_BOT_TOKEN` to the installed app's bot token in your environment.
5. Run the commands below from the repository root. Replace `C_CODE`, `F_CANVAS`,
   and `Ct_VIEW` with IDs returned by previous commands.

The CLI makes Web API calls and exits. Receiving events, view interactions, or
agent slash commands requires a running app, such as the
[Socket Mode handler example](../socketmode_handler/socketmode_handler.go).
Register the appropriate handlers before starting its event loop.

All actions accept `-timeout`, which defaults to `30s`, and stop on Ctrl-C.
Run `go run ./examples/slack_code -help` for the full flag list.
When Slack returns scope diagnostics with an error, the CLI also prints the
required scopes (`needed`) and the token's scopes (`provided`). The Go client
exposes these as `SlackErrorResponse.Needed` and `.Provided` through `errors.As`.

## Create a channel and set properties

```sh
go run ./examples/slack_code -action create -name 'slack-go SDK review'
```

The response contains `channel_id`. Use it as `C_CODE` below. To link creation to
an existing message, supply both `-origin-channel C_SOURCE` and
`-origin-ts 1717171717.123456` using that message's actual channel and timestamp.

```sh
go run ./examples/slack_code -action properties -channel C_CODE \
  -resource-url https://github.com/slack-go/slack
```

This sets an informational context-bar item and the external resource URL.
The context-bar array replaces the current set. Add `-summary-ts` with the
timestamp of an existing summary message to record that message as the summary.
The example uses informational items; typed `code_channel_action` handling is
pending a documented payload schema.

To read channel metadata with the Go client, call `GetConversationInfo` or
`GetConversationInfoContext` and inspect `channel.Properties.AgentSession` and
`channel.Properties.CodeChannel`. These pointers are nil when Slack omits the
corresponding properties. The received code-channel metadata includes context-bar
items with their owning bot user IDs, the summary message, and managed views.

## Verify channel metadata live

Inspect an existing code channel:

```sh
go run ./examples/slack_code -action info -channel C_CODE
```

The output includes the channel's ID, name, and archive state, followed by its
decoded `properties`. Use these to check that you are inspecting the intended
channel. The command reads the channel without changing it.

| JSON path under `properties` | What to verify |
| --- | --- |
| `record_channel.record_type` | `agent_channel` identifies a code channel. |
| `agent_session.status` | The current session status, such as `active`. |
| `agent_session.title` | The session title shown in Slack. |
| `agent_session.encoded_agent_bot_user_ids` | The assigned agents' bot user IDs. The older `agent_bot_user_ids` field may be empty or absent. |
| `agent_session.encoded_initiator_user_id` | The initiator ID, when Slack supplies this observed field. |
| `agent_session.origin_link` | The source channel and message timestamp, when the channel was created with an origin link. |
| `code_channel.context_bar_items` | The items displayed at the top of the channel, including their `bot_user_id`, labels, URLs, and action types. |
| `code_channel.summary_message` | The summary message timestamp and optional parent thread timestamp, after setting a summary. |
| `code_channel.agent_session_views` | The views attached to the channel, with their IDs, types, names, backing file or canvas IDs, versions, owner IDs, and timestamps when present. |

Some fields only appear after the corresponding feature is used. An omitted
origin link or summary is expected when neither has been set; it is not a failed
decode. A newly created code channel can have an empty `code_channel` object.

To verify an origin link, create a test channel from an existing message:

```sh
go run ./examples/slack_code -action create -name 'Metadata verification' \
  -origin-channel C_SOURCE -origin-ts SOURCE_MESSAGE_TS
go run ./examples/slack_code -action info -channel C_NEW_CODE
```

Replace `C_NEW_CODE` with the returned channel ID. The decoded `origin_link`
should match the source channel and timestamp you supplied.

To verify a summary, choose an existing message in the code channel and record
its timestamp:

```sh
go run ./examples/slack_code -action properties -channel C_CODE \
  -summary-ts SUMMARY_MESSAGE_TS
go run ./examples/slack_code -action info -channel C_CODE
```

The `properties` action also replaces the context bar with the sample repository
link. Its received `bot_user_id` should identify the app's bot. If the summary
message is a thread reply, include its parent timestamp:

```sh
go run ./examples/slack_code -action properties -channel C_CODE \
  -summary-ts REPLY_MESSAGE_TS -summary-thread-ts PARENT_MESSAGE_TS
go run ./examples/slack_code -action info -channel C_CODE
```

Check both `message_ts` and `thread_ts` in `code_channel.summary_message`.

After publishing the views and canvas in the following sections, run both:

```sh
go run ./examples/slack_code -action info -channel C_CODE
go run ./examples/slack_code -action list-views -channel C_CODE
```

Compare `properties.code_channel.agent_session_views` from `info` with
`session_views` from `list-views`. Both should expose the attached views, including
the built-in diff view with ID `code`. A canvas's backing ID can appear as
`canvas_id` in channel properties and `file_id` in the view-list response.
`date_updated` is an observed channel-property field. Slack may omit
`content_version` on a canvas's channel-property entry; the decoded integer then
appears as zero.

## Publish and list views

Publish the bundled HTML report, then run the same command again to update it:

```sh
go run ./examples/slack_code -action set-view -channel C_CODE \
  -view-key review/report.html -name 'Review status' \
  -content-file examples/slack_code/report.html
go run ./examples/slack_code -action list-views -channel C_CODE
```

Calls with the same `view_key` update the existing view. Compare its `view_id`
and `content_version` in the responses. Inspect `session_views` for the unified
view list, and also check the legacy `views` array. Live responses can contain
views in `session_views` while `views` is empty. Session entries include `type`
and an optional display `name`; the built-in diff view uses `view_id: "code"`.

To load resources from an additional HTTPS origin in an HTML view, add
`-resource-domain https://cdn.example.com`.

Publish a Block Kit view or the bundled unified diff:

```sh
go run ./examples/slack_code -action set-view -channel C_CODE \
  -type block_kit -view-key review/status -name Status \
  -content-file examples/slack_code/blocks.json
go run ./examples/slack_code -action set-view -channel C_CODE \
  -type diff -content-file examples/slack_code/changes.diff
```

The Block Kit file must be a JSON array of blocks. A diff view is unique within
the channel and does not need a view key. Slack permits up to five views per
channel.

To display a pull request's diff, download the unified diff with the GitHub CLI
and publish it as a diff view:

```sh
gh pr diff 1623 --repo slack-go/slack --color never > /tmp/slack-code-pr.diff
go run ./examples/slack_code -action set-view -channel C_CODE \
  -type diff -content-file /tmp/slack-code-pr.diff
```

Replace the repository and PR number with your own. This replaces the channel's
existing diff view.

Slack's [method reference](https://docs.slack.dev/reference/methods/agents.conversations.setView/)
also documents `type: "pull_request"` with `pr_url`, but a live run returned
`invalid_view_type` for that type. The CLI retains `-type pull_request -pr-url`
for compatibility with the documented API and reports a diff-view workaround
when Slack rejects it. The reference does not explain the availability difference.

## Attach, read, and revise a canvas

Use an existing canvas, or create a standalone one in a workspace that supports
standalone canvases:

```sh
go run ./examples/slack_code -action create-canvas -name 'Review plan' \
  -content-file examples/slack_code/plan.md
```

The response contains `canvas_id`. Use it as `F_CANVAS` below. For `create-canvas`,
add `-channel C_CODE` to create a canvas tab directly with write access. The
[`canvases.create` API](https://docs.slack.dev/reference/methods/canvases.create/)
requires a channel when creating a canvas in a free workspace.

Attach the standalone canvas with comment access:

```sh
go run ./examples/slack_code -action set-view -channel C_CODE \
  -type canvas -view-key review/plan -name Plan -canvas F_CANVAS -access comment
go run ./examples/slack_code -action get-canvas -channel C_CODE -canvas F_CANVAS
```

Add a comment in Slack, then read again. The response includes the markdown,
quoted text, comment threads, replies, and flags indicating truncated results.
Add `-include-resolved` to include resolved threads.

Copy the sample plan, edit the copy, then send the revised document:

```sh
cp examples/slack_code/plan.md /tmp/slack-code-plan.md
```

After editing `/tmp/slack-code-plan.md`:

```sh
go run ./examples/slack_code -action set-canvas-content -channel C_CODE \
  -canvas F_CANVAS -content-file /tmp/slack-code-plan.md
```

[`setCanvasContent`](https://docs.slack.dev/reference/methods/agents.conversations.setCanvasContent/)
returns `sections_changed_count`. Unchanged sections keep their IDs and comment
anchors. Sending the same markdown again returns zero changed sections.

## Register agent commands

The app needs a configured `features.code_channels.slash_command_url`, Socket
Mode, or hosted-app delivery before commands can be registered. See
[agent slash commands](https://docs.slack.dev/ai/slack-code/agent-slash-commands/).

```sh
go run ./examples/slack_code -action commands -channel C_CODE
go run ./examples/slack_code -action commands -channel C_CODE -clear
```

The first command registers `/review-status` and `/summarize`; the second clears
this agent's commands. Registration replaces the current command set. Your
running app must handle invocations of those names; this CLI only registers them.

## Remove views and archive

List the channel's current views, then remove one using its returned `view_id`:

```sh
go run ./examples/slack_code -action list-views -channel C_CODE
go run ./examples/slack_code -action remove-view -channel C_CODE -view-id Ct_VIEW
```

Alternatively, remove a view by its stable key:

```sh
go run ./examples/slack_code -action remove-view -channel C_CODE -view-key review/report.html
```

Remove each view once, using either its ID or its key. Check both `session_views`
and legacy `views` before concluding that the channel has no views. Slack returns
`view_not_found` when no view matches the supplied selector, including when the
view has already been removed. Publish another view with `set-view` to continue
testing removal, then list again to get its current ID.

Removing HTML, diff, or Block Kit views also deletes their backing content files.
Removing a canvas view removes the tab and its view metadata.

When finished, archive the channel:

```sh
go run ./examples/slack_code -action archive -channel C_CODE
```

Add `-summary-ts` with the timestamp of an existing summary message to record it
while archiving. The channel history remains accessible after archival.
