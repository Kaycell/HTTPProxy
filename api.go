package main

import (
	"log"
	"net/http"
)

type api struct {
	p *proxy
}

func newApi(p *proxy) *api {
	return &api{p: p}
}

func (api *api) run(port string) {
	http.HandleFunc("/allow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("KOUKOU")
		rgx := r.PostFormValue("regex")
		if rgx == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := api.p.addToEndpointList(rgx, true)
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	})
	http.HandleFunc("/forbid", func(w http.ResponseWriter, r *http.Request) {
		rgx := r.PostFormValue("regex")
		if rgx == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := api.p.addToEndpointList(rgx, false)
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
