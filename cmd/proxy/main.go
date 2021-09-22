package main

import (
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

func main() {
	baUser := os.Getenv("BASIC_USER")
	if baUser == "" {
		baUser = "user"
	}
	baPassword := os.Getenv("BASIC_PASSWORD")
	if baPassword == "" {
		baPassword = "user"
	}
	allowedPath := os.Getenv("ALLOWED_PATH")
	if allowedPath == "" {
		allowedPath = "/metrics"
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if r.URL.Path != allowedPath {
			ctx.Logf("invalid path %s", r.URL.Path)
			return r, goproxy.NewResponse(r, "", 500, "")
		}
		r.BasicAuth()
		h := r.Header.Get("Proxy-Authorization")
		user, password, _ := parseBasicAuth(h)
		if user == baUser && password == baPassword {
			return r, nil
		}

		ctx.Logf("invalid proxy auth %s:%s", user, password)
		return r, goproxy.NewResponse(r, "", 500, "")
	})
	log.Fatal(http.ListenAndServe(":2222", proxy))
}
