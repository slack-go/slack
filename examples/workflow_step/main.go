package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

type (
	appContext struct {
		slack  *slack.Client
		config configuration
	}
	configuration struct {
		botToken      string
		signingSecret string
	}
	SecretsVerifierMiddleware struct {
		handler http.Handler
	}
)

const (
	APIBaseURL = "/api/v1"
	// MyExampleWorkflowStepCallbackID is configured in slack (api.slack.com/apps).
	// Select your app or create a new one. Then choose menu "Workflow Steps"...
	MyExampleWorkflowStepCallbackID = "example-step"
)

var appCtx appContext

func main() {
	appCtx.config.botToken = os.Getenv("SLACK_BOT_TOKEN")
	appCtx.config.signingSecret = os.Getenv("SLACK_SIGNING_SECRET")

	appCtx.slack = slack.New(appCtx.config.botToken)

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s/interaction", APIBaseURL), handleInteraction)
	mux.HandleFunc(fmt.Sprintf("%s/%s", APIBaseURL, MyExampleWorkflowStepCallbackID), handleMyWorkflowStep)
	middleware := NewSecretsVerifierMiddleware(mux)

	log.Printf("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware))
}
