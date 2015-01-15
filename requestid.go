package context

// Borrowed from goji framework (see https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go)
// and adapted for gohttp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

const keyRequestID = "request-id"

var (
	reqidPrefix  string
	reqidCounter uint64
)

/*
  A quick note on the statistics here: we're trying to calculate the chance that
  two randomly generated base62 prefixes will collide. We use the formula from
  http://en.wikipedia.org/wiki/Birthday_problem
  P[m, n] \approx 1 - e^{-m^2/2n}
  We ballpark an upper bound for $m$ by imagining (for whatever reason) a server
  that restarts every second over 10 years, for $m = 86400 * 365 * 10 = 315360000$
  For a $k$ character base-62 identifier, we have $n(k) = 62^k$
  Plugging this in, we find $P[m, n(10)] \approx 5.75%$, which is good enough for
  our purposes, and is surely more than anyone would ever need in practice -- a
  process that is rebooted a handful of times a day for a hundred years has less
  than a millionth of a percent chance of generating two colliding IDs.
*/

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	reqidPrefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		myid := atomic.AddUint64(&reqidCounter, 1)
		SetRequestID(r, fmt.Sprintf("%s-%06d", reqidPrefix, myid))

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

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
	r.Header.Set("X-Request-Id", id)

	// store id in request context also
	Set(r, keyRequestID, id)
}
