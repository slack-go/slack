package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/slack-go/slack"
)

func (v *SecretsVerifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	sv, err := slack.NewSecretsVerifier(r.Header, appCtx.config.signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	v.handler.ServeHTTP(w, r)
}

func NewSecretsVerifierMiddleware(h http.Handler) *SecretsVerifierMiddleware {
	return &SecretsVerifierMiddleware{h}
}
