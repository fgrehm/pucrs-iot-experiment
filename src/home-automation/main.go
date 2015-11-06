package main

import (
	"fmt"
	"net/http"
	"os"
)

var hostname string

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from `%s`!</h1>", hostname)
}

func main() {
	hostname, _ = os.Hostname()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
