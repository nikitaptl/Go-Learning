package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		Hello, world! <br />
		<img src="/data/img/myImage.png">
	`))
}

func main() {
	http.HandleFunc("/", handler)

	staticHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("./static")))
	http.Handle("/data/", staticHandler)

	addr := ":8080"
	fmt.Println("Starting server at", addr)
	http.ListenAndServe(addr, nil)
}
