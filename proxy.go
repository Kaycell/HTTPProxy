package main

import (
	"bufio"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

type proxy struct {
	requestLogger  *log.Logger
	routeWhiteList []*regexp.Regexp
}

// newProxy create a new instance of proxy.
// It set request logger using rLogPath as output file or os.Stdout by default.
// It also check whitelist file regexs validity
func newProxy(rLogPath string, whitePath string) *proxy {
	var l *log.Logger
	var rwl []*regexp.Regexp

	if rLogPath != "" {
		f, err := os.OpenFile(rLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		l = log.New(f, "REQUEST: ", log.Ldate|log.Ltime)
	} else {
		l = log.New(os.Stdout, "REQUEST: ", log.Ldate|log.Ltime)
	}

	if whitePath != "" {
		fWhite, err := os.Open(whitePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer fWhite.Close()
		s := bufio.NewScanner(fWhite)
		for s.Scan() {
			rgx, err := regexp.Compile(s.Text())
			if err != nil {
				log.Fatalln(err)
			}
			rwl = append(rwl, rgx)
		}
	}

	return &proxy{requestLogger: l, routeWhiteList: rwl}
}

// checkWhiteList return true if r match with one regexp in whitelist file
// or if any whitelist was set
func (p *proxy) checkWhiteList(r string) bool {
	if p.routeWhiteList == nil {
		log.Println("NO WHITELIST DEFINED")
		return true
	}
	for _, rgx := range p.routeWhiteList {
		if rgx.MatchString(r) {
			return true
		}
	}
	return false
}

func (p *proxy) run(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if host == "" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if !p.checkWhiteList(host) {
			p.requestLogger.Println(host, "FORBIDDEN")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		rp := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   host,
		})
		t := time.Now()
		rp.ServeHTTP(w, r)
		p.requestLogger.Println(host, time.Since(t))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
	}
}
