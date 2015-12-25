package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
)

type proxy struct {
	requestLogger  *log.Logger
	routeWhiteList []*regexp.Regexp
	routeBlackList []*regexp.Regexp
}

func newProxy(requestHandle io.Writer) *proxy {
	l := log.New(requestHandle, "REQUEST: ", log.Ldate|log.Ltime)
	return &proxy{requestLogger: l}
}

func (p *proxy) run(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if host == "" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		rp := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   host,
		})
		t := time.Now()
		rp.ServeHTTP(w, r)
		p.requestLogger.Println(r.URL.Path, host, time.Since(t))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
