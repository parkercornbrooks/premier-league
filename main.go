package main

import (
	"log"
	"net/http"
)

const PORT = ":8000"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
