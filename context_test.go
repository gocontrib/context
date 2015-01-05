package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/go-supertest"
	. "github.com/franela/goblin"
	"github.com/gohttp/app"
	. "github.com/onsi/gomega"
)

func TestContext(t *testing.T) {

	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	describe := g.Describe
	it := g.It

	describe("mapargs", func() {
		it("should work with map", func() {
			m1 := make(map[string]interface{})
			m1["a"] = 1
			m := mapargs(m1)
			Expect(m["a"]).Should(Equal(1))
		})

		it("should work with pairs", func() {
			m := mapargs(1, 2)
			Expect(m[1]).Should(Equal(2))
		})
	})

	describe("using Context middleware", func() {
		it("should be easy as pie", func(done Done) {

			const key = "key"
			const value = "value"

			a := app.New()
			a.Use(New(key, value))
			a.Get("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				v, ok := Get(req, key).(string)
				if ok {
					w.Write([]byte(v))
				} else {
					w.Write([]byte("error"))
				}
			}))

			server := httptest.NewServer(a)
			defer server.Close()

			NewRequest(server.URL).
				Get("/").
				Expect(200, value, done)
		})
	})
}
