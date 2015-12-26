package main

import (
	"flag"
	"strconv"
)

func main() {
	port := flag.Int("p", 8000, "Proxy listening port")
	rLogPath := flag.String("rl", "", "Request logfile")
	whitePath := flag.String("w", "", "Path to whitelist file")

	flag.Parse()

	p := newProxy(*rLogPath, *whitePath)
	p.run(strconv.Itoa(*port))
}
