package main

import (
	"flag"
	"strconv"
)

func main() {
	port := flag.Int("p", 8000, "Proxy listening port")
	rLogPath := flag.String("rl", "", "Request logfile")
	whitePath := flag.String("w", "", "Path to whitelist file")
	blackPath := flag.String("b", "", "Path to blacklist file")

	flag.Parse()

	p := newProxy(*rLogPath, *whitePath, *blackPath)
	p.run(strconv.Itoa(*port))
}
