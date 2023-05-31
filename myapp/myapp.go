package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("myapp"))
	})
	log.Println("http server start...")
	http.ListenAndServe(":80", nil)
}
