package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "Port of the server")
	dataDir := flag.String("")
	flag.Parse()
	fmt.Println("Starting the blob server on port " + *port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "Hello, %q from the blob server.", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
