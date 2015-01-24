package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func ParseResponse(path string, values url.Values, intf interface{}, debug bool) error {
	resp, err := http.PostForm(SLACK_API+path, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var decoder *json.Decoder
	if debug {
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Println(string(response))
		decoder = json.NewDecoder(bytes.NewReader(response))
	} else {
		decoder = json.NewDecoder(resp.Body)
	}
	if err := decoder.Decode(&intf); err != nil {
		return err
	}
	return nil
}
