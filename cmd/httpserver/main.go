package main

import (
	"log"
	"net/http"

	"github.com/quii/mockingjay-server-two/adapters/httpserver"
)

func main() {
	app := new(httpserver.App)

	go func() {
		if err := http.ListenAndServe(":8081", http.HandlerFunc(app.ConfigHandler)); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":8080", http.HandlerFunc(app.StubHandler)); err != nil {
		log.Fatal(err)
	}
}
