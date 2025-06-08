# Events Example

This is a very simple example but should give you a glimpse of how to use the events API.

## How to enable this

1. Disable socket mode in the app if it's enabled: this will reveal the `Request URL`.
2. Set up the events `Request URL` in a way that matches the endpoint in
   [events.go](./events.go).

   You can find this [here](https://api.slack.com/apps/<appid>/event-subscriptions).
3. Set up the events you want to be subscribed to.
4. Copy the bot token and signing secret and set up the environment variables (per code).
5. Run the example:
   ```bash
   go run events.go
   ```
