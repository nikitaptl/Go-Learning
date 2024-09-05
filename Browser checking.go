package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("RequestID", "fdsf3423qrfs")

	fmt.Fprintln(w, "Your browser is ", r.UserAgent())
	fmt.Fprintln(w, "You accept", r.Header.Get("Accept"))
}

func main() {
	http.HandleFunc("/", handler)

	addr := ":8080"
	fmt.Println("Starting server at", addr)
	http.ListenAndServe(addr, nil)
}
