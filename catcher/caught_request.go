package catcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var bodyFormatters = map[string]func([]byte) ([]byte, error){
	"application/json": jsonPrettyPrinter,
}

// CaughtRequest represents all the data we collect about a request that
// we catch.
type CaughtRequest struct {
	ID            int64       `json:"-" db:"id"`
	Time          time.Time   `json:"time" db:"when"`
	Host          string      `json:"host" db:"host"`
	Method        string      `json:"method" db:"method"`
	Path          string      `json:"path" db:"path"`
	Headers       http.Header `json:"headers" db:"-"`
	ContentLength int64       `json:"content_length" db:"content_length"`
	RemoteAddr    string      `json:"remote_addr" db:"remote_addr"`
	Form          url.Values  `json:"form_values" db:"-"`
	Body          string      `json:"body" db:"body"`
	RawRequest    string      `json:"raw_request" db:"raw_request"`
}

func convertRequest(req *http.Request) *CaughtRequest {
	raw_request, _ := httputil.DumpRequest(req, true)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
	}

	host := hostWithoutPort(req.Host)

	// Pretty-print the body, if we can.
	prettyBody := string(body)
	if formatter, ok := bodyFormatters[req.Header.Get("Content-Type")]; ok {
		newBody, err := formatter(body)
		if err != nil {
			fmt.Printf("Error formatting body: %v", err)
		} else {
			prettyBody = string(newBody)
		}
	}

	r := &CaughtRequest{
		Time:          time.Now(),
		Host:          host,
		Method:        req.Method,
		Path:          req.RequestURI,
		Headers:       req.Header,
		ContentLength: req.ContentLength,
		RemoteAddr:    hostWithoutPort(req.RemoteAddr),
		Form:          req.PostForm,
		Body:          prettyBody,
		RawRequest:    string(raw_request),
	}
	return r
}

func jsonPrettyPrinter(body []byte) ([]byte, error) {
	var value interface{}

	if err := json.Unmarshal(body, &value); err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(value, "", "  ")
}
