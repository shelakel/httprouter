package httprouter

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

var ErrContextNotSet = fmt.Errorf("Context not set.")

var contextsLock = new(sync.Mutex)
var contexts map[*http.Request]context.Context = make(map[*http.Request]context.Context, 0)

// Get the Context associated with the request.
//
// Panics if there's no Context associated with the request.
func Context(r *http.Request) context.Context {
	contextsLock.Lock()
	if ctx, ok := contexts[r]; ok {
		contextsLock.Unlock()
		return ctx
	}
	contextsLock.Unlock()
	panic(ErrContextNotSet)
}

// Set the Context associated with the request.
func SetContext(r *http.Request, ctx context.Context) {
	contextsLock.Lock()
	contexts[r] = ctx
	contextsLock.Unlock()
}

// Removes the Context from the request.
func unsetContext(r *http.Request) {
	contextsLock.Lock()
	delete(contexts, r)
	contextsLock.Unlock()
}

// Get the "params" map[string]string from
// the Context associated with the request.
//
// Panics if there's no Context associated with the request.
func Params(r *http.Request) map[string]string {
	ctx := Context(r)
	ps := ctx.Value("params")
	if ps == nil {
		ps = map[string]string{}
		SetContext(r, context.WithValue(ctx, "params", ps))
	}
	ps1 := ps.(map[string]string)
	return ps1
}
