package context

import "net/http"
import "github.com/gorilla/context"

// New creates context middleware to pass specified map or key-value pairs to next handlers
func New(args ...interface{}) func(http.Handler) http.Handler {
	ctx := mapargs(args...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set request context
			for key := range ctx {
				context.Set(r, key, ctx[key])
			}

			// serve
			next.ServeHTTP(w, r)

			// clear request context
			context.Clear(r)
		})
	}
}

// Set stores a value for a given key in a given request.
func Set(r *http.Request, key, val interface{}) {
	context.Set(r, key, val)
}

// Get returns a value stored for a given key in a given request.
func Get(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

const keyRequestID = "request-id"

// GetRequestID returns id of request.
func GetRequestID(r *http.Request) string {
	var v = Get(r, keyRequestID)
	var s, ok = v.(string)
	if ok {
		return s
	}
	return ""
}

// SetRequestID assigns request id.
func SetRequestID(r *http.Request, id string) {
	Set(r, keyRequestID, id)
}

func mapargs(args ...interface{}) map[interface{}]interface{} {
	var (
		i   = 0
		n   = len(args)
		res = make(map[interface{}]interface{})
	)

	for i < n {

		// do not use reflection since it is slow
		m, ok := args[i].(map[string]interface{})
		if ok {
			for k := range m {
				res[k] = m[k]
			}
			i++
			continue
		}

		key := args[i]
		i++
		res[key] = args[i]
		i++
	}

	return res
}
