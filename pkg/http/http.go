package http

import (
	"net"
	"net/http"
	"net/url"
)

// Request represents the HTTP request received by APISIX.
// We don't use net/http's Request because it doesn't suit our purpose.
// Take `Request.Header` as an example:
//
// 1. We need to record any change to the request headers. As the Request.Header
// is not an interface, there is not way to inject our special tracker.
//
// 2. As the author of fasthttp pointed out, "headers are stored in a map[string][]string.
// So the server must parse all the headers, ...". The official API is suboptimal, which
// is even worse in our case as it is not a real HTTP server.
type Request interface {
	// ID returns the request id
	ID() uint32

	// SrcIP returns the client's IP
	SrcIP() net.IP
	// Method returns the HTTP method (GET, POST, PUT, etc.)
	Method() string
	// Path returns the path part of the client's URI (without query string and the other parts)
	// It won't be equal to the one in the Request-Line sent by the client if it has
	// been rewritten by APISIX
	Path() []byte
	// SetPath is the setter for Path
	SetPath([]byte)
	// Header returns the HTTP headers
	Header() Header
	// Args returns the query string
	Args() url.Values

	// Var returns the value of a Nginx variable, like `r.Var("request_time")`
	//
	// To fetch the value, the runner will look up the request's cache first. If not found,
	// the runner will ask it from the APISIX. If the RPC call is failed, an error in
	// pkg/common.ErrConnClosed type is returned.
	Var(name string) ([]byte, error)
}

// Header is like http.Header, but only implements the subset of its methods
type Header interface {
	// Set sets the header entries associated with key to the single element value.
	// It replaces any existing values associated with key.
	// The key is case insensitive
	Set(key, value string)
	// Del deletes the values associated with key. The key is case insensitive
	Del(key string)
	// Get gets the first value associated with the given key.
	// If there are no values associated with the key, Get returns "".
	// It is case insensitive
	Get(key string) string
	// View returns the internal structure. It is expected for read operations. Any write operation
	// won't be recorded
	View() http.Header

	// TODO: support Add
}
