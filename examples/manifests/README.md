# Manifest examples

This example shows how to interact with the
new [manifest endpoints](https://api.slack.com/reference/manifests#manifest_apis). These endpoints require a special set
of tokens called `configuration tokens`. Refer to
the [relevant documentation](https://api.slack.com/authentication/config-tokens) for how to create these tokens.

For examples on how to use configuration tokens, see the [tokens example](../tokens).

## Usage info

The manifest endpoints allow you to configure your application programmatically instead of manually creating
a `manifest.yaml` file and uploading it on your Slack application's dashboard.

A manifest should follow a specific structure and has a handful of required fields. These are describe in
the [manifest documentation](https://api.slack.com/reference/manifests#fields), but Slack additionally returns very
informative error messages for malformed templates to help you pin down what the issue is. The library itself does not
attempt to perform any form of validation on your manifest.

**Note that each configuration token may only be used once before being invalidated. Again refer to the tokens example
for more information.**

## Available methods

- ``Slack.CreateManifest()``
- ``Slack.DeleteManifest()``
- ``Slack.ExportManifest()``
- ``Slack.UpdateManifest()``

## Example details

The example code here only shows how to _update_ an application using a manifest. The other available methods are either
identical in usage or trivial to use, so no full example is provided for them.

The example doesn't rotate the configuration tokens after updating the manifest. **You should almost always do this**.
Your access token is invalidated after sending a request, and rotating your tokens will allow you to make another
request in the future. This example does not do this explicitly as it would just repeat the tokens example. For sake of
simplicity, it only focuses on the manifest part.
