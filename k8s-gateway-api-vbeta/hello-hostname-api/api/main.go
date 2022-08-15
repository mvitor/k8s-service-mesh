package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", HandleGet)

	http.ListenAndServe(":8080", nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "GoLang Hello from Pod: %s", hostname)
}
