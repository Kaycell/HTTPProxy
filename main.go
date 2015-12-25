package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func main() {
	port := flag.String("p", "8000", "Proxy listening port")

	flag.Parse()

	if _, err := strconv.Atoi(*port); err != nil {
		log.Fatal("Invalid port: ", *port)
	}

	p := newProxy(os.Stdout)
	p.run(*port)
}
