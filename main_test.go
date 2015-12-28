package main

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	proxy := newProxy("", "", "")
	api := newApi(proxy)
	go api.run("8001")
	go proxy.run("8000")
	os.Exit(m.Run())
}

func TestLocalhost(t *testing.T) {
	resp, err := http.Get("http://localhost:8000")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("Received non-403 response: %d\n", resp.StatusCode)
	}
}

func TestSimpleForward(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
	req.Host = "doc.api.lovebystep.com"
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func TestSimpleBlacklist(t *testing.T) {
	http.PostForm("http://localhost:8001/forbid",
		url.Values{"regex": {"doc.api.lovebystep.com"}})
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
	req.Host = "doc.api.lovebystep.com"
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("Received non-403 response: %d\n", resp.StatusCode)
	}
}
