package main

import (
	"flag"
	"strconv"
)

func main() {
	port := flag.Int("p", 8000, "Proxy listening port")
	apiPort := flag.Int("api-port", 8001, "Proxy admin api listening port")
	rLogPath := flag.String("rl", "", "Request logfile")
	whitePath := flag.String("w", "", "Path to whitelist file")
	blackPath := flag.String("b", "", "Path to blacklist file")

	flag.Parse()

	proxy := newProxy(*rLogPath, *whitePath, *blackPath)
	api := newApi(proxy)
	go proxy.run(strconv.Itoa(*port))
	api.run(strconv.Itoa(*apiPort))
}
