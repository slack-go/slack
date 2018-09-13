package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
)

// SecretsVerifier contains the information needed to verify that the request comes from Slack
type SecretsVerifier struct {
	slackSig    string
	timeStamp   string
	requestBody string
}

// NewSecretVerifierFromHeader returns a new SecretsVerifier object in exchange for an http.Header object
func NewSecretVerifierFromHeader(header http.Header) (SecretsVerifier, error) {
	if header["X-Slack-Signature"][0] == "" || header["X-Slack-Request-Timestamp"][0] == "" {
		return SecretsVerifier{}, errors.New("headers are empty, cannot create SecretsVerifier")
	}

	return SecretsVerifier{
		slackSig:  header["X-Slack-Signature"][0],
		timeStamp: header["X-Slack-Request-Timestamp"][0],
	}, nil
}

func (v *SecretsVerifier) Write(body []byte) (n int, err error) {
	v.requestBody = string(body)
	return len(body), nil
}

// Ensure compares the signature sent from Slack with the actual computed hash to judge validity
func (v SecretsVerifier) Ensure(signingSecret string) error {
	message := fmt.Sprintf("v0:%v:%v", v.timeStamp, v.requestBody)

	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(message))

	actualSignature := "v0=" + string(hex.EncodeToString(mac.Sum(nil)))
	fmt.Printf("actual: %s expected: %s", actualSignature, v.slackSig)
	if actualSignature == v.slackSig {
		fmt.Printf("bingo")
		return nil
	}

	return errors.New("invalid token")
}
