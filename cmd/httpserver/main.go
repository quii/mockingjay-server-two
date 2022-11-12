package main

import (
	"log"
	"net/http"

	"github.com/quii/mockingjay-server-two/adapters/httpserver"
)

func main() {
	app := httpserver.New()

	go func() {
		if err := http.ListenAndServe(":8081", app.AdminRouter); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(":8080", http.HandlerFunc(app.StubHandler)); err != nil {
		log.Fatal(err)
	}
}
