# Events Example

The following are the settings you'll need to configure in the Slack UI to get this example to work. In particular, you must disable socket mode to reveal the `reqeust_url` text field in the Slack UI.

```
_metadata:
  major_version: 1
  minor_version: 1
display_information:
  name: Name
  description: Description
  background_color: "<hexcode>"
features:
  app_home:
    home_tab_enabled: N/A
    messages_tab_enabled: N/A
    messages_tab_read_only_enabled: N/A
  bot_user:
    display_name: Name
    always_online: N/A
oauth_config:
  scopes:
    bot:
      - app_mentions:read *
      - channels:history *
settings:
  allowed_ip_address_ranges:  # <--- Required
    - xxx.xxx.xxx.xxx/32
    - xxx.xxx.xxx.xxx/32
    - xxx.xxx.xxx.xxx/32
  event_subscriptions:
    request_url: https://example.com/example-endpoint
    bot_events:
      - app_mention
      - message.channels
  org_deploy_enabled: N/A
  socket_mode_enabled: false  # <--- Required

```
