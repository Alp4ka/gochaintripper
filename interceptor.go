package gochaintripper

import "net/http"

type Interceptor interface {
	RoundTripWithTransport(
		transport http.RoundTripper,
		req *http.Request,
	) (resp *http.Response, err error)
}

// WrappedInterceptor wraps Interceptor to act like a http.RoundTripper.
type WrappedInterceptor struct {
	transport   http.RoundTripper
	interceptor Interceptor
}

// WrapInterceptor - creates a new WrappedInterceptor with the given transport and interceptor.
func WrapInterceptor(transport http.RoundTripper, interceptor Interceptor) *WrappedInterceptor {
	return &WrappedInterceptor{
		transport:   transport,
		interceptor: interceptor,
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (i *WrappedInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	if i.interceptor == nil {
		return i.transport.RoundTrip(r)
	}

	return i.interceptor.RoundTripWithTransport(i.transport, r)
}

var _ http.RoundTripper = (*WrappedInterceptor)(nil)
