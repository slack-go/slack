package slack

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionCompleteSuccessOutputs(t *testing.T) {
	methods := []struct {
		name string
		call func(*Client, string, ...FunctionCompleteSuccessRequestOption) error
	}{
		{"FunctionCompleteSuccess", (*Client).FunctionCompleteSuccess},
		{"FunctionCompleteSuccessContext", func(api *Client, id string, options ...FunctionCompleteSuccessRequestOption) error {
			return api.FunctionCompleteSuccessContext(t.Context(), id, options...)
		}},
	}
	tests := []struct {
		name    string
		options []FunctionCompleteSuccessRequestOption
		want    string
	}{
		{name: "no outputs", want: `{}`},
		{name: "nil outputs", options: []FunctionCompleteSuccessRequestOption{FunctionCompleteSuccessRequestOptionOutput(nil)}, want: `{}`},
		{name: "empty outputs", options: []FunctionCompleteSuccessRequestOption{FunctionCompleteSuccessRequestOptionOutput(map[string]string{})}, want: `{}`},
		{name: "populated outputs", options: []FunctionCompleteSuccessRequestOption{FunctionCompleteSuccessRequestOptionOutput(map[string]string{"customer_id": "U123"})}, want: `{"customer_id":"U123"}`},
		{name: "custom nil outputs", options: []FunctionCompleteSuccessRequestOption{func(r *FunctionCompleteSuccessRequest) error {
			r.Outputs = nil
			return nil
		}}, want: `{}`},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					var body []byte
					api := New("testing-token", OptionHTTPClient(httpClientFunc(func(r *http.Request) (*http.Response, error) {
						defer r.Body.Close()
						if r.URL.Path != "/api/functions.completeSuccess" {
							t.Errorf("unexpected path: %s", r.URL.Path)
						}
						var err error
						body, err = io.ReadAll(r.Body)
						if err != nil {
							return nil, err
						}
						return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
					})))
					err := method.call(api, "Fx123", test.options...)
					require.NoError(t, err)
					require.JSONEq(t, `{"function_execution_id":"Fx123","outputs":`+test.want+`}`, string(body))
				})
			}
		})
	}
}

func TestFunctionCompleteSuccessOptionError(t *testing.T) {
	wantErr := errors.New("invalid function output")
	// A request would fail on this invalid URL if the option error were ignored.
	api := New("testing-token", OptionAPIURL(":"))
	err := api.FunctionCompleteSuccess("function-success", func(*FunctionCompleteSuccessRequest) error {
		return wantErr
	})
	require.ErrorIs(t, err, wantErr)
}

func postHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	t.Helper()
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Error(err)
			return
		}

		var req FunctionCompleteSuccessRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			t.Error(err)
			return
		}

		switch req.FunctionExecutionID {
		case "function-success":
			postSuccess(rw)
		case "function-failure":
			postFailure(rw)
		}
	}
}

func postSuccess(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
    "ok": true
	}`)
	rw.Write(response)
}

func postFailure(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
			"ok": false,
			"error": "function_execution_not_found"
	}`)
	rw.Write(response)
	rw.WriteHeader(http.StatusInternalServerError)
}

func TestFunctionComplete(t *testing.T) {
	http.HandleFunc("/functions.completeSuccess", postHandler(t))

	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.FunctionCompleteSuccess("function-success")
	if err != nil {
		t.Error(err)
	}

	err = api.FunctionCompleteSuccess("function-failure")
	if err == nil {
		t.Fail()
	}

	err = api.FunctionCompleteSuccessContext(context.Background(), "function-success")
	if err != nil {
		t.Error(err)
	}

	err = api.FunctionCompleteSuccessContext(context.Background(), "function-failure")
	if err == nil {
		t.Fail()
	}
}
