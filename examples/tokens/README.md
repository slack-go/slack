# Tokens examples

The refresh token endpoint can be used to update
your [configuration tokenset](https://api.slack.com/authentication/config-tokens). These tokens may only be used **once
** before being invalidated, and are only valid for up to **12 hours**.

Once a token has been used, or before it expires, you can use the `RotateTokens()` method to obtain a fresh set to use
for the next request. Depending on your use-case you may want to store these somewhere for a future run, so they are
only returned by the method call. If you wish to update the tokens inside the active Slack client, this can be done
using `UpdateConfigTokens()`.
