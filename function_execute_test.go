package slack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func postHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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
			postSuccess(rw, r)
		case "function-failure":
			postFailure(rw, r)
		}
	}
}

func postSuccess(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
    "ok": true
	}`)
	rw.Write(response)
}

func postFailure(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
			"ok": false,
			"error": "function_execution_not_found"
	}`)
	rw.Write(response)
	rw.WriteHeader(500)
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
