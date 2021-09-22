package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"prometheus-proxier/internal/ascii"

	"github.com/elazarl/goproxy"
)

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !ascii.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

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
