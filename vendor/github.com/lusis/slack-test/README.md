# slacktest

[![Build Status](https://travis-ci.org/lusis/slack-test.svg?branch=master)](https://travis-ci.org/lusis/slack-test)

This is a very basic golang library for testing your slack RTM chatbots

Depending on your mechanism for building slackbots in go and how broken out your message parsing logic is, you can normally test that part pretty cleanly.
However testing your bots RESPONSES are a bit harder.

This library attempts to make that a tad bit easier but in a slightly opinionated way.

The current most popular slack library for golang is [nlopes/slack](https://github.com/nlopes/slack). Conviently the author has made overriding the slack API endpoint a feature. This allows us to use our fake slack server to inspect the chat process.

## Limitations

Right now the test server is VERY limited. It currently handles the following API endpoints

- `rtm.start`
- `chat.postMessage`
- `channels.list`
- `groups.list`
- `users.info`
- `bots.info`

Additional endpoints are welcome.

## Example usage

You can see an example in the `examples` directory of how to you might test it

If you just want to play around:

```go
package main

import (
    "log"
    slacktest "github.com/lusis/slack-test"
    slackbot "github.com/lusis/go-slackbot"
    slack "github.com/nlopes/slack"
)

func globalMessageHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
    bot.Reply(evt, "I see your message", slackbot.WithoutTyping)
}

func main() {
    // Setup our test server
    s := slacktest.NewTestServer()
    // Set a custom name for our bot
    s.SetBotName("MyBotName")
    // ensure that anything using the slack api library uses our custom server
    slack.APIURL = "http://" + s.ServerAddr + "/"
    // start the test server
    go s.Start()
    // create a new slackbot. Token is irrelevant here
    bot := slackbot.New("ABCEDFG")
    // add a message handler
    bot.Hear("this is a channel message").MessageHandler(globalMessageHandler)
    // start the bot
    go bot.Run()
    //send a message to a channel
    s.SendMessageToChannel("#random", "this is a channel message")
    for m := range s.SeenFeed {
        log.Printf("saw message in slack: %s", m)
    }
}
```

Output:

```shell
# go run main.go
2017/09/26 10:53:49 {"type":"message","channel":"#random","text":"this is a channel message","ts":"1506437629","pinned_to":null}
2017/09/26 10:53:49 {"id":1,"channel":"#random","text":"I see your message","type":"message"}
#
```

You can see that our bot handled the message correctly.

## Usage in tests

Currently it's not as ergonomic as I'd like. So much depends on how modular your bot code is in being able to run the same message handling code against a test instance. In the `examples` directory there are a couple of test cases.

Additionally, you definitely need to call `time.Sleep` for now in your test code to give time for the messages to work through the various channels and populate. I'd like to add a safer subscription mechanism in the future with a proper timeout mechanism.

If you want to, you can test the existing example in `examples/go-slackbot`:

```shell
# cd examples/go-slackbot
# go test -v
=== RUN   TestGlobalMessageHandler
--- PASS: TestGlobalMessageHandler (2.00s)
=== RUN   TestHelloMessageHandler
--- PASS: TestHelloMessageHandler (2.00s)
PASS
ok      github.com/lusis/slack-test/examples/go-slackbot        4.005s
#
```

## testing an actual RTM session

This gets tricky. You can look at the existing `rtm_test.go` but here's a documented example:

```go
package foo
import (
    "testing"
    "time"
    "github.com/nlopes/slack"
    "github.com/stretchr/testify/assert"
)

func TestRTMDirectMessage(t *testing.T) {
    // Let's skip this when we want short/quick tests
    if testing.Short() {
        t.Skip("skipping timered test")
    }
    // how long should we wait for this test to finish?
    maxWait := 5 * time.Second
    // start our test server
    s := NewTestServer()
    go s.Start()
    // set our slack API to the mock server
    slack.APIURL = s.GetAPIURL()
    api := slack.New("ABCDEFG")
    // rtm instance
    rtm := api.NewRTM()
    go rtm.ManageConnection()
    // create a channel to pass our results from the next goroutine
    // that is a goroutine doing the normal range over rtm.IncomingEvents
    messageChan := make(chan (*slack.MessageEvent), 1)
    go func() {
        for msg := range rtm.IncomingEvents {
            switch ev := msg.Data.(type) {
            case *slack.MessageEvent:
                messageChan <- ev
            }
        }
    }()
    // since we want to test direct messages, let's send one to the bot
    s.SendDirectMessageToBot(t.Name())
    // now we block this test
    select {
    // if we get a slack.MessageEvent, perform some assertions
    case m := <-messageChan:
        assert.Equal(t, "D024BE91L", m.Channel)
        assert.Equal(t, t.Name(), m.Text)
        break
    // if we hit our timeout, fail the test
    case <-time.After(maxWait):
        assert.FailNow(t, "did not get direct message in time")
    }
}
```