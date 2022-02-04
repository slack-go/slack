#WorkflowStep

Have you ever wanted to run an app from a Slack workflow? This sample app shows you how it works.

Slack describes some of the basics here:
https://api.slack.com/workflows/steps
https://api.slack.com/tutorials/workflow-builder-steps


1. Start the example app localy on port 8080


2. Use ngrok to expose your app to the internet

```shell
    ./ngrok http 8080
```
Copy the https forwarding URL and paste it into the app manifest down below (event_subscription request_url and interactivity request_url)


3. Create a new Slack App at api.slack.com/apps from an app manifest

The manifest of a sample Slack App looks like this:
```yaml
display_information:
  name: Workflowstep-Example
features:
  bot_user:
    display_name: Workflowstep-Example
    always_online: false
  workflow_steps:
    - name: Example Step
      callback_id: example-step
oauth_config:
  scopes:
    bot:
      - workflow.steps:execute
settings:
  event_subscriptions:
    request_url: https://*****.ngrok.io/api/v1/example-step
    bot_events:
      - workflow_step_execute
  interactivity:
    is_enabled: true
    request_url: https://*****.ngrok.io/api/v1/interaction
  org_deploy_enabled: false
  socket_mode_enabled: false
  token_rotation_enabled: false
```

("Interactivity" and "Enable Events" should be turned on)

4. Slack Workflow (**paid plan required!**)
   1. Create a new Workflow at app.slack.com/workflow-builder
   2. give it a name
   3. select "Planned date & time"
   4. add another step and select "Example Step" from App Workflowstep-Example
   5. configure your app and hit save
   6. don't forget to publish your workflow