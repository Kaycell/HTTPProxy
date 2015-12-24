package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		host := req.Host
		log.Println("Host:", host)
		if host == "" {
			w.WriteHeader(502)
			return
		}

		println("--->", host, os.Args[1], req.URL.String())

		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   host,
		})

		proxy.ServeHTTP(w, req)
	})

	http.ListenAndServe(":"+os.Args[1], nil)
}
