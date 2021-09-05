package ngsi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//Request is an interface to be used when passing header and body information for NGSI-LD API requests
type Request interface {
	BodyReader() io.Reader
	DecodeBodyInto(v interface{}) error
	Request() *http.Request
}

func newRequestWrapper(req *http.Request) Request {
	rw := &requestWrapper{request: req}

	if req.Body != nil {
		// Read the request body and store it in the wrapper
		rw.body, _ = ioutil.ReadAll(req.Body)
		// Restore the body so that others can read it
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rw.body))
	}

	return rw
}

type requestWrapper struct {
	request *http.Request
	body    []byte
}

func (r *requestWrapper) Request() *http.Request {
	return r.request
}

func (r *requestWrapper) BodyReader() io.Reader {
	// Return a new reader with a copy of our stored request body
	return ioutil.NopCloser(bytes.NewBuffer(r.body))
}

func (r *requestWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(r.BodyReader()).Decode(v)
}
