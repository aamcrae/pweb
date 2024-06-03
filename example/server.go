package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 8100, "Server port number")
var base = flag.String("base", "", "Base directory")

func main() {
	flag.Parse()
	fmt.Printf("Server: localhost:%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), http.FileServer(http.Dir(*base))); err != nil {
		log.Fatalf("Server: %v", err)
	}
}
