package httprouter

import (
	"fmt"
	"net/http"
	"sync"
)

var ErrRequestNil = fmt.Errorf("*http.Request is nil")

var paramsMap = make(map[*http.Request]map[string]string, 0)
var paramsLock = new(sync.Mutex)

// SetParams associates parameters with a request.
var SetParams = defaultSetParams

// UnsetParams unassociates parameters with a request.
var UnsetParams = defaultUnsetParams

// Params gets the parameters associated with the request. Returns nil if no parameters are associated with the request or the request is nil.
var Params = defaultParams

// ResetParams resets SetParams, UnsetParams and Params to the built-in functions.
func ResetParams() {
	SetParams = defaultSetParams
	UnsetParams = defaultUnsetParams
	Params = defaultParams
}

func defaultSetParams(r *http.Request, params map[string]string) {
	paramsLock.Lock()
	paramsMap[r] = params
	paramsLock.Unlock()
}

func defaultUnsetParams(r *http.Request) {
	paramsLock.Lock()
	delete(paramsMap, r)
	paramsLock.Unlock()
}

func defaultParams(r *http.Request) map[string]string {
	if r == nil {
		return nil
	}
	paramsLock.Lock()
	if params, ok := paramsMap[r]; ok {
		paramsLock.Unlock()
		return params
	}
	paramsLock.Unlock()
	return nil
}

func defaultInitializer(w http.ResponseWriter, r *http.Request, params map[string]string, next http.Handler) {
	SetParams(r, params)
	defer UnsetParams(r)
	next.ServeHTTP(w, r)
}
