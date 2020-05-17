// Package main contains a basic HTTP server using the net/http package.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.String("port", ":8888", "http server port number. Default to :8888.")

func main() {
	log.Println("Running a http server in Go.")
	http.HandleFunc("/showHeaders", showHeaders)
	http.ListenAndServe(*port, nil)
}

func showHeaders(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
