# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.18.0-rc2] - 2026-01-28

### Added

- **Audit Logs example** - New example demonstrating how to use the Audit Logs API. ([#1144])
- **Admin Conversations API support** - Comprehensive support for `admin.conversations.*`
  methods including core operations (archive, unarchive, create, delete, rename, invite,
  search, lookup, getTeams, convertToPrivate, convertToPublic, disconnectShared, setTeams),
  bulk operations (bulkArchive, bulkDelete, bulkMove), preferences, retention management,
  restrict access controls, and EKM channel info. ([#1329])

### Changed

- **BREAKING**: Removed deprecated `UploadFile`, `UploadFileContext`, and
  `FileUploadParameters`. The `files.upload` API was discontinued by Slack on November
  12, 2025. ([#1481])
- **BREAKING**: Renamed `UploadFileV2` → `UploadFile`, `UploadFileV2Context` →
  `UploadFileContext`, and `UploadFileV2Parameters` → `UploadFileParameters`. The "V2"
  suffix is no longer needed now that the old API is removed. ([#1481])

### Fixed

- **File upload error wrapping** - `UploadFile` now wraps errors with the step name
  (`GetUploadURLExternal`, `UploadToURL`, or `CompleteUploadExternal`) so callers can
  identify which of the three upload steps failed. ([#1491])
- **Audit Logs API endpoint** - Fixed `GetAuditLogs` to use the correct endpoint
  (`api.slack.com`) instead of the regular API endpoint (`slack.com/api`). The Audit
  Logs API requires a different base URL. Added `OptionAuditAPIURL` for testing. ([#1144])
- **Socket mode websocket dial debugging** - Added debug logging when a custom dialer is
  used including HTTP response status on dial failures. This helps diagnose proxy/TLS
  issues like "bad handshake" errors. ([#1360])
- **`MsgOptionPostMessageParameters` now passes `MetaData`** - Previously, metadata was
  silently dropped when using `PostMessageParameters`. ([#1343])

## [0.18.0-rc1] - 2026-01-26

### Added

- **Huddle support** - New `HuddleRoom`, `HuddleParticipantEvent`, and `HuddleRecording`
  types for handling Slack huddle events (`huddle_thread` subtype messages).
- **Call block data parsing** - `CallBlock` now includes full call data when retrieved
  from Slack messages, with new `CallBlockData`, `CallBlockDataV1`, and `CallBlockIconURLs`
  types. ([#897])
- **Chat Streaming API support** - New streaming API for real-time chat interactions
  with example usage. ([#1506])
- **Data Access API support** - Full support for Slack's Data Access API with
  example implementation. ([#1439])
- **Cursor-based pagination for `GetUsers`** - More efficient user retrieval
  with cursor pagination. ([#1465])
- **`GetAllConversations` with pagination** - Retrieve all conversations with
  automatic pagination handling, including rate limit and server error handling. ([#1463])
- **Table blocks support** - Parse and create table blocks with proper
  unmarshaling. ([#1490], [#1511])
- **Context actions block support** - New `context_actions` block type. ([#1495])
- **Workflow button block element** - Support for `workflow_button` in block
  elements. ([#1499])
- **`loading_messages` parameter for `SetAssistantThreadsStatus`** - Optional
  parameter to customize loading state messages. ([#1489])
- **Attachment image fields** - Added `ImageBytes`, `ImageHeight`, and `ImageWidth`
  fields to attachments. ([#1516])
- **`RecordChannel` to conversation properties** - New property for conversation
  metadata. ([#1513])
- **Title argument for `CreateChannelCanvas`** - Canvas creation now supports
  custom titles. ([#1483])
- **`PostEphemeral` handler for slacktest** - Audit outgoing ephemeral messages
  in test environments. ([#1517])
- **`PreviewImageName` for remote files** - Customize preview image filename
  instead of using the default `preview.jpg`.

### Fixed

- **`PublishView` no longer sends empty hash** - Prevents unnecessary payload
  when hash is empty. ([#1515])
- **`ImageBlockElement` validation** - Now properly validates that either
  `imageURL` or `SlackFile` is provided. ([#1488])
- **Rich text section channel return** - Correctly returns channel for section
  channel rich text elements. ([#1472])
- **`KickUserFromConversation` error handling** - Errors are now properly parsed
  as a map structure. ([#1471])

### Changed

- **BREAKING**: `GetReactions` now returns `ReactedItem` instead of `[]ItemReaction`.
  This aligns the response with the actual Slack API, which includes the item itself
  (message, file, or file_comment) alongside reactions. To migrate, use `resp.Reactions`
  to access the slice of reactions. ([#1480])
- **BREAKING**: `Settings` struct fields `Interactivity` and `EventSubscriptions`
  are now pointers, allowing them to be omitted when empty. ([#1461])
- Minimum Go version bumped to 1.24. ([#1504])

## [0.17.3] - 2025-07-04

Previous release. See [GitHub releases](https://github.com/slack-go/slack/releases/tag/v0.17.3)
for details.

[#897]: https://github.com/slack-go/slack/issues/897
[#1144]: https://github.com/slack-go/slack/issues/1144
[#1329]: https://github.com/slack-go/slack/issues/1329
[#1343]: https://github.com/slack-go/slack/issues/1343
[#1360]: https://github.com/slack-go/slack/issues/1360
[#1439]: https://github.com/slack-go/slack/pull/1439
[#1461]: https://github.com/slack-go/slack/pull/1461
[#1463]: https://github.com/slack-go/slack/pull/1463
[#1465]: https://github.com/slack-go/slack/pull/1465
[#1471]: https://github.com/slack-go/slack/pull/1471
[#1472]: https://github.com/slack-go/slack/pull/1472
[#1480]: https://github.com/slack-go/slack/pull/1480
[#1483]: https://github.com/slack-go/slack/pull/1483
[#1488]: https://github.com/slack-go/slack/pull/1488
[#1489]: https://github.com/slack-go/slack/pull/1489
[#1490]: https://github.com/slack-go/slack/pull/1490
[#1491]: https://github.com/slack-go/slack/issues/1491
[#1495]: https://github.com/slack-go/slack/pull/1495
[#1499]: https://github.com/slack-go/slack/pull/1499
[#1504]: https://github.com/slack-go/slack/pull/1504
[#1506]: https://github.com/slack-go/slack/pull/1506
[#1511]: https://github.com/slack-go/slack/pull/1511
[#1513]: https://github.com/slack-go/slack/pull/1513
[#1515]: https://github.com/slack-go/slack/pull/1515
[#1516]: https://github.com/slack-go/slack/pull/1516
[#1517]: https://github.com/slack-go/slack/pull/1517

[Unreleased]: https://github.com/slack-go/slack/compare/v0.18.0-rc2...HEAD
[0.18.0-rc2]: https://github.com/slack-go/slack/releases/tag/v0.18.0-rc2
[0.18.0-rc1]: https://github.com/slack-go/slack/releases/tag/v0.18.0-rc1
[0.17.3]: https://github.com/slack-go/slack/releases/tag/v0.17.3
