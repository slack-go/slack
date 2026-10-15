# Lists example

This program calls all 12 [Slack Lists methods](https://docs.slack.dev/reference/methods/?family=lists).
Each action runs against a real workspace. It prints the raw response before decoding,
then the relevant typed fields. Use `-debug` for Slack client logs.

## Setup

1. Use a paid Slack workspace with Lists enabled.
2. In the app's **OAuth & Permissions** page, add `lists:read` and `lists:write`
   to the bot token scopes or user token scopes.
3. Install or reinstall the app after changing scopes.
4. Set `SLACK_BOT_TOKEN` to the installed app's bot token (`xoxb-`). Alternatively,
   set `SLACK_USER_TOKEN` to a user token (`xoxp-`). If both are set, the bot token wins.
5. For an existing list, use a token whose acting user or bot can access that list.

These methods do not require an org token, Enterprise Grid, an agent app, event
subscriptions, Socket Mode, or interactivity configuration. The program calls the
Web API directly; no message or command in Slack is needed to trigger it.

Read and export actions need `lists:read`. Write actions need `lists:write`.
The `update` action also reads the item back, so it needs both scopes.
See Slack's [create](https://docs.slack.dev/reference/methods/slackLists.create/),
[list items](https://docs.slack.dev/reference/methods/slackLists.items.list/), and
[export](https://docs.slack.dev/reference/methods/slackLists.download.start/) references.

## Create and edit items

Run commands from the repository root. Create a list with a text column named
Task and a checkbox column named Done:

```sh
go run ./examples/lists -action create -name 'slack-go Lists example'
```

The output gives the list ID and both column IDs. Replace `F_LIST`, `Col_TASK`,
and `Col_DONE` below with those values. Add an item:

```sh
go run ./examples/lists -action add -list F_LIST -column Col_TASK -text 'Review the Lists API'
```

Replace `Rec_ITEM` with the returned item ID. Edit its text, check the checkbox,
then clear it. Each update reads the item back:

```sh
go run ./examples/lists -action update -list F_LIST -item Rec_ITEM -column Col_TASK -text 'Review complete'
go run ./examples/lists -action update -list F_LIST -item Rec_ITEM -column Col_DONE -checked=true
go run ./examples/lists -action update -list F_LIST -item Rec_ITEM -column Col_DONE -checked=false
```

The decoded checkbox should change from `true` to `false`. Compare it with the
raw `slackLists.items.info` response. Open the list in Slack to check the text and
checkbox visually. With a bot token, use `grant` below to give yourself access.

## Read and paginate

```sh
go run ./examples/lists -action info -list F_LIST -item Rec_ITEM
go run ./examples/lists -action list -list F_LIST -limit 2
go run ./examples/lists -action list -list F_LIST -limit 2 -cursor 'RETURNED_CURSOR'
go run ./examples/lists -action list -list F_LIST -limit 2 -all
```

Add at least three items to check pagination with a page size of two. A nonempty
`Next cursor` means another page is available. `-all` follows cursors until Slack
returns an empty cursor. Add `-archived` to read archived items. The `list` action
also prints the list's column IDs, which is useful for an existing list.

## Rename and share

```sh
go run ./examples/lists -action rename -list F_LIST -name 'Reviewed tasks'
go run ./examples/lists -action grant -list F_LIST -user U_USER -access write
go run ./examples/lists -action grant -list F_LIST -channel C_CHANNEL -access read
go run ./examples/lists -action revoke -list F_LIST -channel C_CHANNEL
```

Use either `-user` or `-channel`; each accepts comma-separated IDs. Access levels
are `read`, `write`, and `owner`. `owner` transfers ownership to a user and requires
the current owner's token. See [access rules](https://docs.slack.dev/reference/methods/slackLists.access.set/).

## Export

```sh
go run ./examples/lists -action export -list F_LIST
go run ./examples/lists -action download -list F_LIST -job Le_JOB
```

Replace `Le_JOB` with the job ID from `export`. `download` prints the status and
download URL. If the job is not complete, run `download` again later.

For JSON exports, use the same format and inclusion flags on both calls:

```sh
go run ./examples/lists -action export -list F_LIST -format json -threads -attachments
go run ./examples/lists -action download -list F_LIST -job Le_JOB -format json -threads -attachments
```

Add `-archived` to `export` to include archived items. The example returns the
download URL; it does not fetch the exported file.

## Remove sample items

```sh
go run ./examples/lists -action delete -list F_LIST -item Rec_ITEM
go run ./examples/lists -action delete-multiple -list F_LIST -items Rec_SECOND,Rec_THIRD
```

These actions delete the specified items. Remove the sample list itself in Slack.

## Report results

Send the command, raw response, decoded output, and any `Slack response messages`
for a failed action. For a successful run, compare the IDs, cell values, and
pagination cursor in both outputs. Do not include tokens or private list contents.
