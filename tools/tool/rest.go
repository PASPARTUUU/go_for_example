package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/PASPARTUUU/go_for_example/tools/errpath"
)

// SendRequest -
func SendRequest(method string, url string, headers map[string]string, payload interface{}, response interface{}) error {
	var err error
	var b []byte

	if t := reflect.TypeOf(payload); t != nil {
		if t.Kind() == reflect.String {
			if payload == "" {
				b = nil
			} else {
				b = []byte(strings.Replace(strings.Replace(strings.Replace(fmt.Sprintf("%v", payload), " ", "", -1), "\n", "", -1), "\t", "", -1))
			}
		} else {
			b, err = json.Marshal(payload)
			if err != nil {
				return errpath.Err(err, "failed to marshal a payload")
			}
		}
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return errpath.Err(err, "failed to create an http request")
	}

	// req.Header.Add("Authorization","Bearer"+token)
	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errpath.Err(err, fmt.Sprintf("failed to make a %s request", method))
	}
	defer resp.Body.Close()

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errpath.Errorf("%s request to %s failed with status: %d", method, url, resp.StatusCode)
		}
		return errpath.Errorf("%s request to %s failed with status: %d and body: %s", method, url, resp.StatusCode, string(body))
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return errpath.Err(err, "decoding error")
		}
	}

	return nil
}

// SetURLParams -
func SetURLParams(url string, params map[string]string) string {

	for _, char := range url {
		if string(char) == "?" {
			fmt.Println("warn symbol '?' exist")
		}
	}
	url += "?"

	var i int = 0
	for key, val := range params {
		p := key + "=" + val
		if i < len(params)-1 {
			p += "&"
		}
		url += p
		i++
	}

	return url
}
