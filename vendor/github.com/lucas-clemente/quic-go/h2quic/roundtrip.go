package h2quic

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/lex/httplex"
)

// RoundTripper implements the http.RoundTripper interface
type RoundTripper struct {
	mutex sync.Mutex

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool

	// TLSClientConfig specifies the TLS configuration to use with
	// tls.Client. If nil, the default configuration is used.
	TLSClientConfig *tls.Config

	clients map[string]http.RoundTripper
}

var _ http.RoundTripper = &RoundTripper{}

// RoundTrip does a round trip
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		closeRequestBody(req)
		return nil, errors.New("quic: nil Request.URL")
	}
	if req.URL.Host == "" {
		closeRequestBody(req)
		return nil, errors.New("quic: no Host in request URL")
	}
	if req.Header == nil {
		closeRequestBody(req)
		return nil, errors.New("quic: nil Request.Header")
	}

	if req.URL.Scheme == "https" {
		for k, vv := range req.Header {
			if !httplex.ValidHeaderFieldName(k) {
				return nil, fmt.Errorf("quic: invalid http header field name %q", k)
			}
			for _, v := range vv {
				if !httplex.ValidHeaderFieldValue(v) {
					return nil, fmt.Errorf("quic: invalid http header field value %q for key %v", v, k)
				}
			}
		}
	} else {
		closeRequestBody(req)
		return nil, fmt.Errorf("quic: unsupported protocol scheme: %s", req.URL.Scheme)
	}

	if req.Method != "" && !validMethod(req.Method) {
		closeRequestBody(req)
		return nil, fmt.Errorf("quic: invalid method %q", req.Method)
	}

	hostname := authorityAddr("https", hostnameFromRequest(req))
	return r.getClient(hostname).RoundTrip(req)
}

func (r *RoundTripper) getClient(hostname string) http.RoundTripper {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.clients == nil {
		r.clients = make(map[string]http.RoundTripper)
	}

	client, ok := r.clients[hostname]
	if !ok {
		client = newClient(r.TLSClientConfig, hostname, &roundTripperOpts{DisableCompression: r.DisableCompression})
		r.clients[hostname] = client
	}
	return client
}

func closeRequestBody(req *http.Request) {
	if req.Body != nil {
		req.Body.Close()
	}
}

func validMethod(method string) bool {
	/*
				     Method         = "OPTIONS"                ; Section 9.2
		   		                    | "GET"                    ; Section 9.3
		   		                    | "HEAD"                   ; Section 9.4
		   		                    | "POST"                   ; Section 9.5
		   		                    | "PUT"                    ; Section 9.6
		   		                    | "DELETE"                 ; Section 9.7
		   		                    | "TRACE"                  ; Section 9.8
		   		                    | "CONNECT"                ; Section 9.9
		   		                    | extension-method
		   		   extension-method = token
		   		     token          = 1*<any CHAR except CTLs or separators>
	*/
	return len(method) > 0 && strings.IndexFunc(method, isNotToken) == -1
}

// copied from net/http/http.go
func isNotToken(r rune) bool {
	return !httplex.IsTokenRune(r)
}
