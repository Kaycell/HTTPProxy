package main

import (
	"bufio"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

type proxy struct {
	requestLogger     *log.Logger
	endpointWhiteList []*regexp.Regexp
	endpointBlackList []*regexp.Regexp
	mutex             sync.Mutex
}

// newProxy creates a new instance of proxy.
// It sets request logger using rLogPath as output file or os.Stdout by default.
// If whitePath of blackPath is not empty they are parsed to set endpoint lists.
func newProxy(rLogPath string, whitePath string, blackPath string) *proxy {
	var p proxy
	var l *log.Logger

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
	p.requestLogger = l

	if whitePath != "" {
		err := p.addEndpointListFromFile(whitePath, true)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if blackPath != "" {
		err := p.addEndpointListFromFile(blackPath, false)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &p
}

// addToEndpointListFromFile reads file line by line and calls
// addToEndpointList for each
// use t to choose list type: true for whitelist false for blacklist
func (p *proxy) addEndpointListFromFile(path string, t bool) error {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			err = p.addToEndpointList(s.Text(), t)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// addToEndpointList compiles regex and adds it to an endpointList
// if regex is valid
// use t to choose list type: true for whitelist false for blacklist
func (p *proxy) addToEndpointList(r string, t bool) error {
	rgx, err := regexp.Compile(r)
	if err == nil {
		p.mutex.Lock()
		if t {
			p.endpointWhiteList = append(p.endpointWhiteList, rgx)
		} else {
			p.endpointBlackList = append(p.endpointBlackList, rgx)
		}
		p.mutex.Unlock()
	}
	return err
}

// checkEndpointList looks if r is in whitelist or blackllist
// returns true if endpoint is allowed
func (p *proxy) checkEndpointList(e string) bool {
	if p.endpointBlackList == nil && p.endpointWhiteList == nil {
		return true
	}

	for _, rgx := range p.endpointBlackList {
		if rgx.MatchString(e) {
			return false
		}
	}

	if p.endpointWhiteList == nil {
		return true
	}
	for _, rgx := range p.endpointWhiteList {
		if rgx.MatchString(e) {
			return true
		}
	}
	return false
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if host == "" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if !p.checkEndpointList(host) {
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
}

// run sets a handlefunc which log, authorize and forward requests
// then it listen and serve on specified port
func (p *proxy) run(port string) {
	p.addToEndpointList("localhost", false)
	p.addToEndpointList("127.0.0.1", false)
	log.Fatal(http.ListenAndServe(":"+port, p))
}
