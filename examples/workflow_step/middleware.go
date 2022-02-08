package main

import (
	"bytes"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
)

func (v *SecretsVerifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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
