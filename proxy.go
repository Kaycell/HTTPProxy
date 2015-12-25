package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type proxy struct {
	requestLogger *log.Logger
}

func newProxy(requestHandle io.Writer) *proxy {
	l := log.New(requestHandle, "REQUEST: ", log.Ldate|log.Ltime)
	return &proxy{requestLogger: l}
}

func (p *proxy) run(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		if host == "" {
			w.WriteHeader(502)
			return
		}
		p.requestLogger.Println(host)
		pr := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   host,
		})
		pr.ServeHTTP(w, req)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
