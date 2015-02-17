package httprouter

import (
	"fmt"
	"net/http"
	"sync"
)

var ErrRequestNil = fmt.Errorf("*http.Request is nil")

var paramsMap = make(map[*http.Request]map[string]string, 0)
var paramsLock = new(sync.Mutex)

// SetParams associates parameters with a request. Pass nil to unassociate params with the request.
//
// Passing an empty request will panic.
var SetParams = defaultSetParams

// Params gets the parameters associated with the request. Returns nil if no parameters are associated with the request.
//
// Passing an empty request will panic.
var Params = defaultParams

// ResetParams resets SetParams and Params to the built-in functions.
func ResetParams() {
	SetParams = defaultSetParams
	Params = defaultParams
}

func defaultSetParams(r *http.Request, params map[string]string) {
	if r == nil {
		panic(ErrRequestNil)
	}
	paramsLock.Lock()
	if params != nil {
		paramsMap[r] = params
	} else {
		delete(paramsMap, r)
	}
	paramsLock.Unlock()
}

func defaultParams(r *http.Request) map[string]string {
	if r == nil {
		panic(ErrRequestNil)
	}
	paramsLock.Lock()
	if params, ok := paramsMap[r]; ok {
		paramsLock.Unlock()
		return params
	}
	paramsLock.Unlock()
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request, params map[string]string, next http.Handler) {
	SetParams(r, params)
	defer SetParams(r, nil)
	next.ServeHTTP(w, r)
}
