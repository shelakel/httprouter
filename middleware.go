package httprouter

import "net/http"

type Middleware []func(http.Handler) http.Handler

// NewMiddleware creates a new middleware chain.
func NewMiddleware(middleware ...func(http.Handler) http.Handler) Middleware {
	if middleware == nil || len(middleware) == 0 {
		return Middleware{}
	}
	return append(make(Middleware, len(middleware)), middleware...)
}

// Chains the middleware to the http.Handler endpoint.
func (mw Middleware) Then(h http.Handler) http.Handler {
	var final http.Handler
	if h != nil {
		final = h
	} else {
		final = http.DefaultServeMux
	}
	for i := len(mw) - 1; i >= 0; i-- {
		final = mw[i](final)
	}
	return final
}

// Chains the middleware to the http.HandlerFunc endpoint.
func (mw Middleware) ThenFunc(h func(http.ResponseWriter, *http.Request)) http.Handler {
	return mw.Then(http.HandlerFunc(h))
}

// Add a func(http.Handler)http.Handler middleware function.
func (mw *Middleware) Use(middleware ...func(http.Handler) http.Handler) {
	*mw = append(*mw, middleware...)
}
