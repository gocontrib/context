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
