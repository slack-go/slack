# Verify interactive element fields

This example posts one test message and listens for interactions through Socket Mode. It tests eight fields on buttons, checkboxes, conversation multi-selects, and rich text inputs. It leaves the message in the channel for comparison between users.

## Configure the app

- Export `SLACK_BOT_TOKEN` (`xoxb-`, bot scope `chat:write`) and `SLACK_APP_TOKEN` (`xapp-`, app-level scope `connections:write`). The repository's `.envrc` can supply them.
- Enable Socket Mode and Interactivity & Shortcuts. No request URL or event subscriptions are needed.
- Install the app, then invite its bot to the test channel. Reinstall after changing scopes.
- Stop other Socket Mode processes for the same app during this test, so they cannot receive its interactions.

These fields do not require an org-level user token or Enterprise Grid. The Slackbot prompt path requires Slackbot AI access for the person who clicks the button. Any app can send `agent_prompt`; it does not require the agent app manifest feature. Without that user access, Slack sends the normal interaction instead. See the [button reference](https://docs.slack.dev/reference/block-kit/block-elements/button-element/).

## Run

From the repository root, use an existing test channel and your own user ID:

```sh
go run ./examples/block_element_fields -channel C0123456789 -user U0123456789
```

The default `-case all` posts all message elements together. Keep the process running while you interact with them. Add `-debug` for library diagnostics. To isolate a rejection, rerun with `-case button`, `-case agent`, `-case checkboxes`, or `-case modal`. Each run posts one message; `-user` is only required for `all`, `button`, and `checkboxes`.

After **Submit**, the modal closes and the program prints the selected conversation and entered text in the terminal.

## Check each field

| Field | Check in Slack |
|---|---|
| Button `accessibility_label` | As the `-user` user, focus **Visible test button** with a screen reader. It should announce **Verify the button accessibility label**. |
| Button `visible_to_user_ids` | The `-user` user should see **Visible test button**. A different member of the same channel should not see it. Click it and confirm a `field_check_button` action appears in the output. |
| Button `agent_prompt` | Click **Test Slackbot prompt**. With Slackbot AI access, Slackbot should open with the test prompt. Without access, the output should show the `field_check_agent` fallback action. The fallback alone does not test the AI path. |
| Button `agent_prompt_display` | On the resulting Slackbot message, check for **Verify Block Kit prompt fields**. Expand the message to inspect the full prompt. |
| Checkboxes `visible_to_user_ids` | Only the `-user` user should see **Visible test checkbox**. Toggle it and confirm the `field_check_checkboxes` action. Compare with a second channel member. |
| Multi-select `default_to_current_conversation` | While viewing the test channel, click **Open input tests**. The modal should already select that channel. Leave it unchanged, submit, and check `includes launch channel=true`. |
| Rich text input `min_lines` | Open the modal and check that its empty editor shows two visible lines. |
| Rich text input `max_lines` | Enter at least six lines. The editor should grow to four visible lines, then scroll. Submit and check that the output still contains all entered text. These fields set visible height, not a limit on submitted lines. |

The multi-select behavior is documented in the [conversation multi-select reference](https://docs.slack.dev/reference/block-kit/block-elements/multi-select-menu-element/). The editor height behavior is documented in the [rich text input reference](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-input-element/).

## Read the output

The program prints outgoing blocks, raw `chat.postMessage` and `views.open` response bodies, and Socket Mode envelopes before the library parses them. Tokens, response URLs, trigger IDs, and credential-bearing WebSocket URLs are redacted. No HTTP request headers are logged. Payloads still contain workspace IDs and content entered in the modal.

For each field Slack returns, `Go preserved=true` confirms that decoding and re-encoding retained the returned value. A decoding error or `Go preserved=false` stops the example with a nonzero exit status. If Slack omits a field, the example reports it as unverified. Values that Slack changes remain visible for comparison with the outgoing blocks; the example does not treat them as a Go decoding failure. Request acceptance alone does not prove retention or correct UI behavior. Complete the table above.

Paste back the per-field output, any `response_metadata.messages` errors, and the UI observations. Include whether the clicking user had Slackbot AI access. Press Ctrl-C when finished.
