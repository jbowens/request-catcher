package catcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// CaughtRequest represents all the data we collect about a request that
// catch.
type CaughtRequest struct {
	Time          time.Time   `json:"time"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Headers       http.Header `json:"headers"`
	ContentLength int64       `json:"content_length"`
	RemoteAddr    string      `json:"remote_addr"`
	Form          url.Values  `json:"form_values"`
	Body          string      `json:"body"`
	RawRequest    string      `json:"raw_request"`
}

func convertRequest(req *http.Request) *CaughtRequest {

	raw_request, _ := httputil.DumpRequest(req, true)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
	}

	r := &CaughtRequest{
		Time:          time.Now(),
		Method:        req.Method,
		Path:          req.RequestURI,
		Headers:       req.Header,
		ContentLength: req.ContentLength,
		RemoteAddr:    req.RemoteAddr,
		Form:          req.PostForm,
		Body:          string(body),
		RawRequest:    string(raw_request),
	}
	return r
}
