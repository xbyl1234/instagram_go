package main

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("http: %s\n", r.RemoteAddr)
	fmt.Printf("get: %s\n", r.URL.RequestURI())
	fmt.Printf("header:\n")
	for k, v := range r.Header {
		fmt.Printf("---\t%s\t\t:\t%s\n", k, v)
	}
	fmt.Printf("\n\n")
	w.Write([]byte("123"))
}

func main() {
	http.HandleFunc("/", HelloHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
