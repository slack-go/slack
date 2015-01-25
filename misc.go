package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func fileUploadReq(path, fpath string, values url.Values) (*http.Request, error) {
	fullpath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	wr := multipart.NewWriter(body)

	ioWriter, err := wr.CreateFormFile("file", filepath.Base(fullpath))
	if err != nil {
		wr.Close()
		return nil, err
	}
	bytes, err := io.Copy(ioWriter, file)
	if err != nil {
		wr.Close()
		return nil, err
	}
	// Close the multipart writer or the footer won't be written
	wr.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if bytes != stat.Size() {
		return nil, errors.New("could not read the whole file")
	}
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", wr.FormDataContentType())
	req.URL.RawQuery = (values).Encode()
	return req, nil
}

func parseResponseBody(body io.ReadCloser, intf *interface{}, debug bool) error {
	var decoder *json.Decoder
	if debug {
		response, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		log.Println(string(response))
		decoder = json.NewDecoder(bytes.NewReader(response))
	} else {
		decoder = json.NewDecoder(body)
	}
	if err := decoder.Decode(&intf); err != nil {
		return err
	}
	return nil

}

func parseResponseMultipart(path string, filepath string, values url.Values, intf interface{}, debug bool) error {
	req, err := fileUploadReq(SLACK_API+path, filepath, values)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return parseResponseBody(resp.Body, &intf, debug)
}

func parseResponse(path string, values url.Values, intf interface{}, debug bool) error {
	resp, err := http.PostForm(SLACK_API+path, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return parseResponseBody(resp.Body, &intf, debug)
}
