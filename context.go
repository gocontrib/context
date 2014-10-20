package gohttp

import "net/http"
import "github.com/gorilla/context"

// Context middleware to pass specified map or key-value pairs to next handlers
func Context(args ...interface{}) func(http.Handler) http.Handler {
	ctx := mapargs(args...)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			// set request context
			for key := range ctx {
				context.Set(req, key, ctx[key])
			}

			// serve
			h.ServeHTTP(w, req)

			// clear request context
			context.Clear(req)
		})
	}
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
