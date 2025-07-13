package interceptors

import (
	"net/http"
)

// ChainTripper joins passed interceptors and implements http.RoundTripper interface.
type ChainTripper struct {
	wrappedChain []*WrappedInterceptor
}

func NewInterceptors(transport http.RoundTripper, chain ...Interceptor) *ChainTripper {
	wrappedChain := make([]*WrappedInterceptor, len(chain)+1)
	wrappedChain[len(chain)] = WrapInterceptor(transport, nil)

	for i := len(chain) - 1; i >= 0; i-- {
		wrappedChain[i] = WrapInterceptor(wrappedChain[i+1], chain[i])
	}

	return &ChainTripper{
		wrappedChain: wrappedChain,
	}
}

func (i *ChainTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return i.wrappedChain[0].RoundTrip(req)
}

var _ http.RoundTripper = (*ChainTripper)(nil)
