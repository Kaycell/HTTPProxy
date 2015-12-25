package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	port := flag.Int("p", 8000, "Proxy listening port")
	reqLogfile := flag.String("rl", "", "Request logfile")

	flag.Parse()

	var rLogWriter io.Writer
	if *reqLogfile != "" {
		rLogWriter, err := os.OpenFile(*reqLogfile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer rLogWriter.Close()
	} else {
		rLogWriter = os.Stdout
	}

	p := newProxy(rLogWriter)
	p.run(strconv.Itoa(*port))
}
