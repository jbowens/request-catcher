package catcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var bodyFormatters = map[string]func([]byte) ([]byte, error){
	"application/json": jsonPrettyPrinter,
}

type RequestData struct {
	Time          time.Time   `json:"time"`
	Host          string      `json:"host"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Headers       http.Header `json:"headers"`
	ContentLength int64       `json:"content_length"`
	RemoteAddr    string      `json:"remote_addr"`
	Form          url.Values  `json:"form_values"`
	Body          string      `json:"body"`
	RawRequest    string      `json:"raw_request"`
}

type CaughtRequest struct {
	*http.Request
	time time.Time
	raw  []byte
}

func convertRequest(req *http.Request) *CaughtRequest {
	// Dumping a request will replace req.Body with an in-memory
	// cached version of the request body. So it's okay to read
	// from req.Body even after the client connection is gone.
	raw, _ := httputil.DumpRequest(req, true)
	return &CaughtRequest{time: time.Now(), raw: raw, Request: req}
}

func (req *CaughtRequest) MarshalJSON() ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	host := hostWithoutPort(req.Host)

	// Pretty-print the body, if we can.
	var prettyBody string
	if formatter, ok := bodyFormatters[req.Header.Get("Content-Type")]; ok {
		if newBody, err := formatter(body); err == nil {
			prettyBody = string(newBody)
		}
	}
	if prettyBody == "" {
		prettyBody = string(body)
	}

	return json.Marshal(&RequestData{
		Time:          req.time,
		Host:          host,
		Method:        req.Method,
		Path:          req.RequestURI,
		Headers:       req.Header,
		ContentLength: req.ContentLength,
		RemoteAddr:    hostWithoutPort(req.RemoteAddr),
		Form:          req.PostForm,
		Body:          prettyBody,
		RawRequest:    string(req.raw),
	})
}

func jsonPrettyPrinter(body []byte) ([]byte, error) {
	var value interface{}

	if err := json.Unmarshal(body, &value); err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(value, "", "  ")
}

func hostWithoutPort(host string) string {
	if sepIndex := strings.IndexRune(host, ':'); sepIndex != -1 {
		host = host[:sepIndex]
	}
	return host
}
